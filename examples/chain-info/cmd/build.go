package cmd

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/ignite/cli/v29/ignite/pkg/cache"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

// ExecuteBuild executes the build subcommand.
func ExecuteBuild(ctx context.Context, cmd *plugin.ExecutedCommand, c *chain.Chain) error {
	flags := plugin.Flags(cmd.Flags)

	output, err := flags.GetString(flagOutput)
	if err != nil {
		return errors.Errorf("could not get --%s flag: %s", flagOutput, err)
	}

	tempDir, err := os.MkdirTemp(os.TempDir(), "buildcache-*")
	if err != nil {
		return errors.Errorf("could not create a temp dir: %s", err)
	}
	cacheStorage, err := cache.NewStorage(path.Join(tempDir, "cacheStorage.db"))
	if err != nil {
		return errors.Errorf("could not prepare a cache storage for building chain: %s", err)
	}
	binaryName, err := c.Build(ctx, cacheStorage, nil, output, false, true)
	if err != nil {
		return errors.Errorf("building chain failed with error: %s", err)
	}
	fmt.Printf("Chain built successfully at %s\n", binaryName)
	return nil
}
