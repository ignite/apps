package networktypes

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	projecttypes "github.com/tendermint/spn/x/project/types"
)

// Project represents the project of a chain on SPN.
type Project struct {
	ID                 uint64    `json:"ID"`
	Name               string    `json:"Name"`
	CoordinatorID      uint64    `json:"CoordinatorID"`
	MainnetID          uint64    `json:"MainnetID"`
	MainnetInitialized bool      `json:"MainnetInitialized"`
	TotalSupply        sdk.Coins `json:"TotalSupply"`
	AllocatedShares    string    `json:"AllocatedShares"`
	Metadata           string    `json:"Metadata"`
}

// ToProject converts a project data from SPN and returns a Project object.
func ToProject(project projecttypes.Project) Project {
	return Project{
		ID:                 project.ProjectID,
		Name:               project.ProjectName,
		CoordinatorID:      project.CoordinatorID,
		MainnetID:          project.MainnetID,
		MainnetInitialized: project.MainnetInitialized,
		TotalSupply:        project.TotalSupply,
		AllocatedShares:    project.AllocatedShares.String(),
		Metadata:           string(project.Metadata),
	}
}

// MainnetAccount represents the project mainnet account of a chain on SPN.
type MainnetAccount struct {
	Address string              `json:"Address"`
	Shares  projecttypes.Shares `json:"Shares"`
}

// ToMainnetAccount converts a mainnet account data from SPN and returns a MainnetAccount object.
func ToMainnetAccount(acc projecttypes.MainnetAccount) MainnetAccount {
	return MainnetAccount{
		Address: acc.Address,
		Shares:  acc.Shares,
	}
}

// ProjectChains represents the chains of a project on SPN.
type ProjectChains struct {
	ProjectID uint64   `json:"ProjectID"`
	Chains    []uint64 `json:"Chains"`
}

// ToProjectChains converts a project chains data from SPN and returns a ProjectChains object.
func ToProjectChains(c projecttypes.ProjectChains) ProjectChains {
	return ProjectChains{
		ProjectID: c.ProjectID,
		Chains:    c.Chains,
	}
}
