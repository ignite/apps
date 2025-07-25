package cmd

import (
	"fmt"
	"time"

	"github.com/ignite/cli/v29/ignite/pkg/ctxticker"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/ignite/apps/network/network/networkchain"
	"github.com/ignite/apps/network/network/xchisel"
)

func NewNetworkTool() *cobra.Command {
	c := &cobra.Command{
		Use:   "tool [command]",
		Short: "Commands to run subsidiary tools",
	}
	c.AddCommand(NewNetworkToolProxyTunnel())
	return c
}

func NewNetworkToolProxyTunnel() *cobra.Command {
	c := &cobra.Command{
		Use:   "proxy-tunnel [config-file]",
		Short: "Setup a proxy tunnel via HTTP",
		Long: `Starts an HTTP proxy server and HTTP proxy clients for each node that
needs HTTP tunneling.

HTTP tunneling is activated **ONLY** if SPN_CONFIG_FILE has "tunneled_peers"
field inside with a list of tunneled peers/nodes.

If you're using SPN as coordinator and do not want to allow HTTP tunneling
feature at all, you can prevent "spn.yml" file to being generated by not
approving validator requests that has HTTP tunneling enabled instead of plain
TCP connections.`,
		Args: cobra.ExactArgs(1),
		RunE: networkToolProxyTunnelHandler,
	}
	return c
}

const tunnelRerunDelay = 5 * time.Second

func networkToolProxyTunnelHandler(cmd *cobra.Command, args []string) error {
	spnConfig, err := networkchain.GetSPNConfig(args[0])
	if err != nil {
		return fmt.Errorf("failed to open spn config file: %w", err)
	}
	// exit if there aren't tunneled validators in the network
	if len(spnConfig.TunneledPeers) == 0 {
		return nil
	}

	g, ctx := errgroup.WithContext(cmd.Context())
	for _, peer := range spnConfig.TunneledPeers {
		if peer.Name == networkchain.HTTPTunnelChisel {
			peer := peer
			g.Go(func() error {
				return ctxticker.DoNow(ctx, tunnelRerunDelay, func() error {
					fmt.Printf("Starting chisel client, tunnelAddress:%s, localPort:%s\n", peer.Address, peer.LocalPort)
					err := xchisel.StartClient(ctx, peer.Address, peer.LocalPort, "26656")
					if err != nil {
						fmt.Printf(
							"Failed to start chisel client, tunnelAddress:%s, localPort:%s, reason:%v\n",
							peer.Address, peer.LocalPort, err,
						)
					}
					return nil
				})
			})
		}
	}

	g.Go(func() error {
		return ctxticker.DoNow(ctx, tunnelRerunDelay, func() error {
			fmt.Printf("Starting chisel server, port:%s\n", xchisel.DefaultServerPort)
			err := xchisel.StartServer(ctx, xchisel.DefaultServerPort)
			if err != nil {
				fmt.Printf(
					"Failed to start chisel server, port:%s, reason:%v\n",
					xchisel.DefaultServerPort, err,
				)
			}
			return nil
		})
	})
	return g.Wait()
}
