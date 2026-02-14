package compile_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile"
)

type ModelInfoTestSuite struct {
	suite.Suite
}

func TestModelInfoTestSuite(t *testing.T) {
	suite.Run(t, new(ModelInfoTestSuite))
}

// ── ExtractModelInfo ────────────────────────────────────────────────────────

func (suite *ModelInfoTestSuite) TestExtractModelInfo_SimpleModel() {
	model := yaml.Model{
		Name: "Organization",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeUUID},
			"Code": {Type: yaml.ModelFieldTypeString},
			"Name": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
			"code":    {Fields: []string{"Code"}},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	info := compile.ExtractModelInfo(model)

	suite.Equal("Organization", info.Name)
	suite.Equal("organizations", info.CollectionName)
	suite.Equal("organization", info.SnakeName)
	suite.Equal("ID", info.PrimaryIDField)
	suite.Equal(yaml.ModelFieldTypeUUID, info.PrimaryIDType)
	suite.Len(info.Filters, 0)
	suite.True(info.HasCodeIdentifier())
	suite.Equal("Code", info.CodeIdentifierField())
}

func (suite *ModelInfoTestSuite) TestExtractModelInfo_WithForOneRelation() {
	model := yaml.Model{
		Name: "Project",
		Fields: map[string]yaml.ModelField{
			"ID":          {Type: yaml.ModelFieldTypeUUID},
			"Code":        {Type: yaml.ModelFieldTypeString},
			"Name":        {Type: yaml.ModelFieldTypeString},
			"Description": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
			"code":    {Fields: []string{"Code"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Organization": {Type: "ForOne"},
		},
	}

	info := compile.ExtractModelInfo(model)

	suite.Equal("Project", info.Name)
	suite.Equal("projects", info.CollectionName)
	suite.Len(info.Filters, 1)
	suite.Equal("Organization", info.Filters[0].RelationName)
	suite.Equal("organization_id", info.Filters[0].ParamName)
	suite.Equal("organizationID", info.Filters[0].GoFieldName)
	suite.True(info.HasCodeIdentifier())
}

func (suite *ModelInfoTestSuite) TestExtractModelInfo_WithMultipleForOneRelations() {
	model := yaml.Model{
		Name: "PluginVersion",
		Fields: map[string]yaml.ModelField{
			"ID":          {Type: yaml.ModelFieldTypeUUID},
			"Description": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"FormatVersion": {Type: "ForOne"},
			"Plugin":        {Type: "ForOne"},
		},
	}

	info := compile.ExtractModelInfo(model)

	suite.Equal("PluginVersion", info.Name)
	suite.Equal("plugin-versions", info.CollectionName)
	suite.Len(info.Filters, 2)
	suite.False(info.HasCodeIdentifier())
	suite.Equal("", info.CodeIdentifierField())

	// Filters should be sorted by relation name
	suite.Equal("FormatVersion", info.Filters[0].RelationName)
	suite.Equal("format_version_id", info.Filters[0].ParamName)
	suite.Equal("formatVersionID", info.Filters[0].GoFieldName)
	suite.Equal("Plugin", info.Filters[1].RelationName)
	suite.Equal("plugin_id", info.Filters[1].ParamName)
	suite.Equal("pluginID", info.Filters[1].GoFieldName)
}

func (suite *ModelInfoTestSuite) TestExtractModelInfo_NoCodeIdentifier() {
	model := yaml.Model{
		Name: "Task",
		Fields: map[string]yaml.ModelField{
			"ID":     {Type: yaml.ModelFieldTypeUUID},
			"Title":  {Type: yaml.ModelFieldTypeString},
			"Status": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Project": {Type: "ForOne"},
		},
	}

	info := compile.ExtractModelInfo(model)

	suite.Equal("Task", info.Name)
	suite.False(info.HasCodeIdentifier())
	suite.Equal("", info.CodeIdentifierField())
	suite.Len(info.Filters, 1)
}

func (suite *ModelInfoTestSuite) TestExtractModelInfo_WithForOnePoly() {
	model := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeAutoIncrement},
			"Text": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Person", "Company"},
			},
		},
	}

	info := compile.ExtractModelInfo(model)

	suite.Equal("Comment", info.Name)
	suite.Len(info.Filters, 1)
	suite.Equal("Commentable", info.Filters[0].RelationName)
	suite.Equal("commentable_id", info.Filters[0].ParamName)
	suite.Equal("commentableID", info.Filters[0].GoFieldName)
}

func (suite *ModelInfoTestSuite) TestExtractModelInfo_WithForOnePoly_Through() {
	model := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeAutoIncrement},
			"Text": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Target": {
				Type:    "ForOnePoly",
				For:     []string{"Person", "Company"},
				Through: "Targetable",
			},
		},
	}

	info := compile.ExtractModelInfo(model)

	suite.Len(info.Filters, 1)
	suite.Equal("Targetable", info.Filters[0].RelationName)
	suite.Equal("targetable_id", info.Filters[0].ParamName)
	suite.Equal("targetableID", info.Filters[0].GoFieldName)
}

func (suite *ModelInfoTestSuite) TestExtractModelInfo_FieldsSorted() {
	model := yaml.Model{
		Name: "Widget",
		Fields: map[string]yaml.ModelField{
			"ZName":       {Type: yaml.ModelFieldTypeString},
			"ID":          {Type: yaml.ModelFieldTypeUUID},
			"Description": {Type: yaml.ModelFieldTypeString},
			"ALabel":      {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	info := compile.ExtractModelInfo(model)

	// Regular fields should be sorted alphabetically
	suite.Equal("ALabel", info.Fields[0].Name)
	suite.Equal("Description", info.Fields[1].Name)
	suite.Equal("ID", info.Fields[2].Name)
	suite.Equal("ZName", info.Fields[3].Name)
}

func (suite *ModelInfoTestSuite) TestExtractModelInfo_OptionalField() {
	model := yaml.Model{
		Name: "Item",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeUUID},
			"Name": {Type: yaml.ModelFieldTypeString, Attributes: []string{"optional"}},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{},
	}

	info := compile.ExtractModelInfo(model)

	// Find the Name field (has "optional" attribute)
	var nameField compile.FieldInfo
	for _, f := range info.Fields {
		if f.Name == "Name" {
			nameField = f
			break
		}
	}
	suite.True(nameField.IsOptional)

	// Find the ID field (no "optional" attribute)
	var idField compile.FieldInfo
	for _, f := range info.Fields {
		if f.Name == "ID" {
			idField = f
			break
		}
	}
	suite.False(idField.IsOptional)
}

func (suite *ModelInfoTestSuite) TestExtractModelInfo_FKFields() {
	model := yaml.Model{
		Name: "Project",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeUUID},
			"Name": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Organization": {Type: "ForOne"},
		},
	}

	info := compile.ExtractModelInfo(model)

	// Find the FK field
	var fkField compile.FieldInfo
	for _, f := range info.Fields {
		if f.Name == "OrganizationID" {
			fkField = f
			break
		}
	}
	suite.True(fkField.IsFK)
	suite.Equal("organization_id", fkField.SnakeName)
}

func (suite *ModelInfoTestSuite) TestExtractModelInfo_PolyFKFields() {
	model := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeAutoIncrement},
			"Text": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Person", "Company"},
			},
		},
	}

	info := compile.ExtractModelInfo(model)

	// Should have CommentableID (FK) and CommentableType (Poly)
	var fkIDField, fkTypeField compile.FieldInfo
	for _, f := range info.Fields {
		if f.Name == "CommentableID" {
			fkIDField = f
		}
		if f.Name == "CommentableType" {
			fkTypeField = f
		}
	}
	suite.True(fkIDField.IsFK)
	suite.True(fkIDField.IsPoly)
	suite.True(fkTypeField.IsFK)
	suite.True(fkTypeField.IsPoly)
}

// ── HasMany relations should NOT create filters ─────────────────────────────

func (suite *ModelInfoTestSuite) TestExtractModelInfo_HasManyDoesNotCreateFilter() {
	model := yaml.Model{
		Name: "Organization",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeUUID},
			"Name": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Project": {Type: "HasMany"},
		},
	}

	info := compile.ExtractModelInfo(model)

	suite.Len(info.Filters, 0)
}
