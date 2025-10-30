package template

import (
	"embed"
	"io/fs"
	"math"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	configchain "github.com/ignite/cli/v29/ignite/config/chain"
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

	if err := updateConfigYaml(chain); err != nil {
		return nil, errors.Errorf("failed to update config.yaml: %w", err)
	}

	g.RunFn(commandsModify(appPath, binaryName))
	g.RunFn(rootModify(appPath, binaryName))
	g.RunFn(appModify(appPath))
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

const defaultValPower = 1

// updateConfigYaml updates the default bond tokens.
// this is required as the chain uses 18 decimals.
func updateConfigYaml(c *chain.Chain) error {
	igniteConfig, err := c.Config()
	if err != nil {
		return err
	}

	coins := sdk.NewCoin(igniteConfig.DefaultDenom, sdkmath.NewInt((defaultValPower * int64(math.Pow10(18)))))
	igniteConfig.Validators[0].Bonded = coins.String()
	for i, account := range igniteConfig.Accounts {
		if account.Name == igniteConfig.Validators[0].Name {
			igniteConfig.Accounts[i].Coins = []string{coins.String()}
		}
	}

	if err := configchain.Save(*igniteConfig, c.ConfigPath()); err != nil {
		return err
	}

	return nil
}
