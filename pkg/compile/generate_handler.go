package compile

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/naming"
)

// GenerateHandler generates a Go Gin handler file for a model.
func GenerateHandler(config cfg.CompileConfig, info ModelInfo) string {
	var b strings.Builder

	repoImport := config.RepoConfig.PackagePath
	repoPkg := config.RepoConfig.PackageName
	modelsPkg := config.ModelsConfig.PackageName
	modelsImport := config.ModelsConfig.PackagePath
	structuresPkg := config.StructuresConfig.PackageName
	structuresImport := config.StructuresConfig.PackagePath
	handlersPkg := config.HandlersConfig.PackageName

	// Package declaration
	b.WriteString(fmt.Sprintf("package %s\n\n", handlersPkg))

	// Imports
	b.WriteString("import (\n")
	b.WriteString("\t\"net/http\"\n\n")
	b.WriteString("\t\"github.com/gin-gonic/gin\"\n")
	b.WriteString(fmt.Sprintf("\t\"%s\"\n", modelsImport))
	b.WriteString(fmt.Sprintf("\t\"%s\"\n", repoImport))
	b.WriteString(fmt.Sprintf("\t\"%s\"\n", structuresImport))
	b.WriteString(")\n\n")

	// Handler struct
	b.WriteString(fmt.Sprintf("// %sHandler handles HTTP requests for %s resources.\n", info.Name, info.Name))
	b.WriteString(fmt.Sprintf("type %sHandler struct {\n", info.Name))
	b.WriteString(fmt.Sprintf("\trepo %s.%sRepository\n", repoPkg, info.Name))
	b.WriteString("}\n\n")

	// Constructor
	b.WriteString(fmt.Sprintf("// New%sHandler creates a new %sHandler.\n", info.Name, info.Name))
	b.WriteString(fmt.Sprintf("func New%sHandler(r %s.%sRepository) *%sHandler {\n", info.Name, repoPkg, info.Name, info.Name))
	b.WriteString(fmt.Sprintf("\treturn &%sHandler{repo: r}\n", info.Name))
	b.WriteString("}\n\n")

	// List handler
	writeListHandler(&b, info, structuresPkg)

	// Per-identifier Get handlers
	writeGetByIdentifierHandlers(&b, info, structuresPkg)

	// Create handler
	writeCreateHandler(&b, info, modelsPkg, structuresPkg)

	// Update handler
	writeUpdateHandler(&b, info, modelsPkg, structuresPkg)

	// Delete handler
	writeDeleteHandler(&b, info)

	// RegisterRoutes
	writeRegisterRoutes(&b, info)

	// Response mapper
	writeResponseMapper(&b, info, modelsPkg, structuresPkg)

	return b.String()
}

func writeListHandler(b *strings.Builder, info ModelInfo, structuresPkg string) {
	b.WriteString(fmt.Sprintf("// List handles GET /%s\n", info.CollectionName))
	b.WriteString(fmt.Sprintf("func (h *%sHandler) List(c *gin.Context) {\n", info.Name))
	b.WriteString("\tctx := c.Request.Context()\n\n")

	// Parse optional filter params from query string
	for _, filter := range info.Filters {
		b.WriteString(fmt.Sprintf("\tvar %s *string\n", filter.GoFieldName))
		b.WriteString(fmt.Sprintf("\tif v := c.Query(\"%s\"); v != \"\" {\n", filter.ParamName))
		b.WriteString(fmt.Sprintf("\t\t%s = &v\n", filter.GoFieldName))
		b.WriteString("\t}\n\n")
	}

	// Call repo
	repoArgs := "ctx"
	for _, filter := range info.Filters {
		repoArgs += ", " + filter.GoFieldName
	}
	b.WriteString(fmt.Sprintf("\titems, err := h.repo.GetAll(%s)\n", repoArgs))
	b.WriteString("\tif err != nil {\n")
	b.WriteString(fmt.Sprintf("\t\tc.JSON(http.StatusInternalServerError, gin.H{\"error\": \"Failed to fetch %s\"})\n", info.CollectionName))
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n\n")

	// Build response using DTO structs
	b.WriteString(fmt.Sprintf("\tresponses := make([]%s.%sResponse, len(items))\n", structuresPkg, info.Name))
	b.WriteString("\tfor i, item := range items {\n")
	b.WriteString(fmt.Sprintf("\t\tresponses[i] = to%sResponse(item)\n", info.Name))
	b.WriteString("\t}\n\n")

	b.WriteString("\tc.JSON(http.StatusOK, gin.H{\"data\": responses})\n")
	b.WriteString("}\n\n")
}

func writeGetByIdentifierHandlers(b *strings.Builder, info ModelInfo, structuresPkg string) {
	// Primary identifier: GetByID
	b.WriteString(fmt.Sprintf("// GetByID handles GET /%s/:id\n", info.CollectionName))
	b.WriteString(fmt.Sprintf("func (h *%sHandler) GetByID(c *gin.Context) {\n", info.Name))
	b.WriteString("\tctx := c.Request.Context()\n")
	b.WriteString("\tid := c.Param(\"id\")\n\n")

	b.WriteString("\tif id == \"\" {\n")
	b.WriteString(fmt.Sprintf("\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": \"%s ID is required\"})\n", info.Name))
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n\n")

	b.WriteString("\titem, err := h.repo.GetByID(ctx, id)\n")
	writeGetErrorHandling(b, info)

	b.WriteString(fmt.Sprintf("\tc.JSON(http.StatusOK, gin.H{\"data\": to%sResponse(*item)})\n", info.Name))
	b.WriteString("}\n\n")

	// Secondary identifiers, sorted for deterministic output
	secNames := make([]string, 0, len(info.SecondaryIdentifiers))
	for name := range info.SecondaryIdentifiers {
		secNames = append(secNames, name)
	}
	sort.Strings(secNames)

	for _, idName := range secNames {
		fields := info.SecondaryIdentifiers[idName]
		if len(fields) == 0 {
			continue
		}
		paramName := naming.ToKebabCase(idName)
		methodName := "GetBy" + strings.ToUpper(idName[:1]) + idName[1:]
		fieldParam := strings.ToLower(fields[0][:1]) + fields[0][1:]
		// Handle all-uppercase field names like "ID"
		if strings.ToUpper(fields[0]) == fields[0] {
			fieldParam = strings.ToLower(fields[0])
		}

		b.WriteString(fmt.Sprintf("// %s handles GET /%s/by-%s/:%s\n", methodName, info.CollectionName, paramName, fieldParam))
		b.WriteString(fmt.Sprintf("func (h *%sHandler) %s(c *gin.Context) {\n", info.Name, methodName))
		b.WriteString("\tctx := c.Request.Context()\n")
		b.WriteString(fmt.Sprintf("\t%s := c.Param(\"%s\")\n\n", fieldParam, fieldParam))

		b.WriteString(fmt.Sprintf("\tif %s == \"\" {\n", fieldParam))
		b.WriteString(fmt.Sprintf("\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": \"%s %s is required\"})\n", info.Name, fields[0]))
		b.WriteString("\t\treturn\n")
		b.WriteString("\t}\n\n")

		b.WriteString(fmt.Sprintf("\titem, err := h.repo.%s(ctx, %s)\n", methodName, fieldParam))
		writeGetErrorHandling(b, info)

		b.WriteString(fmt.Sprintf("\tc.JSON(http.StatusOK, gin.H{\"data\": to%sResponse(*item)})\n", info.Name))
		b.WriteString("}\n\n")
	}
}

func writeGetErrorHandling(b *strings.Builder, info ModelInfo) {
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\tif err.Error() == \"no rows in result set\" {\n")
	b.WriteString(fmt.Sprintf("\t\t\tc.JSON(http.StatusNotFound, gin.H{\"error\": \"%s not found\"})\n", info.Name))
	b.WriteString("\t\t} else {\n")
	b.WriteString(fmt.Sprintf("\t\t\tc.JSON(http.StatusInternalServerError, gin.H{\"error\": \"Failed to fetch %s\"})\n", naming.ToSnakeCase(info.Name)))
	b.WriteString("\t\t}\n")
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n\n")
}

func writeCreateHandler(b *strings.Builder, info ModelInfo, modelsPkg string, structuresPkg string) {
	b.WriteString(fmt.Sprintf("// Create handles POST /%s\n", info.CollectionName))
	b.WriteString(fmt.Sprintf("func (h *%sHandler) Create(c *gin.Context) {\n", info.Name))
	b.WriteString("\tctx := c.Request.Context()\n\n")

	b.WriteString(fmt.Sprintf("\tvar input %s.%s\n", modelsPkg, info.Name))
	b.WriteString("\tif err := c.ShouldBindJSON(&input); err != nil {\n")
	b.WriteString("\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": err.Error()})\n")
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n\n")

	b.WriteString("\tcreated, err := h.repo.Create(ctx, &input)\n")
	b.WriteString("\tif err != nil {\n")
	b.WriteString(fmt.Sprintf("\t\tc.JSON(http.StatusInternalServerError, gin.H{\"error\": \"Failed to create %s\"})\n", naming.ToSnakeCase(info.Name)))
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n\n")

	b.WriteString(fmt.Sprintf("\tc.JSON(http.StatusCreated, gin.H{\"data\": to%sResponse(*created)})\n", info.Name))
	b.WriteString("}\n\n")
}

func writeUpdateHandler(b *strings.Builder, info ModelInfo, modelsPkg string, structuresPkg string) {
	b.WriteString(fmt.Sprintf("// Update handles PUT /%s/:id\n", info.CollectionName))
	b.WriteString(fmt.Sprintf("func (h *%sHandler) Update(c *gin.Context) {\n", info.Name))
	b.WriteString("\tctx := c.Request.Context()\n")
	b.WriteString("\tid := c.Param(\"id\")\n\n")

	b.WriteString("\tif id == \"\" {\n")
	b.WriteString(fmt.Sprintf("\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": \"%s ID is required\"})\n", info.Name))
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n\n")

	b.WriteString(fmt.Sprintf("\tvar input %s.%s\n", modelsPkg, info.Name))
	b.WriteString("\tif err := c.ShouldBindJSON(&input); err != nil {\n")
	b.WriteString("\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": err.Error()})\n")
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n\n")

	b.WriteString("\tupdated, err := h.repo.Update(ctx, id, &input)\n")
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\tif err.Error() == \"no rows in result set\" {\n")
	b.WriteString(fmt.Sprintf("\t\t\tc.JSON(http.StatusNotFound, gin.H{\"error\": \"%s not found\"})\n", info.Name))
	b.WriteString("\t\t} else {\n")
	b.WriteString(fmt.Sprintf("\t\t\tc.JSON(http.StatusInternalServerError, gin.H{\"error\": \"Failed to update %s\"})\n", naming.ToSnakeCase(info.Name)))
	b.WriteString("\t\t}\n")
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n\n")

	b.WriteString(fmt.Sprintf("\tc.JSON(http.StatusOK, gin.H{\"data\": to%sResponse(*updated)})\n", info.Name))
	b.WriteString("}\n\n")
}

func writeDeleteHandler(b *strings.Builder, info ModelInfo) {
	b.WriteString(fmt.Sprintf("// Delete handles DELETE /%s/:id\n", info.CollectionName))
	b.WriteString(fmt.Sprintf("func (h *%sHandler) Delete(c *gin.Context) {\n", info.Name))
	b.WriteString("\tctx := c.Request.Context()\n")
	b.WriteString("\tid := c.Param(\"id\")\n\n")

	b.WriteString("\tif id == \"\" {\n")
	b.WriteString(fmt.Sprintf("\t\tc.JSON(http.StatusBadRequest, gin.H{\"error\": \"%s ID is required\"})\n", info.Name))
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n\n")

	b.WriteString("\terr := h.repo.Delete(ctx, id)\n")
	b.WriteString("\tif err != nil {\n")
	b.WriteString("\t\tif err.Error() == \"no rows in result set\" {\n")
	b.WriteString(fmt.Sprintf("\t\t\tc.JSON(http.StatusNotFound, gin.H{\"error\": \"%s not found\"})\n", info.Name))
	b.WriteString("\t\t} else {\n")
	b.WriteString(fmt.Sprintf("\t\t\tc.JSON(http.StatusInternalServerError, gin.H{\"error\": \"Failed to delete %s\"})\n", naming.ToSnakeCase(info.Name)))
	b.WriteString("\t\t}\n")
	b.WriteString("\t\treturn\n")
	b.WriteString("\t}\n\n")

	b.WriteString("\tc.JSON(http.StatusNoContent, nil)\n")
	b.WriteString("}\n\n")
}

func writeRegisterRoutes(b *strings.Builder, info ModelInfo) {
	b.WriteString(fmt.Sprintf("// RegisterRoutes registers all CRUD routes for %s.\n", info.Name))
	b.WriteString(fmt.Sprintf("func (h *%sHandler) RegisterRoutes(rg *gin.RouterGroup) {\n", info.Name))
	b.WriteString(fmt.Sprintf("\trg.GET(\"/%s\", h.List)\n", info.CollectionName))
	b.WriteString(fmt.Sprintf("\trg.POST(\"/%s\", h.Create)\n", info.CollectionName))
	b.WriteString(fmt.Sprintf("\trg.GET(\"/%s/:id\", h.GetByID)\n", info.CollectionName))
	b.WriteString(fmt.Sprintf("\trg.PUT(\"/%s/:id\", h.Update)\n", info.CollectionName))
	b.WriteString(fmt.Sprintf("\trg.DELETE(\"/%s/:id\", h.Delete)\n", info.CollectionName))

	// Secondary identifier routes
	secNames := make([]string, 0, len(info.SecondaryIdentifiers))
	for name := range info.SecondaryIdentifiers {
		secNames = append(secNames, name)
	}
	sort.Strings(secNames)

	for _, idName := range secNames {
		fields := info.SecondaryIdentifiers[idName]
		if len(fields) == 0 {
			continue
		}
		paramName := naming.ToKebabCase(idName)
		methodName := "GetBy" + strings.ToUpper(idName[:1]) + idName[1:]
		fieldParam := strings.ToLower(fields[0][:1]) + fields[0][1:]
		if strings.ToUpper(fields[0]) == fields[0] {
			fieldParam = strings.ToLower(fields[0])
		}
		b.WriteString(fmt.Sprintf("\trg.GET(\"/%s/by-%s/:%s\", h.%s)\n", info.CollectionName, paramName, fieldParam, methodName))
	}

	b.WriteString("}\n")
}

func writeResponseMapper(b *strings.Builder, info ModelInfo, modelsPkg string, structuresPkg string) {
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("// to%sResponse converts a %s model to a %sResponse DTO.\n", info.Name, info.Name, info.Name))
	b.WriteString(fmt.Sprintf("func to%sResponse(m %s.%s) %s.%sResponse {\n", info.Name, modelsPkg, info.Name, structuresPkg, info.Name))
	b.WriteString(fmt.Sprintf("\treturn %s.%sResponse{\n", structuresPkg, info.Name))

	for _, field := range info.Fields {
		b.WriteString(fmt.Sprintf("\t\t%s: m.%s,\n", field.Name, field.Name))
	}

	b.WriteString("\t}\n")
	b.WriteString("}\n")
}
