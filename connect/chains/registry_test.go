package chains

import (
	"reflect"
	"testing"

	"github.com/ignite/cli/v29/ignite/pkg/chainregistry"
)

func TestCleanGRPCEntries(t *testing.T) {
	tests := []struct {
		name     string
		entries  []chainregistry.APIProvider
		expected []chainregistry.APIProvider
	}{
		{
			name: "Clean entries with http:// and https://",
			entries: []chainregistry.APIProvider{
				{Address: "http://example1.com:1234/"},
				{Address: "https://example2.com:5678"},
				{Address: "http://example3.com:9012"},
			},
			expected: []chainregistry.APIProvider{
				{Address: "example1.com:1234"},
				{Address: "example2.com:5678"},
				{Address: "example3.com:9012"},
			},
		},
		{
			name: "Remove trailing slashes",
			entries: []chainregistry.APIProvider{
				{Address: "http://example.com:1234/"},
				{Address: "https://example.com:5678/"},
			},
			expected: []chainregistry.APIProvider{
				{Address: "example.com:1234"},
				{Address: "example.com:5678"},
			},
		},
		{
			name: "Filter entries without ports",
			entries: []chainregistry.APIProvider{
				{Address: "example.com"},
				{Address: "example.com:1234"},
			},
			expected: []chainregistry.APIProvider{
				{Address: "example.com:1234"},
			},
		},
		{
			name:     "Empty input",
			entries:  []chainregistry.APIProvider{},
			expected: []chainregistry.APIProvider{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Copy the entries so that modifications in one test case
			// don't affect others
			input := make([]chainregistry.APIProvider, len(tt.entries))
			copy(input, tt.entries)

			result := cleanGRPCEntries(input)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("cleanGRPCEntries() = %v, want %v", result, tt.expected)
			}
		})
	}
}
