package cmd

import (
	"fmt"
	"hash/fnv"
	"path"
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/ignite/apps/marketplace/pkg/apps"
	"github.com/ignite/apps/marketplace/pkg/xgithub"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/spf13/cobra"
)

const igniteCLIPackage = "github.com/ignite/cli"

var (
	linkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Underline(true)
	installaitonStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("9")).
				MarginLeft(7)
	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true)
)

// NewInfo creates a new info command that shows the details of an ignite application repository.
func NewInfo() *cobra.Command {
	return &cobra.Command{
		Use:   "info [package-url]",
		Short: "Show the details of an ignite application repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			githubToken, _ := cmd.Flags().GetString(githubTokenFlag)

			session := cliui.New(cliui.StartSpinner())
			defer session.End()

			session.StartSpinner("🔎 Fetching repository details from GitHub...")

			client := xgithub.NewClient(githubToken)
			repo, err := apps.GetRepositoryDetails(cmd.Context(), client, args[0])
			if err != nil {
				return err
			}

			session.StopSpinner()

			printRepoDetails(session, repo)

			return nil
		},
	}
}

func printRepoDetails(sess *cliui.Session, repo *apps.AppRepositoryDetails) {
	sess.Println("Description:", repo.Description)
	sess.Print("Tags:")
	for _, tag := range repo.Tags {
		sess.Print(lipgloss.NewStyle().Background(colorFromText(tag)).Render(tag), " ")
	}
	sess.Println()
	sess.Println("Stars ⭐️:", repo.Stars)
	sess.Println("License ⚖️ :", repo.License)
	sess.Printf("Updated At 🕒: %s %s\n", repo.UpdatedAt.Format(time.DateTime), updatedAtStyle.Render("("+humanize.Time(repo.UpdatedAt)+")"))
	sess.Println("URL 🌎: ", linkStyle.Render(repo.URL))
	sess.Println("Apps 🔥:")
	printAppsTable(sess, repo)
}

func colorFromText(text string) lipgloss.Color {
	h := fnv.New64a()
	h.Write([]byte(text))
	return lipgloss.Color(strconv.FormatUint(h.Sum64()%16, 10))
}

func printAppsTable(sess *cliui.Session, repo *apps.AppRepositoryDetails) {
	for i, app := range repo.Apps {
		sess.Println("\tName:", app.Name)
		sess.Println("\tDescription:", app.Description)
		sess.Println("\tPath:", app.Path)
		sess.Println("\tGo Version:", app.GoVersion)
		sess.Println("\tIgnite Version:", app.IgniteVersion)
		sess.Println(installaitonStyle.Render(fmt.Sprintf(
			"🚀 Install via: %s",
			commandStyle.Render(fmt.Sprintf("ignite app -g install %s", path.Join(repo.PackageURL, app.Path))),
		)))

		if i < len(repo.Apps)-1 {
			sess.Println()
		}
	}
}
