package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ignite/cli/v29/ignite/pkg/cmdrunner/step"
	envtest "github.com/ignite/cli/v29/integration"
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

	var (
		app     = env.Scaffold("github.com/apps/feeapp", "--fee-abstraction")
		servers = app.RandomizeServerPorts()
	)

	// check if fee abstraction file was scaffolded
	require.FileExists(filepath.Join(app.SourcePath(), "app/feeabs.go"))

	// check the chains is up
	app.EnsureSteady()

	ctx, cancel := context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
	defer cancel()

	go func() {
		app.Serve("serve app", envtest.ExecCtx(ctx))
	}()

	// Wait for the server to be up before running the client tests
	require.NoError(env.IsAppServed(ctx, servers.API))
}
