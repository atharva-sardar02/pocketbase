package apis

import (
	"context"
	"net/http"

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
}

// AIQueryResponse represents the response from the AI query endpoint.
type AIQueryResponse struct {
	Filter     string        `json:"filter"`
	Results    []interface{} `json:"results,omitempty"`
	TotalItems int           `json:"totalItems,omitempty"`
	Page       int           `json:"page,omitempty"`
	PerPage    int           `json:"perPage,omitempty"`
	Error      string        `json:"error,omitempty"`
}

// bindAIQueryApi registers the AI query api endpoints.
func bindAIQueryApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	subGroup := rg.Group("/ai/query").Bind(RequireAuth())
	subGroup.POST("", aiQuery)
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

	// Extract schema
	schema := ai.ExtractSchema(e.App, collection)

	// Build prompts
	systemPrompt := ai.BuildSystemPrompt(schema)
	userPrompt := ai.BuildUserPrompt(req.Query)

	// Create OpenAI client
	client := ai.NewOpenAIClient(settings.AI)

	// Call LLM
	ctx := context.Background()
	filter, err := client.SendCompletion(ctx, systemPrompt, userPrompt)
	if err != nil {
		return e.BadRequestError("Failed to generate filter from query.", err)
	}

	// Check for INVALID_QUERY response
	if filter == "INVALID_QUERY" {
		return e.BadRequestError("The query could not be expressed as a filter.", nil)
	}

	// Trim whitespace
	filter = trimFilter(filter)

	// Validate filter
	if err := ai.ValidateFilter(filter, collection); err != nil {
		return e.BadRequestError("Generated filter is invalid.", err)
	}

	// Build response
	response := AIQueryResponse{
		Filter: filter,
	}

	// Optionally execute filter and return results
	if req.Execute {
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

