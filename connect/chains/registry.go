package chains

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ignite/cli/v29/ignite/pkg/chainregistry"
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

	apiResponseType := struct {
		Chain *chainregistry.Chain `json:"chain"`
	}{
		Chain: chain,
	}

	if err := json.Unmarshal(body, &apiResponseType); err != nil {
		return fmt.Errorf("failed to unmarshal cosmos.directory API response: %w", err)
	}

	chain.APIs.Grpc = cleanGRPCEntries(chain.APIs.Grpc)

	return nil
}

func cleanGRPCEntries(entries []chainregistry.APIProvider) []chainregistry.APIProvider {
	cleanEntries := make([]chainregistry.APIProvider, 0)
	for _, api := range entries {
		// clean-up the http(s):// prefix
		if idx := strings.Index(api.Address, "://"); idx != -1 {
			api.Address = api.Address[idx+3:]
		}
		// remove trailing slashes
		api.Address = strings.TrimSuffix(api.Address, "/")

		// remove addresses without a port
		if !strings.Contains(api.Address, ":") {
			continue
		}

		cleanEntries = append(cleanEntries, api)
	}

	return cleanEntries
}
