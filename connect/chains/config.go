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
	out, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	configDir, err := ConfigDir()
	if err != nil {
		return err
	}

	connectConfigPath := path.Join(configDir, configName)
	if err := os.WriteFile(connectConfigPath, out, 0o644); err != nil {
		return fmt.Errorf("error saving config: %w", err)
	}

	return nil
}

func ReadConfig() (*Config, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return nil, err
	}

	connectConfigPath := path.Join(configDir, configName)
	if _, err := os.Stat(connectConfigPath); os.IsNotExist(err) {
		return &Config{map[string]*ChainConfig{}}, ErrConfigNotFound
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

func ConfigDir() (string, error) {
	igniteConfigDir, err := igniteconfig.DirPath()
	if err != nil {
		return "", fmt.Errorf("failed to get ignite config directory: %w", err)
	}

	dir := path.Join(igniteConfigDir, "connect")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return dir, nil
}
