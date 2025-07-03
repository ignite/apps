package integration_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	pluginsconfig "github.com/ignite/cli/v29/ignite/config/plugins"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/services/plugin"
	envtest "github.com/ignite/cli/v29/integration"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	var (
		require = require.New(t)
		env     = envtest.New(t)
	)

	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "connect")

	env.Must(env.Exec("install connect app globally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", "-g", pluginPath),
		)),
	))

	assertGlobalPlugins(t, []pluginsconfig.Plugin{
		{
			Path: pluginPath,
		},
	})

	buf := &bytes.Buffer{}
	env.Must(env.Exec("run connect",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"connect",
			),
			step.Stdout(buf),
			step.Stderr(buf),
		)),
	))

	require.Contains(buf.String(), "Connect allows you to interact with any Cosmos SDK based blockchain.")
}

func assertGlobalPlugins(t *testing.T, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfgPath, err := plugin.PluginsPath()
	require.NoError(t, err)
	cfg, err := pluginsconfig.ParseDir(cfgPath)
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected global apps")
}
