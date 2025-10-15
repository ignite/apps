package template

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodule"
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

// appConfigStakingModify modifies the app to add the blanked x/staking modules.
func appConfigStakingModify(appPath string) genny.RunFn {
	appConfigModify := func(r *genny.Runner) error {
		configPath := filepath.Join(appPath, module.PathAppConfigGo)
		f, err := r.Disk.Find(configPath)
		if err != nil {
			return err
		}

		// replace staking blank import
		content := strings.Replace(f.String(), "github.com/cosmos/cosmos-sdk/x/staking", "github.com/evstack/ev-abci/modules/staking", 1)

		return r.File(genny.NewFileS(configPath, content))
	}

	appGoModify := func(r *genny.Runner) error {
		configPath := filepath.Join(appPath, module.PathAppGo)
		f, err := r.Disk.Find(configPath)
		if err != nil {
			return err
		}

		content := strings.ReplaceAll(f.String(), "github.com/cosmos/cosmos-sdk/x/staking/keeper", "github.com/evstack/ev-abci/modules/staking/keeper")

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
		err := appConfigModify(r)
		err = errors.Join(err, exportModify(r))
		err = errors.Join(err, appGoModify(r))

		return err
	}
}

// updateDependencies makes sure the correct dependencies are added to the go.mod files.
// ev-abci expects evolve v1 to be used.
func updateDependencies(appPath string) error {
	gomod, err := gomodule.ParseAt(appPath)
	if err != nil {
		return errors.Errorf("failed to parse go.mod: %w", err)
	}

	gomod.AddNewRequire(EvABCIPackage, EvABCIVersion, false)
	gomod.AddNewRequire(EvNodePackage, EvNodeVersion, false)

	// add local-da as go tool dependency (useful for local development)
	if err := gomod.AddTool(EvNodeDaCmd); err != nil {
		return errors.Errorf("failed to add local-da tool: %w", err)
	}

	// add required replaces
	gomod.AddReplace(GoHeaderPackage, "", GoHeaderPackageFork, GoHeaderVersionFork)

	// save go.mod
	data, err := gomod.Format()
	if err != nil {
		return errors.Errorf("failed to format go.mod: %w", err)
	}

	return os.WriteFile(filepath.Join(appPath, "go.mod"), data, 0o644)
}

// replaceLegacyAddCommands replaces the legacy `AddCommands` with a temporary `AddCommandsWithStartCmdOptions` boilerplate.
// Atfterwards, we let the same xast function replace the `AddCommandsWithStartCmdOptions` argument.
func replaceLegacyAddCommands(content string) string {
	return strings.Replace(content, "server.AddCommands(", ServerAddCommandsWithStartCmdOptions+"(", 1)
}
