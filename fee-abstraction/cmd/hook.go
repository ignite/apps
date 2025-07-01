package cmd

import (
	"context"

	"github.com/blang/semver/v4"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
	"github.com/ignite/cli/v29/ignite/services/chain"
	"github.com/ignite/cli/v29/ignite/services/plugin"

	"github.com/ignite/apps/fee-abstraction/services/scaffolder"
)

// ExecuteScaffoldPreHook executes the scaffold pre hook.
func ExecuteScaffoldPreHook(h *plugin.ExecutedHook) error {
	var (
		flags       = plugin.Flags(h.Hook.Flags)
		noModule, _ = flags.GetBool(flagNoModule)
		name        = h.ExecutedCommand.Args[0]
	)
	if !noModule && name == feeAbsModuleName {
		return errors.Errorf("cannot scaffold module with name %s", feeAbsModuleName)
	}
	return nil
}

// ExecuteScaffoldChainPostHook executes the scaffold chain post hook.
func ExecuteScaffoldChainPostHook(ctx context.Context, h *plugin.ExecutedHook) error {
	var (
		flags           = plugin.Flags(h.Hook.Flags)
		feeAbsModule, _ = flags.GetBool(flagFeeAbsModule)
		name            = h.ExecutedCommand.Args[0]
	)
	if !feeAbsModule {
		return nil
	}

	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

	pathInfo, err := gomodulepath.Parse(name)
	if err != nil {
		return err
	}

	version := getVersion(flags)
	semVersion, err := semver.Parse(version)
	if err != nil {
		return err
	}

	c, err := newChain(pathInfo.Root, flags, chain.WithOutputer(session), chain.CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	sc, err := scaffolder.New(c, session)
	if err != nil {
		return err
	}

	sm, err := sc.AddFeeAbstraction(ctx, placeholder.New(), scaffolder.WithVersion(semVersion))
	if err != nil {
		return err
	}

	modificationsStr, err := sourceModificationToString(sm)
	if err != nil {
		return err
	}

	session.Println(modificationsStr)
	session.Printf("\n🎉 Fee Abstraction added (`%[1]v`).\n\n", c.AppPath())

	return nil
}
