package databaseHelper

import (
	"strings"
	"unicode"
)

// CamelToSnakeMap converts all keys in a map from camelCase to snake_case.
// This is needed because JSON unmarshal produces camelCase keys (from json tags),
// but GORM's Updates(map) expects snake_case column names.
func CamelToSnakeMap(input map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{}, len(input))
	for k, v := range input {
		output[camelToSnake(k)] = v
	}
	return output
}

func camelToSnake(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
