package naming_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/naming"
)

type NamingTestSuite struct {
	suite.Suite
}

func TestNamingTestSuite(t *testing.T) {
	suite.Run(t, new(NamingTestSuite))
}

// ── ToSnakeCase ─────────────────────────────────────────────────────────────

func (suite *NamingTestSuite) TestToSnakeCase_Simple() {
	suite.Equal("organization", naming.ToSnakeCase("Organization"))
}

func (suite *NamingTestSuite) TestToSnakeCase_MultiWord() {
	suite.Equal("plugin_version", naming.ToSnakeCase("PluginVersion"))
}

func (suite *NamingTestSuite) TestToSnakeCase_Acronym() {
	suite.Equal("id", naming.ToSnakeCase("ID"))
}

func (suite *NamingTestSuite) TestToSnakeCase_AcronymInMiddle() {
	suite.Equal("specification_version_id", naming.ToSnakeCase("SpecificationVersionID"))
}

func (suite *NamingTestSuite) TestToSnakeCase_Empty() {
	suite.Equal("", naming.ToSnakeCase(""))
}

func (suite *NamingTestSuite) TestToSnakeCase_SingleChar() {
	suite.Equal("a", naming.ToSnakeCase("A"))
}

func (suite *NamingTestSuite) TestToSnakeCase_ThreeWords() {
	suite.Equal("format_version_id", naming.ToSnakeCase("FormatVersionID"))
}

// ── ToKebabCase ─────────────────────────────────────────────────────────────

func (suite *NamingTestSuite) TestToKebabCase_Simple() {
	suite.Equal("organization", naming.ToKebabCase("Organization"))
}

func (suite *NamingTestSuite) TestToKebabCase_MultiWord() {
	suite.Equal("plugin-version", naming.ToKebabCase("PluginVersion"))
}

func (suite *NamingTestSuite) TestToKebabCase_ThreeWords() {
	suite.Equal("specification-version", naming.ToKebabCase("SpecificationVersion"))
}

// ── Pluralize ───────────────────────────────────────────────────────────────

func (suite *NamingTestSuite) TestPluralize_Simple() {
	suite.Equal("organizations", naming.Pluralize("organization"))
}

func (suite *NamingTestSuite) TestPluralize_EndsWithS() {
	suite.Equal("statuses", naming.Pluralize("status"))
}

func (suite *NamingTestSuite) TestPluralize_EndsWithX() {
	suite.Equal("indexes", naming.Pluralize("index"))
}

func (suite *NamingTestSuite) TestPluralize_EndsWithConsonantY() {
	suite.Equal("entities", naming.Pluralize("entity"))
}

func (suite *NamingTestSuite) TestPluralize_EndsWithVowelY() {
	suite.Equal("keys", naming.Pluralize("key"))
}

func (suite *NamingTestSuite) TestPluralize_Empty() {
	suite.Equal("", naming.Pluralize(""))
}

// ── CollectionName ──────────────────────────────────────────────────────────

func (suite *NamingTestSuite) TestCollectionName_Simple() {
	suite.Equal("organizations", naming.CollectionName("Organization"))
}

func (suite *NamingTestSuite) TestCollectionName_MultiWord() {
	suite.Equal("plugin-versions", naming.CollectionName("PluginVersion"))
}

func (suite *NamingTestSuite) TestCollectionName_ThreeWords() {
	suite.Equal("specification-versions", naming.CollectionName("SpecificationVersion"))
}

// ── QueryParamName ──────────────────────────────────────────────────────────

func (suite *NamingTestSuite) TestQueryParamName_Simple() {
	suite.Equal("organization_id", naming.QueryParamName("Organization"))
}

func (suite *NamingTestSuite) TestQueryParamName_MultiWord() {
	suite.Equal("format_version_id", naming.QueryParamName("FormatVersion"))
}

// ── FilterFieldName ─────────────────────────────────────────────────────────

func (suite *NamingTestSuite) TestFilterFieldName_Simple() {
	suite.Equal("organizationID", naming.FilterFieldName("Organization"))
}

func (suite *NamingTestSuite) TestFilterFieldName_MultiWord() {
	suite.Equal("formatVersionID", naming.FilterFieldName("FormatVersion"))
}
