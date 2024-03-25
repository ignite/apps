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

func TestMarketplace(t *testing.T) {
	var (
		require = require.New(t)
		env     = envtest.New(t)
	)

	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "marketplace")

	env.Must(env.Exec("add marketplace plugin globally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", "-g", pluginPath),
		)),
	))

	// One local plugin expected
	assertGlobalPlugins(t, []pluginsconfig.Plugin{{Path: pluginPath}})

	listOutput := &bytes.Buffer{}
	env.Must(env.Exec("run marketplace list",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"marketplace",
				"list",
			),
			// all test outputs are going to the stdErr for no reason, but
			// it's ok when we run the app. The output goes to stdout.
			step.Stderr(listOutput),
			step.Stdout(listOutput),
		)),
	))
	gotList := listOutput.String()
	require.True(strings.HasPrefix(gotList, "📦"), "unexpected output: %s", gotList)

	infoOutput := &bytes.Buffer{}
	env.Must(env.Exec("run marketplace info",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"marketplace",
				"info",
				"github.com/ignite/apps",
			),
			// all test outputs are going to the stdErr for no reason, but
			// it's ok when we run the app. The output goes to stdout.
			step.Stderr(infoOutput),
			step.Stdout(infoOutput),
		)),
	))
	gotInfo := infoOutput.String()
	require.Contains(gotInfo, "Description:\tIgnite Apps")
}

func assertGlobalPlugins(t *testing.T, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfgPath, err := plugin.PluginsPath()
	require.NoError(t, err)
	cfg, err := pluginsconfig.ParseDir(cfgPath)
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected global apps")
}
