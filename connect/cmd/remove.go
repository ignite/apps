package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/ignite/cli/v29/ignite/services/plugin"

	"github.com/ignite/apps/connect/chains"
)

func RemoveHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	if len(cmd.Args) < 1 {
		return errors.New("usage: connect remove <chain>")
	}

	cfg, err := chains.ReadConfig()
	if errors.Is(err, chains.ErrConfigNotFound) {
		return nil
	} else if err != nil {
		return nil
	}

	chainName := cmd.Args[0]
	if _, ok := cfg.Chains[chainName]; !ok {
		return errors.New("chain not found")
	}

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

	fmt.Printf("Chain %s successfully removed!\n", chainName)
	return nil
}
