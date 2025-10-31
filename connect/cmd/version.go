package cmd

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

func VersionHandler(_ context.Context, _ *plugin.ExecutedCommand) error {
	version, ok := debug.ReadBuildInfo()
	if !ok {
		return errors.New("failed to get hubl version")
	}

	fmt.Println(strings.TrimSpace(version.Main.Version))
	return nil
}
