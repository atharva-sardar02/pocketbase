package ai

import (
	"fmt"
	"strings"

	"github.com/pocketbase/pocketbase/core"
)

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

