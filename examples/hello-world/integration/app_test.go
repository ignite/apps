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
	)

	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "hello-world")

	env.Must(env.Exec("install hello-world app globally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", "-g", pluginPath),
		)),
	))

	// One local plugin expected
	assertGlobalPlugins(t, []pluginsconfig.Plugin{{Path: pluginPath}})

	buf := &bytes.Buffer{}
	env.Must(env.Exec("run hello-world",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"hello-world",
			),
			step.Stdout(buf),
		)),
	))
	require.Equal("Hello, world!\n", buf.String())
}

func assertGlobalPlugins(t *testing.T, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfgPath, err := plugin.PluginsPath()
	require.NoError(t, err)
	cfg, err := pluginsconfig.ParseDir(cfgPath)
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected global apps")
}
