package template

import (
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/templates/module"
	modulecreate "github.com/ignite/cli/v29/ignite/templates/module/create"
)

// appConfigModify modifies the application app_config.go to use GnoVM.
func appConfigModify(appPath string) genny.RunFn {
	return func(r *genny.Runner) error {
		configPath := filepath.Join(appPath, module.PathAppConfigGo)
		f, err := r.Disk.Find(configPath)
		if err != nil {
			return err
		}

		// change imports
		content, err := xast.AppendImports(f.String(),
			xast.WithNamedImport("_", "github.com/ignite/gnovm/x/gnovm/module"),
			xast.WithNamedImport("gnovmmoduletypes", "github.com/ignite/gnovm/x/gnovm/types"),
		)
		if err != nil {
			return err
		}

		// module account permissions
		content, err = xast.ModifyGlobalArrayVar(
			content,
			"moduleAccPerms",
			xast.AppendGlobalArrayValue(`{Account: gnovmmoduletypes.ModuleName, Permissions: []string{authtypes.Minter, authtypes.Burner}}`),
		)
		if err != nil {
			return err
		}

		content, err = modulecreate.AddModuleToAppConfig(content, "gnovm")
		if err != nil {
			return err
		}

		return r.File(genny.NewFileS(configPath, content))
	}
}
