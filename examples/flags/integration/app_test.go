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

func TestFlags(t *testing.T) {
	var (
		require = require.New(t)
		env     = envtest.New(t)
	)

	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "flags")

	env.Must(env.Exec("install flags app globally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", "-g", pluginPath),
		)),
	))

	// One local plugin expected
	assertGlobalPlugins(t, []pluginsconfig.Plugin{{Path: pluginPath}})

	buf := &bytes.Buffer{}
	env.Must(env.Exec("run hello",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"flags",
				"hello",
				"--name",
				"Test",
			),
			step.Stdout(buf),
		)),
	))
	require.Equal("Hello, Test!\n", buf.String())

	buf.Reset()
	env.Must(env.Exec("run cowsay",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"flags",
				"cowsay",
				"--name",
				"Test",
			),
			step.Stdout(buf),
		)),
	))
	require.Contains(buf.String(), "Hello, Test!")
}

func assertGlobalPlugins(t *testing.T, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfgPath, err := plugin.PluginsPath()
	require.NoError(t, err)
	cfg, err := pluginsconfig.ParseDir(cfgPath)
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected global apps")
}
