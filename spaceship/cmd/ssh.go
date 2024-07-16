package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	ignitecmd "github.com/ignite/cli/v28/ignite/cmd"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/spaceship/pkg/ssh"
)

func ExecuteSSHDevelopment(ctx context.Context, chain *plugin.ChainInfo) error {
	var (
		host = "danilopantani@127.0.0.1"          // arg host or URI
		key  = "/Users/danilopantani/.ssh/id_rsa" // flag key
		// user = "danilopantani"                    // flag user
		// password = ""                          // flag password
		// port     = "22" // flag port
		// keyPassword = args[5] // flag key password
		// keyRaw      = args[6] // flag key raw
		args = []string{"chain", "build"} // arg ignite cmd
	)

	c, err := ssh.New(host, ssh.WithKey(key), ssh.WithWorkspace(chain.ChainId))
	if err != nil {
		return err
	}
	if err := c.Connect(ctx); err != nil {
		return err
	}
	defer c.Close()

	srcPath, err := c.UploadSource(ctx, chain.AppPath)
	if err != nil {
		return err
	}

	args = append(args, "-p", srcPath)
	out, err := c.RunIgniteCommand(ctx, args...)
	if err != nil {
		return err
	}

	fmt.Println(out)
	return nil
}

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

		localDir  = filepath.Join(os.TempDir(), "spaceship")
		binOutput = filepath.Join(localDir, "bin")
		chainBin  = fmt.Sprintf("%s/%sd", binOutput, chain.ChainId)
		chainHome = filepath.Join(localDir, "home")
	)

	// we are using the ignite chain build command to build the app.
	igniteChainBuildCmd := ignitecmd.NewChainBuild()
	igniteChainBuildCmd.SetArgs([]string{"-p", chain.AppPath, "-o", binOutput, "-y"})
	if err := igniteChainBuildCmd.ExecuteContext(ctx); err != nil {
		return err
	}

	// init the chain
	igniteChainInitCmd := ignitecmd.NewChainInit()
	igniteChainInitCmd.SetArgs([]string{"-p", chain.AppPath, "-h", chainHome, "-y"})
	if err := igniteChainInitCmd.ExecuteContext(ctx); err != nil {
		return err
	}

	c, err := ssh.New(host, ssh.WithKey(key))
	if err != nil {
		return err
	}
	if err := c.Connect(ctx); err != nil {
		return err
	}
	defer c.Close()

	binPath, err := c.UploadBinary(ctx, chainBin)
	if err != nil {
		return err
	}

	fmt.Println(binPath)
	return nil
}
