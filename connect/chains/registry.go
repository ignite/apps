package chains

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ignite/cli/v28/ignite/pkg/chainregistry"
)

const (
	repoURL               = "https://github.com/cosmos/chain-registry"
	cosmosDirectoryAPIURL = "https://chains.cosmos.directory"
)

type ChainRegistry struct {
	Chains map[string]chainregistry.Chain
	Assets map[string]chainregistry.Asset
}

func NewChainRegistry() *ChainRegistry {
	return &ChainRegistry{
		Chains: make(map[string]chainregistry.Chain),
		Assets: make(map[string]chainregistry.Asset),
	}
}

// FetchChains fetches the list of chains from the cosmos.directory API
// Note, the output chainregistry.Chain doesn't contain the full list of fields
func (r *ChainRegistry) FetchChains() error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, cosmosDirectoryAPIURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch chains: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var cdOutput map[string]json.RawMessage
	if err := json.Unmarshal(body, &cdOutput); err != nil {
		return fmt.Errorf("failed to unmarshal cosmos.directory API response: %w", err)
	}

	rawChains, ok := cdOutput["chains"]
	if !ok {
		return fmt.Errorf("failed to get chains from response: cosmos.directory API may have changed")
	}

	var chains []chainregistry.Chain
	if err := json.Unmarshal(rawChains, &chains); err != nil {
		return fmt.Errorf("failed to unmarshal chains: %w", err)
	}

	for _, c := range chains {
		r.Chains[c.ChainName] = c
	}

	return nil
}

// EnrichChain fetches the full chain information from the cosmos.directory API
func EnrichChain(chain *chainregistry.Chain) error {
	baseURL := fmt.Sprintf("%s/%s", cosmosDirectoryAPIURL, chain.ChainName)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch chains: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var cdOutput map[string]json.RawMessage
	if err := json.Unmarshal(body, &cdOutput); err != nil {
		return fmt.Errorf("failed to unmarshal cosmos.directory API response: %w", err)
	}

	rawChain, ok := cdOutput["chain"]
	if !ok {
		return fmt.Errorf("failed to get chains from response: cosmos.directory API may have changed")
	}

	if err := json.Unmarshal(rawChain, chain); err != nil {
		return fmt.Errorf("failed to unmarshal %s chain: %w", chain.ChainName, err)
	}

	return nil
}
