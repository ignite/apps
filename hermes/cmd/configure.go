package cmd

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/gookit/color"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/cliquiz"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/manifoldco/promptui"
	"github.com/spf13/pflag"

	"github.com/ignite/apps/hermes/pkg/hermes"
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

	flagChainBPortID                    = "chain-b-port-id"
	flagChainBEventSourceMode           = "chain-b-event-source-mode"
	flagChainBEventSourceURL            = "chain-b-event-source-url"
	flagChainBEventSourceBatchDelay     = "chain-b-event-source-batch-delay"
	flagChainBRPCTimeout                = "chain-b-rpc-timeout"
	flagChainBAccountPrefix             = "chain-b-account-prefix"
	flagChainBAddressType               = "chain-b-address-types"
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
	flagChainBFaucet                    = "chain-b-faucet"
	flagChainBCCVConsumerChain          = "chain-b-ccv-consumer-chain"
	flagChainBTrustedNode               = "chain-b-trusted-node"
	flagChainBMemoPrefix                = "chain-b-memo-prefix"

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

	mnemonicEntropySize = 256
)

func ConfigureHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	args := cmd.Args
	flags, err := cmd.NewFlags()
	if err != nil {
		return err
	}

	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	session.StartSpinner("Generating Hermes config")

	var (
		chainAID = args[0]
		chainBID = args[3]

		generateWallets, _ = flags.GetBool(flagGenerateWallets)
		overwriteConfig, _ = flags.GetBool(flagOverwriteConfig)
		chainAPortID, _    = flags.GetString(flagChainAPortID)
		chainAFaucet, _    = flags.GetString(flagChainAFaucet)
		chainBPortID, _    = flags.GetString(flagChainBPortID)
		chainBFaucet, _    = flags.GetString(flagChainBFaucet)
		customCfg          = getConfig(flags)
	)

	var hermesCfg *hermes.Config
	if customCfg != "" {
		hermesCfg, err = hermes.LoadConfig(customCfg)
		if err != nil {
			return err
		}
	} else {
		hermesCfg, err = newHermesConfig(flags, args, customCfg)
		if err != nil {
			return err
		}
	}
	cfgPath, err := hermesCfg.ConfigPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(cfgPath); overwriteConfig || os.IsNotExist(err) {
		if err := hermesCfg.Save(); err != nil {
			return err
		}
	} else {
		if err := session.AskConfirm(fmt.Sprintf(
			"Hermes %s <-> %s config already exist at %s. Do you want to reuse this config file",
			chainAID,
			chainBID,
			cfgPath,
		)); err != nil {
			if !errors.Is(err, promptui.ErrAbort) {
				return err
			}
			if err := hermesCfg.Save(); err != nil {
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

	if err := ensureAccount(
		ctx,
		session,
		hermesCfg,
		h,
		chainAID,
		chainAFaucet,
		cfgPath,
		generateWallets,
	); err != nil {
		return err
	}

	if err := ensureAccount(
		ctx,
		session,
		hermesCfg,
		h,
		chainBID,
		chainBFaucet,
		cfgPath,
		generateWallets,
	); err != nil {
		return err
	}

	session.StartSpinner("Creating clients")

	// create client A
	var (
		bufClientAResult = bytes.Buffer{}
		clientAResult    = hermes.ClientResult{}
	)
	if err := h.CreateClient(
		ctx,
		chainAID,
		chainBID,
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&bufClientAResult),
		hermes.WithJSONOutput(),
	); err != nil {
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
	if err := h.CreateClient(
		ctx,
		chainBID,
		chainAID,
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&bufClientBResult),
		hermes.WithJSONOutput(),
	); err != nil {
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
	if err := h.CreateConnection(
		ctx,
		chainAID,
		clientAResult.CreateClient.ClientID,
		clientBResult.CreateClient.ClientID,
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&bufConnection),
		hermes.WithJSONOutput(),
	); err != nil {
		return err
	}
	if err := hermes.UnmarshalResult(bufConnection.Bytes(), &connection); err != nil {
		return err
	}

	_ = session.Println(color.Green.Sprintf(
		"Connection '%s (%s) <-> %s (%s)' created",
		chainAID,
		connection.ASide.ConnectionID,
		chainBID,
		connection.BSide.ConnectionID,
	))
	session.StartSpinner("Creating channel")

	// create and query channel
	var (
		bufChannel = bytes.Buffer{}
		channel    = hermes.ConnectionResult{}
	)
	if err := h.CreateChannel(
		ctx,
		chainAID,
		connection.ASide.ConnectionID,
		chainAPortID,
		chainBPortID,
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&bufChannel),
		hermes.WithJSONOutput(),
	); err != nil {
		return err
	}
	if err := hermes.UnmarshalResult(bufChannel.Bytes(), &channel); err != nil {
		return err
	}

	_ = session.Println(color.Green.Sprintf(
		"Channel '%s (%s) <-> %s (%s)' created",
		chainAID,
		channel.ASide.ChannelID,
		chainBID,
		channel.BSide.ChannelID,
	))

	return nil
}

// ensureAccount ensures the account exists and get found if the faucet is set.
func ensureAccount(
	ctx context.Context,
	session *cliui.Session,
	hCfg *hermes.Config,
	h *hermes.Hermes,
	chainID,
	faucetAddr,
	cfgPath string,
	generateWallets bool,
) error {
	chainAddr, err := verifyChainKeys(ctx, session, h, chainID, cfgPath, generateWallets)
	if err != nil {
		return err
	}
	chain, err := hCfg.Chains.Get(chainID)
	if err != nil {
		return err
	}
	balance, err := chain.Balance(ctx, chainAddr)
	if err != nil {
		return err
	}
	if balance.Empty() && faucetAddr == "" {
		return errors.Errorf(
			"wallet %s balance is empty, please add funds or provide the faucet address flag (--%s or --%s)",
			chainAddr,
			flagChainAFaucet,
			flagChainBFaucet)
	}
	if faucetAddr != "" {
		_, err := chain.TryRetrieve(ctx, chainAddr, faucetAddr)
		if err != nil {
			return err
		}
	}
	return nil
}

// verifyChainKeys verifies if the Hermes has a key for the specific chain,
// if not,  ask for the user to create one.
func verifyChainKeys(
	ctx context.Context,
	session *cliui.Session,
	h *hermes.Hermes,
	chainID,
	cfgPath string,
	generateWallets bool,
) (string, error) {
GetKey:
	var (
		bufKeysChain    = bytes.Buffer{}
		keysChainResult = hermes.KeysListResult{}
	)
	if err := h.KeysList(
		ctx,
		chainID,
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&bufKeysChain),
		hermes.WithJSONOutput(),
	); err != nil {
		return "", err
	}
	if err := hermes.UnmarshalResult(bufKeysChain.Bytes(), &keysChainResult); err != nil {
		return "", err
	}
	if keysChainResult.Wallet.Account == "" {
		var mnemonic string
		if !generateWallets {
			if err := session.Ask(cliquiz.NewQuestion(
				fmt.Sprintf(
					"Chain %s doesn't have a default Hermes key. Type your mnemonic to continue or type enter to generate a new one:",
					chainID,
				),
				&mnemonic,
			)); err != nil {
				return "", err
			}
		}

		if mnemonic == "" {
			entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
			if err != nil {
				return "", err
			}
			mnemonic, err = bip39.NewMnemonic(entropySeed)
			if err != nil {
				return "", err
			}
			_ = session.Printf(
				"%s %s\n",
				color.Yellow.Sprint("New mnemonic generated:"),
				color.Blue.Sprint(mnemonic),
			)
		}

		if !bip39.IsMnemonicValid(mnemonic) {
			return "", errors.Errorf("invalid mnemonic: %s", mnemonic)
		}

		bufKeysChainAdd := bytes.Buffer{}
		if err := h.AddMnemonic(
			ctx,
			chainID,
			mnemonic,
			hermes.WithConfigFile(cfgPath),
			hermes.WithStdOut(&bufKeysChainAdd),
			hermes.WithJSONOutput(),
		); err != nil {
			return "", err
		}
		if err := hermes.ValidateResult(bufKeysChainAdd.Bytes()); err != nil {
			return "", err
		}
		_ = session.Println(color.Yellow.Sprintf("Chain %s key created", chainID))

		goto GetKey
	}
	_ = session.Println(color.Green.Sprintf("Chain %s relayer wallet: %s", chainID, keysChainResult.Wallet.Account))
	return keysChainResult.Wallet.Account, nil
}

// newHermesConfig create a new hermes config based in the cmd args.
func newHermesConfig(flags *pflag.FlagSet, args []string, customCfg string) (*hermes.Config, error) {
	// if a custom config was set, save it in the ignite hermes folder
	if customCfg != "" {
		c, err := hermes.LoadConfig(customCfg)
		if err != nil {
			return nil, err
		}
		return c, c.Save()
	}

	// Create the default hermes config
	var (
		telemetryEnabled, _                         = flags.GetBool(flagTelemetryEnabled)
		telemetryHost, _                            = flags.GetString(flagTelemetryHost)
		telemetryPort, _                            = flags.GetUint64(flagTelemetryPort)
		restEnabled, _                              = flags.GetBool(flagRestEnabled)
		restHost, _                                 = flags.GetString(flagRestHost)
		restPort, _                                 = flags.GetUint64(flagRestPort)
		modeChannelsEnabled, _                      = flags.GetBool(flagModeChannelsEnabled)
		modeClientsEnabled, _                       = flags.GetBool(flagModeClientsEnabled)
		modeClientsMisbehaviour, _                  = flags.GetBool(flagModeClientsMisbehaviour)
		modeClientsRefresh, _                       = flags.GetBool(flagModeClientsRefresh)
		modeConnectionsEnabled, _                   = flags.GetBool(flagModeConnectionsEnabled)
		modePacketsEnabled, _                       = flags.GetBool(flagModePacketsEnabled)
		modePacketsClearInterval, _                 = flags.GetUint64(flagModePacketsClearInterval)
		modePacketsClearOnStart, _                  = flags.GetBool(flagModePacketsClearOnStart)
		modePacketsTxConfirmation, _                = flags.GetBool(flagModePacketsTxConfirmation)
		modePacketsAutoRegisterCounterpartyPayee, _ = flags.GetBool(flagAutoRegisterCounterpartyPayee)
	)

	c := hermes.DefaultConfig(
		hermes.WithTelemetryEnabled(telemetryEnabled),
		hermes.WithTelemetryHost(telemetryHost),
		hermes.WithTelemetryPort(telemetryPort),
		hermes.WithRestEnabled(restEnabled),
		hermes.WithRestHost(restHost),
		hermes.WithRestPort(restPort),
		hermes.WithModeChannelsEnabled(modeChannelsEnabled),
		hermes.WithModeClientsEnabled(modeClientsEnabled),
		hermes.WithModeClientsMisbehaviour(modeClientsMisbehaviour),
		hermes.WithModeClientsRefresh(modeClientsRefresh),
		hermes.WithModeConnectionsEnabled(modeConnectionsEnabled),
		hermes.WithModePacketsEnabled(modePacketsEnabled),
		hermes.WithModePacketsClearInterval(modePacketsClearInterval),
		hermes.WithModePacketsClearOnStart(modePacketsClearOnStart),
		hermes.WithModePacketsTxConfirmation(modePacketsTxConfirmation),
		hermes.WithAutoRegisterCounterpartyPayee(modePacketsAutoRegisterCounterpartyPayee),
	)

	// Add chain A into the config
	var (
		chainAID       = args[0]
		chainARPCAddr  = args[1]
		chainAGRPCAddr = args[2]

		chainAEventSourceMode, _           = flags.GetString(flagChainAEventSourceMode)
		chainAEventSourceURL, _            = flags.GetString(flagChainAEventSourceURL)
		chainAEventSourceBatchDelay, _     = flags.GetString(flagChainAEventSourceBatchDelay)
		chainARPCTimeout, _                = flags.GetString(flagChainARPCTimeout)
		chainAAccountPrefix, _             = flags.GetString(flagChainAAccountPrefix)
		chainAAddressType, _               = flags.GetString(flagChainAAddressType)
		chainAKeyName, _                   = flags.GetString(flagChainAKeyName)
		chainAStorePrefix, _               = flags.GetString(flagChainAStorePrefix)
		chainADefaultGas, _                = flags.GetUint64(flagChainADefaultGas)
		chainAMaxGas, _                    = flags.GetUint64(flagChainAMaxGas)
		chainAGasPrice, _                  = flags.GetString(flagChainAGasPrice)
		chainAGasMultiplier, _             = flags.GetString(flagChainAGasMultiplier)
		chainAMaxMsgNum, _                 = flags.GetUint64(flagChainAMaxMsgNum)
		chainAMaxTxSize, _                 = flags.GetUint64(flagChainAMaxTxSize)
		chainAClockDrift, _                = flags.GetString(flagChainAClockDrift)
		chainAMaxBlockTime, _              = flags.GetString(flagChainAMaxBlockTime)
		chainATrustingPeriod, _            = flags.GetString(flagChainATrustingPeriod)
		chainATrustThresholdNumerator, _   = flags.GetUint64(flagChainATrustThresholdNumerator)
		chainATrustThresholdDenominator, _ = flags.GetUint64(flagChainATrustThresholdDenominator)
		chainACCVConsumerChain, _          = flags.GetBool(flagChainACCVConsumerChain)
		chainATrustedNode, _               = flags.GetBool(flagChainATrustedNode)
		chainAMemoPrefix, _                = flags.GetString(flagChainAMemoPrefix)
	)

	chainAGasMulti := new(big.Float)
	chainAGasMulti, ok := chainAGasMulti.SetString(chainAGasMultiplier)
	if !ok {
		return nil, errors.Errorf("invalid chain A gas multiplier: %s", chainAGasMultiplier)
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
		optChainA = append(optChainA, hermes.WithChainAccountPrefix(chainAAccountPrefix))
	}
	if chainAAddressType != "" {
		optChainA = append(optChainA, hermes.WithChainAddressType(chainAAddressType))
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
		gasPrice, err := sdk.ParseDecCoin(chainAGasPrice)
		if err != nil {
			return nil, err
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
	if chainACCVConsumerChain {
		optChainA = append(optChainA, hermes.WithChainCCVConsumerChain(chainACCVConsumerChain))
	}
	if chainATrustedNode {
		optChainA = append(optChainA, hermes.WithChainTrustedNode(chainATrustedNode))
	}
	if chainAMemoPrefix != "" {
		optChainA = append(optChainA, hermes.WithChainMemoPrefix(chainAMemoPrefix))
	}

	_, err := c.AddChain(chainAID, chainARPCAddr, chainAGRPCAddr, optChainA...)
	if err != nil {
		return nil, err
	}

	// Add chain B into the config
	var (
		chainBID       = args[3]
		chainBRPCAddr  = args[4]
		chainBGRPCAddr = args[5]

		chainBEventSourceMode, _           = flags.GetString(flagChainBEventSourceMode)
		chainBEventSourceURL, _            = flags.GetString(flagChainBEventSourceURL)
		chainBEventSourceBatchDelay, _     = flags.GetString(flagChainBEventSourceBatchDelay)
		chainBRPCTimeout, _                = flags.GetString(flagChainBRPCTimeout)
		chainBAccountPrefix, _             = flags.GetString(flagChainBAccountPrefix)
		chainBAddressType, _               = flags.GetString(flagChainBAddressType)
		chainBKeyName, _                   = flags.GetString(flagChainBKeyName)
		chainBStorePrefix, _               = flags.GetString(flagChainBStorePrefix)
		chainBDefaultGas, _                = flags.GetUint64(flagChainBDefaultGas)
		chainBMaxGas, _                    = flags.GetUint64(flagChainBMaxGas)
		chainBGasPrice, _                  = flags.GetString(flagChainBGasPrice)
		chainBGasMultiplier, _             = flags.GetString(flagChainBGasMultiplier)
		chainBMaxMsgNum, _                 = flags.GetUint64(flagChainBMaxMsgNum)
		chainBMaxTxSize, _                 = flags.GetUint64(flagChainBMaxTxSize)
		chainBClockDrift, _                = flags.GetString(flagChainBClockDrift)
		chainBMaxBlockTime, _              = flags.GetString(flagChainBMaxBlockTime)
		chainBTrustingPeriod, _            = flags.GetString(flagChainBTrustingPeriod)
		chainBTrustThresholdNumerator, _   = flags.GetUint64(flagChainBTrustThresholdNumerator)
		chainBTrustThresholdDenominator, _ = flags.GetUint64(flagChainBTrustThresholdDenominator)
		chainBCCVConsumerChain, _          = flags.GetBool(flagChainBCCVConsumerChain)
		chainBTrustedNode, _               = flags.GetBool(flagChainBTrustedNode)
		chainBMemoPrefix, _                = flags.GetString(flagChainBMemoPrefix)
	)

	chainBGasMulti := new(big.Float)
	chainBGasMulti, ok = chainBGasMulti.SetString(chainBGasMultiplier)
	if !ok {
		return nil, errors.Errorf("invalid chain B gas multiplier: %s", chainBGasMultiplier)
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
		optChainB = append(optChainB, hermes.WithChainAccountPrefix(chainBAccountPrefix))
	}
	if chainBAddressType != "" {
		optChainB = append(optChainB, hermes.WithChainAddressType(chainBAddressType))
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
		gasPrice, err := sdk.ParseDecCoin(chainBGasPrice)
		if err != nil {
			return nil, err
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
	if chainBCCVConsumerChain {
		optChainB = append(optChainB, hermes.WithChainCCVConsumerChain(chainBCCVConsumerChain))
	}
	if chainBTrustedNode {
		optChainB = append(optChainB, hermes.WithChainTrustedNode(chainBTrustedNode))
	}
	if chainBMemoPrefix != "" {
		optChainB = append(optChainB, hermes.WithChainMemoPrefix(chainBMemoPrefix))
	}

	_, err = c.AddChain(chainBID, chainBRPCAddr, chainBGRPCAddr, optChainB...)
	if err != nil {
		return nil, err
	}

	return c, nil
}