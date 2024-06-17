package cmd

import (
	"math"
	"os"
	"path/filepath"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	cmprivval "github.com/cometbft/cometbft/privval"
	cmttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	configchain "github.com/ignite/cli/v28/ignite/config/chain"
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/colors"
	"github.com/ignite/cli/v28/ignite/services/chain"
)

const defaultValPower = 1000

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

			rc, err := chain.New(absPath, chain.CollectEvents(session.EventBus()))
			if err != nil {
				return err
			}

			// use val power to set validator power in ignite config.yaml
			igniteConfig, err := rc.Config()
			if err != nil {
				return err
			}

			coins := sdk.NewCoin("stake", sdkmath.NewInt((defaultValPower * int64(math.Pow10(6)))))
			igniteConfig.Validators[0].Bonded = coins.String()
			for i, account := range igniteConfig.Accounts {
				if account.Name == igniteConfig.Validators[0].Name {
					igniteConfig.Accounts[i].Coins = []string{coins.String()}
				}
			}

			if err := configchain.Save(*igniteConfig, rc.ConfigPath()); err != nil {
				return err
			}

			if err := rc.Init(cmd.Context(), chain.InitArgsAll); err != nil {
				return err
			}

			home, err := rc.Home()
			if err != nil {
				return err
			}

			// modify genesis (add sequencer)
			genesisPath, err := rc.GenesisPath()
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
					Power:   defaultValPower,
				},
			}

			if err := genesis.SaveAs(genesisPath); err != nil {
				return err
			}

			return session.Printf("🗃 Initialized. Checkout your rollkit chain's home (data) directory: %s\n", colors.Info(home))
		},
	}

	c.Flags().StringP(flagPath, "p", ".", "path of the app")
	c.Flags().Bool("local-da", false, "this flag does nothing but is kept for backward compatibility")

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
	if err = cmtjson.Unmarshal(keyJSONBytes, &pvKey); err != nil {
		return nil, errors.Errorf("error reading PrivValidator key from %v: %s", keyFilePath, err)
	}
	return pvKey.PubKey, nil
}
