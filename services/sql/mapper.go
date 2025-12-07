package sql

import (
	"regexp"
	"strings"

	"github.com/pocketbase/pocketbase/core"
)

// PocketBaseFieldType represents a PocketBase field type
type PocketBaseFieldType string

const (
	FieldTypeText     PocketBaseFieldType = "text"
	FieldTypeNumber   PocketBaseFieldType = "number"
	FieldTypeBool     PocketBaseFieldType = "bool"
	FieldTypeEmail    PocketBaseFieldType = "email"
	FieldTypeURL      PocketBaseFieldType = "url"
	FieldTypeDate     PocketBaseFieldType = "date"
	FieldTypeSelect   PocketBaseFieldType = "select"
	FieldTypeRelation PocketBaseFieldType = "relation"
	FieldTypeFile     PocketBaseFieldType = "file"
	FieldTypeJSON     PocketBaseFieldType = "json"
	FieldTypeEditor   PocketBaseFieldType = "editor"
)

// MapSQLType maps a SQL type to a PocketBase field type
func MapSQLType(sqlType string) PocketBaseFieldType {
	sqlType = strings.ToUpper(strings.TrimSpace(sqlType))

	// Handle type with size specifier (e.g., VARCHAR(255))
	if idx := strings.Index(sqlType, "("); idx != -1 {
		baseType := sqlType[:idx]
		switch baseType {
		case "VARCHAR", "CHAR", "NVARCHAR", "NCHAR":
			return FieldTypeText
		case "DECIMAL", "NUMERIC":
			return FieldTypeNumber
		}
	}

	switch sqlType {
	// Text types
	case "TEXT", "VARCHAR", "CHAR", "NVARCHAR", "NCHAR", "CLOB", "STRING":
		return FieldTypeText

	// Number types
	case "INTEGER", "INT", "SMALLINT", "BIGINT", "TINYINT":
		return FieldTypeNumber
	case "REAL", "FLOAT", "DOUBLE", "DECIMAL", "NUMERIC", "NUMBER":
		return FieldTypeNumber

	// Boolean types
	case "BOOLEAN", "BOOL", "BIT":
		return FieldTypeBool

	// Date/Time types
	case "DATE", "DATETIME", "TIMESTAMP", "TIME":
		return FieldTypeDate

	// Special text types
	case "EMAIL":
		return FieldTypeEmail
	case "URL", "URI":
		return FieldTypeURL

	// Structured types
	case "JSON", "JSONB":
		return FieldTypeJSON
	case "BLOB", "BINARY", "VARBINARY", "FILE":
		return FieldTypeFile

	// Enum/Select type
	case "ENUM":
		return FieldTypeSelect

	// Foreign key (detected separately)
	case "FOREIGN KEY", "REFERENCES":
		return FieldTypeRelation

	default:
		// Default to text for unknown types
		return FieldTypeText
	}
}

// MapColumnToField converts a SQL column definition to a PocketBase field
func MapColumnToField(col ColumnDef) core.Field {
	fieldType := MapSQLType(col.Type)

	switch fieldType {
	case FieldTypeText:
		return &core.TextField{
			Name:     col.Name,
			Required: col.Required,
		}

	case FieldTypeNumber:
		return &core.NumberField{
			Name:     col.Name,
			Required: col.Required,
		}

	case FieldTypeBool:
		return &core.BoolField{
			Name:     col.Name,
			Required: col.Required,
		}

	case FieldTypeEmail:
		return &core.EmailField{
			Name:     col.Name,
			Required: col.Required,
		}

	case FieldTypeURL:
		return &core.URLField{
			Name:     col.Name,
			Required: col.Required,
		}

	case FieldTypeDate:
		return &core.DateField{
			Name:     col.Name,
			Required: col.Required,
		}

	case FieldTypeJSON:
		return &core.JSONField{
			Name:     col.Name,
			Required: col.Required,
		}

	case FieldTypeFile:
		return &core.FileField{
			Name:     col.Name,
			Required: col.Required,
		}

	case FieldTypeSelect:
		return &core.SelectField{
			Name:     col.Name,
			Required: col.Required,
			Values:   col.Options,
		}

	case FieldTypeRelation:
		return &core.RelationField{
			Name:         col.Name,
			Required:     col.Required,
			CollectionId: col.Reference, // Will need to be resolved to actual ID
		}

	default:
		return &core.TextField{
			Name:     col.Name,
			Required: col.Required,
		}
	}
}

// MapPocketBaseTypeToSQL converts a PocketBase field type to SQL type
func MapPocketBaseTypeToSQL(fieldType string) string {
	switch fieldType {
	case "text":
		return "TEXT"
	case "number":
		return "REAL"
	case "bool":
		return "BOOLEAN"
	case "email":
		return "TEXT"
	case "url":
		return "TEXT"
	case "date", "autodate":
		return "DATETIME"
	case "select":
		return "TEXT"
	case "relation":
		return "TEXT" // Stores ID as text
	case "file":
		return "TEXT" // Stores filename as text
	case "json":
		return "JSON"
	case "editor":
		return "TEXT"
	case "geopoint":
		return "TEXT" // JSON format: {"lat": x, "lon": y}
	default:
		return "TEXT"
	}
}

// ParseEnumOptions extracts options from an ENUM type definition
// e.g., "ENUM('active', 'inactive', 'pending')" returns ["active", "inactive", "pending"]
func ParseEnumOptions(typeDef string) []string {
	// Match ENUM('opt1', 'opt2', ...) or CHECK(field IN ('opt1', 'opt2', ...))
	enumRegex := regexp.MustCompile(`(?i)(?:ENUM|IN)\s*\(([^)]+)\)`)
	match := enumRegex.FindStringSubmatch(typeDef)
	if len(match) < 2 {
		return nil
	}

	optionsStr := match[1]
	var options []string

	// Parse quoted options
	optRegex := regexp.MustCompile(`'([^']*)'|"([^"]*)"`)
	matches := optRegex.FindAllStringSubmatch(optionsStr, -1)
	for _, m := range matches {
		if m[1] != "" {
			options = append(options, m[1])
		} else if m[2] != "" {
			options = append(options, m[2])
		}
	}

	return options
}

// DetectFieldTypeFromConstraints analyzes constraints to better determine field type
func DetectFieldTypeFromConstraints(col ColumnDef, constraints string) PocketBaseFieldType {
	constraints = strings.ToUpper(constraints)

	// Check for email pattern
	if strings.Contains(constraints, "EMAIL") ||
		strings.Contains(constraints, "LIKE '%@%'") {
		return FieldTypeEmail
	}

	// Check for URL pattern
	if strings.Contains(constraints, "URL") ||
		strings.Contains(constraints, "LIKE 'HTTP%'") {
		return FieldTypeURL
	}

	// Check for ENUM/CHECK constraints (select type)
	if strings.Contains(constraints, "IN (") ||
		strings.Contains(constraints, "ENUM") {
		return FieldTypeSelect
	}

	// Check for foreign key
	if strings.Contains(constraints, "REFERENCES") ||
		strings.Contains(constraints, "FOREIGN KEY") {
		return FieldTypeRelation
	}

	// Fall back to basic type mapping
	return MapSQLType(col.Type)
}

// ColumnDefToSchemaField converts ColumnDef to a structured field definition
// that can be used when creating a PocketBase collection
type SchemaField struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Required   bool     `json:"required"`
	Options    []string `json:"options,omitempty"`    // For select fields
	Collection string   `json:"collection,omitempty"` // For relation fields
}

// ConvertColumnsToSchemaFields converts SQL column definitions to PocketBase schema fields
func ConvertColumnsToSchemaFields(columns []ColumnDef) []SchemaField {
	var fields []SchemaField

	for _, col := range columns {
		// Skip id, created, updated - PocketBase adds these automatically
		lowerName := strings.ToLower(col.Name)
		if lowerName == "id" || lowerName == "created" || lowerName == "updated" {
			continue
		}

		fieldType := MapSQLType(col.Type)

		field := SchemaField{
			Name:     col.Name,
			Type:     string(fieldType),
			Required: col.Required,
		}

		// Handle select options
		if fieldType == FieldTypeSelect && len(col.Options) > 0 {
			field.Options = col.Options
		}

		// Handle relation reference
		if col.Reference != "" {
			field.Type = string(FieldTypeRelation)
			field.Collection = col.Reference
		}

		fields = append(fields, field)
	}

	return fields
}

