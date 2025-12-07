package sql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
)

// ExecutionResult represents the result of executing a SQL statement
type ExecutionResult struct {
	Type         StatementType    `json:"type"`
	Success      bool             `json:"success"`
	Message      string           `json:"message"`
	RowsAffected int64            `json:"rowsAffected,omitempty"`
	LastInsertID string           `json:"lastInsertId,omitempty"`
	Columns      []string         `json:"columns,omitempty"`
	Rows         []map[string]any `json:"rows,omitempty"`
	TotalRows    int              `json:"totalRows,omitempty"`
	ExecutionMs  int64            `json:"executionMs"`
}

// Executor handles SQL statement execution against PocketBase
type Executor struct {
	app     core.App
	timeout time.Duration
}

// NewExecutor creates a new SQL executor
func NewExecutor(app core.App) *Executor {
	return &Executor{
		app:     app,
		timeout: 30 * time.Second,
	}
}

// SetTimeout sets the execution timeout
func (e *Executor) SetTimeout(d time.Duration) {
	e.timeout = d
}

// Execute parses and executes a SQL statement
func (e *Executor) Execute(ctx context.Context, sqlStr string) (*ExecutionResult, error) {
	start := time.Now()

	// Parse the SQL statement
	stmt, err := ParseSQL(sqlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SQL: %w", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	var result *ExecutionResult

	switch stmt.Type {
	case StatementSelect:
		result, err = e.executeSelect(ctx, stmt)
	case StatementInsert:
		result, err = e.executeInsert(ctx, stmt)
	case StatementUpdate:
		result, err = e.executeUpdate(ctx, stmt)
	case StatementDelete:
		result, err = e.executeDelete(ctx, stmt)
	case StatementCreateTable:
		result, err = e.executeCreateTable(ctx, stmt)
	case StatementAlterTable:
		result, err = e.executeAlterTable(ctx, stmt)
	case StatementDropTable:
		result, err = e.executeDropTable(ctx, stmt)
	default:
		return nil, errors.New("unsupported SQL statement type")
	}

	if err != nil {
		return nil, err
	}

	result.ExecutionMs = time.Since(start).Milliseconds()
	return result, nil
}

// executeSelect executes a SELECT statement directly against SQLite
func (e *Executor) executeSelect(ctx context.Context, stmt *SQLStatement) (*ExecutionResult, error) {
	db := e.app.DB()

	// Use dbx.NewQuery to execute raw SQL
	var results []map[string]any
	rows, err := db.NewQuery(stmt.Raw).Rows()
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	// Scan results
	for rows.Next() {
		// Create slice of interface{} to scan into
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert to map
		row := make(map[string]any)
		for i, col := range columns {
			val := values[i]
			// Convert []byte to string for readability
			if b, ok := val.([]byte); ok {
				val = string(b)
			}
			row[col] = val
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return &ExecutionResult{
		Type:      StatementSelect,
		Success:   true,
		Message:   fmt.Sprintf("Query returned %d rows", len(results)),
		Columns:   columns,
		Rows:      results,
		TotalRows: len(results),
	}, nil
}

// executeInsert executes an INSERT statement via PocketBase Records API
// Supports multi-row INSERT
func (e *Executor) executeInsert(ctx context.Context, stmt *SQLStatement) (*ExecutionResult, error) {
	if len(stmt.Tables) == 0 {
		return nil, errors.New("no table specified for INSERT")
	}

	tableName := stmt.Tables[0]

	// Find the collection
	collection, err := e.app.FindCachedCollectionByNameOrId(tableName)
	if err != nil {
		return nil, fmt.Errorf("collection '%s' not found", tableName)
	}

	// Determine which values to insert (MultiValues for multi-row, Values for single-row)
	var rowsToInsert []map[string]any
	if len(stmt.MultiValues) > 0 {
		rowsToInsert = stmt.MultiValues
	} else if len(stmt.Values) > 0 {
		rowsToInsert = []map[string]any{stmt.Values}
	} else {
		return nil, errors.New("no values provided for INSERT")
	}

	var insertedCount int64
	var lastInsertID string

	for _, rowValues := range rowsToInsert {
		// Create a new record
		record := core.NewRecord(collection)

		// Set values from INSERT statement
		for colName, value := range rowValues {
			// Skip system fields
			if colName == "id" || colName == "created" || colName == "updated" {
				continue
			}
			record.Set(colName, value)
		}

		// Save the record
		if err := e.app.Save(record); err != nil {
			return nil, fmt.Errorf("failed to insert record: %w", err)
		}
		insertedCount++
		lastInsertID = record.Id
	}

	message := "1 row inserted"
	if insertedCount > 1 {
		message = fmt.Sprintf("%d rows inserted", insertedCount)
	}

	return &ExecutionResult{
		Type:         StatementInsert,
		Success:      true,
		Message:      message,
		RowsAffected: insertedCount,
		LastInsertID: lastInsertID,
	}, nil
}

// executeUpdate executes an UPDATE statement via PocketBase Records API
func (e *Executor) executeUpdate(ctx context.Context, stmt *SQLStatement) (*ExecutionResult, error) {
	if len(stmt.Tables) == 0 {
		return nil, errors.New("no table specified for UPDATE")
	}

	tableName := stmt.Tables[0]

	// Find the collection
	collection, err := e.app.FindCachedCollectionByNameOrId(tableName)
	if err != nil {
		return nil, fmt.Errorf("collection '%s' not found", tableName)
	}

	// Convert SQL WHERE to PocketBase filter
	filter := convertWhereToFilter(stmt.Where)

	// Find matching records
	records, err := e.app.FindRecordsByFilter(collection.Id, filter, "", 0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to find records: %w", err)
	}

	// Update each record
	var updated int64
	for _, record := range records {
		for colName, value := range stmt.SetClauses {
			// Skip system fields (except updated which PocketBase handles)
			if colName == "id" || colName == "created" {
				continue
			}
			// Handle special SQL expressions
			if valStr, ok := value.(string); ok {
				// Check for datetime('now') expressions
				if strings.Contains(strings.ToLower(valStr), "datetime('now')") {
					value = time.Now().UTC().Format(time.RFC3339)
				}
			}
			record.Set(colName, value)
		}

		if err := e.app.Save(record); err != nil {
			return nil, fmt.Errorf("failed to update record %s: %w", record.Id, err)
		}
		updated++
	}

	return &ExecutionResult{
		Type:         StatementUpdate,
		Success:      true,
		Message:      fmt.Sprintf("%d row(s) updated", updated),
		RowsAffected: updated,
	}, nil
}

// executeDelete executes a DELETE statement via PocketBase Records API
func (e *Executor) executeDelete(ctx context.Context, stmt *SQLStatement) (*ExecutionResult, error) {
	if len(stmt.Tables) == 0 {
		return nil, errors.New("no table specified for DELETE")
	}

	tableName := stmt.Tables[0]

	// Find the collection
	collection, err := e.app.FindCachedCollectionByNameOrId(tableName)
	if err != nil {
		return nil, fmt.Errorf("collection '%s' not found", tableName)
	}

	// Convert SQL WHERE to PocketBase filter
	filter := convertWhereToFilter(stmt.Where)

	if filter == "" {
		return nil, errors.New("DELETE without WHERE clause is not allowed for safety")
	}

	// Find matching records
	records, err := e.app.FindRecordsByFilter(collection.Id, filter, "", 0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to find records: %w", err)
	}

	// Delete each record
	var deleted int64
	for _, record := range records {
		if err := e.app.Delete(record); err != nil {
			return nil, fmt.Errorf("failed to delete record %s: %w", record.Id, err)
		}
		deleted++
	}

	return &ExecutionResult{
		Type:         StatementDelete,
		Success:      true,
		Message:      fmt.Sprintf("%d row(s) deleted", deleted),
		RowsAffected: deleted,
	}, nil
}

// executeCreateTable creates a new PocketBase collection
func (e *Executor) executeCreateTable(ctx context.Context, stmt *SQLStatement) (*ExecutionResult, error) {
	if len(stmt.Tables) == 0 {
		return nil, errors.New("no table name specified for CREATE TABLE")
	}

	tableName := stmt.Tables[0]

	// Check if collection already exists
	if existing, _ := e.app.FindCachedCollectionByNameOrId(tableName); existing != nil {
		return nil, fmt.Errorf("collection '%s' already exists", tableName)
	}

	// Create new collection
	collection := core.NewCollection(core.CollectionTypeBase, tableName)

	// Add fields from column definitions
	for _, col := range stmt.Columns {
		// Skip id, created, updated - PocketBase adds these automatically
		lowerName := strings.ToLower(col.Name)
		if lowerName == "id" || lowerName == "created" || lowerName == "updated" {
			continue
		}

		field := MapColumnToField(col)
		if field != nil {
			collection.Fields.Add(field)
		}
	}

	// Set default rules (allow all for authenticated users)
	emptyRule := ""
	collection.ListRule = &emptyRule
	collection.ViewRule = &emptyRule
	collection.CreateRule = &emptyRule
	collection.UpdateRule = &emptyRule
	collection.DeleteRule = &emptyRule

	// Save the collection
	if err := e.app.Save(collection); err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	return &ExecutionResult{
		Type:    StatementCreateTable,
		Success: true,
		Message: fmt.Sprintf("Collection '%s' created successfully", tableName),
	}, nil
}

// executeAlterTable modifies a PocketBase collection
func (e *Executor) executeAlterTable(ctx context.Context, stmt *SQLStatement) (*ExecutionResult, error) {
	if len(stmt.Tables) == 0 {
		return nil, errors.New("no table name specified for ALTER TABLE")
	}

	tableName := stmt.Tables[0]

	// Find the collection
	collection, err := e.app.FindCachedCollectionByNameOrId(tableName)
	if err != nil {
		return nil, fmt.Errorf("collection '%s' not found", tableName)
	}

	// Add new columns
	for _, col := range stmt.Columns {
		field := MapColumnToField(col)
		if field != nil {
			collection.Fields.Add(field)
		}
	}

	// Save the collection
	if err := e.app.Save(collection); err != nil {
		return nil, fmt.Errorf("failed to alter collection: %w", err)
	}

	return &ExecutionResult{
		Type:    StatementAlterTable,
		Success: true,
		Message: fmt.Sprintf("Collection '%s' altered successfully", tableName),
	}, nil
}

// executeDropTable deletes a PocketBase collection
func (e *Executor) executeDropTable(ctx context.Context, stmt *SQLStatement) (*ExecutionResult, error) {
	if len(stmt.Tables) == 0 {
		return nil, errors.New("no table name specified for DROP TABLE")
	}

	tableName := stmt.Tables[0]

	// Find the collection
	collection, err := e.app.FindCachedCollectionByNameOrId(tableName)
	if err != nil {
		return nil, fmt.Errorf("collection '%s' not found", tableName)
	}

	// Delete the collection
	if err := e.app.Delete(collection); err != nil {
		return nil, fmt.Errorf("failed to drop collection: %w", err)
	}

	return &ExecutionResult{
		Type:    StatementDropTable,
		Success: true,
		Message: fmt.Sprintf("Collection '%s' dropped successfully", tableName),
	}, nil
}

// convertWhereToFilter converts a SQL WHERE clause to PocketBase filter syntax
// This is a simplified conversion - complex SQL may not convert perfectly
func convertWhereToFilter(where string) string {
	if where == "" {
		return ""
	}

	filter := where

	// Convert SQL operators to PocketBase operators
	// Note: This is simplified - real conversion would be more complex
	
	// Replace AND with &&
	filter = strings.ReplaceAll(filter, " AND ", " && ")
	filter = strings.ReplaceAll(filter, " and ", " && ")
	
	// Replace OR with ||
	filter = strings.ReplaceAll(filter, " OR ", " || ")
	filter = strings.ReplaceAll(filter, " or ", " || ")

	// Replace <> with !=
	filter = strings.ReplaceAll(filter, "<>", "!=")

	// Replace LIKE with ~
	filter = strings.ReplaceAll(filter, " LIKE ", " ~ ")
	filter = strings.ReplaceAll(filter, " like ", " ~ ")

	// Replace NOT LIKE with !~
	filter = strings.ReplaceAll(filter, " NOT LIKE ", " !~ ")
	filter = strings.ReplaceAll(filter, " not like ", " !~ ")

	// Replace single quotes with double quotes for string values
	// This is a simplified approach - might not work for all cases
	filter = strings.ReplaceAll(filter, "'", "\"")

	// Handle IS NULL / IS NOT NULL
	filter = strings.ReplaceAll(filter, " IS NULL", " = null")
	filter = strings.ReplaceAll(filter, " is null", " = null")
	filter = strings.ReplaceAll(filter, " IS NOT NULL", " != null")
	filter = strings.ReplaceAll(filter, " is not null", " != null")

	return filter
}

// ExecuteRawSQL executes raw SQL directly against the database
// WARNING: This bypasses PocketBase's collection/record APIs
// Use with caution - mainly for SELECT queries
func (e *Executor) ExecuteRawSQL(ctx context.Context, sqlStr string) (*ExecutionResult, error) {
	start := time.Now()
	db := e.app.DB()

	// Determine if it's a query (SELECT) or exec (INSERT/UPDATE/DELETE)
	upperSQL := strings.TrimSpace(strings.ToUpper(sqlStr))

	if strings.HasPrefix(upperSQL, "SELECT") {
		rows, err := db.NewQuery(sqlStr).Rows()
		if err != nil {
			return nil, fmt.Errorf("query execution failed: %w", err)
		}
		defer rows.Close()

		columns, _ := rows.Columns()
		var results []map[string]any

		for rows.Next() {
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				continue
			}

			row := make(map[string]any)
			for i, col := range columns {
				val := values[i]
				if b, ok := val.([]byte); ok {
					val = string(b)
				}
				row[col] = val
			}
			results = append(results, row)
		}

		return &ExecutionResult{
			Type:        StatementSelect,
			Success:     true,
			Message:     fmt.Sprintf("Query returned %d rows", len(results)),
			Columns:     columns,
			Rows:        results,
			TotalRows:   len(results),
			ExecutionMs: time.Since(start).Milliseconds(),
		}, nil
	}

	// For non-SELECT statements, use Execute
	result, err := db.NewQuery(sqlStr).Execute()
	if err != nil {
		return nil, fmt.Errorf("execution failed: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	lastID, _ := result.LastInsertId()

	return &ExecutionResult{
		Type:         StatementUnknown,
		Success:      true,
		Message:      fmt.Sprintf("%d row(s) affected", rowsAffected),
		RowsAffected: rowsAffected,
		LastInsertID: fmt.Sprintf("%d", lastID),
		ExecutionMs:  time.Since(start).Milliseconds(),
	}, nil
}

// GenerateID generates a unique ID for new records
func GenerateID() string {
	return security.RandomString(15)
}

// SplitStatements splits a SQL string into individual statements
// Handles semicolons inside strings properly
func SplitStatements(sqlStr string) []string {
	var statements []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	
	for i := 0; i < len(sqlStr); i++ {
		char := sqlStr[i]
		
		// Handle escape sequences
		if (inSingleQuote || inDoubleQuote) && char == '\\' && i+1 < len(sqlStr) {
			current.WriteByte(char)
			current.WriteByte(sqlStr[i+1])
			i++
			continue
		}
		
		// Toggle quote states
		if char == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
		} else if char == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
		}
		
		// Check for statement separator (semicolon outside quotes)
		if char == ';' && !inSingleQuote && !inDoubleQuote {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			current.Reset()
			continue
		}
		
		current.WriteByte(char)
	}
	
	// Add the last statement if any
	stmt := strings.TrimSpace(current.String())
	if stmt != "" {
		statements = append(statements, stmt)
	}
	
	return statements
}

// MultiExecutionResult represents results from multiple SQL statements
type MultiExecutionResult struct {
	TotalStatements int                `json:"totalStatements"`
	Successful      int                `json:"successful"`
	Failed          int                `json:"failed"`
	Results         []*ExecutionResult `json:"results"`
	TotalMs         int64              `json:"totalMs"`
}

// ExecuteMultiple executes multiple SQL statements in sequence
func (e *Executor) ExecuteMultiple(ctx context.Context, sqlStr string) (*MultiExecutionResult, error) {
	start := time.Now()
	
	statements := SplitStatements(sqlStr)
	if len(statements) == 0 {
		return nil, errors.New("no SQL statements provided")
	}
	
	// If only one statement, just execute it normally
	if len(statements) == 1 {
		result, err := e.Execute(ctx, statements[0])
		if err != nil {
			return &MultiExecutionResult{
				TotalStatements: 1,
				Successful:      0,
				Failed:          1,
				Results: []*ExecutionResult{{
					Success: false,
					Message: err.Error(),
				}},
				TotalMs: time.Since(start).Milliseconds(),
			}, nil
		}
		return &MultiExecutionResult{
			TotalStatements: 1,
			Successful:      1,
			Failed:          0,
			Results:         []*ExecutionResult{result},
			TotalMs:         time.Since(start).Milliseconds(),
		}, nil
	}
	
	// Execute multiple statements
	multiResult := &MultiExecutionResult{
		TotalStatements: len(statements),
		Results:         make([]*ExecutionResult, 0, len(statements)),
	}
	
	for _, stmt := range statements {
		result, err := e.Execute(ctx, stmt)
		if err != nil {
			multiResult.Failed++
			multiResult.Results = append(multiResult.Results, &ExecutionResult{
				Success: false,
				Message: err.Error(),
			})
			// Continue executing remaining statements
			continue
		}
		multiResult.Successful++
		multiResult.Results = append(multiResult.Results, result)
	}
	
	multiResult.TotalMs = time.Since(start).Milliseconds()
	return multiResult, nil
}

// ValidateStatement checks if a SQL statement is safe to execute
func ValidateStatement(stmt *SQLStatement) error {
	if stmt == nil {
		return errors.New("nil statement")
	}

	// Don't allow empty table names for non-SELECT statements
	if len(stmt.Tables) == 0 && stmt.Type != StatementSelect {
		return errors.New("no table specified")
	}

	// Require WHERE for DELETE (safety check)
	if stmt.Type == StatementDelete && stmt.Where == "" {
		return errors.New("DELETE without WHERE clause is not allowed")
	}

	// Allow SELECT on system collections, but block modifications
	if stmt.Type == StatementSelect {
		return nil // SELECT is always safe
	}

	// Don't allow modifying system collections (INSERT, UPDATE, DELETE, CREATE, ALTER, DROP)
	for _, table := range stmt.Tables {
		if strings.HasPrefix(table, "_") {
			return fmt.Errorf("cannot modify system collection: %s", table)
		}
	}

	return nil
}

