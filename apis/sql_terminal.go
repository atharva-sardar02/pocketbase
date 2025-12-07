package apis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/services/ai"
	"github.com/pocketbase/pocketbase/services/sql"
	"github.com/pocketbase/pocketbase/tools/router"
)

// SQLExecuteRequest represents the request body for SQL execution
type SQLExecuteRequest struct {
	SQL     string `json:"sql"`
	Confirm bool   `json:"confirm"` // Required for destructive operations
}

// SQLExecuteResponse represents the response from SQL execution
type SQLExecuteResponse struct {
	Success      bool             `json:"success"`
	Type         string           `json:"type"`
	Message      string           `json:"message"`
	Columns      []string         `json:"columns,omitempty"`
	Rows         []map[string]any `json:"rows,omitempty"`
	TotalRows    int              `json:"totalRows,omitempty"`
	RowsAffected int64            `json:"rowsAffected,omitempty"`
	ExecutionMs  int64            `json:"executionMs"`
	Error        string           `json:"error,omitempty"`
	// Multi-statement fields
	IsMulti          bool                   `json:"isMulti,omitempty"`
	TotalStatements  int                    `json:"totalStatements,omitempty"`
	SuccessfulCount  int                    `json:"successfulCount,omitempty"`
	FailedCount      int                    `json:"failedCount,omitempty"`
	Results          []*SQLExecuteResponse  `json:"results,omitempty"`
}

// SQLAIRequest represents the request body for AI-powered SQL generation
type SQLAIRequest struct {
	Query   string `json:"query"`   // Natural language query
	Execute bool   `json:"execute"` // Whether to execute the generated SQL
	Confirm bool   `json:"confirm"` // Required for destructive operations
}

// SQLAIResponse represents the response from AI SQL generation
type SQLAIResponse struct {
	SQL          string           `json:"sql"`
	IsDestructive bool            `json:"isDestructive"`
	RequiresConfirm bool          `json:"requiresConfirm"`
	Executed     bool             `json:"executed"`
	Result       *SQLExecuteResponse `json:"result,omitempty"`
}

// SQLSchemaResponse represents the database schema for the schema browser
type SQLSchemaResponse struct {
	Collections []CollectionSchema `json:"collections"`
}

// CollectionSchema represents a collection's schema
type CollectionSchema struct {
	Name   string        `json:"name"`
	Type   string        `json:"type"`
	Fields []FieldSchema `json:"fields"`
}

// FieldSchema represents a field's schema
type FieldSchema struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Required bool     `json:"required"`
	Options  []string `json:"options,omitempty"`  // For select fields
	Relation string   `json:"relation,omitempty"` // For relation fields
}

// bindSQLTerminalApi registers the SQL terminal API endpoints
func bindSQLTerminalApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	subGroup := rg.Group("/sql").Bind(RequireAuth())
	
	// Execute raw SQL
	subGroup.POST("/execute", sqlExecute)
	
	// AI-powered SQL generation
	subGroup.POST("/ai", sqlAI)
	
	// Get database schema for schema browser
	subGroup.GET("/schema", sqlSchema)
}

// sqlExecute handles direct SQL execution (supports multiple statements)
func sqlExecute(e *core.RequestEvent) error {
	// Parse request body
	var req SQLExecuteRequest
	if err := e.BindBody(&req); err != nil {
		return e.BadRequestError("An error occurred while loading the submitted data.", err)
	}

	// Validate required fields
	if req.SQL == "" {
		return e.BadRequestError("SQL statement is required.", nil)
	}

	// Split into multiple statements
	statements := sql.SplitStatements(req.SQL)
	if len(statements) == 0 {
		return e.BadRequestError("No valid SQL statements found.", nil)
	}

	executor := sql.NewExecutor(e.App)
	ctx := context.Background()

	// Handle single statement (original behavior)
	if len(statements) == 1 {
		stmt, err := sql.ParseSQL(statements[0])
		if err != nil {
			return e.BadRequestError("Invalid SQL syntax.", err)
		}

		if err := sql.ValidateStatement(stmt); err != nil {
			return e.BadRequestError(err.Error(), nil)
		}

		if stmt.RequiresConfirmation() && !req.Confirm {
			return e.JSON(http.StatusOK, SQLExecuteResponse{
				Success: false,
				Type:    string(stmt.Type),
				Message: "This operation requires confirmation. Set 'confirm: true' to proceed.",
				Error:   "confirmation_required",
			})
		}

		result, err := executor.Execute(ctx, statements[0])
		if err != nil {
			return e.BadRequestError("SQL execution failed.", err)
		}

		return e.JSON(http.StatusOK, SQLExecuteResponse{
			Success:      result.Success,
			Type:         string(result.Type),
			Message:      result.Message,
			Columns:      result.Columns,
			Rows:         result.Rows,
			TotalRows:    result.TotalRows,
			RowsAffected: result.RowsAffected,
			ExecutionMs:  result.ExecutionMs,
		})
	}

	// Handle multiple statements
	// First, validate all statements and check for confirmation requirements
	needsConfirm := false
	for _, stmtStr := range statements {
		stmt, err := sql.ParseSQL(stmtStr)
		if err != nil {
			return e.BadRequestError("Invalid SQL syntax in statement: "+stmtStr, err)
		}
		if err := sql.ValidateStatement(stmt); err != nil {
			return e.BadRequestError(err.Error(), nil)
		}
		if stmt.RequiresConfirmation() {
			needsConfirm = true
		}
	}

	if needsConfirm && !req.Confirm {
		return e.JSON(http.StatusOK, SQLExecuteResponse{
			Success:         false,
			IsMulti:         true,
			TotalStatements: len(statements),
			Message:         "One or more statements require confirmation. Set 'confirm: true' to proceed.",
			Error:           "confirmation_required",
		})
	}

	// Execute all statements
	multiResult, err := executor.ExecuteMultiple(ctx, req.SQL)
	if err != nil {
		return e.BadRequestError("SQL execution failed.", err)
	}

	// Convert results to response format
	var results []*SQLExecuteResponse
	var lastSelectResult *sql.ExecutionResult
	var totalRowsAffected int64

	for _, r := range multiResult.Results {
		resp := &SQLExecuteResponse{
			Success:      r.Success,
			Type:         string(r.Type),
			Message:      r.Message,
			Columns:      r.Columns,
			Rows:         r.Rows,
			TotalRows:    r.TotalRows,
			RowsAffected: r.RowsAffected,
			ExecutionMs:  r.ExecutionMs,
		}
		results = append(results, resp)
		totalRowsAffected += r.RowsAffected
		
		// Track last SELECT for displaying results
		if r.Type == sql.StatementSelect && len(r.Rows) > 0 {
			lastSelectResult = r
		}
	}

	response := SQLExecuteResponse{
		Success:         multiResult.Failed == 0,
		IsMulti:         true,
		TotalStatements: multiResult.TotalStatements,
		SuccessfulCount: multiResult.Successful,
		FailedCount:     multiResult.Failed,
		Message:         generateMultiMessage(multiResult),
		RowsAffected:    totalRowsAffected,
		ExecutionMs:     multiResult.TotalMs,
		Results:         results,
	}

	// If there was a SELECT, include its results at the top level for easy display
	if lastSelectResult != nil {
		response.Columns = lastSelectResult.Columns
		response.Rows = lastSelectResult.Rows
		response.TotalRows = lastSelectResult.TotalRows
	}

	return e.JSON(http.StatusOK, response)
}

// generateMultiMessage creates a summary message for multi-statement execution
func generateMultiMessage(result *sql.MultiExecutionResult) string {
	if result.Failed == 0 {
		return fmt.Sprintf("All %d statements executed successfully", result.TotalStatements)
	}
	return fmt.Sprintf("%d of %d statements executed successfully, %d failed", 
		result.Successful, result.TotalStatements, result.Failed)
}

// sqlAI handles AI-powered SQL generation
func sqlAI(e *core.RequestEvent) error {
	// Parse request body
	var req SQLAIRequest
	if err := e.BindBody(&req); err != nil {
		return e.BadRequestError("An error occurred while loading the submitted data.", err)
	}

	// Validate required fields
	if req.Query == "" {
		return e.BadRequestError("Query is required.", nil)
	}

	// Load AI settings
	settings := e.App.Settings()
	if !settings.AI.Enabled {
		return e.BadRequestError("AI Query feature is not enabled.", nil)
	}

	// Extract full database schema for SQL generation
	schema := ai.ExtractAllSchemas(e.App)

	// Build prompts using SQL terminal mode
	systemPrompt := ai.BuildSQLTerminalPrompt(schema)
	userPrompt := ai.BuildUserPrompt(req.Query)

	// Create OpenAI client
	client := ai.NewOpenAIClient(settings.AI)

	// Call LLM
	ctx := context.Background()
	generatedSQL, err := client.SendCompletion(ctx, systemPrompt, userPrompt)
	if err != nil {
		return e.BadRequestError("Failed to generate SQL.", err)
	}

	// Clean the generated SQL
	generatedSQL = cleanSQLResponse(generatedSQL)

	// Split into multiple statements to support multi-statement SQL
	statements := sql.SplitStatements(generatedSQL)
	
	// Check for destructive operations in any statement
	isDestructive := false
	requiresConfirm := false
	for _, stmtStr := range statements {
		stmt, parseErr := sql.ParseSQL(stmtStr)
		if parseErr == nil {
			if stmt.IsDestructive() {
				isDestructive = true
			}
			if stmt.RequiresConfirmation() {
				requiresConfirm = true
			}
		}
	}

	response := SQLAIResponse{
		SQL:             generatedSQL,
		IsDestructive:   isDestructive,
		RequiresConfirm: requiresConfirm,
		Executed:        false,
	}

	// Optionally execute the generated SQL
	if req.Execute {
		executor := sql.NewExecutor(e.App)

		// Handle single statement
		if len(statements) == 1 {
			stmt, parseErr := sql.ParseSQL(statements[0])
			
			// Check confirmation for destructive operations
			if parseErr == nil && stmt.RequiresConfirmation() && !req.Confirm {
				response.Result = &SQLExecuteResponse{
					Success: false,
					Type:    string(stmt.Type),
					Message: "This operation requires confirmation. Set 'confirm: true' to proceed.",
					Error:   "confirmation_required",
				}
				return e.JSON(http.StatusOK, response)
			}

			// Validate statement
			if parseErr != nil {
				response.Result = &SQLExecuteResponse{
					Success: false,
					Message: "Invalid SQL syntax.",
					Error:   parseErr.Error(),
				}
				return e.JSON(http.StatusOK, response)
			}

			if err := sql.ValidateStatement(stmt); err != nil {
				response.Result = &SQLExecuteResponse{
					Success: false,
					Message: err.Error(),
					Error:   "validation_failed",
				}
				return e.JSON(http.StatusOK, response)
			}

			// Execute the SQL
			result, execErr := executor.Execute(ctx, statements[0])
			if execErr != nil {
				response.Result = &SQLExecuteResponse{
					Success: false,
					Message: "SQL execution failed.",
					Error:   execErr.Error(),
				}
				return e.JSON(http.StatusOK, response)
			}

			response.Executed = true
			response.Result = &SQLExecuteResponse{
				Success:      result.Success,
				Type:         string(result.Type),
				Message:      result.Message,
				Columns:      result.Columns,
				Rows:         result.Rows,
				TotalRows:    result.TotalRows,
				RowsAffected: result.RowsAffected,
				ExecutionMs:  result.ExecutionMs,
			}
		} else {
			// Handle multiple statements
			// First, validate all statements and check for confirmation requirements
			needsConfirm := false
			for _, stmtStr := range statements {
				stmt, parseErr := sql.ParseSQL(stmtStr)
				if parseErr != nil {
					response.Result = &SQLExecuteResponse{
						Success: false,
						Message: "Invalid SQL syntax in statement: " + stmtStr,
						Error:   parseErr.Error(),
					}
					return e.JSON(http.StatusOK, response)
				}
				if err := sql.ValidateStatement(stmt); err != nil {
					response.Result = &SQLExecuteResponse{
						Success: false,
						Message: err.Error(),
						Error:   "validation_failed",
					}
					return e.JSON(http.StatusOK, response)
				}
				if stmt.RequiresConfirmation() {
					needsConfirm = true
				}
			}

			if needsConfirm && !req.Confirm {
				response.Result = &SQLExecuteResponse{
					Success:         false,
					IsMulti:         true,
					TotalStatements: len(statements),
					Message:         "One or more statements require confirmation. Set 'confirm: true' to proceed.",
					Error:           "confirmation_required",
				}
				return e.JSON(http.StatusOK, response)
			}

			// Execute all statements
			multiResult, execErr := executor.ExecuteMultiple(ctx, generatedSQL)
			if execErr != nil {
				response.Result = &SQLExecuteResponse{
					Success: false,
					Message: "SQL execution failed.",
					Error:   execErr.Error(),
				}
				return e.JSON(http.StatusOK, response)
			}

			// Convert results to response format
			var results []*SQLExecuteResponse
			var lastSelectResult *sql.ExecutionResult
			var totalRowsAffected int64

			for _, r := range multiResult.Results {
				resp := &SQLExecuteResponse{
					Success:      r.Success,
					Type:         string(r.Type),
					Message:      r.Message,
					Columns:      r.Columns,
					Rows:         r.Rows,
					TotalRows:    r.TotalRows,
					RowsAffected: r.RowsAffected,
					ExecutionMs:  r.ExecutionMs,
				}
				results = append(results, resp)
				totalRowsAffected += r.RowsAffected

				// Track last SELECT for displaying results
				if r.Type == sql.StatementSelect && len(r.Rows) > 0 {
					lastSelectResult = r
				}
			}

			response.Executed = true
			response.Result = &SQLExecuteResponse{
				Success:         multiResult.Failed == 0,
				IsMulti:         true,
				TotalStatements: multiResult.TotalStatements,
				SuccessfulCount: multiResult.Successful,
				FailedCount:     multiResult.Failed,
				Message:         generateMultiMessage(multiResult),
				RowsAffected:    totalRowsAffected,
				ExecutionMs:     multiResult.TotalMs,
				Results:         results,
			}

			// If there was a SELECT, include its results at the top level for easy display
			if lastSelectResult != nil {
				response.Result.Columns = lastSelectResult.Columns
				response.Result.Rows = lastSelectResult.Rows
				response.Result.TotalRows = lastSelectResult.TotalRows
			}
		}
	}

	return e.JSON(http.StatusOK, response)
}

// sqlSchema returns the database schema for the schema browser
func sqlSchema(e *core.RequestEvent) error {
	collections, err := e.App.FindAllCollections()
	if err != nil {
		return e.BadRequestError("Failed to load collections.", err)
	}

	var schemas []CollectionSchema
	for _, coll := range collections {
		// Skip system collections
		if len(coll.Name) > 0 && coll.Name[0] == '_' {
			continue
		}

		schema := CollectionSchema{
			Name: coll.Name,
			Type: coll.Type,
		}

		for _, field := range coll.Fields {
			if field.GetHidden() {
				continue
			}

			fieldSchema := FieldSchema{
				Name: field.GetName(),
				Type: field.Type(),
			}

			// Handle select field options
			if selectField, ok := field.(*core.SelectField); ok {
				fieldSchema.Options = selectField.Values
			}

			// Handle relation field
			if relField, ok := field.(*core.RelationField); ok {
				if relField.CollectionId != "" {
					if relColl, err := e.App.FindCachedCollectionByNameOrId(relField.CollectionId); err == nil {
						fieldSchema.Relation = relColl.Name
					}
				}
			}

			schema.Fields = append(schema.Fields, fieldSchema)
		}

		schemas = append(schemas, schema)
	}

	return e.JSON(http.StatusOK, SQLSchemaResponse{
		Collections: schemas,
	})
}
