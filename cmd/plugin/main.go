package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile/cfg"
)

type CompileConfigEntry struct {
	PackagePath string `json:"PackagePath"`
}

type PluginConfig struct {
	InputPath  string             `json:"inputPath"`
	OutputPath string             `json:"outputPath"`
	Config     PluginConfigFields `json:"config"`
	Verbose    bool               `json:"verbose,omitempty"`
}

type PluginConfigFields struct {
	Handlers      CompileConfigEntry `json:"handlers"`
	Repo          CompileConfigEntry `json:"repo"`
	Models        CompileConfigEntry `json:"models"`
	Structures    CompileConfigEntry `json:"structures"`
	ExcludeModels []string           `json:"excludeModels,omitempty"`
}

const (
	ErrMissingConfig       = 3
	ErrInvalidConfig       = 4
	ErrInputPathRequired   = 12
	ErrOutputPathRequired  = 13
	ErrPackagePathRequired = 14
	ErrCompileFailed       = 1
)

func logInfo(verbose bool, format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stdout, format+"\n", args...)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: plugin-morphe-crud-go-gin <config>")
		fmt.Fprintln(os.Stderr, "  config: JSON string with inputPath, outputPath, and config parameters")
		os.Exit(ErrMissingConfig)
	}

	rawConfig := os.Args[1]
	var pluginConfig PluginConfig
	if err := json.Unmarshal([]byte(rawConfig), &pluginConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config JSON:", err)
		os.Exit(ErrInvalidConfig)
	}

	if pluginConfig.InputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Input path is required")
		os.Exit(ErrInputPathRequired)
	}

	if pluginConfig.OutputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Output path is required")
		os.Exit(ErrOutputPathRequired)
	}

	if pluginConfig.Config.Handlers.PackagePath == "" {
		fmt.Fprintln(os.Stderr, "Error: Handlers package path is required")
		os.Exit(ErrPackagePathRequired)
	}

	if pluginConfig.Config.Repo.PackagePath == "" {
		fmt.Fprintln(os.Stderr, "Error: Repo package path is required")
		os.Exit(ErrPackagePathRequired)
	}

	if pluginConfig.Config.Models.PackagePath == "" {
		fmt.Fprintln(os.Stderr, "Error: Models package path is required")
		os.Exit(ErrPackagePathRequired)
	}

	inputAbs, err := filepath.Abs(pluginConfig.InputPath)
	if err == nil {
		pluginConfig.InputPath = inputAbs
	}

	outputAbs, err := filepath.Abs(pluginConfig.OutputPath)
	if err == nil {
		pluginConfig.OutputPath = outputAbs
	}

	logInfo(pluginConfig.Verbose, "Processing Morphe registry from: '%s'", pluginConfig.InputPath)
	logInfo(pluginConfig.Verbose, "Output Go API to: '%s'", pluginConfig.OutputPath)

	compileConfig := cfg.CompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      filepath.Join(pluginConfig.InputPath, "enums"),
			RegistryModelsDirPath:     filepath.Join(pluginConfig.InputPath, "models"),
			RegistryStructuresDirPath: filepath.Join(pluginConfig.InputPath, "structures"),
			RegistryEntitiesDirPath:   filepath.Join(pluginConfig.InputPath, "entities"),
		},
		HandlersConfig: cfg.HandlersConfig{
			PackagePath: pluginConfig.Config.Handlers.PackagePath,
		},
		RepoConfig: cfg.RepoConfig{
			PackagePath: pluginConfig.Config.Repo.PackagePath,
		},
		ModelsConfig: cfg.ModelsConfig{
			PackagePath: pluginConfig.Config.Models.PackagePath,
		},
		StructuresConfig: cfg.StructuresConfig{
			PackagePath: pluginConfig.Config.Structures.PackagePath,
		},
		ExcludeModels:       pluginConfig.Config.ExcludeModels,
		OutputDirPath:       pluginConfig.OutputPath,
		StructOutputDirPath: filepath.Join(pluginConfig.InputPath, "structures"),
	}

	logInfo(pluginConfig.Verbose, "Starting compilation process...")
	compileErr := compile.MorpheToGinAPI(compileConfig)
	if compileErr != nil {
		fmt.Fprintln(os.Stderr, "Compilation failed:", compileErr)
		os.Exit(ErrCompileFailed)
	}

	logInfo(pluginConfig.Verbose, "Compilation completed successfully")
	os.Exit(0)
}
