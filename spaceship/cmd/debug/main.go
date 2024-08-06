package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/spaceship/cmd"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	var (
		args      = os.Args
		ctx       = context.Background()
		chainInfo = &plugin.ChainInfo{
			AppPath: filepath.Join(home, "Desktop/go/src/github.com/ignite/mars"),
			ChainId: "mars",
		}
		c = &plugin.ExecutedCommand{
			Use:    args[1],
			Path:   "ignite spaceship " + args[1],
			Args:   []string{fmt.Sprintf("%s@127.0.0.1", filepath.Base(home))},
			OsArgs: os.Args,
			With:   nil,
			Flags: []*plugin.Flag{
				{
					Name:      "key",
					Shorthand: "k",
					Usage:     "ssh key",
					Type:      plugin.FlagTypeString,
					Value:     filepath.Join(home, ".ssh/id_rsa"),
				},
				{
					Name:      "init-chain",
					Shorthand: "i",
					Usage:     "run init chain and create the home folder",
					Type:      plugin.FlagTypeBool,
					Value:     "true",
				},
			},
		}
	)
	switch args[1] {
	case "deploy":
		if err := cmd.ExecuteSSHDeploy(ctx, c, chainInfo); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	case "log":
		if err := cmd.ExecuteSSHLog(ctx, c, chainInfo); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	case "status":
		if err := cmd.ExecuteSSHStatus(ctx, c, chainInfo); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	case "restart":
		if err := cmd.ExecuteSSHRestart(ctx, c, chainInfo); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	case "stop":
		if err := cmd.ExecuteSSHSStop(ctx, c, chainInfo); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s", args[1])
		return
	}
}
