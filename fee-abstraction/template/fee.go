package template

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/xast"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/templates/field/plushhelpers"
	"github.com/ignite/cli/v28/ignite/templates/module"
)

const funcRegisterIBCFeeAbs = `
	modules := map[string]appmodule.AppModule{
		ibcexported.ModuleName:      ibc.AppModule{},
		ibctransfertypes.ModuleName: ibctransfer.AppModule{},
		ibcfeetypes.ModuleName:      ibcfee.AppModule{},
		icatypes.ModuleName:         icamodule.AppModule{},
		capabilitytypes.ModuleName:  capability.AppModule{},
		ibctm.ModuleName:            ibctm.AppModule{},
		solomachine.ModuleName:      solomachine.AppModule{},
		feeabstypes.ModuleName:      feeabsmodule.AppModule{},
	}

    for name, m := range modules {
		module.CoreAppModuleBasicAdaptor(name, m).RegisterInterfaces(registry)
	}

	return modules`

// Options fee abstraction scaffold options.
type Options struct {
	BinaryName string
	AppPath    string
	Home       string
}

//go:embed files/* files/**/*
var fsAppFeeAbs embed.FS

// NewFeeAbstractionGenerator returns the generator to scaffold a fee abstraction integration inside an app.
func NewFeeAbstractionGenerator(replacer placeholder.Replacer, opts *Options) (*genny.Generator, error) {
	var (
		g         = genny.New()
		appFeeAbs = xgenny.NewEmbedWalker(
			fsAppFeeAbs,
			"files/",
			opts.AppPath,
		)
	)
	if err := g.Box(appFeeAbs); err != nil {
		return g, err
	}

	ctx := plush.NewContext()
	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))

	g.RunFn(appModify(replacer, opts))
	g.RunFn(appConfigModify(replacer, opts))
	g.RunFn(ibcModify(replacer, opts))

	return g, nil
}

// appModify app.go modification when adding fee abstraction integration.
func appModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		appPath := filepath.Join(opts.AppPath, module.PathAppGo)
		f, err := r.Disk.Find(appPath)
		if err != nil {
			return err
		}

		// Import
		content, err := xast.AppendImports(f.String(),
			xast.WithLastNamedImport("feeabskeeper", "github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/keeper"),
		)
		if err != nil {
			return err
		}

		// Keeper declaration
		template := `
// Fee Abstraction
FeeAbsKeeper		feeabskeeper.Keeper
ScopedFeeAbsKeeper	capabilitykeeper.ScopedKeeper

%[1]v`
		replacement := fmt.Sprintf(template, module.PlaceholderSgAppKeeperDeclaration)
		content = replacer.Replace(content, module.PlaceholderSgAppKeeperDeclaration, replacement)

		return r.File(genny.NewFileS(appPath, content))
	}
}

// appConfigModify app_config.go modification when adding fee abstraction integration.
func appConfigModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		configPath := filepath.Join(opts.AppPath, module.PathAppConfigGo)
		f, err := r.Disk.Find(configPath)
		if err != nil {
			return err
		}

		// Import
		content, err := xast.AppendImports(f.String(),
			xast.WithLastNamedImport("feeabstypes", "github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/types"),
		)
		if err != nil {
			return err
		}

		// Init genesis / begin block / end block
		template := `feeabstypes.ModuleName,
%[1]v`
		replacement := fmt.Sprintf(template, module.PlaceholderSgAppInitGenesis)
		content = replacer.Replace(content, module.PlaceholderSgAppInitGenesis, replacement)
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppBeginBlockers)
		content = replacer.Replace(content, module.PlaceholderSgAppBeginBlockers, replacement)
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppEndBlockers)
		content = replacer.Replace(content, module.PlaceholderSgAppEndBlockers, replacement)

		// Mac Perms
		template = `{Account: feeabstypes.ModuleName},
%[1]v`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppMaccPerms)
		content = replacer.Replace(content, module.PlaceholderSgAppMaccPerms, replacement)

		return r.File(genny.NewFileS(configPath, content))
	}
}

// ibcModify ibc.go modification when adding fee abstraction integration.
func ibcModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		ibcPath := filepath.Join(opts.AppPath, "app/ibc.go")
		f, err := r.Disk.Find(ibcPath)
		if err != nil {
			return err
		}

		// Import
		content, err := xast.AppendImports(f.String(),
			xast.WithLastNamedImport("feeabsmodule", "github.com/osmosis-labs/fee-abstraction/v8/x/feeabs"),
			xast.WithLastNamedImport("feeabstypes", "github.com/osmosis-labs/fee-abstraction/v8/x/feeabs/types"),
		)
		if err != nil {
			return err
		}

		// create fee abstraction module
		templateIBCModule := `feeAbsStack, err := app.registerFeeAbstractionModules()
if err != nil {
	return err
}
ibcRouter.AddRoute(feeabstypes.ModuleName, feeAbsStack)

%[1]v`
		replacementIBCModule := fmt.Sprintf(templateIBCModule, module.PlaceholderIBCNewModule)
		content = replacer.Replace(content, module.PlaceholderIBCNewModule, replacementIBCModule)

		content, err = xast.ModifyFunction(content,
			"RegisterIBC",
			xast.ReplaceFuncBody(
				funcRegisterIBCFeeAbs,
			),
		)
		if err != nil {
			return err
		}

		return r.File(genny.NewFileS(ibcPath, content))
	}
}
