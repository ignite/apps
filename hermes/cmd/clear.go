package cmd

import (
	"context"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/services/plugin"

	"github.com/ignite/apps/hermes/pkg/hermes"
)

func ClearBinariesHandler(context.Context, *plugin.ExecutedCommand) error {
	session := cliui.New()
	defer session.End()

	session.StartSpinner("Clearing previous binaries")

	path, err := hermes.ClearBinCachePath()
	if err != nil {
		return err
	}

	session.StopSpinner()
	return session.Printf("All previous binaries were deleted from: %s", path)
}

func ClearConfigsHandler(context.Context, *plugin.ExecutedCommand) error {
	session := cliui.New()
	defer session.End()

	session.StartSpinner("Clearing previous config paths")

	path, err := hermes.ClearConfigPath()
	if err != nil {
		return err
	}

	session.StopSpinner()
	return session.Printf("All previous configurations were deleted from: %s", path)
}
