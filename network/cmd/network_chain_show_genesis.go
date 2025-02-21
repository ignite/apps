package cmd

import (
	"os"
	"path/filepath"

	"github.com/ignite/cli/v28/ignite/pkg/chaincmd"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v28/ignite/pkg/xos"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/network/network/networkchain"
	"github.com/ignite/apps/network/network/networktypes"
)

func newNetworkChainShowGenesis() *cobra.Command {
	c := &cobra.Command{
		Use:   "genesis [launch-id]",
		Short: "Show the chain genesis file",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainShowGenesisHandler,
	}

	flagSetClearCache(c)
	c.Flags().String(flagOut, "./genesis.json", "path to output Genesis file")

	return c
}

func networkChainShowGenesisHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	out, _ := cmd.Flags().GetString(flagOut)

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	nb, launchID, err := networkChainLaunch(cmd, args, session)
	if err != nil {
		return err
	}
	n, err := nb.Network()
	if err != nil {
		return err
	}

	chainLaunch, err := n.ChainLaunch(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	networkOptions := []networkchain.Option{
		networkchain.WithKeyringBackend(chaincmd.KeyringBackendTest),
	}

	c, err := nb.Chain(networkchain.SourceLaunch(chainLaunch), networkOptions...)
	if err != nil {
		return err
	}

	// generate the genesis in a temp dir
	tmpHome, err := os.MkdirTemp("", "*-spn")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpHome)

	c.SetHome(tmpHome)

	if err := prepareFromGenesisInformation(
		cmd,
		cacheStorage,
		launchID,
		n,
		c,
		chainLaunch,
	); err != nil {
		return err
	}

	// get the new genesis path
	genesisPath, err := c.GenesisPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(out), 0o744); err != nil {
		return err
	}

	if err := xos.Rename(genesisPath, out); err != nil {
		return err
	}

	if chainLaunch.Metadata.Cli.Version != "" && !chainLaunch.Metadata.IsCurrentVersion() {
		session.Printf(`⚠️ chain %d has been published with a different version of the plugin (%s, current version is %s)
this may result in a genesis that is different from other validators' genesis
for chain launch, please update the plugin to the same version\n`,
			launchID,
			chainLaunch.Metadata.Cli.Version,
			networktypes.Version,
		)
	}

	return session.Printf("%s Genesis generated: %s\n", icons.Bullet, out)
}
