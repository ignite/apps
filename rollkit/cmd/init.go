package cmd

import (
	"os"
	"path/filepath"

	"github.com/cometbft/cometbft/crypto"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	cmprivval "github.com/cometbft/cometbft/privval"
	cmttypes "github.com/cometbft/cometbft/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v28/ignite/services/chain"
)

func NewRollkitInit() *cobra.Command {
	c := &cobra.Command{
		Use:   "init",
		Short: "Init rollkit support",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New()
			defer session.End()

			appPath, err := cmd.Flags().GetString(flagPath)
			if err != nil {
				return err
			}
			absPath, err := filepath.Abs(appPath)
			if err != nil {
				return err
			}

			rollkitChain, err := chain.New(absPath)
			if err != nil {
				return err
			}

			if err := rollkitChain.Init(cmd.Context(), chain.InitArgsAll); err != nil {
				return err
			}

			home, err := rollkitChain.Home()
			if err != nil {
				return err
			}

			// modify genesis (add sequencer)
			genesisPath, err := rollkitChain.GenesisPath()
			if err != nil {
				return err
			}

			genesis, err := genutiltypes.AppGenesisFromFile(genesisPath)
			if err != nil {
				return err
			}

			pubKey, err := getPubKey(home)
			if err != nil {
				return err
			}

			genesis.Consensus.Validators = []cmttypes.GenesisValidator{
				{
					Name:    "Rollkit Sequencer",
					Address: pubKey.Address(),
					PubKey:  pubKey,
					Power:   1000,
				},
			}

			if err := genesis.SaveAs(genesisPath); err != nil {
				return err
			}

			return session.Printf("🗃  Initialized. Checkout your rollkit chain's home (data) directory: %s\n", colors.Info(home))
		},
	}

	c.Flags().StringP(flagPath, "p", ".", "path of the app")

	return c
}

// getPubKey returns the validator public key.
func getPubKey(chainHome string) (crypto.PubKey, error) {
	keyFilePath := filepath.Join(chainHome, "config", "priv_validator_key.json")
	keyJSONBytes, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}
	var pvKey cmprivval.FilePVKey
	err = cmtjson.Unmarshal(keyJSONBytes, &pvKey)
	if err != nil {
		return nil, errors.Errorf("error reading PrivValidator key from %v: %s", keyFilePath, err)
	}
	return pvKey.PubKey, nil
}
