package apis

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/types"
)

// bindMetricsApi registers the metrics dashboard api endpoints.
func bindMetricsApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	sub := rg.Group("/metrics").Bind(RequireSuperuserAuth())
	sub.GET("/overview", metricsOverview)
	sub.GET("/requests", metricsRequests)
	sub.GET("/latency", metricsLatency)
	sub.GET("/errors", metricsErrors)
	sub.GET("/endpoints", metricsEndpoints)
	sub.GET("/collections", metricsCollections)
}

// OverviewResponse represents the overview metrics response.
type OverviewResponse struct {
	TotalRequests int64   `json:"totalRequests"`
	AvgLatency    float64 `json:"avgLatency"`
	ErrorRate     float64 `json:"errorRate"`
	DatabaseSize  int64   `json:"databaseSize"`
	TotalErrors   int64   `json:"totalErrors"`
	Period        string  `json:"period"`
}

// RequestsTimeSeriesItem represents a single point in the requests time-series.
type RequestsTimeSeriesItem struct {
	Date  types.DateTime `json:"date" db:"date"`
	Total int            `json:"total" db:"total"`
}

// LatencyTimeSeriesItem represents a single point in the latency time-series.
type LatencyTimeSeriesItem struct {
	Date types.DateTime `json:"date" db:"date"`
	Avg  float64        `json:"avg" db:"avg"`
	P50  float64        `json:"p50" db:"p50"`
	P95  float64        `json:"p95" db:"p95"`
	P99  float64        `json:"p99" db:"p99"`
}

// ErrorsTimeSeriesItem represents error counts by status code.
type ErrorsTimeSeriesItem struct {
	Date        types.DateTime `json:"date" db:"date"`
	Status4xx   int            `json:"status4xx" db:"status4xx"`
	Status5xx   int            `json:"status5xx" db:"status5xx"`
	TotalErrors int            `json:"totalErrors" db:"total_errors"`
}

// EndpointStats represents statistics for a single endpoint.
type EndpointStats struct {
	Endpoint   string  `json:"endpoint" db:"endpoint"`
	Count      int     `json:"count" db:"count"`
	AvgLatency float64 `json:"avgLatency" db:"avg_latency"`
}

// CollectionStats represents record count for a collection.
type CollectionStats struct {
	Name        string `json:"name"`
	RecordCount int64  `json:"recordCount"`
	Type        string `json:"type"`
}

// metricsOverview returns overview statistics.
// GET /api/metrics/overview?period=24h
func metricsOverview(e *core.RequestEvent) error {
	period := e.Request.URL.Query().Get("period")
	if period == "" {
		period = "24h"
	}

	duration, err := parsePeriod(period)
	if err != nil {
		return e.BadRequestError("Invalid period format. Use: 1h, 6h, 24h, 7d", err)
	}

	since := time.Now().UTC().Add(-duration)
	sinceStr := since.Format(types.DefaultDateLayout)

	// Get total requests count
	var totalRequests int64
	err = e.App.AuxDB().NewQuery(`
		SELECT COUNT(*) 
		FROM {{_logs}} 
		WHERE [[created]] >= {:since}
		AND json_extract([[data]], '$.type') = 'request'
	`).Bind(dbx.Params{"since": sinceStr}).Row(&totalRequests)
	if err != nil {
		totalRequests = 0
	}

	// Get average latency
	var avgLatency float64
	err = e.App.AuxDB().NewQuery(`
		SELECT COALESCE(AVG(json_extract([[data]], '$.execTime')), 0)
		FROM {{_logs}} 
		WHERE [[created]] >= {:since}
		AND json_extract([[data]], '$.type') = 'request'
		AND json_extract([[data]], '$.execTime') IS NOT NULL
	`).Bind(dbx.Params{"since": sinceStr}).Row(&avgLatency)
	if err != nil {
		avgLatency = 0
	}

	// Get error count (4xx and 5xx status codes)
	var totalErrors int64
	err = e.App.AuxDB().NewQuery(`
		SELECT COUNT(*)
		FROM {{_logs}} 
		WHERE [[created]] >= {:since}
		AND json_extract([[data]], '$.type') = 'request'
		AND json_extract([[data]], '$.status') >= 400
	`).Bind(dbx.Params{"since": sinceStr}).Row(&totalErrors)
	if err != nil {
		totalErrors = 0
	}

	// Calculate error rate
	var errorRate float64
	if totalRequests > 0 {
		errorRate = float64(totalErrors) / float64(totalRequests) * 100
	}

	// Get database size using PRAGMA
	var pageCount, pageSize int64
	err = e.App.DB().NewQuery("PRAGMA page_count").Row(&pageCount)
	if err != nil {
		pageCount = 0
	}
	err = e.App.DB().NewQuery("PRAGMA page_size").Row(&pageSize)
	if err != nil {
		pageSize = 0
	}
	databaseSize := pageCount * pageSize

	response := OverviewResponse{
		TotalRequests: totalRequests,
		AvgLatency:    avgLatency,
		ErrorRate:     errorRate,
		DatabaseSize:  databaseSize,
		TotalErrors:   totalErrors,
		Period:        period,
	}

	return e.JSON(http.StatusOK, response)
}

// metricsRequests returns time-series data of requests.
// GET /api/metrics/requests?period=24h&interval=1h
func metricsRequests(e *core.RequestEvent) error {
	period := e.Request.URL.Query().Get("period")
	if period == "" {
		period = "24h"
	}

	interval := e.Request.URL.Query().Get("interval")
	if interval == "" {
		interval = "1h"
	}

	duration, err := parsePeriod(period)
	if err != nil {
		return e.BadRequestError("Invalid period format", err)
	}

	since := time.Now().UTC().Add(-duration)
	sinceStr := since.Format(types.DefaultDateLayout)

	// Determine grouping format based on interval
	groupFormat := getGroupFormat(interval)

	var result []RequestsTimeSeriesItem
	err = e.App.AuxDB().NewQuery(`
		SELECT 
			strftime('` + groupFormat + `', [[created]]) as date,
			COUNT(*) as total
		FROM {{_logs}}
		WHERE [[created]] >= {:since}
		AND json_extract([[data]], '$.type') = 'request'
		GROUP BY date
		ORDER BY date ASC
	`).Bind(dbx.Params{"since": sinceStr}).All(&result)

	if err != nil {
		return e.BadRequestError("Failed to fetch requests data", err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"data":     result,
		"period":   period,
		"interval": interval,
	})
}

// metricsLatency returns latency percentiles over time.
// GET /api/metrics/latency?period=24h&interval=1h
func metricsLatency(e *core.RequestEvent) error {
	period := e.Request.URL.Query().Get("period")
	if period == "" {
		period = "24h"
	}

	interval := e.Request.URL.Query().Get("interval")
	if interval == "" {
		interval = "1h"
	}

	duration, err := parsePeriod(period)
	if err != nil {
		return e.BadRequestError("Invalid period format", err)
	}

	since := time.Now().UTC().Add(-duration)
	sinceStr := since.Format(types.DefaultDateLayout)

	groupFormat := getGroupFormat(interval)

	// First, get all latency values grouped by time period
	type rawLatency struct {
		Date     string  `db:"date"`
		ExecTime float64 `db:"exec_time"`
	}
	var rawData []rawLatency
	err = e.App.AuxDB().NewQuery(`
		SELECT 
			strftime('` + groupFormat + `', [[created]]) as date,
			json_extract([[data]], '$.execTime') as exec_time
		FROM {{_logs}}
		WHERE [[created]] >= {:since}
		AND json_extract([[data]], '$.type') = 'request'
		AND json_extract([[data]], '$.execTime') IS NOT NULL
		ORDER BY date ASC
	`).Bind(dbx.Params{"since": sinceStr}).All(&rawData)

	if err != nil {
		return e.BadRequestError("Failed to fetch latency data", err)
	}

	// Group by date and calculate percentiles
	grouped := make(map[string][]float64)
	for _, item := range rawData {
		grouped[item.Date] = append(grouped[item.Date], item.ExecTime)
	}

	var result []LatencyTimeSeriesItem
	for date, latencies := range grouped {
		sort.Float64s(latencies)
		n := len(latencies)
		if n == 0 {
			continue
		}

		// Calculate average
		var sum float64
		for _, v := range latencies {
			sum += v
		}
		avg := sum / float64(n)

		// Calculate percentiles
		p50 := percentile(latencies, 0.50)
		p95 := percentile(latencies, 0.95)
		p99 := percentile(latencies, 0.99)

		parsedDate, _ := types.ParseDateTime(date)
		result = append(result, LatencyTimeSeriesItem{
			Date: parsedDate,
			Avg:  avg,
			P50:  p50,
			P95:  p95,
			P99:  p99,
		})
	}

	// Sort by date
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.Time().Before(result[j].Date.Time())
	})

	return e.JSON(http.StatusOK, map[string]any{
		"data":     result,
		"period":   period,
		"interval": interval,
	})
}

// metricsErrors returns error counts by status code over time.
// GET /api/metrics/errors?period=24h&interval=1h
func metricsErrors(e *core.RequestEvent) error {
	period := e.Request.URL.Query().Get("period")
	if period == "" {
		period = "24h"
	}

	interval := e.Request.URL.Query().Get("interval")
	if interval == "" {
		interval = "1h"
	}

	duration, err := parsePeriod(period)
	if err != nil {
		return e.BadRequestError("Invalid period format", err)
	}

	since := time.Now().UTC().Add(-duration)
	sinceStr := since.Format(types.DefaultDateLayout)

	groupFormat := getGroupFormat(interval)

	var result []ErrorsTimeSeriesItem
	err = e.App.AuxDB().NewQuery(`
		SELECT 
			strftime('` + groupFormat + `', [[created]]) as date,
			SUM(CASE WHEN json_extract([[data]], '$.status') >= 400 AND json_extract([[data]], '$.status') < 500 THEN 1 ELSE 0 END) as status4xx,
			SUM(CASE WHEN json_extract([[data]], '$.status') >= 500 THEN 1 ELSE 0 END) as status5xx,
			SUM(CASE WHEN json_extract([[data]], '$.status') >= 400 THEN 1 ELSE 0 END) as total_errors
		FROM {{_logs}}
		WHERE [[created]] >= {:since}
		AND json_extract([[data]], '$.type') = 'request'
		GROUP BY date
		ORDER BY date ASC
	`).Bind(dbx.Params{"since": sinceStr}).All(&result)

	if err != nil {
		return e.BadRequestError("Failed to fetch errors data", err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"data":     result,
		"period":   period,
		"interval": interval,
	})
}

// metricsEndpoints returns top endpoints by request count.
// GET /api/metrics/endpoints?period=24h&limit=10
func metricsEndpoints(e *core.RequestEvent) error {
	period := e.Request.URL.Query().Get("period")
	if period == "" {
		period = "24h"
	}

	limitStr := e.Request.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	duration, err := parsePeriod(period)
	if err != nil {
		return e.BadRequestError("Invalid period format", err)
	}

	since := time.Now().UTC().Add(-duration)
	sinceStr := since.Format(types.DefaultDateLayout)

	var result []EndpointStats
	err = e.App.AuxDB().NewQuery(`
		SELECT 
			json_extract([[data]], '$.url') as endpoint,
			COUNT(*) as count,
			COALESCE(AVG(json_extract([[data]], '$.execTime')), 0) as avg_latency
		FROM {{_logs}}
		WHERE [[created]] >= {:since}
		AND json_extract([[data]], '$.type') = 'request'
		AND json_extract([[data]], '$.url') IS NOT NULL
		GROUP BY endpoint
		ORDER BY count DESC
		LIMIT {:limit}
	`).Bind(dbx.Params{"since": sinceStr, "limit": limit}).All(&result)

	if err != nil {
		return e.BadRequestError("Failed to fetch endpoints data", err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"data":   result,
		"period": period,
		"limit":  limit,
	})
}

// metricsCollections returns record counts for all collections.
// GET /api/metrics/collections
func metricsCollections(e *core.RequestEvent) error {
	collections, err := e.App.FindAllCollections()
	if err != nil {
		return e.BadRequestError("Failed to fetch collections", err)
	}

	var result []CollectionStats
	for _, col := range collections {
		// Skip system collections
		if col.System {
			continue
		}

		count, err := e.App.CountRecords(col.Id)
		if err != nil {
			count = 0
		}

		result = append(result, CollectionStats{
			Name:        col.Name,
			RecordCount: count,
			Type:        col.Type,
		})
	}

	// Sort by record count descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].RecordCount > result[j].RecordCount
	})

	return e.JSON(http.StatusOK, map[string]any{
		"data":  result,
		"total": len(result),
	})
}

// parsePeriod parses a period string like "1h", "24h", "7d" into a time.Duration.
func parsePeriod(period string) (time.Duration, error) {
	if len(period) < 2 {
		return 0, strconv.ErrSyntax
	}

	unit := period[len(period)-1]
	valueStr := period[:len(period)-1]

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	}

	switch unit {
	case 'h':
		return time.Duration(value) * time.Hour, nil
	case 'd':
		return time.Duration(value) * 24 * time.Hour, nil
	case 'm':
		return time.Duration(value) * time.Minute, nil
	default:
		return 0, strconv.ErrSyntax
	}
}

// getGroupFormat returns the SQLite strftime format based on the interval.
func getGroupFormat(interval string) string {
	switch interval {
	case "1m", "5m", "15m":
		return "%Y-%m-%d %H:%M:00"
	case "1h":
		return "%Y-%m-%d %H:00:00"
	case "1d":
		return "%Y-%m-%d 00:00:00"
	default:
		return "%Y-%m-%d %H:00:00"
	}
}

// percentile calculates the percentile value from a sorted slice.
func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	if len(sorted) == 1 {
		return sorted[0]
	}
	
	index := p * float64(len(sorted)-1)
	lower := int(index)
	upper := lower + 1
	
	if upper >= len(sorted) {
		return sorted[len(sorted)-1]
	}
	
	weight := index - float64(lower)
	return sorted[lower]*(1-weight) + sorted[upper]*weight
}
