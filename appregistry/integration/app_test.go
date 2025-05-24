package integration_test

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	pluginsconfig "github.com/ignite/cli/v28/ignite/config/plugins"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	envtest "github.com/ignite/cli/v28/integration"
	"github.com/stretchr/testify/require"
)

var flagBranch = flag.String("branch", "", "The app branch to use")

func TestAppRegistry(t *testing.T) {
	flag.Parse()
	var (
		require   = require.New(t)
		env       = envtest.New(t)
		gitBranch = getGitBranch(t)
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
				"--branch", gitBranch,
			),
			// all test outputs are going to the stdErr for no reason, but
			// it's ok when we run the app. The output goes to stdout.
			step.Stderr(listOutput),
			step.Stdout(listOutput),
		)),
	))
	gotList := listOutput.String()
	require.True(strings.Contains(gotList, "App Registry (id: appregistry):"), "unexpected output: %s", gotList)

	infoOutput := &bytes.Buffer{}
	env.Must(env.Exec("run appregistry info",
		step.NewSteps(step.New(
			step.Exec(
				envtest.IgniteApp,
				"appregistry",
				"info",
				"explorer",
				"--branch", gitBranch,
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

func getGitBranch(t *testing.T) string {
	t.Helper()
	gitBranch := *flagBranch
	if gitBranch != "" {
		return gitBranch
	}

	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	require.NoError(t, err)
	branch := strings.TrimSpace(string(output))
	return branch
}

func assertGlobalPlugins(t *testing.T, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfgPath, err := plugin.PluginsPath()
	require.NoError(t, err)
	cfg, err := pluginsconfig.ParseDir(cfgPath)
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected global apps")
}
