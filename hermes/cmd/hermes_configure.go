package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gookit/color"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/cliquiz"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"relayer/pkg/hermes"
)

const (
	flagChainAPortID                    = "chain-a-port-id"
	flagChainAEventSourceMode           = "chain-a-event-source-mode"
	flagChainAEventSourceURL            = "chain-a-event-source-url"
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

	flagChainBPortID                    = "chain-b-port-id"
	flagChainBEventSourceMode           = "chain-b-event-source-mode"
	flagChainBEventSourceURL            = "chain-b-event-source-url"
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
		Short: "Configure the Hermes realyer creating the config file, client, channels and connection",
		Args:  cobra.ExactArgs(6),
		RunE:  hermesConfigureHandler,
	}

	c.Flags().String(flagChainAPortID, "transfer", "Port ID of the chain A")
	c.Flags().String(flagChainBPortID, "transfer", "Port ID of the chain B")
	c.Flags().String(flagChainAEventSourceURL, "", "WS event source url of the chain A")
	c.Flags().String(flagChainBEventSourceURL, "", "WS event source url of the chain B")
	c.Flags().String(flagChainAEventSourceMode, "push", "WS event source mode of the chain A (event source url should be set to use this flag)")
	c.Flags().String(flagChainBEventSourceMode, "push", "WS event source mode of the chain B (event source url should be set to use this flag)")
	c.Flags().String(flagChainAEventSourceBatchDelay, "500ms", "WS event source batch delay time of the chain A (event source url should be set to use this flag)")
	c.Flags().String(flagChainBEventSourceBatchDelay, "500ms", "WS event source batch delay time of the chain B (event source url should be set to use this flag)")
	c.Flags().String(flagChainARPCTimeout, "15s", "RPC timeout of the chain A")
	c.Flags().String(flagChainBRPCTimeout, "15s", "RPC timeout of the chain B")
	c.Flags().String(flagChainAAccountPrefix, "cosmos", "address prefix of the chain A")
	c.Flags().String(flagChainBAccountPrefix, "cosmos", "address prefix of the chain B")
	c.Flags().String(flagChainAKeyName, "wallet", "hermes account name of the chain A")
	c.Flags().String(flagChainBKeyName, "wallet", "hermes account name of the chain B")
	c.Flags().String(flagChainAStorePrefix, "ibc", "store prefix of the chain A")
	c.Flags().String(flagChainBStorePrefix, "ibc", "store prefix of the chain B")
	c.Flags().Uint64(flagChainADefaultGas, 100000, "default gas used for transactions on chain A")
	c.Flags().Uint64(flagChainBDefaultGas, 100000, "default gas used for transactions on chain B")
	c.Flags().Uint64(flagChainAMaxGas, 10000000, "max gas used for transactions on chain A")
	c.Flags().Uint64(flagChainBMaxGas, 10000000, "max gas used for transactions on chain B")
	c.Flags().String(flagChainAGasPrice, "1stake", "gas price used for transactions on chain A")
	c.Flags().String(flagChainBGasPrice, "1stake", "gas price used for transactions on chain B")
	c.Flags().String(flagChainAGasMultiplier, "1.1", "gas multiplier used for transactions on chain A")
	c.Flags().String(flagChainBGasMultiplier, "1.1", "gas multiplier used for transactions on chain B")
	c.Flags().Uint64(flagChainAMaxMsgNum, 30, "max message number used for transactions on chain A")
	c.Flags().Uint64(flagChainBMaxMsgNum, 30, "max message number used for transactions on chain B")
	c.Flags().Uint64(flagChainAMaxTxSize, 2097152, "max transaction size on chain A")
	c.Flags().Uint64(flagChainBMaxTxSize, 2097152, "max transaction size on chain B")
	c.Flags().String(flagChainAClockDrift, "5s", "clock drift of the chain A")
	c.Flags().String(flagChainBClockDrift, "5s", "clock drift of the chain B")
	c.Flags().String(flagChainAMaxBlockTime, "10s", "maximum block time of the chain A")
	c.Flags().String(flagChainBMaxBlockTime, "10s", "maximum block time of the chain B")
	c.Flags().String(flagChainATrustingPeriod, "14days", "trusting period of the chain A")
	c.Flags().String(flagChainBTrustingPeriod, "14days", "trusting period of the chain B")
	c.Flags().Uint64(flagChainATrustThresholdNumerator, 1, "trusting threshold numerator of the chain A")
	c.Flags().Uint64(flagChainBTrustThresholdNumerator, 1, "trusting threshold numerator of the chain B")
	c.Flags().Uint64(flagChainATrustThresholdDenominator, 3, "trusting threshold denominator of the chain A")
	c.Flags().Uint64(flagChainBTrustThresholdDenominator, 3, "trusting threshold denominator of the chain B")

	c.Flags().Bool(flagTelemetryEnabled, false, "enable hermes telemetry")
	c.Flags().String(flagTelemetryHost, "127.0.0.1", "hermes telemetry host")
	c.Flags().Uint64(flagTelemetryPort, 3001, "hermes telemetry port")
	c.Flags().Bool(flagModeChannelsEnabled, true, "enable hermes channels")
	c.Flags().Bool(flagModeClientsEnabled, true, "enable hermes clients")
	c.Flags().Bool(flagModeClientsMisbehaviour, true, "enable hermes clients misbehaviour")
	c.Flags().Bool(flagModeClientsRefresh, true, "enable hermes client refresh time")
	c.Flags().Bool(flagModeConnectionsEnabled, true, "enable hermes connections")
	c.Flags().Bool(flagModePacketsEnabled, true, "enable hermes packets")
	c.Flags().Uint64(flagModePacketsClearInterval, 100, "hermes packet clear interval")
	c.Flags().Bool(flagModePacketsClearOnStart, true, "enable hermes packets clear on start")
	c.Flags().Bool(flagModePacketsTxConfirmation, true, "hermes packet transaction confirmation")

	return c
}

func hermesConfigureHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	session.StartSpinner("Generating Hermes config")

	var (
		chainAID = args[0]
		chainBID = args[3]

		chainAPortID, _ = cmd.Flags().GetString(flagChainAPortID)
		chainBPortID, _ = cmd.Flags().GetString(flagChainBPortID)
		customCfg       = getConfig(cmd)
	)

	cfgName := strings.Join([]string{args[0], args[3]}, hermes.ConfigNameSeparator)
	cfgPath, err := hermes.ConfigPath(cfgName)
	if err != nil {
		return err
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		if err := hermesCreateConfig(cmd, args, customCfg); err != nil {
			return err
		}
	} else {
		if err := session.AskConfirm(fmt.Sprintf(
			"Hermes %s <-> %s config already exist at %s. Do you want to reuse this config file?",
			chainAID,
			chainBID,
			cfgPath,
		)); err != nil {
			if !errors.Is(err, promptui.ErrAbort) {
				return err
			}
			if err := hermesCreateConfig(cmd, args, customCfg); err != nil {
				return err
			}
		}
	}

	_ = session.Println(color.Green.Sprintf("Hermes config created at %s", cfgPath))

	h, err := hermes.New()
	if err != nil {
		return err
	}
	defer h.Cleanup()

	session.StartSpinner("Verifying chain keys")

	if err := verifyChainKeys(cmd.Context(), session, h, chainAID, cfgPath); err != nil {
		return err
	}

	if err := verifyChainKeys(cmd.Context(), session, h, chainBID, cfgPath); err != nil {
		return err
	}

	session.StartSpinner("Creating clients")

	// create client A
	var (
		bufClientAResult = bytes.Buffer{}
		clientAResult    = hermes.ClientResult{}
	)
	err = h.CreateClient(
		cmd.Context(),
		chainAID,
		chainBID,
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&bufClientAResult),
	)
	if err != nil {
		return err
	}
	if err := hermes.UnmarshalResult(bufClientAResult.Bytes(), &clientAResult); err != nil {
		return err
	}

	_ = session.Println(color.Green.Sprintf(
		"Client '%s' created (%s -> %s)",
		clientAResult.CreateClient.ClientID,
		chainAID,
		chainBID,
	))

	// create client B
	var (
		bufClientBResult = bytes.Buffer{}
		clientBResult    = hermes.ClientResult{}
	)
	err = h.CreateClient(
		cmd.Context(),
		chainBID,
		chainAID,
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&bufClientBResult),
	)
	if err != nil {
		return err
	}
	if err := hermes.UnmarshalResult(bufClientBResult.Bytes(), &clientBResult); err != nil {
		return err
	}

	_ = session.Println(color.Green.Sprintf(
		"Client %s' created (%s -> %s)",
		clientBResult.CreateClient.ClientID,
		chainBID,
		chainAID,
	))
	session.StartSpinner("Creating connection")

	// create connection
	var (
		bufConnection = bytes.Buffer{}
		connection    = hermes.ConnectionResult{}
	)
	err = h.CreateConnection(
		cmd.Context(),
		chainAID,
		clientAResult.CreateClient.ClientID,
		clientBResult.CreateClient.ClientID,
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&bufConnection),
	)
	if err != nil {
		return err
	}
	if err := hermes.UnmarshalResult(bufConnection.Bytes(), &connection); err != nil {
		return err
	}

	_ = session.Println(color.Green.Sprintf(
		"Connection %s <-> %s created",
		connection.ASide.ConnectionID,
		connection.BSide.ConnectionID,
	))
	session.StartSpinner("Creating channel")

	// create and query channel
	var (
		bufChannel = bytes.Buffer{}
		channel    = hermes.ConnectionResult{}
	)
	err = h.CreateChannel(
		cmd.Context(),
		chainAID,
		connection.ASide.ConnectionID,
		chainAPortID,
		chainBPortID,
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&bufChannel),
	)
	if err != nil {
		return err
	}

	if err := hermes.UnmarshalResult(bufChannel.Bytes(), &channel); err != nil {
		return err
	}

	_ = session.Println(color.Green.Sprintf(
		"Channel '%s <-> %s' created",
		chainAID,
		chainBID,
	))

	return nil
}

func verifyChainKeys(ctx context.Context, session *cliui.Session, h *hermes.Hermes, chainID, cfgPath string) error {
	var (
		bufKeysChain    = bytes.Buffer{}
		keysChainResult = hermes.KeysListResult{}
	)
	if err := h.KeysList(
		ctx,
		chainID,
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&bufKeysChain),
	); err != nil {
		return err
	}
	if err := hermes.UnmarshalResult(bufKeysChain.Bytes(), &keysChainResult); err != nil {
		return err
	}
	if keysChainResult.Wallet.Account == "" {
		var chainAMnemonic string
		if err := session.Ask(cliquiz.NewQuestion(
			fmt.Sprintf("Chain %s doesn't have a default Hermes key. Type your mnemonic to continue:", chainID),
			&chainAMnemonic,
			cliquiz.Required(),
		)); err != nil {
			return err
		}

		bufKeysChainAdd := bytes.Buffer{}
		err := h.AddMnemonic(
			ctx,
			chainID,
			chainAMnemonic,
			hermes.WithConfigFile(cfgPath),
			hermes.WithStdOut(&bufKeysChainAdd),
		)
		if err != nil {
			return err
		}
		if err := hermes.ValidateResult(bufKeysChainAdd.Bytes()); err != nil {
			return err
		}

		return session.Println(color.Yellow.Sprintf("Chain %s key created", chainID))
	}
	return nil
}

func hermesCreateConfig(cmd *cobra.Command, args []string, customCfg string) error {
	// if a custom config was set, save it in the ignite hermes folder
	if customCfg != "" {
		c, err := hermes.LoadConfig(customCfg)
		if err != nil {
			return err
		}
		return c.Save()
	}

	// Create the default hermes config
	var (
		telemetryEnabled, _          = cmd.Flags().GetBool(flagTelemetryEnabled)
		telemetryHost, _             = cmd.Flags().GetString(flagTelemetryHost)
		telemetryPort, _             = cmd.Flags().GetUint64(flagTelemetryPort)
		modeChannelsEnabled, _       = cmd.Flags().GetBool(flagModeChannelsEnabled)
		modeClientsEnabled, _        = cmd.Flags().GetBool(flagModeClientsEnabled)
		modeClientsMisbehaviour, _   = cmd.Flags().GetBool(flagModeClientsMisbehaviour)
		modeClientsRefresh, _        = cmd.Flags().GetBool(flagModeClientsRefresh)
		modeConnectionsEnabled, _    = cmd.Flags().GetBool(flagModeConnectionsEnabled)
		modePacketsEnabled, _        = cmd.Flags().GetBool(flagModePacketsEnabled)
		modePacketsClearInterval, _  = cmd.Flags().GetUint64(flagModePacketsClearInterval)
		modePacketsClearOnStart, _   = cmd.Flags().GetBool(flagModePacketsClearOnStart)
		modePacketsTxConfirmation, _ = cmd.Flags().GetBool(flagModePacketsTxConfirmation)
	)

	c := hermes.DefaultConfig(
		hermes.WithTelemetryEnabled(telemetryEnabled),
		hermes.WithTelemetryHost(telemetryHost),
		hermes.WithTelemetryPort(telemetryPort),
		hermes.WithModeChannelsEnabled(modeChannelsEnabled),
		hermes.WithModeClientsEnabled(modeClientsEnabled),
		hermes.WithModeClientsMisbehaviour(modeClientsMisbehaviour),
		hermes.WithModeClientsRefresh(modeClientsRefresh),
		hermes.WithModeConnectionsEnabled(modeConnectionsEnabled),
		hermes.WithModePacketsEnabled(modePacketsEnabled),
		hermes.WithModePacketsClearInterval(modePacketsClearInterval),
		hermes.WithModePacketsClearOnStart(modePacketsClearOnStart),
		hermes.WithModePacketsTxConfirmation(modePacketsTxConfirmation),
	)

	// Add chain A into the config
	var (
		chainAID       = args[0]
		chainARPCAddr  = args[1]
		chainAGRPCAddr = args[2]

		chainAEventSourceMode, _           = cmd.Flags().GetString(flagChainAEventSourceMode)
		chainAEventSourceURL, _            = cmd.Flags().GetString(flagChainAEventSourceURL)
		chainAEventSourceBatchDelay, _     = cmd.Flags().GetString(flagChainAEventSourceBatchDelay)
		chainARPCTimeout, _                = cmd.Flags().GetString(flagChainARPCTimeout)
		chainAAccountPrefix, _             = cmd.Flags().GetString(flagChainAAccountPrefix)
		chainAKeyName, _                   = cmd.Flags().GetString(flagChainAKeyName)
		chainAStorePrefix, _               = cmd.Flags().GetString(flagChainAStorePrefix)
		chainADefaultGas, _                = cmd.Flags().GetUint64(flagChainADefaultGas)
		chainAMaxGas, _                    = cmd.Flags().GetUint64(flagChainAMaxGas)
		chainAGasPrice, _                  = cmd.Flags().GetString(flagChainAGasPrice)
		chainAGasMultiplier, _             = cmd.Flags().GetString(flagChainAGasMultiplier)
		chainAMaxMsgNum, _                 = cmd.Flags().GetUint64(flagChainAMaxMsgNum)
		chainAMaxTxSize, _                 = cmd.Flags().GetUint64(flagChainAMaxTxSize)
		chainAClockDrift, _                = cmd.Flags().GetString(flagChainAClockDrift)
		chainAMaxBlockTime, _              = cmd.Flags().GetString(flagChainAMaxBlockTime)
		chainATrustingPeriod, _            = cmd.Flags().GetString(flagChainATrustingPeriod)
		chainATrustThresholdNumerator, _   = cmd.Flags().GetUint64(flagChainATrustThresholdNumerator)
		chainATrustThresholdDenominator, _ = cmd.Flags().GetUint64(flagChainATrustThresholdDenominator)
	)

	chainAGasMulti := new(big.Float)
	chainAGasMulti, ok := chainAGasMulti.SetString(chainAGasMultiplier)
	if !ok {
		return fmt.Errorf("invalid chain A gas multiplier: %s", chainAGasMultiplier)
	}

	optChainA := []hermes.ChainOption{
		hermes.WithChainTrustThreshold(chainATrustThresholdNumerator, chainATrustThresholdDenominator),
		hermes.WithChainGasMultiplier(chainAGasMulti),
	}
	if chainAEventSourceURL != "" {
		optChainA = append(optChainA, hermes.WithChainEventSource(
			chainAEventSourceMode,
			chainAEventSourceURL,
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

	// Add chain B into the config
	var (
		chainBID       = args[3]
		chainBRPCAddr  = args[4]
		chainBGRPCAddr = args[5]

		chainBEventSourceMode, _           = cmd.Flags().GetString(flagChainBEventSourceMode)
		chainBEventSourceURL, _            = cmd.Flags().GetString(flagChainBEventSourceURL)
		chainBEventSourceBatchDelay, _     = cmd.Flags().GetString(flagChainBEventSourceBatchDelay)
		chainBRPCTimeout, _                = cmd.Flags().GetString(flagChainBRPCTimeout)
		chainBAccountPrefix, _             = cmd.Flags().GetString(flagChainBAccountPrefix)
		chainBKeyName, _                   = cmd.Flags().GetString(flagChainBKeyName)
		chainBStorePrefix, _               = cmd.Flags().GetString(flagChainBStorePrefix)
		chainBDefaultGas, _                = cmd.Flags().GetUint64(flagChainBDefaultGas)
		chainBMaxGas, _                    = cmd.Flags().GetUint64(flagChainBMaxGas)
		chainBGasPrice, _                  = cmd.Flags().GetString(flagChainBGasPrice)
		chainBGasMultiplier, _             = cmd.Flags().GetString(flagChainBGasMultiplier)
		chainBMaxMsgNum, _                 = cmd.Flags().GetUint64(flagChainBMaxMsgNum)
		chainBMaxTxSize, _                 = cmd.Flags().GetUint64(flagChainBMaxTxSize)
		chainBClockDrift, _                = cmd.Flags().GetString(flagChainBClockDrift)
		chainBMaxBlockTime, _              = cmd.Flags().GetString(flagChainBMaxBlockTime)
		chainBTrustingPeriod, _            = cmd.Flags().GetString(flagChainBTrustingPeriod)
		chainBTrustThresholdNumerator, _   = cmd.Flags().GetUint64(flagChainBTrustThresholdNumerator)
		chainBTrustThresholdDenominator, _ = cmd.Flags().GetUint64(flagChainBTrustThresholdDenominator)
	)

	chainBGasMulti := new(big.Float)
	chainBGasMulti, ok = chainBGasMulti.SetString(chainBGasMultiplier)
	if !ok {
		return fmt.Errorf("invalid chain B gas multiplier: %s", chainBGasMultiplier)
	}

	optChainB := []hermes.ChainOption{
		hermes.WithChainTrustThreshold(chainBTrustThresholdNumerator, chainBTrustThresholdDenominator),
		hermes.WithChainGasMultiplier(chainBGasMulti),
	}
	if chainBEventSourceURL != "" {
		optChainB = append(optChainB, hermes.WithChainEventSource(
			chainBEventSourceMode,
			chainBEventSourceURL,
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

	if err := c.AddChain(chainBID, chainBRPCAddr, chainBGRPCAddr, optChainB...); err != nil {
		return err
	}

	return c.Save()
}
