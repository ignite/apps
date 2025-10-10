package template

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodule"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
)

//go:embed files/* files/**/*
var fsAppEvm embed.FS

// NewEVMGenerator returns the generator to scaffold a evm integration inside an app.
func NewEVMGenerator(chain *chain.Chain) (*genny.Generator, error) {
	appEvm, err := fs.Sub(fsAppEvm, "files")
	if err != nil {
		return nil, errors.Errorf("fail to generate sub: %w", err)
	}

	g := genny.New()

	if err := g.OnlyFS(appEvm, nil, nil); err != nil {
		return g, err
	}

	appPath := chain.AppPath()
	modPath, _, err := gomodulepath.Find(appPath)
	if err != nil {
		return nil, err
	}

	ctx := plush.NewContext()
	ctx.Set("ModulePath", modPath.RawPath)
	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))

	binaryName, err := chain.Binary()
	if err != nil {
		return nil, err
	}

	if err := updateDependencies(appPath); err != nil {
		return nil, errors.Errorf("failed to update go.mod: %w", err)
	}

	g.RunFn(commandsModify(appPath, binaryName))
	g.RunFn(rootModify(appPath, binaryName))
	g.RunFn(appModify(appPath, binaryName))
	g.RunFn(appConfigModify(appPath))
	g.RunFn(ibcModify(appPath))

	return g, nil
}

// updateDependencies makes sure the correct dependencies are added to the go.mod files.
func updateDependencies(appPath string) error {
	gomod, err := gomodule.ParseAt(appPath)
	if err != nil {
		return errors.Errorf("failed to parse go.mod: %w", err)
	}

	gomod.AddNewRequire(CosmosEVMPackage, CosmosEVMVersion, false)

	// add required replaces
	gomod.AddReplace(EthereumGoEthereumPackage, "", CosmosGoEthereumPackage, CosmosGoEthereumVersion)

	// save go.mod
	data, err := gomod.Format()
	if err != nil {
		return errors.Errorf("failed to format go.mod: %w", err)
	}

	return os.WriteFile(filepath.Join(appPath, "go.mod"), data, 0o644)
}
