package naming

import (
	"strings"
	"unicode"
)

// ToSnakeCase converts PascalCase to snake_case.
// e.g., "PluginVersion" -> "plugin_version"
func ToSnakeCase(s string) string {
	return splitAndJoin(s, "_")
}

// ToKebabCase converts PascalCase to kebab-case.
// e.g., "PluginVersion" -> "plugin-version"
func ToKebabCase(s string) string {
	return splitAndJoin(s, "-")
}

// Pluralize adds a simple "s" suffix.
// Handles common cases for typical model names.
func Pluralize(s string) string {
	if s == "" {
		return s
	}
	// Handle common English pluralization patterns
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "z") {
		return s + "es"
	}
	if strings.HasSuffix(s, "y") && len(s) > 1 {
		prev := s[len(s)-2]
		if prev != 'a' && prev != 'e' && prev != 'i' && prev != 'o' && prev != 'u' {
			return s[:len(s)-1] + "ies"
		}
	}
	return s + "s"
}

// CollectionName converts a PascalCase model name to a URL-friendly pluralized kebab-case collection name.
// e.g., "PluginVersion" -> "plugin-versions"
func CollectionName(modelName string) string {
	return Pluralize(ToKebabCase(modelName))
}

// QueryParamName converts a PascalCase relationship name to a snake_case query parameter with _id suffix.
// e.g., "Organization" -> "organization_id"
func QueryParamName(relName string) string {
	return ToSnakeCase(relName) + "_id"
}

// FilterFieldName converts a relationship name to a Go parameter name for filters.
// e.g., "Organization" -> "organizationID"
func FilterFieldName(relName string) string {
	return lowerFirst(relName) + "ID"
}

// splitAndJoin splits PascalCase into words and joins with the given separator.
func splitAndJoin(s string, sep string) string {
	if s == "" {
		return s
	}

	var words []string
	wordStart := 0

	runes := []rune(s)
	for i := 1; i < len(runes); i++ {
		if unicode.IsUpper(runes[i]) {
			// Check if this starts a new word
			// Handle sequences of uppercase (e.g., "ID" stays together, "IDPrimary" splits as "ID" + "Primary")
			if unicode.IsLower(runes[i-1]) {
				words = append(words, string(runes[wordStart:i]))
				wordStart = i
			} else if i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				words = append(words, string(runes[wordStart:i]))
				wordStart = i
			}
		}
	}
	words = append(words, string(runes[wordStart:]))

	for i, w := range words {
		words[i] = strings.ToLower(w)
	}

	return strings.Join(words, sep)
}

// lowerFirst lowercases the first character of a string.
func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
