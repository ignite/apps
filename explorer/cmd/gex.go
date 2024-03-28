package cmd

import (
	"context"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/explorer/gex"
	"github.com/ignite/apps/explorer/pkg/xurl"
)

// ExecuteGex executes explorer gex subcommand.
func ExecuteGex(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	flags, err := cmd.NewFlags()
	if err != nil {
		return err
	}

	rpcAddress, _ := flags.GetString(flagRPCAddress)
	if err != nil {
		return errors.Errorf("could not get --%s flag: %s", flagRPCAddress, err)
	}

	rpcURL, err := xurl.Parse(rpcAddress)
	if err != nil {
		return errors.Wrapf(err, "failed to parse RPC URL %s", rpcAddress)
	}

	g, err := gex.New(
		gex.WithHost(rpcURL.Hostname()),
		gex.WithPort(rpcURL.Port()),
		gex.WithSSL(xurl.IsSSL(rpcURL)),
	)
	if err != nil {
		return errors.Wrap(err, "failed to initialize Gex")
	}
	defer g.Cleanup()

	return g.Run(ctx)
}
