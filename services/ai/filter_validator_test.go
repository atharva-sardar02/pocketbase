package ai

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/stretchr/testify/assert"
)

func TestValidateFilter_ValidSimple(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.TextField{Name: "status"})

	err := ValidateFilter(`status = "active"`, collection)
	assert.NoError(t, err)
}

func TestValidateFilter_ValidComplex(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "orders")
	collection.Fields.Add(&core.TextField{Name: "status"})
	collection.Fields.Add(&core.NumberField{Name: "total"})

	err := ValidateFilter(`status = "active" && total > 100`, collection)
	assert.NoError(t, err)
}

func TestValidateFilter_UnknownField(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.TextField{Name: "status"})
	collection.Fields.Add(&core.TextField{Name: "name"})

	err := ValidateFilter(`invalid_field = "x"`, collection)
	assert.Error(t, err)
	
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "invalid_field", validationErr.Field)
	assert.Contains(t, validationErr.Message, "unknown field")
	assert.Contains(t, validationErr.Message, "status")
	assert.Contains(t, validationErr.Message, "name")
}

func TestValidateFilter_InvalidOperator(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.TextField{Name: "name"})

	// Text field cannot use > operator
	err := ValidateFilter(`name > 100`, collection)
	assert.Error(t, err)
	
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "name", validationErr.Field)
	assert.Contains(t, validationErr.Message, "comparison operators")
}

func TestValidateFilter_MalformedSyntax(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.TextField{Name: "status"})

	// Malformed filter
	err := ValidateFilter(`status = "active" &&`, collection)
	assert.Error(t, err)
	
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Contains(t, validationErr.Message, "malformed filter syntax")
}

func TestValidateFilter_DatetimeMacros(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.DateField{Name: "created"})

	// Valid datetime macro usage
	err := ValidateFilter(`created >= @now - 86400`, collection)
	assert.NoError(t, err)

	// Valid with other datetime macros
	err = ValidateFilter(`created >= @now - 604800 && created <= @now`, collection)
	assert.NoError(t, err)
}

func TestValidateFilter_RelationFields(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	usersCollection := core.NewCollection(core.CollectionTypeBase, "users")
	usersCollection.Fields.Add(&core.TextField{Name: "name"})
	app.Save(usersCollection)
	app.ReloadCachedCollections()

	postsCollection := core.NewCollection(core.CollectionTypeBase, "posts")
	postsCollection.Fields.Add(&core.RelationField{
		Name:         "author",
		CollectionId: usersCollection.Id,
	})

	// Valid relation field usage
	err := ValidateFilter(`author = "user123"`, postsCollection)
	assert.NoError(t, err)
}

func TestValidateFilter_NumberFieldOperators(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "orders")
	collection.Fields.Add(&core.NumberField{Name: "total"})

	// Valid operators for number fields
	err := ValidateFilter(`total > 100`, collection)
	assert.NoError(t, err)

	err = ValidateFilter(`total >= 100 && total <= 1000`, collection)
	assert.NoError(t, err)

	// Invalid: string operators on number field
	err = ValidateFilter(`total ~ "100"`, collection)
	assert.Error(t, err)
	
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Contains(t, validationErr.Message, "string operators")
}

func TestValidateFilter_BoolFieldOperators(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.BoolField{Name: "published"})

	// Valid operators for bool fields
	err := ValidateFilter(`published = true`, collection)
	assert.NoError(t, err)

	err = ValidateFilter(`published != false`, collection)
	assert.NoError(t, err)

	// Invalid: comparison operators on bool field
	err = ValidateFilter(`published > true`, collection)
	assert.Error(t, err)
	
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Contains(t, validationErr.Message, "only = and !=")
}

func TestValidateFilter_SelectFieldOperators(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.SelectField{
		Name:   "status",
		Values: []string{"active", "inactive", "draft"},
	})

	// Valid operators for select fields
	err := ValidateFilter(`status = "active"`, collection)
	assert.NoError(t, err)

	// Invalid: comparison operators on select field
	err = ValidateFilter(`status > "active"`, collection)
	assert.Error(t, err)
	
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Contains(t, validationErr.Message, "comparison operators")
}

func TestValidateFilter_MultiSelectArrayOperators(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.SelectField{
		Name:     "tags",
		Values:   []string{"urgent", "important", "normal"},
		MaxSelect: 5, // Multi-select
	})

	// Valid array operators for multi-select
	err := ValidateFilter(`tags ?= "urgent"`, collection)
	assert.NoError(t, err)

	err = ValidateFilter(`tags ?~ "imp"`, collection)
	assert.NoError(t, err)

	// Invalid: array operators on single-select
	singleSelectCollection := core.NewCollection(core.CollectionTypeBase, "posts2")
	singleSelectCollection.Fields.Add(&core.SelectField{
		Name:     "status",
		Values:   []string{"active", "inactive"},
		MaxSelect: 1, // Single-select
	})

	err = ValidateFilter(`status ?= "active"`, singleSelectCollection)
	assert.Error(t, err)
	
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Contains(t, validationErr.Message, "array operators")
}

func TestValidateFilter_DateFieldOperators(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.DateField{Name: "created"})

	// Valid operators for date fields
	err := ValidateFilter(`created > @now - 86400`, collection)
	assert.NoError(t, err)

	err = ValidateFilter(`created >= @now - 604800 && created <= @now`, collection)
	assert.NoError(t, err)

	// Invalid: string operators on date field
	err = ValidateFilter(`created ~ "2024"`, collection)
	assert.Error(t, err)
	
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Contains(t, validationErr.Message, "string operators")
}

func TestValidateFilter_TextFieldOperators(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.TextField{Name: "title"})

	// Valid operators for text fields
	err := ValidateFilter(`title = "test"`, collection)
	assert.NoError(t, err)

	err = ValidateFilter(`title ~ "test"`, collection)
	assert.NoError(t, err)

	err = ValidateFilter(`title !~ "spam"`, collection)
	assert.NoError(t, err)

	// Invalid: comparison operators on text field
	err = ValidateFilter(`title > "test"`, collection)
	assert.Error(t, err)
}

func TestValidateFilter_EmptyFilter(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.TextField{Name: "status"})

	// Empty filter is valid
	err := ValidateFilter("", collection)
	assert.NoError(t, err)

	err = ValidateFilter("   ", collection)
	assert.NoError(t, err)
}

func TestValidateFilter_NilCollection(t *testing.T) {
	err := ValidateFilter(`status = "active"`, nil)
	assert.Error(t, err)
	
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Contains(t, validationErr.Message, "collection is nil")
}

func TestValidateFilter_ComplexExpression(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "orders")
	collection.Fields.Add(&core.TextField{Name: "status"})
	collection.Fields.Add(&core.NumberField{Name: "total"})
	collection.Fields.Add(&core.DateField{Name: "created"})

	// Complex valid expression
	err := ValidateFilter(`(status = "pending" || status = "processing") && total > 100 && created >= @now - 604800`, collection)
	assert.NoError(t, err)
}

func TestValidateFilter_HiddenFields(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "users")
	collection.Fields.Add(&core.TextField{Name: "name"})
	passwordField := &core.PasswordField{Name: "password"}
	passwordField.SetHidden(true)
	collection.Fields.Add(passwordField)

	// Hidden fields should not be queryable
	err := ValidateFilter(`password = "test"`, collection)
	assert.Error(t, err)
	
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "password", validationErr.Field)
}

func TestExtractFieldNames(t *testing.T) {
	fields, err := ExtractFieldNames(`status = "active" && total > 100`)
	assert.NoError(t, err)
	assert.Contains(t, fields, "status")
	assert.Contains(t, fields, "total")
}

func TestExtractFieldNames_WithDatetimeMacros(t *testing.T) {
	fields, err := ExtractFieldNames(`created >= @now - 86400`)
	assert.NoError(t, err)
	assert.Contains(t, fields, "created")
	assert.NotContains(t, fields, "@now")
}

func TestExtractFieldNames_ComplexExpression(t *testing.T) {
	fields, err := ExtractFieldNames(`(status = "pending" || status = "processing") && total > 100`)
	assert.NoError(t, err)
	assert.Contains(t, fields, "status")
	assert.Contains(t, fields, "total")
	
	// Should not contain duplicates
	statusCount := 0
	for _, f := range fields {
		if f == "status" {
			statusCount++
		}
	}
	assert.Equal(t, 1, statusCount, "status should appear only once")
}

