package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite/cli/v28/ignite/config/chain"
	"github.com/ignite/cli/v28/ignite/config/chain/base"
	v1 "github.com/ignite/cli/v28/ignite/config/chain/v1"
	pluginsconfig "github.com/ignite/cli/v28/ignite/config/plugins"
	"github.com/ignite/cli/v28/ignite/pkg/availableport"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner"
	"github.com/ignite/cli/v28/ignite/pkg/cmdrunner/step"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	yamlmap "github.com/ignite/cli/v28/ignite/pkg/yaml"
	"github.com/ignite/cli/v28/ignite/services/plugin"
	envtest "github.com/ignite/cli/v28/integration"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const (
	relayerMnemonic = "great immense still pill defense fetch pencil slow purchase symptom speed arm shoot fence have divorce cigar rapid hen vehicle pear evolve correct nerve"
)

var (
	bobName    = "bob"
	marsConfig = v1.Config{
		Config: base.Config{
			Version: 1,
			Build: base.Build{
				Proto: base.Proto{
					Path:            "proto",
					ThirdPartyPaths: []string{"third_party/proto", "proto_vendor"},
				},
			},
			Accounts: []base.Account{
				{
					Name:     "alice",
					Coins:    []string{"100000000000token", "10000000000000000000stake"},
					Mnemonic: "slide moment original seven milk crawl help text kick fluid boring awkward doll wonder sure fragile plate grid hard next casual expire okay body",
				},
				{
					Name:     "bob",
					Coins:    []string{"100000000000token", "10000000000000000000stake"},
					Mnemonic: "trap possible liquid elite embody host segment fantasy swim cable digital eager tiny broom burden diary earn hen grow engine pigeon fringe claim program",
				},
				{
					Name:     "relayer",
					Coins:    []string{"100000000000token", "1000000000000000000000stake"},
					Mnemonic: relayerMnemonic,
				},
			},
			Faucet: base.Faucet{
				Name:  &bobName,
				Coins: []string{"500token", "100000000stake"},
				Host:  ":4501",
			},
			Genesis: yamlmap.Map{"chain_id": "mars-1"},
		},
		Validators: []v1.Validator{
			{
				Name:   "alice",
				Bonded: "100000000stake",
				Client: yamlmap.Map{"keyring-backend": keyring.BackendTest},
				App: yamlmap.Map{
					"api":      yamlmap.Map{"address": ":1318"},
					"grpc":     yamlmap.Map{"address": ":9092"},
					"grpc-web": yamlmap.Map{"address": ":9093"},
				},
				Config: yamlmap.Map{
					"p2p": yamlmap.Map{"laddr": ":26658"},
					"rpc": yamlmap.Map{"laddr": ":26658", "pprof_laddr": ":6061"},
				},
				Home: "$HOME/.mars",
			},
		},
	}
	earthConfig = v1.Config{
		Config: base.Config{
			Version: 1,
			Build: base.Build{
				Proto: base.Proto{
					Path:            "proto",
					ThirdPartyPaths: []string{"third_party/proto", "proto_vendor"},
				},
			},
			Accounts: []base.Account{
				{
					Name:     "alice",
					Coins:    []string{"100000000000token", "10000000000000000000stake"},
					Mnemonic: "slide moment original seven milk crawl help text kick fluid boring awkward doll wonder sure fragile plate grid hard next casual expire okay body",
				},
				{
					Name:     "bob",
					Coins:    []string{"100000000000token", "10000000000000000000stake"},
					Mnemonic: "trap possible liquid elite embody host segment fantasy swim cable digital eager tiny broom burden diary earn hen grow engine pigeon fringe claim program",
				},
				{
					Name:     "relayer",
					Coins:    []string{"100000000000token", "1000000000000000000000stake"},
					Mnemonic: relayerMnemonic,
				},
			},
			Faucet: base.Faucet{
				Name:  &bobName,
				Coins: []string{"500token", "100000000stake"},
				Host:  ":4500",
			},
			Genesis: yamlmap.Map{"chain_id": "earth-1"},
		},
		Validators: []v1.Validator{
			{
				Name:   "alice",
				Bonded: "100000000stake",
				Client: yamlmap.Map{"keyring-backend": keyring.BackendTest},
				App: yamlmap.Map{
					"api":      yamlmap.Map{"address": ":1317"},
					"grpc":     yamlmap.Map{"address": ":9090"},
					"grpc-web": yamlmap.Map{"address": ":9091"},
				},
				Config: yamlmap.Map{
					"p2p": yamlmap.Map{"laddr": ":26656"},
					"rpc": yamlmap.Map{"laddr": ":26656", "pprof_laddr": ":6060"},
				},
				Home: "$HOME/.earth",
			},
		},
	}

	nameOnRecvIbcPostPacket = "OnRecvIbcPostPacket"
	funcOnRecvIbcPostPacket = `
	packetAck.PostId, err = k.PostSeq.Next(ctx)
	if err != nil {
		return packetAck, err
	}
	return packetAck, k.Post.Set(ctx, packetAck.PostId, types.Post{Title: data.Title, Content: data.Content})`

	nameOnAcknowledgementIbcPostPacket = "OnAcknowledgementIbcPostPacket"
	funcOnAcknowledgementIbcPostPacket = `
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		// We will not treat acknowledgment error in this tutorial
		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.IbcPostPacketAck
		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		seq, err := k.SentPostSeq.Next(ctx)
		if err != nil {
			return err
		}

		return k.SentPost.Set(ctx, seq,
			types.SentPost{
				PostId: packetAck.PostId,
				Title:  data.Title,
				Chain:  packet.DestinationPort + "-" + packet.DestinationChannel,
			},
		)
	default:
		return errors.New("the counter-party module does not implement the correct acknowledgment format")
	}`

	nameOnTimeoutIbcPostPacket = "OnTimeoutIbcPostPacket"
	funcOnTimeoutIbcPostPacket = `
	seq, err := k.TimeoutPostSeq.Next(ctx)
	if err != nil {
		return err
	}

	return k.TimeoutPost.Set(ctx, seq,
		types.TimeoutPost{
			Title: data.Title,
			Chain: packet.DestinationPort + "-" + packet.DestinationChannel,
		},
	)`
)

type (
	QueryChannels struct {
		Channels []struct {
			ChannelId      string   `json:"channel_id"`
			ConnectionHops []string `json:"connection_hops"`
			Counterparty   struct {
				ChannelId string `json:"channel_id"`
				PortId    string `json:"port_id"`
			} `json:"counterparty"`
			Ordering string `json:"ordering"`
			PortId   string `json:"port_id"`
			State    string `json:"state"`
			Version  string `json:"version"`
		} `json:"channels"`
	}

	QueryBalances struct {
		Balances sdk.Coins `json:"balances"`
	}
)

func runChain(
	t *testing.T,
	ctx context.Context,
	env envtest.Env,
	app envtest.App,
	cfg v1.Config,
	tmpDir string,
	ports []uint,
) (api, rpc, grpc, faucet string) {
	t.Helper()
	if len(ports) < 7 {
		t.Fatalf("invalid number of ports %d", len(ports))
	}

	var (
		chainID   = cfg.Genesis["chain_id"].(string)
		chainPath = filepath.Join(tmpDir, chainID)
		homePath  = filepath.Join(chainPath, "home")
		cfgPath   = filepath.Join(chainPath, chain.ConfigFilenames[0])
	)
	require.NoError(t, os.MkdirAll(chainPath, os.ModePerm))

	genAddr := func(port uint) string {
		return fmt.Sprintf(":%d", port)
	}

	cfg.Validators[0].Home = homePath

	cfg.Faucet.Host = genAddr(ports[0])
	cfg.Validators[0].App["api"] = yamlmap.Map{"address": genAddr(ports[1])}
	cfg.Validators[0].App["grpc"] = yamlmap.Map{"address": genAddr(ports[2])}
	cfg.Validators[0].App["grpc-web"] = yamlmap.Map{"address": genAddr(ports[3])}
	cfg.Validators[0].Config["p2p"] = yamlmap.Map{"laddr": genAddr(ports[4])}
	cfg.Validators[0].Config["rpc"] = yamlmap.Map{
		"laddr":       genAddr(ports[5]),
		"pprof_laddr": genAddr(ports[6]),
	}

	file, err := os.Create(cfgPath)
	require.NoError(t, err)
	require.NoError(t, yaml.NewEncoder(file).Encode(cfg))
	require.NoError(t, file.Close())

	app.SetConfigPath(cfgPath)
	app.SetHomePath(homePath)
	go func() {
		env.Must(app.Serve("should serve chain", envtest.ExecCtx(ctx)))
	}()

	genHTTPAddr := func(port uint) string {
		return fmt.Sprintf("http://127.0.0.1:%d", port)
	}
	return genHTTPAddr(ports[1]), genHTTPAddr(ports[5]), genHTTPAddr(ports[2]), genHTTPAddr(ports[0])
}

func TestCustomIBCTx(t *testing.T) {
	t.Skip("skip till new CLI version > v28.4.0 released")

	var (
		name        = "blog"
		env         = envtest.New(t)
		app         = env.Scaffold(fmt.Sprintf("github.com/ignite/%s", name), "--no-module")
		tmpDir      = t.TempDir()
		ctx, cancel = context.WithCancel(env.Ctx())
	)
	t.Cleanup(func() {
		cancel()
		time.Sleep(5 * time.Second)
		require.NoError(t, os.RemoveAll(tmpDir))
	})

	dir, err := os.Getwd()
	require.NoError(t, err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "hermes")

	env.Must(env.Exec("install hermes app locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", pluginPath),
			step.Workdir(app.SourcePath()),
		)),
	))

	// One local plugin expected
	assertLocalPlugins(t, app, []pluginsconfig.Plugin{{Path: pluginPath}})
	assertGlobalPlugins(t, nil)

	// prepare the chain
	env.Must(env.Exec("create an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"module",
				"blog",
				"--ibc",
				"--require-registration",
				"--yes",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a post type list in an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"list",
				"post",
				"title",
				"content",
				"--no-message",
				"--module",
				"blog",
				"--yes",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a sentPost type list in an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"list",
				"sentPost",
				"postID:uint",
				"title",
				"chain",
				"--no-message",
				"--module",
				"blog",
				"--yes",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a timeoutPost type list in an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"list",
				"timeoutPost",
				"title",
				"chain",
				"--no-message",
				"--module",
				"blog",
				"--yes",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	env.Must(env.Exec("create a ibcPost package in an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"packet",
				"ibcPost",
				"title",
				"content",
				"--ack",
				"postID:uint",
				"--module",
				"blog",
				"--yes",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	// wait for CLI version > v28.4.0
	// ibcPostPath := filepath.Join(app.SourcePath(), "x/blog/keeper/ibc_post.go")
	// content, err := os.ReadFile(ibcPostPath)
	// require.NoError(t, err)
	// fileContent, err := xast.ModifyFunction(string(content), nameOnRecvIbcPostPacket, xast.ReplaceFuncBody(funcOnRecvIbcPostPacket))
	// require.NoError(t, err)
	// fileContent, err = xast.ModifyFunction(fileContent, nameOnAcknowledgementIbcPostPacket, xast.ReplaceFuncBody(funcOnAcknowledgementIbcPostPacket))
	// require.NoError(t, err)
	// fileContent, err = xast.ModifyFunction(fileContent, nameOnTimeoutIbcPostPacket, xast.ReplaceFuncBody(funcOnTimeoutIbcPostPacket))
	// require.NoError(t, err)
	// require.NoError(t, os.WriteFile(ibcPostPath, []byte(fileContent), 0o644))

	// serve both chains.
	ports, err := availableport.Find(
		14,
		availableport.WithMinPort(4000),
		availableport.WithMaxPort(5000),
	)
	require.NoError(t, err)
	earthAPI, earthRPC, earthGRPC, earthFaucet := runChain(t, ctx, env, app, earthConfig, tmpDir, ports[:7])
	earthChainID := earthConfig.Genesis["chain_id"].(string)
	earthHome := earthConfig.Validators[0].Home
	marsAPI, marsRPC, marsGRPC, marsFaucet := runChain(t, ctx, env, app, marsConfig, tmpDir, ports[7:])
	marsChainID := marsConfig.Genesis["chain_id"].(string)
	marsHome := marsConfig.Validators[0].Home

	// check the chains is up
	stepsCheckChains := step.NewSteps(
		step.New(
			step.Exec(
				app.Binary(),
				"config",
				"output", "json",
			),
			step.PreExec(func() error {
				if err := env.IsAppServed(ctx, earthAPI); err != nil {
					return err
				}
				return env.IsAppServed(ctx, marsAPI)
			}),
			step.Workdir(app.SourcePath()),
		),
	)
	env.Exec("waiting the chain is up", stepsCheckChains, envtest.ExecRetry())

	env.Must(env.Exec("configure the hermes relayer app",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"relayer",
				"hermes",
				"configure",
				earthChainID,
				earthRPC,
				earthGRPC,
				marsChainID,
				marsRPC,
				marsGRPC,
				"--chain-a-faucet", earthFaucet,
				"--chain-b-faucet", marsFaucet,
				"--chain-a-port-id", name,
				"--chain-b-port-id", name,
				"--channel-version", name+"-1",
				"--generate-wallets",
				"--overwrite-config",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	go func() {
		env.Must(env.Exec("run the hermes relayer",
			step.NewSteps(step.New(
				step.Exec(envtest.IgniteApp, "relayer", "hermes", "start", earthChainID, marsChainID),
				step.Workdir(app.SourcePath()),
			)),
			envtest.ExecCtx(ctx),
		))
	}()
	time.Sleep(3 * time.Second)

	var (
		queryOutput   = &bytes.Buffer{}
		queryResponse QueryChannels
	)
	env.Must(env.Exec("verify if the channel was created", step.NewSteps(
		step.New(
			step.Stdout(queryOutput),
			step.Exec(
				app.Binary(),
				"q",
				"ibc",
				"channel",
				"channels",
				"--node", earthRPC,
				"--log_format", "json",
				"--output", "json",
			),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				if err := json.Unmarshal(queryOutput.Bytes(), &queryResponse); err != nil {
					return errors.Errorf("unmarshling tx response: %s", err)
				}
				if len(queryResponse.Channels) == 0 ||
					len(queryResponse.Channels[0].ConnectionHops) == 0 {
					return errors.Errorf("channel not found")
				}
				if queryResponse.Channels[0].State != "STATE_OPEN" {
					return errors.Errorf("channel is not open")
				}
				return nil
			}),
		),
	)))

	var (
		sender     = "alice"
		txOutput   = &bytes.Buffer{}
		txResponse struct {
			Code   int
			RawLog string `json:"raw_log"`
			TxHash string `json:"txhash"`
		}
	)

	stepsTx := step.NewSteps(
		step.New(
			step.Stdout(txOutput),
			step.Exec(
				app.Binary(),
				"tx",
				"blog",
				"send-ibc-post",
				"blog",
				"channel-0",
				"Hello",
				"Hello Mars, I'm Alice from Earth",
				"--from", sender,
				"--node", earthRPC,
				"--home", earthHome,
				"--chain-id", earthChainID,
				"--output", "json",
				"--log_format", "json",
				"--keyring-backend", "test",
				"--yes",
			),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				if err := json.Unmarshal(txOutput.Bytes(), &txResponse); err != nil {
					return errors.Errorf("unmarshling tx response: %w", err)
				}
				return cmdrunner.New().Run(ctx, step.New(
					step.Exec(
						app.Binary(),
						"q",
						"tx",
						txResponse.TxHash,
						"--node", earthRPC,
						"--home", earthHome,
						"--chain-id", earthChainID,
						"--output", "json",
						"--log_format", "json",
					),
					step.Stdout(txOutput),
					step.PreExec(func() error {
						txOutput.Reset()
						return nil
					}),
					step.PostExec(func(execErr error) error {
						if execErr != nil {
							return execErr
						}
						if err := json.Unmarshal(txOutput.Bytes(), &txResponse); err != nil {
							return err
						}
						return nil
					}),
				))
			}),
		),
	)
	if !env.Exec("send an IBC transfer", stepsTx, envtest.ExecRetry()) {
		t.FailNow()
	}
	require.Equalf(t, 0, txResponse.Code,
		"tx failed code=%d log=%s", txResponse.Code, txResponse.RawLog)

	var (
		balanceOutput   = &bytes.Buffer{}
		balanceResponse QueryBalances
	)
	env.Must(env.Exec("check ibc balance", step.NewSteps(
		step.New(
			step.Stdout(balanceOutput),
			step.Exec(
				app.Binary(),
				"q",
				"blog",
				"list-post",
				"--node", marsRPC,
				"--home", marsHome,
				"--log_format", "json",
				"--output", "json",
			),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				if err := json.Unmarshal(balanceOutput.Bytes(), &balanceResponse); err != nil {
					return errors.Errorf("unmarshalling tx response: %w", err)
				}
				if balanceResponse.Balances.Empty() {
					return errors.Errorf("empty balances")
				}
				if !strings.HasPrefix(balanceResponse.Balances[0].Denom, "ibc/") {
					return errors.Errorf("invalid ibc balance: %v", balanceResponse.Balances[0])
				}
				return nil
			}),
		),
	)))
}

func TestTransferIBCTx(t *testing.T) {
	var (
		name        = "blog"
		env         = envtest.New(t)
		app         = env.Scaffold(fmt.Sprintf("github.com/apps/%s", name), "--no-module")
		tmpDir      = t.TempDir()
		ctx, cancel = context.WithCancel(env.Ctx())
	)
	t.Cleanup(func() {
		cancel()
		time.Sleep(5 * time.Second)
		require.NoError(t, os.RemoveAll(tmpDir))
	})

	dir, err := os.Getwd()
	require.NoError(t, err)
	pluginPath := filepath.Join(filepath.Dir(filepath.Dir(dir)), "hermes")

	env.Must(env.Exec("install hermes app locally",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp, "app", "install", pluginPath),
			step.Workdir(app.SourcePath()),
		)),
	))

	// One local plugin expected
	assertLocalPlugins(t, app, []pluginsconfig.Plugin{{Path: pluginPath}})
	assertGlobalPlugins(t, nil)

	// prepare the chain
	env.Must(env.Exec("create an IBC module",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"s",
				"module",
				"blog",
				"--ibc",
				"--require-registration",
				"--yes",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	// serve both chains.
	ports, err := availableport.Find(
		14,
		availableport.WithMinPort(4000),
		availableport.WithMaxPort(5000),
	)
	require.NoError(t, err)
	earthAPI, earthRPC, earthGRPC, earthFaucet := runChain(t, ctx, env, app, earthConfig, tmpDir, ports[:7])
	earthChainID := earthConfig.Genesis["chain_id"].(string)
	earthHome := earthConfig.Validators[0].Home
	marsAPI, marsRPC, marsGRPC, marsFaucet := runChain(t, ctx, env, app, marsConfig, tmpDir, ports[7:])
	marsChainID := marsConfig.Genesis["chain_id"].(string)
	marsHome := marsConfig.Validators[0].Home

	// check the chains is up
	stepsCheckChains := step.NewSteps(
		step.New(
			step.Exec(
				app.Binary(),
				"config",
				"output", "json",
			),
			step.PreExec(func() error {
				if err := env.IsAppServed(ctx, earthAPI); err != nil {
					return err
				}
				return env.IsAppServed(ctx, marsAPI)
			}),
			step.Workdir(app.SourcePath()),
		),
	)
	env.Exec("waiting the chain is up", stepsCheckChains, envtest.ExecRetry())

	env.Must(env.Exec("configure the hermes relayer app",
		step.NewSteps(step.New(
			step.Exec(envtest.IgniteApp,
				"relayer",
				"hermes",
				"configure",
				earthChainID,
				earthRPC,
				earthGRPC,
				marsChainID,
				marsRPC,
				marsGRPC,
				"--chain-a-faucet", earthFaucet,
				"--chain-b-faucet", marsFaucet,
				"--generate-wallets",
				"--overwrite-config",
			),
			step.Workdir(app.SourcePath()),
		)),
	))

	go func() {
		env.Must(env.Exec("run the hermes relayer",
			step.NewSteps(step.New(
				step.Exec(envtest.IgniteApp, "relayer", "hermes", "start", earthChainID, marsChainID),
				step.Workdir(app.SourcePath()),
			)),
			envtest.ExecCtx(ctx),
		))
	}()
	time.Sleep(3 * time.Second)

	var (
		queryOutput   = &bytes.Buffer{}
		queryResponse QueryChannels
	)
	env.Must(env.Exec("verify if the channel was created", step.NewSteps(
		step.New(
			step.Stdout(queryOutput),
			step.Exec(
				app.Binary(),
				"q",
				"ibc",
				"channel",
				"channels",
				"--node", earthRPC,
				"--log_format", "json",
				"--output", "json",
			),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				if err := json.Unmarshal(queryOutput.Bytes(), &queryResponse); err != nil {
					return errors.Errorf("unmarshling tx response: %s", err)
				}
				if len(queryResponse.Channels) == 0 ||
					len(queryResponse.Channels[0].ConnectionHops) == 0 {
					return errors.Errorf("channel not found")
				}
				if queryResponse.Channels[0].State != "STATE_OPEN" {
					return errors.Errorf("channel is not open")
				}
				return nil
			}),
		),
	)))

	var (
		sender       = "alice"
		receiverAddr = "cosmos1nrksk5swk6lnmlq670a8kwxmsjnu0ezqts39sa"
		txOutput     = &bytes.Buffer{}
		txResponse   struct {
			Code   int
			RawLog string `json:"raw_log"`
			TxHash string `json:"txhash"`
		}
	)

	stepsTx := step.NewSteps(
		step.New(
			step.Stdout(txOutput),
			step.Exec(
				app.Binary(),
				"tx",
				"ibc-transfer",
				"transfer",
				"transfer",
				"channel-0",
				receiverAddr,
				"100000stake",
				"--from", sender,
				"--node", earthRPC,
				"--home", earthHome,
				"--chain-id", earthChainID,
				"--output", "json",
				"--log_format", "json",
				"--keyring-backend", "test",
				"--yes",
			),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				if err := json.Unmarshal(txOutput.Bytes(), &txResponse); err != nil {
					return errors.Errorf("unmarshling tx response: %w", err)
				}
				return cmdrunner.New().Run(ctx, step.New(
					step.Exec(
						app.Binary(),
						"q",
						"tx",
						txResponse.TxHash,
						"--node", earthRPC,
						"--home", earthHome,
						"--chain-id", earthChainID,
						"--output", "json",
						"--log_format", "json",
					),
					step.Stdout(txOutput),
					step.PreExec(func() error {
						txOutput.Reset()
						return nil
					}),
					step.PostExec(func(execErr error) error {
						if execErr != nil {
							return execErr
						}
						if err := json.Unmarshal(txOutput.Bytes(), &txResponse); err != nil {
							return err
						}
						return nil
					}),
				))
			}),
		),
	)
	if !env.Exec("send an IBC transfer", stepsTx, envtest.ExecRetry()) {
		t.FailNow()
	}
	require.Equalf(t, 0, txResponse.Code,
		"tx failed code=%d log=%s", txResponse.Code, txResponse.RawLog)

	var (
		balanceOutput   = &bytes.Buffer{}
		balanceResponse QueryBalances
	)
	env.Must(env.Exec("check ibc balance", step.NewSteps(
		step.New(
			step.Stdout(balanceOutput),
			step.Exec(
				app.Binary(),
				"q",
				"bank",
				"balances",
				receiverAddr,
				"--node", marsRPC,
				"--home", marsHome,
				"--log_format", "json",
				"--output", "json",
			),
			step.PostExec(func(execErr error) error {
				if execErr != nil {
					return execErr
				}
				if err := json.Unmarshal(balanceOutput.Bytes(), &balanceResponse); err != nil {
					return errors.Errorf("unmarshalling tx response: %w", err)
				}
				if balanceResponse.Balances.Empty() {
					return errors.Errorf("empty balances")
				}
				if !strings.HasPrefix(balanceResponse.Balances[0].Denom, "ibc/") {
					return errors.Errorf("invalid ibc balance: %v", balanceResponse.Balances[0])
				}
				return nil
			}),
		),
	)))
}

func assertLocalPlugins(t *testing.T, app envtest.App, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfg, err := pluginsconfig.ParseDir(app.SourcePath())
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected local apps")
}

func assertGlobalPlugins(t *testing.T, expectedPlugins []pluginsconfig.Plugin) {
	t.Helper()
	cfgPath, err := plugin.PluginsPath()
	require.NoError(t, err)
	cfg, err := pluginsconfig.ParseDir(cfgPath)
	require.NoError(t, err)
	require.ElementsMatch(t, expectedPlugins, cfg.Apps, "unexpected global apps")
}
