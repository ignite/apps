package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateCosmosApp(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) string
		wantErr bool
	}{
		{
			name: "valid app",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				if err := os.MkdirAll(filepath.Join(tmpDir, "app"), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644); err != nil {
					t.Fatal(err)
				}
				return tmpDir
			},
			wantErr: false,
		},
		{
			name: "no go.mod",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				if err := os.MkdirAll(filepath.Join(tmpDir, "app"), 0755); err != nil {
					t.Fatal(err)
				}
				return tmpDir
			},
			wantErr: true,
		},
		{
			name: "no app directory",
			setup: func(t *testing.T) string {
				tmpDir := t.TempDir()
				if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644); err != nil {
					t.Fatal(err)
				}
				return tmpDir
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appPath := tt.setup(t)
			err := validateCosmosApp(appPath)

			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
			} else if !tt.wantErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}

func TestGetModuleInfo(t *testing.T) {
	tests := []struct {
		name           string
		goModContent   string
		expectedModule string
		expectedApp    string
		wantErr        bool
	}{
		{
			name: "simple module",
			goModContent: `module github.com/example/test

go 1.21`,
			expectedModule: "github.com/example/test",
			expectedApp:    "test",
			wantErr:        false,
		},
		{
			name: "complex module path",
			goModContent: `module github.com/cosmos/gaia/v15

go 1.21

require (
	github.com/cosmos/cosmos-sdk v0.50.0
)`,
			expectedModule: "github.com/cosmos/gaia/v15",
			expectedApp:    "v15",
			wantErr:        false,
		},
		{
			name: "no module declaration",
			goModContent: `go 1.21

require (
	github.com/cosmos/cosmos-sdk v0.50.0
)`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			goModPath := filepath.Join(tmpDir, "go.mod")
			if err := os.WriteFile(goModPath, []byte(tt.goModContent), 0644); err != nil {
				t.Fatal(err)
			}

			modulePath, appName, err := getModuleInfo(tmpDir)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("expected no error but got: %v", err)
				return
			}

			if modulePath != tt.expectedModule {
				t.Errorf("expected module path %q but got %q", tt.expectedModule, modulePath)
			}

			if appName != tt.expectedApp {
				t.Errorf("expected app name %q but got %q", tt.expectedApp, appName)
			}
		})
	}
}

func TestUpdateGoMod(t *testing.T) {
	tests := []struct {
		name           string
		initialContent string
		wantContains   []string
	}{
		{
			name: "basic go.mod",
			initialContent: `module github.com/example/test

go 1.21

require (
	github.com/cosmos/cosmos-sdk v0.50.0
)`,
			wantContains: []string{
				"replace github.com/ethereum/go-ethereum => github.com/cosmos/go-ethereum",
				"github.com/cosmos/evm",
				"github.com/ethereum/go-ethereum",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			goModPath := filepath.Join(tmpDir, "go.mod")

			if err := os.WriteFile(goModPath, []byte(tt.initialContent), 0644); err != nil {
				t.Fatal(err)
			}

			if err := updateGoMod(tmpDir); err != nil {
				t.Errorf("updateGoMod failed: %v", err)
				return
			}

			content, err := os.ReadFile(goModPath)
			if err != nil {
				t.Fatal(err)
			}

			contentStr := string(content)
			for _, want := range tt.wantContains {
				if !strings.Contains(contentStr, want) {
					t.Errorf("expected updated go.mod to contain %q, but it doesn't. Content:\n%s", want, contentStr)
				}
			}
		})
	}
}

func TestGenerateAnteFile(t *testing.T) {
	modulePath := "github.com/example/test"
	appName := "test"

	content := generateAnteFile(modulePath, appName)

	expectedContains := []string{
		"package app",
		"import (",
		"func (app *App) setAnteHandler",
		"appante \"github.com/example/test/app/ante\"",
		"app.SetAnteHandler(appante.NewAnteHandler(options))",
	}

	for _, expected := range expectedContains {
		if !strings.Contains(content, expected) {
			t.Errorf("expected generated ante file to contain %q, but it doesn't", expected)
		}
	}
}

func TestGenerateAnteHandlerFile(t *testing.T) {
	content := generateAnteHandlerFile()

	expectedContains := []string{
		"package ante",
		"func NewAnteHandler",
		"newMonoEVMAnteHandler",
		"newCosmosAnteHandler",
		"ExtensionOptionsEthereumTx",
		"ExtensionOptionDynamicFeeTx",
	}

	for _, expected := range expectedContains {
		if !strings.Contains(content, expected) {
			t.Errorf("expected generated ante handler file to contain %q, but it doesn't", expected)
		}
	}
}

func TestGenerateEVMFile(t *testing.T) {
	appName := "test"

	content := generateEVMFile(appName)

	expectedContains := []string{
		"package app",
		"func (app *App) registerEVMModules",
		"func (app *App) postRegisterEVMModules",
		"func getCustomEVMActivators",
		"func getEVMChainID",
		"func cosmosChainIDToEVMChainID",
		"func RegisterEVM",
		"func ProvideMsgEthereumTxCustomGetSigner",
		"func (app *App) GetStoreKeysMap",
	}

	for _, expected := range expectedContains {
		if !strings.Contains(content, expected) {
			t.Errorf("expected generated EVM file to contain %q, but it doesn't", expected)
		}
	}
}
