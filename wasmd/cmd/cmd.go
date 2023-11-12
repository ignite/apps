package cmd

import "github.com/spf13/cobra"

const (
	flagPath = "path"
	flagYes  = "yes"
)

const (
	msgCommitPrefix = "Your saved project changes have not been committed.\nTo enable reverting to your current state, commit your saved changes."
	msgCommitPrompt = "Do you want to proceed without committing your saved changes"
)

func getYes(cmd *cobra.Command) (ok bool) {
	ok, _ = cmd.Flags().GetBool(flagYes)
	return
}

func flagGetPath(cmd *cobra.Command) (path string) {
	path, _ = cmd.Flags().GetString(flagPath)
	return
}
