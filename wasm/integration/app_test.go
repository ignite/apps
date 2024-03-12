package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	pluginsconfig "github.com/ignite/cli/v28/ignite/config/plugins"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	envtest "github.com/ignite/cli/v28/integration"
	"github.com/stretchr/testify/require"
)

func TestWasm(t *testing.T) {
	t.Skip("this tests will only work after we release a new ignite version (>= v29)")

	var (
		require     = require.New(t)
		env         = envtest.New(t)
		app         = env.Scaffold("github.com/apps/wasm-app")
		servers     = app.RandomizeServerPorts()
		ctx, cancel = context.WithCancel(env.Ctx())
	)
	defer cancel()

	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "wasm")

	env.Must(env.Exec("install wasm app locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", pluginPath),
			step.Workdir(app.SourcePath()),
		)),
	))

	// One local plugin expected
	assertLocalPlugins(t, app, []pluginsconfig.Plugin{
		{
			Path: pluginPath,
		},
	})
	assertGlobalPlugins(t, nil)

	env.Must(env.Exec("run wasm",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"wasm",
				"add",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	// sign tx to add an item to the list.
	steps := step.NewSteps(
		step.New(
			step.Exec(
				app.Binary(),
				"config",
				"output", "json",
			),
			step.Workdir(app.SourcePath()),
			step.PreExec(func() error {
				return env.IsAppServed(ctx, servers.API)
			}),
		),
		step.New(
			step.Workdir(app.SourcePath()),
			step.PreExec(func() error {
				err := env.IsAppServed(ctx, servers.API)
				return err
			}),
			step.Exec(
				envtest.IgniteApp,
				"wasm",
				"config",
			),
		),
	)

	isBodyRetrieved := false
	go func() {
		defer cancel()
		isBodyRetrieved = env.Exec("add wasm to the config", steps, envtest.ExecRetry())
	}()

	env.Must(app.Serve("should serve", envtest.ExecCtx(ctx)))

	if !isBodyRetrieved {
		t.FailNow()
	}
}

func assertLocalPlugins(t *testing.T, app envtest.App, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfg, err := pluginsconfig.ParseDir(app.SourcePath())
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected local apps")
}

func assertGlobalPlugins(t *testing.T, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfgPath, err := plugin.PluginsPath()
	require.NoError(t, err)
	cfg, err := pluginsconfig.ParseDir(cfgPath)
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected global apps")
}
