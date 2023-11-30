package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const appsSearchQuery = "topic:ignite-cli-app language:go stars:>10"

var githubAccessToken = os.Getenv("GITHUB_ACCESS_TOKEN")

// NewMarketplace creates a new marketplace command that holds
// some other sub commands related to running marketplace like
// list and info.
func NewMarketplace() *cobra.Command {
	c := &cobra.Command{
		Use:     "marketplace [command]",
		Aliases: []string{"mp"},
		Short:   "Run marketplace commands",
		Long: `Marketplace is a command line tool that helps you to search for ignite apps
			using GitHub search API. It also helps you to get more information about an app.
			Please note that Github API has a very limited rate limit for unauthenticated requests
			so it's recommended to set GITHUB_ACCESS_TOKEN environment variable to your GitHub access token
			if you want to use marketplace commands frequently.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// add sub commands.
	c.AddCommand(
		NewList(),
		NewInfo(),
	)

	return c
}
