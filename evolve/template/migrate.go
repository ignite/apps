package template

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/templates/module"
)

// migrateFromCometModify modifies the app to add the migration from cometbft commands and modules.
func migrateFromCometModify(appPath string) genny.RunFn {
	replacer := placeholder.New()

	appConfigModify := func(r *genny.Runner) error {
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

		// end block for migrationmngr
		endBlockerTemplate := `migrationmngrtypes.ModuleName,
%[1]v`
		endBlockerReplacement := fmt.Sprintf(endBlockerTemplate, module.PlaceholderSgAppEndBlockers)
		content = replacer.Replace(content, module.PlaceholderSgAppEndBlockers, endBlockerReplacement)

		// replace staking blank import
		content = strings.Replace(content, "github.com/cosmos/cosmos-sdk/x/staking", "github.com/evstack/ev-abci/modules/staking", 1)

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
