package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	ignitecmd "github.com/ignite/cli/v28/ignite/cmd"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/spaceship/pkg/ssh"
	"github.com/ignite/apps/spaceship/templates/script"
)

const (
	flagPort        = "port"
	flagUser        = "user"
	flagPassword    = "password"
	flagKey         = "key"
	flagRawKey      = "raw-key"
	flagKeyPassword = "key-password"
	flagInitChain   = "init-chain"
)

func executeSSH(cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) (*ssh.SSH, error) {
	args := cmd.Args
	if len(args) < 1 {
		return nil, errors.New("must specify unless a uri host")
	}
	flags, err := cmd.NewFlags()
	if err != nil {
		return nil, err
	}

	var (
		host = args[0]

		user, _        = flags.GetString(flagUser)
		port, _        = flags.GetString(flagPort)
		password, _    = flags.GetString(flagPassword)
		key, _         = flags.GetString(flagKey)
		rawKey, _      = flags.GetString(flagRawKey)
		keyPassword, _ = flags.GetString(flagKeyPassword)
	)

	// Connect to the SSH.
	c, err := ssh.New(
		host,
		ssh.WithUser(user),
		ssh.WithPort(port),
		ssh.WithPassword(password),
		ssh.WithKey(key),
		ssh.WithRawKey(rawKey),
		ssh.WithKeyPassword(keyPassword),
		ssh.WithWorkspace(chain.ChainId),
	)
	if err != nil {
		return nil, err
	}

	return c, c.Connect()
}

// ExecuteSSHStatus executes the ssh status subcommand.
func ExecuteSSHStatus(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	c, err := executeSSH(cmd, chain)
	if err != nil {
		return err
	}
	defer c.Close()
	status, err := c.Status(ctx)
	if err != nil {
		return err
	}
	fmt.Println(status)
	return nil
}

// ExecuteSSHSStop executes the ssh stop subcommand.
func ExecuteSSHSStop(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	c, err := executeSSH(cmd, chain)
	if err != nil {
		return err
	}
	defer c.Close()
	stop, err := c.Stop(ctx)
	if err != nil {
		return err
	}
	fmt.Println(stop)
	return nil
}

// ExecuteSSHRestart executes the ssh restart subcommand.
func ExecuteSSHRestart(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	c, err := executeSSH(cmd, chain)
	if err != nil {
		return err
	}
	defer c.Close()
	restart, err := c.Restart(ctx)
	if err != nil {
		return err
	}
	fmt.Println(restart)
	return nil
}

// ExecuteSSHLog executes the ssh log subcommand.
func ExecuteSSHLog(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	c, err := executeSSH(cmd, chain)
	if err != nil {
		return err
	}
	defer c.Close()
	log, err := c.LatestLog()
	if err != nil {
		return err
	}
	fmt.Println(string(log))
	return nil
}

// ExecuteSSHDeploy executes the ssh deploy subcommand.
func ExecuteSSHDeploy(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	flags, err := cmd.NewFlags()

	localDir, err := os.MkdirTemp(os.TempDir(), "spaceship")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(localDir)
	}()

	var (
		initChain, _ = flags.GetBool(flagInitChain)

		localChainHome = filepath.Join(localDir, "home")
		localBinOutput = filepath.Join(localDir, "bin")
		localChainBin  = fmt.Sprintf("%s/%sd", localBinOutput, chain.ChainId)
	)

	c, err := executeSSH(cmd, chain)
	if err != nil {
		return err
	}
	defer c.Close()

	//os, err := c.Target(ctx)
	//if err != nil {
	//	return err
	//}

	// We are using the ignite chain build command to build the app.
	igniteChainBuildCmd := ignitecmd.NewChainBuild()
	// igniteChainBuildCmd.SetArgs([]string{"-p", chain.AppPath, "-o", localBinOutput, "--release", "--release.targets", os})
	igniteChainBuildCmd.SetArgs([]string{"-p", chain.AppPath, "-o", localBinOutput})
	if err := igniteChainBuildCmd.ExecuteContext(ctx); err != nil {
		return err
	}

	// Upload the built binary.
	binPath, err := c.UploadBinary(localChainBin)
	if err != nil {
		return err
	}

	home := c.Home()
	if initChain || !c.HasInitialized(ctx) {
		// Init the chain.
		igniteChainInitCmd := ignitecmd.NewChainInit()
		igniteChainInitCmd.SetArgs([]string{"-p", chain.AppPath, "--home", localChainHome})
		if err := igniteChainInitCmd.ExecuteContext(ctx); err != nil {
			return err
		}

		home, err = c.UploadHome(ctx, localChainHome)
		if err != nil {
			return err
		}
	}

	// Create the runner script.
	localRunScriptPath, err := script.NewRunScript(c.Workspace(), c.Log(), home, binPath, localDir)
	if err != nil {
		return err
	}

	if _, err := c.UploadRunnerScript(localRunScriptPath); err != nil {
		return err
	}

	start, err := c.Start(ctx)
	if err != nil {
		return err
	}
	fmt.Println(start)
	return nil
}
