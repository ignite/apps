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

		// Import rollkitmngr module
		content, err := xast.AppendImports(f.String(),
			xast.WithNamedImport("rollkitmngrmodule", "github.com/evstack/ev-abci/modules/rollkitmngr/module"),
			xast.WithNamedImport("rollkitmngrtypes", "github.com/evstack/ev-abci/modules/rollkitmngr/types"),
			xast.WithNamedImport("_", "github.com/evstack/ev-abci/modules/rollkitmngr"),
		)
		if err != nil {
			return err
		}

		// end block for rollkitmngr
		template := `rollkitmngrtypes.ModuleName,
%[1]v`
		replacement := fmt.Sprintf(template, module.PlaceholderSgAppEndBlockers)
		content = replacer.Replace(content, module.PlaceholderSgAppEndBlockers, replacement)

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
