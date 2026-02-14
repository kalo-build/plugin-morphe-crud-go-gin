package compile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/go-util/assertfile"
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/internal/testutils"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile/cfg"
)

type CompileTestSuite struct {
	assertfile.FileSuite

	TestDirPath            string
	TestGroundTruthDirPath string

	ModelsDirPath     string
	EnumsDirPath      string
	StructuresDirPath string
	EntitiesDirPath   string
}

func TestCompileTestSuite(t *testing.T) {
	suite.Run(t, new(CompileTestSuite))
}

func (suite *CompileTestSuite) SetupTest() {
	suite.TestDirPath = testutils.GetTestDirPath()
	suite.TestGroundTruthDirPath = filepath.Join(suite.TestDirPath, "ground-truth", "compile-minimal")

	suite.ModelsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "models")
	suite.EnumsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "enums")
	suite.StructuresDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "structures")
	suite.EntitiesDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "entities")
}

func (suite *CompileTestSuite) TearDownTest() {
	suite.TestDirPath = ""
}

func (suite *CompileTestSuite) TestMorpheToGinAPI() {
	workingDirPath := filepath.Join(suite.TestDirPath, "working")
	suite.Nil(os.Mkdir(workingDirPath, 0755))
	defer os.RemoveAll(workingDirPath)

	structOutputDirPath := filepath.Join(workingDirPath, "structures")

	config := cfg.CompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      suite.EnumsDirPath,
			RegistryModelsDirPath:     suite.ModelsDirPath,
			RegistryStructuresDirPath: suite.StructuresDirPath,
			RegistryEntitiesDirPath:   suite.EntitiesDirPath,
		},
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
		OutputDirPath:       workingDirPath,
		StructOutputDirPath: structOutputDirPath,
	}

	compileErr := compile.MorpheToGinAPI(config)

	suite.NoError(compileErr)

	// ── Verify handlers written directly to output directory ────────────
	// Repo interfaces are generated separately by plugin-morpherepo-go,
	// so this plugin only generates handler files.
	gtHandlersDirPath := filepath.Join(suite.TestGroundTruthDirPath, "handlers")

	handlerOrgPath := filepath.Join(workingDirPath, "organization_handler.go")
	gtHandlerOrgPath := filepath.Join(gtHandlersDirPath, "organization_handler.go")
	suite.FileExists(handlerOrgPath)
	suite.FileEquals(handlerOrgPath, gtHandlerOrgPath)

	handlerProjectPath := filepath.Join(workingDirPath, "project_handler.go")
	gtHandlerProjectPath := filepath.Join(gtHandlersDirPath, "project_handler.go")
	suite.FileExists(handlerProjectPath)
	suite.FileEquals(handlerProjectPath, gtHandlerProjectPath)

	handlerTaskPath := filepath.Join(workingDirPath, "task_handler.go")
	gtHandlerTaskPath := filepath.Join(gtHandlersDirPath, "task_handler.go")
	suite.FileExists(handlerTaskPath)
	suite.FileEquals(handlerTaskPath, gtHandlerTaskPath)

	routesPath := filepath.Join(workingDirPath, "routes.go")
	gtRoutesPath := filepath.Join(gtHandlersDirPath, "routes.go")
	suite.FileExists(routesPath)
	suite.FileEquals(routesPath, gtRoutesPath)

	// ── Verify no helpers.go (isUUID removed) ───────────────────────────
	_, helpersErr := os.Stat(filepath.Join(workingDirPath, "helpers.go"))
	suite.True(os.IsNotExist(helpersErr))

	// ── Verify structures directory ─────────────────────────────────────
	gtStructDirPath := filepath.Join(suite.TestGroundTruthDirPath, "structures")
	suite.DirExists(structOutputDirPath)

	orgStructPath := filepath.Join(structOutputDirPath, "organization_response.str")
	gtOrgStructPath := filepath.Join(gtStructDirPath, "organization_response.str")
	suite.FileExists(orgStructPath)
	suite.FileEquals(orgStructPath, gtOrgStructPath)

	projectStructPath := filepath.Join(structOutputDirPath, "project_response.str")
	gtProjectStructPath := filepath.Join(gtStructDirPath, "project_response.str")
	suite.FileExists(projectStructPath)
	suite.FileEquals(projectStructPath, gtProjectStructPath)

	taskStructPath := filepath.Join(structOutputDirPath, "task_response.str")
	gtTaskStructPath := filepath.Join(gtStructDirPath, "task_response.str")
	suite.FileExists(taskStructPath)
	suite.FileEquals(taskStructPath, gtTaskStructPath)
}

func (suite *CompileTestSuite) TestMorpheToGinAPI_ExcludeModels() {
	workingDirPath := filepath.Join(suite.TestDirPath, "working-exclude")
	suite.Nil(os.Mkdir(workingDirPath, 0755))
	defer os.RemoveAll(workingDirPath)

	structOutputDirPath := filepath.Join(workingDirPath, "structures")

	config := cfg.CompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      suite.EnumsDirPath,
			RegistryModelsDirPath:     suite.ModelsDirPath,
			RegistryStructuresDirPath: suite.StructuresDirPath,
			RegistryEntitiesDirPath:   suite.EntitiesDirPath,
		},
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
		ExcludeModels:       []string{"Task"},
		OutputDirPath:       workingDirPath,
		StructOutputDirPath: structOutputDirPath,
	}

	compileErr := compile.MorpheToGinAPI(config)

	suite.NoError(compileErr)

	// Organization and Project handlers should exist (written directly to output dir)
	suite.FileExists(filepath.Join(workingDirPath, "organization_handler.go"))
	suite.FileExists(filepath.Join(workingDirPath, "project_handler.go"))

	suite.FileExists(filepath.Join(structOutputDirPath, "organization_response.str"))
	suite.FileExists(filepath.Join(structOutputDirPath, "project_response.str"))

	// Task should NOT exist
	_, taskHandlerErr := os.Stat(filepath.Join(workingDirPath, "task_handler.go"))
	suite.True(os.IsNotExist(taskHandlerErr))
	_, taskStructErr := os.Stat(filepath.Join(structOutputDirPath, "task_response.str"))
	suite.True(os.IsNotExist(taskStructErr))

	// Routes should reference only Organization and Project
	routesPath := filepath.Join(workingDirPath, "routes.go")
	suite.FileExists(routesPath)
	routesContent, readErr := os.ReadFile(routesPath)
	suite.NoError(readErr)
	suite.Contains(string(routesContent), "Organization")
	suite.Contains(string(routesContent), "Project")
	suite.NotContains(string(routesContent), "Task")
}
