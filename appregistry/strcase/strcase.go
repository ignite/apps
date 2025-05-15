package strcase

import (
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
)

func ToKebab(s string) string {
	return strcase.ToKebab(s)
}

func ToSnake(s string) string {
	return strcase.ToSnake(s)
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func ToUpperCamel(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			for j, r := range runes {
				if j == 0 {
					runes[j] = unicode.ToUpper(r)
				} else {
					runes[j] = unicode.ToLower(r)
				}
			}
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

func ToLowerCamel(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	return strings.ToLower(s)
}
