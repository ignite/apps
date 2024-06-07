package cmd

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

const (
	githubTokenFlag = "github-token"
)

// NewMarketplace creates a new marketplace command that holds
// some other sub commands related to running marketplace like
// list and info.
func NewMarketplace() *cobra.Command {
	c := &cobra.Command{
		Use:     "marketplace [command]",
		Aliases: []string{"mp"},
		Short:   "Browse the Ignite App Marketplace",
		Long: `Marketplace is a command line tool that helps you to search for ignite apps.
It also helps you to get more information about an app.
Please note this command uses the Github API that a very limited rate limit for unauthenticated requests
so it's recommended to use the --github-token flag you want to use marketplace commands frequently.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	c.AddCommand(
		NewListCmd(),
		NewDetailsCmd(),
		NewInstallCmd(),
	)

	c.PersistentFlags().String(githubTokenFlag, "", "GitHub access token")

	return c
}

func init() {
	lipgloss.SetColorProfile(termenv.TrueColor)
}
