package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/pkg/errors"

	"github.com/ignite/apps/explorer/gex"
)

const maxNumArgs = 1

// ExecuteGex executes explorer gex subcommand.
func ExecuteGex(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	argc := len(cmd.Args)
	if argc > maxNumArgs {
		return fmt.Errorf("accepts at most %d arg(s), received %d", maxNumArgs, argc)
	}

	ssl := false
	host := "localhost"
	port := "26657"

	if argc == 1 {
		rpcURL, err := url.Parse(cmd.Args[0])
		if err != nil {
			return errors.Wrapf(err, "failed to parse RPC URL %s", cmd.Args[0])
		}

		ssl = rpcURL.Scheme == "https"
		host = rpcURL.Hostname()
		port = rpcURL.Port()
		if port == "" {
			if ssl {
				port = "443"
			} else {
				port = "80"
			}
		}
	}

	g, err := gex.New()
	if err != nil {
		return errors.Wrap(err, "failed to initialize Gex")
	}
	return g.Run(ctx, os.Stdout, os.Stderr, host, port, ssl)
}
