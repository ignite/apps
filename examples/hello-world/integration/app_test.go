package integration_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	pluginsconfig "github.com/ignite/cli/v28/ignite/config/plugins"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	envtest "github.com/ignite/cli/v28/integration"
)

func TestHelloWorld(t *testing.T) {
	var (
		require = require.New(t)
		env     = envtest.New(t)
		app     = env.Scaffold("github.com/test/test")
	)

	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "hello-world")

	env.Must(env.Exec("install hello-world app locally",
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

	buf := &bytes.Buffer{}
	env.Must(env.Exec("run hello-world",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"hello-world",
			),
			step.Workdir(app.SourcePath()),
			step.Stdout(buf),
		)),
	))
	require.Equal("Hello, world!\n", buf.String())
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
