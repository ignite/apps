package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	ignitecmd "github.com/ignite/cli/v28/ignite/cmd"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/spaceship/pkg/ssh"
	"github.com/ignite/apps/spaceship/templates/script"
)

// ExecuteSSHDeploy executes the ssh deploy subcommand.
func ExecuteSSHDeploy(ctx context.Context, chain *plugin.ChainInfo) error {
	// args := os.Args[2:]
	var (
		host = "danilopantani@127.0.0.1"          // arg host or URI
		key  = "/Users/danilopantani/.ssh/id_rsa" // flag key
		// user = "danilopantani"                    // flag user
		// password = ""                          // flag password
		// port     = "22" // flag port
		// keyPassword = args[5] // flag key password
		// keyRaw      = args[6] // flag key raw
		initChain = true // flag key raw

		localDir       = filepath.Join(os.TempDir(), "spaceship", chain.ChainId)
		localChainHome = filepath.Join(localDir, "home")
		localBinOutput = filepath.Join(localDir, "bin")
		localChainBin  = fmt.Sprintf("%s/%sd", localBinOutput, chain.ChainId)
	)

	// Connect to the SSH.
	c, err := ssh.New(host, ssh.WithKey(key), ssh.WithWorkspace(chain.ChainId))
	if err != nil {
		return err
	}
	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Close()

	// We are using the ignite chain build command to build the app.
	igniteChainBuildCmd := ignitecmd.NewChainBuild()
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
	if initChain {
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
	localRunScriptPath, err := script.NewRunScript(c.Workspace(), home, binPath, localDir)
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

	status, err := c.Status(ctx)
	if err != nil {
		return err
	}
	fmt.Println(status)
	return nil
}
