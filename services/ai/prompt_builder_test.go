package ai

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/stretchr/testify/assert"
)

func TestBuildSystemPrompt_IncludesSchema(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.TextField{Name: "title"})
	collection.Fields.Add(&core.TextField{Name: "status"})

	schema := ExtractSchema(app, collection)
	prompt := BuildSystemPrompt(schema)

	assert.Contains(t, prompt, "Collection: posts")
	assert.Contains(t, prompt, "title (text)")
	assert.Contains(t, prompt, "status (text)")
}

func TestBuildSystemPrompt_IncludesSyntaxRules(t *testing.T) {
	schema := "Collection: test\nFields:\n  name (text)"
	prompt := BuildSystemPrompt(schema)

	assert.Contains(t, prompt, "Use = for exact match")
	assert.Contains(t, prompt, "Use != for not equals")
	assert.Contains(t, prompt, "Use > < >= <=")
	assert.Contains(t, prompt, "Use ~ for contains")
	assert.Contains(t, prompt, "Use && for AND")
	assert.Contains(t, prompt, "Use || for OR")
	assert.Contains(t, prompt, "Use () for grouping")
}

func TestBuildSystemPrompt_IncludesExamples(t *testing.T) {
	schema := "Collection: test\nFields:\n  name (text)"
	prompt := BuildSystemPrompt(schema)

	assert.Contains(t, prompt, "User: \"active users\"")
	assert.Contains(t, prompt, "Filter: status = \"active\"")
	assert.Contains(t, prompt, "User: \"orders over 100 dollars from this week\"")
	assert.Contains(t, prompt, "Filter: total > 100 && created >= @now - 604800")
}

func TestBuildSystemPrompt_IncludesDatetimeMacros(t *testing.T) {
	schema := "Collection: test\nFields:\n  created (date)"
	prompt := BuildSystemPrompt(schema)

	assert.Contains(t, prompt, "@now")
	assert.Contains(t, prompt, "@second")
	assert.Contains(t, prompt, "@minute")
	assert.Contains(t, prompt, "@hour")
	assert.Contains(t, prompt, "@weekday")
	assert.Contains(t, prompt, "@day")
	assert.Contains(t, prompt, "@month")
	assert.Contains(t, prompt, "@year")
	assert.Contains(t, prompt, "@now - 604800")
}

func TestBuildUserPrompt_WrapsQuery(t *testing.T) {
	query := "show me active users"
	prompt := BuildUserPrompt(query)

	assert.Equal(t, "USER QUERY: show me active users", prompt)
}

func TestBuildSystemPrompt_SchemaInjection(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "orders")
	collection.Fields.Add(&core.TextField{Name: "status"})
	collection.Fields.Add(&core.NumberField{Name: "total"})

	schema := ExtractSchema(app, collection)
	prompt := BuildSystemPrompt(schema)

	// Verify schema is injected (not the placeholder)
	assert.Contains(t, prompt, "Collection: orders")
	assert.NotContains(t, prompt, "{schema}")
}

func TestBuildSystemPrompt_ResponseFormat(t *testing.T) {
	schema := "Collection: test\nFields:\n  name (text)"
	prompt := BuildSystemPrompt(schema)

	// Verify it instructs to return only the filter
	assert.Contains(t, prompt, "Respond with ONLY the filter expression")
	assert.Contains(t, prompt, "no explanation")
	assert.Contains(t, prompt, "INVALID_QUERY")
}

func TestBuildSystemPrompt_FieldNameRules(t *testing.T) {
	schema := "Collection: test\nFields:\n  name (text)"
	prompt := BuildSystemPrompt(schema)

	// Verify field name rules
	assert.Contains(t, prompt, "Field names are case-sensitive")
	assert.Contains(t, prompt, "Do NOT include quotes around field names")
	assert.Contains(t, prompt, "Use the exact field names from the schema")
}

func TestBuildSystemPrompt_StringValueRules(t *testing.T) {
	schema := "Collection: test\nFields:\n  name (text)"
	prompt := BuildSystemPrompt(schema)

	// Verify string value rules
	assert.Contains(t, prompt, "Wrap string values in double quotes")
}

func TestBuildSystemPrompt_ArrayOperators(t *testing.T) {
	schema := "Collection: test\nFields:\n  tags (select)"
	prompt := BuildSystemPrompt(schema)

	// Verify array operators
	assert.Contains(t, prompt, "Use ?= for any equals (arrays)")
	assert.Contains(t, prompt, "Use ?~ for any contains (arrays)")
}

// ============================================================================
// V2 Dual Output and SQL Terminal Prompt Tests
// ============================================================================

func TestBuildDualOutputPrompt_IncludesSchema(t *testing.T) {
	schema := "DATABASE SCHEMA\n===============\nTABLES:\n  orders:\n    total NUMBER\n    status ENUM(pending, completed)"
	prompt := BuildDualOutputPrompt(schema)

	assert.Contains(t, prompt, "orders:")
	assert.Contains(t, prompt, "total NUMBER")
	assert.Contains(t, prompt, "DATABASE SCHEMA")
}

func TestBuildDualOutputPrompt_IncludesJSONFormat(t *testing.T) {
	schema := "DATABASE SCHEMA\nTABLES:\n  test"
	prompt := BuildDualOutputPrompt(schema)

	// Verify JSON output format instructions
	assert.Contains(t, prompt, `{"filter":`)
	assert.Contains(t, prompt, `"sql":`)
	assert.Contains(t, prompt, `"requiresSQL":`)
}

func TestBuildDualOutputPrompt_IncludesBothSyntax(t *testing.T) {
	schema := "DATABASE SCHEMA\nTABLES:\n  test"
	prompt := BuildDualOutputPrompt(schema)

	// Verify PocketBase filter syntax
	assert.Contains(t, prompt, "POCKETBASE FILTER SYNTAX")
	assert.Contains(t, prompt, "Use = for exact match")

	// Verify SQL syntax
	assert.Contains(t, prompt, "SQL SYNTAX")
	assert.Contains(t, prompt, "SQLite SQL syntax")
	assert.Contains(t, prompt, "LIKE '%value%'")
}

func TestBuildDualOutputPrompt_IncludesRequiresSQLGuidance(t *testing.T) {
	schema := "DATABASE SCHEMA\nTABLES:\n  test"
	prompt := BuildDualOutputPrompt(schema)

	// Verify guidance on when SQL is required
	assert.Contains(t, prompt, "WHEN requiresSQL IS TRUE")
	assert.Contains(t, prompt, "JOINs across multiple tables")
	assert.Contains(t, prompt, "Aggregate functions")
	assert.Contains(t, prompt, "GROUP BY")
	assert.Contains(t, prompt, "Subqueries")
}

func TestBuildDualOutputPrompt_IncludesExamples(t *testing.T) {
	schema := "DATABASE SCHEMA\nTABLES:\n  test"
	prompt := BuildDualOutputPrompt(schema)

	// Verify both simple and complex examples
	assert.Contains(t, prompt, "active users")
	assert.Contains(t, prompt, "count of orders by customer")
	assert.Contains(t, prompt, "orders with customer names")
	assert.Contains(t, prompt, "JOIN")
	assert.Contains(t, prompt, "GROUP BY")
}

func TestBuildSQLTerminalPrompt_IncludesSchema(t *testing.T) {
	schema := "DATABASE SCHEMA\n===============\nTABLES:\n  products:\n    name TEXT\n    price NUMBER"
	prompt := BuildSQLTerminalPrompt(schema)

	assert.Contains(t, prompt, "products:")
	assert.Contains(t, prompt, "name TEXT")
	assert.Contains(t, prompt, "price NUMBER")
}

func TestBuildSQLTerminalPrompt_IncludesSQLCapabilities(t *testing.T) {
	schema := "DATABASE SCHEMA\nTABLES:\n  test"
	prompt := BuildSQLTerminalPrompt(schema)

	// Verify SQL capabilities
	assert.Contains(t, prompt, "SQL CAPABILITIES")
	assert.Contains(t, prompt, "SELECT")
	assert.Contains(t, prompt, "INSERT")
	assert.Contains(t, prompt, "UPDATE")
	assert.Contains(t, prompt, "DELETE")
	assert.Contains(t, prompt, "CREATE TABLE")
	assert.Contains(t, prompt, "ALTER TABLE")
	assert.Contains(t, prompt, "DROP TABLE")
}

func TestBuildSQLTerminalPrompt_IncludesTableStructure(t *testing.T) {
	schema := "DATABASE SCHEMA\nTABLES:\n  test"
	prompt := BuildSQLTerminalPrompt(schema)

	// Verify table structure info
	assert.Contains(t, prompt, "TABLE STRUCTURE")
	assert.Contains(t, prompt, "id (TEXT PRIMARY KEY)")
	assert.Contains(t, prompt, "created (DATETIME)")
	assert.Contains(t, prompt, "updated (DATETIME)")
}

func TestBuildSQLTerminalPrompt_IncludesJoinSyntax(t *testing.T) {
	schema := "DATABASE SCHEMA\nTABLES:\n  test"
	prompt := BuildSQLTerminalPrompt(schema)

	// Verify JOIN syntax
	assert.Contains(t, prompt, "JOIN SYNTAX")
	assert.Contains(t, prompt, "table.field")
}

func TestBuildSQLTerminalPrompt_IncludesExamples(t *testing.T) {
	schema := "DATABASE SCHEMA\nTABLES:\n  test"
	prompt := BuildSQLTerminalPrompt(schema)

	// Verify various SQL examples
	assert.Contains(t, prompt, "SELECT * FROM")
	assert.Contains(t, prompt, "datetime('now'")
	assert.Contains(t, prompt, "COUNT(*)")
	assert.Contains(t, prompt, "SUM(")
	assert.Contains(t, prompt, "strftime(")
	assert.Contains(t, prompt, "CREATE TABLE")
	assert.Contains(t, prompt, "INSERT INTO")
	assert.Contains(t, prompt, "UPDATE")
	assert.Contains(t, prompt, "DELETE FROM")
}

func TestBuildPromptForMode_Filter(t *testing.T) {
	schema := "Collection: test\nFields:\n  name (text)"
	prompt := BuildPromptForMode(schema, PromptModeFilter)

	// Should use the basic filter template
	assert.Contains(t, prompt, "PocketBase filter query generator")
	assert.Contains(t, prompt, "FILTER SYNTAX RULES")
	assert.NotContains(t, prompt, "SQL SYNTAX")
}

func TestBuildPromptForMode_Dual(t *testing.T) {
	schema := "DATABASE SCHEMA\nTABLES:\n  test"
	prompt := BuildPromptForMode(schema, PromptModeDual)

	// Should use the dual output template
	assert.Contains(t, prompt, "POCKETBASE FILTER SYNTAX")
	assert.Contains(t, prompt, "SQL SYNTAX")
	assert.Contains(t, prompt, `"requiresSQL"`)
}

func TestBuildPromptForMode_SQL(t *testing.T) {
	schema := "DATABASE SCHEMA\nTABLES:\n  test"
	prompt := BuildPromptForMode(schema, PromptModeSQL)

	// Should use the SQL terminal template
	assert.Contains(t, prompt, "SQL query generator")
	assert.Contains(t, prompt, "SQL CAPABILITIES")
	assert.Contains(t, prompt, "CREATE TABLE")
}

func TestBuildPromptForMode_DefaultToFilter(t *testing.T) {
	schema := "Collection: test\nFields:\n  name (text)"
	prompt := BuildPromptForMode(schema, "unknown")

	// Should default to filter mode
	assert.Contains(t, prompt, "PocketBase filter query generator")
}

