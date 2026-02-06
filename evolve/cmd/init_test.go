package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestPrependFieldToGenesis(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		fieldName     string
		fieldValue    string
		expectedStart string
		wantErr       bool
	}{
		{
			name: "prepend to non-empty genesis",
			input: `{
  "genesis_time": "2024-01-01T00:00:00Z",
  "chain_id": "test-chain"
}`,
			fieldName:  "da_epoch_forced_inclusion",
			fieldValue: "0",
			expectedStart: `{
  "da_epoch_forced_inclusion": 0,
  "genesis_time": "2024-01-01T00:00:00Z",`,
			wantErr: false,
		},
		{
			name:       "prepend to empty genesis",
			input:      `{}`,
			fieldName:  "da_epoch_forced_inclusion",
			fieldValue: "0",
			expectedStart: `{
  "da_epoch_forced_inclusion": 0
}`,
			wantErr: false,
		},
		{
			name: "prepend to genesis with whitespace",
			input: `{

  "app_state": {}
}`,
			fieldName:  "da_epoch_forced_inclusion",
			fieldValue: "0",
			expectedStart: `{
  "da_epoch_forced_inclusion": 0,

  "app_state": {}`,
			wantErr: false,
		},
		{
			name:       "invalid genesis - no opening brace",
			input:      `"chain_id": "test"}`,
			fieldName:  "da_epoch_forced_inclusion",
			fieldValue: "0",
			wantErr:    true,
		},
		{
			name: "prepend string value",
			input: `{
  "existing": "value"
}`,
			fieldName:  "new_field",
			fieldValue: `"string_value"`,
			expectedStart: `{
  "new_field": "string_value",
  "existing": "value"`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			genesisPath := filepath.Join(tmpDir, "genesis.json")

			err := os.WriteFile(genesisPath, []byte(tt.input), 0o644)
			assert.NilError(t, err)

			err = prependFieldToGenesis(genesisPath, tt.fieldName, tt.fieldValue)

			if tt.wantErr {
				assert.ErrorContains(t, err, "")
				return
			}

			assert.NilError(t, err)

			result, err := os.ReadFile(genesisPath)
			assert.NilError(t, err)

			resultStr := string(result)
			if !strings.HasPrefix(resultStr, tt.expectedStart) {
				t.Errorf("result does not start with expected content\nGot:\n%s\n\nExpected to start with:\n%s", resultStr, tt.expectedStart)
			}

			if !strings.HasSuffix(resultStr, "}") {
				t.Errorf("result does not end with closing brace: %s", resultStr)
			}
		})
	}
}
