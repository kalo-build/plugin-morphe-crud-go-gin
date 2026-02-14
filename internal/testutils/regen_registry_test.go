package testutils_test

import (
	"path/filepath"
	"testing"

	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile/cfg"
)

func TestRegenRegistry(t *testing.T) {
	t.Skip("Run manually to regenerate registry generated files")

	registryRoot := filepath.Join("..", "..", "..", "kalo-plugin-registry", "morphe", "registry")
	outputPath := filepath.Join("..", "..", "..", "kalo-plugin-registry", "internal", "generated")
	structOutputPath := filepath.Join(registryRoot, "structures")

	config := cfg.CompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      filepath.Join(registryRoot, "enums"),
			RegistryModelsDirPath:     filepath.Join(registryRoot, "models"),
			RegistryStructuresDirPath: filepath.Join(registryRoot, "structures"),
			RegistryEntitiesDirPath:   filepath.Join(registryRoot, "entities"),
		},
		HandlersConfig: cfg.HandlersConfig{
			PackagePath: "github.com/kalo-build/kalo-plugin-registry/internal/generated/handlers",
		},
		RepoConfig: cfg.RepoConfig{
			PackagePath: "github.com/kalo-build/kalo-plugin-registry/internal/generated/repo",
		},
		ModelsConfig: cfg.ModelsConfig{
			PackagePath: "github.com/kalo-build/kalo-plugin-registry/internal/types/models",
		},
		StructuresConfig: cfg.StructuresConfig{
			PackagePath: "github.com/kalo-build/kalo-plugin-registry/internal/types/structures",
		},
		ExcludeModels:       []string{"User"},
		OutputDirPath:       outputPath,
		StructOutputDirPath: structOutputPath,
	}

	if err := compile.MorpheToGinAPI(config); err != nil {
		t.Fatal(err)
	}

	t.Log("Registry generated files written to:", outputPath)
	t.Log("Registry .str files written to:", structOutputPath)
}
