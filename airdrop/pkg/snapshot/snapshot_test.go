package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParseSnapshot(t *testing.T) {
	var (
		accAddr1 = sdk.AccAddress(rand.Str(32)).String()
		accAddr2 = sdk.AccAddress(rand.Str(32)).String()
		accAddr3 = sdk.AccAddress(rand.Str(32)).String()
	)
	sampleConfig := Snapshot{
		NumberAccounts: 3,
		Accounts: Accounts{
			accAddr1: {
				Address:        accAddr1,
				Staked:         math.NewInt(1),
				UnbondingStake: math.NewInt(53425),
				Balance: sdk.NewCoins(
					sdk.NewCoin("stake", math.NewInt(200000000)),
					sdk.NewCoin("token", math.NewInt(20000)),
				),
			},
			accAddr2: {
				Address:        accAddr2,
				Staked:         math.NewInt(41581980),
				UnbondingStake: math.NewInt(50000),
				Balance: sdk.NewCoins(
					sdk.NewCoin("stake", math.NewInt(100000000)),
					sdk.NewCoin("token", math.NewInt(10000)),
				),
			},
			accAddr3: {
				Address:        accAddr3,
				Staked:         math.NewInt(1),
				UnbondingStake: math.NewInt(6985000),
				Balance:        sdk.NewCoins(),
			},
		},
	}
	yamlData, err := json.Marshal(&sampleConfig)
	require.NoError(t, err)
	sampleConfigPath := filepath.Join(t.TempDir(), "config.yml")
	err = os.WriteFile(sampleConfigPath, yamlData, 0o644)
	require.NoError(t, err)

	tests := []struct {
		name     string
		filename string
		want     Snapshot
		err      error
	}{
		{
			name:     "valid snapshot file",
			filename: sampleConfigPath,
			want:     sampleConfig,
		},
		{
			name:     "valid snapshot file",
			filename: "invalid_file_path",
			want:     sampleConfig,
			err:      errors.New("open invalid_file_path: no such file or directory"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSnapshot(tt.filename)
			if tt.err != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.err.Error())
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.want, got)
		})
	}
}

func TestNewAccount(t *testing.T) {
	var (
		address = sdk.AccAddress(rand.Str(10)).String()
		got     = newAccount(address)
	)
	require.Equal(t, address, got.Address)
	require.Equal(t, math.ZeroInt(), got.Staked)
	require.Equal(t, math.ZeroInt(), got.UnbondingStake)
	require.Equal(t, sdk.NewCoins(), got.Balance)
}

func TestAccounts_getAccount(t *testing.T) {
	var (
		accAddr1 = sdk.AccAddress(rand.Str(32)).String()
		accAddr2 = sdk.AccAddress(rand.Str(32)).String()
		accAddr3 = sdk.AccAddress(rand.Str(32)).String()
	)

	sampleAccounts := Accounts{
		accAddr1: {
			Address:        accAddr1,
			Staked:         math.NewInt(10),
			UnbondingStake: math.NewInt(10),
			Balance:        sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(10))),
		},
		accAddr2: {
			Address:        accAddr2,
			Staked:         math.NewInt(12),
			UnbondingStake: math.NewInt(12),
			Balance:        sdk.NewCoins(sdk.NewCoin("uatom", math.NewInt(12))),
		},
	}
	tests := []struct {
		name    string
		a       Accounts
		address string
		want    Account
	}{
		{
			name:    "already exist address 1",
			a:       sampleAccounts,
			address: accAddr1,
			want:    sampleAccounts[accAddr1],
		},
		{
			name:    "already exist address 2",
			a:       sampleAccounts,
			address: accAddr2,
			want:    sampleAccounts[accAddr2],
		},
		{
			name:    "not exist address",
			a:       sampleAccounts,
			address: accAddr3,
			want: Account{
				Address:        accAddr3,
				Staked:         math.ZeroInt(),
				UnbondingStake: math.ZeroInt(),
				Balance:        sdk.NewCoins(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.getAccount(tt.address)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestAccounts_ExcludeAddress(t *testing.T) {
	var (
		accAddr1 = sdk.AccAddress(rand.Str(32)).String()
		accAddr2 = sdk.AccAddress(rand.Str(32)).String()
		accAddr3 = sdk.AccAddress(rand.Str(32)).String()
	)

	tests := []struct {
		name    string
		accs    Accounts
		address string
		want    Accounts
	}{
		{
			name: "exclude an address",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
			address: accAddr1,
			want: Accounts{
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
		},
		{
			name: "exclude an non exist address",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
			address: sdk.AccAddress(rand.Str(32)).String(),
			want: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
		},
		{
			name: "exclude the last address",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
			},
			address: accAddr1,
			want:    Accounts{},
		},
		{
			name:    "exclude an address from a empty list",
			accs:    Accounts{},
			address: accAddr1,
			want:    Accounts{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.accs.excludeAddress(tt.address)
			require.EqualValues(t, tt.want, tt.accs)
		})
	}
}

func TestAccounts_ExcludeAddresses(t *testing.T) {
	var (
		accAddr1 = sdk.AccAddress(rand.Str(32)).String()
		accAddr2 = sdk.AccAddress(rand.Str(32)).String()
		accAddr3 = sdk.AccAddress(rand.Str(32)).String()
	)

	tests := []struct {
		name      string
		accs      Accounts
		addresses []string
		want      Accounts
	}{
		{
			name: "exclude one address",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
			addresses: []string{accAddr1},
			want: Accounts{
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
		},
		{
			name: "exclude two address",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
			addresses: []string{accAddr1, accAddr2},
			want: Accounts{
				accAddr3: {Address: accAddr3},
			},
		},
		{
			name: "exclude all address",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
			addresses: []string{accAddr1, accAddr2, accAddr3},
			want:      Accounts{},
		},
		{
			name: "exclude all address and a non exiting",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
			addresses: []string{accAddr1, accAddr2, accAddr3, sdk.AccAddress(rand.Str(32)).String()},
			want:      Accounts{},
		},
		{
			name: "exclude an non exist address",
			accs: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
			addresses: []string{sdk.AccAddress(rand.Str(32)).String()},
			want: Accounts{
				accAddr1: {Address: accAddr1},
				accAddr2: {Address: accAddr2},
				accAddr3: {Address: accAddr3},
			},
		},
		{
			name:      "exclude an address from a empty list",
			accs:      Accounts{},
			addresses: []string{sdk.AccAddress(rand.Str(32)).String()},
			want:      Accounts{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.accs.excludeAddresses(tt.addresses...)
			require.EqualValues(t, tt.want, tt.accs)
		})
	}
}

func TestAccounts_FilterDenom(t *testing.T) {
	var (
		accAddr1 = sdk.AccAddress(rand.Str(32)).String()
		accAddr2 = sdk.AccAddress(rand.Str(32)).String()
		accAddr3 = sdk.AccAddress(rand.Str(32)).String()
	)

	tests := []struct {
		name  string
		accs  Accounts
		denom string
		want  Accounts
	}{
		{
			name: "filter uatom denom",
			accs: Accounts{
				accAddr1: {
					Address: accAddr1,
					Balance: sdk.NewCoins(
						sdk.NewCoin("uatom", math.NewInt(1000)),
						sdk.NewCoin("token", math.NewInt(2000)),
						sdk.NewCoin("stake", math.NewInt(3000)),
					),
				},
				accAddr2: {
					Address: accAddr2,
					Balance: sdk.NewCoins(
						sdk.NewCoin("stake", math.NewInt(3000)),
					),
				},
				accAddr3: {
					Address: accAddr3,
					Balance: sdk.NewCoins(
						sdk.NewCoin("uatom", math.NewInt(1000)),
						sdk.NewCoin("stake", math.NewInt(3000)),
					),
				},
			},
			denom: "uatom",
			want: Accounts{
				accAddr1: {
					Address: accAddr1,
					Balance: sdk.NewCoins(
						sdk.NewCoin("uatom", math.NewInt(1000)),
					),
				},
				accAddr2: {
					Address: accAddr2,
					Balance: sdk.NewCoins(),
				},
				accAddr3: {
					Address: accAddr3,
					Balance: sdk.NewCoins(
						sdk.NewCoin("uatom", math.NewInt(1000)),
					),
				},
			},
		},
		{
			name: "filter non exiting denom",
			accs: Accounts{
				accAddr1: {
					Address: accAddr1,
					Balance: sdk.NewCoins(
						sdk.NewCoin("uatom", math.NewInt(1000)),
						sdk.NewCoin("token", math.NewInt(2000)),
						sdk.NewCoin("stake", math.NewInt(3000)),
					),
				},
				accAddr2: {
					Address: accAddr2,
					Balance: sdk.NewCoins(
						sdk.NewCoin("stake", math.NewInt(3000)),
					),
				},
			},
			denom: "void",
			want: Accounts{
				accAddr1: {
					Address: accAddr1,
					Balance: sdk.NewCoins(),
				},
				accAddr2: {
					Address: accAddr2,
					Balance: sdk.NewCoins(),
				},
			},
		},
		{
			name: "filter empty balances",
			accs: Accounts{
				accAddr1: {
					Address: accAddr1,
					Balance: sdk.NewCoins(),
				},
				accAddr2: {
					Address: accAddr2,
				},
			},
			denom: "uatom",
			want: Accounts{
				accAddr1: {
					Address: accAddr1,
					Balance: sdk.NewCoins(),
				},
				accAddr2: {
					Address: accAddr2,
					Balance: sdk.NewCoins(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.accs.filterDenom(tt.denom)
			require.EqualValues(t, tt.want, tt.accs)
		})
	}
}

func TestAccount_TotalStake(t *testing.T) {
	tests := []struct {
		name           string
		staked         math.Int
		unbondingStake math.Int
		want           math.Int
	}{
		{
			name:           "staked and unbounding stake",
			unbondingStake: math.NewInt(100),
			staked:         math.NewInt(999),
			want:           math.NewInt(1099),
		},
		{
			name:           "zero staked",
			unbondingStake: math.NewInt(0),
			staked:         math.NewInt(999),
			want:           math.NewInt(999),
		},
		{
			name:           "zero unbounding stake",
			unbondingStake: math.NewInt(100),
			staked:         math.NewInt(0),
			want:           math.NewInt(100),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Account{
				Staked:         tt.staked,
				UnbondingStake: tt.unbondingStake,
			}
			got := a.totalStake()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestAccount_BalanceAmount(t *testing.T) {
	tests := []struct {
		name           string
		balance        sdk.Coins
		staked         math.Int
		unboudingStake math.Int
		want           math.Int
	}{
		{
			name: "test 3 denom and stake",
			balance: sdk.NewCoins(
				sdk.NewCoin("uatom", math.NewInt(1000)),
				sdk.NewCoin("token", math.NewInt(2000)),
				sdk.NewCoin("stake", math.NewInt(3000)),
			),
			staked:         math.NewInt(100),
			unboudingStake: math.NewInt(0),
			want:           math.NewInt(6100),
		},
		{
			name: "test 2 denom",
			balance: sdk.NewCoins(
				sdk.NewCoin("uatom", math.NewInt(1000)),
				sdk.NewCoin("stake", math.NewInt(3000)),
			),
			staked:         math.NewInt(33),
			unboudingStake: math.NewInt(22),
			want:           math.NewInt(4055),
		},
		{
			name: "test 1 denom",
			balance: sdk.NewCoins(
				sdk.NewCoin("uatom", math.NewInt(1000)),
			),
			staked:         math.NewInt(0),
			unboudingStake: math.NewInt(12),
			want:           math.NewInt(1012),
		},
		{
			name:           "test no denom",
			balance:        sdk.NewCoins(),
			staked:         math.NewInt(1),
			unboudingStake: math.Int{},
			want:           math.NewInt(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Account{Balance: tt.balance, Staked: tt.staked, UnbondingStake: tt.unboudingStake}
			got := a.balanceAmount()
			require.Equal(t, tt.want, got)
		})
	}
}
