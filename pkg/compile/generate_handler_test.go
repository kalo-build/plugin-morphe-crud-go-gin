package compile_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile/cfg"
)

type GenerateHandlerTestSuite struct {
	suite.Suite
}

func TestGenerateHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GenerateHandlerTestSuite))
}

func (suite *GenerateHandlerTestSuite) getConfig() cfg.CompileConfig {
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

// ── Package & imports ───────────────────────────────────────────────────────

func (suite *GenerateHandlerTestSuite) TestGenerateHandler_PackageAndImports() {
	config := suite.getConfig()
	info := compile.ExtractModelInfo(yaml.Model{
		Name:        "Organization",
		Fields:      map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Identifiers: map[string]yaml.ModelIdentifier{"primary": {Fields: []string{"ID"}}},
		Related:     map[string]yaml.ModelRelation{},
	})

	output := compile.GenerateHandler(config, info)

	suite.Contains(output, "package handlers")
	suite.Contains(output, `"net/http"`)
	suite.Contains(output, `"github.com/gin-gonic/gin"`)
	suite.Contains(output, `"github.com/test/app/internal/types/models"`)
	suite.Contains(output, `"github.com/test/app/internal/generated/repo"`)
	suite.Contains(output, `"github.com/test/app/internal/types/structures"`)
}

// ── Handler struct & constructor ────────────────────────────────────────────

func (suite *GenerateHandlerTestSuite) TestGenerateHandler_StructAndConstructor() {
	config := suite.getConfig()
	info := compile.ExtractModelInfo(yaml.Model{
		Name:        "Organization",
		Fields:      map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Identifiers: map[string]yaml.ModelIdentifier{"primary": {Fields: []string{"ID"}}},
		Related:     map[string]yaml.ModelRelation{},
	})

	output := compile.GenerateHandler(config, info)

	suite.Contains(output, "type OrganizationHandler struct {")
	suite.Contains(output, "repo repo.OrganizationRepository")
	suite.Contains(output, "func NewOrganizationHandler(r repo.OrganizationRepository) *OrganizationHandler {")
}

// ── List handler ────────────────────────────────────────────────────────────

func (suite *GenerateHandlerTestSuite) TestGenerateHandler_ListNoFilters() {
	config := suite.getConfig()
	info := compile.ExtractModelInfo(yaml.Model{
		Name:        "Organization",
		Fields:      map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Identifiers: map[string]yaml.ModelIdentifier{"primary": {Fields: []string{"ID"}}},
		Related:     map[string]yaml.ModelRelation{},
	})

	output := compile.GenerateHandler(config, info)

	suite.Contains(output, "func (h *OrganizationHandler) List(c *gin.Context) {")
	suite.Contains(output, "h.repo.GetAll(ctx)")
	suite.Contains(output, "toOrganizationResponse(item)")
	suite.Contains(output, "[]structures.OrganizationResponse")
}

func (suite *GenerateHandlerTestSuite) TestGenerateHandler_ListWithFilters() {
	config := suite.getConfig()
	info := compile.ExtractModelInfo(yaml.Model{
		Name:        "Project",
		Fields:      map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Identifiers: map[string]yaml.ModelIdentifier{"primary": {Fields: []string{"ID"}}},
		Related: map[string]yaml.ModelRelation{
			"Organization": {Type: "ForOne"},
		},
	})

	output := compile.GenerateHandler(config, info)

	suite.Contains(output, "func (h *ProjectHandler) List(c *gin.Context) {")
	suite.Contains(output, `c.Query("organization_id")`)
	suite.Contains(output, "var organizationID *string")
	suite.Contains(output, "h.repo.GetAll(ctx, organizationID)")
}

// ── Per-identifier Get handlers ─────────────────────────────────────────────

func (suite *GenerateHandlerTestSuite) TestGenerateHandler_GetByID() {
	config := suite.getConfig()
	info := compile.ExtractModelInfo(yaml.Model{
		Name:        "Organization",
		Fields:      map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Identifiers: map[string]yaml.ModelIdentifier{"primary": {Fields: []string{"ID"}}},
		Related:     map[string]yaml.ModelRelation{},
	})

	output := compile.GenerateHandler(config, info)

	suite.Contains(output, "func (h *OrganizationHandler) GetByID(c *gin.Context) {")
	suite.Contains(output, "h.repo.GetByID(ctx, id)")
	suite.NotContains(output, "isUUID")
}

func (suite *GenerateHandlerTestSuite) TestGenerateHandler_GetByCode() {
	config := suite.getConfig()
	info := compile.ExtractModelInfo(yaml.Model{
		Name: "Organization",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeUUID},
			"Code": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
			"code":    {Fields: []string{"Code"}},
		},
		Related: map[string]yaml.ModelRelation{},
	})

	output := compile.GenerateHandler(config, info)

	suite.Contains(output, "func (h *OrganizationHandler) GetByID(c *gin.Context) {")
	suite.Contains(output, "func (h *OrganizationHandler) GetByCode(c *gin.Context) {")
	suite.Contains(output, "h.repo.GetByCode(ctx, code)")
	suite.NotContains(output, "isUUID")
}

func (suite *GenerateHandlerTestSuite) TestGenerateHandler_PrimaryOnlyNoSecondaryHandlers() {
	config := suite.getConfig()
	info := compile.ExtractModelInfo(yaml.Model{
		Name:        "Task",
		Fields:      map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Identifiers: map[string]yaml.ModelIdentifier{"primary": {Fields: []string{"ID"}}},
		Related:     map[string]yaml.ModelRelation{},
	})

	output := compile.GenerateHandler(config, info)

	suite.Contains(output, "func (h *TaskHandler) GetByID(c *gin.Context) {")
	suite.NotContains(output, "GetByCode")
}

// ── Create handler ──────────────────────────────────────────────────────────

func (suite *GenerateHandlerTestSuite) TestGenerateHandler_Create() {
	config := suite.getConfig()
	info := compile.ExtractModelInfo(yaml.Model{
		Name:        "Organization",
		Fields:      map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Identifiers: map[string]yaml.ModelIdentifier{"primary": {Fields: []string{"ID"}}},
		Related:     map[string]yaml.ModelRelation{},
	})

	output := compile.GenerateHandler(config, info)

	suite.Contains(output, "func (h *OrganizationHandler) Create(c *gin.Context) {")
	suite.Contains(output, "var input models.Organization")
	suite.Contains(output, "c.ShouldBindJSON(&input)")
	suite.Contains(output, "h.repo.Create(ctx, &input)")
	suite.Contains(output, "http.StatusCreated")
}

// ── Update handler ──────────────────────────────────────────────────────────

func (suite *GenerateHandlerTestSuite) TestGenerateHandler_Update() {
	config := suite.getConfig()
	info := compile.ExtractModelInfo(yaml.Model{
		Name:        "Organization",
		Fields:      map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Identifiers: map[string]yaml.ModelIdentifier{"primary": {Fields: []string{"ID"}}},
		Related:     map[string]yaml.ModelRelation{},
	})

	output := compile.GenerateHandler(config, info)

	suite.Contains(output, "func (h *OrganizationHandler) Update(c *gin.Context) {")
	suite.Contains(output, "h.repo.Update(ctx, id, &input)")
}

// ── Delete handler ──────────────────────────────────────────────────────────

func (suite *GenerateHandlerTestSuite) TestGenerateHandler_Delete() {
	config := suite.getConfig()
	info := compile.ExtractModelInfo(yaml.Model{
		Name:        "Organization",
		Fields:      map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Identifiers: map[string]yaml.ModelIdentifier{"primary": {Fields: []string{"ID"}}},
		Related:     map[string]yaml.ModelRelation{},
	})

	output := compile.GenerateHandler(config, info)

	suite.Contains(output, "func (h *OrganizationHandler) Delete(c *gin.Context) {")
	suite.Contains(output, "h.repo.Delete(ctx, id)")
	suite.Contains(output, "http.StatusNoContent")
}

// ── RegisterRoutes ──────────────────────────────────────────────────────────

func (suite *GenerateHandlerTestSuite) TestGenerateHandler_RegisterRoutes_PrimaryOnly() {
	config := suite.getConfig()
	info := compile.ExtractModelInfo(yaml.Model{
		Name:        "Task",
		Fields:      map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Identifiers: map[string]yaml.ModelIdentifier{"primary": {Fields: []string{"ID"}}},
		Related:     map[string]yaml.ModelRelation{},
	})

	output := compile.GenerateHandler(config, info)

	suite.Contains(output, "func (h *TaskHandler) RegisterRoutes(rg *gin.RouterGroup) {")
	suite.Contains(output, `rg.GET("/tasks", h.List)`)
	suite.Contains(output, `rg.POST("/tasks", h.Create)`)
	suite.Contains(output, `rg.GET("/tasks/:id", h.GetByID)`)
	suite.Contains(output, `rg.PUT("/tasks/:id", h.Update)`)
	suite.Contains(output, `rg.DELETE("/tasks/:id", h.Delete)`)
	suite.NotContains(output, "by-code")
}

func (suite *GenerateHandlerTestSuite) TestGenerateHandler_RegisterRoutes_WithCodeIdentifier() {
	config := suite.getConfig()
	info := compile.ExtractModelInfo(yaml.Model{
		Name: "Organization",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeUUID},
			"Code": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
			"code":    {Fields: []string{"Code"}},
		},
		Related: map[string]yaml.ModelRelation{},
	})

	output := compile.GenerateHandler(config, info)

	suite.Contains(output, `rg.GET("/organizations/:id", h.GetByID)`)
	suite.Contains(output, `rg.GET("/organizations/by-code/:code", h.GetByCode)`)
}

func (suite *GenerateHandlerTestSuite) TestGenerateHandler_RegisterRoutes_MultiWord() {
	config := suite.getConfig()
	info := compile.ExtractModelInfo(yaml.Model{
		Name:        "PluginVersion",
		Fields:      map[string]yaml.ModelField{"ID": {Type: yaml.ModelFieldTypeUUID}},
		Identifiers: map[string]yaml.ModelIdentifier{"primary": {Fields: []string{"ID"}}},
		Related:     map[string]yaml.ModelRelation{},
	})

	output := compile.GenerateHandler(config, info)

	suite.Contains(output, `rg.GET("/plugin-versions", h.List)`)
	suite.Contains(output, `rg.GET("/plugin-versions/:id", h.GetByID)`)
}

// ── Response mapper ─────────────────────────────────────────────────────────

func (suite *GenerateHandlerTestSuite) TestGenerateHandler_ResponseMapper_DTO() {
	config := suite.getConfig()
	info := compile.ExtractModelInfo(yaml.Model{
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

	output := compile.GenerateHandler(config, info)

	suite.Contains(output, "func toOrganizationResponse(m models.Organization) structures.OrganizationResponse {")
	suite.Contains(output, "structures.OrganizationResponse{")
	suite.Contains(output, "ID: m.ID")
	suite.Contains(output, "Code: m.Code")
	suite.Contains(output, "Name: m.Name")
}

func (suite *GenerateHandlerTestSuite) TestGenerateHandler_ResponseMapper_WithFK() {
	config := suite.getConfig()
	info := compile.ExtractModelInfo(yaml.Model{
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

	output := compile.GenerateHandler(config, info)

	suite.Contains(output, "func toProjectResponse(m models.Project) structures.ProjectResponse {")
	suite.Contains(output, "OrganizationID: m.OrganizationID")
}
