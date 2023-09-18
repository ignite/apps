package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ignite/apps/hermes/pkg/hermes"
)

// NewHermesStart start the hermes relayer.
func NewHermesStart() *cobra.Command {
	c := &cobra.Command{
		Use:   "start [chain-a-id] [chain-a-rpc]",
		Short: "Start the Hermes realyer",
		Args:  cobra.ExactArgs(2),
		RunE:  hermesStartHandler,
	}

	return c
}

func hermesStartHandler(cmd *cobra.Command, args []string) (err error) {
	var (
		customCfg = getConfig(cmd)
		cfgName   = strings.Join(args, hermes.ConfigNameSeparator)
	)

	cfgPath := customCfg
	if cfgPath == "" {
		cfgPath, err = hermes.ConfigPath(cfgName)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return fmt.Errorf("config file (%s) not exist, try to configure you relayer first", cfgPath)
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
		hermes.WithStdErr(os.Stderr),
	)
}
