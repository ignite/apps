// Package scaffolder initializes Ignite CLI apps and modifies existing ones
// to add more features in a later time.
package scaffolder

import (
	"context"
	"path/filepath"

	"github.com/ignite/cli/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/ignite/pkg/cosmosanalysis"
	"github.com/ignite/cli/ignite/pkg/cosmosver"
	"github.com/ignite/cli/ignite/pkg/gocmd"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/version"
)

const (
	wasmImport = "github.com/CosmWasm/wasmd"
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
	if err := version.AssertSupportedCosmosSDKVersion(ver); err != nil {
		return Scaffolder{}, err
	}

	if err := cosmosanalysis.IsChainPath(path); err != nil {
		return Scaffolder{}, err
	}

	return Scaffolder{
		Version: ver,
		path:    path,
		modpath: modpath,
	}, nil
}

func finish(ctx context.Context, path, gomodPath string) error {
	if err := gocmd.ModTidy(ctx, path); err != nil {
		return err
	}
	return gocmd.Fmt(ctx, path)
}

func (s Scaffolder) installWasm(ctx context.Context, version string) error {
	return cmdrunner.
		New().
		Run(
			ctx,
			step.New(step.Exec(gocmd.Name(), "get", gocmd.PackageLiteral(wasmImport, version))),
		)
}
