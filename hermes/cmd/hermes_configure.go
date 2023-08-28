package cmd

import (
	"github.com/spf13/cobra"
)

// NewHermesConfigure configure the hermes relayer and create the config file.
func NewHermesConfigure() *cobra.Command {
	c := &cobra.Command{
		Use:   "configure [launch-id]",
		Short: "",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE:  hermesConfigureHandler,
	}

	return c
}

func hermesConfigureHandler(cmd *cobra.Command, args []string) error {
	return nil
}
