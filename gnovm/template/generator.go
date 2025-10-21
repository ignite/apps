package template

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodule"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
)

//go:embed files/* files/**/*
var fsAppGnoVM embed.FS

// NewGnoVMGenerator returns the generator to scaffold a gnoVM integration inside an app.
func NewGnoVMGenerator(chain *chain.Chain) (*genny.Generator, error) {
	appGnoVM, err := fs.Sub(fsAppGnoVM, "files")
	if err != nil {
		return nil, errors.Errorf("fail to generate sub: %w", err)
	}

	g := genny.New()

	if err := g.OnlyFS(appGnoVM, nil, nil); err != nil {
		return g, err
	}

	appPath := chain.AppPath()
	modPath, _, err := gomodulepath.Find(appPath)
	if err != nil {
		return nil, err
	}

	ctx := plush.NewContext()
	ctx.Set("ModulePath", modPath.RawPath)
	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))

	if err := updateDependencies(appPath); err != nil {
		return nil, errors.Errorf("failed to update go.mod: %w", err)
	}

	g.RunFn(appModify(appPath))
	g.RunFn(appConfigModify(appPath))

	return g, nil
}

// updateDependencies makes sure the correct dependencies are added to the go.mod files.
func updateDependencies(appPath string) error {
	gomod, err := gomodule.ParseAt(appPath)
	if err != nil {
		return errors.Errorf("failed to parse go.mod: %w", err)
	}

	// add required dependencies
	if err := gomod.AddRequire(GnoVMModulePackage, GnoVMModuleVersion); err != nil {
		return errors.Errorf("failed to add gnovm module dependencie: %w", err)
	}

	// add temporary replaces
	if err := gomod.AddReplace(GnolangPackage, "", GnolangForkPackage, GnolangForkVersion); err != nil {
		return errors.Errorf("failed to add gnolang dependencie replace: %w", err)
	}

	// save go.mod
	data, err := gomod.Format()
	if err != nil {
		return errors.Errorf("failed to format go.mod: %w", err)
	}

	return os.WriteFile(filepath.Join(appPath, "go.mod"), data, 0o644)
}
