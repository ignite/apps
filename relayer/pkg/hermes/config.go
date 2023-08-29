package hermes

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type (
	Config struct {
		Chains    []Chain   `toml:"chains" json:"chains"`
		Global    Global    `toml:"global" json:"global"`
		Telemetry Telemetry `toml:"telemetry" json:"telemetry"`
		Mode      Mode      `toml:"mode" json:"mode"`
	}

	Chain struct {
		Id             string         `toml:"id" json:"id"`
		RpcAddr        string         `toml:"rpc_addr" json:"rpc_addr"`
		GrpcAddr       string         `toml:"grpc_addr" json:"grpc_addr"`
		EventSource    EventSource    `toml:"event_source,inline" json:"event_source"`
		RpcTimeout     string         `toml:"rpc_timeout" json:"rpc_timeout"`
		AccountPrefix  string         `toml:"account_prefix" json:"account_prefix"`
		KeyName        string         `toml:"key_name" json:"key_name"`
		StorePrefix    string         `toml:"store_prefix" json:"store_prefix"`
		DefaultGas     int            `toml:"default_gas" json:"default_gas"`
		MaxGas         int            `toml:"max_gas" json:"max_gas"`
		GasPrice       GasPrice       `toml:"gas_price,inline" json:"gas_price"`
		GasMultiplier  float64        `toml:"gas_multiplier" json:"gas_multiplier"`
		MaxMsgNum      int            `toml:"max_msg_num" json:"max_msg_num"`
		MaxTxSize      int            `toml:"max_tx_size" json:"max_tx_size"`
		ClockDrift     string         `toml:"clock_drift" json:"clock_drift"`
		MaxBlockTime   string         `toml:"max_block_time" json:"max_block_time"`
		TrustingPeriod string         `toml:"trusting_period" json:"trusting_period"`
		TrustThreshold TrustThreshold `toml:"trust_threshold,inline" json:"trust_threshold"`
		AddressType    AddressType    `toml:"address_type,inline" json:"address_type"`
	}

	EventSource struct {
		BatchDelay string `toml:"batch_delay" json:"batch_delay"`
		Mode       string `toml:"mode" json:"mode"`
		Url        string `toml:"url" json:"url"`
	}

	GasPrice struct {
		Denom string  `toml:"denom" json:"denom"`
		Price float64 `toml:"price" json:"price"`
	}

	TrustThreshold struct {
		Denominator string `toml:"denominator" json:"denominator"`
		Numerator   string `toml:"numerator" json:"numerator"`
	}

	AddressType struct {
		Derivation string `toml:"derivation" json:"derivation"`
	}

	Global struct {
		LogLevel string `toml:"log_level" json:"log_level"`
	}

	Telemetry struct {
		Enabled bool   `toml:"enabled" json:"enabled"`
		Host    string `toml:"host" json:"host"`
		Port    int    `toml:"port" json:"port"`
	}

	Mode struct {
		Channels    Channels    `toml:"channels" json:"channels"`
		Clients     Clients     `toml:"clients" json:"clients"`
		Connections Connections `toml:"connections" json:"connections"`
		Packets     Packets     `toml:"packets" json:"packets"`
	}

	Channels struct {
		Enabled bool `toml:"enabled" json:"enabled"`
	}

	Clients struct {
		Enabled      bool `toml:"enabled" json:"enabled"`
		Misbehaviour bool `toml:"misbehaviour" json:"misbehaviour"`
		Refresh      bool `toml:"refresh" json:"refresh"`
	}

	Connections struct {
		Enabled bool `toml:"enabled" json:"enabled"`
	}

	Packets struct {
		ClearInterval  int  `toml:"clear_interval" json:"clear_interval"`
		ClearOnStart   bool `toml:"clear_on_start" json:"clear_on_start"`
		Enabled        bool `toml:"enabled" json:"enabled"`
		TxConfirmation bool `toml:"tx_confirmation" json:"tx_confirmation"`
	}

	// ChainOption configures chain hermes configs.
	ChainOption func(*Chain)
)

func WithEventSource(mode, url, batchDelay string) ChainOption {
	return func(c *Chain) {
		c.EventSource = EventSource{
			BatchDelay: batchDelay,
			Mode:       mode,
			Url:        url,
		}
	}
}

func WithRPCTimeout(timeout string) ChainOption {
	return func(c *Chain) {
		c.RpcTimeout = timeout
	}
}

func WithAccountPrefix(prefix string) ChainOption {
	return func(c *Chain) {
		c.AccountPrefix = prefix
	}
}

func WithKeyName(key string) ChainOption {
	return func(c *Chain) {
		c.KeyName = key
	}
}

func WithStorePrefix(prefix string) ChainOption {
	return func(c *Chain) {
		c.StorePrefix = prefix
	}
}

func WithDefaultGas(defaultGas int) ChainOption {
	return func(c *Chain) {
		c.DefaultGas = defaultGas
	}
}

func WithMaxGas(maxGas int) ChainOption {
	return func(c *Chain) {
		c.MaxGas = maxGas
	}
}

func WithGasPrice(price float64, denom string) ChainOption {
	return func(c *Chain) {
		c.GasPrice = GasPrice{
			Denom: denom,
			Price: price,
		}
	}
}

func WithGasMultiplier(gasMultipler float64) ChainOption {
	return func(c *Chain) {
		c.GasMultiplier = gasMultipler
	}
}

func WithMaxMsgNum(maxMsg int) ChainOption {
	return func(c *Chain) {
		c.MaxMsgNum = maxMsg
	}
}

func WithMaxTxSize(size int) ChainOption {
	return func(c *Chain) {
		c.MaxTxSize = size
	}
}

func WithClockDrift(clock string) ChainOption {
	return func(c *Chain) {
		c.ClockDrift = clock
	}
}

func WithMaxBlockTime(maxBlockTime string) ChainOption {
	return func(c *Chain) {
		c.MaxBlockTime = maxBlockTime
	}
}

func WithTrustingPeriod(trustingPeriod string) ChainOption {
	return func(c *Chain) {
		c.TrustingPeriod = trustingPeriod
	}
}

func WithTrustThreshold(numerator, denominator string) ChainOption {
	return func(c *Chain) {
		c.TrustThreshold = TrustThreshold{
			Denominator: denominator,
			Numerator:   numerator,
		}
	}
}

func WithAddressType(derivation string) ChainOption {
	return func(c *Chain) {
		c.AddressType = AddressType{Derivation: derivation}
	}
}

func (c *Config) AddChain(chainID, rpcAddr, grpcAddr string, options ...ChainOption) error {
	rpcUrl, err := url.Parse(rpcAddr)
	if err != nil {
		return err
	}

	chain := Chain{
		Id:       chainID,
		RpcAddr:  rpcAddr,
		GrpcAddr: grpcAddr,
		EventSource: EventSource{
			BatchDelay: "500ms",
			Mode:       "push",
			Url:        fmt.Sprintf("ws://%s", rpcUrl.Host),
		},
		RpcTimeout:    "15s",
		AccountPrefix: "cosmos",
		KeyName:       "wallet",
		StorePrefix:   "ibc",
		DefaultGas:    100000,
		MaxGas:        10000000,
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
	return nil
}

func (c *Config) Save() error {
	configPath, err := ConfigPath()
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

func Parse(path string) (cfg Config, err error) {
	err = toml.Unmarshal([]byte(path), &cfg)
	return
}

func DefaultConfig() *Config {
	return &Config{
		Chains: []Chain{},
		Global: Global{
			LogLevel: "info",
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
			Enabled: true,
			Host:    "127.0.0.1",
			Port:    3001,
		},
	}
}

func ConfigPath() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userHomeDir, ".hermes", "config.toml"), nil
}
