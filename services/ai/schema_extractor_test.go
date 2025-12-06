package ai

import (
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tests"
	"github.com/stretchr/testify/assert"
)

func TestExtractSchema_TextFields(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.TextField{Name: "title"})
	collection.Fields.Add(&core.TextField{Name: "description"})

	schema := ExtractSchema(app, collection)

	assert.Contains(t, schema, "Collection: posts")
	assert.Contains(t, schema, "title (text)")
	assert.Contains(t, schema, "description (text)")
}

func TestExtractSchema_NumberFields(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "orders")
	collection.Fields.Add(&core.NumberField{Name: "total"})
	collection.Fields.Add(&core.NumberField{Name: "quantity"})

	schema := ExtractSchema(app, collection)

	assert.Contains(t, schema, "Collection: orders")
	assert.Contains(t, schema, "total (number)")
	assert.Contains(t, schema, "quantity (number)")
}

func TestExtractSchema_SelectFields(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "posts")
	collection.Fields.Add(&core.SelectField{
		Name:   "status",
		Values: []string{"active", "inactive", "draft"},
	})
	collection.Fields.Add(&core.SelectField{
		Name:   "priority",
		Values: []string{"high", "medium", "low"},
	})

	schema := ExtractSchema(app, collection)

	assert.Contains(t, schema, "status (select: active|inactive|draft)")
	assert.Contains(t, schema, "priority (select: high|medium|low)")
}

func TestExtractSchema_RelationFields(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create a users collection first
	usersCollection := core.NewCollection(core.CollectionTypeBase, "users")
	usersCollection.Fields.Add(&core.TextField{Name: "name"})
	app.Save(usersCollection)

	// Reload cache to ensure the collection is available
	app.ReloadCachedCollections()

	// Create posts collection with relation to users
	postsCollection := core.NewCollection(core.CollectionTypeBase, "posts")
	postsCollection.Fields.Add(&core.RelationField{
		Name:         "author",
		CollectionId: usersCollection.Id,
	})

	schema := ExtractSchema(app, postsCollection)

	// Should contain relation field - either with collection name (if resolved) or ID (if not)
	assert.Contains(t, schema, "author (relation →")
	// Should contain either the collection name or ID
	assert.True(t, 
		strings.Contains(schema, "author (relation → users") || 
		strings.Contains(schema, "author (relation → "+usersCollection.Id),
		"Schema should contain relation to users collection name or ID")
}

func TestExtractSchema_AllFieldTypes(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "test")
	collection.Fields.Add(&core.TextField{Name: "title"})
	collection.Fields.Add(&core.NumberField{Name: "count"})
	collection.Fields.Add(&core.BoolField{Name: "published"})
	collection.Fields.Add(&core.EmailField{Name: "email"})
	collection.Fields.Add(&core.URLField{Name: "website"})
	collection.Fields.Add(&core.DateField{Name: "created"})
	collection.Fields.Add(&core.SelectField{Name: "status", Values: []string{"active", "inactive"}})
	collection.Fields.Add(&core.FileField{Name: "attachment"})
	collection.Fields.Add(&core.JSONField{Name: "metadata"})
	collection.Fields.Add(&core.EditorField{Name: "content"})
	collection.Fields.Add(&core.GeoPointField{Name: "location"})

	schema := ExtractSchema(app, collection)

	assert.Contains(t, schema, "title (text)")
	assert.Contains(t, schema, "count (number)")
	assert.Contains(t, schema, "published (bool)")
	assert.Contains(t, schema, "email (email)")
	assert.Contains(t, schema, "website (url)")
	assert.Contains(t, schema, "created (date)")
	assert.Contains(t, schema, "status (select: active|inactive)")
	assert.Contains(t, schema, "attachment (file)")
	assert.Contains(t, schema, "metadata (json)")
	assert.Contains(t, schema, "content (editor)")
	assert.Contains(t, schema, "location (geopoint)")
}

func TestExtractSchema_EmptyCollection(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "empty")
	// Remove default id field for truly empty collection
	collection.Fields.RemoveByName("id")

	schema := ExtractSchema(app, collection)

	assert.Contains(t, schema, "Collection: empty")
	assert.Contains(t, schema, "Fields: (none)")
}

func TestExtractSchema_HiddenFields(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "test")
	collection.Fields.Add(&core.TextField{Name: "title"})
	passwordField := &core.PasswordField{Name: "password"}
	passwordField.SetHidden(true)
	collection.Fields.Add(passwordField)
	tokenField := &core.TextField{Name: "tokenKey"}
	tokenField.SetHidden(true)
	collection.Fields.Add(tokenField)

	schema := ExtractSchema(app, collection)

	assert.Contains(t, schema, "title (text)")
	assert.NotContains(t, schema, "password")
	assert.NotContains(t, schema, "tokenKey")
}

func TestExtractSchema_NilCollection(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	schema := ExtractSchema(app, nil)

	assert.Equal(t, "No collection provided", schema)
}

func TestExtractSchema_SelectFieldWithoutOptions(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	collection := core.NewCollection(core.CollectionTypeBase, "test")
	collection.Fields.Add(&core.SelectField{
		Name:   "status",
		Values: []string{},
	})

	schema := ExtractSchema(app, collection)

	assert.Contains(t, schema, "status (select)")
	assert.NotContains(t, schema, "status (select: )")
}

