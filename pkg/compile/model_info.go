package compile

import (
	"sort"

	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/naming"
)

// ModelInfo holds derived information about a Morphe model needed for code generation.
type ModelInfo struct {
	// Name is the PascalCase model name (e.g., "Organization").
	Name string

	// CollectionName is the URL path segment (e.g., "organizations").
	CollectionName string

	// SnakeName is the snake_case name (e.g., "organization").
	SnakeName string

	// PrimaryIDField is the field name used as primary identifier (e.g., "ID").
	PrimaryIDField string

	// PrimaryIDType is the Morphe field type of the primary identifier (e.g., "UUID").
	PrimaryIDType yaml.ModelFieldType

	// SecondaryIdentifiers are non-primary identifiers (e.g., "code" -> ["Code"]).
	// Key is the identifier name, value is the list of field names.
	SecondaryIdentifiers map[string][]string

	// Filters are derived from ForOne/ForOnePoly relationships.
	// Each filter represents an optional query parameter on GetAll.
	Filters []FilterInfo

	// Fields are all model fields in sorted order.
	Fields []FieldInfo
}

// FilterInfo describes a filter parameter derived from a relationship.
type FilterInfo struct {
	// RelationName is the PascalCase relationship name (e.g., "Organization").
	RelationName string

	// ParamName is the query parameter name (e.g., "organization_id").
	ParamName string

	// GoFieldName is the Go parameter name (e.g., "organizationID").
	GoFieldName string
}

// FieldInfo describes a model field.
type FieldInfo struct {
	// Name is the PascalCase field name (e.g., "DisplayName").
	Name string

	// Type is the Morphe field type (e.g., "String").
	Type yaml.ModelFieldType

	// SnakeName is the snake_case name for JSON output (e.g., "display_name").
	SnakeName string

	// IsOptional indicates whether the field has the "optional" attribute.
	IsOptional bool

	// IsFK indicates this is a foreign key field derived from a relationship.
	IsFK bool

	// IsPoly indicates this is a polymorphic type field.
	IsPoly bool
}

// ExtractModelInfo analyzes a Morphe model and extracts the information needed for code generation.
func ExtractModelInfo(model yaml.Model) ModelInfo {
	info := ModelInfo{
		Name:                 model.Name,
		CollectionName:       naming.CollectionName(model.Name),
		SnakeName:            naming.ToSnakeCase(model.Name),
		SecondaryIdentifiers: make(map[string][]string),
	}

	// Extract identifiers
	for idName, id := range model.Identifiers {
		if idName == "primary" {
			if len(id.Fields) > 0 {
				info.PrimaryIDField = id.Fields[0]
				if field, ok := model.Fields[id.Fields[0]]; ok {
					info.PrimaryIDType = field.Type
				}
			}
		} else {
			info.SecondaryIdentifiers[idName] = id.Fields
		}
	}

	// Extract filters from ForOne/ForOnePoly relationships
	relNames := sortedKeys(model.Related)
	for _, relName := range relNames {
		rel := model.Related[relName]
		switch rel.Type {
		case "ForOne":
			info.Filters = append(info.Filters, FilterInfo{
				RelationName: relName,
				ParamName:    naming.QueryParamName(relName),
				GoFieldName:  naming.FilterFieldName(relName),
			})
		case "ForOnePoly":
			// For polymorphic relations, use the "through" name if available
			filterName := relName
			if rel.Through != "" {
				filterName = rel.Through
			}
			info.Filters = append(info.Filters, FilterInfo{
				RelationName: filterName,
				ParamName:    naming.QueryParamName(filterName),
				GoFieldName:  naming.FilterFieldName(filterName),
			})
		}
	}

	// Extract fields in sorted order
	fieldNames := sortedKeys(model.Fields)
	for _, fieldName := range fieldNames {
		field := model.Fields[fieldName]
		isOptional := false
		for _, attr := range field.Attributes {
			if attr == "optional" {
				isOptional = true
				break
			}
		}
		info.Fields = append(info.Fields, FieldInfo{
			Name:       fieldName,
			Type:       field.Type,
			SnakeName:  naming.ToSnakeCase(fieldName),
			IsOptional: isOptional,
		})
	}

	// Add FK fields from ForOne relationships
	for _, relName := range relNames {
		rel := model.Related[relName]
		switch rel.Type {
		case "ForOne":
			info.Fields = append(info.Fields, FieldInfo{
				Name:       relName + "ID",
				SnakeName:  naming.ToSnakeCase(relName) + "_id",
				IsFK:       true,
				IsOptional: true,
			})
		case "ForOnePoly":
			through := relName
			if rel.Through != "" {
				through = rel.Through
			}
			info.Fields = append(info.Fields, FieldInfo{
				Name:      through + "ID",
				SnakeName: naming.ToSnakeCase(through) + "_id",
				IsFK:      true,
				IsPoly:    true,
			})
			info.Fields = append(info.Fields, FieldInfo{
				Name:      through + "Type",
				SnakeName: naming.ToSnakeCase(through) + "_type",
				IsFK:      true,
				IsPoly:    true,
			})
		}
	}

	return info
}

// HasCodeIdentifier returns true if the model has a "code" secondary identifier.
func (m ModelInfo) HasCodeIdentifier() bool {
	_, ok := m.SecondaryIdentifiers["code"]
	return ok
}

// CodeIdentifierField returns the field name for the code identifier, or empty string.
func (m ModelInfo) CodeIdentifierField() string {
	fields, ok := m.SecondaryIdentifiers["code"]
	if !ok || len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
