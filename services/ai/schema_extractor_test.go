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

// ============================================================================
// V2 Multi-Collection Schema Extraction Tests
// ============================================================================

func TestExtractAllSchemas_MultipleCollections(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create customers collection
	customers := core.NewCollection(core.CollectionTypeBase, "customers")
	customers.Fields.Add(&core.TextField{Name: "name"})
	customers.Fields.Add(&core.EmailField{Name: "email"})
	app.Save(customers)

	// Create orders collection with relation to customers
	orders := core.NewCollection(core.CollectionTypeBase, "orders")
	orders.Fields.Add(&core.NumberField{Name: "total"})
	orders.Fields.Add(&core.SelectField{Name: "status", Values: []string{"pending", "completed"}})
	orders.Fields.Add(&core.RelationField{Name: "customer", CollectionId: customers.Id})
	app.Save(orders)

	app.ReloadCachedCollections()

	schema := ExtractAllSchemas(app)

	// Should contain header
	assert.Contains(t, schema, "DATABASE SCHEMA")
	assert.Contains(t, schema, "TABLES:")

	// Should contain both collections
	assert.Contains(t, schema, "customers:")
	assert.Contains(t, schema, "orders:")

	// Should contain SQL-like field descriptions
	assert.Contains(t, schema, "name TEXT")
	assert.Contains(t, schema, "email TEXT (email)")
	assert.Contains(t, schema, "total NUMBER")
	assert.Contains(t, schema, "ENUM(pending, completed)")

	// Should contain relationships section
	assert.Contains(t, schema, "RELATIONSHIPS")
	assert.Contains(t, schema, "orders.customer → customers.id")
}

func TestExtractAllSchemas_NoCollections(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Don't create any collections - note: the test app might have default collections
	// So we test for the schema header being present
	schema := ExtractAllSchemas(app)

	// Should still have header
	assert.Contains(t, schema, "DATABASE SCHEMA")
}

func TestExtractAllSchemas_NilApp(t *testing.T) {
	schema := ExtractAllSchemas(nil)
	assert.Equal(t, "No app provided", schema)
}

func TestExtractRelationships(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create users collection
	users := core.NewCollection(core.CollectionTypeBase, "users")
	users.Fields.Add(&core.TextField{Name: "name"})
	app.Save(users)

	// Create posts collection with relation to users
	posts := core.NewCollection(core.CollectionTypeBase, "posts")
	posts.Fields.Add(&core.TextField{Name: "title"})
	posts.Fields.Add(&core.RelationField{
		Name:         "author",
		CollectionId: users.Id,
		MaxSelect:    1, // Single relation
	})
	app.Save(posts)

	// Create comments collection with relations to both posts and users
	comments := core.NewCollection(core.CollectionTypeBase, "comments")
	comments.Fields.Add(&core.TextField{Name: "body"})
	comments.Fields.Add(&core.RelationField{
		Name:         "post",
		CollectionId: posts.Id,
		MaxSelect:    1,
	})
	comments.Fields.Add(&core.RelationField{
		Name:         "commenter",
		CollectionId: users.Id,
		MaxSelect:    1,
	})
	app.Save(comments)

	app.ReloadCachedCollections()

	relationships := ExtractRelationships(app)

	// Should find 3 relationships
	assert.GreaterOrEqual(t, len(relationships), 3, "Should have at least 3 relationships")

	// Check for specific relationships
	foundPostsAuthor := false
	foundCommentsPost := false
	foundCommentsUser := false

	for _, rel := range relationships {
		if rel.FromCollection == "posts" && rel.FromField == "author" && rel.ToCollection == "users" {
			foundPostsAuthor = true
			assert.False(t, rel.IsMultiple, "posts.author should be single relation")
		}
		if rel.FromCollection == "comments" && rel.FromField == "post" && rel.ToCollection == "posts" {
			foundCommentsPost = true
		}
		if rel.FromCollection == "comments" && rel.FromField == "commenter" && rel.ToCollection == "users" {
			foundCommentsUser = true
		}
	}

	assert.True(t, foundPostsAuthor, "Should find posts.author → users relationship")
	assert.True(t, foundCommentsPost, "Should find comments.post → posts relationship")
	assert.True(t, foundCommentsUser, "Should find comments.commenter → users relationship")
}

func TestExtractRelationships_MultipleRelation(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create tags collection
	tags := core.NewCollection(core.CollectionTypeBase, "tags")
	tags.Fields.Add(&core.TextField{Name: "name"})
	app.Save(tags)

	// Create posts with multi-select relation to tags
	posts := core.NewCollection(core.CollectionTypeBase, "posts")
	posts.Fields.Add(&core.TextField{Name: "title"})
	posts.Fields.Add(&core.RelationField{
		Name:         "tags",
		CollectionId: tags.Id,
		MaxSelect:    0, // 0 means unlimited = multiple
	})
	app.Save(posts)

	app.ReloadCachedCollections()

	relationships := ExtractRelationships(app)

	// Find the posts.tags relationship
	for _, rel := range relationships {
		if rel.FromCollection == "posts" && rel.FromField == "tags" {
			assert.True(t, rel.IsMultiple, "posts.tags should be multiple relation")
			assert.Equal(t, "tags", rel.ToCollection)
			return
		}
	}
	t.Error("Should find posts.tags relationship")
}

func TestExtractSchemaForCollection(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create customers collection
	customers := core.NewCollection(core.CollectionTypeBase, "customers")
	customers.Fields.Add(&core.TextField{Name: "name"})
	customers.Fields.Add(&core.EmailField{Name: "email"})
	app.Save(customers)

	// Create orders collection with relation to customers
	orders := core.NewCollection(core.CollectionTypeBase, "orders")
	orders.Fields.Add(&core.NumberField{Name: "total"})
	orders.Fields.Add(&core.RelationField{Name: "customer", CollectionId: customers.Id})
	app.Save(orders)

	app.ReloadCachedCollections()

	// Get schema for orders (should include related customers)
	schema := ExtractSchemaForCollection(app, "orders")

	// Should contain primary table
	assert.Contains(t, schema, "PRIMARY TABLE:")
	assert.Contains(t, schema, "orders:")
	assert.Contains(t, schema, "total NUMBER")

	// Should contain related table
	assert.Contains(t, schema, "RELATED TABLES:")
	assert.Contains(t, schema, "customers:")
	assert.Contains(t, schema, "name TEXT")

	// Should contain relationships
	assert.Contains(t, schema, "RELATIONSHIPS:")
	assert.Contains(t, schema, "orders.customer → customers.id")
}

func TestExtractSchemaForCollection_NotFound(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	schema := ExtractSchemaForCollection(app, "nonexistent")
	assert.Contains(t, schema, "Collection 'nonexistent' not found")
}

func TestExtractSchemaForCollection_NoRelations(t *testing.T) {
	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// Create a collection without relations
	products := core.NewCollection(core.CollectionTypeBase, "products")
	products.Fields.Add(&core.TextField{Name: "name"})
	products.Fields.Add(&core.NumberField{Name: "price"})
	app.Save(products)

	app.ReloadCachedCollections()

	schema := ExtractSchemaForCollection(app, "products")

	// Should contain primary table
	assert.Contains(t, schema, "PRIMARY TABLE:")
	assert.Contains(t, schema, "products:")

	// Should NOT contain related tables section
	assert.NotContains(t, schema, "RELATED TABLES:")
}

