package cmd

import (
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/spf13/cobra"

	"github.com/ignite/apps/marketplace/pkg/xgithub"
	"github.com/ignite/apps/marketplace/registry"
)

var (
	linkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Underline(true)
	installaitonStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("9")).
				MarginLeft(15)
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

			session := cliui.New(cliui.StartSpinnerWithText("🔎 Fetching repository details from GitHub..."))
			defer session.End()

			client := xgithub.NewClient(githubToken)
			repo, err := registry.GetRepositoryDetails(cmd.Context(), client, args[0])
			if err != nil {
				return err
			}

			session.StopSpinner()
			printRepoDetails(repo)
			return nil
		},
	}
}

func printRepoDetails(repo *registry.AppRepositoryDetails) {
	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	printItem := func(s string, v interface{}) {
		fmt.Fprintf(w, "%s:\t%v\n", s, v)
	}

	printItem("Description", repo.Description)
	tags := make([]string, len(repo.Tags))
	for i, tag := range repo.Tags {
		tags[i] = lipgloss.NewStyle().Background(colorFromText(tag)).Render(tag)
	}
	printItem("Tags", strings.Join(tags, " "))
	printItem("Stars", strconv.Itoa(repo.Stars))
	printItem("License", repo.License)
	printItem("URL", linkStyle.Render(repo.URL))
	printItem("Apps", "")
	w.Flush()

	printAppsTable(repo)
}

func printAppsTable(repo *registry.AppRepositoryDetails) {
	printItem := func(w io.Writer, s string, v interface{}) {
		fmt.Fprintf(w, "\t%s:\t%v\n", s, v)
	}

	for _, app := range repo.Apps {
		w := &tabwriter.Writer{}
		w.Init(os.Stdout, 16, 8, 0, '\t', 0)

		printItem(w, "Name", app.Name)
		printItem(w, "Description", app.Description)
		printItem(w, "Path", app.Path)
		printItem(w, "Go Version", app.GoVersion)
		printItem(w, "Ignite Version", app.IgniteVersion)

		fmt.Fprintln(w, installaitonStyle.Render(fmt.Sprintf(
			"🚀 Install via: %s",
			commandStyle.Render(fmt.Sprintf("ignite app -g install %s", path.Join(repo.PackageURL, app.Path))),
		)))

		w.Flush()
	}
}

func colorFromText(text string) lipgloss.Color {
	h := fnv.New64a()
	h.Write([]byte(text))
	return lipgloss.Color(strconv.FormatUint(h.Sum64()%16, 10))
}
