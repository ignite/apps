package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"relayer/pkg/hermes"
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
	h, err := hermes.New()
	if err != nil {
		return err
	}
	defer h.Cleanup()

	// Create the default config and add chains
	c := hermes.DefaultConfig()
	cfgPath, err := c.ConfigPath()
	if err != nil {
		return err
	}

	return h.Run(cmd.Context(), os.Stdout, os.Stderr, cfgPath, args...)
}
