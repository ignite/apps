package testutil

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/ignite/apps/network/network/mocks"
)

// Suite is a mocks container, used to write less code for tests setup.
type Suite struct {
	ChainMock                *mocks.Chain
	CosmosClientMock         *mocks.CosmosClient
	LaunchQueryMock          *mocks.LaunchClient
	ProjectQueryMock         *mocks.ProjectClient
	ProfileQueryMock         *mocks.ProfileClient
	RewardClient             *mocks.RewardClient
	StakingClient            *mocks.StakingClient
	BankClient               *mocks.BankClient
	MonitoringConsumerClient *mocks.MonitoringcClient
}

// AssertAllMocks asserts all suite mocks expectations.
func (s *Suite) AssertAllMocks(t *testing.T) {
	t.Helper()
	s.ChainMock.AssertExpectations(t)
	s.ProfileQueryMock.AssertExpectations(t)
	s.LaunchQueryMock.AssertExpectations(t)
	s.CosmosClientMock.AssertExpectations(t)
	s.ProjectQueryMock.AssertExpectations(t)
	s.RewardClient.AssertExpectations(t)
	s.StakingClient.AssertExpectations(t)
	s.MonitoringConsumerClient.AssertExpectations(t)
	s.BankClient.AssertExpectations(t)
}

// NewSuite creates new suite with mocks.
func NewSuite() Suite {
	cosmos := new(mocks.CosmosClient)
	cosmos.On("Context").Return(client.Context{})
	return Suite{
		ChainMock:                new(mocks.Chain),
		CosmosClientMock:         cosmos,
		LaunchQueryMock:          new(mocks.LaunchClient),
		ProjectQueryMock:         new(mocks.ProjectClient),
		ProfileQueryMock:         new(mocks.ProfileClient),
		RewardClient:             new(mocks.RewardClient),
		StakingClient:            new(mocks.StakingClient),
		BankClient:               new(mocks.BankClient),
		MonitoringConsumerClient: new(mocks.MonitoringcClient),
	}
}
