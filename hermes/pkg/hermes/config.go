package hermes

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/v28/ignite/pkg/cosmosfaucet"
	"github.com/pelletier/go-toml/v2"
)

const (
	// ConfigNameSeparator config file chain name separator.
	ConfigNameSeparator = "_"
)

type (
	// Config represents the Hermes config struct.
	Config struct {
		Chains    Chains    `toml:"chains" json:"chains"`
		Global    Global    `toml:"global" json:"global"`
		Telemetry Telemetry `toml:"telemetry" json:"telemetry"`
		Mode      Mode      `toml:"mode" json:"mode"`
	}

	// Chains represents a list of chains.
	Chains []Chain

	// Chain represents the chain into the Hermes config struct.
	Chain struct {
		ID             string         `toml:"id" json:"id"`
		RPCAddr        string         `toml:"rpc_addr" json:"rpc_addr"`
		GRPCAddr       string         `toml:"grpc_addr" json:"grpc_addr"`
		EventSource    EventSource    `toml:"event_source,inline" json:"event_source"`
		RPCTimeout     string         `toml:"rpc_timeout" json:"rpc_timeout"`
		AccountPrefix  string         `toml:"account_prefix" json:"account_prefix"`
		KeyName        string         `toml:"key_name" json:"key_name"`
		StorePrefix    string         `toml:"store_prefix" json:"store_prefix"`
		DefaultGas     uint64         `toml:"default_gas" json:"default_gas"`
		MaxGas         uint64         `toml:"max_gas" json:"max_gas"`
		GasPrice       GasPrice       `toml:"gas_price,inline" json:"gas_price"`
		GasMultiplier  float64        `toml:"gas_multiplier" json:"gas_multiplier"`
		MaxMsgNum      uint64         `toml:"max_msg_num" json:"max_msg_num"`
		MaxTxSize      uint64         `toml:"max_tx_size" json:"max_tx_size"`
		ClockDrift     string         `toml:"clock_drift" json:"clock_drift"`
		MaxBlockTime   string         `toml:"max_block_time" json:"max_block_time"`
		TrustingPeriod string         `toml:"trusting_period" json:"trusting_period"`
		TrustThreshold TrustThreshold `toml:"trust_threshold,inline" json:"trust_threshold"`
		AddressType    AddressType    `toml:"address_type,inline" json:"address_type"`
	}

	// EventSource represents the chain event source into the Hermes config struct.
	EventSource struct {
		BatchDelay string `toml:"batch_delay" json:"batch_delay"`
		Mode       string `toml:"mode" json:"mode"`
		URL        string `toml:"url" json:"url"`
	}

	// GasPrice represents the chain gas price into the Hermes config struct.
	GasPrice struct {
		Denom string  `toml:"denom" json:"denom"`
		Price float64 `toml:"price" json:"price"`
	}

	// TrustThreshold represents the chain trust threshold into the Hermes config struct.
	TrustThreshold struct {
		Denominator string `toml:"denominator" json:"denominator"`
		Numerator   string `toml:"numerator" json:"numerator"`
	}

	// AddressType represents the chain address type into the Hermes config struct.
	AddressType struct {
		Derivation string `toml:"derivation" json:"derivation"`
	}

	// Global represents the global values into the Hermes config struct.
	Global struct {
		LogLevel string `toml:"log_level" json:"log_level"`
	}

	// Telemetry represents the telemetry into the Hermes config struct.
	Telemetry struct {
		Enabled bool   `toml:"enabled" json:"enabled"`
		Host    string `toml:"host" json:"host"`
		Port    uint64 `toml:"port" json:"port"`
	}

	// Mode represents the mode into the Hermes config struct.
	Mode struct {
		Channels    Channels    `toml:"channels" json:"channels"`
		Clients     Clients     `toml:"clients" json:"clients"`
		Connections Connections `toml:"connections" json:"connections"`
		Packets     Packets     `toml:"packets" json:"packets"`
	}

	// Channels represents the mode channels into the Hermes config struct.
	Channels struct {
		Enabled bool `toml:"enabled" json:"enabled"`
	}

	// Clients represents the mode clients into the Hermes config struct.
	Clients struct {
		Enabled      bool `toml:"enabled" json:"enabled"`
		Misbehaviour bool `toml:"misbehaviour" json:"misbehaviour"`
		Refresh      bool `toml:"refresh" json:"refresh"`
	}

	// Connections represents the mode connections into the Hermes config struct.
	Connections struct {
		Enabled bool `toml:"enabled" json:"enabled"`
	}

	// Packets represents the mode packets into the Hermes config struct.
	Packets struct {
		ClearInterval  uint64 `toml:"clear_interval" json:"clear_interval"`
		ClearOnStart   bool   `toml:"clear_on_start" json:"clear_on_start"`
		Enabled        bool   `toml:"enabled" json:"enabled"`
		TxConfirmation bool   `toml:"tx_confirmation" json:"tx_confirmation"`
	}

	// ChainOption configures chain hermes configs.
	ChainOption func(*Chain)
	// ConfigOption configures hermes configs.
	ConfigOption func(*Config)
)

// Get returns the chain by chain id.
func (c Chains) Get(chainID string) (Chain, error) {
	for _, chain := range c {
		if chain.ID == chainID {
			return chain, nil
		}
	}
	return Chain{}, fmt.Errorf("chain %s not exist", chainID)
}

// Remove delete the Hermes config file.
func (c *Config) Remove() error {
	configPath, err := c.ConfigPath()
	if err != nil {
		return err
	}
	return os.RemoveAll(configPath)
}

// Save create and save a new Hermes config file.
func (c *Config) Save() error {
	configPath, err := c.ConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return err
	}

	file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	return toml.NewEncoder(file).Encode(c)
}

// ConfigName returns the config file name based on the chains inside the config file.
func (c *Config) ConfigName() (string, error) {
	if len(c.Chains) < 2 {
		return "", errors.New("cannot create a config file without unless two chains")
	}
	names := make([]string, 0)
	for _, chain := range c.Chains {
		names = append(names, chain.ID)
	}
	return strings.Join(names, ConfigNameSeparator), nil
}

// ConfigPath return the config file path.
func (c *Config) ConfigPath() (string, error) {
	cfgName, err := c.ConfigName()
	if err != nil {
		return "", err
	}
	return ConfigPath(cfgName)
}

// ConfigPath generates a config file path.
func ConfigPath(cfgName string) (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(
		userHomeDir,
		".ignite",
		"relayer",
		"hermes",
		cfgName,
	), nil
}

// LoadConfig loads a config from the path.
func LoadConfig(cfgPath string) (*Config, error) {
	cfgBytes, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}
	var cfg *Config
	return cfg, toml.Unmarshal(cfgBytes, cfg)
}

// DefaultConfigPath returns the default Hermes config path.
func DefaultConfigPath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userHomeDir, ".hermes", "config.toml"), nil
}

// WithTelemetryEnabled set telemetry enable into the Hermes config.
func WithTelemetryEnabled(enabled bool) ConfigOption {
	return func(c *Config) {
		c.Telemetry.Enabled = enabled
	}
}

// WithTelemetryHost set TelemetryHost into the Hermes config.
func WithTelemetryHost(host string) ConfigOption {
	return func(c *Config) {
		c.Telemetry.Host = host
	}
}

// WithTelemetryPort set TelemetryPort into the Hermes config.
func WithTelemetryPort(port uint64) ConfigOption {
	return func(c *Config) {
		c.Telemetry.Port = port
	}
}

// WithModeChannelsEnabled set ModeChannelsEnabled into the Hermes config.
func WithModeChannelsEnabled(enabled bool) ConfigOption {
	return func(c *Config) {
		c.Mode.Channels.Enabled = enabled
	}
}

// WithModeClientsEnabled set ModeClientsEnabled into the Hermes config.
func WithModeClientsEnabled(enabled bool) ConfigOption {
	return func(c *Config) {
		c.Mode.Clients.Enabled = enabled
	}
}

// WithModeClientsMisbehaviour set ModeClientsMisbehaviour into the Hermes config.
func WithModeClientsMisbehaviour(misbehaviour bool) ConfigOption {
	return func(c *Config) {
		c.Mode.Clients.Misbehaviour = misbehaviour
	}
}

// WithModeClientsRefresh set ModeClientsRefresh into the Hermes config.
func WithModeClientsRefresh(refresh bool) ConfigOption {
	return func(c *Config) {
		c.Mode.Clients.Refresh = refresh
	}
}

// WithModeConnectionsEnabled set ModeConnectionsEnabled into the Hermes config.
func WithModeConnectionsEnabled(enabled bool) ConfigOption {
	return func(c *Config) {
		c.Mode.Connections.Enabled = enabled
	}
}

// WithModePacketsEnabled set ModePacketsEnabled into the Hermes config.
func WithModePacketsEnabled(enabled bool) ConfigOption {
	return func(c *Config) {
		c.Mode.Packets.Enabled = enabled
	}
}

// WithModePacketsClearInterval set ModePacketsClearInterval into the Hermes config.
func WithModePacketsClearInterval(clearInterval uint64) ConfigOption {
	return func(c *Config) {
		c.Mode.Packets.ClearInterval = clearInterval
	}
}

// WithModePacketsClearOnStart set ModePacketsClearOnStart into the Hermes config.
func WithModePacketsClearOnStart(clearOnStart bool) ConfigOption {
	return func(c *Config) {
		c.Mode.Packets.ClearOnStart = clearOnStart
	}
}

// WithModePacketsTxConfirmation set ModePacketsTxConfirmation into the Hermes config.
func WithModePacketsTxConfirmation(txConfirmation bool) ConfigOption {
	return func(c *Config) {
		c.Mode.Packets.TxConfirmation = txConfirmation
	}
}

// DefaultConfig returns a default configuration struct for Hermes.
func DefaultConfig(options ...ConfigOption) *Config {
	cfg := &Config{
		Chains: []Chain{},
		Global: Global{
			LogLevel: "error",
		},
		Mode: Mode{
			Channels: Channels{
				Enabled: true,
			},
			Clients: Clients{
				Enabled:      true,
				Misbehaviour: true,
				Refresh:      true,
			},
			Connections: Connections{
				Enabled: true,
			},
			Packets: Packets{
				ClearInterval:  100,
				ClearOnStart:   true,
				Enabled:        true,
				TxConfirmation: true,
			},
		},
		Telemetry: Telemetry{
			Enabled: false,
			Host:    "127.0.0.1",
			Port:    3001,
		},
	}
	for _, o := range options {
		o(cfg)
	}
	return cfg
}

// WithChainEventSource set event source into the chain config.
func WithChainEventSource(mode, url, batchDelay string) ChainOption {
	return func(c *Chain) {
		c.EventSource = EventSource{
			BatchDelay: batchDelay,
			Mode:       mode,
			URL:        url,
		}
	}
}

// WithChainRPCTimeout set the chain rpc timeout into the Hermes config.
func WithChainRPCTimeout(timeout string) ChainOption {
	return func(c *Chain) {
		c.RPCTimeout = timeout
	}
}

// WithChainAccountPrefix set the chain account prefix into the Hermes config.
func WithChainAccountPrefix(prefix string) ChainOption {
	return func(c *Chain) {
		c.AccountPrefix = prefix
	}
}

// WithChainKeyName set the chain key name into the Hermes config.
func WithChainKeyName(key string) ChainOption {
	return func(c *Chain) {
		c.KeyName = key
	}
}

// WithChainStorePrefix set the chain store prefix into the Hermes config.
func WithChainStorePrefix(prefix string) ChainOption {
	return func(c *Chain) {
		c.StorePrefix = prefix
	}
}

// WithChainDefaultGas set the chain default gas into the Hermes config.
func WithChainDefaultGas(defaultGas uint64) ChainOption {
	return func(c *Chain) {
		c.DefaultGas = defaultGas
	}
}

// WithChainMaxGas set the chain max gas into the Hermes config.
func WithChainMaxGas(maxGas uint64) ChainOption {
	return func(c *Chain) {
		c.MaxGas = maxGas
	}
}

// WithChainGasPrice set the chain gas price into the Hermes config.
func WithChainGasPrice(price sdk.Coin) ChainOption {
	return func(c *Chain) {
		f, _ := price.Amount.BigInt().Float64()
		c.GasPrice = GasPrice{
			Denom: price.Denom,
			Price: f,
		}
	}
}

// WithChainGasMultiplier set the chain gas multiplier into the Hermes config.
func WithChainGasMultiplier(gasMultiplier *big.Float) ChainOption {
	return func(c *Chain) {
		c.GasMultiplier, _ = gasMultiplier.Float64()
	}
}

// WithChainMaxMsgNum set the chain max mesage number into the Hermes config.
func WithChainMaxMsgNum(maxMsg uint64) ChainOption {
	return func(c *Chain) {
		c.MaxMsgNum = maxMsg
	}
}

// WithChainMaxTxSize set the chain maximum transaction size into the Hermes config.
func WithChainMaxTxSize(size uint64) ChainOption {
	return func(c *Chain) {
		c.MaxTxSize = size
	}
}

// WithChainClockDrift set the chain clock drift into the Hermes config.
func WithChainClockDrift(clock string) ChainOption {
	return func(c *Chain) {
		c.ClockDrift = clock
	}
}

// WithChainMaxBlockTime set the chain block time into the Hermes config.
func WithChainMaxBlockTime(maxBlockTime string) ChainOption {
	return func(c *Chain) {
		c.MaxBlockTime = maxBlockTime
	}
}

// WithChainTrustingPeriod set the chain trusting period into the Hermes config.
func WithChainTrustingPeriod(trustingPeriod string) ChainOption {
	return func(c *Chain) {
		c.TrustingPeriod = trustingPeriod
	}
}

// WithChainTrustThreshold set the chain trust threshold into the Hermes config.
func WithChainTrustThreshold(numerator, denominator uint64) ChainOption {
	return func(c *Chain) {
		c.TrustThreshold = TrustThreshold{
			Denominator: strconv.FormatUint(denominator, 10),
			Numerator:   strconv.FormatUint(numerator, 10),
		}
	}
}

// WithChainAddressPrefix set the chain address prefix into the Hermes config.
func WithChainAddressPrefix(derivation string) ChainOption {
	return func(c *Chain) {
		c.AddressType = AddressType{Derivation: derivation}
	}
}

// AddChain adds a new chain into the Hermes config.
func (c *Config) AddChain(chainID, rpcAddr, grpcAddr string, options ...ChainOption) (Chain, error) {
	rpcURL, err := url.Parse(rpcAddr)
	if err != nil {
		return Chain{}, err
	}

	chain := Chain{
		ID:       chainID,
		RPCAddr:  rpcAddr,
		GRPCAddr: grpcAddr,
		EventSource: EventSource{
			BatchDelay: "500ms",
			Mode:       "push",
			URL:        fmt.Sprintf("ws://%s/websocket", rpcURL.Host),
		},
		RPCTimeout:    "15s",
		AccountPrefix: "cosmos",
		KeyName:       "wallet",
		StorePrefix:   "ibc",
		DefaultGas:    1000,
		MaxGas:        100000,
		GasPrice: GasPrice{
			Denom: "stake",
			Price: 0.01,
		},
		GasMultiplier:  1.1,
		MaxMsgNum:      30,
		MaxTxSize:      2097152,
		ClockDrift:     "5s",
		MaxBlockTime:   "10s",
		TrustingPeriod: "14days",
		TrustThreshold: TrustThreshold{
			Denominator: "3",
			Numerator:   "1",
		},
		AddressType: AddressType{
			Derivation: "cosmos",
		},
	}
	for _, o := range options {
		o(&chain)
	}

	c.Chains = append(c.Chains, chain)
	return chain, nil
}

// Balance returns the total account balance.
func (c *Chain) Balance(ctx context.Context, rpcAddress, addr string) (sdk.Coins, error) {
	client, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(rpcAddress))
	if err != nil {
		return nil, err
	}

	queryClient := banktypes.NewQueryClient(client.Context())
	res, err := queryClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{Address: addr})
	if err != nil {
		return nil, err
	}

	return res.Balances, nil
}

// TryRetrieve tries to receive some coins to the account and returns the total balance.
func (c *Chain) TryRetrieve(ctx context.Context, addr, faucetAddr string) (sdk.Coins, error) {
	if err := cosmosfaucet.TryRetrieve(ctx, c.ID, c.RPCAddr, faucetAddr, addr); err != nil {
		return nil, err
	}
	return c.Balance(ctx, c.RPCAddr, addr)
}
