package cmd

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

const (
	flagGithubToken = "github-token"
	flagBranch      = "branch"
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
		NewValidateCmd(),
		NewInstallCmd(),
	)

	c.PersistentFlags().String(flagGithubToken, "", "GitHub access token")
	c.PersistentFlags().StringP(flagBranch, "b", "main", "The app branch to use (default: main)")

	return c
}

func getBranchFlag(cmd *cobra.Command) string {
	if branch, _ := cmd.Flags().GetString(flagBranch); branch != "" {
		return branch
	}
	return "main"
}

func getGitHubToken(cmd *cobra.Command) string {
	if githubToken, _ := cmd.Flags().GetString(flagGithubToken); githubToken != "" {
		return githubToken
	}
	if envToken := os.Getenv("GITHUB_TOKEN"); envToken != "" {
		return envToken
	}
	return ""
}

func init() {
	lipgloss.SetColorProfile(termenv.TrueColor)
}
