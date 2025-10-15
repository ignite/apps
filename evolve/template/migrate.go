package template

import (
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/templates/module"
)

// commandsMigrateModify adds the evolve migrate command to the application.
func commandsMigrateModify(appPath, binaryName string) genny.RunFn {
	return func(r *genny.Runner) error {
		cmdPath := filepath.Join(appPath, "cmd", binaryName, "cmd", "commands.go")
		f, err := r.Disk.Find(cmdPath)
		if err != nil {
			return err
		}

		content, err := xast.AppendImports(
			f.String(),
			xast.WithNamedImport("abciserver", "github.com/evstack/ev-abci/server"),
		)
		if err != nil {
			return err
		}

		// add migrate command
		alreadyAdded := false // to avoid adding the migrate command multiple times as there are multiple calls to `rootCmd.AddCommand`
		content, err = xast.ModifyCaller(content, "rootCmd.AddCommand", func(args []string) ([]string, error) {
			if !alreadyAdded {
				args = append(args, evolveV1MigrateCmd)
				alreadyAdded = true
			}

			return args, nil
		})

		return r.File(genny.NewFileS(cmdPath, content))
	}
}

// appConfigMigrateModify modifies the app to add the migration from cometbft commands and modules.
func appConfigMigrateModify(appPath string) genny.RunFn {
	replacer := placeholder.New()

	return func(r *genny.Runner) error {
		configPath := filepath.Join(appPath, module.PathAppConfigGo)
		f, err := r.Disk.Find(configPath)
		if err != nil {
			return err
		}

		// Import migrationmngr module
		content, err := xast.AppendImports(f.String(),
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

		return r.File(genny.NewFileS(configPath, content))
	}
}
