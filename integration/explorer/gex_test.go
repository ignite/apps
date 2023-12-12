package explorer_test

import (
	"os"
	"path/filepath"
	"testing"

	pluginsconfig "github.com/ignite/cli/v28/ignite/config/plugins"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	envtest "github.com/ignite/cli/v28/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGexExplorer(t *testing.T) {
	var (
		require = require.New(t)
		assert  = assert.New(t)
		env     = envtest.New(t)
		app     = env.Scaffold("github.com/test/explorer")

		assertPlugins = func(expectedLocalPlugins, expectedGlobalPlugins []pluginsconfig.Plugin) {
			localCfg, err := pluginsconfig.ParseDir(app.SourcePath())
			require.NoError(err)
			assert.ElementsMatch(expectedLocalPlugins, localCfg.Apps, "unexpected local apps")

			globalCfgPath, err := plugin.PluginsPath()
			require.NoError(err)
			globalCfg, err := pluginsconfig.ParseDir(globalCfgPath)
			require.NoError(err)
			assert.ElementsMatch(expectedGlobalPlugins, globalCfg.Apps, "unexpected global apps")
		}
	)

	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "explorer")

	env.Must(env.Exec("add explorer plugin",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", pluginPath),
			step.Workdir(app.SourcePath()),
		)),
	))

	// one local plugin expected
	assertPlugins(
		[]pluginsconfig.Plugin{
			{
				Path: pluginPath,
			},
		},
		nil,
	)

	env.Must(env.Exec("run gex explorer help",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"e",
				"gex",
				"--help",
			),
			step.Workdir(app.SourcePath()),
		)),
	))
}
