package cmd

import (
	"strings"

	"github.com/blang/semver/v4"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/hermes/pkg/hermes"
)

const (
	flagChainAPortID                = "chain-a-port-id"
	flagChainAEventSourceMode       = "chain-a-event-source-mode"
	flagChainAEventSourceURL        = "chain-a-event-source-url"
	flagChainAEventSourceBatchDelay = "chain-a-event-source-batch-delay"
	flagChainARPCTimeout            = "chain-a-rpc-timeout"
	flagChainAAccountPrefix         = "chain-a-account-prefix"
	flagChainAAddressType           = "chain-a-address-types"
	flagChainAKeyName               = "chain-a-key-name"
	flagChainAKeyStoreType          = "chain-a-key-store-type"
	flagChainAStorePrefix           = "chain-a-store-prefix"
	flagChainADefaultGas            = "chain-a-default-gas"
	flagChainAMaxGas                = "chain-a-max-gas"
	flagChainAGasPrice              = "chain-a-gas-price"
	flagChainAGasMultiplier         = "chain-a-gas-multiplier"
	flagChainAMaxMsgNum             = "chain-a-max-msg-num"
	flagChainAMaxTxSize             = "chain-a-tx-size"
	flagChainAClockDrift            = "chain-a-clock-drift"
	flagChainAMaxBlockTime          = "chain-a-max-block-time"
	flagChainATrustingPeriod        = "chain-a-trusting-period"
	flagChainATrustThreshold        = "chain-a-trust-threshold"
	flagChainAFaucet                = "chain-a-faucet"
	flagChainACCVConsumerChain      = "chain-a-ccv-consumer-chain"
	flagChainATrustedNode           = "chain-a-trusted-node"
	flagChainAMemoPrefix            = "chain-a-memo-prefix"
	flagChainAType                  = "chain-a-type"
	flagChainASequentialBatchTx     = "chain-a-sequential-batch-tx"

	flagChainBPortID                = "chain-b-port-id"
	flagChainBEventSourceMode       = "chain-b-event-source-mode"
	flagChainBEventSourceURL        = "chain-b-event-source-url"
	flagChainBEventSourceBatchDelay = "chain-b-event-source-batch-delay"
	flagChainBRPCTimeout            = "chain-b-rpc-timeout"
	flagChainBAccountPrefix         = "chain-b-account-prefix"
	flagChainBAddressType           = "chain-b-address-types"
	flagChainBKeyName               = "chain-b-key-name"
	flagChainBKeyStoreType          = "chain-b-key-store-type"
	flagChainBStorePrefix           = "chain-b-store-prefix"
	flagChainBDefaultGas            = "chain-b-default-gas"
	flagChainBMaxGas                = "chain-b-max-gas"
	flagChainBGasPrice              = "chain-b-gas-price"
	flagChainBGasMultiplier         = "chain-b-gas-multiplier"
	flagChainBMaxMsgNum             = "chain-b-max-msg-num"
	flagChainBMaxTxSize             = "chain-b-tx-size"
	flagChainBClockDrift            = "chain-b-clock-drift"
	flagChainBMaxBlockTime          = "chain-b-max-block-time"
	flagChainBTrustingPeriod        = "chain-b-trusting-period"
	flagChainBTrustThreshold        = "chain-b-trust-threshold"
	flagChainBFaucet                = "chain-b-faucet"
	flagChainBCCVConsumerChain      = "chain-b-ccv-consumer-chain"
	flagChainBTrustedNode           = "chain-b-trusted-node"
	flagChainBMemoPrefix            = "chain-b-memo-prefix"
	flagChainBType                  = "chain-b-type"
	flagChainBSequentialBatchTx     = "chain-b-sequential-batch-tx"

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
	flagChannelVersion                = "channel-version"

	flagConfig        = "config"
	flagHermesVersion = "version"

	mnemonicEntropySize = 256
)

func getConfig(flags plugin.Flags) string {
	config, _ := flags.GetString(flagConfig)
	return config
}

func getVersion(flags plugin.Flags) (string, error) {
	version, _ := flags.GetString(flagHermesVersion)
	if version == "" {
		version = hermes.DefaultVersion
	}
	sv, err := semver.Parse(strings.TrimPrefix(version, "v"))
	if err != nil {
		return version, errors.Wrapf(err, "invalid version format %s", version)
	}
	return "v" + sv.String(), nil
}
