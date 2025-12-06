package ai

import (
	"regexp"
	"strings"
)

// TokenizeFilter extracts field names and operators from a filter expression.
// This is a basic tokenizer for validation purposes.
func TokenizeFilter(filter string) ([]string, []string, error) {
	// Remove whitespace for easier parsing
	filter = strings.TrimSpace(filter)
	if filter == "" {
		return nil, nil, nil
	}

	var fields []string
	var operators []string

	// Remove quoted strings first (they contain values, not field names)
	quotedStringPattern := regexp.MustCompile(`"[^"]*"`)
	filterWithoutQuotes := quotedStringPattern.ReplaceAllString(filter, "")

	// Regex to match field names (alphanumeric + underscore, not starting with @)
	// Field names appear before operators, so we look for patterns like "fieldName operator"
	fieldPattern := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\s*(>=|<=|!=|!~|\?=|\?~|&&|\|\||[=<>~])`)
	
	// Regex to match operators
	operatorPattern := regexp.MustCompile(`(>=|<=|!=|!~|\?=|\?~|&&|\|\||[=<>~])`)

	// Extract field names that appear before operators
	fieldMatches := fieldPattern.FindAllStringSubmatch(filterWithoutQuotes, -1)
	for _, match := range fieldMatches {
		if len(match) > 1 {
			fieldName := match[1]
			// Skip datetime macros
			if strings.HasPrefix(fieldName, "@") {
				continue
			}
			// Skip boolean operators
			if fieldName == "AND" || fieldName == "OR" {
				continue
			}
			// Skip common keywords
			if fieldName == "null" || fieldName == "true" || fieldName == "false" {
				continue
			}
			fields = append(fields, fieldName)
		}
	}

	// Extract operators
	operatorMatches := operatorPattern.FindAllString(filter, -1)
	operators = append(operators, operatorMatches...)

	return fields, operators, nil
}

// ExtractFieldNames extracts all field names from a filter expression.
func ExtractFieldNames(filter string) ([]string, error) {
	fields, _, err := TokenizeFilter(filter)
	if err != nil {
		return nil, err
	}

	// Remove duplicates
	seen := make(map[string]bool)
	var uniqueFields []string
	for _, field := range fields {
		if !seen[field] {
			seen[field] = true
			uniqueFields = append(uniqueFields, field)
		}
	}

	return uniqueFields, nil
}

