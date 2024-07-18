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

	"github.com/ignite/apps/hermes/pkg/hermes"
)

func ConfigureHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	var (
		args  = cmd.Args
		flags = cmd.Flags
	)

	session := cliui.New(cliui.StartSpinnerWithText("Generating Hermes config"))
	defer session.End()

	var (
		chainAID = args[0]
		chainBID = args[3]

		generateWallets, _ = getFlag[bool](flags, flagGenerateWallets)
		overwriteConfig, _ = getFlag[bool](flags, flagOverwriteConfig)
		chainAPortID, _    = getFlag[string](flags, flagChainAPortID)
		chainAFaucet, _    = getFlag[string](flags, flagChainAFaucet)
		chainBPortID, _    = getFlag[string](flags, flagChainBPortID)
		chainBFaucet, _    = getFlag[string](flags, flagChainBFaucet)
		channelVersion, _  = getFlag[string](flags, flagChannelVersion)
		customCfg          = getConfig(flags)
	)

	var (
		hermesCfg *hermes.Config
		err       error
	)
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
		session.StopSpinner()
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

	session.StopSpinner()
	_ = session.Println(color.Green.Sprintf("Hermes config created at %s", cfgPath))

	h, err := hermes.New()
	if err != nil {
		return err
	}
	defer h.Cleanup()

	session.StartSpinner(fmt.Sprintf("Verifying chain A (%s) keys", chainAID))

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

	session.StartSpinner(fmt.Sprintf("Verifying chain B (%s) keys", chainBID))

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

	// create client A
	session.StartSpinner("Creating client A")
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

	session.StopSpinner()
	_ = session.Println(color.Green.Sprintf(
		"Client '%s' created (%s -> %s)",
		clientAResult.CreateClient.ClientID,
		chainAID,
		chainBID,
	))

	// create client B
	session.StartSpinner("Creating client B")
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

	session.StopSpinner()
	_ = session.Println(color.Green.Sprintf(
		"Client %s' created (%s -> %s)",
		clientBResult.CreateClient.ClientID,
		chainBID,
		chainAID,
	))

	// create connection
	session.StartSpinner("Creating connection")
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

	session.StopSpinner()
	_ = session.Println(color.Green.Sprintf(
		"Connection '%s (%s) <-> %s (%s)' created",
		chainAID,
		connection.ASide.ConnectionID,
		chainBID,
		connection.BSide.ConnectionID,
	))

	// create and query channel
	session.StartSpinner("Creating channel")
	var (
		bufChannel     = bytes.Buffer{}
		channel        = hermes.ConnectionResult{}
		createChanOpts = []hermes.Option{
			hermes.WithConfigFile(cfgPath),
			hermes.WithStdOut(&bufChannel),
			hermes.WithJSONOutput(),
		}
	)
	if channelVersion != "" {
		createChanOpts = append(createChanOpts, hermes.WithFlags(hermes.Flags{flagChannelVersion: channelVersion}))
	}

	if err := h.CreateChannel(
		ctx,
		chainAID,
		connection.ASide.ConnectionID,
		chainAPortID,
		chainBPortID,
		createChanOpts...,
	); err != nil {
		return err
	}
	if err := hermes.UnmarshalResult(bufChannel.Bytes(), &channel); err != nil {
		return err
	}

	session.StopSpinner()
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

	session.StartSpinner(fmt.Sprintf("verifying %s balance", chainAddr))

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
		session.StartSpinner(fmt.Sprintf("requesting faucet balance for %s", chainAddr))
		newBalance, err := chain.TryRetrieve(ctx, chainAddr, faucetAddr)
		if err != nil {
			return err
		}

		session.StopSpinner()
		_ = session.Printf(
			"%s %s\n",
			color.Green.Sprint("New balance from faucet:"),
			color.Yellow.Sprint(newBalance.String()),
		)
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
			session.StopSpinner()
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

			session.StopSpinner()
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

		session.StopSpinner()
		_ = session.Println(color.Green.Sprintf("Chain %s key created", chainID))

		goto GetKey
	}

	session.StopSpinner()
	_ = session.Printf(
		"%s %s\n",
		color.Green.Sprintf("Chain %s relayer wallet:", chainID),
		color.Yellow.Sprint(keysChainResult.Wallet.Account),
	)

	return keysChainResult.Wallet.Account, nil
}

// newHermesConfig create a new hermes config based in the cmd args.
func newHermesConfig(flags []*plugin.Flag, args []string, customCfg string) (*hermes.Config, error) {
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
		telemetryEnabled, _                         = getFlag[bool](flags, flagTelemetryEnabled)
		telemetryHost, _                            = getFlag[string](flags, flagTelemetryHost)
		telemetryPort, _                            = getFlag[uint64](flags, flagTelemetryPort)
		restEnabled, _                              = getFlag[bool](flags, flagRestEnabled)
		restHost, _                                 = getFlag[string](flags, flagRestHost)
		restPort, _                                 = getFlag[uint64](flags, flagRestPort)
		modeChannelsEnabled, _                      = getFlag[bool](flags, flagModeChannelsEnabled)
		modeClientsEnabled, _                       = getFlag[bool](flags, flagModeClientsEnabled)
		modeClientsMisbehaviour, _                  = getFlag[bool](flags, flagModeClientsMisbehaviour)
		modeClientsRefresh, _                       = getFlag[bool](flags, flagModeClientsRefresh)
		modeConnectionsEnabled, _                   = getFlag[bool](flags, flagModeConnectionsEnabled)
		modePacketsEnabled, _                       = getFlag[bool](flags, flagModePacketsEnabled)
		modePacketsClearInterval, _                 = getFlag[uint64](flags, flagModePacketsClearInterval)
		modePacketsClearOnStart, _                  = getFlag[bool](flags, flagModePacketsClearOnStart)
		modePacketsTxConfirmation, _                = getFlag[bool](flags, flagModePacketsTxConfirmation)
		modePacketsAutoRegisterCounterpartyPayee, _ = getFlag[bool](flags, flagAutoRegisterCounterpartyPayee)
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

		chainAEventSourceMode, _           = getFlag[string](flags, flagChainAEventSourceMode)
		chainAEventSourceURL, _            = getFlag[string](flags, flagChainAEventSourceURL)
		chainAEventSourceBatchDelay, _     = getFlag[string](flags, flagChainAEventSourceBatchDelay)
		chainARPCTimeout, _                = getFlag[string](flags, flagChainARPCTimeout)
		chainAAccountPrefix, _             = getFlag[string](flags, flagChainAAccountPrefix)
		chainAAddressType, _               = getFlag[string](flags, flagChainAAddressType)
		chainAKeyName, _                   = getFlag[string](flags, flagChainAKeyName)
		chainAKeyStoreType, _              = getFlag[string](flags, flagChainAKeyStoreType)
		chainAStorePrefix, _               = getFlag[string](flags, flagChainAStorePrefix)
		chainADefaultGas, _                = getFlag[uint64](flags, flagChainADefaultGas)
		chainAMaxGas, _                    = getFlag[uint64](flags, flagChainAMaxGas)
		chainAGasPrice, _                  = getFlag[string](flags, flagChainAGasPrice)
		chainAGasMultiplier, _             = getFlag[string](flags, flagChainAGasMultiplier)
		chainAMaxMsgNum, _                 = getFlag[uint64](flags, flagChainAMaxMsgNum)
		chainAMaxTxSize, _                 = getFlag[uint64](flags, flagChainAMaxTxSize)
		chainAClockDrift, _                = getFlag[string](flags, flagChainAClockDrift)
		chainAMaxBlockTime, _              = getFlag[string](flags, flagChainAMaxBlockTime)
		chainATrustingPeriod, _            = getFlag[string](flags, flagChainATrustingPeriod)
		chainATrustThresholdNumerator, _   = getFlag[uint64](flags, flagChainATrustThresholdNumerator)
		chainATrustThresholdDenominator, _ = getFlag[uint64](flags, flagChainATrustThresholdDenominator)
		chainACCVConsumerChain, _          = getFlag[bool](flags, flagChainACCVConsumerChain)
		chainATrustedNode, _               = getFlag[bool](flags, flagChainATrustedNode)
		chainAMemoPrefix, _                = getFlag[string](flags, flagChainAMemoPrefix)
		chainAType, _                      = getFlag[string](flags, flagChainAType)
		chainASequentialBatchTx, _         = getFlag[bool](flags, flagChainASequentialBatchTx)
	)

	fmt.Println("aefaefaefeaf _ " + chainAGasMultiplier)
	chainAGasMulti := new(big.Float)
	chainAGasMulti, ok := chainAGasMulti.SetString(chainAGasMultiplier)
	if !ok {
		return nil, errors.Errorf("invalid chain A gas multiplier: %s", chainAGasMultiplier)
	}

	optChainA := []hermes.ChainOption{
		hermes.WithChainTrustThreshold(chainATrustThresholdNumerator, chainATrustThresholdDenominator),
		hermes.WithChainGasMultiplier(chainAGasMulti),
		hermes.WithChainCCVConsumerChain(chainACCVConsumerChain),
		hermes.WithChainTrustedNode(chainATrustedNode),
		hermes.WithChainSequentialBatchTx(chainASequentialBatchTx),
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
	if chainAKeyStoreType != "" {
		optChainA = append(optChainA, hermes.WithChainKeyStoreType(chainAKeyStoreType))
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
	if chainAMemoPrefix != "" {
		optChainA = append(optChainA, hermes.WithChainMemoPrefix(chainAMemoPrefix))
	}
	if chainAType != "" {
		optChainA = append(optChainA, hermes.WithChainType(chainAType))
	}

	if _, err := c.AddChain(chainAID, chainARPCAddr, chainAGRPCAddr, optChainA...); err != nil {
		return nil, err
	}

	// Add chain B into the config
	var (
		chainBID       = args[3]
		chainBRPCAddr  = args[4]
		chainBGRPCAddr = args[5]

		chainBEventSourceMode, _           = getFlag[string](flags, flagChainBEventSourceMode)
		chainBEventSourceURL, _            = getFlag[string](flags, flagChainBEventSourceURL)
		chainBEventSourceBatchDelay, _     = getFlag[string](flags, flagChainBEventSourceBatchDelay)
		chainBRPCTimeout, _                = getFlag[string](flags, flagChainBRPCTimeout)
		chainBAccountPrefix, _             = getFlag[string](flags, flagChainBAccountPrefix)
		chainBAddressType, _               = getFlag[string](flags, flagChainBAddressType)
		chainBKeyName, _                   = getFlag[string](flags, flagChainBKeyName)
		chainBKeyStoreType, _              = getFlag[string](flags, flagChainBKeyStoreType)
		chainBStorePrefix, _               = getFlag[string](flags, flagChainBStorePrefix)
		chainBDefaultGas, _                = getFlag[uint64](flags, flagChainBDefaultGas)
		chainBMaxGas, _                    = getFlag[uint64](flags, flagChainBMaxGas)
		chainBGasPrice, _                  = getFlag[string](flags, flagChainBGasPrice)
		chainBGasMultiplier, _             = getFlag[string](flags, flagChainBGasMultiplier)
		chainBMaxMsgNum, _                 = getFlag[uint64](flags, flagChainBMaxMsgNum)
		chainBMaxTxSize, _                 = getFlag[uint64](flags, flagChainBMaxTxSize)
		chainBClockDrift, _                = getFlag[string](flags, flagChainBClockDrift)
		chainBMaxBlockTime, _              = getFlag[string](flags, flagChainBMaxBlockTime)
		chainBTrustingPeriod, _            = getFlag[string](flags, flagChainBTrustingPeriod)
		chainBTrustThresholdNumerator, _   = getFlag[uint64](flags, flagChainBTrustThresholdNumerator)
		chainBTrustThresholdDenominator, _ = getFlag[uint64](flags, flagChainBTrustThresholdDenominator)
		chainBCCVConsumerChain, _          = getFlag[bool](flags, flagChainBCCVConsumerChain)
		chainBTrustedNode, _               = getFlag[bool](flags, flagChainBTrustedNode)
		chainBMemoPrefix, _                = getFlag[string](flags, flagChainBMemoPrefix)
		chainBType, _                      = getFlag[string](flags, flagChainBType)
		chainBSequentialBatchTx, _         = getFlag[bool](flags, flagChainBSequentialBatchTx)
	)

	chainBGasMulti := new(big.Float)
	chainBGasMulti, ok = chainBGasMulti.SetString(chainBGasMultiplier)
	if !ok {
		return nil, errors.Errorf("invalid chain B gas multiplier: %s", chainBGasMultiplier)
	}

	optChainB := []hermes.ChainOption{
		hermes.WithChainTrustThreshold(chainBTrustThresholdNumerator, chainBTrustThresholdDenominator),
		hermes.WithChainGasMultiplier(chainBGasMulti),
		hermes.WithChainCCVConsumerChain(chainBCCVConsumerChain),
		hermes.WithChainTrustedNode(chainBTrustedNode),
		hermes.WithChainSequentialBatchTx(chainBSequentialBatchTx),
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
	if chainBKeyStoreType != "" {
		optChainB = append(optChainB, hermes.WithChainKeyStoreType(chainBKeyStoreType))
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
	if chainBMemoPrefix != "" {
		optChainB = append(optChainB, hermes.WithChainMemoPrefix(chainBMemoPrefix))
	}
	if chainBType != "" {
		optChainB = append(optChainB, hermes.WithChainType(chainBType))
	}

	if _, err := c.AddChain(chainBID, chainBRPCAddr, chainBGRPCAddr, optChainB...); err != nil {
		return nil, err
	}

	return c, nil
}
