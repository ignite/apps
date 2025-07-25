package integration_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	pluginsconfig "github.com/ignite/cli/v29/ignite/config/plugins"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/services/plugin"
	envtest "github.com/ignite/cli/v29/integration"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestGexExplorer(t *testing.T) {
	var (
		require     = require.New(t)
		env         = envtest.New(t)
		app         = env.ScaffoldApp("github.com/test/explorer")
		servers     = app.RandomizeServerPorts()
		ctx, cancel = context.WithCancel(env.Ctx())
	)

	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "explorer")

	env.Must(env.Exec("add explorer plugin locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", pluginPath),
			step.Workdir(app.SourcePath()),
		)),
	))

	// One local plugin expected
	assertLocalPlugins(t, app, []pluginsconfig.Plugin{{Path: pluginPath}})
	assertGlobalPlugins(t, nil)

	b := &bytes.Buffer{}
	steps := step.NewSteps(
		step.New(
			step.Stdout(b),
			step.Workdir(app.SourcePath()),
			step.PreExec(func() error {
				return env.IsAppServed(ctx, servers.API)
			}),
			step.Exec(envtest.IgniteApp, "e", "gex", "--rpc-address", servers.RPC),
			step.InExec(func() error {
				defer cancel()
				time.Sleep(5 * time.Second)
				return nil
			}),
		),
	)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		env.Must(env.Exec("run gex", steps, envtest.ExecRetry(), envtest.ExecCtx(ctx)))
		return nil
	})

	env.Must(app.Serve("serve app", envtest.ExecCtx(ctx)))
	require.NoError(g.Wait())
}

func TestPingPubExplorer(t *testing.T) {
	var (
		require     = require.New(t)
		env         = envtest.New(t)
		app         = env.ScaffoldApp("github.com/test/pingpub-explorer")
		servers     = app.RandomizeServerPorts()
		ctx, cancel = context.WithCancel(env.Ctx())
	)

	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "explorer")

	env.Must(env.Exec("add explorer plugin locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", pluginPath),
			step.Workdir(app.SourcePath()),
		)),
	))

	// One local plugin expected
	assertLocalPlugins(t, app, []pluginsconfig.Plugin{{Path: pluginPath}})
	assertGlobalPlugins(t, nil)

	// Skip test if yarn is not installed
	if _, err := os.Stat("/usr/bin/yarn"); os.IsNotExist(err) {
		t.Skip("yarn not installed, skipping pingpub test")
	}

	execErr := &bytes.Buffer{}
	steps := step.NewSteps(
		step.New(
			step.Stderr(execErr),
			step.Workdir(app.SourcePath()),
			step.PreExec(func() error {
				return env.IsAppServed(ctx, servers.API)
			}),
			step.Exec(envtest.IgniteApp, "e", "pingpub", "--path", app.SourcePath()),
			step.InExec(func() error {
				time.Sleep(15 * time.Second) // Give more time for pingpub to start
				cancel()
				return nil
			}),
		),
	)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		env.Must(env.Exec("run pingpub", steps, envtest.ExecRetry(), envtest.ExecCtx(ctx)))
		return nil
	})

	app.MustServe(ctx)
	require.NoError(g.Wait())

	require.Empty(execErr.String())

	// Check if explorer directory was created
	explorerDir := filepath.Join(app.SourcePath(), "explorer", "ping-pub")
	_, err = os.Stat(explorerDir)
	require.NoError(err, "Explorer directory should be created")

	// Check if configuration file was created
	chainName := filepath.Base(app.SourcePath())
	configFile := filepath.Join(explorerDir, "chains", "mainnet", chainName+".json")
	_, err = os.Stat(configFile)
	require.NoError(err, "Configuration file should be created")
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
