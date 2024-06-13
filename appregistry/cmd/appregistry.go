package cmd

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

const (
	githubTokenFlag = "github-token"
)

// NewAppRegistry creates a new app registry command that holds
// some other sub commands related to running appregistry like
// list and info.
func NewAppRegistry() *cobra.Command {
	c := &cobra.Command{
		Use:     "appregistry [command]",
		Aliases: []string{"mp"},
		Short:   "Browse the Ignite App Registry App",
		Long: `AppRegistry is a command line tool that helps you to search for ignite apps.
It also helps you to get more information about an app.
Please note this command uses the Github API that a very limited rate limit for unauthenticated requests
so it's recommended to use the --github-token flag you want to use appregistry commands frequently.`,
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
