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
	)
	switch args[1] {
	case "aws":
		if err := cmd.ExecuteAWS(ctx, nil); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	case "ssh":
		switch args[2] {
		case "deploy":
			if err := cmd.ExecuteSSHDeploy(ctx, chainInfo); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		default:
			fmt.Fprintf(os.Stderr, "unknown ssh command: %s", args[2])
			return
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s", args[1])
		return
	}
}
