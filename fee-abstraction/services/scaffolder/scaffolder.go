package scaffolder

import (
	"context"
	"fmt"

	"github.com/blang/semver/v4"
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
	errNewCosmosSDKVersionStr = "Your chain has been scaffolded with the new version (%s) of Cosmos SDK greater than %s"
)

// Scaffolder is fee abstraction app scaffolder.
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

// assertSupportedCosmosSDKVersion asserts that a Cosmos SDK version is supported by Ignite CLI.
func assertSupportedCosmosSDKVersion(v cosmosver.Version) error {
	v0501 := cosmosver.StargateFiftyVersion
	v0510, err := cosmosver.Parse("0.51.0")
	if err != nil {
		return err
	}
	switch {
	case v.LT(v0501):
		return errors.Errorf(errOldCosmosSDKVersionStr, v.String())
	case v.GTE(v0510):
		return errors.Errorf(errNewCosmosSDKVersionStr, v.String(), v0510.String())
	}
	return nil
}

// finish finalize the scaffolded code downloading the fee abstraction and formatting the code.
func finish(ctx context.Context, session *cliui.Session, path string, version semver.Version) error {
	// Add fee-abstraction module to the go.mod
	session.StartSpinner("Downloading fee abstraction module...")

	pkgVersion, ok := versionMap[version.Major]
	if !ok {
		return errors.Errorf("version %d not supported", version.Major)
	}
	URL := fmt.Sprintf("%s@v%s", pkgVersion, version.String())
	if err := gocmd.Get(ctx, path, []string{URL}); err != nil {
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
