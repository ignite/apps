package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"relayer/pkg/hermes"
)

// NewHermesStart start the hermes relayer.
func NewHermesStart() *cobra.Command {
	c := &cobra.Command{
		Use:   "start [chain-a-id] [chain-a-rpc]",
		Short: "",
		Long:  ``,
		Args:  cobra.ExactArgs(2),
		RunE:  hermesStartHandler,
	}

	return c
}

func hermesStartHandler(cmd *cobra.Command, args []string) error {
	cfgName := strings.Join(args, hermes.ConfigNameSeparator)
	cfg, err := hermes.LoadConfig(cfgName)
	if err != nil {
		return err
	}

	cfgPath, err := cfg.ConfigPath()
	if err != nil {
		return err
	}

	h, err := hermes.New()
	if err != nil {
		return err
	}
	defer h.Cleanup()

	return h.Start(
		cmd.Context(),
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(os.Stdout),
	)
}
