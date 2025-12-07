# Dashboard Feature

## Overview

The Dashboard is a real-time metrics and monitoring feature for PocketBase that provides visual insights into API performance, error rates, and collection statistics. It displays data from PocketBase's built-in `_logs` table through interactive charts and summary cards.

## Features

### Core Features (V3)
- **Overview Metrics** - Total requests, average latency, error rate, and database size
- **Requests Chart** - Time-series line chart showing request volume over time
- **Latency Percentiles** - Multi-line chart showing p50/p95/p99 latency percentiles
- **Top Endpoints** - Horizontal bar chart of most-accessed endpoints
- **Collection Statistics** - Table showing record counts per collection
- **Auto-Refresh** - Automatic data refresh every 30 seconds (configurable)
- **Time Period Selection** - View metrics for 1 hour, 6 hours, 24 hours, or 7 days

## Access

### Admin UI

1. Log in to the PocketBase Admin dashboard as a superuser
2. Click **Dashboard** in the main sidebar (chart icon)
3. Or navigate directly to `http://127.0.0.1:8090/_/#/dashboard`

### Requirements

- **Authentication**: Superuser access required
- **Data Source**: Uses data from the `_logs` table (no additional setup required)

## API Reference

All dashboard endpoints require superuser authentication.

### Overview Metrics

```
GET /api/metrics/overview?period=24h
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `period` | string | `24h` | Time period: `1h`, `6h`, `24h`, or `7d` |

**Response:**
```json
{
  "totalRequests": 15420,
  "avgLatency": 45.23,
  "errorRate": 2.5,
  "databaseSize": 52428800,
  "totalErrors": 385,
  "period": "24h"
}
```

**Fields:**
| Field | Type | Description |
|-------|------|-------------|
| `totalRequests` | integer | Total number of requests in the period |
| `avgLatency` | float | Average response time in milliseconds |
| `errorRate` | float | Percentage of requests with 4xx/5xx status |
| `databaseSize` | integer | SQLite database size in bytes |
| `totalErrors` | integer | Count of error responses |
| `period` | string | The requested time period |

---

### Requests Time-Series

```
GET /api/metrics/requests?period=24h&interval=1h
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `period` | string | `24h` | Time period |
| `interval` | string | `1h` | Grouping interval: `1m`, `5m`, `15m`, `1h`, `1d` |

**Response:**
```json
{
  "data": [
    { "date": "2025-12-07 10:00:00.000Z", "total": 142 },
    { "date": "2025-12-07 11:00:00.000Z", "total": 238 }
  ],
  "period": "24h",
  "interval": "1h"
}
```

---

### Latency Percentiles

```
GET /api/metrics/latency?period=24h&interval=1h
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `period` | string | `24h` | Time period |
| `interval` | string | `1h` | Grouping interval |

**Response:**
```json
{
  "data": [
    {
      "date": "2025-12-07 10:00:00.000Z",
      "avg": 45.5,
      "p50": 32.0,
      "p95": 120.5,
      "p99": 250.0
    }
  ],
  "period": "24h",
  "interval": "1h"
}
```

**Fields:**
| Field | Type | Description |
|-------|------|-------------|
| `avg` | float | Average latency in milliseconds |
| `p50` | float | 50th percentile (median) latency |
| `p95` | float | 95th percentile latency |
| `p99` | float | 99th percentile latency |

---

### Error Statistics

```
GET /api/metrics/errors?period=24h&interval=1h
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `period` | string | `24h` | Time period |
| `interval` | string | `1h` | Grouping interval |

**Response:**
```json
{
  "data": [
    {
      "date": "2025-12-07 10:00:00.000Z",
      "status4xx": 15,
      "status5xx": 2,
      "totalErrors": 17
    }
  ],
  "period": "24h",
  "interval": "1h"
}
```

**Fields:**
| Field | Type | Description |
|-------|------|-------------|
| `status4xx` | integer | Count of 4xx client errors |
| `status5xx` | integer | Count of 5xx server errors |
| `totalErrors` | integer | Total error count |

---

### Top Endpoints

```
GET /api/metrics/endpoints?period=24h&limit=10
```

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `period` | string | `24h` | Time period |
| `limit` | integer | `10` | Maximum number of endpoints to return |

**Response:**
```json
{
  "data": [
    {
      "endpoint": "/api/collections/users/records",
      "count": 1542,
      "avgLatency": 38.5
    },
    {
      "endpoint": "/api/collections/posts/records",
      "count": 892,
      "avgLatency": 45.2
    }
  ],
  "period": "24h",
  "limit": 10
}
```

---

### Collection Statistics

```
GET /api/metrics/collections
```

**Response:**
```json
{
  "data": [
    {
      "name": "users",
      "recordCount": 1542,
      "type": "auth"
    },
    {
      "name": "posts",
      "recordCount": 8923,
      "type": "base"
    }
  ],
  "total": 5
}
```

**Note:** System collections are excluded from the response.

---

## UI Components

### MetricCard
Displays a single metric value with icon, title, and subtitle.

### RequestsChart
Line chart showing request volume over time using Chart.js.

### LatencyChart
Multi-line chart showing latency percentiles (avg, p50, p95, p99).

### EndpointsChart
Horizontal bar chart showing top endpoints by request count.

### CollectionsTable
Table displaying collection names, types, and record counts.

## Example Usage

### JavaScript (using fetch)

```javascript
const token = 'YOUR_AUTH_TOKEN';

// Fetch overview metrics
const response = await fetch('http://127.0.0.1:8090/api/metrics/overview?period=24h', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const data = await response.json();
console.log(`Total Requests: ${data.totalRequests}`);
console.log(`Avg Latency: ${data.avgLatency.toFixed(2)}ms`);
console.log(`Error Rate: ${data.errorRate.toFixed(2)}%`);
console.log(`Database Size: ${(data.databaseSize / 1024 / 1024).toFixed(2)}MB`);
```

### cURL

```bash
# Get overview metrics
curl -X GET "http://127.0.0.1:8090/api/metrics/overview?period=24h" \
  -H "Authorization: Bearer YOUR_AUTH_TOKEN"

# Get requests time-series
curl -X GET "http://127.0.0.1:8090/api/metrics/requests?period=7d&interval=1d" \
  -H "Authorization: Bearer YOUR_AUTH_TOKEN"

# Get top endpoints
curl -X GET "http://127.0.0.1:8090/api/metrics/endpoints?period=24h&limit=5" \
  -H "Authorization: Bearer YOUR_AUTH_TOKEN"
```

## Architecture

### Data Source

The Dashboard reads metrics from PocketBase's built-in `_logs` table, which automatically stores request logs. No additional database tables or configuration are required.

**Key log fields used:**
- `created` - Timestamp for time-series grouping
- `data.type` - Filter for 'request' type logs
- `data.status` - HTTP status code for error rate calculation
- `data.execTime` - Execution time for latency metrics
- `data.url` - Request URL for endpoint statistics

### Backend Components

```
apis/metrics.go           # 6 API endpoints
apis/metrics_test.go      # 17 test cases
apis/base.go              # Route registration
```

### Frontend Components

```
ui/src/pages/Dashboard.svelte                      # Main page
ui/src/stores/dashboard.js                         # State management
ui/src/components/dashboard/MetricCard.svelte      # Metric card
ui/src/components/dashboard/RequestsChart.svelte   # Requests chart
ui/src/components/dashboard/LatencyChart.svelte    # Latency chart
ui/src/components/dashboard/EndpointsChart.svelte  # Endpoints chart
ui/src/components/dashboard/CollectionsTable.svelte # Collections table
ui/src/scss/_dashboard.scss                        # Styles
```

## Performance Considerations

1. **Time Period Selection** - Use shorter periods (1h, 6h) for faster queries on systems with large log tables
2. **Auto-Refresh** - Default 30-second refresh can be adjusted or disabled for lower server load
3. **Log Retention** - Consider implementing log retention policies for systems with high request volumes
4. **Index Usage** - The `_logs` table is automatically indexed by PocketBase

## Troubleshooting

### Dashboard Not Loading

**Problem**: Dashboard shows loading state indefinitely.

**Solutions**:
1. Ensure you're logged in as a superuser
2. Check browser console for API errors
3. Verify the PocketBase server is running
4. Check network connectivity

### No Data Displayed

**Problem**: Charts and metrics show zero values.

**Possible Causes**:
- No requests logged in the selected time period
- Logs table is empty (new installation)
- Selected period is too short

**Solutions**:
1. Select a longer time period (e.g., 7 days)
2. Make some API requests to generate log data
3. Check that logging is enabled in PocketBase

### Slow Dashboard Performance

**Problem**: Dashboard takes long to load.

**Solutions**:
1. Use shorter time periods (1h instead of 7d)
2. Increase refresh interval or disable auto-refresh
3. Consider implementing log retention to reduce table size

## Security

- All endpoints require superuser authentication
- No sensitive data is exposed through metrics
- Log queries use parameterized SQL to prevent injection

## Related Documentation

- [AI Query Feature](./AI_QUERY_FEATURE.md)
- [SQL Terminal Feature](./SQL_TERMINAL_FEATURE.md)
- [Data Import Feature](./IMPORT_FEATURE.md)
- [PocketBase Logs](https://pocketbase.io/docs/logs)
