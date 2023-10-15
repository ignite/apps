package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/ignite/apps/airdrop/pkg/formula"
)

var (
	ErrInvalidConfig         = errors.New("invalid config file")
	ErrInvalidSnapshotConfig = errors.New("invalid snapshot config")
)

type (
	// Config defines a struct with the fields that are common to all config.
	Config struct {
		AirdropToken string     `json:"airdrop_token" yaml:"airdrop_token"`
		DustWallet   uint64     `json:"dust_wallet" yaml:"dust_wallet"`
		Snapshots    []Snapshot `json:"snapshots" yaml:"snapshots"`
	}

	// Snapshot defines a struct with the fields that are common to all config snapshot.
	Snapshot struct {
		Type     string        `json:"type" yaml:"type"`
		Denom    string        `json:"denom" yaml:"denom"`
		Formula  formula.Value `json:"formula" yaml:"formula"`
		Excluded []string      `json:"excluded" yaml:"excluded"`
	}
)

// ParseConfig expects to find and parse a config file.
func ParseConfig(filename string) (c Config, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return c, err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&c); err != nil {
		return c, err
	}
	return c, c.validate()
}

// validate validates the config parameters.
func (c Config) validate() error {
	if c.AirdropToken == "" {
		return errors.Wrap(ErrInvalidConfig, "airdrop token type not defined")
	}
	if len(c.Snapshots) == 0 {
		return errors.Wrap(ErrInvalidConfig, "snapshots not defined")
	}
	for _, snapshot := range c.Snapshots {
		if snapshot.Type == "" {
			return errors.Wrap(ErrInvalidSnapshotConfig, "snapshot type not defined")
		}
		if snapshot.Denom == "" {
			return errors.Wrap(ErrInvalidSnapshotConfig, "snapshot denom not defined")
		}
		if snapshot.Formula.Type == "" {
			return errors.Wrap(ErrInvalidSnapshotConfig, "snapshot formula type not defined")
		}
		if snapshot.Formula.Value == 0 {
			return errors.Wrap(ErrInvalidSnapshotConfig, "invalid snapshot formula value")
		}
	}
	return nil
}
