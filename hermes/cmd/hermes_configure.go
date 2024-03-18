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
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

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
	c.Flags().Bool(flagChainACCVConsumerChain, false, "only specify true if the chain A is a CCV consumer")
	c.Flags().Bool(flagChainBCCVConsumerChain, false, "only specify true if the chain B is a CCV consumer")
	c.Flags().String(flagChainAEventSourceURL, "", "WS event source url of the chain A")
	c.Flags().String(flagChainBEventSourceURL, "", "WS event source url of the chain B")
	c.Flags().String(flagChainAEventSourceMode, "push", "WS event source mode of the chain A (event source url should be set to use this flag)")
	c.Flags().String(flagChainBEventSourceMode, "push", "WS event source mode of the chain B (event source url should be set to use this flag)")
	c.Flags().String(flagChainAEventSourceBatchDelay, "500ms", "WS event source batch delay time of the chain A (event source url should be set to use this flag)")
	c.Flags().String(flagChainBEventSourceBatchDelay, "500ms", "WS event source batch delay time of the chain B (event source url should be set to use this flag)")
	c.Flags().String(flagChainARPCTimeout, "10s", "RPC timeout of the chain A")
	c.Flags().String(flagChainBRPCTimeout, "10s", "RPC timeout of the chain B")
	c.Flags().Bool(flagChainATrustedNode, false, "enable trusted node on the chain A")
	c.Flags().Bool(flagChainBTrustedNode, false, "enable trusted node on the chain B")
	c.Flags().String(flagChainAAccountPrefix, "cosmos", "account prefix of the chain A")
	c.Flags().String(flagChainBAccountPrefix, "cosmos", "account prefix of the chain B")
	c.Flags().String(flagChainAKeyName, "wallet", "hermes account name of the chain A")
	c.Flags().String(flagChainBKeyName, "wallet", "hermes account name of the chain B")
	c.Flags().String(flagChainAAddressType, "cosmos", "address type of the chain A")
	c.Flags().String(flagChainBAddressType, "cosmos", "address type of the chain B")
	c.Flags().String(flagChainAStorePrefix, "ibc", "store prefix of the chain A")
	c.Flags().String(flagChainBStorePrefix, "ibc", "store prefix of the chain B")
	c.Flags().Uint64(flagChainADefaultGas, 100000, "default gas used for transactions on chain A")
	c.Flags().Uint64(flagChainBDefaultGas, 100000, "default gas used for transactions on chain B")
	c.Flags().Uint64(flagChainAMaxGas, 400000, "max gas used for transactions on chain A")
	c.Flags().Uint64(flagChainBMaxGas, 400000, "max gas used for transactions on chain B")
	c.Flags().String(flagChainAGasPrice, "0.025stake", "gas price used for transactions on chain A")
	c.Flags().String(flagChainBGasPrice, "0.025stake", "gas price used for transactions on chain B")
	c.Flags().String(flagChainAGasMultiplier, "1.1", "gas multiplier used for transactions on chain A")
	c.Flags().String(flagChainBGasMultiplier, "1.1", "gas multiplier used for transactions on chain B")
	c.Flags().Uint64(flagChainAMaxMsgNum, 30, "max message number used for transactions on chain A")
	c.Flags().Uint64(flagChainBMaxMsgNum, 30, "max message number used for transactions on chain B")
	c.Flags().Uint64(flagChainAMaxTxSize, 2097152, "max transaction size on chain A")
	c.Flags().Uint64(flagChainBMaxTxSize, 2097152, "max transaction size on chain B")
	c.Flags().String(flagChainAClockDrift, "5s", "clock drift of the chain A")
	c.Flags().String(flagChainBClockDrift, "5s", "clock drift of the chain B")
	c.Flags().String(flagChainAMaxBlockTime, "30s", "maximum block time of the chain A")
	c.Flags().String(flagChainBMaxBlockTime, "30s", "maximum block time of the chain B")
	c.Flags().String(flagChainATrustingPeriod, "14days", "trusting period of the chain A")
	c.Flags().String(flagChainBTrustingPeriod, "14days", "trusting period of the chain B")
	c.Flags().Uint64(flagChainATrustThresholdNumerator, 2, "trusting threshold numerator of the chain A")
	c.Flags().Uint64(flagChainBTrustThresholdNumerator, 2, "trusting threshold numerator of the chain B")
	c.Flags().Uint64(flagChainATrustThresholdDenominator, 3, "trusting threshold denominator of the chain A")
	c.Flags().Uint64(flagChainBTrustThresholdDenominator, 3, "trusting threshold denominator of the chain B")
	c.Flags().String(flagChainAMemoPrefix, "", "memo prefix of the chain A")
	c.Flags().String(flagChainBMemoPrefix, "", "memo prefix of the chain B")
	c.Flags().String(flagChainAFaucet, "", "faucet URL of the chain A")
	c.Flags().String(flagChainBFaucet, "", "faucet URL of the chain B")

	c.Flags().Bool(flagTelemetryEnabled, false, "enable hermes telemetry")
	c.Flags().String(flagTelemetryHost, "127.0.0.1", "hermes telemetry host")
	c.Flags().Uint64(flagTelemetryPort, 3001, "hermes telemetry port")
	c.Flags().Bool(flagRestEnabled, false, "enable hermes rest")
	c.Flags().String(flagRestHost, "127.0.0.1", "hermes rest host")
	c.Flags().Uint64(flagRestPort, 3000, "hermes rest port")
	c.Flags().Bool(flagModeChannelsEnabled, true, "enable hermes channels")
	c.Flags().Bool(flagModeClientsEnabled, true, "enable hermes clients")
	c.Flags().Bool(flagModeClientsMisbehaviour, true, "enable hermes clients misbehaviour")
	c.Flags().Bool(flagModeClientsRefresh, true, "enable hermes client refresh time")
	c.Flags().Bool(flagModeConnectionsEnabled, true, "enable hermes connections")
	c.Flags().Bool(flagModePacketsEnabled, true, "enable hermes packets")
	c.Flags().Uint64(flagModePacketsClearInterval, 100, "hermes packet clear interval")
	c.Flags().Bool(flagModePacketsClearOnStart, true, "enable hermes packets clear on start")
	c.Flags().Bool(flagModePacketsTxConfirmation, true, "hermes packet transaction confirmation")
	c.Flags().Bool(flagAutoRegisterCounterpartyPayee, false, "auto register the counterparty payee on a destination chain to the relayer's address on the source chain")
	c.Flags().Bool(flagGenerateWallets, false, "automatically generate wallets if they do not exist")
	c.Flags().Bool(flagOverwriteConfig, false, "overwrite the current config if it already exists")

	return c
}

func hermesConfigureHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	session.StartSpinner("Generating Hermes config")

	var (
		chainAID = args[0]
		chainBID = args[3]

		generateWallets, _ = cmd.Flags().GetBool(flagGenerateWallets)
		overwriteConfig, _ = cmd.Flags().GetBool(flagOverwriteConfig)
		chainAPortID, _    = cmd.Flags().GetString(flagChainAPortID)
		chainAFaucet, _    = cmd.Flags().GetString(flagChainAFaucet)
		chainBPortID, _    = cmd.Flags().GetString(flagChainBPortID)
		chainBFaucet, _    = cmd.Flags().GetString(flagChainBFaucet)
		customCfg          = getConfig(cmd)
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
		hermesCfg, err = newHermesConfig(cmd, args, customCfg)
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
		cmd.Context(),
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
		cmd.Context(),
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
		cmd.Context(),
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
		cmd.Context(),
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
		cmd.Context(),
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
		cmd.Context(),
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
func newHermesConfig(cmd *cobra.Command, args []string, customCfg string) (*hermes.Config, error) {
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
		telemetryEnabled, _                         = cmd.Flags().GetBool(flagTelemetryEnabled)
		telemetryHost, _                            = cmd.Flags().GetString(flagTelemetryHost)
		telemetryPort, _                            = cmd.Flags().GetUint64(flagTelemetryPort)
		restEnabled, _                              = cmd.Flags().GetBool(flagRestEnabled)
		restHost, _                                 = cmd.Flags().GetString(flagRestHost)
		restPort, _                                 = cmd.Flags().GetUint64(flagRestPort)
		modeChannelsEnabled, _                      = cmd.Flags().GetBool(flagModeChannelsEnabled)
		modeClientsEnabled, _                       = cmd.Flags().GetBool(flagModeClientsEnabled)
		modeClientsMisbehaviour, _                  = cmd.Flags().GetBool(flagModeClientsMisbehaviour)
		modeClientsRefresh, _                       = cmd.Flags().GetBool(flagModeClientsRefresh)
		modeConnectionsEnabled, _                   = cmd.Flags().GetBool(flagModeConnectionsEnabled)
		modePacketsEnabled, _                       = cmd.Flags().GetBool(flagModePacketsEnabled)
		modePacketsClearInterval, _                 = cmd.Flags().GetUint64(flagModePacketsClearInterval)
		modePacketsClearOnStart, _                  = cmd.Flags().GetBool(flagModePacketsClearOnStart)
		modePacketsTxConfirmation, _                = cmd.Flags().GetBool(flagModePacketsTxConfirmation)
		modePacketsAutoRegisterCounterpartyPayee, _ = cmd.Flags().GetBool(flagAutoRegisterCounterpartyPayee)
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

		chainAEventSourceMode, _           = cmd.Flags().GetString(flagChainAEventSourceMode)
		chainAEventSourceURL, _            = cmd.Flags().GetString(flagChainAEventSourceURL)
		chainAEventSourceBatchDelay, _     = cmd.Flags().GetString(flagChainAEventSourceBatchDelay)
		chainARPCTimeout, _                = cmd.Flags().GetString(flagChainARPCTimeout)
		chainAAccountPrefix, _             = cmd.Flags().GetString(flagChainAAccountPrefix)
		chainAAddressType, _               = cmd.Flags().GetString(flagChainAAddressType)
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
		chainACCVConsumerChain, _          = cmd.Flags().GetBool(flagChainACCVConsumerChain)
		chainATrustedNode, _               = cmd.Flags().GetBool(flagChainATrustedNode)
		chainAMemoPrefix, _                = cmd.Flags().GetString(flagChainAMemoPrefix)
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

		chainBEventSourceMode, _           = cmd.Flags().GetString(flagChainBEventSourceMode)
		chainBEventSourceURL, _            = cmd.Flags().GetString(flagChainBEventSourceURL)
		chainBEventSourceBatchDelay, _     = cmd.Flags().GetString(flagChainBEventSourceBatchDelay)
		chainBRPCTimeout, _                = cmd.Flags().GetString(flagChainBRPCTimeout)
		chainBAccountPrefix, _             = cmd.Flags().GetString(flagChainBAccountPrefix)
		chainBAddressType, _               = cmd.Flags().GetString(flagChainBAddressType)
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
		chainBCCVConsumerChain, _          = cmd.Flags().GetBool(flagChainBCCVConsumerChain)
		chainBTrustedNode, _               = cmd.Flags().GetBool(flagChainBTrustedNode)
		chainBMemoPrefix, _                = cmd.Flags().GetString(flagChainBMemoPrefix)
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
