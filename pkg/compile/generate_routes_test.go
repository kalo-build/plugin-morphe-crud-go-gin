package compile_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile/cfg"
)

type GenerateRoutesTestSuite struct {
	suite.Suite
}

func TestGenerateRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(GenerateRoutesTestSuite))
}

func (suite *GenerateRoutesTestSuite) getConfig() cfg.CompileConfig {
	config := cfg.CompileConfig{
		HandlersConfig: cfg.HandlersConfig{
			PackagePath: "github.com/test/app/internal/generated/handlers",
		},
		RepoConfig: cfg.RepoConfig{
			PackagePath: "github.com/test/app/internal/generated/repo",
		},
		ModelsConfig: cfg.ModelsConfig{
			PackagePath: "github.com/test/app/internal/types/models",
		},
		StructuresConfig: cfg.StructuresConfig{
			PackagePath: "github.com/test/app/internal/types/structures",
		},
	}
	config.Resolve()
	return config
}

func (suite *GenerateRoutesTestSuite) getModelInfos() []compile.ModelInfo {
	org := compile.ExtractModelInfo(yaml.Model{
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
	})

	project := compile.ExtractModelInfo(yaml.Model{
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
	})

	return []compile.ModelInfo{org, project}
}

// ── Package & imports ───────────────────────────────────────────────────────

func (suite *GenerateRoutesTestSuite) TestGenerateRoutes_PackageAndImports() {
	config := suite.getConfig()
	models := suite.getModelInfos()

	output := compile.GenerateRoutes(config, models)

	suite.Contains(output, "package handlers")
	suite.Contains(output, `"github.com/gin-gonic/gin"`)
	suite.Contains(output, `"github.com/test/app/internal/generated/repo"`)
}

// ── Repos struct ────────────────────────────────────────────────────────────

func (suite *GenerateRoutesTestSuite) TestGenerateRoutes_ReposStruct() {
	config := suite.getConfig()
	models := suite.getModelInfos()

	output := compile.GenerateRoutes(config, models)

	suite.Contains(output, "type Repos struct {")
	suite.Contains(output, "Organization repo.OrganizationRepository")
	suite.Contains(output, "Project repo.ProjectRepository")
}

// ── RegisterAllRoutes ───────────────────────────────────────────────────────

func (suite *GenerateRoutesTestSuite) TestGenerateRoutes_RegisterAllRoutes() {
	config := suite.getConfig()
	models := suite.getModelInfos()

	output := compile.GenerateRoutes(config, models)

	suite.Contains(output, "func RegisterAllRoutes(rg *gin.RouterGroup, repos Repos) {")
	suite.Contains(output, "NewOrganizationHandler(repos.Organization).RegisterRoutes(rg)")
	suite.Contains(output, "NewProjectHandler(repos.Project).RegisterRoutes(rg)")
}

// ── Single model ────────────────────────────────────────────────────────────

func (suite *GenerateRoutesTestSuite) TestGenerateRoutes_SingleModel() {
	config := suite.getConfig()
	models := []compile.ModelInfo{suite.getModelInfos()[0]}

	output := compile.GenerateRoutes(config, models)

	suite.Contains(output, "Organization repo.OrganizationRepository")
	suite.Contains(output, "NewOrganizationHandler(repos.Organization).RegisterRoutes(rg)")
	suite.NotContains(output, "Project")
}

// ── Empty models list (should still produce valid structure) ─────────────────

func (suite *GenerateRoutesTestSuite) TestGenerateRoutes_EmptyModels() {
	config := suite.getConfig()

	output := compile.GenerateRoutes(config, []compile.ModelInfo{})

	suite.Contains(output, "package handlers")
	suite.Contains(output, "type Repos struct {")
	suite.Contains(output, "func RegisterAllRoutes(rg *gin.RouterGroup, repos Repos) {")
}
