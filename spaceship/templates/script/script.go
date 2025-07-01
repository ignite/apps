package script

import (
	"embed"
	"io/fs"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/xgenny"
	"github.com/ignite/cli/v29/ignite/templates/field/plushhelpers"
)

//go:embed files/*
var fsRunScript embed.FS

// NewRunScripts returns the generator to scaffold a chain and faucet run script.
func NewRunScripts(path, log, home, binDirPath, chainBinPath, faucetBinPath, account, denoms, output string) error {
	runScript, err := fs.Sub(fsRunScript, "files")
	if err != nil {
		return err
	}

	g := genny.New()
	if err := g.OnlyFS(runScript, nil, nil); err != nil {
		return errors.Errorf("generator fs: %w", err)
	}

	ctx := plush.NewContext()
	ctx.Set("path", path)
	ctx.Set("log", log)
	ctx.Set("home", home)
	ctx.Set("chainBinPath", chainBinPath)
	ctx.Set("faucetBinPath", faucetBinPath)
	ctx.Set("binDirPath", binDirPath)
	ctx.Set("binary", filepath.Base(chainBinPath))
	ctx.Set("account", account)
	ctx.Set("denoms", denoms)

	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))

	_, err = xgenny.NewRunner(ctx, output).RunAndApply(g)
	return err
}
