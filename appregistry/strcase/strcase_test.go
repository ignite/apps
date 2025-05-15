package strcase

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_toLowerCamel(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "camel Case",
			arg:  "camel Case",
			want: "camel case",
		},
		{
			name: "snake_case",
			arg:  "snake_case",
			want: "snake case",
		},
		{
			name: "Pascal case",
			arg:  "Pascal case",
			want: "pascal case",
		},
		{
			name: "kebab-case",
			arg:  "kebab-cAse",
			want: "kebab case",
		},
		{
			name: "Title Case",
			arg:  "Title Case",
			want: "title case",
		},
		{
			name: "UPPER CASE",
			arg:  "UPPER CASE",
			want: "upper case",
		},
		{
			name: "single",
			arg:  "single",
			want: "single",
		},
		{
			name: "empty string",
			arg:  "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToLowerCamel(tt.arg)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_toUpperCamel(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "camel Case",
			arg:  "camel Case",
			want: "Camel Case",
		},
		{
			name: "snake_case",
			arg:  "snAke_caSe",
			want: "SnAke CaSe",
		},
		{
			name: "Pascal case",
			arg:  "Pascal Case",
			want: "Pascal Case",
		},
		{
			name: "kebab-case",
			arg:  "keBab-case",
			want: "KeBab Case",
		},
		{
			name: "Title Case",
			arg:  "Title CAse",
			want: "Title CAse",
		},
		{
			name: "upper case",
			arg:  "UPPER CASE",
			want: "UPPER CASE",
		},
		{
			name: "single",
			arg:  "single",
			want: "Single",
		},
		{
			name: "Single with uppers",
			arg:  "singleWithUppers",
			want: "SingleWithUppers",
		},
		{
			name: "empty string",
			arg:  "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToUpperCamel(tt.arg)
			require.Equal(t, tt.want, got)
		})
	}
}
