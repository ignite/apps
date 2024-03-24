package wasm

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/plush/v4"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/pkg/placeholder"
	"github.com/ignite/cli/v28/ignite/pkg/xgenny"
	"github.com/ignite/cli/v28/ignite/templates/field/plushhelpers"
	"github.com/ignite/cli/v28/ignite/templates/module"

	"github.com/ignite/apps/wasm/pkg/goanalysis"
)

const funcRegisterIBCWasm = `
	modules := map[string]appmodule.AppModule{
		ibcexported.ModuleName:      ibc.AppModule{},
		ibctransfertypes.ModuleName: ibctransfer.AppModule{},
		ibcfeetypes.ModuleName:      ibcfee.AppModule{},
		icatypes.ModuleName:         icamodule.AppModule{},
		capabilitytypes.ModuleName:  capability.AppModule{},
		ibctm.ModuleName:            ibctm.AppModule{},
		solomachine.ModuleName:      solomachine.AppModule{},
		wasmtypes.ModuleName:        wasm.AppModule{},
	}

	for name, m := range modules {
		module.CoreAppModuleBasicAdaptor(name, m).RegisterInterfaces(registry)
	}

	return modules`

//go:embed files/* files/**/*
var fsAppWasm embed.FS

// Options wasm scaffold options.
type Options struct {
	BinaryName string
	AppPath    string
	Home       string
}

// Validate that options are usable.
func (opts *Options) Validate() error {
	return nil
}

// NewWasmGenerator returns the generator to scaffold a wasm integration inside an app.
func NewWasmGenerator(replacer placeholder.Replacer, opts *Options) (*genny.Generator, error) {
	var (
		g       = genny.New()
		appWasm = xgenny.NewEmbedWalker(
			fsAppWasm,
			"files/",
			opts.AppPath,
		)
	)
	if err := g.Box(appWasm); err != nil {
		return g, err
	}

	ctx := plush.NewContext()
	plushhelpers.ExtendPlushContext(ctx)
	g.Transformer(xgenny.Transformer(ctx))

	g.RunFn(appModify(replacer, opts))
	g.RunFn(appConfigModify(replacer, opts))
	g.RunFn(ibcModify(replacer, opts))
	g.RunFn(cmdModify(opts))

	return g, nil
}

// appConfigModify app_config.go modification when adding wasm integration.
func appConfigModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		configPath := filepath.Join(opts.AppPath, module.PathAppConfigGo)
		fConfig, err := r.Disk.Find(configPath)
		if err != nil {
			return err
		}

		// Import
		template := `wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
%[1]v`
		replacement := fmt.Sprintf(template, module.PlaceholderSgAppModuleImport)
		content := replacer.Replace(fConfig.String(), module.PlaceholderSgAppModuleImport, replacement)

		// Init genesis / begin block / end block
		template = `wasmtypes.ModuleName,
%[1]v`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppInitGenesis)
		content = replacer.Replace(content, module.PlaceholderSgAppInitGenesis, replacement)
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppBeginBlockers)
		content = replacer.Replace(content, module.PlaceholderSgAppBeginBlockers, replacement)
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppEndBlockers)
		content = replacer.Replace(content, module.PlaceholderSgAppEndBlockers, replacement)

		// Mac Perms
		template = `{Account: wasmtypes.ModuleName, Permissions: []string{authtypes.Burner}},
%[1]v`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppMaccPerms)
		content = replacer.Replace(content, module.PlaceholderSgAppMaccPerms, replacement)

		return r.File(genny.NewFileS(configPath, content))
	}
}

// appModify app.go modification when adding wasm integration.
func appModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		appPath := filepath.Join(opts.AppPath, module.PathAppGo)
		f, err := r.Disk.Find(appPath)
		if err != nil {
			return err
		}

		// Import
		template := `wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
%[1]v`
		replacement := fmt.Sprintf(template, module.PlaceholderSgAppModuleImport)
		content := replacer.Replace(f.String(), module.PlaceholderSgAppModuleImport, replacement)

		// Keeper declaration
		template = `
// CosmWasm
WasmKeeper       wasmkeeper.Keeper
ScopedWasmKeeper capabilitykeeper.ScopedKeeper

%[1]v`
		replacement = fmt.Sprintf(template, module.PlaceholderSgAppKeeperDeclaration)
		content = replacer.Replace(content, module.PlaceholderSgAppKeeperDeclaration, replacement)

		content, err = goanalysis.ReplaceReturn(
			content,
			"New",
			"app",
			"app.WasmKeeper.InitializePinnedCodes(app.NewUncachedContext(true, tmproto.Header{}))",
		)
		if err != nil {
			return err
		}

		return r.File(genny.NewFileS(appPath, content))
	}
}

// ibcModify ibc.go modification when adding wasm integration.
func ibcModify(replacer placeholder.Replacer, opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		ibcPath := filepath.Join(opts.AppPath, "app/ibc.go")
		f, err := r.Disk.Find(ibcPath)
		if err != nil {
			return err
		}

		if !strings.Contains(f.String(), "registerIBCModules(appOpts servertypes.AppOptions) error") {
			return errors.Errorf("chain does not support wasm integration (CLI >= v28 and Cosmos SDK >= v0.50). See the ignite migration guide")
		}

		// Import
		templateImport := `"github.com/CosmWasm/wasmd/x/wasm"
wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
%[1]v`
		replacementImport := fmt.Sprintf(templateImport, module.PlaceholderIBCImport)
		content := replacer.Replace(f.String(), module.PlaceholderIBCImport, replacementImport)

		// create wasm module
		templateIBCModule := `wasmStack, err := app.registerWasmModules(appOpts)
if err != nil {
	return err
}
ibcRouter.AddRoute(wasmtypes.ModuleName, wasmStack)

%[1]v`
		replacementIBCModule := fmt.Sprintf(templateIBCModule, module.PlaceholderIBCNewModule)
		content = replacer.Replace(content, module.PlaceholderIBCNewModule, replacementIBCModule)

		content, err = goanalysis.ReplaceCode(content, "RegisterIBC", funcRegisterIBCWasm)
		if err != nil {
			return err
		}

		return r.File(genny.NewFileS(ibcPath, content))
	}
}

// cmdModify cmd.go modification when adding wasm integration.
func cmdModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		cmdPath := filepath.Join(opts.AppPath, "cmd", opts.BinaryName, "cmd/commands.go")
		f, err := r.Disk.Find(cmdPath)
		if err != nil {
			return err
		}

		content, err := goanalysis.AppendCode(
			f.String(),
			"initRootCmd",
			"wasmcli.ExtendUnsafeResetAllCmd(rootCmd)",
		)
		if err != nil {
			return err
		}

		content, err = goanalysis.AppendCode(
			content,
			"addModuleInitFlags",
			"wasm.AddModuleInitFlags(startCmd)",
		)
		if err != nil {
			return err
		}

		content, err = goanalysis.AppendImports(
			content,
			"github.com/CosmWasm/wasmd/x/wasm",
			"wasmcli github.com/CosmWasm/wasmd/x/wasm/client/cli",
		)
		if err != nil {
			return err
		}

		return r.File(genny.NewFileS(cmdPath, content))
	}
}
