package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v29/ignite/services/plugin"

	"github.com/ignite/apps/spaceship/cmd"
	"github.com/ignite/apps/spaceship/pkg/ssh"
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
			AppPath:    filepath.Join(home, "Desktop/go/src/github.com/ignite/mars"),
			ConfigPath: filepath.Join(home, "Desktop/go/src/github.com/ignite/mars/config.yml"),
			ChainId:    "mars",
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
			},
		}
	)
	switch args[1] {
	case "deploy":
		c.Flags = append(c.Flags,
			&plugin.Flag{
				Name:      "init-chain",
				Shorthand: "i",
				Usage:     "run init chain and create the home folder",
				Type:      plugin.FlagTypeBool,
				Value:     "true",
			},
			&plugin.Flag{
				Name:      "init-chain",
				Shorthand: "i",
				Usage:     "run init chain and create the home folder",
				Type:      plugin.FlagTypeBool,
				Value:     "true",
			},
			&plugin.Flag{
				Name:         "faucet",
				Shorthand:    "f",
				Usage:        "create a chain faucet",
				Type:         plugin.FlagTypeBool,
				DefaultValue: "true",
			},
		)
		if err := cmd.ExecuteSSHDeploy(ctx, c, chainInfo); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	case "log":
		c.Flags = append(c.Flags,
			&plugin.Flag{
				Name:  "real-time",
				Usage: "show the logs in the real time",
				Type:  plugin.FlagTypeBool,
				Value: "true",
			},
			&plugin.Flag{
				Name:      "app",
				Shorthand: "a",
				Usage: fmt.Sprintf(
					"the app to show the log (%s)",
					strings.Join(ssh.LogTypes(), ","),
				),
				Type:         plugin.FlagTypeString,
				DefaultValue: ssh.LogChain.String(),
				Value:        ssh.LogChain.String(),
			},
		)
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
	case "faucet":
		switch args[2] {
		case "status":
			if err := cmd.ExecuteSSHFaucetStatus(ctx, c, chainInfo); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		case "start":
			if err := cmd.ExecuteSSHFaucetStart(ctx, c, chainInfo); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		case "restart":
			if err := cmd.ExecuteSSHFaucetRestart(ctx, c, chainInfo); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		case "stop":
			if err := cmd.ExecuteSSHSFaucetStop(ctx, c, chainInfo); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}

		default:
			fmt.Fprintf(os.Stderr, "unknown faucet command: %s", args[2])
			return
		}
	}
}
