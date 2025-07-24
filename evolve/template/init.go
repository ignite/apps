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
)

// commandsStartModify modifies the application start to use evolve.
func commandsStartModify(appPath, binaryName string, version cosmosver.Version) genny.RunFn {
	return func(r *genny.Runner) error {
		cmdPath := filepath.Join(appPath, "cmd", binaryName, "cmd/commands.go")
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

// commandsGenesisModify modifies the application genesis command to use evolve.
func commandsGenesisModify(appPath, binaryName string) genny.RunFn {
	return func(r *genny.Runner) error {
		cmdPath := filepath.Join(appPath, "cmd", binaryName, "cmd/commands.go")
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
		alreadyAdded := false // to avoid adding the migrate command multiple times as there are multiple calls to `rootCmd.AddCommand`
		content, err = xast.ModifyCaller(content, "rootCmd.AddCommand", func(args []string) ([]string, error) {
			if strings.Contains(args[0], "InitCmd") {
				args[0] = "genesisCmd"
			}

			// add migrate command
			if !alreadyAdded {
				args = append(args, evolveV1MigrateCmd)
				alreadyAdded = true
			}

			return args, nil
		})

		return r.File(genny.NewFileS(cmdPath, content))
	}
}

// updateDependencies makes sure the correct dependencies are added to the go.mod files.
// ev-abci expects evolve v1 to be used.
func updateDependencies(appPath string) error {
	gomod, err := gomodule.ParseAt(appPath)
	if err != nil {
		return errors.Errorf("failed to parse go.mod: %w", err)
	}

	gomod.AddNewRequire(GoExecPackage, GoExecVersion, false)
	gomod.AddNewRequire(EvNodePackage, EvNodeVersion, false)

	// add local-da as go tool dependency (useful for local development)
	if err := gomod.AddTool(EvNodeDaCmd); err != nil {
		return errors.Errorf("failed to add local-da tool: %w", err)
	}

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
