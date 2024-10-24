package cmd

import (
	"path/filepath"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v28/ignite/pkg/goenv"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/network/network"
	"github.com/ignite/apps/network/network/networkchain"
)

// NewNetworkChainInstall returns a new command to install a chain's binary by the launch id.
func NewNetworkChainInstall() *cobra.Command {
	c := &cobra.Command{
		Use:   "install [launch-id]",
		Short: "Install chain binary for a launch",
		Args:  cobra.ExactArgs(1),
		RunE:  networkChainInstallHandler,
	}

	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetCheckDependencies())
	return c
}

func networkChainInstallHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse launch ID
	launchID, err := network.ParseID(args[0])
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

	var networkOptions []networkchain.Option

	if flagGetCheckDependencies(cmd) {
		networkOptions = append(networkOptions, networkchain.CheckDependencies())
	}

	c, err := nb.Chain(networkchain.SourceLaunch(chainLaunch), networkOptions...)
	if err != nil {
		return err
	}

	binaryName, err := c.Build(cmd.Context(), cacheStorage)
	if err != nil {
		return err
	}
	binaryPath := filepath.Join(goenv.Bin(), binaryName)

	session.Printf("%s Binary installed\n", icons.OK)
	session.Printf("%s Binary's name: %s\n", icons.Info, colors.Info(binaryName))
	session.Printf("%s Binary's path: %s\n", icons.Info, colors.Info(binaryPath))

	return nil
}
