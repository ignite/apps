package cmd

import (
	"context"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	"golang.org/x/sync/errgroup"

	"github.com/ignite/apps/spaceship/pkg/ssh"
)

func ExecuteChainSSHLog(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	return ExecuteSSHLog(ctx, ssh.LogChain, cmd, chain)
}

func ExecuteFaucetSSHLog(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	return ExecuteSSHLog(ctx, ssh.LogFaucet, cmd, chain)
}

// ExecuteSSHLog executes the ssh log subcommand.
func ExecuteSSHLog(ctx context.Context, logType ssh.LogType, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusConnecting))
	defer session.End()

	var (
		flags       = plugin.Flags(cmd.Flags)
		lines, _    = flags.GetInt(flagLines)
		realTime, _ = flags.GetBool(flagRealTime)
	)

	c, err := executeSSH(cmd, chain)
	if err != nil {
		return err
	}
	defer c.Close()

	if !c.HasRunnerScript(ctx) {
		return ErrServerNotInitialized
	}

	logs, err := c.LatestLog(logType, lines)
	if err != nil {
		return err
	}
	_ = session.Println(logs)

	if realTime {
		// Create a buffered channel to receive log lines.
		logChannel := make(chan string, 100)
		g, ctx := errgroup.WithContext(ctx)

		// Start the FollowLog method in a goroutine using errgroup
		g.Go(func() error {
			return c.FollowLog(ctx, logType, logChannel)
		})

		// Start a goroutine to consume log lines
		g.Go(func() error {
			for {
				select {
				case line := <-logChannel:
					_ = session.Print(line)
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})

		// Wait for all goroutines to complete
		if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
			return err
		}
	}

	return nil
}
