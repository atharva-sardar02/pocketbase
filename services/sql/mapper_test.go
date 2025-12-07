package sql

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/stretchr/testify/assert"
)

func TestMapSQLType_Text(t *testing.T) {
	testCases := []struct {
		sqlType  string
		expected PocketBaseFieldType
	}{
		{"TEXT", FieldTypeText},
		{"VARCHAR", FieldTypeText},
		{"VARCHAR(255)", FieldTypeText},
		{"CHAR", FieldTypeText},
		{"NVARCHAR", FieldTypeText},
		{"STRING", FieldTypeText},
		{"CLOB", FieldTypeText},
	}

	for _, tc := range testCases {
		t.Run(tc.sqlType, func(t *testing.T) {
			result := MapSQLType(tc.sqlType)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMapSQLType_Number(t *testing.T) {
	testCases := []struct {
		sqlType  string
		expected PocketBaseFieldType
	}{
		{"INTEGER", FieldTypeNumber},
		{"INT", FieldTypeNumber},
		{"SMALLINT", FieldTypeNumber},
		{"BIGINT", FieldTypeNumber},
		{"REAL", FieldTypeNumber},
		{"FLOAT", FieldTypeNumber},
		{"DOUBLE", FieldTypeNumber},
		{"DECIMAL", FieldTypeNumber},
		{"NUMERIC", FieldTypeNumber},
		{"DECIMAL(10,2)", FieldTypeNumber},
	}

	for _, tc := range testCases {
		t.Run(tc.sqlType, func(t *testing.T) {
			result := MapSQLType(tc.sqlType)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMapSQLType_Bool(t *testing.T) {
	testCases := []string{"BOOLEAN", "BOOL", "BIT"}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			result := MapSQLType(tc)
			assert.Equal(t, FieldTypeBool, result)
		})
	}
}

func TestMapSQLType_Date(t *testing.T) {
	testCases := []string{"DATE", "DATETIME", "TIMESTAMP", "TIME"}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			result := MapSQLType(tc)
			assert.Equal(t, FieldTypeDate, result)
		})
	}
}

func TestMapSQLType_Special(t *testing.T) {
	assert.Equal(t, FieldTypeEmail, MapSQLType("EMAIL"))
	assert.Equal(t, FieldTypeURL, MapSQLType("URL"))
	assert.Equal(t, FieldTypeJSON, MapSQLType("JSON"))
	assert.Equal(t, FieldTypeFile, MapSQLType("BLOB"))
	assert.Equal(t, FieldTypeSelect, MapSQLType("ENUM"))
}

func TestMapSQLType_CaseInsensitive(t *testing.T) {
	assert.Equal(t, MapSQLType("text"), MapSQLType("TEXT"))
	assert.Equal(t, MapSQLType("Integer"), MapSQLType("INTEGER"))
}

func TestMapSQLType_Unknown(t *testing.T) {
	// Unknown types should default to text
	assert.Equal(t, FieldTypeText, MapSQLType("UNKNOWN_TYPE"))
}

func TestMapColumnToField_Text(t *testing.T) {
	col := ColumnDef{
		Name:     "title",
		Type:     "TEXT",
		Required: true,
	}

	field := MapColumnToField(col)

	textField, ok := field.(*core.TextField)
	assert.True(t, ok)
	assert.Equal(t, "title", textField.Name)
	assert.True(t, textField.Required)
}

func TestMapColumnToField_Number(t *testing.T) {
	col := ColumnDef{
		Name:     "price",
		Type:     "REAL",
		Required: false,
	}

	field := MapColumnToField(col)

	numField, ok := field.(*core.NumberField)
	assert.True(t, ok)
	assert.Equal(t, "price", numField.Name)
	assert.False(t, numField.Required)
}

func TestMapColumnToField_Bool(t *testing.T) {
	col := ColumnDef{
		Name: "is_active",
		Type: "BOOLEAN",
	}

	field := MapColumnToField(col)

	boolField, ok := field.(*core.BoolField)
	assert.True(t, ok)
	assert.Equal(t, "is_active", boolField.Name)
}

func TestMapColumnToField_Date(t *testing.T) {
	col := ColumnDef{
		Name: "created_at",
		Type: "DATETIME",
	}

	field := MapColumnToField(col)

	dateField, ok := field.(*core.DateField)
	assert.True(t, ok)
	assert.Equal(t, "created_at", dateField.Name)
}

func TestMapColumnToField_Select(t *testing.T) {
	col := ColumnDef{
		Name:    "status",
		Type:    "ENUM",
		Options: []string{"active", "inactive", "pending"},
	}

	field := MapColumnToField(col)

	selectField, ok := field.(*core.SelectField)
	assert.True(t, ok)
	assert.Equal(t, "status", selectField.Name)
	assert.Equal(t, []string{"active", "inactive", "pending"}, selectField.Values)
}

func TestMapColumnToField_Relation(t *testing.T) {
	col := ColumnDef{
		Name:      "customer",
		Type:      "TEXT",
		Reference: "customers",
	}

	field := MapColumnToField(col)

	relField, ok := field.(*core.RelationField)
	assert.True(t, ok)
	assert.Equal(t, "customer", relField.Name)
	assert.Equal(t, "customers", relField.CollectionId)
}

func TestMapPocketBaseTypeToSQL(t *testing.T) {
	testCases := []struct {
		pbType   string
		expected string
	}{
		{"text", "TEXT"},
		{"number", "REAL"},
		{"bool", "BOOLEAN"},
		{"date", "DATETIME"},
		{"select", "TEXT"},
		{"relation", "TEXT"},
		{"json", "JSON"},
		{"email", "TEXT"},
		{"url", "TEXT"},
	}

	for _, tc := range testCases {
		t.Run(tc.pbType, func(t *testing.T) {
			result := MapPocketBaseTypeToSQL(tc.pbType)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestParseEnumOptions(t *testing.T) {
	// ENUM syntax
	options := ParseEnumOptions("ENUM('active', 'inactive', 'pending')")
	assert.Equal(t, []string{"active", "inactive", "pending"}, options)

	// CHECK IN syntax
	options = ParseEnumOptions("CHECK(status IN ('active', 'inactive'))")
	assert.Equal(t, []string{"active", "inactive"}, options)

	// With double quotes
	options = ParseEnumOptions(`ENUM("yes", "no")`)
	assert.Equal(t, []string{"yes", "no"}, options)

	// No options
	options = ParseEnumOptions("TEXT")
	assert.Nil(t, options)
}

func TestDetectFieldTypeFromConstraints(t *testing.T) {
	col := ColumnDef{Name: "field", Type: "TEXT"}

	// Email detection
	assert.Equal(t, FieldTypeEmail, DetectFieldTypeFromConstraints(col, "email format"))

	// URL detection
	assert.Equal(t, FieldTypeURL, DetectFieldTypeFromConstraints(col, "LIKE 'http%'"))

	// Select detection
	assert.Equal(t, FieldTypeSelect, DetectFieldTypeFromConstraints(col, "IN ('a', 'b')"))

	// Relation detection
	assert.Equal(t, FieldTypeRelation, DetectFieldTypeFromConstraints(col, "REFERENCES users"))

	// Default
	assert.Equal(t, FieldTypeText, DetectFieldTypeFromConstraints(col, ""))
}

func TestConvertColumnsToSchemaFields(t *testing.T) {
	columns := []ColumnDef{
		{Name: "id", Type: "TEXT"},      // Should be skipped
		{Name: "created", Type: "DATETIME"}, // Should be skipped
		{Name: "name", Type: "TEXT", Required: true},
		{Name: "age", Type: "INTEGER"},
		{Name: "status", Type: "ENUM", Options: []string{"active", "inactive"}},
		{Name: "customer", Type: "TEXT", Reference: "customers"},
	}

	fields := ConvertColumnsToSchemaFields(columns)

	// Should skip id and created
	assert.Len(t, fields, 4)

	// Check name field
	assert.Equal(t, "name", fields[0].Name)
	assert.Equal(t, "text", fields[0].Type)
	assert.True(t, fields[0].Required)

	// Check age field
	assert.Equal(t, "age", fields[1].Name)
	assert.Equal(t, "number", fields[1].Type)

	// Check status field (select)
	assert.Equal(t, "status", fields[2].Name)
	assert.Equal(t, "select", fields[2].Type)
	assert.Equal(t, []string{"active", "inactive"}, fields[2].Options)

	// Check customer field (relation)
	assert.Equal(t, "customer", fields[3].Name)
	assert.Equal(t, "relation", fields[3].Type)
	assert.Equal(t, "customers", fields[3].Collection)
}

