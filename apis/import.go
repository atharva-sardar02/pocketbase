package apis

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

// bindImportApi registers the data import api endpoints.
func bindImportApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	sub := rg.Group("/import").Bind(RequireSuperuserAuth())
	sub.POST("/preview", importPreview)
	sub.POST("/validate", importValidate)
	sub.POST("/execute", importExecute)
}

// PreviewRequest represents the request body for preview endpoint.
type PreviewRequest struct {
	Data      string `json:"data"`      // Base64 encoded file content or raw text
	Format    string `json:"format"`    // "csv" or "json"
	Delimiter string `json:"delimiter"` // For CSV: "," or "\t" (default: ",")
}

// PreviewResponse represents the preview endpoint response.
type PreviewResponse struct {
	Headers    []string          `json:"headers"`
	SampleRows [][]string        `json:"sampleRows"`
	TotalRows  int               `json:"totalRows"`
	Format     string            `json:"format"`
	Errors     []string          `json:"errors,omitempty"`
}

// ValidateRequest represents the request body for validation endpoint.
type ValidateRequest struct {
	Collection string            `json:"collection"`
	Mapping    map[string]string `json:"mapping"` // source column -> target field
}

// ValidateResponse represents the validation endpoint response.
type ValidateResponse struct {
	Valid          bool              `json:"valid"`
	Errors         []string          `json:"errors,omitempty"`
	FieldTypes     map[string]string `json:"fieldTypes,omitempty"`
	RequiredFields []string          `json:"requiredFields,omitempty"`
}

// ExecuteRequest represents the request body for import execution.
type ExecuteRequest struct {
	Collection string            `json:"collection"`
	Data       string            `json:"data"`      // Base64 encoded file content or raw text
	Format     string            `json:"format"`    // "csv" or "json"
	Delimiter  string            `json:"delimiter"` // For CSV
	Mapping    map[string]string `json:"mapping"`   // source column -> target field
	SkipHeader bool              `json:"skipHeader"` // Skip first row for CSV (default: true)
}

// ExecuteResponse represents the import execution response.
type ExecuteResponse struct {
	TotalRows    int            `json:"totalRows"`
	SuccessCount int            `json:"successCount"`
	FailureCount int            `json:"failureCount"`
	Errors       []ImportError  `json:"errors,omitempty"`
}

// ImportError represents an error for a specific row during import.
type ImportError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// importPreview parses uploaded data and returns headers + sample rows.
// POST /api/import/preview
func importPreview(e *core.RequestEvent) error {
	var req PreviewRequest
	if err := e.BindBody(&req); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}

	if req.Data == "" {
		return e.BadRequestError("Data is required", nil)
	}

	if req.Format == "" {
		req.Format = detectFormat(req.Data)
	}

	var response PreviewResponse
	response.Format = req.Format

	switch strings.ToLower(req.Format) {
	case "csv":
		headers, rows, total, err := parseCSV(req.Data, req.Delimiter, 5)
		if err != nil {
			response.Errors = append(response.Errors, err.Error())
		} else {
			response.Headers = headers
			response.SampleRows = rows
			response.TotalRows = total
		}
	case "json":
		headers, rows, total, err := parseJSON(req.Data, 5)
		if err != nil {
			response.Errors = append(response.Errors, err.Error())
		} else {
			response.Headers = headers
			response.SampleRows = rows
			response.TotalRows = total
		}
	default:
		return e.BadRequestError("Unsupported format. Use 'csv' or 'json'", nil)
	}

	return e.JSON(http.StatusOK, response)
}

// importValidate validates the field mapping against the collection schema.
// POST /api/import/validate
func importValidate(e *core.RequestEvent) error {
	var req ValidateRequest
	if err := e.BindBody(&req); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}

	if req.Collection == "" {
		return e.BadRequestError("Collection is required", nil)
	}

	// Find the collection
	collection, err := e.App.FindCollectionByNameOrId(req.Collection)
	if err != nil {
		return e.BadRequestError("Collection not found: "+req.Collection, err)
	}

	response := ValidateResponse{
		Valid:      true,
		FieldTypes: make(map[string]string),
	}

	// Get collection fields
	fieldMap := make(map[string]bool)
	for _, field := range collection.Fields {
		fieldMap[field.GetName()] = true
		response.FieldTypes[field.GetName()] = field.Type()
	}

	// Validate mapping
	if len(req.Mapping) > 0 {
		for sourceCol, targetField := range req.Mapping {
			if targetField == "" || targetField == "-" {
				continue // Skip unmapped columns
			}

			if !fieldMap[targetField] {
				response.Valid = false
				response.Errors = append(response.Errors, 
					fmt.Sprintf("Field '%s' (mapped from '%s') does not exist in collection", targetField, sourceCol))
			}
		}
	}

	return e.JSON(http.StatusOK, response)
}

// importExecute performs the bulk import.
// POST /api/import/execute
func importExecute(e *core.RequestEvent) error {
	var req ExecuteRequest
	if err := e.BindBody(&req); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}

	if req.Collection == "" {
		return e.BadRequestError("Collection is required", nil)
	}

	if req.Data == "" {
		return e.BadRequestError("Data is required", nil)
	}

	if len(req.Mapping) == 0 {
		return e.BadRequestError("Mapping is required", nil)
	}

	// Find the collection
	collection, err := e.App.FindCollectionByNameOrId(req.Collection)
	if err != nil {
		return e.BadRequestError("Collection not found: "+req.Collection, err)
	}

	// Parse the data
	var headers []string
	var allRows [][]string

	switch strings.ToLower(req.Format) {
	case "csv":
		h, rows, _, err := parseCSV(req.Data, req.Delimiter, -1) // -1 means all rows
		if err != nil {
			return e.BadRequestError("Failed to parse CSV: "+err.Error(), err)
		}
		headers = h
		allRows = rows
	case "json":
		h, rows, _, err := parseJSON(req.Data, -1) // -1 means all rows
		if err != nil {
			return e.BadRequestError("Failed to parse JSON: "+err.Error(), err)
		}
		headers = h
		allRows = rows
	default:
		return e.BadRequestError("Unsupported format. Use 'csv' or 'json'", nil)
	}

	// Create header index map
	headerIndex := make(map[string]int)
	for i, h := range headers {
		headerIndex[h] = i
	}

	response := ExecuteResponse{
		TotalRows: len(allRows),
	}

	// Process each row
	for rowNum, row := range allRows {
		// Create record data
		recordData := make(map[string]any)
		
		for sourceCol, targetField := range req.Mapping {
			if targetField == "" || targetField == "-" {
				continue
			}

			idx, exists := headerIndex[sourceCol]
			if !exists {
				continue
			}

			if idx < len(row) {
				value := strings.TrimSpace(row[idx])
				if value != "" {
					recordData[targetField] = value
				}
			}
		}

		// Create the record
		record := core.NewRecord(collection)
		for key, value := range recordData {
			record.Set(key, value)
		}

		// Save the record
		if err := e.App.Save(record); err != nil {
			response.FailureCount++
			response.Errors = append(response.Errors, ImportError{
				Row:     rowNum + 1, // 1-indexed for user display
				Message: err.Error(),
				Data:    recordData,
			})
		} else {
			response.SuccessCount++
		}
	}

	return e.JSON(http.StatusOK, response)
}

// parseCSV parses CSV data and returns headers, sample rows, and total count.
// If sampleSize is -1, returns all rows.
func parseCSV(data string, delimiter string, sampleSize int) ([]string, [][]string, int, error) {
	reader := csv.NewReader(strings.NewReader(data))
	
	// Set delimiter
	if delimiter == "\\t" || delimiter == "\t" {
		reader.Comma = '\t'
	} else if delimiter != "" && len(delimiter) == 1 {
		reader.Comma = rune(delimiter[0])
	}
	// Default is comma

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, nil, 0, fmt.Errorf("CSV file is empty")
	}

	// First row is headers
	headers := records[0]
	
	// Remaining rows are data
	dataRows := records[1:]
	totalRows := len(dataRows)

	// Get sample rows
	var sampleRows [][]string
	if sampleSize == -1 {
		sampleRows = dataRows
	} else {
		limit := sampleSize
		if limit > len(dataRows) {
			limit = len(dataRows)
		}
		sampleRows = dataRows[:limit]
	}

	return headers, sampleRows, totalRows, nil
}

// parseJSON parses JSON array data and returns headers, sample rows, and total count.
// If sampleSize is -1, returns all rows.
func parseJSON(data string, sampleSize int) ([]string, [][]string, int, error) {
	// Try to parse as array of objects
	var objects []map[string]any
	
	decoder := json.NewDecoder(bytes.NewReader([]byte(data)))
	decoder.UseNumber()
	
	if err := decoder.Decode(&objects); err != nil {
		return nil, nil, 0, fmt.Errorf("failed to parse JSON: expected array of objects: %w", err)
	}

	if len(objects) == 0 {
		return nil, nil, 0, fmt.Errorf("JSON array is empty")
	}

	// Extract all unique keys as headers (from all objects)
	headerSet := make(map[string]bool)
	for _, obj := range objects {
		for key := range obj {
			headerSet[key] = true
		}
	}

	// Convert to sorted slice for consistent ordering
	headers := make([]string, 0, len(headerSet))
	for key := range headerSet {
		headers = append(headers, key)
	}
	// Sort alphabetically for consistent output
	sortStrings(headers)

	// Convert objects to rows
	totalRows := len(objects)
	var rows [][]string

	limit := totalRows
	if sampleSize != -1 && sampleSize < limit {
		limit = sampleSize
	}

	for i := 0; i < limit; i++ {
		row := make([]string, len(headers))
		for j, header := range headers {
			if val, exists := objects[i][header]; exists {
				row[j] = valueToString(val)
			}
		}
		rows = append(rows, row)
	}

	return headers, rows, totalRows, nil
}

// detectFormat tries to detect the file format from content.
func detectFormat(data string) string {
	trimmed := strings.TrimSpace(data)
	if strings.HasPrefix(trimmed, "[") || strings.HasPrefix(trimmed, "{") {
		return "json"
	}
	return "csv"
}

// valueToString converts any value to string for display.
func valueToString(val any) string {
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	case float64:
		return fmt.Sprintf("%v", v)
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		// For complex types, marshal to JSON
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(b)
	}
}

// sortStrings sorts a string slice in place.
func sortStrings(s []string) {
	for i := 0; i < len(s)-1; i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] > s[j] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}

// GetCollectionFields returns all field names and their types for a collection.
// This is a helper function for the import wizard to show available fields.
func getCollectionFields(app core.App, collectionNameOrId string) (map[string]string, error) {
	collection, err := app.FindCollectionByNameOrId(collectionNameOrId)
	if err != nil {
		return nil, err
	}

	fields := make(map[string]string)
	for _, field := range collection.Fields {
		fields[field.GetName()] = field.Type()
	}

	return fields, nil
}
