package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/ignite/apps/connect/chains"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

func RemoveHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	cfg, err := chains.ReadConfig()
	if errors.Is(err, chains.ErrConfigNotFound) {
		return nil
	} else if err != nil {
		return nil
	}

	chainName := cmd.Args[0]

	// delete config
	delete(cfg.Chains, chainName)
	if err := cfg.Save(); err != nil {
		return err
	}

	// delete proto files
	configDir, err := chains.ConfigDir()
	if err != nil {
		return err
	}

	_ = os.Remove(path.Join(configDir, fmt.Sprintf("%s.fds", chainName)))
	_ = os.Remove(path.Join(configDir, fmt.Sprintf("%s.autocli", chainName)))

	fmt.Printf("Chain %s successfully removed from Connect\n", chainName)
	return nil
}
