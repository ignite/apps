package cmd

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/ignite/cli/v28/ignite/pkg/cache"
	"github.com/ignite/cli/v28/ignite/services/chain"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	"github.com/spf13/pflag"
)

// ExecuteBuild executes the build subcommand.
func ExecuteBuild(ctx context.Context, cmd *plugin.ExecutedCommand, c *chain.Chain) error {
	flags, err := cmd.NewFlags()
	if err != nil {
		return err
	}
	output, err := getOutputFlag(flags)
	if err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp(os.TempDir(), "buildcache-*")
	if err != nil {
		return fmt.Errorf("could not create a temp dir: %w", err)
	}
	cacheStorage, err := cache.NewStorage(path.Join(tempDir, "cacheStorage.db"))
	if err != nil {
		return fmt.Errorf("could not prepare a cache storage for building chain: %w", err)
	}
	binaryName, err := c.Build(ctx, cacheStorage, nil, output, false, true)
	if err != nil {
		return fmt.Errorf("building chain failed with error: %w", err)
	}
	fmt.Printf("Chain built successfully at %s\n", binaryName)
	return nil
}

func getOutputFlag(flags *pflag.FlagSet) (string, error) {
	out, err := flags.GetString("output")
	if err != nil {
		return "", fmt.Errorf("could not get --output flag: %w", err)
	}

	if out == "" {
		return ".", nil
	}
	return out, nil
}
