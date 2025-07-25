package integration_test

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	pluginsconfig "github.com/ignite/cli/v29/ignite/config/plugins"
	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v29/ignite/pkg/goenv"
	"github.com/ignite/cli/v29/ignite/services/plugin"
	envtest "github.com/ignite/cli/v29/integration"
	"github.com/stretchr/testify/require"
)

func TestRollkit(t *testing.T) {
	var (
		require = require.New(t)
		env     = envtest.New(t)
		app     = env.ScaffoldApp("github.com/apps/rollkit")
	)

	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "rollkit")

	env.Must(env.Exec("install rollkit app locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", pluginPath),
			step.Workdir(app.SourcePath()),
		)),
	))

	// One local plugin expected
	assertLocalPlugins(t, app, []pluginsconfig.Plugin{{Path: pluginPath}})
	assertGlobalPlugins(t, nil)

	env.Must(env.Exec("run rollkit add",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"rollkit",
				"add",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	buf := &bytes.Buffer{}
	bin := path.Join(goenv.Bin(), app.Binary())
	env.Must(env.Exec("check rollkitd", step.NewSteps(
		step.New(
			step.Exec(
				envtest.IgniteApp,
				"chain",
				"build",
			),
			step.Workdir(app.SourcePath()),
		),
		step.New(
			step.Exec(bin, "start", "--help"),
			step.Stdout(buf),
			step.Workdir(app.SourcePath()),
		)),
	))

	if !strings.Contains(buf.String(), "--rollkit.da") {
		t.Errorf("rollkitd doesn't contain --rollkit flags: %s", buf.String())
	}

	buf.Reset()

	env.Must(env.Exec("run rollkit init", step.NewSteps(
		step.New(
			step.Exec(
				envtest.IgniteApp,
				"rollkit",
				"init",
			),
			step.PostExec(func(exitErr error) error {
				return os.Remove(bin)
			}),
			step.Workdir(app.SourcePath()),
			step.Stdout(buf),
		),
	)))

	if !strings.Contains(buf.String(), "Initialized. Checkout your rollkit chain's home") {
		t.Errorf("ignite rollkit init has failed: %s", buf.String())
	}
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
