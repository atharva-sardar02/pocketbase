package ai

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ganigeorgiev/fexpr"
	"github.com/pocketbase/pocketbase/core"
)

// ValidationError represents a filter validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("field %q: %s", e.Field, e.Message)
	}
	return e.Message
}

// ValidateFilter validates a filter expression against a collection schema.
func ValidateFilter(filter string, collection *core.Collection) error {
	if collection == nil {
		return &ValidationError{Message: "collection is nil"}
	}

	filter = strings.TrimSpace(filter)
	if filter == "" {
		return nil // Empty filter is valid
	}

	// Check if filter contains datetime macros
	// If it does, we'll be more lenient with parsing since fexpr may not handle
	// datetime macro arithmetic directly
	hasDatetimeMacros := regexp.MustCompile(`@\w+`).MatchString(filter)

	// Try to parse the filter using fexpr to catch syntax errors
	// For datetime macros, we'll replace them with placeholders first
	var parseFilter string
	if hasDatetimeMacros {
		parseFilter = preprocessDatetimeMacros(filter)
	} else {
		parseFilter = filter
	}

	_, err := fexpr.Parse(parseFilter)
	if err != nil {
		// If parsing fails and we have datetime macros, try a simpler validation
		// by just checking if basic syntax is present
		if hasDatetimeMacros {
			// For datetime macros, just validate that the structure is reasonable
			// (has operators, field names, etc.) without strict fexpr parsing
			if !isBasicValidStructure(filter) {
				return &ValidationError{
					Message: fmt.Sprintf("malformed filter syntax: %v", err),
				}
			}
		} else {
			return &ValidationError{
				Message: fmt.Sprintf("malformed filter syntax: %v", err),
			}
		}
	}

	// Extract field names from the filter
	fieldNames, err := ExtractFieldNames(filter)
	if err != nil {
		return &ValidationError{
			Message: fmt.Sprintf("failed to extract field names: %v", err),
		}
	}

	// Build a map of available field names for quick lookup
	availableFields := make(map[string]core.Field)
	for _, field := range collection.Fields {
		// Skip hidden fields (they're not queryable)
		if field.GetHidden() {
			continue
		}
		availableFields[field.GetName()] = field
	}

	// Validate each field exists in the collection
	for _, fieldName := range fieldNames {
		field, exists := availableFields[fieldName]
		if !exists {
			// Build list of available fields for error message
			var availableNames []string
			for name := range availableFields {
				availableNames = append(availableNames, name)
			}
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("unknown field. Available fields: %s", strings.Join(availableNames, ", ")),
			}
		}

		// Validate operator compatibility with field type
		if err := validateFieldOperators(filter, fieldName, field); err != nil {
			return err
		}
	}

	return nil
}

// validateFieldOperators checks if the operators used with a field are compatible with its type.
func validateFieldOperators(filter string, fieldName string, field core.Field) error {
	fieldType := field.Type()

	// Extract operators used with this field
	// This is a simplified check - in a real implementation, we'd parse the AST
	// For now, we'll check if incompatible operators appear near the field name

	// Patterns that indicate comparison operators (>, <, >=, <=)
	comparisonPattern := fmt.Sprintf(`%s\s*(>=|<=|>|<)`, regexp.QuoteMeta(fieldName))
	comparisonRegex := regexp.MustCompile(comparisonPattern)

	// Patterns that indicate string operators (~, !~)
	stringPattern := fmt.Sprintf(`%s\s*(~|!~)`, regexp.QuoteMeta(fieldName))
	stringRegex := regexp.MustCompile(stringPattern)

	// Patterns that indicate array operators (?=, ?~)
	arrayPattern := fmt.Sprintf(`%s\s*(\?=|\?~)`, regexp.QuoteMeta(fieldName))
	arrayRegex := regexp.MustCompile(arrayPattern)

	// Check for incompatible operators based on field type
	switch fieldType {
	case core.FieldTypeText, core.FieldTypeEmail, core.FieldTypeURL, core.FieldTypeEditor:
		// Text fields can use =, !=, ~, !~, but not >, <, >=, <=
		if comparisonRegex.MatchString(filter) {
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("invalid operator for field type '%s': comparison operators (>, <, >=, <=) cannot be used with text fields", fieldType),
			}
		}
		// Array operators are only for select/relation fields with multiple values
		if arrayRegex.MatchString(filter) {
			// Check if it's actually a multi-select field
			if selectField, ok := field.(*core.SelectField); ok {
				if selectField.MaxSelect <= 1 {
					return &ValidationError{
						Field:   fieldName,
						Message: fmt.Sprintf("invalid operator for field type '%s': array operators (?=, ?~) can only be used with multi-select fields", fieldType),
					}
				}
			} else {
				return &ValidationError{
					Field:   fieldName,
					Message: fmt.Sprintf("invalid operator for field type '%s': array operators (?=, ?~) can only be used with select or relation fields", fieldType),
				}
			}
		}

	case core.FieldTypeNumber:
		// Number fields can use =, !=, >, <, >=, <=, but not ~, !~
		if stringRegex.MatchString(filter) {
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("invalid operator for field type '%s': string operators (~, !~) cannot be used with number fields", fieldType),
			}
		}

	case core.FieldTypeBool:
		// Bool fields can only use =, !=
		if comparisonRegex.MatchString(filter) || stringRegex.MatchString(filter) {
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("invalid operator for field type '%s': only = and != operators are allowed", fieldType),
			}
		}

	case core.FieldTypeDate, core.FieldTypeAutodate:
		// Date fields can use =, !=, >, <, >=, <=, but not ~, !~
		if stringRegex.MatchString(filter) {
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("invalid operator for field type '%s': string operators (~, !~) cannot be used with date fields", fieldType),
			}
		}

	case core.FieldTypeSelect:
		// Select fields can use =, !=, ?=, ?~
		// Array operators are only valid for multi-select
		if selectField, ok := field.(*core.SelectField); ok {
			if arrayRegex.MatchString(filter) && selectField.MaxSelect <= 1 {
				return &ValidationError{
					Field:   fieldName,
					Message: fmt.Sprintf("invalid operator for field type '%s': array operators (?=, ?~) can only be used with multi-select fields (maxSelect > 1)", fieldType),
				}
			}
		}
		// Comparison operators are not valid for select
		if comparisonRegex.MatchString(filter) {
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("invalid operator for field type '%s': comparison operators (>, <, >=, <=) cannot be used with select fields", fieldType),
			}
		}

	case core.FieldTypeRelation:
		// Relation fields can use =, !=, ?=, ?~
		// Array operators are only valid for multi-relation
		if relationField, ok := field.(*core.RelationField); ok {
			if arrayRegex.MatchString(filter) && relationField.MaxSelect <= 1 {
				return &ValidationError{
					Field:   fieldName,
					Message: fmt.Sprintf("invalid operator for field type '%s': array operators (?=, ?~) can only be used with multi-relation fields (maxSelect > 1)", fieldType),
				}
			}
		}
		// Comparison operators are not valid for relation
		if comparisonRegex.MatchString(filter) {
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("invalid operator for field type '%s': comparison operators (>, <, >=, <=) cannot be used with relation fields", fieldType),
			}
		}

	case core.FieldTypeFile, core.FieldTypeJSON, core.FieldTypeGeoPoint:
		// These field types have limited operator support
		if comparisonRegex.MatchString(filter) || stringRegex.MatchString(filter) {
			return &ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("invalid operator for field type '%s': only = and != operators are allowed", fieldType),
			}
		}
	}

	return nil
}

// preprocessDatetimeMacros replaces datetime macro arithmetic with placeholders
// to allow fexpr parsing. Patterns like "@now - 86400" become a simple identifier.
func preprocessDatetimeMacros(filter string) string {
	// Pattern: @macro - number or @macro + number
	// Replace with a simple identifier that fexpr can parse
	datetimeArithmeticPattern := regexp.MustCompile(`@(\w+)\s*([+\-])\s*(\d+)`)
	return datetimeArithmeticPattern.ReplaceAllString(filter, "datetime_macro_$1")
}

// isBasicValidStructure performs a basic syntax check without strict fexpr parsing.
// This is used as a fallback when datetime macros are present.
func isBasicValidStructure(filter string) bool {
	// Check for basic operators
	hasOperators := regexp.MustCompile(`(>=|<=|!=|!~|\?=|\?~|&&|\|\||[=<>~])`).MatchString(filter)
	if !hasOperators {
		return false
	}
	// Check for balanced parentheses (basic check)
	openCount := strings.Count(filter, "(")
	closeCount := strings.Count(filter, ")")
	return openCount == closeCount
}

