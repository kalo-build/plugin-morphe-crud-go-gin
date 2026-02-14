package compile

import (
	"fmt"
	"strings"

	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile/cfg"
)

// GenerateRoutes generates a routes.go file that registers all handlers.
func GenerateRoutes(config cfg.CompileConfig, models []ModelInfo) string {
	var b strings.Builder

	repoPkg := config.RepoConfig.PackageName
	repoImport := config.RepoConfig.PackagePath
	handlersPkg := config.HandlersConfig.PackageName

	// Package declaration
	b.WriteString(fmt.Sprintf("package %s\n\n", handlersPkg))

	// Imports
	b.WriteString("import (\n")
	b.WriteString("\t\"github.com/gin-gonic/gin\"\n")
	b.WriteString(fmt.Sprintf("\t\"%s\"\n", repoImport))
	b.WriteString(")\n\n")

	// Repos struct
	b.WriteString("// Repos holds all repository implementations needed by the generated handlers.\n")
	b.WriteString("type Repos struct {\n")
	for _, m := range models {
		b.WriteString(fmt.Sprintf("\t%s %s.%sRepository\n", m.Name, repoPkg, m.Name))
	}
	b.WriteString("}\n\n")

	// RegisterAllRoutes function
	b.WriteString("// RegisterAllRoutes registers all generated CRUD handlers on the given router group.\n")
	b.WriteString("func RegisterAllRoutes(rg *gin.RouterGroup, repos Repos) {\n")
	for _, m := range models {
		b.WriteString(fmt.Sprintf("\tNew%sHandler(repos.%s).RegisterRoutes(rg)\n", m.Name, m.Name))
	}
	b.WriteString("}\n")

	return b.String()
}
