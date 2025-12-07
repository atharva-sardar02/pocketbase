package sql

import (
	"errors"
	"regexp"
	"strings"
)

// StatementType represents the type of SQL statement
type StatementType string

const (
	StatementSelect      StatementType = "SELECT"
	StatementInsert      StatementType = "INSERT"
	StatementUpdate      StatementType = "UPDATE"
	StatementDelete      StatementType = "DELETE"
	StatementCreateTable StatementType = "CREATE TABLE"
	StatementAlterTable  StatementType = "ALTER TABLE"
	StatementDropTable   StatementType = "DROP TABLE"
	StatementUnknown     StatementType = "UNKNOWN"
)

// ColumnDef represents a column definition in CREATE TABLE
type ColumnDef struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Required   bool     `json:"required"`
	Reference  string   `json:"reference"`  // For foreign key references
	Options    []string `json:"options"`    // For ENUM/CHECK constraints
	Default    string   `json:"default"`    // Default value
	PrimaryKey bool     `json:"primaryKey"`
}

// SQLStatement represents a parsed SQL statement
type SQLStatement struct {
	Type        StatementType        `json:"type"`
	Raw         string               `json:"raw"`
	Tables      []string             `json:"tables"`
	Columns     []ColumnDef          `json:"columns"`
	Where       string               `json:"where"`
	Values      map[string]any       `json:"values"`
	MultiValues []map[string]any     `json:"multiValues"` // For multi-row INSERT
	SetClauses  map[string]any       `json:"setClauses"`
	Joins       []JoinClause         `json:"joins"`
	GroupBy     []string             `json:"groupBy"`
	OrderBy     []OrderByClause      `json:"orderBy"`
	Limit       int                  `json:"limit"`
	Offset      int                  `json:"offset"`
}

// JoinClause represents a JOIN in a SELECT statement
type JoinClause struct {
	Type      string `json:"type"` // INNER, LEFT, RIGHT, FULL
	Table     string `json:"table"`
	Condition string `json:"condition"`
}

// OrderByClause represents an ORDER BY clause
type OrderByClause struct {
	Column string `json:"column"`
	Desc   bool   `json:"desc"`
}

// ParseSQL parses a SQL statement and returns a structured representation
func ParseSQL(sql string) (*SQLStatement, error) {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return nil, errors.New("empty SQL statement")
	}

	// Remove trailing semicolon
	sql = strings.TrimSuffix(sql, ";")

	// Determine statement type
	upperSQL := strings.ToUpper(sql)
	
	stmt := &SQLStatement{
		Raw:        sql,
		Values:     make(map[string]any),
		SetClauses: make(map[string]any),
	}

	switch {
	case strings.HasPrefix(upperSQL, "SELECT"):
		stmt.Type = StatementSelect
		return parseSelect(stmt, sql)
	case strings.HasPrefix(upperSQL, "INSERT"):
		stmt.Type = StatementInsert
		return parseInsert(stmt, sql)
	case strings.HasPrefix(upperSQL, "UPDATE"):
		stmt.Type = StatementUpdate
		return parseUpdate(stmt, sql)
	case strings.HasPrefix(upperSQL, "DELETE"):
		stmt.Type = StatementDelete
		return parseDelete(stmt, sql)
	case strings.HasPrefix(upperSQL, "CREATE TABLE"):
		stmt.Type = StatementCreateTable
		return parseCreateTable(stmt, sql)
	case strings.HasPrefix(upperSQL, "ALTER TABLE"):
		stmt.Type = StatementAlterTable
		return parseAlterTable(stmt, sql)
	case strings.HasPrefix(upperSQL, "DROP TABLE"):
		stmt.Type = StatementDropTable
		return parseDropTable(stmt, sql)
	default:
		stmt.Type = StatementUnknown
		return stmt, errors.New("unsupported SQL statement type")
	}
}

// parseSelect parses a SELECT statement
func parseSelect(stmt *SQLStatement, sql string) (*SQLStatement, error) {
	// Extract table names from FROM clause
	fromRegex := regexp.MustCompile(`(?i)\bFROM\s+(\w+)`)
	fromMatch := fromRegex.FindStringSubmatch(sql)
	if len(fromMatch) > 1 {
		stmt.Tables = append(stmt.Tables, fromMatch[1])
	}

	// Extract JOIN tables
	joinRegex := regexp.MustCompile(`(?i)\b(INNER|LEFT|RIGHT|FULL)?\s*JOIN\s+(\w+)\s+(?:\w+\s+)?ON\s+([^JOIN]+?)(?:\s+(?:INNER|LEFT|RIGHT|FULL)?\s*JOIN|\s+WHERE|\s+GROUP|\s+ORDER|\s+LIMIT|$)`)
	joinMatches := joinRegex.FindAllStringSubmatch(sql, -1)
	for _, match := range joinMatches {
		joinType := "INNER"
		if match[1] != "" {
			joinType = strings.ToUpper(match[1])
		}
		stmt.Joins = append(stmt.Joins, JoinClause{
			Type:      joinType,
			Table:     match[2],
			Condition: strings.TrimSpace(match[3]),
		})
		stmt.Tables = append(stmt.Tables, match[2])
	}

	// Extract WHERE clause
	whereRegex := regexp.MustCompile(`(?i)\bWHERE\s+(.+?)(?:\s+GROUP\s+BY|\s+ORDER\s+BY|\s+LIMIT|$)`)
	whereMatch := whereRegex.FindStringSubmatch(sql)
	if len(whereMatch) > 1 {
		stmt.Where = strings.TrimSpace(whereMatch[1])
	}

	// Extract GROUP BY
	groupByRegex := regexp.MustCompile(`(?i)\bGROUP\s+BY\s+(.+?)(?:\s+HAVING|\s+ORDER\s+BY|\s+LIMIT|$)`)
	groupByMatch := groupByRegex.FindStringSubmatch(sql)
	if len(groupByMatch) > 1 {
		columns := strings.Split(groupByMatch[1], ",")
		for _, col := range columns {
			stmt.GroupBy = append(stmt.GroupBy, strings.TrimSpace(col))
		}
	}

	// Extract ORDER BY
	orderByRegex := regexp.MustCompile(`(?i)\bORDER\s+BY\s+(.+?)(?:\s+LIMIT|$)`)
	orderByMatch := orderByRegex.FindStringSubmatch(sql)
	if len(orderByMatch) > 1 {
		columns := strings.Split(orderByMatch[1], ",")
		for _, col := range columns {
			col = strings.TrimSpace(col)
			desc := strings.HasSuffix(strings.ToUpper(col), " DESC")
			colName := strings.TrimSuffix(col, " DESC")
			colName = strings.TrimSuffix(colName, " ASC")
			stmt.OrderBy = append(stmt.OrderBy, OrderByClause{
				Column: strings.TrimSpace(colName),
				Desc:   desc,
			})
		}
	}

	return stmt, nil
}

// parseInsert parses an INSERT statement (supports multi-row INSERT)
func parseInsert(stmt *SQLStatement, sql string) (*SQLStatement, error) {
	// Extract table name
	tableRegex := regexp.MustCompile(`(?i)\bINSERT\s+INTO\s+(\w+)`)
	tableMatch := tableRegex.FindStringSubmatch(sql)
	if len(tableMatch) > 1 {
		stmt.Tables = append(stmt.Tables, tableMatch[1])
	}

	// Extract column names
	colNamesRegex := regexp.MustCompile(`(?i)\bINSERT\s+INTO\s+\w+\s*\(([^)]+)\)`)
	colNamesMatch := colNamesRegex.FindStringSubmatch(sql)
	if len(colNamesMatch) < 2 {
		return stmt, nil
	}
	
	columnNames := strings.Split(colNamesMatch[1], ",")
	for i := range columnNames {
		columnNames[i] = strings.TrimSpace(columnNames[i])
	}

	// Extract all VALUES rows
	// Find the VALUES keyword and everything after it
	valuesIdx := strings.Index(strings.ToUpper(sql), "VALUES")
	if valuesIdx == -1 {
		return stmt, nil
	}
	
	valuesSection := sql[valuesIdx+6:] // Skip "VALUES"
	
	// Parse each row of values (handles nested parentheses for function calls)
	valueRows := parseMultipleValueRows(valuesSection)
	
	for _, rowValues := range valueRows {
		rowMap := make(map[string]any)
		for i, col := range columnNames {
			if i < len(rowValues) {
				rowMap[col] = parseValue(rowValues[i])
			}
		}
		stmt.MultiValues = append(stmt.MultiValues, rowMap)
	}
	
	// For backward compatibility, set Values to first row if present
	if len(stmt.MultiValues) > 0 {
		stmt.Values = stmt.MultiValues[0]
	}

	return stmt, nil
}

// parseMultipleValueRows parses multiple (val1, val2), (val3, val4) rows
func parseMultipleValueRows(valuesSection string) [][]string {
	var rows [][]string
	var currentRow []string
	var currentValue strings.Builder
	depth := 0
	inString := false
	stringChar := byte(0)
	
	for i := 0; i < len(valuesSection); i++ {
		c := valuesSection[i]
		
		// Handle string literals
		if (c == '\'' || c == '"') && (i == 0 || valuesSection[i-1] != '\\') {
			if !inString {
				inString = true
				stringChar = c
			} else if c == stringChar {
				inString = false
			}
			currentValue.WriteByte(c)
			continue
		}
		
		if inString {
			currentValue.WriteByte(c)
			continue
		}
		
		switch c {
		case '(':
			if depth == 0 {
				// Start of a new row
				currentValue.Reset()
			} else {
				currentValue.WriteByte(c)
			}
			depth++
		case ')':
			depth--
			if depth == 0 {
				// End of a row
				val := strings.TrimSpace(currentValue.String())
				if val != "" {
					currentRow = append(currentRow, val)
				}
				if len(currentRow) > 0 {
					rows = append(rows, currentRow)
				}
				currentRow = nil
				currentValue.Reset()
			} else {
				currentValue.WriteByte(c)
			}
		case ',':
			if depth == 1 {
				// Separator between values in a row
				val := strings.TrimSpace(currentValue.String())
				currentRow = append(currentRow, val)
				currentValue.Reset()
			} else if depth == 0 {
				// Separator between rows, skip
			} else {
				currentValue.WriteByte(c)
			}
		default:
			if depth > 0 {
				currentValue.WriteByte(c)
			}
		}
	}
	
	return rows
}

// parseUpdate parses an UPDATE statement
func parseUpdate(stmt *SQLStatement, sql string) (*SQLStatement, error) {
	// Extract table name
	tableRegex := regexp.MustCompile(`(?i)\bUPDATE\s+(\w+)`)
	tableMatch := tableRegex.FindStringSubmatch(sql)
	if len(tableMatch) > 1 {
		stmt.Tables = append(stmt.Tables, tableMatch[1])
	}

	// Extract SET clauses
	setRegex := regexp.MustCompile(`(?i)\bSET\s+(.+?)(?:\s+WHERE|$)`)
	setMatch := setRegex.FindStringSubmatch(sql)
	if len(setMatch) > 1 {
		assignments := splitAssignments(setMatch[1])
		for _, assignment := range assignments {
			parts := strings.SplitN(assignment, "=", 2)
			if len(parts) == 2 {
				colName := strings.TrimSpace(parts[0])
				value := parseValue(strings.TrimSpace(parts[1]))
				stmt.SetClauses[colName] = value
			}
		}
	}

	// Extract WHERE clause
	whereRegex := regexp.MustCompile(`(?i)\bWHERE\s+(.+)$`)
	whereMatch := whereRegex.FindStringSubmatch(sql)
	if len(whereMatch) > 1 {
		stmt.Where = strings.TrimSpace(whereMatch[1])
	}

	return stmt, nil
}

// parseDelete parses a DELETE statement
func parseDelete(stmt *SQLStatement, sql string) (*SQLStatement, error) {
	// Extract table name
	tableRegex := regexp.MustCompile(`(?i)\bDELETE\s+FROM\s+(\w+)`)
	tableMatch := tableRegex.FindStringSubmatch(sql)
	if len(tableMatch) > 1 {
		stmt.Tables = append(stmt.Tables, tableMatch[1])
	}

	// Extract WHERE clause
	whereRegex := regexp.MustCompile(`(?i)\bWHERE\s+(.+)$`)
	whereMatch := whereRegex.FindStringSubmatch(sql)
	if len(whereMatch) > 1 {
		stmt.Where = strings.TrimSpace(whereMatch[1])
	}

	return stmt, nil
}

// parseCreateTable parses a CREATE TABLE statement
func parseCreateTable(stmt *SQLStatement, sql string) (*SQLStatement, error) {
	// Extract table name
	tableRegex := regexp.MustCompile(`(?i)\bCREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?(\w+)`)
	tableMatch := tableRegex.FindStringSubmatch(sql)
	if len(tableMatch) > 1 {
		stmt.Tables = append(stmt.Tables, tableMatch[1])
	}

	// Extract column definitions - use (?s) to make . match newlines
	columnsRegex := regexp.MustCompile(`(?s)\((.+)\)`)
	columnsMatch := columnsRegex.FindStringSubmatch(sql)
	if len(columnsMatch) > 1 {
		// Normalize whitespace (replace newlines with spaces)
		columnStr := strings.ReplaceAll(columnsMatch[1], "\n", " ")
		columnStr = strings.ReplaceAll(columnStr, "\r", " ")
		columnStr = strings.TrimSpace(columnStr)
		
		columnDefs := splitColumnDefs(columnStr)
		for _, colDef := range columnDefs {
			col := parseColumnDef(colDef)
			if col != nil {
				stmt.Columns = append(stmt.Columns, *col)
			}
		}
	}

	return stmt, nil
}

// parseAlterTable parses an ALTER TABLE statement
func parseAlterTable(stmt *SQLStatement, sql string) (*SQLStatement, error) {
	// Extract table name
	tableRegex := regexp.MustCompile(`(?i)\bALTER\s+TABLE\s+(\w+)`)
	tableMatch := tableRegex.FindStringSubmatch(sql)
	if len(tableMatch) > 1 {
		stmt.Tables = append(stmt.Tables, tableMatch[1])
	}

	// Extract ADD COLUMN
	addColRegex := regexp.MustCompile(`(?i)\bADD\s+(?:COLUMN\s+)?(\w+\s+\w+.*)`)
	addColMatch := addColRegex.FindStringSubmatch(sql)
	if len(addColMatch) > 1 {
		col := parseColumnDef(addColMatch[1])
		if col != nil {
			stmt.Columns = append(stmt.Columns, *col)
		}
	}

	return stmt, nil
}

// parseDropTable parses a DROP TABLE statement
func parseDropTable(stmt *SQLStatement, sql string) (*SQLStatement, error) {
	// Extract table name
	tableRegex := regexp.MustCompile(`(?i)\bDROP\s+TABLE\s+(?:IF\s+EXISTS\s+)?(\w+)`)
	tableMatch := tableRegex.FindStringSubmatch(sql)
	if len(tableMatch) > 1 {
		stmt.Tables = append(stmt.Tables, tableMatch[1])
	}

	return stmt, nil
}

// parseColumnDef parses a column definition string into a ColumnDef
func parseColumnDef(def string) *ColumnDef {
	def = strings.TrimSpace(def)
	if def == "" {
		return nil
	}

	// Skip table constraints
	upperDef := strings.ToUpper(def)
	if strings.HasPrefix(upperDef, "PRIMARY KEY") ||
		strings.HasPrefix(upperDef, "FOREIGN KEY") ||
		strings.HasPrefix(upperDef, "UNIQUE") ||
		strings.HasPrefix(upperDef, "CHECK") ||
		strings.HasPrefix(upperDef, "CONSTRAINT") {
		return nil
	}

	parts := strings.Fields(def)
	if len(parts) < 2 {
		return nil
	}

	col := &ColumnDef{
		Name: parts[0],
		Type: strings.ToUpper(parts[1]),
	}

	// Check for NOT NULL
	if strings.Contains(strings.ToUpper(def), "NOT NULL") {
		col.Required = true
	}

	// Check for PRIMARY KEY
	if strings.Contains(strings.ToUpper(def), "PRIMARY KEY") {
		col.PrimaryKey = true
	}

	// Check for REFERENCES (foreign key)
	refRegex := regexp.MustCompile(`(?i)\bREFERENCES\s+(\w+)`)
	refMatch := refRegex.FindStringSubmatch(def)
	if len(refMatch) > 1 {
		col.Reference = refMatch[1]
	}

	// Check for DEFAULT
	defaultRegex := regexp.MustCompile(`(?i)\bDEFAULT\s+(.+?)(?:\s+NOT|\s+PRIMARY|\s+REFERENCES|$)`)
	defaultMatch := defaultRegex.FindStringSubmatch(def)
	if len(defaultMatch) > 1 {
		col.Default = strings.TrimSpace(defaultMatch[1])
	}

	return col
}

// splitColumnDefs splits column definitions, handling nested parentheses
func splitColumnDefs(s string) []string {
	var result []string
	var current strings.Builder
	depth := 0

	for _, char := range s {
		switch char {
		case '(':
			depth++
			current.WriteRune(char)
		case ')':
			depth--
			current.WriteRune(char)
		case ',':
			if depth == 0 {
				result = append(result, strings.TrimSpace(current.String()))
				current.Reset()
			} else {
				current.WriteRune(char)
			}
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		result = append(result, strings.TrimSpace(current.String()))
	}

	return result
}

// splitValues splits values in an INSERT statement, handling quoted strings
func splitValues(s string) []string {
	var result []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, char := range s {
		switch {
		case (char == '\'' || char == '"') && !inQuote:
			inQuote = true
			quoteChar = char
			current.WriteRune(char)
		case char == quoteChar && inQuote:
			inQuote = false
			quoteChar = 0
			current.WriteRune(char)
		case char == ',' && !inQuote:
			result = append(result, strings.TrimSpace(current.String()))
			current.Reset()
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		result = append(result, strings.TrimSpace(current.String()))
	}

	return result
}

// splitAssignments splits SET assignments, handling quoted strings
func splitAssignments(s string) []string {
	var result []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)
	depth := 0

	for _, char := range s {
		switch {
		case (char == '\'' || char == '"') && !inQuote:
			inQuote = true
			quoteChar = char
			current.WriteRune(char)
		case char == quoteChar && inQuote:
			inQuote = false
			quoteChar = 0
			current.WriteRune(char)
		case char == '(' && !inQuote:
			depth++
			current.WriteRune(char)
		case char == ')' && !inQuote:
			depth--
			current.WriteRune(char)
		case char == ',' && !inQuote && depth == 0:
			result = append(result, strings.TrimSpace(current.String()))
			current.Reset()
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		result = append(result, strings.TrimSpace(current.String()))
	}

	return result
}

// parseValue converts a SQL value string to a Go value
func parseValue(s string) any {
	s = strings.TrimSpace(s)

	// Check for NULL
	if strings.ToUpper(s) == "NULL" {
		return nil
	}

	// Check for quoted string
	if (strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'")) ||
		(strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) {
		return s[1 : len(s)-1]
	}

	// Check for boolean
	upper := strings.ToUpper(s)
	if upper == "TRUE" || upper == "1" {
		return true
	}
	if upper == "FALSE" || upper == "0" {
		return false
	}

	// Return as-is (could be number, function call, etc.)
	return s
}

// IsReadOnly returns true if the statement doesn't modify data
func (s *SQLStatement) IsReadOnly() bool {
	return s.Type == StatementSelect
}

// IsDestructive returns true if the statement could delete data
func (s *SQLStatement) IsDestructive() bool {
	return s.Type == StatementDelete || s.Type == StatementDropTable
}

// RequiresConfirmation returns true if the statement should prompt for confirmation
func (s *SQLStatement) RequiresConfirmation() bool {
	return s.IsDestructive() || s.Type == StatementUpdate
}
