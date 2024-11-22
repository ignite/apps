package cmd

import (
	"context"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

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

	flags := plugin.Flags(cmd.Flags)
	faucetPort, _ := flags.GetUint64(flagFaucetPort)
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

	faucetRestart, err := c.FaucetRestart(ctx)
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
