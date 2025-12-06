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

