package cmd

import (
	"context"

	"github.com/ignite/cli/v28/ignite/pkg/availableport"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

func faucetPort(f []*plugin.Flag) (uint64, error) {
	flags := plugin.Flags(f)
	port, err := flags.GetUint64(flagFaucetPort)
	if err != nil {
		return 0, err
	}
	if port == 0 {
		return availablePort()
	}
	return port, nil
}

func availablePort() (uint64, error) {
	ports, err := availableport.Find(
		1,
		availableport.WithMinPort(8000),
		availableport.WithMaxPort(9000),
	)
	if err != nil {
		return 0, err
	}
	if len(ports) == 0 {
		return 0, errors.New("no available ports")
	}
	return uint64(ports[0]), nil
}

// ExecuteSSHFaucetStatus executes the ssh faucet status subcommand.
func ExecuteSSHFaucetStatus(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusConnecting))
	defer session.End()

	c, err := executeSSH(cmd, chain)
	if err != nil {
		return err
	}
	defer c.Close()

	if !c.HasFaucetScript(ctx) {
		return ErrServerNotInitialized
	}

	stopStatus, err := c.FaucetStatus(ctx)
	if err != nil {
		return err
	}

	return session.Println(stopStatus)
}

// ExecuteSSHFaucetStart executes the ssh faucet start subcommand.
func ExecuteSSHFaucetStart(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusConnecting))
	defer session.End()

	c, err := executeSSH(cmd, chain)
	if err != nil {
		return err
	}
	defer c.Close()

	if !c.HasFaucetScript(ctx) {
		return ErrServerNotInitialized
	}

	faucetPort, err := faucetPort(cmd.Flags)
	if err != nil {
		return err
	}
	faucetRestart, err := c.FaucetStart(ctx, faucetPort)
	if err != nil {
		return err
	}

	return session.Println(faucetRestart)
}

// ExecuteSSHFaucetRestart executes the ssh faucet restart subcommand.
func ExecuteSSHFaucetRestart(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusConnecting))
	defer session.End()

	c, err := executeSSH(cmd, chain)
	if err != nil {
		return err
	}
	defer c.Close()

	if !c.HasFaucetScript(ctx) {
		return ErrServerNotInitialized
	}

	faucetPort, err := faucetPort(cmd.Flags)
	if err != nil {
		return err
	}
	faucetRestart, err := c.FaucetRestart(ctx, faucetPort)
	if err != nil {
		return err
	}

	return session.Println(faucetRestart)
}

// ExecuteSSHSFaucetStop executes the ssh faucet stop subcommand.
func ExecuteSSHSFaucetStop(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusConnecting))
	defer session.End()

	c, err := executeSSH(cmd, chain)
	if err != nil {
		return err
	}
	defer c.Close()

	if !c.HasFaucetScript(ctx) {
		return ErrServerNotInitialized
	}

	faucetStop, err := c.FaucetStop(ctx)
	if err != nil {
		return err
	}

	return session.Println(faucetStop)
}
