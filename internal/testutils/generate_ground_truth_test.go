package testutils_test

import (
	"path/filepath"
	"testing"

	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/internal/testutils"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile/cfg"
)

// TestGenerateGroundTruth is a helper test to regenerate ground truth files.
// Run with: go test ./internal/testutils/ -run TestGenerateGroundTruth -v
func TestGenerateGroundTruth(t *testing.T) {
	// To regenerate, temporarily comment out the t.Skip below:
	t.Skip("Run manually to regenerate ground truth: go test ./internal/testutils/ -run TestGenerateGroundTruth -v -count=1")

	testDirPath := testutils.GetTestDirPath()
	groundTruthPath := filepath.Join(testDirPath, "ground-truth", "compile-minimal")

	config := cfg.CompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      filepath.Join(testDirPath, "registry", "minimal", "enums"),
			RegistryModelsDirPath:     filepath.Join(testDirPath, "registry", "minimal", "models"),
			RegistryStructuresDirPath: filepath.Join(testDirPath, "registry", "minimal", "structures"),
			RegistryEntitiesDirPath:   filepath.Join(testDirPath, "registry", "minimal", "entities"),
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
		OutputDirPath:       groundTruthPath,
		StructOutputDirPath: filepath.Join(groundTruthPath, "structures"),
	}

	err := compile.MorpheToGinAPI(config)
	if err != nil {
		t.Fatalf("failed to generate ground truth: %v", err)
	}

	t.Logf("Ground truth generated at: %s", groundTruthPath)
}
