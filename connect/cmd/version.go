package cmd

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/ignite/cli/v28/ignite/services/plugin"
)

func VersionHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	version, ok := debug.ReadBuildInfo()
	if !ok {
		return errors.New("failed to get hubl version")
	}

	fmt.Println(strings.TrimSpace(version.Main.Version))
	return nil
}
