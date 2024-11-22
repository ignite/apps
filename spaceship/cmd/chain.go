package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gookit/color"
	ignitecmd "github.com/ignite/cli/v28/ignite/cmd"
	config "github.com/ignite/cli/v28/ignite/config/chain"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/schollz/progressbar/v3"

	"github.com/ignite/apps/spaceship/pkg/tarball"
	"github.com/ignite/apps/spaceship/templates/script"
)

// ExecuteSSHStatus executes the ssh status subcommand.
func ExecuteSSHStatus(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusConnecting))
	defer session.End()

	c, err := executeSSH(cmd, chain)
	if err != nil {
		return err
	}
	defer c.Close()

	if !c.HasRunnerScript(ctx) {
		return ErrServerNotInitialized
	}

	chainStatus, err := c.Status(ctx)
	if err != nil {
		return err
	}
	_ = session.Println(chainStatus)

	stopStatus, err := c.FaucetStatus(ctx)
	if err != nil {
		return err
	}

	return session.Println(stopStatus)
}

// ExecuteSSHRestart executes the ssh restart subcommand.
func ExecuteSSHRestart(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusConnecting))
	defer session.End()

	c, err := executeSSH(cmd, chain)
	if err != nil {
		return err
	}
	defer c.Close()

	if !c.HasRunnerScript(ctx) {
		return ErrServerNotInitialized
	}

	chainRestart, err := c.Restart(ctx)
	if err != nil {
		return err
	}
	_ = session.Println(chainRestart)

	faucetRestart, err := c.FaucetRestart(ctx)
	if err != nil {
		return err
	}

	return session.Println(faucetRestart)
}

// ExecuteSSHSStop executes the ssh stop subcommand.
func ExecuteSSHSStop(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusConnecting))
	defer session.End()

	c, err := executeSSH(cmd, chain)
	if err != nil {
		return err
	}
	defer c.Close()

	if !c.HasRunnerScript(ctx) {
		return ErrServerNotInitialized
	}

	chainStop, err := c.Stop(ctx)
	if err != nil {
		return err
	}
	_ = session.Println(chainStop)

	faucetStop, err := c.FaucetStop(ctx)
	if err != nil {
		return err
	}

	return session.Println(faucetStop)
}

// ExecuteSSHDeploy executes the ssh deploy subcommand.
func ExecuteSSHDeploy(ctx context.Context, cmd *plugin.ExecutedCommand, chain *plugin.ChainInfo) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusConnecting))
	defer session.End()

	flags := plugin.Flags(cmd.Flags)

	localDir, err := os.MkdirTemp(os.TempDir(), "spaceship")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(localDir)
	}()

	var (
		initChain, _  = flags.GetBool(flagInitChain)
		faucet, _     = flags.GetBool(flagFaucet)
		faucetPort, _ = flags.GetUint64(flagFaucetPort)

		localChainHome = filepath.Join(localDir, "home")
		localBinOutput = filepath.Join(localDir, "bin")
	)

	c, err := executeSSH(cmd, chain)
	if err != nil {
		return err
	}
	defer c.Close()

	_ = session.Println(color.Yellow.Sprint("Building chain binary using Ignite:"))

	target, err := c.Target(ctx)
	if err != nil {
		return err
	}
	targetName := strings.ReplaceAll(target, ":", "_")

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
		"--verbose",
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
			targetName,
		)
	)
	extracted, err := tarball.ExtractFile(ctx, localChainTarball, localBinOutput, binName)
	if err != nil {
		return err
	}
	if len(extracted) == 0 {
		return errors.Errorf("zero files extracted from the tarball %s", localChainTarball)
	}

	bar.Describe("Uploading chain binary")
	binPath, err := c.UploadBinary(extracted[0], progressCallback)
	if err != nil {
		return err
	}
	_ = session.Println()
	_ = session.Println(color.Yellow.Sprintf("Chain binary uploaded to '%s'\n", binPath))

	home := c.Home()
	if initChain || !c.HasGenesis(ctx) {
		_ = session.Println(color.Yellow.Sprint("Initializing the chain home folder using Ignite:"))

		igniteChainInitCmd := ignitecmd.NewChainInit()
		igniteChainInitCmd.SetArgs([]string{"-p", chain.AppPath, "--home", localChainHome})
		if err := igniteChainInitCmd.ExecuteContext(ctx); err != nil {
			return err
		}

		bar.Describe("Uploading chain home folder")
		homeFiles, err := c.UploadHome(ctx, localChainHome, progressCallback)
		if err != nil {
			return err
		}
		_ = session.Println(color.Yellow.Sprintf("Uploaded files: \n- %s\n", strings.Join(homeFiles, "\n- ")))
	}

	chainCfg, err := chainConfig(chain)
	if err != nil {
		return err
	}

	denom := "token"
	if len(chainCfg.Faucet.Coins) > 0 {
		coin, err := sdk.ParseCoinNormalized(chainCfg.Faucet.Coins[0])
		if err != nil {
			return err
		}
		denom = coin.Denom
	}

	// Create the chain and faucet runner scripts.
	scriptsDir := filepath.Join(localDir, "scripts")
	if err := script.NewRunScripts(
		c.Workspace(),
		c.Log(),
		home,
		binPath,
		*chainCfg.Faucet.Name,
		denom,
		scriptsDir,
	); err != nil {
		return err
	}

	bar.Describe("Uploading runner script")
	if err := c.UploadScripts(ctx, scriptsDir, progressCallback); err != nil {
		return err
	}

	bar.Describe("Uploading faucet binary")
	if _, err := c.UploadFaucetBinary(ctx, targetName, progressCallback); err != nil {
		return err
	}

	_ = session.Println(color.Yellow.Sprintf("Running chain %s", binName))
	start, err := c.Start(ctx)
	if err != nil {
		return err
	}
	_ = session.Println("")
	_ = session.Println(color.Blue.Sprint(start))

	if faucet {
		_ = session.Println(color.Yellow.Sprintf("Running chain %s faucet", binName))
		faucetStart, err := c.FaucetStart(ctx, faucetPort)
		if err != nil {
			return err
		}
		_ = session.Println("")
		return session.Println(color.Blue.Sprint(faucetStart))
	}
	return nil
}

// chainConfig retrieves and parses the configuration for the given chain.
// It first attempts to load the config from chain's specified path, if unsuccessful,
// it attempts to locate and load the default config file.
func chainConfig(chain *plugin.ChainInfo) (*config.Config, error) {
	cfg, err := config.ParseFile(chain.ConfigPath)
	if err == nil {
		return cfg, nil
	}

	cfgPath, err := config.LocateDefault(chain.AppPath)
	if err != nil {
		return nil, err
	}
	return config.ParseFile(cfgPath)
}
