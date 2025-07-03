package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	authv1betav1 "cosmossdk.io/api/cosmos/auth/v1beta1"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc"

	"github.com/ignite/cli/v29/ignite/pkg/chainregistry"
	"github.com/ignite/cli/v29/ignite/services/plugin"

	"github.com/ignite/apps/connect/chains"
)

type addCmdModel struct {
	spinner spinner.Model
	err     error
	page    int

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
	return tea.Batch() // no preloading needed
}

func (m *addCmdModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	totalSize := len(m.chain.APIs.Grpc)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n", "right":
			if m.page < totalSize/pageSize {
				m.page++
			}
		case "p", "left":
			if m.page > 0 {
				m.page--
			}
		case "up":
			if m.selectedIndex > 0 {
				m.selectedIndex--
			}

			if m.selectedIndex%pageSize == pageSize-1 {
				m.page--
				m.selectedIndex = 0
			}
		case "down":
			if m.page*pageSize+m.selectedIndex == totalSize-1 {
				return m, cmd
			}

			m.selectedIndex++

			if m.selectedIndex%pageSize == 0 {
				m.page++
				m.selectedIndex = 0
			}
		case "enter":
			selectedIndex := m.page*pageSize + m.selectedIndex
			m.selectedEndpoint = m.chain.APIs.Grpc[selectedIndex].Address
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
	totalSize := len(items)

	start := m.page * pageSize
	if start >= totalSize {
		start = 0
		m.page = 0
	}

	end := start + pageSize
	if end > totalSize {
		end = totalSize
	}

	out := "\033[K" // clear current line before printing
	out += "Select endpoint:\n"
	for i, item := range items[start:end] {
		out += "\033[K" // clear current line before printing
		if i == m.selectedIndex {
			out += "\033[32m✓ \033[1m" + item + "\033[0m\n"
		} else {
			out += fmt.Sprintf("  %s\n", item)
		}
	}

	out += "\033[K" // clear current line before printing
	out += "(press 'n'/'right' for next, 'p'/'left' for prev, 'enter' to add chain, 'q'/'ctrl+c' to quit)\n"

	return out
}

func AddHandler(ctx context.Context, cmd *plugin.ExecutedCommand) error {
	chainRegistry := chains.NewChainRegistry()
	if err := chainRegistry.FetchChains(); err != nil {
		return err
	}

	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return errors.New("usage: connect add <chain> [endpoint]")
	} else if len(cmd.Args) == 2 { // support custom chains
		return initChain(ctx, chainregistry.Chain{ChainName: cmd.Args[0]}, cmd.Args[1])
	}

	chain, ok := chainRegistry.Chains[cmd.Args[0]]
	if !ok {
		return fmt.Errorf("chain %s not found", cmd.Args[0])
	}

	model := newAddCmdModel(chain)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return err
	}

	fmt.Println("Selected endpoint:", model.selectedEndpoint)
	return initChain(context.TODO(), model.chain, model.selectedEndpoint)
}

func initChain(ctx context.Context, chain chainregistry.Chain, endpoint string) error {
	cfg, err := chains.ReadConfig()
	if err != nil && !errors.Is(err, chains.ErrConfigNotFound) {
		return err
	}

	// add chain to cfg
	chainCfg := &chains.ChainConfig{
		ChainID:      chain.ChainID,
		Bech32Prefix: chain.Bech32Prefix,
		GRPCEndpoint: endpoint,
	}

	fmt.Println("Initializing chain...")
	conn, err := chains.NewConn(chain.ChainName, chainCfg)
	if err != nil {
		return err
	}

	if err := conn.Load(ctx); err != nil {
		return err
	}

	// if no bech32 prefix is set, fetch it from the chain
	if chainCfg.Bech32Prefix == "" {
		client, err := conn.Connect()
		if err != nil {
			return err
		}

		chainCfg.Bech32Prefix, err = getAddressPrefix(ctx, client)
		if err != nil {
			return err
		}
	}

	cfg.Chains[chain.ChainName] = chainCfg
	if err := cfg.Save(); err != nil {
		return err
	}

	fmt.Printf("%s is ready to Connect!\n", strings.Title(chain.ChainName)) //nolintlint:staticcheck // strings.Title has a better API.
	return nil
}

// getAddressPrefix returns the address prefix of the chain.
func getAddressPrefix(ctx context.Context, conn grpc.ClientConnInterface) (string, error) {
	authClient := authv1betav1.NewQueryClient(conn)
	resp, err := authClient.Bech32Prefix(ctx, &authv1betav1.Bech32PrefixRequest{})
	if err != nil {
		return "", err
	}

	if resp == nil || resp.Bech32Prefix == "" {
		return "", errors.New("bech32 account address prefix is not set")
	}

	return resp.Bech32Prefix, nil
}
