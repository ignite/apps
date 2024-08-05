package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ignitecmd "github.com/ignite/cli/v28/ignite/cmd"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/schollz/progressbar/v3"

	"github.com/ignite/apps/spaceship/pkg/ssh"
	"github.com/ignite/apps/spaceship/pkg/tarball"
	"github.com/ignite/apps/spaceship/templates/script"
)

var ErrServerNotInitialized = errors.New("server not initialized")

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

	var (
		host = args[0]

		flags          = plugin.Flags(cmd.Flags)
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

	if !c.HasRunnerScript(ctx) {
		return ErrServerNotInitialized
	}

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

	if !c.HasRunnerScript(ctx) {
		return ErrServerNotInitialized
	}

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

	if !c.HasRunnerScript(ctx) {
		return ErrServerNotInitialized
	}

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

	if !c.HasRunnerScript(ctx) {
		return ErrServerNotInitialized
	}

	log, err := c.LatestLog()
	if err != nil {
		return err
	}
	fmt.Println(string(log))
	return nil
}

// ExecuteSSHDeploy executes the ssh deploy subcommand.
func ExecuteSSHDeploy(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	flags := plugin.Flags(cmd.Flags)

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
	)

	c, err := executeSSH(cmd, chain)
	if err != nil {
		return err
	}
	defer c.Close()

	target, err := c.Target(ctx)
	if err != nil {
		return err
	}

	// We are using the ignite chain build command to build the app.
	igniteChainBuildCmd := ignitecmd.NewChainBuild()
	igniteChainBuildCmd.SetArgs([]string{
		"-p",
		chain.AppPath,
		"-o",
		localBinOutput,
		"--release",
		"--release.targets",
		target,
		"-v",
	})
	if err := igniteChainBuildCmd.ExecuteContext(ctx); err != nil {
		return err
	}

	// Define a progress callback function.
	bar := progressbar.DefaultBytes(1, "uploading")
	progressCallback := func(bytesUploaded int64, totalBytes int64) error {
		if bar.GetMax64() != totalBytes {
			bar.ChangeMax64(totalBytes)
			bar.Reset()
		}
		if err := bar.Set64(bytesUploaded); err != nil {
			return err
		}
		if bytesUploaded == totalBytes {
			if err := bar.Finish(); err != nil {
				return err
			}
		}
		return nil
	}

	// Extract and upload the built binary.
	var (
		binName           = fmt.Sprintf("%sd", chain.ChainId)
		localChainTarball = fmt.Sprintf(
			"%s/%s_%s.tar.gz",
			localBinOutput,
			chain.ChainId,
			strings.ReplaceAll(target, ":", "_"),
		)
	)
	extracted, err := tarball.Extract(ctx, localChainTarball, localBinOutput, binName)
	if err != nil {
		return err
	}
	if len(extracted) == 0 {
		return errors.Errorf("zero files extracted from the tarball %s", localChainTarball)
	}

	bar.Describe("uploading chain binary")
	binPath, err := c.UploadBinary(extracted[0], progressCallback)
	if err != nil {
		return err
	}

	home := c.Home()
	if initChain || !c.HasGenesis(ctx) {
		// Init the chain.
		igniteChainInitCmd := ignitecmd.NewChainInit()
		igniteChainInitCmd.SetArgs([]string{"-p", chain.AppPath, "--home", localChainHome}) // TODO add verbose flag after merge and backport this one https://github.com/ignite/cli/pull/4286
		if err := igniteChainInitCmd.ExecuteContext(ctx); err != nil {
			return err
		}

		bar.Describe("uploading chain home folder")
		home, err = c.UploadHome(ctx, localChainHome, progressCallback)
		if err != nil {
			return err
		}
	}

	// Create the runner script.
	localRunScriptPath, err := script.NewRunScript(c.Workspace(), c.Log(), home, binPath, localDir)
	if err != nil {
		return err
	}

	bar.Describe("uploading runner script")
	if _, err := c.UploadRunnerScript(localRunScriptPath, progressCallback); err != nil {
		return err
	}

	start, err := c.Start(ctx)
	if err != nil {
		return err
	}
	fmt.Println(start)

	return nil
}
