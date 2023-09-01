package cmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"os"

	"github.com/spf13/cobra"

	"relayer/pkg/hermes"
)

const (
	flagChainAEventSourceMode           = "chain-a-event-source-mode"
	flagChainAEventSourceUrl            = "chain-a-event-source-url"
	flagChainAEventSourceBatchDelay     = "chain-a-event-source-batch-delay"
	flagChainARPCTimeout                = "chain-a-rpc-timeout"
	flagChainAAccountPrefix             = "chain-a-account-prefix"
	flagChainAKeyName                   = "chain-a-key-name"
	flagChainAStorePrefix               = "chain-a-store-prefix"
	flagChainADefaultGas                = "chain-a-default-gas"
	flagChainAMaxGas                    = "chain-a-max-gas"
	flagChainAGasPrice                  = "chain-a-gas-price"
	flagChainAGasMultiplier             = "chain-a-gas-multiplier"
	flagChainAMaxMsgNum                 = "chain-a-max-msg-num"
	flagChainAMaxTxSize                 = "chain-a-tx-size"
	flagChainAClockDrift                = "chain-a-clock-drift"
	flagChainAMaxBlockTime              = "chain-a-max-block-time"
	flagChainATrustingPeriod            = "chain-a-trusting-period"
	flagChainATrustThresholdNumerator   = "chain-a-trust-threshold-numerator"
	flagChainATrustThresholdDenominator = "chain-a-trust-threshold-denominator"

	flagChainBEventSourceMode           = "chain-b-event-source-mode"
	flagChainBEventSourceUrl            = "chain-b-event-source-url"
	flagChainBEventSourceBatchDelay     = "chain-b-event-source-batch-delay"
	flagChainBRPCTimeout                = "chain-b-rpc-timeout"
	flagChainBAccountPrefix             = "chain-b-account-prefix"
	flagChainBKeyName                   = "chain-b-key-name"
	flagChainBStorePrefix               = "chain-b-store-prefix"
	flagChainBDefaultGas                = "chain-b-default-gas"
	flagChainBMaxGas                    = "chain-b-max-gas"
	flagChainBGasPrice                  = "chain-b-gas-price"
	flagChainBGasMultiplier             = "chain-b-gas-multiplier"
	flagChainBMaxMsgNum                 = "chain-b-max-msg-num"
	flagChainBMaxTxSize                 = "chain-b-tx-size"
	flagChainBClockDrift                = "chain-b-clock-drift"
	flagChainBMaxBlockTime              = "chain-b-max-block-time"
	flagChainBTrustingPeriod            = "chain-b-trusting-period"
	flagChainBTrustThresholdNumerator   = "chain-b-trust-threshold-numerator"
	flagChainBTrustThresholdDenominator = "chain-b-trust-threshold-denominator"

	flagTelemetryEnabled          = "telemetry_enabled"
	flagTelemetryHost             = "telemetry_host"
	flagTelemetryPort             = "telemetry_port"
	flagModeChannelsEnabled       = "mode_channels_enabled"
	flagModeClientsEnabled        = "mode_clients_enabled"
	flagModeClientsMisbehaviour   = "mode_clients_misbehaviour"
	flagModeClientsRefresh        = "mode_clients_refresh"
	flagModeConnectionsEnabled    = "mode_connections_enabled"
	flagModePacketsEnabled        = "mode_packets_enabled"
	flagModePacketsClearInterval  = "mode_packets_clear_interval"
	flagModePacketsClearOnStart   = "mode_packets_clear_on_start"
	flagModePacketsTxConfirmation = "mode_packets_tx_confirmation"
)

// NewHermesConfigure configure the hermes relayer and create the config file.
func NewHermesConfigure() *cobra.Command {
	c := &cobra.Command{
		Use:   "configure [chain-a-id] [chain-a-rpc] [chain-a-grpc] [chain-b-id] [chain-b-rpc] [chain-b-grpc]",
		Short: "",
		Long:  ``,
		Args:  cobra.ExactArgs(6),
		RunE:  hermesConfigureHandler,
	}

	c.Flags().String(flagChainAEventSourceMode, "push", "WS event source mode of the chain A")
	c.Flags().String(flagChainBEventSourceMode, "push", "WS event source mode of the chain B")
	c.Flags().String(flagChainAEventSourceUrl, "", "WS event source url of the chain A")
	c.Flags().String(flagChainBEventSourceUrl, "", "WS event source url of the chain B")
	c.Flags().String(flagChainAEventSourceBatchDelay, "500ms", "WS event source batch delay time of the chain A")
	c.Flags().String(flagChainBEventSourceBatchDelay, "500ms", "WS event source batch delay time of the chain B")
	c.Flags().String(flagChainARPCTimeout, "", "RPC timeout of the chain A")
	c.Flags().String(flagChainBRPCTimeout, "", "RPC timeout of the chain B")
	c.Flags().String(flagChainAAccountPrefix, "", "address prefix of the chain A")
	c.Flags().String(flagChainBAccountPrefix, "", "address prefix of the chain B")
	c.Flags().String(flagChainAKeyName, "", "hermes account name of the chain A")
	c.Flags().String(flagChainBKeyName, "", "hermes account name of the chain B")
	c.Flags().String(flagChainAStorePrefix, "", "store prefix of the chain A")
	c.Flags().String(flagChainBStorePrefix, "", "store prefix of the chain B")
	c.Flags().Uint64(flagChainADefaultGas, 0, "default gas used for transactions on chain A")
	c.Flags().Uint64(flagChainBDefaultGas, 0, "default gas used for transactions on chain B")
	c.Flags().Uint64(flagChainAMaxGas, 0, "max gas used for transactions on chain A")
	c.Flags().Uint64(flagChainBMaxGas, 0, "max gas used for transactions on chain B")
	c.Flags().String(flagChainAGasPrice, "", "gas price used for transactions on chain A")
	c.Flags().String(flagChainBGasPrice, "", "gas price used for transactions on chain B")
	c.Flags().Float64(flagChainAGasMultiplier, 0, "gas multiplier used for transactions on chain A")
	c.Flags().Float64(flagChainBGasMultiplier, 0, "gas multiplier used for transactions on chain B")
	c.Flags().Uint64(flagChainAMaxMsgNum, 0, "max message number used for transactions on chain A")
	c.Flags().Uint64(flagChainBMaxMsgNum, 0, "max message number used for transactions on chain B")
	c.Flags().Uint64(flagChainAMaxTxSize, 0, "max transaction size on chain A")
	c.Flags().Uint64(flagChainBMaxTxSize, 0, "max transaction size on chain B")
	c.Flags().String(flagChainAClockDrift, "", "clock drift of the chain A")
	c.Flags().String(flagChainBClockDrift, "", "clock drift of the chain B")
	c.Flags().String(flagChainAMaxBlockTime, "", "maximum block time of the chain A")
	c.Flags().String(flagChainBMaxBlockTime, "", "maximum block time of the chain B")
	c.Flags().String(flagChainATrustingPeriod, "", "trusting period of the chain A")
	c.Flags().String(flagChainBTrustingPeriod, "", "trusting period of the chain B")
	c.Flags().Uint64(flagChainATrustThresholdNumerator, 1, "trusting threshold numerator of the chain A")
	c.Flags().Uint64(flagChainBTrustThresholdNumerator, 1, "trusting threshold numerator of the chain B")
	c.Flags().Uint64(flagChainATrustThresholdDenominator, 3, "trusting threshold denominator of the chain A")
	c.Flags().Uint64(flagChainBTrustThresholdDenominator, 3, "trusting threshold denominator of the chain B")

	c.Flags().String(flagTelemetryEnabled, "", "enable hermes telemetry")
	c.Flags().String(flagTelemetryHost, "", "hermes telemetry host")
	c.Flags().String(flagTelemetryPort, "", "hermes telemetry port")
	c.Flags().String(flagModeChannelsEnabled, "", "enable hermes channels")
	c.Flags().String(flagModeClientsEnabled, "", "enable hermes clients")
	c.Flags().String(flagModeClientsMisbehaviour, "", "enable hermes clients misbehaviour")
	c.Flags().String(flagModeClientsRefresh, "", "hermes client refresh time")
	c.Flags().String(flagModeConnectionsEnabled, "", "enable hermes connections")
	c.Flags().String(flagModePacketsEnabled, "", "enable hermes packets")
	c.Flags().Uint64(flagModePacketsClearInterval, 0, "hermes packet clear interval")
	c.Flags().String(flagModePacketsClearOnStart, "", "enable hermes packets clear on start")
	c.Flags().String(flagModePacketsTxConfirmation, "", "hermes packet transaction confirmation")

	return c
}

func hermesConfigureHandler(cmd *cobra.Command, args []string) error {
	// Create the default config and add chains
	var (
		c = hermes.DefaultConfig()

		chainAID       = args[0]
		chainARPCAddr  = args[1]
		chainAGRPCAddr = args[2]

		chainAEventSourceMode, _           = cmd.Flags().GetString(flagChainAEventSourceMode)
		chainAEventSourceUrl, _            = cmd.Flags().GetString(flagChainAEventSourceUrl)
		chainAEventSourceBatchDelay, _     = cmd.Flags().GetString(flagChainAEventSourceBatchDelay)
		chainARPCTimeout, _                = cmd.Flags().GetString(flagChainARPCTimeout)
		chainAAccountPrefix, _             = cmd.Flags().GetString(flagChainAAccountPrefix)
		chainAKeyName, _                   = cmd.Flags().GetString(flagChainAKeyName)
		chainAStorePrefix, _               = cmd.Flags().GetString(flagChainAStorePrefix)
		chainADefaultGas, _                = cmd.Flags().GetUint64(flagChainADefaultGas)
		chainAMaxGas, _                    = cmd.Flags().GetUint64(flagChainAMaxGas)
		chainAGasPrice, _                  = cmd.Flags().GetString(flagChainAGasPrice)
		chainAGasMultiplier, _             = cmd.Flags().GetFloat64(flagChainAGasMultiplier)
		chainAMaxMsgNum, _                 = cmd.Flags().GetUint64(flagChainAMaxMsgNum)
		chainAMaxTxSize, _                 = cmd.Flags().GetUint64(flagChainAMaxTxSize)
		chainAClockDrift, _                = cmd.Flags().GetString(flagChainAClockDrift)
		chainAMaxBlockTime, _              = cmd.Flags().GetString(flagChainAMaxBlockTime)
		chainATrustingPeriod, _            = cmd.Flags().GetString(flagChainATrustingPeriod)
		chainATrustThresholdNumerator, _   = cmd.Flags().GetUint64(flagChainATrustThresholdNumerator)
		chainATrustThresholdDenominator, _ = cmd.Flags().GetUint64(flagChainATrustThresholdDenominator)
	)

	optChainA := []hermes.ChainOption{
		hermes.WithChainTrustThreshold(chainATrustThresholdNumerator, chainATrustThresholdDenominator),
	}
	if chainAEventSourceMode != "" {
		optChainA = append(optChainA, hermes.WithChainEventSource(
			chainAEventSourceMode,
			chainAEventSourceUrl,
			chainAEventSourceBatchDelay,
		))
	}
	if chainARPCTimeout != "" {
		optChainA = append(optChainA, hermes.WithChainRPCTimeout(chainARPCTimeout))
	}
	if chainAAccountPrefix != "" {
		optChainA = append(optChainA, hermes.WithChainAddressPrefix(chainAAccountPrefix))
	}
	if chainAKeyName != "" {
		optChainA = append(optChainA, hermes.WithChainKeyName(chainAKeyName))
	}
	if chainAStorePrefix != "" {
		optChainA = append(optChainA, hermes.WithChainStorePrefix(chainAStorePrefix))
	}
	if chainADefaultGas > 0 {
		optChainA = append(optChainA, hermes.WithChainDefaultGas(chainADefaultGas))
	}
	if chainAMaxGas > 0 {
		optChainA = append(optChainA, hermes.WithChainMaxGas(chainAMaxGas))
	}
	if chainAGasPrice != "" {
		gasPrice, err := sdk.ParseCoinNormalized(chainAGasPrice)
		if err != nil {
			return err
		}
		optChainA = append(optChainA, hermes.WithChainGasPrice(gasPrice))
	}
	if chainAGasMultiplier > 0 {
		optChainA = append(optChainA, hermes.WithChainGasMultiplier(chainAGasMultiplier))
	}
	if chainAMaxMsgNum > 0 {
		optChainA = append(optChainA, hermes.WithChainMaxMsgNum(chainAMaxMsgNum))
	}
	if chainAMaxTxSize > 0 {
		optChainA = append(optChainA, hermes.WithChainMaxTxSize(chainAMaxTxSize))
	}
	if chainAClockDrift != "" {
		optChainA = append(optChainA, hermes.WithChainClockDrift(chainAClockDrift))
	}
	if chainAMaxBlockTime != "" {
		optChainA = append(optChainA, hermes.WithChainMaxBlockTime(chainAMaxBlockTime))
	}
	if chainATrustingPeriod != "" {
		optChainA = append(optChainA, hermes.WithChainTrustingPeriod(chainATrustingPeriod))
	}

	err := c.AddChain(chainAID, chainARPCAddr, chainAGRPCAddr, optChainA...)
	if err != nil {
		return err
	}

	var (
		chainBID       = args[3]
		chainBRPCAddr  = args[4]
		chainBGRPCAddr = args[5]

		chainBEventSourceMode, _           = cmd.Flags().GetString(flagChainBEventSourceMode)
		chainBEventSourceUrl, _            = cmd.Flags().GetString(flagChainBEventSourceUrl)
		chainBEventSourceBatchDelay, _     = cmd.Flags().GetString(flagChainBEventSourceBatchDelay)
		chainBRPCTimeout, _                = cmd.Flags().GetString(flagChainBRPCTimeout)
		chainBAccountPrefix, _             = cmd.Flags().GetString(flagChainBAccountPrefix)
		chainBKeyName, _                   = cmd.Flags().GetString(flagChainBKeyName)
		chainBStorePrefix, _               = cmd.Flags().GetString(flagChainBStorePrefix)
		chainBDefaultGas, _                = cmd.Flags().GetUint64(flagChainBDefaultGas)
		chainBMaxGas, _                    = cmd.Flags().GetUint64(flagChainBMaxGas)
		chainBGasPrice, _                  = cmd.Flags().GetString(flagChainBGasPrice)
		chainBGasMultiplier, _             = cmd.Flags().GetFloat64(flagChainBGasMultiplier)
		chainBMaxMsgNum, _                 = cmd.Flags().GetUint64(flagChainBMaxMsgNum)
		chainBMaxTxSize, _                 = cmd.Flags().GetUint64(flagChainBMaxTxSize)
		chainBClockDrift, _                = cmd.Flags().GetString(flagChainBClockDrift)
		chainBMaxBlockTime, _              = cmd.Flags().GetString(flagChainBMaxBlockTime)
		chainBTrustingPeriod, _            = cmd.Flags().GetString(flagChainBTrustingPeriod)
		chainBTrustThresholdNumerator, _   = cmd.Flags().GetUint64(flagChainBTrustThresholdNumerator)
		chainBTrustThresholdDenominator, _ = cmd.Flags().GetUint64(flagChainBTrustThresholdDenominator)
	)

	optChainB := []hermes.ChainOption{
		hermes.WithChainTrustThreshold(chainBTrustThresholdNumerator, chainBTrustThresholdDenominator),
	}
	if chainBEventSourceMode != "" {
		optChainB = append(optChainB, hermes.WithChainEventSource(
			chainBEventSourceMode,
			chainBEventSourceUrl,
			chainBEventSourceBatchDelay,
		))
	}
	if chainBRPCTimeout != "" {
		optChainB = append(optChainB, hermes.WithChainRPCTimeout(chainBRPCTimeout))
	}
	if chainBAccountPrefix != "" {
		optChainB = append(optChainB, hermes.WithChainAddressPrefix(chainBAccountPrefix))
	}
	if chainBKeyName != "" {
		optChainB = append(optChainB, hermes.WithChainKeyName(chainBKeyName))
	}
	if chainBStorePrefix != "" {
		optChainB = append(optChainB, hermes.WithChainStorePrefix(chainBStorePrefix))
	}
	if chainBDefaultGas > 0 {
		optChainB = append(optChainB, hermes.WithChainDefaultGas(chainBDefaultGas))
	}
	if chainBMaxGas > 0 {
		optChainB = append(optChainB, hermes.WithChainMaxGas(chainBMaxGas))
	}
	if chainBGasPrice != "" {
		gasPrice, err := sdk.ParseCoinNormalized(chainBGasPrice)
		if err != nil {
			return err
		}
		optChainB = append(optChainB, hermes.WithChainGasPrice(gasPrice))
	}
	if chainBGasMultiplier > 0 {
		optChainB = append(optChainB, hermes.WithChainGasMultiplier(chainBGasMultiplier))
	}
	if chainBMaxMsgNum > 0 {
		optChainB = append(optChainB, hermes.WithChainMaxMsgNum(chainBMaxMsgNum))
	}
	if chainBMaxTxSize > 0 {
		optChainB = append(optChainB, hermes.WithChainMaxTxSize(chainBMaxTxSize))
	}
	if chainBClockDrift != "" {
		optChainB = append(optChainB, hermes.WithChainClockDrift(chainBClockDrift))
	}
	if chainBMaxBlockTime != "" {
		optChainB = append(optChainB, hermes.WithChainMaxBlockTime(chainBMaxBlockTime))
	}
	if chainBTrustingPeriod != "" {
		optChainB = append(optChainB, hermes.WithChainTrustingPeriod(chainBTrustingPeriod))
	}

	err = c.AddChain(chainBID, chainBRPCAddr, chainBGRPCAddr, optChainB...)
	if err != nil {
		return err
	}

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

	// Create the default config and add chains
	cfgPath, err := c.ConfigPath()
	if err != nil {
		return err
	}

	return h.Run(cmd.Context(), os.Stdout, os.Stderr, cfgPath, args...)
}
