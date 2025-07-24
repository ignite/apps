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

func TestEvolve(t *testing.T) {
	var (
		require = require.New(t)
		env     = envtest.New(t)
		app     = env.ScaffoldApp("github.com/apps/evolve")
	)

	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "evolve")

	env.Must(env.Exec("install evolve app locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", pluginPath),
			step.Workdir(app.SourcePath()),
		)),
	))

	// One local plugin expected
	assertLocalPlugins(t, app, []pluginsconfig.Plugin{{Path: pluginPath}})
	assertGlobalPlugins(t, nil)

	env.Must(env.Exec("run evolve add",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"evolve",
				"add",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	buf := &bytes.Buffer{}
	bin := path.Join(goenv.Bin(), app.Binary())
	env.Must(env.Exec("check evolved", step.NewSteps(
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
		t.Errorf("evolved doesn't contain --rollkit flags: %s", buf.String())
	}

	buf.Reset()

	env.Must(env.Exec("run evolve init", step.NewSteps(
		step.New(
			step.Exec(
				envtest.IgniteApp,
				"evolve",
				"init",
			),
			step.PostExec(func(exitErr error) error {
				return os.Remove(bin)
			}),
			step.Workdir(app.SourcePath()),
			step.Stdout(buf),
		),
	)))

	if !strings.Contains(buf.String(), "Initialized. Checkout your evolve chain's home") {
		t.Errorf("ignite evolve init has failed: %s", buf.String())
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
