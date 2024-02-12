package integration_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	pluginsconfig "github.com/ignite/cli/v28/ignite/config/plugins"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	envtest "github.com/ignite/cli/v28/integration"
)

func TestHealthMonitor(t *testing.T) {
	var (
		require = require.New(t)
		env     = envtest.New(t)
		app     = env.Scaffold("github.com/test/test")
		ports   = app.RandomizeServerPorts()
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
	assertLocalPlugins(t, app, []pluginsconfig.Plugin{
		{
			Path: pluginPath,
		},
	})
	assertGlobalPlugins(t, app, nil)

	go func() {
		env.Must(app.Serve("serve example chain"))
	}()

	// Wait for the chain to start
	err = env.IsAppServed(context.Background(), ports.API)
	require.NoError(err)

	buf := &bytes.Buffer{}
	go func() {
		env.Must(env.Exec("run health-monitor",
			step.NewSteps(step.New(
				step.Exec(
					envtest.IgniteApp,
					"health-monitor",
					"monitor",
					"--rpc-address", ports.RPC,
					"--refresh-duration", "1s",
				),
				step.Workdir(app.SourcePath()),
				step.Stdout(buf),
			)),
		))
	}()
	time.Sleep(time.Second * 2)
	require.Contains(buf.String(), "Chain ID: test")
	require.Contains(buf.String(), "Version:")
	require.Contains(buf.String(), "Height:")
	require.Contains(buf.String(), "Latest Block Hash:")

}

func assertLocalPlugins(t *testing.T, app envtest.App, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfg, err := pluginsconfig.ParseDir(app.SourcePath())
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected local apps")
}

func assertGlobalPlugins(t *testing.T, app envtest.App, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfgPath, err := plugin.PluginsPath()
	require.NoError(t, err)
	cfg, err := pluginsconfig.ParseDir(cfgPath)
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected global apps")
}
