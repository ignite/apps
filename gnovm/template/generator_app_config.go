package template

import (
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
	"github.com/ignite/cli/v29/ignite/templates/module"
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

		replacer := placeholder.New()

		// init genesis / begin block / end block configuration
		template := `gnovmmoduletypes.ModuleName,
%[1]v`
		replacement := fmt.Sprintf(template, module.PlaceholderSgAppBeginBlockers)
		content = replacer.Replace(content, module.PlaceholderSgAppBeginBlockers, replacement)

		replacement = fmt.Sprintf(template, module.PlaceholderSgAppEndBlockers)
		content = replacer.Replace(content, module.PlaceholderSgAppEndBlockers, replacement)

		replacement = fmt.Sprintf(template, module.PlaceholderSgAppInitGenesis)
		content = replacer.Replace(content, module.PlaceholderSgAppInitGenesis, replacement)

		// add module config for depinject
		moduleConfigTemplate := `{
			Name:   gnovmmoduletypes.ModuleName,
			Config: appconfig.WrapAny(&gnovmmoduletypes.Module{}),
		},
		%[1]v`
		moduleConfigReplacement := fmt.Sprintf(moduleConfigTemplate, module.PlaceholderSgAppModuleConfig)
		content = replacer.Replace(content, module.PlaceholderSgAppModuleConfig, moduleConfigReplacement)

		return r.File(genny.NewFileS(configPath, content))
	}
}
