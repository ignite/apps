package scaffolder

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosanalysis"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/gocmd"
	"github.com/ignite/cli/v28/ignite/services/chain"
)

const (
	errOldCosmosSDKVersionStr = `Your chain has been scaffolded with an older version of Cosmos SDK: %s

Please, follow the migration guide to upgrade your chain to the latest version at https://docs.ignite.com/migration`
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
	if err := assertSupportedCosmosSDKVersion(c.Version); err != nil {
		return Scaffolder{}, err
	}
	return Scaffolder{chain: c, session: session}, nil
}

func hasWasm(appPath string) bool {
	if _, err := os.Stat(filepath.Join(appPath, "app/wasm.go")); err == nil {
		return true
	}
	return false
}

// assertSupportedCosmosSDKVersion asserts that a Cosmos SDK version is supported by Ignite CLI.
func assertSupportedCosmosSDKVersion(v cosmosver.Version) error {
	if v.LT(cosmosver.StargateFiftyVersion) {
		return errors.Errorf(errOldCosmosSDKVersionStr, v)
	}
	return nil
}

func finish(ctx context.Context, session *cliui.Session, path string) error {
	// Add wasmd to the go.mod
	session.StartSpinner("Downloading wasmd module...")
	if err := gocmd.Get(ctx, path, []string{wasmRepo}); err != nil {
		return err
	}

	if err := gocmd.ModTidy(ctx, path); err != nil {
		return err
	}

	session.StartSpinner("Formatting code...")
	if err := gocmd.Fmt(ctx, path); err != nil {
		return err
	}

	// TODO this function will be available only in the next cli version
	// _ = gocmd.GoImports(ctx, path) // goimports installation could fail, so ignore the error

	return nil
}
