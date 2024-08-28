package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
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

	var (
		app         = env.Scaffold("github.com/apps/feeabs", "--fee-abstraction")
		servers     = app.RandomizeServerPorts()
		ctx, cancel = context.WithCancel(env.Ctx())
	)
	defer cancel()

	require.FileExists(filepath.Join(app.SourcePath(), "app/feeabs.go"))

	// check the chains is up
	stepUp := step.NewSteps(
		step.New(
			step.Exec(
				app.Binary(),
				"config",
				"output", "json",
			),
			step.PreExec(func() error {
				return env.IsAppServed(ctx, servers.API)
			}),
			step.Workdir(app.SourcePath()),
		),
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		env.Exec("waiting the chain is up", stepUp, envtest.ExecRetry())
		wg.Done()
	}()

	env.Must(app.Serve("should serve", envtest.ExecCtx(ctx)))
	wg.Wait()
}
