package template

import (
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/templates/module"
	modulecreate "github.com/ignite/cli/v29/ignite/templates/module/create"
)

// appConfigModify modifies the application app_config.go to use EVM.
func appConfigModify(appPath string) genny.RunFn {
	return func(r *genny.Runner) error {
		configPath := filepath.Join(appPath, module.PathAppConfigGo)
		f, err := r.Disk.Find(configPath)
		if err != nil {
			return err
		}

		// change imports
		content, err := xast.AppendImports(f.String(),
			xast.WithNamedImport("erc20moduletypes", "github.com/cosmos/evm/x/erc20/types"),
			xast.WithNamedImport("feemarketmoduletypes", "github.com/cosmos/evm/x/feemarket/types"),
			xast.WithNamedImport("evmmoduletypes", "github.com/cosmos/evm/x/vm/types"),
		)
		if err != nil {
			return err
		}

		// module account permissions
		content, err = xast.ModifyGlobalArrayVar(
			content,
			"moduleAccPerms",
			xast.AppendGlobalArrayValue(`{Account: evmmoduletypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}}`),
			xast.AppendGlobalArrayValue(`{Account: erc20moduletypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}}`),
			xast.AppendGlobalArrayValue(`{Account: feemarketmoduletypes.ModuleName}`),
		)
		if err != nil {
			return err
		}

		content, err = modulecreate.AddModuleToAppConfig(content, "erc20", modulecreate.SkipConfigEntry())
		if err != nil {
			return err
		}

		content, err = modulecreate.AddModuleToAppConfig(content, "feemarket", modulecreate.SkipConfigEntry())
		if err != nil {
			return err
		}

		content, err = modulecreate.AddModuleToAppConfig(content,
			"evm",
			modulecreate.SkipConfigEntry(),
			modulecreate.SpecifyModuleEntry("PreBlockers", "InitGenesis", "BeginBlockers", "EndBlockers"),
		)
		if err != nil {
			return err
		}

		// Add new function
		content, err = xast.AppendFunction(content, `// getBlockAccAddrs returns the list of block accounts addresses.
// it appends the addresses of the static precompiles to the blockAccAddrs slice.
func getBlockAccAddrs() []string {
			for _, precompile := range evmmoduletypes.AvailableStaticPrecompiles {
				blockAccAddrs = append(blockAccAddrs, precompile)
			}

			return blockAccAddrs
}`)
		if err != nil {
			return err
		}

		// Use getBlockAccAddrs function
		content = strings.Replace(content, "BlockedModuleAccountsOverride: blockAccAddrs,", "BlockedModuleAccountsOverride: getBlockAccAddrs(),", 1)

		return r.File(genny.NewFileS(configPath, content))
	}
}
