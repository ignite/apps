package networktypes

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	launchtypes "github.com/ignite/network/x/launch/types"
	"github.com/pkg/errors"
)

// GenesisInformation represents all information for a chain to construct the genesis.
// This structure indexes accounts and validators by their address for better performance.
type GenesisInformation struct {
	// make sure to use slices for the following because slices are ordered.
	// they later used to create a Genesis so, having them ordered is important to
	// be able to produce a deterministic Genesis.

	GenesisAccounts   []GenesisAccount
	VestingAccounts   []VestingAccount
	GenesisValidators []GenesisValidator
	ParamChanges      []ParamChange
}

// GenesisAccount represents an account with initial coin allocation for the chain for the chain genesis.
type GenesisAccount struct {
	Address string    `json:"Address,omitempty"`
	Coins   sdk.Coins `json:"Coins,omitempty"`
}

// VestingAccount represents a vesting account with initial coin allocation  and vesting option for the chain genesis.
// VestingAccount supports currently only delayed vesting option.
type VestingAccount struct {
	Address      string    `json:"Address,omitempty"`
	TotalBalance sdk.Coins `json:"TotalBalance,omitempty"`
	Vesting      sdk.Coins `json:"Vesting,omitempty"`
	EndTime      int64     `json:"EndTime,omitempty"`
}

// GenesisValidator represents a genesis validator associated with a gentx in the chain genesis.
type GenesisValidator struct {
	Address        string           `json:"Address,omitempty"`
	Gentx          []byte           `json:"Gentx,omitempty"`
	Peer           launchtypes.Peer `json:"Peer,omitempty"`
	SelfDelegation sdk.Coin         `json:"SelfDelegation,omitempty"`
}

// ParamChange represents a parameter change to be applied to the chain genesis.
type ParamChange struct {
	Module string `json:"Module,omitempty"`
	Param  string `json:"Param,omitempty"`
	Value  []byte `json:"Value,omitempty"`
}

// ToGenesisAccount converts genesis account from SPN.
func ToGenesisAccount(acc launchtypes.GenesisAccount) GenesisAccount {
	return GenesisAccount{
		Address: acc.Address,
		Coins:   acc.Coins,
	}
}

// ToVestingAccount converts vesting account from SPN.
func ToVestingAccount(acc launchtypes.VestingAccount) (VestingAccount, error) {
	delayedVesting := acc.VestingOptions.GetDelayedVesting()
	if delayedVesting == nil {
		return VestingAccount{}, errors.New("only delayed vesting option is supported")
	}

	return VestingAccount{
		Address:      acc.Address,
		TotalBalance: delayedVesting.TotalBalance,
		Vesting:      delayedVesting.Vesting,
		EndTime:      delayedVesting.EndTime.Unix(),
	}, nil
}

// ToGenesisValidator converts genesis validator from SPN.
func ToGenesisValidator(val launchtypes.GenesisValidator) GenesisValidator {
	return GenesisValidator{
		Address:        val.Address,
		Gentx:          val.GenTx,
		Peer:           val.Peer,
		SelfDelegation: val.SelfDelegation,
	}
}

// ToParamChange converts param change from SPN.
func ToParamChange(pc launchtypes.ParamChange) ParamChange {
	return ParamChange{
		Param:  pc.Param,
		Module: pc.Module,
		Value:  pc.Value,
	}
}

// NewGenesisInformation initializes  new GenesisInformation.
func NewGenesisInformation(
	genAccs []GenesisAccount,
	vestingAccs []VestingAccount,
	genVals []GenesisValidator,
	paramChanges []ParamChange,
) (gi GenesisInformation) {
	return GenesisInformation{
		GenesisAccounts:   genAccs,
		VestingAccounts:   vestingAccs,
		GenesisValidators: genVals,
		ParamChanges:      paramChanges,
	}
}

// ContainsGenesisAccount returns true if GenesisInformation contains given address.
// Returns index if true, -1 if false.
func (gi GenesisInformation) ContainsGenesisAccount(address string) (bool, int) {
	for i, account := range gi.GenesisAccounts {
		if account.Address == address {
			return true, i
		}
	}
	return false, -1
}

// ContainsVestingAccount returns true if GenesisInformation contains given address.
// Returns index if true, -1 if false.
func (gi GenesisInformation) ContainsVestingAccount(address string) (bool, int) {
	for i, account := range gi.VestingAccounts {
		if account.Address == address {
			return true, i
		}
	}
	return false, -1
}

// ContainsGenesisValidator returns true if GenesisInformation contains given address.
// Returns index if true, -1 if false.
func (gi GenesisInformation) ContainsGenesisValidator(address string) (bool, int) {
	for i, account := range gi.GenesisValidators {
		if account.Address == address {
			return true, i
		}
	}
	return false, -1
}

// ContainsParamChange returns true if GenesisInformation contains given module-param pair.
// Returns index if true, -1 if false.
func (gi GenesisInformation) ContainsParamChange(module, param string) (bool, int) {
	for i, paramChange := range gi.ParamChanges {
		if paramChange.Module == module && paramChange.Param == param {
			return true, i
		}
	}
	return false, -1
}

func (gi *GenesisInformation) AddGenesisAccount(acc GenesisAccount) {
	gi.GenesisAccounts = append(gi.GenesisAccounts, acc)
}

func (gi *GenesisInformation) AddVestingAccount(acc VestingAccount) {
	gi.VestingAccounts = append(gi.VestingAccounts, acc)
}

func (gi *GenesisInformation) AddGenesisValidator(val GenesisValidator) {
	gi.GenesisValidators = append(gi.GenesisValidators, val)
}

func (gi *GenesisInformation) RemoveGenesisAccount(address string) {
	for i, account := range gi.GenesisAccounts {
		if account.Address == address {
			gi.GenesisAccounts = append(gi.GenesisAccounts[:i], gi.GenesisAccounts[i+1:]...)
		}
	}
}

func (gi *GenesisInformation) RemoveVestingAccount(address string) {
	for i, account := range gi.VestingAccounts {
		if account.Address == address {
			gi.VestingAccounts = append(gi.VestingAccounts[:i], gi.VestingAccounts[i+1:]...)
		}
	}
}

func (gi *GenesisInformation) RemoveGenesisValidator(address string) {
	for i, account := range gi.GenesisValidators {
		if account.Address == address {
			gi.GenesisValidators = append(gi.GenesisValidators[:i], gi.GenesisValidators[i+1:]...)
		}
	}
}

// AddParamChange adds a ParamChange to the GenesisInformation.
// Appends if entry does not exist.  Updates if it already exists.
func (gi *GenesisInformation) AddParamChange(pc ParamChange) {
	contains, index := gi.ContainsParamChange(pc.Module, pc.Param)
	if contains {
		gi.ParamChanges[index] = pc
		return
	}
	gi.ParamChanges = append(gi.ParamChanges, pc)
}

// ApplyRequest applies to the genesisInformation the changes implied by the approval of a request.
func (gi GenesisInformation) ApplyRequest(request Request) (GenesisInformation, error) {
	switch requestContent := request.Content.Content.(type) {
	case *launchtypes.RequestContent_GenesisAccount:
		// new genesis account in the genesis
		ga := ToGenesisAccount(*requestContent.GenesisAccount)
		genExist, _ := gi.ContainsGenesisAccount(ga.Address)
		vestingExist, _ := gi.ContainsVestingAccount(ga.Address)
		if genExist || vestingExist {
			return gi, NewWrappedErrInvalidRequest(request.RequestID, "genesis account already in genesis")
		}
		gi.AddGenesisAccount(ga)

	case *launchtypes.RequestContent_VestingAccount:
		// new vesting account in the genesis
		va, err := ToVestingAccount(*requestContent.VestingAccount)
		if err != nil {
			// we don't treat this error as errInvalidRequests
			// because it can occur if we don't support this format of vesting account
			// but the request is still correct
			return gi, err
		}

		genExist, _ := gi.ContainsGenesisAccount(va.Address)
		vestingExist, _ := gi.ContainsVestingAccount(va.Address)
		if genExist || vestingExist {
			return gi, NewWrappedErrInvalidRequest(request.RequestID, "vesting account already in genesis")
		}
		gi.AddVestingAccount(va)

	case *launchtypes.RequestContent_AccountRemoval:
		// account removed from the genesis
		ar := requestContent.AccountRemoval
		genExist, _ := gi.ContainsGenesisAccount(ar.Address)
		vestingExist, _ := gi.ContainsVestingAccount(ar.Address)
		if !genExist && !vestingExist {
			return gi, NewWrappedErrInvalidRequest(request.RequestID, "account can't be removed because it doesn't exist")
		}
		gi.RemoveGenesisAccount(ar.Address)
		gi.RemoveVestingAccount(ar.Address)

	case *launchtypes.RequestContent_GenesisValidator:
		// new genesis validator in the genesis
		gv := ToGenesisValidator(*requestContent.GenesisValidator)
		if contains, _ := gi.ContainsGenesisValidator(gv.Address); contains {
			return gi, NewWrappedErrInvalidRequest(request.RequestID, "genesis validator already in genesis")
		}
		gi.AddGenesisValidator(gv)

	case *launchtypes.RequestContent_ValidatorRemoval:
		// validator removed from the genesis
		vr := requestContent.ValidatorRemoval
		if contains, _ := gi.ContainsGenesisValidator(vr.ValAddress); !contains {
			return gi, NewWrappedErrInvalidRequest(request.RequestID, "genesis validator can't be removed because it doesn't exist")
		}

	case *launchtypes.RequestContent_ParamChange:
		// param changed in genesis file
		pc := ToParamChange(*requestContent.ParamChange)
		gi.AddParamChange(pc)
	}

	return gi, nil
}
