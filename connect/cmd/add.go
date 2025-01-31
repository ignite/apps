package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/ignite/cli/v28/ignite/pkg/chainregistry"
	"github.com/ignite/cli/v28/ignite/services/plugin"

	"github.com/ignite/apps/connect/chains"
)

type addCmdModel struct {
	spinner          spinner.Model
	err              error
	chain            chainregistry.Chain
	selectedEndpoint string
	selectedIndex    int
}

func newAddCmdModel(chain chainregistry.Chain) *addCmdModel {
	s := spinner.NewModel()
	s.Spinner = spinner.Dot

	c := &chain
	if err := chains.EnrichChain(c); err != nil {
		return &addCmdModel{err: err}
	}

	return &addCmdModel{
		spinner:       s,
		chain:         *c,
		selectedIndex: 0,
	}
}

func (m *addCmdModel) Init() tea.Cmd {
	return tea.Batch(fetchChainsCmd)
}

func (m *addCmdModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case "down":
			if m.selectedIndex < len(m.chain.APIs.Grpc)-1 {
				m.selectedIndex++
			}
		case "enter":
			m.selectedEndpoint = m.chain.APIs.Grpc[m.selectedIndex].Address
			return m, tea.Quit
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
	default:
	}

	return m, cmd
}

func (m *addCmdModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n(press 'q' to quit)\n", m.err)
	}

	items := make([]string, len(m.chain.APIs.Grpc))
	for i, endpoint := range m.chain.APIs.Grpc {
		items[i] = endpoint.Address
	}

	out := "\033[K" // clear current line before printing
	out += "Select endpoint:\n"
	for i, item := range items {
		out += "\033[K" // clear current line before printing
		prefix := " "
		if i == m.selectedIndex {
			prefix = ">"
		}

		out += fmt.Sprintf("%s %s\n", prefix, item)
	}

	out += "\033[K" // clear current line before printing
	out += "(press 'up' or 'down' for selecting, 'enter' to add chain, 'q'/'ctrl+c' to quit)\n"

	return out
}

func AddHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	chainRegistry := chains.NewChainRegistry()
	if err := chainRegistry.FetchChains(); err != nil {
		return err
	}

	if len(cmd.Args) < 1 {
		return fmt.Errorf("please provide a chain name as argument")
	}

	chain, ok := chainRegistry.Chains[cmd.Args[0]]
	if !ok {
		return fmt.Errorf("chain %s not found", cmd.Args[0])
	}

	p := tea.NewProgram(newAddCmdModel(chain))
	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
