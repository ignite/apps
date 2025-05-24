package networktypes

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	profiletypes "github.com/ignite/network/x/profile/types"
	projecttypes "github.com/ignite/network/x/project/types"
)

// Validator represents the Validator profile on SPN.
type Validator struct {
	Address           string   `json:"Address"`
	OperatorAddresses []string `json:"OperatorAddresses"`
	Identity          string   `json:"Identity"`
	Website           string   `json:"Website"`
	Details           string   `json:"Details"`
	Moniker           string   `json:"Moniker"`
	SecurityContact   string   `json:"SecurityContact"`
}

func (v Validator) ToProfile(
	projectID uint64,
	vouchers sdk.Coins,
	shares projecttypes.Shares,
) Profile {
	return Profile{
		ProjectID:       projectID,
		Address:         v.Address,
		Identity:        v.Identity,
		Website:         v.Website,
		Details:         v.Details,
		Moniker:         v.Moniker,
		SecurityContact: v.SecurityContact,
		Vouchers:        vouchers,
		Shares:          shares,
	}
}

// ToValidator converts a Validator data from SPN and returns a Validator object.
func ToValidator(val profiletypes.Validator) Validator {
	return Validator{
		Address:           val.Address,
		OperatorAddresses: val.OperatorAddresses,
		Identity:          val.Description.Identity,
		Website:           val.Description.Website,
		Details:           val.Description.Details,
		Moniker:           val.Description.Moniker,
		SecurityContact:   val.Description.SecurityContact,
	}
}

// Coordinator represents the Coordinator profile on SPN.
type Coordinator struct {
	CoordinatorID uint64 `json:"ID"`
	Address       string `json:"Address"`
	Active        bool   `json:"Active"`
	Identity      string `json:"Identity"`
	Website       string `json:"Website"`
	Details       string `json:"Details"`
}

func (c Coordinator) ToProfile(
	projectID uint64,
	vouchers sdk.Coins,
	shares projecttypes.Shares,
) Profile {
	return Profile{
		ProjectID: projectID,
		Address:   c.Address,
		Identity:  c.Identity,
		Website:   c.Website,
		Details:   c.Details,
		Vouchers:  vouchers,
		Shares:    shares,
	}
}

// ToCoordinator converts a Coordinator data from SPN and returns a Coordinator object.
func ToCoordinator(coord profiletypes.Coordinator) Coordinator {
	return Coordinator{
		CoordinatorID: coord.CoordinatorId,
		Address:       coord.Address,
		Active:        coord.Active,
		Identity:      coord.Description.Identity,
		Website:       coord.Description.Website,
		Details:       coord.Description.Details,
	}
}

type (
	// ChainShare represents the share of a chain on SPN.
	ChainShare struct {
		LaunchID uint64    `json:"LaunchID"`
		Shares   sdk.Coins `json:"Shares"`
	}

	// Profile represents the address profile on SPN.
	Profile struct {
		Address         string              `json:"Address"`
		ProjectID       uint64              `json:"ProjectID,omitempty"`
		Identity        string              `json:"Identity,omitempty"`
		Website         string              `json:"Website,omitempty"`
		Details         string              `json:"Details,omitempty"`
		Moniker         string              `json:"Moniker,omitempty"`
		SecurityContact string              `json:"SecurityContact,omitempty"`
		Vouchers        sdk.Coins           `json:"Vouchers,omitempty"`
		Shares          projecttypes.Shares `json:"Shares,omitempty"`
	}

	// ProfileAcc represents the address profile method interface.
	ProfileAcc interface {
		ToProfile(
			projectID uint64,
			vouchers sdk.Coins,
			shares projecttypes.Shares,
		) Profile
	}
)
