package apis

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/services/ai"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/search"
)

// AIQueryRequest represents the request body for the AI query endpoint.
type AIQueryRequest struct {
	Collection string `json:"collection"`
	Query      string `json:"query"`
	Execute    bool   `json:"execute"`
	Page       int    `json:"page"`
	PerPage    int    `json:"perPage"`
	Mode       string `json:"mode"` // "filter" (default), "dual", or "sql"
}

// AIQueryResponse represents the response from the AI query endpoint.
type AIQueryResponse struct {
	Filter      string        `json:"filter"`
	SQL         string        `json:"sql,omitempty"`         // V2: SQL query equivalent
	RequiresSQL bool          `json:"requiresSQL,omitempty"` // V2: true if query requires SQL (JOINs, aggregates)
	CanUseFilter bool         `json:"canUseFilter"`          // V2: true if filter syntax can express this query
	Results     []interface{} `json:"results,omitempty"`
	TotalItems  int           `json:"totalItems,omitempty"`
	Page        int           `json:"page,omitempty"`
	PerPage     int           `json:"perPage,omitempty"`
	Error       string        `json:"error,omitempty"`
}

// DualOutputLLMResponse is the JSON structure expected from LLM in dual mode
type DualOutputLLMResponse struct {
	Filter      string `json:"filter"`
	SQL         string `json:"sql"`
	RequiresSQL bool   `json:"requiresSQL"`
}

// bindAIQueryApi registers the AI query api endpoints.
func bindAIQueryApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	subGroup := rg.Group("/ai/query").Bind(RequireAuth())
	subGroup.POST("", aiQuery)
	// Also register GET for debugging (will return method error)
	subGroup.GET("", func(e *core.RequestEvent) error {
		return e.BadRequestError("This endpoint only accepts POST requests.", nil)
	})
}

func aiQuery(e *core.RequestEvent) error {
	// Parse request body
	var req AIQueryRequest
	if err := e.BindBody(&req); err != nil {
		return e.BadRequestError("An error occurred while loading the submitted data.", err)
	}

	// Validate required fields
	if req.Collection == "" {
		return e.BadRequestError("Collection is required.", nil)
	}
	if req.Query == "" {
		return e.BadRequestError("Query is required.", nil)
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 30
	}
	if req.Mode == "" {
		req.Mode = "filter" // Default to V1 behavior
	}

	// Load AI settings
	settings := e.App.Settings()
	if !settings.AI.Enabled {
		return e.BadRequestError("AI Query feature is not enabled.", nil)
	}

	// Load collection
	collection, err := e.App.FindCachedCollectionByNameOrId(req.Collection)
	if err != nil {
		return e.NotFoundError("Collection not found.", err)
	}

	// Check if user has access to list records in this collection
	requestInfo, err := e.RequestInfo()
	if err != nil {
		return e.BadRequestError("Failed to get request info.", err)
	}

	// Check collection access (similar to recordsList)
	if collection.ListRule == nil && !requestInfo.HasSuperuserAuth() {
		return e.ForbiddenError("Only superusers can perform this action.", nil)
	}

	// Build response
	var response AIQueryResponse
	var filter string

	// Handle based on mode
	switch req.Mode {
	case "dual":
		// V2: Dual output mode - returns both filter AND SQL
		response, filter, err = handleDualOutputMode(e, collection, req.Query)
		if err != nil {
			return err
		}
	case "sql":
		// V2: SQL-only mode - returns just SQL
		response, err = handleSQLOnlyMode(e, req.Query)
		if err != nil {
			return err
		}
		// SQL mode doesn't execute via PocketBase filter - skip execution
		return e.JSON(http.StatusOK, response)
	default:
		// V1: Filter-only mode (default)
		response, filter, err = handleFilterOnlyMode(e, collection, req.Query)
		if err != nil {
			return err
		}
	}

	// Optionally execute filter and return results
	if req.Execute && filter != "" && response.CanUseFilter {
		// Calculate offset
		offset := (req.Page - 1) * req.PerPage

		// Use the same approach as recordsList to respect listRule
		// We'll use RecordQuery with proper field resolver
		query := e.App.RecordQuery(collection)
		fieldsResolver := core.NewRecordFieldResolver(e.App, collection, requestInfo, true)

		// Apply listRule if user is not superuser
		if !requestInfo.HasSuperuserAuth() && collection.ListRule != nil && *collection.ListRule != "" {
			expr, err := search.FilterData(*collection.ListRule).BuildExpr(fieldsResolver)
			if err != nil {
				return e.BadRequestError("Failed to apply collection list rule.", err)
			}
			query.AndWhere(expr)
		}

		// Apply user's filter
		if filter != "" {
			expr, err := search.FilterData(filter).BuildExpr(fieldsResolver)
			if err != nil {
				return e.BadRequestError("Failed to apply filter.", err)
			}
			query.AndWhere(expr)
		}

		// Update query with any necessary joins
		err = fieldsResolver.UpdateQuery(query)
		if err != nil {
			return e.BadRequestError("Failed to prepare query.", err)
		}

		// Apply pagination
		if offset > 0 {
			query.Offset(int64(offset))
		}
		if req.PerPage > 0 {
			query.Limit(int64(req.PerPage))
		}

		// Execute query
		records := []*core.Record{}
		if err := query.All(&records); err != nil {
			return e.BadRequestError("Failed to execute query.", err)
		}

		// Get total count for pagination
		countQuery := e.App.RecordQuery(collection)
		if !requestInfo.HasSuperuserAuth() && collection.ListRule != nil && *collection.ListRule != "" {
			expr, err := search.FilterData(*collection.ListRule).BuildExpr(fieldsResolver)
			if err == nil {
				countQuery.AndWhere(expr)
			}
		}
		if filter != "" {
			expr, err := search.FilterData(filter).BuildExpr(fieldsResolver)
			if err == nil {
				countQuery.AndWhere(expr)
			}
		}
		fieldsResolver.UpdateQuery(countQuery)

		var totalItems int
		countQuery.Select("COUNT(*)").Row(&totalItems)

		// Convert records to JSON-serializable format
		results := make([]interface{}, len(records))
		for i, record := range records {
			results[i] = record.PublicExport()
		}

		response.Results = results
		response.TotalItems = totalItems
		response.Page = req.Page
		response.PerPage = req.PerPage
	}

	return e.JSON(http.StatusOK, response)
}

// handleFilterOnlyMode handles V1 filter-only mode
func handleFilterOnlyMode(e *core.RequestEvent, collection *core.Collection, query string) (AIQueryResponse, string, error) {
	settings := e.App.Settings()

	// Extract schema for single collection
	schema := ai.ExtractSchema(e.App, collection)

	// Build prompts using filter mode
	systemPrompt := ai.BuildSystemPrompt(schema)
	userPrompt := ai.BuildUserPrompt(query)

	// DEBUG: Log prompts
	e.App.Logger().Debug("AI Query Debug (Filter Mode)",
		"schema", schema,
		"systemPrompt_length", len(systemPrompt),
		"userPrompt", userPrompt,
	)

	// Create OpenAI client
	client := ai.NewOpenAIClient(settings.AI)

	// Call LLM
	ctx := context.Background()
	filter, err := client.SendCompletion(ctx, systemPrompt, userPrompt)
	
	// DEBUG: Log response
	e.App.Logger().Debug("AI Query Response (Filter Mode)",
		"filter", filter,
		"error", err,
	)
	
	if err != nil {
		return AIQueryResponse{}, "", e.BadRequestError("Failed to generate filter from query.", err)
	}

	// Check for INVALID_QUERY response
	if filter == "INVALID_QUERY" {
		return AIQueryResponse{}, "", e.BadRequestError("The query could not be expressed as a filter.", nil)
	}

	// Trim whitespace
	filter = trimFilter(filter)

	// Validate filter
	if err := ai.ValidateFilter(filter, collection); err != nil {
		return AIQueryResponse{}, "", e.BadRequestError("Generated filter is invalid.", err)
	}

	// Build response
	response := AIQueryResponse{
		Filter:       filter,
		CanUseFilter: true,
	}

	return response, filter, nil
}

// handleDualOutputMode handles V2 dual output mode (filter + SQL)
func handleDualOutputMode(e *core.RequestEvent, collection *core.Collection, query string) (AIQueryResponse, string, error) {
	settings := e.App.Settings()

	// Extract schema with related collections for JOIN context
	schema := ai.ExtractSchemaForCollection(e.App, collection.Name)

	// Build prompts using dual mode
	systemPrompt := ai.BuildDualOutputPrompt(schema)
	userPrompt := ai.BuildUserPrompt(query)

	// DEBUG: Log prompts
	e.App.Logger().Debug("AI Query Debug (Dual Mode)",
		"schema", schema,
		"systemPrompt_length", len(systemPrompt),
		"userPrompt", userPrompt,
	)

	// Create OpenAI client
	client := ai.NewOpenAIClient(settings.AI)

	// Call LLM
	ctx := context.Background()
	llmResponse, err := client.SendCompletion(ctx, systemPrompt, userPrompt)
	
	// DEBUG: Log response
	e.App.Logger().Debug("AI Query Response (Dual Mode)",
		"llmResponse", llmResponse,
		"error", err,
	)
	
	if err != nil {
		return AIQueryResponse{}, "", e.BadRequestError("Failed to generate query.", err)
	}

	// Parse JSON response from LLM
	var dualOutput DualOutputLLMResponse
	
	// Clean up response - remove markdown code blocks if present
	cleanedResponse := cleanJSONResponse(llmResponse)
	
	if err := json.Unmarshal([]byte(cleanedResponse), &dualOutput); err != nil {
		e.App.Logger().Warn("Failed to parse dual output JSON, falling back to filter mode",
			"response", llmResponse,
			"cleanedResponse", cleanedResponse,
			"error", err,
		)
		// Fall back to treating the response as a plain filter
		return handleFilterOnlyMode(e, collection, query)
	}

	// Validate filter if present
	filter := strings.TrimSpace(dualOutput.Filter)
	canUseFilter := filter != "" && !dualOutput.RequiresSQL
	
	if canUseFilter {
		if err := ai.ValidateFilter(filter, collection); err != nil {
			// Filter is invalid, but we might still have valid SQL
			e.App.Logger().Warn("Generated filter is invalid, SQL may still work",
				"filter", filter,
				"error", err,
			)
			canUseFilter = false
		}
	}

	// Build response
	response := AIQueryResponse{
		Filter:       filter,
		SQL:          strings.TrimSpace(dualOutput.SQL),
		RequiresSQL:  dualOutput.RequiresSQL,
		CanUseFilter: canUseFilter,
	}

	return response, filter, nil
}

// handleSQLOnlyMode handles V2 SQL-only mode (for SQL Terminal)
func handleSQLOnlyMode(e *core.RequestEvent, query string) (AIQueryResponse, error) {
	settings := e.App.Settings()

	// Extract full database schema
	schema := ai.ExtractAllSchemas(e.App)

	// Build prompts using SQL mode
	systemPrompt := ai.BuildSQLTerminalPrompt(schema)
	userPrompt := ai.BuildUserPrompt(query)

	// DEBUG: Log prompts
	e.App.Logger().Debug("AI Query Debug (SQL Mode)",
		"schema_length", len(schema),
		"systemPrompt_length", len(systemPrompt),
		"userPrompt", userPrompt,
	)

	// Create OpenAI client
	client := ai.NewOpenAIClient(settings.AI)

	// Call LLM
	ctx := context.Background()
	sqlQuery, err := client.SendCompletion(ctx, systemPrompt, userPrompt)
	
	// DEBUG: Log response
	e.App.Logger().Debug("AI Query Response (SQL Mode)",
		"sql", sqlQuery,
		"error", err,
	)
	
	if err != nil {
		return AIQueryResponse{}, e.BadRequestError("Failed to generate SQL query.", err)
	}

	// Clean SQL response
	sqlQuery = cleanSQLResponse(sqlQuery)

	// Build response
	response := AIQueryResponse{
		SQL:          sqlQuery,
		RequiresSQL:  true,
		CanUseFilter: false,
	}

	return response, nil
}

// cleanJSONResponse removes markdown code blocks and other formatting from LLM response
func cleanJSONResponse(response string) string {
	response = strings.TrimSpace(response)
	
	// Remove markdown JSON code blocks
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
	}
	
	return strings.TrimSpace(response)
}

// cleanSQLResponse removes markdown code blocks from SQL response
func cleanSQLResponse(response string) string {
	response = strings.TrimSpace(response)
	
	// Remove markdown SQL code blocks
	if strings.HasPrefix(response, "```sql") {
		response = strings.TrimPrefix(response, "```sql")
		response = strings.TrimSuffix(response, "```")
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
	}
	
	return strings.TrimSpace(response)
}

// trimFilter removes leading/trailing whitespace and any explanatory text
// that the LLM might have added.
func trimFilter(filter string) string {
	// Remove common prefixes/suffixes that LLMs might add
	filter = trimPrefix(filter, "Filter:", "filter:", "FILTER:")
	filter = trimPrefix(filter, "Query:", "query:", "QUERY:")
	filter = trimPrefix(filter, "```", "```sql", "```javascript")
	filter = trimSuffix(filter, "```")

	// Trim whitespace
	return trimWhitespace(filter)
}

func trimPrefix(s string, prefixes ...string) string {
	for _, prefix := range prefixes {
		if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
			s = s[len(prefix):]
			break
		}
	}
	return s
}

func trimSuffix(s string, suffixes ...string) string {
	for _, suffix := range suffixes {
		if len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix {
			s = s[:len(s)-len(suffix)]
			break
		}
	}
	return s
}

func trimWhitespace(s string) string {
	// Trim leading and trailing whitespace
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

