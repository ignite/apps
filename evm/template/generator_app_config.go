package template

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/templates/module"
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
			xast.WithNamedImport("erc20types", "github.com/cosmos/evm/x/erc20/types"),
			xast.WithNamedImport("feemarkettypes", "github.com/cosmos/evm/x/feemarket/types"),
			xast.WithNamedImport("evmtypes", "github.com/cosmos/evm/x/vm/types"),
		)
		if err != nil {
			return err
		}

		// module account permissions
		content, err = xast.ModifyGlobalArrayVar(
			content,
			"moduleAccPerms",
			xast.AppendGlobalArrayValue(`{Account: evmtypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}}`),
			xast.AppendGlobalArrayValue(`{Account: erc20types.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}}`),
			xast.AppendGlobalArrayValue(`{Account: feemarkettypes.ModuleName}`),
		)
		if err != nil {
			return err
		}

		replacer := placeholder.New()

		// begin block / end block configuration
		template := `// cosmos evm modules
		erc20types.ModuleName,
		feemarkettypes.ModuleName,
		evmtypes.ModuleName,
%[1]v`
		replacement := fmt.Sprintf(template, module.PlaceholderSgAppBeginBlockers)
		content = replacer.Replace(content, module.PlaceholderSgAppBeginBlockers, replacement)

		replacement = fmt.Sprintf(template, module.PlaceholderSgAppEndBlockers)
		content = replacer.Replace(content, module.PlaceholderSgAppEndBlockers, replacement)

		// init genesis configuration
		content = strings.Replace(content, "genutiltypes.ModuleName,", "", 1) // delete genutil current position

		template = `// cosmos evm modules
		erc20types.ModuleName,
		feemarkettypes.ModuleName,
		evmtypes.ModuleName,
		// moved down because of evm modules
		genutiltypes.ModuleName,
%[1]v`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppInitGenesis)
		content = replacer.Replace(content, module.PlaceholderSgAppInitGenesis, replacement)

		// Add new function
		content, err = xast.AppendFunction(content, `// getBlockAccAddrs returns the list of block accounts addresses.
// it appends the addresses of the static precompiles to the blockAccAddrs slice.
func getBlockAccAddrs() []string {
			for _, precompile := range evmtypes.AvailableStaticPrecompiles {
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
