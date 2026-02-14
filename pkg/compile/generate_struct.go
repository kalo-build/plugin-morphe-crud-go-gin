package compile

import (
	"fmt"
	"strings"
)

// GenerateResponseStruct generates a .str YAML file content for a model's response DTO.
func GenerateResponseStruct(info ModelInfo) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("name: %sResponse\n", info.Name))
	b.WriteString("fields:\n")

	for _, field := range info.Fields {
		fieldType := string(field.Type)
		if field.IsFK {
			// FK fields are UUID references
			fieldType = "UUID"
		}
		if field.IsPoly {
			// Polymorphic type fields are strings
			fieldType = "String"
		}
		b.WriteString(fmt.Sprintf("  %s:\n", field.Name))
		b.WriteString(fmt.Sprintf("    type: %s\n", fieldType))
		if field.IsOptional {
			b.WriteString("    attributes:\n")
			b.WriteString("      - optional\n")
		}
	}

	return b.String()
}
