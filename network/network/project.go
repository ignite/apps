package network

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/cli/v28/ignite/pkg/events"
	projecttypes "github.com/ignite/network/x/project/types"

	"github.com/ignite/apps/network/network/networktypes"
)

type (
	// Prop updates project proposal.
	Prop func(*updateProp)

	// updateProp represents the update project proposal.
	updateProp struct {
		name        string
		metadata    []byte
		totalSupply sdk.Coins
	}
)

// WithProjectName provides a name proposal to update the project.
func WithProjectName(name string) Prop {
	return func(c *updateProp) {
		c.name = name
	}
}

// WithProjectMetadata provides a meta data proposal to update the project.
func WithProjectMetadata(metadata string) Prop {
	return func(c *updateProp) {
		c.metadata = []byte(metadata)
	}
}

// WithProjectTotalSupply provides a total supply proposal to update the project.
func WithProjectTotalSupply(totalSupply sdk.Coins) Prop {
	return func(c *updateProp) {
		c.totalSupply = totalSupply
	}
}

// Project fetches the project from Network.
func (n Network) Project(ctx context.Context, projectID uint64) (networktypes.Project, error) {
	n.ev.Send("Fetching project information", events.ProgressStart())
	res, err := n.projectQuery.GetProject(ctx, &projecttypes.QueryGetProjectRequest{
		ProjectId: projectID,
	})
	if isNotFoundErr(err) {
		return networktypes.Project{}, ErrObjectNotFound
	} else if err != nil {
		return networktypes.Project{}, err
	}
	return networktypes.ToProject(res.Project), nil
}

// Projects fetches the projects from Network.
func (n Network) Projects(ctx context.Context) ([]networktypes.Project, error) {
	var projects []networktypes.Project

	n.ev.Send("Fetching projects information", events.ProgressStart())
	res, err := n.projectQuery.ListProject(ctx, &projecttypes.QueryAllProjectRequest{})
	if err != nil {
		return projects, err
	}

	// Parse fetched projects
	for _, project := range res.Project {
		projects = append(projects, networktypes.ToProject(project))
	}

	return projects, nil
}

// CreateProject creates a project in Network.
func (n Network) CreateProject(ctx context.Context, name, metadata string, totalSupply sdk.Coins) (int64, error) {
	n.ev.Send(fmt.Sprintf("Creating project %s", name), events.ProgressStart())
	addr, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return 0, err
	}

	msgCreateProject := projecttypes.NewMsgCreateProject(
		addr,
		name,
		totalSupply,
		[]byte(metadata),
	)
	res, err := n.cosmos.BroadcastTx(ctx, n.account, msgCreateProject)
	if err != nil {
		return 0, err
	}

	var createProjectRes projecttypes.MsgCreateProjectResponse
	if err := res.Decode(&createProjectRes); err != nil {
		return 0, err
	}

	return int64(createProjectRes.ProjectId), nil
}

// InitializeMainnet Initialize the mainnet of the project.
func (n Network) InitializeMainnet(
	ctx context.Context,
	projectID uint64,
	sourceURL,
	sourceHash string,
	mainnetChainID string,
) (uint64, error) {
	n.ev.Send("Initializing the mainnet project", events.ProgressStart())
	addr, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return 0, err
	}

	msg := projecttypes.NewMsgInitializeMainnet(
		addr,
		projectID,
		sourceURL,
		sourceHash,
		mainnetChainID,
	)

	res, err := n.cosmos.BroadcastTx(ctx, n.account, msg)
	if err != nil {
		return 0, err
	}

	var initMainnetRes projecttypes.MsgInitializeMainnetResponse
	if err := res.Decode(&initMainnetRes); err != nil {
		return 0, err
	}

	n.ev.Send(fmt.Sprintf("Project %d initialized on mainnet", projectID), events.ProgressFinish())

	return initMainnetRes.MainnetId, nil
}

// UpdateProject updates the project name or metadata.
func (n Network) UpdateProject(
	ctx context.Context,
	id uint64,
	props ...Prop,
) error {
	// Apply the options provided by the user
	p := updateProp{}
	for _, apply := range props {
		apply(&p)
	}

	n.ev.Send(fmt.Sprintf("Updating the project %d", id), events.ProgressStart())
	account, err := n.account.Address(networktypes.SPN)
	if err != nil {
		return err
	}

	msgs := make([]sdk.Msg, 0)
	if p.name != "" || len(p.metadata) > 0 {
		msgs = append(msgs, projecttypes.NewMsgEditProject(
			account,
			id,
			p.name,
			p.metadata,
		))
	}
	if !p.totalSupply.Empty() {
		msgs = append(msgs, projecttypes.NewMsgUpdateTotalSupply(
			account,
			id,
			p.totalSupply,
		))
	}

	if _, err := n.cosmos.BroadcastTx(ctx, n.account, msgs...); err != nil {
		return err
	}
	n.ev.Send(fmt.Sprintf("Project %d updated", id), events.ProgressFinish())
	return nil
}
