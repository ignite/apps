package strcase

import (
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
)

// ToKebab converts a string to kebab case (lowercase with hyphens).
// Example: "camelCase" -> "camel-case".
func ToKebab(s string) string {
	return strcase.ToKebab(s)
}

// ToSnake converts a string to snake case (lowercase with underscores).
// Example: "camelCase" -> "camel_case".
func ToSnake(s string) string {
	return strcase.ToSnake(s)
}

// ToLower converts a string to all lowercase.
// Example: "HELLO" -> "hello".
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToUpper converts a string to all uppercase.
// Example: "hello" -> "HELLO".
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ToUpperCamel converts a string into upper camel case with spaces preserved between logical words.
// It only modifies the first rune of each word, keeping the rest intact (case-preserving).
// Example: "hello-world" -> "Hello World".
// Example: "snAke_caSe" -> "SnAke CaSe".
func ToUpperCamel(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) == 0 {
			continue
		}
		runes := []rune(word)
		runes[0] = unicode.ToUpper(runes[0])
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// ToLowerCamel converts a string into lower camel case with spaces between words.
// All characters are converted to lowercase. Hyphens and underscores are replaced with spaces.
// Example: "Hello-World" -> "hello world".
// Example: "UPPER_CASE" -> "upper case".
func ToLowerCamel(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	return strings.ToLower(s)
}
