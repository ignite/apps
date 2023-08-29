package cmd

import (
	"github.com/spf13/cobra"
)

// NewRelayerConfigure configure the hermes relayer and create the config file.
func NewRelayerConfigure() *cobra.Command {
	c := &cobra.Command{
		Use:   "configure [launch-id]",
		Short: "",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE:  relayerConfigureHandler,
	}

	return c
}

func relayerConfigureHandler(cmd *cobra.Command, args []string) error {
	return nil
}
