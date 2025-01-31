package chains

import (
	"fmt"
	"os"
	"path"

	igniteconfig "github.com/ignite/cli/v28/ignite/config"

	"gopkg.in/yaml.v3"
)

var (
	configName = "connect.yaml"

	// ErrConfigNotFound is returned when the config file is not found
	ErrConfigNotFound = fmt.Errorf("config file not found")
)

type Config struct {
	Chains map[string]*ChainConfig `yaml:"chains"`
}

type ChainConfig struct {
	ChainID      string `yaml:"chain_id"`
	Bech32Prefix string `yaml:"bech32_prefix"`
	GRPCEndpoint string `yaml:"grpc_endpoint"`
}

func (c *Config) Save() error {
	igniteConfigDir, err := igniteconfig.DirPath()
	if err != nil {
		return fmt.Errorf("failed to get ignite config directory: %w", err)
	}

	out, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	connectConfigPath := path.Join(igniteConfigDir, "connect", configName)
	if err := os.MkdirAll(path.Dir(connectConfigPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(connectConfigPath, out, 0644); err != nil {
		return fmt.Errorf("error saving config: %w", err)
	}

	return nil
}

func ReadConfig() (*Config, error) {
	igniteConfigDir, err := igniteconfig.DirPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get ignite config directory: %w", err)
	}

	connectConfigPath := path.Join(igniteConfigDir, "connect", configName)
	if _, err := os.Stat(connectConfigPath); os.IsNotExist(err) {
		return &Config{}, ErrConfigNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to check config file: %w", err)
	}

	data, err := os.ReadFile(connectConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &c, nil
}
