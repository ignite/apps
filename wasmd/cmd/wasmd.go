package cmd

import (
	"errors"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/xgit"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const (
	statusWasmInit = "Adding Wasm to your chain..."
)

// NewWasmd creates a new wasmd command that holds
// some other sub commands related to running chain wasmds like gex and etc.
func NewWasmd() *cobra.Command {
	c := &cobra.Command{
		Use:           "wasmd [command]",
		Aliases:       []string{"ws"},
		Short:         "Run chain wasmd commands",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// add sub commands.
	c.AddCommand(NewWasmdInit())
	// c.AddCommand(NewWasmdSC())

	/*

		ignite wasmd init
		ignite wasmd sc(smart contract) new/test/deploy

	*/

	return c
}

func gitChangesConfirmPreRunHandler(cmd *cobra.Command, _ []string) error {
	// Don't confirm when the "--yes" flag is present
	if getYes(cmd) {
		return nil
	}

	appPath := flagGetPath(cmd)
	session := cliui.New()

	defer session.End()

	return confirmWhenUncommittedChanges(session, appPath)
}

func confirmWhenUncommittedChanges(session *cliui.Session, appPath string) error {
	cleanState, err := xgit.AreChangesCommitted(appPath)
	if err != nil {
		return err
	}

	if !cleanState {
		session.Println(msgCommitPrefix)
		if err := session.AskConfirm(msgCommitPrompt); err != nil {
			if errors.Is(err, promptui.ErrAbort) {
				return errors.New("No")
			}

			return err
		}
	}

	return nil
}
