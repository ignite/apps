package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/icons"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/network/network"
)

func newNetworkChainShowPeers() *cobra.Command {
	c := &cobra.Command{
		Use:   "peers [launch-id]",
		Short: "Show peers list of the chain",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainShowPeersHandler,
	}

	c.Flags().String(flagOut, "./peers.txt", "path to output peers list")

	return c
}

func networkChainShowPeersHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	out, _ := cmd.Flags().GetString(flagOut)

	nb, launchID, err := networkChainLaunch(cmd, args, session)
	if err != nil {
		return err
	}
	n, err := nb.Network()
	if err != nil {
		return err
	}

	genVals, err := n.GenesisValidators(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	peers := make([]string, 0)
	for _, acc := range genVals {
		peer, err := network.PeerAddress(acc.Peer)
		if err != nil {
			return err
		}
		peers = append(peers, peer)
	}

	if len(peers) == 0 {
		session.Printf("%s %s\n", icons.Info, "no peers found")
		return nil

	}

	if err := os.MkdirAll(filepath.Dir(out), 0o744); err != nil {
		return err
	}

	b := &bytes.Buffer{}
	peerList := strings.Join(peers, ",")
	fmt.Fprintln(b, peerList)
	if err := os.WriteFile(out, b.Bytes(), 0o644); err != nil {
		return err
	}

	return session.Printf("%s Peer list generated: %s\n", icons.Bullet, out)
}
