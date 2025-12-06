package ai

import (
	"fmt"
	"strings"
)

// BuildSystemPrompt constructs the full system prompt by injecting the schema
// into the system prompt template.
func BuildSystemPrompt(schema string) string {
	template := SystemPromptTemplate
	return strings.Replace(template, "{schema}", schema, 1)
}

// BuildUserPrompt wraps the user query in a standardized format.
func BuildUserPrompt(query string) string {
	return fmt.Sprintf("USER QUERY: %s", query)
}

