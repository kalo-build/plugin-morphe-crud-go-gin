package compile

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-crud-go-gin/pkg/naming"
)

// MorpheToGinAPI is the main compilation entrypoint.
// It reads Morphe models and generates Go Gin handlers, route registration,
// and response DTO .struct files.
func MorpheToGinAPI(config cfg.CompileConfig) error {
	config.Resolve()

	// Load Morphe registry
	r, rErr := registry.LoadMorpheRegistry(registry.LoadMorpheRegistryHooks{}, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return fmt.Errorf("failed to load morphe registry: %w", rErr)
	}

	// Build exclude set
	excludeSet := make(map[string]bool)
	for _, name := range config.ExcludeModels {
		excludeSet[name] = true
	}

	// Extract model info for each model
	var modelInfos []ModelInfo
	allModels := r.GetAllModels()

	// Sort model names for deterministic output
	modelNames := make([]string, 0, len(allModels))
	for name := range allModels {
		modelNames = append(modelNames, name)
	}
	sort.Strings(modelNames)

	for _, name := range modelNames {
		if excludeSet[name] {
			continue
		}
		model := allModels[name]
		info := ExtractModelInfo(model)
		modelInfos = append(modelInfos, info)
	}

	if len(modelInfos) == 0 {
		return fmt.Errorf("no models to generate (all excluded or none found)")
	}

	// Generate handlers directly into the output directory.
	// Repo interfaces are generated separately by plugin-morpherepo-go.
	handlersDirPath := config.OutputDirPath
	if err := os.MkdirAll(handlersDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create handlers directory: %w", err)
	}

	for _, info := range modelInfos {
		code := GenerateHandler(config, info)
		if err := writeFormattedGoFile(handlersDirPath, naming.ToSnakeCase(info.Name)+"_handler.go", code); err != nil {
			return fmt.Errorf("failed to write handler for %s: %w", info.Name, err)
		}
	}

	// Generate routes
	routesCode := GenerateRoutes(config, modelInfos)
	if err := writeFormattedGoFile(handlersDirPath, "routes.go", routesCode); err != nil {
		return fmt.Errorf("failed to write routes: %w", err)
	}

	// Generate response DTO .struct files (if output dir configured)
	if config.StructOutputDirPath != "" {
		if err := os.MkdirAll(config.StructOutputDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create struct output directory: %w", err)
		}
		for _, info := range modelInfos {
			content := GenerateResponseStruct(info)
			fileName := naming.ToSnakeCase(info.Name) + "_response.str"
			filePath := filepath.Join(config.StructOutputDirPath, fileName)
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write struct file for %s: %w", info.Name, err)
			}
		}
	}

	return nil
}

// writeFormattedGoFile formats Go source code and writes it to a file.
func writeFormattedGoFile(dirPath string, fileName string, code string) error {
	formatted, err := format.Source([]byte(code))
	if err != nil {
		// Write unformatted for debugging
		debugPath := filepath.Join(dirPath, fileName+".unformatted")
		_ = os.WriteFile(debugPath, []byte(code), 0644)
		return fmt.Errorf("format error in %s: %w", fileName, err)
	}

	filePath := filepath.Join(dirPath, fileName)
	return os.WriteFile(filePath, formatted, 0644)
}
