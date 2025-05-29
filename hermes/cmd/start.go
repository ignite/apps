package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/hermes/pkg/hermes"
)

func StartHandler(ctx context.Context, cmd *plugin.ExecutedCommand) (err error) {
	var (
		flags     = cmd.Flags
		args      = cmd.Args
		customCfg = getConfig(flags)
		cfgName   = strings.Join(args, hermes.ConfigNameSeparator)
		session   = cliui.New()
	)
	defer session.End()

	hermesVersion, err := getVersion(flags)
	if err != nil {
		return err
	}

	session.StartSpinner("Fetching hermes config")
	cfgPath := customCfg
	if cfgPath == "" {
		cfgPath, err = hermes.ConfigPath(cfgName)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return errors.Errorf("config file (%s) not exist, try to configure you relayer first", cfgPath)
	}

	session.StartSpinner(fmt.Sprintf("Fetching hermes binary %s", hermesVersion))
	h, err := hermes.New(hermesVersion)
	if err != nil {
		return err
	}
	session.StopSpinner()

	return h.Start(
		ctx,
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdIn(os.Stdin),
		hermes.WithStdOut(os.Stdout),
		hermes.WithStdErr(os.Stderr),
	)
}
