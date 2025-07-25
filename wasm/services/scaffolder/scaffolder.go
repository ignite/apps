package scaffolder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/blang/semver/v4"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosanalysis"
	"github.com/ignite/cli/v29/ignite/pkg/gocmd"
	"github.com/ignite/cli/v29/ignite/services/chain"
)

// Scaffolder is Wasm app scaffolder.
type Scaffolder struct {
	chain   *chain.Chain
	session *cliui.Session
}

// New creates a new scaffold app.
func New(c *chain.Chain, session *cliui.Session) (Scaffolder, error) {
	if err := cosmosanalysis.IsChainPath(c.AppPath()); err != nil {
		return Scaffolder{}, err
	}
	return Scaffolder{chain: c, session: session}, nil
}

// hasWasm check if the app already have the wasm integration verifying if the app/wasm.go file exist.
func hasWasm(appPath string) bool {
	if _, err := os.Stat(filepath.Join(appPath, "app/wasm.go")); err == nil {
		return true
	}

	return false
}

// finish finalize the scaffolded code downloading the wasm and formatting the code.
func finish(ctx context.Context, session *cliui.Session, path string, wasmVersion semver.Version) error {
	// Add wasmd to the go.mod
	session.StartSpinner("Downloading wasmd module...")

	wasmURL := fmt.Sprintf("%s@v%s", wasmRepo, wasmVersion.String())
	if err := gocmd.Get(ctx, path, []string{wasmURL}); err != nil {
		return err
	}

	if err := gocmd.ModTidy(ctx, path); err != nil {
		return err
	}

	session.StartSpinner("Formatting code...")
	if err := gocmd.Fmt(ctx, path); err != nil {
		return err
	}

	_ = gocmd.GoImports(ctx, path) // goimports installation could fail, so ignore the error

	return nil
}
