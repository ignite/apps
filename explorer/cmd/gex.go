package cmd

import (
	"explorer/pkg/gex"
	"net/url"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	defaultHost = "localhost"
	defaultPort = "26657"
)

func NewGex() *cobra.Command {
	c := &cobra.Command{
		Use:     "gex [rpc url]",
		Aliases: []string{"g"},
		Short:   "Run gex",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			host := defaultHost
			port := defaultPort
			ssl := false

			if len(args) == 1 {
				rpcURL, err := url.Parse(args[0])
				if err != nil {
					return errors.Wrapf(err, "failed to parse rpc url %s", args[0])
				}

				host = rpcURL.Hostname()
				port = rpcURL.Port()
				ssl = false
				if rpcURL.Scheme == "https" {
					ssl = true
				}
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
				errors.Wrap(err, "failed to initialize gex")
			}

			return g.Run(cmd.Context(), os.Stdout, os.Stderr, host, port, ssl)
		},
	}

	return c
}
