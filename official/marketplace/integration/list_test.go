package explorer_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pluginsconfig "github.com/ignite/cli/v28/ignite/config/plugins"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	envtest "github.com/ignite/cli/v28/integration"
)

func TestMarketplace(t *testing.T) {
	var (
		require = require.New(t)
		assert  = assert.New(t)
		env     = envtest.New(t)
		app     = env.Scaffold("github.com/test/mp")

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
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "marketplace")

	env.Must(env.Exec("add marketplace plugin",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", "-g", pluginPath),
			step.Workdir(app.SourcePath()),
		)),
	))

	// one local plugin expected
	assertPlugins(
		nil,
		[]pluginsconfig.Plugin{
			{
				Path: pluginPath,
			},
		},
	)

	buf := &bytes.Buffer{}
	env.Must(env.Exec("run marketplace list",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"marketplace",
				"list",
			),
			step.Workdir(app.SourcePath()),
			step.Stdout(buf),
		)),
	))
	assert.Condition(func() bool {
		return strings.HasPrefix(buf.String(), "❌") || strings.HasPrefix(buf.String(), "📦")
	}, "unexpected output: %s", buf.String())
}
