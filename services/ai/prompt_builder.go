package ai

import (
	"fmt"
	"strings"
)

// BuildSystemPrompt constructs the full system prompt by injecting the schema
// into the system prompt template. Used for single-collection filter generation.
func BuildSystemPrompt(schema string) string {
	template := SystemPromptTemplate
	return strings.Replace(template, "{schema}", schema, 1)
}

// BuildUserPrompt wraps the user query in a standardized format.
func BuildUserPrompt(query string) string {
	return fmt.Sprintf("USER QUERY: %s", query)
}

// BuildDualOutputPrompt constructs a prompt for generating both Filter AND SQL.
// The schema should be the full database schema from ExtractAllSchemas or ExtractSchemaForCollection.
func BuildDualOutputPrompt(schema string) string {
	template := DualOutputPromptTemplate
	return strings.Replace(template, "{schema}", schema, 1)
}

// BuildSQLTerminalPrompt constructs a prompt for direct SQL generation.
// Used by the SQL Terminal feature.
func BuildSQLTerminalPrompt(schema string) string {
	template := SQLTerminalPromptTemplate
	return strings.Replace(template, "{schema}", schema, 1)
}

// PromptMode indicates which type of query generation to use
type PromptMode string

const (
	// PromptModeFilter generates only PocketBase filter (V1 behavior)
	PromptModeFilter PromptMode = "filter"
	// PromptModeDual generates both filter AND SQL (V2 enhanced)
	PromptModeDual PromptMode = "dual"
	// PromptModeSQL generates only SQL (SQL Terminal)
	PromptModeSQL PromptMode = "sql"
)

// BuildPromptForMode constructs the appropriate system prompt based on the mode
func BuildPromptForMode(schema string, mode PromptMode) string {
	switch mode {
	case PromptModeDual:
		return BuildDualOutputPrompt(schema)
	case PromptModeSQL:
		return BuildSQLTerminalPrompt(schema)
	default:
		return BuildSystemPrompt(schema)
	}
}

