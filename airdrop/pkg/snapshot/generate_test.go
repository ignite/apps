package snapshot

import (
	"encoding/json"
	"testing"

	"cosmossdk.io/math"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	sampleGenesis := `{
  "app_state": {
    "auth": {
      "accounts": [
        {
          "@type": "/cosmos.auth.v1beta1.BaseAccount",
          "address": "cosmos1ph572zegdp0s8uhmqc5q4j8h2t8wtwmevte8zw",
          "pub_key": null,
          "account_number": "0",
          "sequence": "0"
        },
        {
          "@type": "/cosmos.auth.v1beta1.BaseAccount",
          "address": "cosmos1dgez7sgugzf23c6p97vhllem0lahate6gyx06n",
          "pub_key": null,
          "account_number": "0",
          "sequence": "0"
        }
      ]
    },
    "bank": {
      "balances": [
        {
          "address": "cosmos1ph572zegdp0s8uhmqc5q4j8h2t8wtwmevte8zw",
          "coins": [
            {
              "denom": "stake",
              "amount": "200000000"
            },
            {
              "denom": "token",
              "amount": "20000"
            }
          ]
        },
        {
          "address": "cosmos1dgez7sgugzf23c6p97vhllem0lahate6gyx06n",
          "coins": [
            {
              "denom": "stake",
              "amount": "100000000"
            },
            {
              "denom": "token",
              "amount": "10000"
            }
          ]
        }
      ],
      "supply": [

      ],
      "denom_metadata": [

      ]
    },
    "staking": {
      "delegations": [
        {
          "delegator_address": "cosmos1vgkzwm2f6gafr38d68fzr64lrt50jv3a0kstql",
          "shares": "1.002042131346618030",
          "validator_address": "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c"
        },
        {
          "delegator_address": "cosmos1ph572zegdp0s8uhmqc5q4j8h2t8wtwmevte8zw",
          "shares": "1.002042131346618030",
          "validator_address": "cosmosvaloper1ey69r37gfxvxg62sh4r0ktpuc46pzjrm873ae8"
        },
        {
          "delegator_address": "cosmos1dgez7sgugzf23c6p97vhllem0lahate6gyx06n",
          "shares": "41581980.000000000000000000",
          "validator_address": "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42"
        }
      ],
      "exported": false,
      "last_total_power": "0",
      "last_validator_powers": [

      ],
      "redelegations": [

      ],
      "unbonding_delegations": [
        {
          "delegator_address": "cosmos1vgkzwm2f6gafr38d68fzr64lrt50jv3a0kstql",
          "entries": [
            {
              "balance": "2000000",
              "completion_time": "2022-01-25T17:44:03.234988625Z",
              "creation_height": "8947293",
              "initial_balance": "2000000"
            },
            {
              "balance": "4985000",
              "completion_time": "2022-01-25T18:27:58.166611500Z",
              "creation_height": "8947656",
              "initial_balance": "4985000"
            }
          ],
          "validator_address": "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c"
        },
        {
          "delegator_address": "cosmos1ph572zegdp0s8uhmqc5q4j8h2t8wtwmevte8zw",
          "entries": [
            {
              "balance": "53425",
              "completion_time": "2022-01-22T18:39:34.505180406Z",
              "creation_height": "8911974",
              "initial_balance": "53425"
            }
          ],
          "validator_address": "cosmosvaloper1ey69r37gfxvxg62sh4r0ktpuc46pzjrm873ae8"
        },
        {
          "delegator_address": "cosmos1dgez7sgugzf23c6p97vhllem0lahate6gyx06n",
          "entries": [
            {
              "balance": "50000",
              "completion_time": "2022-01-28T17:51:40.681512567Z",
              "creation_height": "8983119",
              "initial_balance": "50000"
            }
          ],
          "validator_address": "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42"
        }
      ],
      "validators": [
        {
          "commission": {
            "commission_rates": {
              "max_change_rate": "1.000000000000000000",
              "max_rate": "1.000000000000000000",
              "rate": "0.030000000000000000"
            },
            "update_time": "2019-11-01T04:08:08.548659287Z"
          },
          "consensus_pubkey": {
            "@type": "/cosmos.crypto.ed25519.PubKey",
            "key": "Roh99RlsnDKHUFYUcQVHk2S84NeZfZdpc+CBb6NREhM="
          },
          "delegator_shares": "6066835864634.794010447346973671",
          "description": {
            "details": "Sunny Aggarwal (@sunnya97) and Dev Ojha (@ValarDragon)",
            "identity": "5B5AB9D8FBBCEDC6",
            "moniker": "Sikka",
            "security_contact": "",
            "website": "sikka.tech"
          },
          "jailed": false,
          "min_self_delegation": "1",
          "operator_address": "cosmosvaloper1ey69r37gfxvxg62sh4r0ktpuc46pzjrm873ae8",
          "status": "BOND_STATUS_BONDED",
          "tokens": "6065622565414",
          "unbonding_height": "7893737",
          "unbonding_time": "2021-10-25T23:53:04.980335653Z"
        },
        {
          "commission": {
            "commission_rates": {
              "max_change_rate": "0.100000000000000000",
              "max_rate": "0.200000000000000000",
              "rate": "0.020000000000000000"
            },
            "update_time": "2021-02-02T06:53:49.809692399Z"
          },
          "consensus_pubkey": {
            "@type": "/cosmos.crypto.ed25519.PubKey",
            "key": "Qajjf1kiAJ0M1UcH1TSUYLP13kgE128Av1XmGQO711c="
          },
          "delegator_shares": "5332026292167.000000000000000000",
          "description": {
            "details": "For all game enthusiasts",
            "identity": "6F3A316294AD9D0B",
            "moniker": "GAME",
            "security_contact": "",
            "website": ""
          },
          "jailed": false,
          "min_self_delegation": "1",
          "operator_address": "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42",
          "status": "BOND_STATUS_BONDED",
          "tokens": "5332026292167",
          "unbonding_height": "0",
          "unbonding_time": "1970-01-01T00:00:00Z"
        },
        {
          "commission": {
            "commission_rates": {
              "max_change_rate": "0.020000000000000000",
              "max_rate": "0.200000000000000000",
              "rate": "0.010000000000000000"
            },
            "update_time": "2021-11-07T12:17:47.672397912Z"
          },
          "consensus_pubkey": {
            "@type": "/cosmos.crypto.ed25519.PubKey",
            "key": "hUDXospsiB6oJVvkRVB2IyanCHs5hiaeqoEWzp9be8w="
          },
          "delegator_shares": "5630381252757.000000000000000000",
          "description": {
            "details": "SG-1 - Your favorite team on Cosmos. We refund downtime slashing to 100%",
            "identity": "48608633F99D1B60",
            "moniker": "SG-1",
            "security_contact": "",
            "website": "https://sg-1.online"
          },
          "jailed": false,
          "min_self_delegation": "1000",
          "operator_address": "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c",
          "status": "BOND_STATUS_BONDED",
          "tokens": "5630381252757",
          "unbonding_height": "0",
          "unbonding_time": "1970-01-01T00:00:00Z"
        }
      ]
    }
  }
}`
	var sampleDoc tmtypes.GenesisDoc
	err := tmjson.Unmarshal([]byte(sampleGenesis), &sampleDoc)
	require.NoError(t, err)

	var sampleGenState map[string]json.RawMessage
	err = json.Unmarshal(sampleDoc.AppState, &sampleGenState)
	require.NoError(t, err)

	var (
		acc1 = "cosmos1ph572zegdp0s8uhmqc5q4j8h2t8wtwmevte8zw"
		acc2 = "cosmos1dgez7sgugzf23c6p97vhllem0lahate6gyx06n"
		acc3 = "cosmos1vgkzwm2f6gafr38d68fzr64lrt50jv3a0kstql"
	)

	invalidGenStateData := `{"invalid_app":{"invalid_key":"invalid_value","invalid_numer":10}}`
	var invalidGenState map[string]json.RawMessage
	err = json.Unmarshal([]byte(invalidGenStateData), &invalidGenState)
	require.NoError(t, err)

	tests := []struct {
		name     string
		genState map[string]json.RawMessage
		want     Snapshot
		err      error
	}{
		{
			name:     "valid genesis state",
			genState: sampleGenState,
			want: Snapshot{
				NumberAccounts: 3,
				Accounts: Accounts{
					acc1: {
						Address:        acc1,
						Staked:         math.NewInt(1),
						UnbondingStake: math.NewInt(53425),
						Balance: sdk.NewCoins(
							sdk.NewCoin("stake", math.NewInt(200000000)),
							sdk.NewCoin("token", math.NewInt(20000)),
						),
					},
					acc2: {
						Address:        acc2,
						Staked:         math.NewInt(41581980),
						UnbondingStake: math.NewInt(50000),
						Balance: sdk.NewCoins(
							sdk.NewCoin("stake", math.NewInt(100000000)),
							sdk.NewCoin("token", math.NewInt(10000)),
						),
					},
					acc3: {
						Address:        acc3,
						Staked:         math.NewInt(1),
						UnbondingStake: math.NewInt(6985000),
						Balance:        sdk.NewCoins(),
					},
				},
			},
		},
		{
			name:     "invalid genesis state",
			genState: invalidGenState,
			want:     Snapshot{NumberAccounts: 0, Accounts: Accounts{}},
		},
		{
			name: "invalid genesis data",
			genState: map[string]json.RawMessage{
				"invalid": []byte("invalid"),
			},
			want: Snapshot{NumberAccounts: 0, Accounts: Accounts{}},
		},
		{
			name:     "nil genesis state",
			genState: nil,
			want:     Snapshot{NumberAccounts: 0, Accounts: Accounts{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Generate(tt.genState)
			if tt.err != nil {
				require.Error(t, err)
				require.ErrorIs(t, tt.err, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
