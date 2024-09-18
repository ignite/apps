package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/v28/integration"
	"github.com/stretchr/testify/require"
)

func TestFeeAbstraction(t *testing.T) {
	var (
		require = require.New(t)
		env     = envtest.New(t)
	)
	dir, err := os.Getwd()
	require.NoError(err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "fee-abstraction")

	env.Must(env.Exec("install fee-abstraction app locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", "-g", pluginPath),
		)),
	))

	app := env.Scaffold("github.com/apps/feeapp", "--fee-abstraction")
	require.FileExists(filepath.Join(app.SourcePath(), "app/feeabs.go"))

	// check the chains is up
	app.EnsureSteady()
}
