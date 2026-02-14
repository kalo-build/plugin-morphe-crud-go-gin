package cfg_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile/cfg"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

// ── PackageNameFromPath ─────────────────────────────────────────────────────

func (suite *ConfigTestSuite) TestPackageNameFromPath_FullPath() {
	suite.Equal("models", cfg.PackageNameFromPath("github.com/org/app/internal/types/models"))
}

func (suite *ConfigTestSuite) TestPackageNameFromPath_ShortPath() {
	suite.Equal("repo", cfg.PackageNameFromPath("github.com/org/app/repo"))
}

func (suite *ConfigTestSuite) TestPackageNameFromPath_SingleSegment() {
	suite.Equal("models", cfg.PackageNameFromPath("models"))
}

func (suite *ConfigTestSuite) TestPackageNameFromPath_Empty() {
	suite.Equal("", cfg.PackageNameFromPath(""))
}

// ── Resolve ─────────────────────────────────────────────────────────────────

func (suite *ConfigTestSuite) TestResolve_DerivesPackageNames() {
	config := cfg.CompileConfig{
		HandlersConfig: cfg.HandlersConfig{
			PackagePath: "github.com/org/app/internal/generated/handlers",
		},
		RepoConfig: cfg.RepoConfig{
			PackagePath: "github.com/org/app/internal/generated/repo",
		},
		ModelsConfig: cfg.ModelsConfig{
			PackagePath: "github.com/org/app/internal/types/models",
		},
		StructuresConfig: cfg.StructuresConfig{
			PackagePath: "github.com/org/app/internal/types/structures",
		},
	}

	config.Resolve()

	suite.Equal("handlers", config.HandlersConfig.PackageName)
	suite.Equal("repo", config.RepoConfig.PackageName)
	suite.Equal("models", config.ModelsConfig.PackageName)
	suite.Equal("structures", config.StructuresConfig.PackageName)
}

func (suite *ConfigTestSuite) TestResolve_PreservesExplicitPackageNames() {
	config := cfg.CompileConfig{
		HandlersConfig: cfg.HandlersConfig{
			PackagePath: "github.com/org/app/internal/generated/handlers",
			PackageName: "myhandlers",
		},
		RepoConfig: cfg.RepoConfig{
			PackagePath: "github.com/org/app/internal/generated/repo",
			PackageName: "myrepo",
		},
		ModelsConfig: cfg.ModelsConfig{
			PackagePath: "github.com/org/app/internal/types/models",
			PackageName: "mymodels",
		},
	}

	config.Resolve()

	suite.Equal("myhandlers", config.HandlersConfig.PackageName)
	suite.Equal("myrepo", config.RepoConfig.PackageName)
	suite.Equal("mymodels", config.ModelsConfig.PackageName)
}
