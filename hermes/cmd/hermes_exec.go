package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"relayer/pkg/hermes"
)

// NewHermesExecute execute hermes relayer commands.
func NewHermesExecute() *cobra.Command {
	c := &cobra.Command{
		Use:   "exec [args...]",
		Short: "",
		Long:  ``,
		Args:  cobra.MinimumNArgs(1),
		RunE:  hermesExecuteHandler,
	}

	return c
}

func hermesExecuteHandler(cmd *cobra.Command, args []string) error {
	//WithEventSource(mode, url, batchDelay string)
	//WithRPCTimeout(timeout string)
	//WithAccountPrefix(prefix string)
	//WithKeyName(key string)
	//WithStorePrefix(prefix string)
	//WithDefaultGas(defaultGas int)
	//WithMaxGas(maxGas int)
	//WithGasPrice(price float64, denom string)
	//WithGasMultiplier(gasMultipler float64)
	//WithMaxMsgNum(maxMsg int)
	//WithMaxTxSize(size int)
	//WithClockDrift(clock string)
	//WithMaxBlockTime(maxBlockTime string)
	//WithTrustingPeriod(trustingPeriod string)
	//WithTrustThreshold(numerator, denominator string)
	//WithAddressPrefix(derivation string)
	
	h, err := hermes.New()
	if err != nil {
		return err
	}
	defer h.Cleanup()
	return h.Run(cmd.Context(), os.Stdout, os.Stderr, "", args...)
}
