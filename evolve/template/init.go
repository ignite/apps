package template

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/templates/module"
)

// commandsStartModify modifies the application start to use evolve.
func commandsStartModify(appPath, binaryName string, version cosmosver.Version) genny.RunFn {
	return func(r *genny.Runner) error {
		cmdPath := filepath.Join(appPath, "cmd", binaryName, "cmd", "commands.go")
		f, err := r.Disk.Find(cmdPath)
		if err != nil {
			return err
		}

		if strings.Contains(f.String(), RollkitV0XStartHandler) {
			return errors.New("rollkit v0.x is already installed. Please remove it before installing evolve v1.x")
		}

		if strings.Contains(f.String(), EvolveV1XStartHandler) {
			return errors.New("ev-abci is already installed.")
		}

		if version.LT(cosmosver.StargateFiftyVersion) {
			return errors.New("Evolve requires Ignite v28+ / Cosmos SDK v0.50+")
		}

		content, err := xast.AppendImports(
			f.String(),
			xast.WithNamedImport("abciserver", "github.com/evstack/ev-abci/server"),
		)
		if err != nil {
			return err
		}

		// replace potential legacy boilerplate present in an ignite v28 chain.
		content = replaceLegacyAddCommands(content)

		// modify the add commands arguments using xast.
		content, err = xast.ModifyCaller(content, ServerAddCommandsWithStartCmdOptions, func(args []string) ([]string, error) {
			return []string{
				"rootCmd",
				"app.DefaultNodeHome",
				"newApp",
				"appExport",
				`server.StartCmdOptions{
				AddFlags: addModuleInitFlags,
				StartCommandHandler: abciserver.StartHandler(),
			}`,
			}, nil
		})

		// add the start command flags
		content, err = xast.ModifyFunction(content,
			"addModuleInitFlags",
			xast.AppendFuncCode("abciserver.AddFlags(startCmd)"),
		)
		if err != nil {
			return err
		}

		return r.File(genny.NewFileS(cmdPath, content))
	}
}

// commandsGenesisInitModify modifies the application genesis init command to use evolve.
// this is only needed when the start command is also modified.
func commandsGenesisInitModify(appPath, binaryName string) genny.RunFn {
	return func(r *genny.Runner) error {
		cmdPath := filepath.Join(appPath, "cmd", binaryName, "cmd", "commands.go")
		f, err := r.Disk.Find(cmdPath)
		if err != nil {
			return err
		}

		content, err := xast.AppendImports(
			f.String(),
			xast.WithNamedImport("evnodeconf", "github.com/evstack/ev-node/pkg/config"),
			xast.WithNamedImport("abciserver", "github.com/evstack/ev-abci/server"),
		)
		if err != nil {
			return err
		}

		// use ast to modify the function that initializes genesisCmd
		content, err = xast.ModifyFunction(content, "initRootCmd",
			xast.AppendFuncAtLine(`
		genesisCmd := genutilcli.InitCmd(basicManager, app.DefaultNodeHome)
		evnodeconf.AddFlags(genesisCmd)
		genesisCmdRunE := genesisCmd.RunE
		genesisCmd.RunE = func(cmd *cobra.Command, args []string) error {
		    if err := genesisCmdRunE(cmd, args); err != nil {
		        return err
		    }
		    return abciserver.InitRunE(cmd, args)
		}
		        `,
				0),
		)
		if err != nil {
			return err
		}

		// modify the add commands arguments using xast.
		content, err = xast.ModifyCaller(content, "rootCmd.AddCommand", func(args []string) ([]string, error) {
			if strings.Contains(args[0], "InitCmd") {
				args[0] = "genesisCmd"
			}

			return args, nil
		})

		return r.File(genny.NewFileS(cmdPath, content))
	}
}

// commandsRollbackModify modifies the application rollback command to use evolve.
func commandsRollbackModify(appPath, binaryName string) genny.RunFn {
	return func(r *genny.Runner) error {
		cmdPath := filepath.Join(appPath, "cmd", binaryName, "cmd", "commands.go")
		f, err := r.Disk.Find(cmdPath)
		if err != nil {
			return err
		}

		// use ast to modify the function that initializes genesisCmd
		content, err := xast.ModifyFunction(f.String(), "initRootCmd",
			xast.AppendFuncCode(`
				// override rollback command
				evNodeRollbackCmd := abciserver.NewRollbackCmd(newApp, app.DefaultNodeHome)
				if currentRollbackCmd, _, err := rootCmd.Find([]string{evNodeRollbackCmd.Name()}); err == nil{
					rootCmd.RemoveCommand(currentRollbackCmd)
				}
				rootCmd.AddCommand(evNodeRollbackCmd)
		    `,
			),
		)
		if err != nil {
			return err
		}

		return r.File(genny.NewFileS(cmdPath, content))
	}
}

// appModify modifies the app to add the blanked staking module and optionally the migration utilities.
func appModify(appPath string, withMigration bool) genny.RunFn {
	replacer := placeholder.New()

	appConfigModify := func(r *genny.Runner, withMigration bool) error {
		configPath := filepath.Join(appPath, module.PathAppConfigGo)
		f, err := r.Disk.Find(configPath)
		if err != nil {
			return err
		}

		content := f.String()

		if withMigration {
			// Import migrationmngr module
			content, err = xast.AppendImports(content,
				xast.WithNamedImport("migrationmngrmodule", "github.com/evstack/ev-abci/modules/migrationmngr/module"),
				xast.WithNamedImport("migrationmngrtypes", "github.com/evstack/ev-abci/modules/migrationmngr/types"),
				xast.WithNamedImport("_", "github.com/evstack/ev-abci/modules/migrationmngr"),
			)
			if err != nil {
				return err
			}

			// add migrationmngr module config for depinject
			moduleConfigTemplate := `{
				Name:   migrationmngrtypes.ModuleName,
				Config: appconfig.WrapAny(&migrationmngrmodule.Module{}),
			},
			%[1]v`
			moduleConfigReplacement := fmt.Sprintf(moduleConfigTemplate, module.PlaceholderSgAppModuleConfig)
			content = replacer.Replace(content, module.PlaceholderSgAppModuleConfig, moduleConfigReplacement)

			// preblocker for migrationmngr
			preBlockerTemplate := `migrationmngrtypes.ModuleName,
						%[1]v`
			preBlockerReplacement := fmt.Sprintf(preBlockerTemplate, "// this line is used by starport scaffolding # stargate/app/preBlockers")
			content = replacer.Replace(content, "// this line is used by starport scaffolding # stargate/app/preBlockers", preBlockerReplacement)

			// end block for migrationmngr
			endBlockerTemplate := `migrationmngrtypes.ModuleName,
%[1]v`
			endBlockerReplacement := fmt.Sprintf(endBlockerTemplate, module.PlaceholderSgAppEndBlockers)
			content = replacer.Replace(content, module.PlaceholderSgAppEndBlockers, endBlockerReplacement)
		}

		// replace staking blank import
		content, err = xast.RemoveImports(content,
			xast.WithNamedImport("_", "github.com/cosmos/cosmos-sdk/x/staking"),
		)
		if err != nil {
			return err
		}

		if content, err = xast.AppendImports(content,
			xast.WithNamedImport("_", "github.com/evstack/ev-abci/modules/staking"),
		); err != nil {
			return err
		}

		return r.File(genny.NewFileS(configPath, content))
	}

	appGoModify := func(r *genny.Runner) error {
		configPath := filepath.Join(appPath, module.PathAppGo)
		f, err := r.Disk.Find(configPath)
		if err != nil {
			return err
		}

		// replace staking import
		content, err := xast.RemoveImports(f.String(),
			xast.WithNamedImport("stakingkeeper", "github.com/cosmos/cosmos-sdk/x/staking/keeper"),
		)
		if err != nil {
			return err
		}

		if content, err = xast.AppendImports(content,
			xast.WithNamedImport("stakingkeeper", "github.com/evstack/ev-abci/modules/staking/keeper"),
		); err != nil {
			return err
		}

		return r.File(genny.NewFileS(configPath, content))
	}

	exportModify := func(r *genny.Runner) error {
		configPath := filepath.Join(appPath, filepath.Join(module.PathAppModule, "export.go"))
		f, err := r.Disk.Find(configPath)
		if err != nil {
			return err
		}

		content := strings.ReplaceAll(f.String(), "staking.WriteValidators(ctx, app.StakingKeeper)", "staking.WriteValidators(ctx, app.StakingKeeper.Keeper)")

		return r.File(genny.NewFileS(configPath, content))
	}

	return func(r *genny.Runner) error {
		err := appConfigModify(r, withMigration)
		err = errors.Join(err, exportModify(r))
		err = errors.Join(err, appGoModify(r))

		return err
	}
}

// replaceLegacyAddCommands replaces the legacy `AddCommands` with a temporary `AddCommandsWithStartCmdOptions` boilerplate.
// Atfterwards, we let the same xast function replace the `AddCommandsWithStartCmdOptions` argument.
func replaceLegacyAddCommands(content string) string {
	return strings.Replace(content, "server.AddCommands(", ServerAddCommandsWithStartCmdOptions+"(", 1)
}
