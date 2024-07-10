package integration_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	pluginsconfig "github.com/ignite/cli/v28/ignite/config/plugins"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	envtest "github.com/ignite/cli/v28/integration"
)

func TestAppRegistry(t *testing.T) {
	var (
		require = require.New(t)
		env     = envtest.New(t)
	)

	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "appregistry")

	env.Must(env.Exec("add appregistry plugin globally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", "-g", pluginPath),
		)),
	))

	// One local plugin expected
	assertGlobalPlugins(t, []pluginsconfig.Plugin{{Path: pluginPath}})

	listOutput := &bytes.Buffer{}
	env.Must(env.Exec("run appregistry list",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"appregistry",
				"list",
			),
			// all test outputs are going to the stdErr for no reason, but
			// it's ok when we run the app. The output goes to stdout.
			step.Stderr(listOutput),
			step.Stdout(listOutput),
		)),
	))
	gotList := listOutput.String()
	require.True(strings.Contains(gotList, "appregistry :"), "unexpected output: %s", gotList)

	infoOutput := &bytes.Buffer{}
	env.Must(env.Exec("run appregistry info",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"appregistry",
				"info",
				"explorer",
			),
			// all test outputs are going to the stdErr for no reason, but
			// it's ok when we run the app. The output goes to stdout.
			step.Stderr(infoOutput),
			step.Stdout(infoOutput),
		)),
	))
	gotInfo := infoOutput.String()
	require.Contains(gotInfo, "Description:\tEasy to use terminal chain explorer for testing your Ignite blockchains")
}

func assertGlobalPlugins(t *testing.T, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfgPath, err := plugin.PluginsPath()
	require.NoError(t, err)
	cfg, err := pluginsconfig.ParseDir(cfgPath)
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected global apps")
}
