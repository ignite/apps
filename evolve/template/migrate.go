package template

import (
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/ignite/cli/v29/ignite/pkg/xast"
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
