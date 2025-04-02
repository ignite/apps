package network

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	launchtypes "github.com/ignite/network/x/launch/types"
	profiletypes "github.com/ignite/network/x/profile/types"
	projecttypes "github.com/ignite/network/x/project/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ignite/apps/network/network/networktypes"
	"github.com/ignite/apps/network/network/testutil"
)

var metadata = []byte(`{"cli":{"version":"1"}}`)

func startGenesisTestServer(filepath string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile(filepath)
		if err != nil {
			panic(err)
		}
		if _, err = w.Write(file); err != nil {
			panic(err)
		}
	}))
}

func startInvalidJSONServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
}

func TestPublish(t *testing.T) {
	t.Run("publish chain without project", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"GetCoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				Coordinator: profiletypes.Coordinator{
					Address:       addr,
					CoordinatorId: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainId: testutil.ChainID,
					SourceUrl:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasProject:     false,
					ProjectId:      0,
					Metadata:       metadata,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchId: testutil.LaunchID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, projectID, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, int64(-1), projectID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with custom account balance", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		accountBalance, err := sdk.ParseCoinsNormalized("1000foo,500bar")
		require.NoError(t, err)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"GetCoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				Coordinator: profiletypes.Coordinator{
					Address:       addr,
					CoordinatorId: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainId: testutil.ChainID,
					SourceUrl:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasProject:     false,
					ProjectId:      0,
					AccountBalance: accountBalance,
					Metadata:       metadata,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchId: testutil.LaunchID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, projectID, publishError := network.Publish(
			context.Background(),
			suite.ChainMock,
			WithAccountBalance(accountBalance),
		)
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, int64(-1), projectID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with pre created project", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"GetCoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				Coordinator: profiletypes.Coordinator{
					Address:       addr,
					CoordinatorId: 1,
				},
			}, nil).
			Once()
		suite.ProjectQueryMock.
			On(
				"GetProject",
				context.Background(),
				&projecttypes.QueryGetProjectRequest{
					ProjectId: testutil.ProjectID,
				},
			).
			Return(nil, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainId: testutil.ChainID,
					SourceUrl:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasProject:     true,
					ProjectId:      testutil.ProjectID,
					Metadata:       metadata,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchId: testutil.LaunchID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, projectID, publishError := network.Publish(context.Background(), suite.ChainMock, WithProject(int64(testutil.ProjectID)))
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, testutil.ProjectID, uint64(projectID))
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with a pre created project with shares", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"GetCoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				Coordinator: profiletypes.Coordinator{
					Address:       addr,
					CoordinatorId: 1,
				},
			}, nil).
			Once()
		suite.ProjectQueryMock.
			On(
				"GetProject",
				context.Background(),
				&projecttypes.QueryGetProjectRequest{
					ProjectId: testutil.ProjectID,
				},
			).
			Return(nil, nil).
			Once()
		suite.ProjectQueryMock.
			On(
				"TotalShares",
				context.Background(),
				&projecttypes.QueryTotalSharesRequest{},
			).
			Return(&projecttypes.QueryTotalSharesResponse{
				TotalShares: 100000,
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				projecttypes.NewMsgMintVouchers(
					addr,
					testutil.ProjectID,
					projecttypes.NewSharesFromCoins(sdk.NewCoins(sdk.NewInt64Coin("foo", 2000), sdk.NewInt64Coin("staking", 50000))),
				),
			).
			Return(testutil.NewResponse(&projecttypes.MsgMintVouchersResponse{}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainId: testutil.ChainID,
					SourceUrl:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasProject:     true,
					ProjectId:      testutil.ProjectID,
					Metadata:       metadata,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchId: testutil.LaunchID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, projectID, publishError := network.Publish(context.Background(), suite.ChainMock, WithProject(int64(testutil.ProjectID)),
			WithPercentageShares([]SharePercent{
				SampleSharePercent(t, "foo", 2, 100),
				SampleSharePercent(t, "staking", 50, 100),
			}),
		)
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, testutil.ProjectID, uint64(projectID))
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with custom genesis url", func(t *testing.T) {
		var (
			account              = testutil.NewTestAccount(t, testutil.TestAccountName)
			customGenesisChainID = "test-custom-1"
			customGenesisHash    = "86167654c1af18c801837d443563fd98b3fe5e8d337e70faad181cdf2100da52"
			gts                  = startGenesisTestServer("mocks/data/genesis.json")
			suite, network       = newSuite(account)
		)
		defer gts.Close()

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"GetCoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				Coordinator: profiletypes.Coordinator{
					Address:       addr,
					CoordinatorId: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainId: customGenesisChainID,
					SourceUrl:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewGenesisURL(
						gts.URL,
						customGenesisHash,
					),
					HasProject: false,
					ProjectId:  0,
					Metadata:   metadata,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchId: testutil.LaunchID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, projectID, publishError := network.Publish(
			context.Background(),
			suite.ChainMock,
			WithCustomGenesisURL(gts.URL),
		)
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, int64(-1), projectID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with custom chain id", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"GetCoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				Coordinator: profiletypes.Coordinator{
					Address:       addr,
					CoordinatorId: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainId: testutil.ChainID,
					SourceUrl:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasProject:     false,
					ProjectId:      0,
					Metadata:       metadata,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchId: testutil.LaunchID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, projectID, publishError := network.Publish(context.Background(), suite.ChainMock, WithChainID(testutil.ChainID))
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, int64(-1), projectID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with custom genesis config", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"GetCoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				Coordinator: profiletypes.Coordinator{
					Address:       addr,
					CoordinatorId: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainId: testutil.ChainID,
					SourceUrl:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewGenesisConfig(
						testutil.ChainConfigYML,
					),
					HasProject: false,
					ProjectId:  0,
					Metadata:   metadata,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchId: testutil.LaunchID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, projectID, publishError := network.Publish(
			context.Background(),
			suite.ChainMock,
			WithCustomGenesisConfig(testutil.ChainConfigYML),
		)
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, int64(-1), projectID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with custom chain id", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"GetCoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				Coordinator: profiletypes.Coordinator{
					Address:       addr,
					CoordinatorId: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainId: testutil.ChainID,
					SourceUrl:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasProject:     false,
					ProjectId:      0,
					Metadata:       metadata,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchId: testutil.LaunchID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, projectID, publishError := network.Publish(context.Background(), suite.ChainMock, WithChainID(testutil.ChainID))
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, int64(-1), projectID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with mainnet", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			gts            = startGenesisTestServer("mocks/data/genesis.json")
			suite, network = newSuite(account)
		)
		defer gts.Close()

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"GetCoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				Coordinator: profiletypes.Coordinator{
					Address:       addr,
					CoordinatorId: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&projecttypes.MsgCreateProject{
					Coordinator: addr,
					ProjectName: testutil.ChainName,
					Metadata:    []byte{},
				},
			).
			Return(testutil.NewResponse(&projecttypes.MsgCreateProjectResponse{
				ProjectId: testutil.ProjectID,
			}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&projecttypes.MsgInitializeMainnet{
					Coordinator:    addr,
					ProjectId:      testutil.ProjectID,
					SourceUrl:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					MainnetChainId: testutil.ChainID,
				},
			).
			Return(testutil.NewResponse(&projecttypes.MsgInitializeMainnetResponse{
				MainnetId: testutil.MainnetID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("Name").Return(testutil.ChainName).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, projectID, publishError := network.Publish(context.Background(), suite.ChainMock, Mainnet())
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, testutil.ProjectID, uint64(projectID))
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain with mainnet, failed to initialize mainnet", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to initialize mainnet")
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"GetCoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				Coordinator: profiletypes.Coordinator{
					Address:       addr,
					CoordinatorId: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&projecttypes.MsgCreateProject{
					Coordinator: addr,
					ProjectName: testutil.ChainName,
					Metadata:    []byte{},
				},
			).
			Return(testutil.NewResponse(&projecttypes.MsgCreateProjectResponse{
				ProjectId: testutil.ProjectID,
			}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&projecttypes.MsgInitializeMainnet{
					Coordinator:    addr,
					ProjectId:      testutil.ProjectID,
					SourceUrl:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					MainnetChainId: testutil.ChainID,
				},
			).
			Return(testutil.NewResponse(&projecttypes.MsgInitializeMainnetResponse{
				MainnetId: testutil.MainnetID,
			}), expectedError).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("Name").Return(testutil.ChainName).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock, Mainnet())
		require.Error(t, publishError)
		require.ErrorIs(t, publishError, expectedError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain with custom genesis, failed to parse custom genesis", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			gts            = startInvalidJSONServer()
			expectedError  = errors.New("JSON field not found")
		)
		defer gts.Close()

		_, _, publishError := network.Publish(
			context.Background(),
			suite.ChainMock,
			WithCustomGenesisURL(gts.URL),
		)
		require.Error(t, publishError)
		require.Equal(t, expectedError.Error(), publishError.Error())
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with coordinator creation", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On("GetCoordinatorByAddress", mock.Anything, &profiletypes.QueryGetCoordinatorByAddressRequest{
				Address: addr,
			}).
			Return(nil, sdkerrors.ErrNotFound).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainId: testutil.ChainID,
					SourceUrl:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasProject:     false,
					ProjectId:      0,
					Metadata:       metadata,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchId: testutil.LaunchID,
			}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&profiletypes.MsgCreateCoordinator{
					Address: addr,
				},
			).
			Return(testutil.NewResponse(&profiletypes.MsgCreateCoordinatorResponse{}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, projectID, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, int64(-1), projectID)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain, failed to fetch coordinator profile", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to fetch coordinator")
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On("GetCoordinatorByAddress", mock.Anything, &profiletypes.QueryGetCoordinatorByAddressRequest{
				Address: addr,
			}).
			Return(nil, expectedError).
			Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.Error(t, publishError)
		require.ErrorIs(t, publishError, expectedError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain, failed to read chain id", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to get chainID")
		)

		suite.ChainMock.
			On("ChainID").
			Return("", expectedError).
			Once()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.Error(t, publishError)
		require.ErrorIs(t, publishError, expectedError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain, failed to fetch existed project", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"GetCoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				Coordinator: profiletypes.Coordinator{
					Address:       addr,
					CoordinatorId: 1,
				},
			}, nil).
			Once()
		suite.ProjectQueryMock.
			On("GetProject", mock.Anything, &projecttypes.QueryGetProjectRequest{
				ProjectId: testutil.ProjectID,
			}).
			Return(nil, sdkerrors.ErrNotFound).
			Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock, WithProject(int64(testutil.ProjectID)))
		require.Error(t, publishError)
		require.ErrorIs(t, publishError, sdkerrors.ErrNotFound)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain, failed to create chain", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to create chain")
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"GetCoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				Coordinator: profiletypes.Coordinator{
					Address:       addr,
					CoordinatorId: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainId: testutil.ChainID,
					SourceUrl:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasProject:     false,
					ProjectId:      0,
					Metadata:       metadata,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchId: testutil.LaunchID,
			}), expectedError).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.Error(t, publishError)
		require.Equal(t, expectedError, publishError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain, failed to cache binary", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to cache binary")
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"GetCoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				Coordinator: profiletypes.Coordinator{
					Address:       addr,
					CoordinatorId: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainId: testutil.ChainID,
					SourceUrl:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasProject:     false,
					ProjectId:      0,
					Metadata:       metadata,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchId: testutil.LaunchID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.
			On("CacheBinary", testutil.LaunchID).
			Return(expectedError).
			Once()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.Error(t, publishError)
		require.ErrorIs(t, publishError, expectedError)
		suite.AssertAllMocks(t)
	})
}
