package cmd

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strconv"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ignite/apps/connect/chains"
	"github.com/ignite/cli/v28/ignite/pkg/chainregistry"
	"github.com/ignite/cli/v28/ignite/services/plugin"
)

const pageSize = 10

// Model for Bubble Tea UI
type model struct {
	spinner  spinner.Model
	fetching bool
	err      error
	page     int
	input    string // track user input for chain number

	chainRegistry *chains.ChainRegistry
	selectedChain chainregistry.Chain
}

func initialModel() *model {
	s := spinner.NewModel()
	s.Spinner = spinner.Dot
	return &model{
		spinner:  s,
		fetching: true,
	}
}

// Messages
type fetchDoneMsg struct {
	reg *chains.ChainRegistry
}
type fetchErrMsg struct{ err error }

// Initialize Bubble Tea program
func (m *model) Init() tea.Cmd {
	return tea.Batch(fetchChainsCmd)
}

// Fetch the chains in the background
func fetchChainsCmd() tea.Msg {
	cr := chains.NewChainRegistry()
	if err := cr.FetchChains(); err != nil {
		return fetchErrMsg{err}
	}
	return fetchDoneMsg{cr}
}

// Update handles messages and updates the model accordingly
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n", "right":
			m.page++
			m.input = "" // Clear input when navigating
		case "p", "left":
			if m.page > 0 {
				m.page--
			}
			m.input = "" // Clear input when navigating
		case "enter":
			if m.input != "" {
				if m.chainRegistry != nil {
					chains := m.chainRegistry.Chains
					if len(chains) > 0 {
						chainsSlice := slices.Collect(maps.Keys(chains))
						start := m.page * pageSize
						if start >= len(chainsSlice) {
							start = 0
							m.page = 0
						}
						end := start + pageSize
						if end > len(chainsSlice) {
							end = len(chainsSlice)
						}

						inputNum, err := strconv.Atoi(m.input)
						if err != nil {
							m.err = fmt.Errorf("invalid input: %s", m.input)
							m.input = ""
							return m, nil
						}

						// Adjust inputNum to be 0-indexed
						selectedIndex := start + inputNum - 1
						if selectedIndex >= start && selectedIndex < end {
							chainKey := chainsSlice[selectedIndex]

							m.selectedChain = chains[chainKey]
							return m, tea.Quit // chain selected, quit
						} else {
							m.err = fmt.Errorf("invalid chain number: %s", m.input)
							m.input = ""
						}
					}
				}
			}
		case "backspace", "delete":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			if len(msg.String()) == 1 && msg.String()[0] >= '0' && msg.String()[0] <= '9' {
				m.input += msg.String()
			}
		}
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
	case fetchDoneMsg:
		m.fetching = false
		m.chainRegistry = msg.reg
	case fetchErrMsg:
		m.err = msg.err
		m.fetching = false
	}
	return m, cmd
}

// View returns the UI as a string
func (m *model) View() string {
	if m.fetching && m.err == nil {
		return fmt.Sprintf("%s Discovering chains... (press 'q' to quit)\n", m.spinner.View())
	}
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n(press 'q' to quit)\n", m.err)
	}

	chains := m.chainRegistry.Chains
	total := len(chains)
	start := m.page * pageSize
	if start >= total {
		start = 0
		m.page = 0
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	chainsSlice := slices.Sorted(maps.Keys(chains))
	out := "\033[K" // clear current line before printing
	out += fmt.Sprintf("Fetched %d chains. Showing %d-%d:\n", total, start+1, end)

	// ANSI escape codes for highlighting
	highlightStart := "\033[32;1m" // Green color and bold
	highlightEnd := "\033[0m"      // Reset to default

	for i, k := range chainsSlice[start:end] {
		chain := m.chainRegistry.Chains[k]
		lineNumber := i + 1 + start

		out += "\033[K" // clear current line before printing
		// Check if the input corresponds to this line number
		inputNum, err := strconv.Atoi(m.input)
		if err == nil && inputNum == lineNumber {
			out += fmt.Sprintf("%s%d. %s (%s)%s\n", highlightStart, lineNumber, chain.PrettyName, chain.ChainName, highlightEnd)
		} else {
			out += fmt.Sprintf("%d. %s (%s)\n", lineNumber, chain.PrettyName, chain.ChainName)
		}
	}

	out += "\033[K" // clear current line before printing
	out += fmt.Sprintf("Select chain: %s\n", m.input)
	out += "(press 'n'/'right' for next, 'p'/'left' for prev, 'enter' to init chain, 'q'/'ctrl+c' to quit)\n"
	return out
}

func DiscoverHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	model := initialModel()
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}

	// if the user selected a chain, execute the add command
	if len(model.selectedChain.ChainName) > 0 {
		return initializeChain(ctx, model.selectedChain)
	}

	return nil
}
