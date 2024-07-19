package cmd

import (
	"strconv"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

const (
	flagChainAPortID                    = "chain-a-port-id"
	flagChainAEventSourceMode           = "chain-a-event-source-mode"
	flagChainAEventSourceURL            = "chain-a-event-source-url"
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
	flagChainAMemoPrefix                = "chain-a-memo-prefix"
	flagChainAType                      = "chain-a-type"
	flagChainASequentialBatchTx         = "chain-a-sequential-batch-tx"

	flagChainBPortID                    = "chain-b-port-id"
	flagChainBEventSourceMode           = "chain-b-event-source-mode"
	flagChainBEventSourceURL            = "chain-b-event-source-url"
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
	flagChainBMemoPrefix                = "chain-b-memo-prefix"
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
	flagChannelVersion                = "channel-version"

	flagConfig = "config"

	mnemonicEntropySize = 256
)

var (
	ErrFlagNotFound  = errors.New("flag not found")
	ErrFlagAssertion = errors.New("flag type assertion to failed")
)

func getConfig(flags []*plugin.Flag) string {
	config, _ := getFlag[string](flags, flagConfig)
	return config
}

// getFlag function to get flag with a generic return type
// TODO remove these helpers for flags after we fix this issue:
// https://github.com/ignite/apps/issues/116
func getFlag[A any](flags []*plugin.Flag, key string) (result A, err error) {
	v, err := getValue(flags, key)
	if err != nil {
		return result, err
	}

	value, ok := v.(A)
	if !ok {
		return result, errors.Wrapf(ErrFlagAssertion, "type assertion to %T failed for field %s", v, key)
	}
	return value, nil
}

func getValue(flags []*plugin.Flag, key string) (interface{}, error) {
	for _, flag := range flags {
		if flag.Name == key {
			return exportToFlagValue(flag)
		}
	}
	return nil, errors.Wrapf(ErrFlagNotFound, "flag %s not found", key)
}

func exportToFlagValue(f *plugin.Flag) (interface{}, error) {
	switch f.Type {
	case plugin.FlagTypeBool:
		v, err := strconv.ParseBool(flagValue(f))
		if err != nil {
			return false, err
		}
		return v, nil
	case plugin.FlagTypeInt:
		v, err := strconv.Atoi(flagValue(f))
		if err != nil {
			return 0, err
		}
		return v, nil
	case plugin.FlagTypeUint:
		v, err := strconv.ParseUint(flagValue(f), 10, 64)
		if err != nil {
			return uint(0), err
		}
		return uint(v), nil
	case plugin.FlagTypeInt64:
		v, err := strconv.ParseInt(flagValue(f), 10, 64)
		if err != nil {
			return int64(0), err
		}
		return v, nil
	case plugin.FlagTypeUint64:
		v, err := strconv.ParseUint(flagValue(f), 10, 64)
		if err != nil {
			return uint64(0), err
		}
		return v, nil
	case plugin.FlagTypeStringSlice:
		v := strings.Trim(flagValue(f), "[]")
		s := strings.Split(v, ",")
		if len(s) == 0 || (len(s) == 1 && s[0] == "") {
			return []string{}, nil
		}

		return s, nil
	default:
		return strings.TrimSpace(flagValue(f)), nil
	}
}

func flagValue(flag *plugin.Flag) string {
	if flag.Value != "" {
		return flag.Value
	}
	return flag.DefaultValue
}
