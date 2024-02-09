// Package scaffolder initializes Ignite CLI apps and modifies existing ones
// to add more features in a later time.
package scaffolder

import (
	"context"
	"path/filepath"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosanalysis"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/gocmd"
	"github.com/ignite/cli/v28/ignite/pkg/gomodulepath"
)

const (
	errOldCosmosSDKVersionStr = `Your chain has been scaffolded with an older version of Cosmos SDK: %s

Please, follow the migration guide to upgrade your chain to the latest version at https://docs.ignite.com/migration`
)

// Scaffolder is Ignite CLI app scaffolder.
type Scaffolder struct {
	// Version of the chain
	Version cosmosver.Version

	// path of the app.
	path string

	// modpath represents the go module path of the app.
	modpath gomodulepath.Path
}

// New creates a new scaffold app.
func New(appPath string) (Scaffolder, error) {
	path, err := filepath.Abs(appPath)
	if err != nil {
		return Scaffolder{}, err
	}

	modpath, path, err := gomodulepath.Find(path)
	if err != nil {
		return Scaffolder{}, err
	}

	ver, err := cosmosver.Detect(path)
	if err != nil {
		return Scaffolder{}, err
	}

	// Make sure that the app was scaffolded with a supported Cosmos SDK version
	if err := AssertSupportedCosmosSDKVersion(ver); err != nil {
		return Scaffolder{}, err
	}

	if err := cosmosanalysis.IsChainPath(path); err != nil {
		return Scaffolder{}, err
	}

	s := Scaffolder{
		Version: ver,
		path:    path,
		modpath: modpath,
	}

	return s, nil
}

// AssertSupportedCosmosSDKVersion asserts that a Cosmos SDK version is supported by Ignite CLI.
func AssertSupportedCosmosSDKVersion(v cosmosver.Version) error {
	if v.LT(cosmosver.StargateFiftyVersion) {
		return errors.Errorf(errOldCosmosSDKVersionStr, v)
	}
	return nil
}

func finish(ctx context.Context, path string) error {
	if err := gocmd.Fmt(ctx, path); err != nil {
		return err
	}

	// TODO this function will be available only in the next cli version
	// _ = gocmd.GoImports(ctx, path) // goimports installation could fail, so ignore the error

	return gocmd.ModTidy(ctx, path)
}
