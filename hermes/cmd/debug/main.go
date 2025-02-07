package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/hermes/cmd"
)

const (
	flagChainAPortID                    = "chain-a-port-id"
	flagChainAEventSourceMode           = "chain-a-event-source-mode"
	flagChainAEventSourceBatchDelay     = "chain-a-event-source-batch-delay"
	flagChainARPCTimeout                = "chain-a-rpc-timeout"
	flagChainAAccountPrefix             = "chain-a-account-prefix"
	flagChainAAddressType               = "chain-a-address-types"
	flagChainAKeyName                   = "chain-a-key-name"
	flagChainAKeyStoreType              = "chain-a-key-store-type"
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
	flagChainAFaucet                    = "chain-a-faucet"
	flagChainACCVConsumerChain          = "chain-a-ccv-consumer-chain"
	flagChainATrustedNode               = "chain-a-trusted-node"
	flagChainAType                      = "chain-a-type"
	flagChainASequentialBatchTx         = "chain-a-sequential-batch-tx"

	flagChainBPortID                    = "chain-b-port-id"
	flagChainBEventSourceMode           = "chain-b-event-source-mode"
	flagChainBEventSourceBatchDelay     = "chain-b-event-source-batch-delay"
	flagChainBRPCTimeout                = "chain-b-rpc-timeout"
	flagChainBAccountPrefix             = "chain-b-account-prefix"
	flagChainBAddressType               = "chain-b-address-types"
	flagChainBKeyName                   = "chain-b-key-name"
	flagChainBKeyStoreType              = "chain-b-key-store-type"
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
	flagChainBFaucet                    = "chain-b-faucet"
	flagChainBCCVConsumerChain          = "chain-b-ccv-consumer-chain"
	flagChainBTrustedNode               = "chain-b-trusted-node"
	flagChainBType                      = "chain-b-type"
	flagChainBSequentialBatchTx         = "chain-b-sequential-batch-tx"

	flagTelemetryEnabled              = "telemetry-enabled"
	flagTelemetryHost                 = "telemetry-host"
	flagTelemetryPort                 = "telemetry-port"
	flagRestEnabled                   = "rest-enabled"
	flagRestHost                      = "rest-host"
	flagRestPort                      = "rest-port"
	flagModeChannelsEnabled           = "mode-channels-enabled"
	flagModeClientsEnabled            = "mode-clients-enabled"
	flagModeClientsMisbehaviour       = "mode-clients-misbehaviour"
	flagModeClientsRefresh            = "mode-clients-refresh"
	flagModeConnectionsEnabled        = "mode-connections-enabled"
	flagModePacketsEnabled            = "mode-packets-enabled"
	flagModePacketsClearInterval      = "mode-packets-clear-interval"
	flagModePacketsClearOnStart       = "mode-packets-clear-on-start"
	flagModePacketsTxConfirmation     = "mode-packets-tx-confirmation"
	flagAutoRegisterCounterpartyPayee = "auto_register_counterparty_payee"
	flagGenerateWallets               = "generate-wallets"
	flagOverwriteConfig               = "overwrite-config"
	flagConfig                        = "config"
)

func main() {
	var (
		args    = os.Args[1:]
		ctx     = context.Background()
		cmdName = args[0]
		c       = &plugin.ExecutedCommand{
			Use:    cmdName,
			Path:   "ignite relayer hermes " + cmdName,
			Args:   args[1:],
			OsArgs: args,
		}
	)
	c.Flags = plugin.Flags{
		{
			Name:       flagConfig,
			Usage:      "set a custom Hermes config file",
			Shorthand:  "c",
			Persistent: true,
			Type:       plugin.FlagTypeString,
		},
	}
	switch cmdName {
	case "configure":
		c.Flags = append(c.Flags, plugin.Flags{
			{Name: flagChainAPortID, DefaultValue: "transfer", Usage: "port ID of the chain A", Type: plugin.FlagTypeString},
			{Name: flagChainBPortID, DefaultValue: "transfer", Usage: "port ID of the chain B", Type: plugin.FlagTypeString},
			{Name: flagChainACCVConsumerChain, DefaultValue: "false", Usage: "only specify true if the chain A is a CCV consumer", Type: plugin.FlagTypeBool},
			{Name: flagChainBCCVConsumerChain, DefaultValue: "false", Usage: "only specify true if the chain B is a CCV consumer", Type: plugin.FlagTypeBool},
			{Name: flagChainAEventSourceMode, DefaultValue: "push", Usage: "WS event source mode of the chain A (event source url should be set to use this flag)", Type: plugin.FlagTypeString},
			{Name: flagChainBEventSourceMode, DefaultValue: "push", Usage: "WS event source mode of the chain B (event source url should be set to use this flag)", Type: plugin.FlagTypeString},
			{Name: flagChainAEventSourceBatchDelay, DefaultValue: "500ms", Usage: "WS event source batch delay time of the chain A (event source url should be set to use this flag)", Type: plugin.FlagTypeString},
			{Name: flagChainBEventSourceBatchDelay, DefaultValue: "500ms", Usage: "WS event source batch delay time of the chain B (event source url should be set to use this flag)", Type: plugin.FlagTypeString},
			{Name: flagChainARPCTimeout, DefaultValue: "10s", Usage: "RPC timeout of the chain A", Type: plugin.FlagTypeString},
			{Name: flagChainBRPCTimeout, DefaultValue: "10s", Usage: "RPC timeout of the chain B", Type: plugin.FlagTypeString},
			{Name: flagChainATrustedNode, DefaultValue: "false", Usage: "enable trusted node on the chain A", Type: plugin.FlagTypeBool},
			{Name: flagChainBTrustedNode, DefaultValue: "false", Usage: "enable trusted node on the chain B", Type: plugin.FlagTypeBool},
			{Name: flagChainAAccountPrefix, DefaultValue: "cosmos", Usage: "account prefix of the chain A", Type: plugin.FlagTypeString},
			{Name: flagChainBAccountPrefix, DefaultValue: "cosmos", Usage: "account prefix of the chain B", Type: plugin.FlagTypeString},
			{Name: flagChainAKeyName, DefaultValue: "wallet", Usage: "hermes account name of the chain A", Type: plugin.FlagTypeString},
			{Name: flagChainBKeyName, DefaultValue: "wallet", Usage: "hermes account name of the chain B", Type: plugin.FlagTypeString},
			{Name: flagChainAAddressType, DefaultValue: "cosmos", Usage: "address type of the chain A", Type: plugin.FlagTypeString},
			{Name: flagChainBAddressType, DefaultValue: "cosmos", Usage: "address type of the chain B", Type: plugin.FlagTypeString},
			{Name: flagChainAKeyStoreType, DefaultValue: "Test", Usage: "address type of the chain B", Type: plugin.FlagTypeString},
			{Name: flagChainBKeyStoreType, DefaultValue: "Test", Usage: "address type of the chain B", Type: plugin.FlagTypeString},
			{Name: flagChainAStorePrefix, DefaultValue: "ibc", Usage: "key store type of the chain A", Type: plugin.FlagTypeString},
			{Name: flagChainBStorePrefix, DefaultValue: "ibc", Usage: "key store type of the chain B", Type: plugin.FlagTypeString},
			{Name: flagChainADefaultGas, DefaultValue: "100000", Usage: "default gas used for transactions on chain A", Type: plugin.FlagTypeUint64},
			{Name: flagChainBDefaultGas, DefaultValue: "100000", Usage: "default gas used for transactions on chain B", Type: plugin.FlagTypeUint64},
			{Name: flagChainAMaxGas, DefaultValue: "400000", Usage: "max gas used for transactions on chain A", Type: plugin.FlagTypeUint64},
			{Name: flagChainBMaxGas, DefaultValue: "400000", Usage: "max gas used for transactions on chain B", Type: plugin.FlagTypeUint64},
			{Name: flagChainAGasPrice, DefaultValue: "0.025stake", Usage: "gas price used for transactions on chain A", Type: plugin.FlagTypeString},
			{Name: flagChainBGasPrice, DefaultValue: "0.025stake", Usage: "gas price used for transactions on chain B", Type: plugin.FlagTypeString},
			{Name: flagChainAGasMultiplier, DefaultValue: "1.1", Usage: "gas multiplier used for transactions on chain A", Type: plugin.FlagTypeString},
			{Name: flagChainBGasMultiplier, DefaultValue: "1.1", Usage: "gas multiplier used for transactions on chain B", Type: plugin.FlagTypeString},
			{Name: flagChainAMaxMsgNum, DefaultValue: "30", Usage: "max message number used for transactions on chain A", Type: plugin.FlagTypeUint64},
			{Name: flagChainBMaxMsgNum, DefaultValue: "30", Usage: "max message number used for transactions on chain B", Type: plugin.FlagTypeUint64},
			{Name: flagChainAMaxTxSize, DefaultValue: "2097152", Usage: "max transaction size on chain A", Type: plugin.FlagTypeUint64},
			{Name: flagChainBMaxTxSize, DefaultValue: "2097152", Usage: "max transaction size on chain B", Type: plugin.FlagTypeUint64},
			{Name: flagChainAClockDrift, DefaultValue: "5s", Usage: "clock drift of the chain A", Type: plugin.FlagTypeString},
			{Name: flagChainBClockDrift, DefaultValue: "5s", Usage: "clock drift of the chain B", Type: plugin.FlagTypeString},
			{Name: flagChainAMaxBlockTime, DefaultValue: "30s", Usage: "maximum block time of the chain A", Type: plugin.FlagTypeString},
			{Name: flagChainBMaxBlockTime, DefaultValue: "30s", Usage: "maximum block time of the chain B", Type: plugin.FlagTypeString},
			{Name: flagChainATrustingPeriod, DefaultValue: "14days", Usage: "trusting period of the chain A", Type: plugin.FlagTypeString},
			{Name: flagChainBTrustingPeriod, DefaultValue: "14days", Usage: "trusting period of the chain B", Type: plugin.FlagTypeString},
			{Name: flagChainATrustThresholdNumerator, DefaultValue: "2", Usage: "trusting threshold numerator of the chain A", Type: plugin.FlagTypeUint64},
			{Name: flagChainBTrustThresholdNumerator, DefaultValue: "2", Usage: "trusting threshold numerator of the chain B", Type: plugin.FlagTypeUint64},
			{Name: flagChainATrustThresholdDenominator, DefaultValue: "3", Usage: "trusting threshold denominator of the chain A", Type: plugin.FlagTypeUint64},
			{Name: flagChainBTrustThresholdDenominator, DefaultValue: "3", Usage: "trusting threshold denominator of the chain B", Type: plugin.FlagTypeUint64},
			{Name: flagChainAType, DefaultValue: "CosmosSdk", Usage: "type of the chain A", Type: plugin.FlagTypeString},
			{Name: flagChainBType, DefaultValue: "CosmosSdk", Usage: "type of the chain B", Type: plugin.FlagTypeString},
			{Name: flagChainASequentialBatchTx, DefaultValue: "false", Usage: "enable sequential batch transaction on the chain A", Type: plugin.FlagTypeBool},
			{Name: flagChainBSequentialBatchTx, DefaultValue: "false", Usage: "enable sequential batch transaction on the chain B", Type: plugin.FlagTypeBool},
			{Name: flagTelemetryEnabled, DefaultValue: "false", Usage: "enable hermes telemetry", Type: plugin.FlagTypeBool},
			{Name: flagTelemetryHost, DefaultValue: "127.0.0.1", Usage: "hermes telemetry host", Type: plugin.FlagTypeString},
			{Name: flagTelemetryPort, DefaultValue: "3001", Usage: "hermes telemetry port", Type: plugin.FlagTypeUint64},
			{Name: flagRestEnabled, DefaultValue: "false", Usage: "enable hermes rest", Type: plugin.FlagTypeBool},
			{Name: flagRestHost, DefaultValue: "127.0.0.1", Usage: "hermes rest host", Type: plugin.FlagTypeString},
			{Name: flagRestPort, DefaultValue: "3000", Usage: "hermes rest port", Type: plugin.FlagTypeUint64},
			{Name: flagModeChannelsEnabled, DefaultValue: "true", Usage: "enable hermes channels", Type: plugin.FlagTypeBool},
			{Name: flagModeClientsEnabled, DefaultValue: "true", Usage: "enable hermes clients", Type: plugin.FlagTypeBool},
			{Name: flagModeClientsMisbehaviour, DefaultValue: "true", Usage: "enable hermes clients misbehaviour", Type: plugin.FlagTypeBool},
			{Name: flagModeClientsRefresh, DefaultValue: "true", Usage: "enable hermes client refresh time", Type: plugin.FlagTypeBool},
			{Name: flagModeConnectionsEnabled, DefaultValue: "true", Usage: "enable hermes connections", Type: plugin.FlagTypeBool},
			{Name: flagModePacketsEnabled, DefaultValue: "true", Usage: "enable hermes packets", Type: plugin.FlagTypeBool},
			{Name: flagModePacketsClearInterval, DefaultValue: "100", Usage: "hermes packet clear interval", Type: plugin.FlagTypeUint64},
			{Name: flagModePacketsClearOnStart, DefaultValue: "true", Usage: "enable hermes packets clear on start", Type: plugin.FlagTypeBool},
			{Name: flagModePacketsTxConfirmation, DefaultValue: "true", Usage: "hermes packet transaction confirmation", Type: plugin.FlagTypeBool},
			{Name: flagAutoRegisterCounterpartyPayee, DefaultValue: "false", Usage: "auto register the counterparty payee on a destination chain to the relayer's address on the source chain", Type: plugin.FlagTypeBool},
			{Name: flagGenerateWallets, DefaultValue: "true", Usage: "automatically generate wallets if they do not exist", Type: plugin.FlagTypeBool},
			{Name: flagOverwriteConfig, DefaultValue: "true", Usage: "overwrite the current config if it already exists", Type: plugin.FlagTypeBool},
			{Name: flagChainAFaucet, DefaultValue: "http://0.0.0.0:4501", Type: plugin.FlagTypeString},
			{Name: flagChainBFaucet, DefaultValue: "http://0.0.0.0:4500", Type: plugin.FlagTypeString},
		}...)
		if err := cmd.ConfigureHandler(ctx, c); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	case "exec":
		if err := cmd.ExecuteHandler(ctx, c); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	case "start":
		if err := cmd.StartHandler(ctx, c); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
	case "keys":
		switch args[1] {
		case "add":
			if err := cmd.KeysAddMnemonicHandler(ctx, c); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		case "file":
			if err := cmd.KeysAddFileHandler(ctx, c); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		case "list":
			if err := cmd.KeysListHandler(ctx, c); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		case "delete":
			if err := cmd.KeysDeleteHandler(ctx, c); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
		default:
			err := errors.Errorf("unknown keys command: %s", args[1])
			fmt.Fprintln(os.Stderr, err)
			return
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown hermes command: %s", cmdName)
		return
	}
}
