package ai

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pocketbase/pocketbase/core"
)

// Relationship represents a foreign key relationship between collections
type Relationship struct {
	FromCollection string
	FromField      string
	ToCollection   string
	ToField        string // Usually "id"
	IsMultiple     bool   // True if MaxSelect > 1
}

// ExtractSchema converts a collection to a prompt-friendly format that describes
// all fields, their types, and relevant metadata.
// The app parameter is used to resolve relation collection names from IDs.
func ExtractSchema(app core.App, collection *core.Collection) string {
	if collection == nil {
		return "No collection provided"
	}

	var parts []string
	parts = append(parts, fmt.Sprintf("Collection: %s", collection.Name))
	parts = append(parts, "")

	if len(collection.Fields) == 0 {
		parts = append(parts, "Fields: (none)")
		return strings.Join(parts, "\n")
	}

	parts = append(parts, "Fields:")
	for _, field := range collection.Fields {
		// Skip hidden system fields (like password, tokenKey)
		if field.GetHidden() {
			continue
		}

		fieldDesc := extractFieldDescription(app, field)
		if fieldDesc != "" {
			parts = append(parts, "  "+fieldDesc)
		}
	}

	return strings.Join(parts, "\n")
}

// ExtractAllSchemas extracts schemas for ALL collections in the database,
// formatted for multi-table query understanding (JOINs, SQL).
func ExtractAllSchemas(app core.App) string {
	if app == nil {
		return "No app provided"
	}

	collections, err := app.FindAllCollections()
	if err != nil || len(collections) == 0 {
		return "No collections found"
	}

	var parts []string
	parts = append(parts, "DATABASE SCHEMA")
	parts = append(parts, "===============")
	parts = append(parts, "")

	// Sort collections by name for consistent output
	sort.Slice(collections, func(i, j int) bool {
		return collections[i].Name < collections[j].Name
	})

	// Extract schema for each collection
	parts = append(parts, "TABLES:")
	parts = append(parts, "-------")
	for _, collection := range collections {
		// Skip system collections (those starting with underscore)
		if strings.HasPrefix(collection.Name, "_") {
			continue
		}
		collectionSchema := extractCollectionSchemaForSQL(app, collection)
		parts = append(parts, collectionSchema)
		parts = append(parts, "")
	}

	// Extract relationships
	relationships := ExtractRelationships(app)
	if len(relationships) > 0 {
		parts = append(parts, "RELATIONSHIPS (Foreign Keys):")
		parts = append(parts, "-----------------------------")
		for _, rel := range relationships {
			relType := "many-to-one"
			if rel.IsMultiple {
				relType = "many-to-many"
			}
			parts = append(parts, fmt.Sprintf("  %s.%s → %s.%s (%s)",
				rel.FromCollection, rel.FromField,
				rel.ToCollection, rel.ToField,
				relType))
		}
		parts = append(parts, "")
	}

	return strings.Join(parts, "\n")
}

// extractCollectionSchemaForSQL formats a collection schema in SQL-like format
func extractCollectionSchemaForSQL(app core.App, collection *core.Collection) string {
	if collection == nil {
		return ""
	}

	var fields []string
	for _, field := range collection.Fields {
		if field.GetHidden() {
			continue
		}
		fieldDesc := extractFieldDescriptionSQL(app, field)
		if fieldDesc != "" {
			fields = append(fields, fieldDesc)
		}
	}

	if len(fields) == 0 {
		return fmt.Sprintf("  %s (no fields)", collection.Name)
	}

	return fmt.Sprintf("  %s:\n    %s", collection.Name, strings.Join(fields, "\n    "))
}

// extractFieldDescriptionSQL formats a field for SQL-like schema output
func extractFieldDescriptionSQL(app core.App, field core.Field) string {
	fieldName := field.GetName()
	fieldType := field.Type()

	switch fieldType {
	case core.FieldTypeText:
		return fmt.Sprintf("%s TEXT", fieldName)
	case core.FieldTypeNumber:
		return fmt.Sprintf("%s NUMBER", fieldName)
	case core.FieldTypeBool:
		return fmt.Sprintf("%s BOOLEAN", fieldName)
	case core.FieldTypeEmail:
		return fmt.Sprintf("%s TEXT (email)", fieldName)
	case core.FieldTypeURL:
		return fmt.Sprintf("%s TEXT (url)", fieldName)
	case core.FieldTypeDate, core.FieldTypeAutodate:
		return fmt.Sprintf("%s DATETIME", fieldName)
	case core.FieldTypeSelect:
		if selectField, ok := field.(*core.SelectField); ok {
			if len(selectField.Values) > 0 {
				options := strings.Join(selectField.Values, ", ")
				return fmt.Sprintf("%s ENUM(%s)", fieldName, options)
			}
		}
		return fmt.Sprintf("%s ENUM", fieldName)
	case core.FieldTypeRelation:
		if relationField, ok := field.(*core.RelationField); ok {
			if relationField.CollectionId != "" && app != nil {
				if relatedCollection, err := app.FindCachedCollectionByNameOrId(relationField.CollectionId); err == nil {
					return fmt.Sprintf("%s FOREIGN KEY → %s.id", fieldName, relatedCollection.Name)
				}
			}
		}
		return fmt.Sprintf("%s FOREIGN KEY", fieldName)
	case core.FieldTypeFile:
		return fmt.Sprintf("%s FILE", fieldName)
	case core.FieldTypeJSON:
		return fmt.Sprintf("%s JSON", fieldName)
	case core.FieldTypeEditor:
		return fmt.Sprintf("%s TEXT (rich)", fieldName)
	case core.FieldTypePassword:
		return ""
	case core.FieldTypeGeoPoint:
		return fmt.Sprintf("%s GEOPOINT", fieldName)
	default:
		return fmt.Sprintf("%s %s", fieldName, strings.ToUpper(fieldType))
	}
}

// ExtractRelationships detects all foreign key relationships between collections
func ExtractRelationships(app core.App) []Relationship {
	if app == nil {
		return nil
	}

	collections, err := app.FindAllCollections()
	if err != nil {
		return nil
	}

	var relationships []Relationship

	for _, collection := range collections {
		// Skip system collections
		if strings.HasPrefix(collection.Name, "_") {
			continue
		}

		for _, field := range collection.Fields {
			if field.Type() == core.FieldTypeRelation {
				if relationField, ok := field.(*core.RelationField); ok {
					if relationField.CollectionId != "" && app != nil {
						if relatedCollection, err := app.FindCachedCollectionByNameOrId(relationField.CollectionId); err == nil {
							isMultiple := relationField.MaxSelect == 0 || relationField.MaxSelect > 1
							relationships = append(relationships, Relationship{
								FromCollection: collection.Name,
								FromField:      field.GetName(),
								ToCollection:   relatedCollection.Name,
								ToField:        "id",
								IsMultiple:     isMultiple,
							})
						}
					}
				}
			}
		}
	}

	return relationships
}

// ExtractSchemaForCollection extracts schema for a specific collection with its related collections
// This is useful for queries that might need JOIN context
func ExtractSchemaForCollection(app core.App, collectionName string) string {
	if app == nil {
		return "No app provided"
	}

	collection, err := app.FindCachedCollectionByNameOrId(collectionName)
	if err != nil {
		return fmt.Sprintf("Collection '%s' not found", collectionName)
	}

	var parts []string
	parts = append(parts, "PRIMARY TABLE:")
	parts = append(parts, "--------------")
	parts = append(parts, extractCollectionSchemaForSQL(app, collection))
	parts = append(parts, "")

	// Find related collections
	relatedCollections := make(map[string]*core.Collection)
	for _, field := range collection.Fields {
		if field.Type() == core.FieldTypeRelation {
			if relationField, ok := field.(*core.RelationField); ok {
				if relationField.CollectionId != "" {
					if relatedCollection, err := app.FindCachedCollectionByNameOrId(relationField.CollectionId); err == nil {
						relatedCollections[relatedCollection.Name] = relatedCollection
					}
				}
			}
		}
	}

	if len(relatedCollections) > 0 {
		parts = append(parts, "RELATED TABLES:")
		parts = append(parts, "---------------")
		for _, relColl := range relatedCollections {
			parts = append(parts, extractCollectionSchemaForSQL(app, relColl))
			parts = append(parts, "")
		}
	}

	// Add relationships for context
	relationships := ExtractRelationships(app)
	relevantRels := []Relationship{}
	for _, rel := range relationships {
		if rel.FromCollection == collectionName || rel.ToCollection == collectionName {
			relevantRels = append(relevantRels, rel)
		}
	}

	if len(relevantRels) > 0 {
		parts = append(parts, "RELATIONSHIPS:")
		parts = append(parts, "--------------")
		for _, rel := range relevantRels {
			parts = append(parts, fmt.Sprintf("  %s.%s → %s.%s", rel.FromCollection, rel.FromField, rel.ToCollection, rel.ToField))
		}
	}

	return strings.Join(parts, "\n")
}

// extractFieldDescription converts a single field to a human-readable description.
func extractFieldDescription(app core.App, field core.Field) string {
	fieldName := field.GetName()
	fieldType := field.Type()

	switch fieldType {
	case core.FieldTypeText:
		return fmt.Sprintf("%s (text)", fieldName)

	case core.FieldTypeNumber:
		return fmt.Sprintf("%s (number)", fieldName)

	case core.FieldTypeBool:
		return fmt.Sprintf("%s (bool)", fieldName)

	case core.FieldTypeEmail:
		return fmt.Sprintf("%s (email)", fieldName)

	case core.FieldTypeURL:
		return fmt.Sprintf("%s (url)", fieldName)

	case core.FieldTypeDate, core.FieldTypeAutodate:
		return fmt.Sprintf("%s (date)", fieldName)

	case core.FieldTypeSelect:
		if selectField, ok := field.(*core.SelectField); ok {
			if len(selectField.Values) > 0 {
				options := strings.Join(selectField.Values, "|")
				return fmt.Sprintf("%s (select: %s)", fieldName, options)
			}
			return fmt.Sprintf("%s (select)", fieldName)
		}
		return fmt.Sprintf("%s (select)", fieldName)

	case core.FieldTypeRelation:
		if relationField, ok := field.(*core.RelationField); ok {
			if relationField.CollectionId != "" {
				// Try to resolve collection name from ID
				if app != nil {
					if relatedCollection, err := app.FindCachedCollectionByNameOrId(relationField.CollectionId); err == nil {
						return fmt.Sprintf("%s (relation → %s)", fieldName, relatedCollection.Name)
					}
				}
				// Fallback to ID if we can't resolve the name
				return fmt.Sprintf("%s (relation → %s)", fieldName, relationField.CollectionId)
			}
			return fmt.Sprintf("%s (relation)", fieldName)
		}
		return fmt.Sprintf("%s (relation)", fieldName)

	case core.FieldTypeFile:
		return fmt.Sprintf("%s (file)", fieldName)

	case core.FieldTypeJSON:
		return fmt.Sprintf("%s (json)", fieldName)

	case core.FieldTypeEditor:
		return fmt.Sprintf("%s (editor)", fieldName)

	case core.FieldTypePassword:
		// Skip password fields in schema (they're hidden anyway)
		return ""

	case core.FieldTypeGeoPoint:
		return fmt.Sprintf("%s (geopoint)", fieldName)

	default:
		return fmt.Sprintf("%s (%s)", fieldName, fieldType)
	}
}

