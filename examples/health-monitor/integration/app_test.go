package integration_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	pluginsconfig "github.com/ignite/cli/v28/ignite/config/plugins"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	envtest "github.com/ignite/cli/v28/integration"
	"github.com/stretchr/testify/require"
)

func TestHealthMonitor(t *testing.T) {
	var (
		require     = require.New(t)
		env         = envtest.New(t)
		app         = env.Scaffold("github.com/apps/health-monitor")
		servers     = app.RandomizeServerPorts()
		ctx, cancel = context.WithCancel(env.Ctx())
	)

	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "health-monitor")

	env.Must(env.Exec("install health-monitor app locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", pluginPath),
			step.Workdir(app.SourcePath()),
		)),
	))

	// One local plugin expected
	assertLocalPlugins(t, app, []pluginsconfig.Plugin{{Path: pluginPath}})
	assertGlobalPlugins(t, nil)

	var (
		isRetrieved         bool
		output              = &bytes.Buffer{}
		stepCtx, stepCancel = context.WithCancel(env.Ctx())
	)
	steps := step.NewSteps(
		step.New(
			step.Stdout(output),
			step.Workdir(app.SourcePath()),
			step.PreExec(func() error {
				return env.IsAppServed(ctx, servers.API)
			}),
			step.Exec(
				envtest.IgniteApp,
				"health-monitor",
				"monitor",
				"--rpc-address", servers.RPC,
				"--refresh-duration", "1s",
			),
			step.InExec(func() error {
				time.Sleep(2 * time.Second)
				stepCancel()
				return nil
			}),
		),
	)

	go func() {
		defer cancel()
		isRetrieved = env.Exec("run health-monitor", steps, envtest.ExecRetry(), envtest.ExecCtx(stepCtx))
	}()

	env.Must(app.Serve("should serve", envtest.ExecCtx(ctx)))

	if !isRetrieved {
		t.FailNow()
	}

	got := output.String()
	require.Contains(got, "Chain ID: healthmonitor")
	require.Contains(got, "Version:")
	require.Contains(got, "Height:")
	require.Contains(got, "Latest Block Hash:")
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
