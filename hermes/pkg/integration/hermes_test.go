package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"relayer/pkg/hermes"
)

func TestHermes(t *testing.T) {
	ctx := context.Background()
	h, err := hermes.New()
	require.NoError(t, err)
	defer h.Cleanup()

	// Create the default config and add chains
	c := hermes.DefaultConfig()
	err = c.AddChain("mars-1", "http://localhost:26649", "http://localhost:9082")
	require.NoError(t, err)

	err = c.AddChain("venus-1", "http://localhost:26659", "http://localhost:9092")
	require.NoError(t, err)

	err = c.Save()
	require.NoError(t, err)

	cfgPath, err := c.ConfigPath()
	require.NoError(t, err)

	// Add hermes keys
	var (
		buf    = bytes.Buffer{}
		result = hermes.Result{}
	)
	err = h.AddMnemonic(
		ctx,
		"mars-1",
		"letter column benefit acoustic evidence false trim cave jump pluck awesome lion",
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&buf),
	)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(buf.Bytes(), &result))
	fmt.Println(result)

	buf = bytes.Buffer{}
	result = hermes.Result{}
	err = h.AddMnemonic(
		ctx,
		"venus-1",
		"jeans payment lock client result enemy bullet rug crush deny month salad",
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&buf),
	)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(buf.Bytes(), &result))
	fmt.Println(result)

	// create clients
	buf = bytes.Buffer{}
	result = hermes.Result{}
	err = h.CreateClient(
		ctx,
		"mars-1",
		"venus-1",
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&buf),
	)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(buf.Bytes(), &result))

	var clientResult1 hermes.ClientResult
	require.NoError(t, json.Unmarshal(result.Result, &clientResult1))
	fmt.Println(clientResult1)

	buf = bytes.Buffer{}
	result = hermes.Result{}
	err = h.CreateClient(
		ctx,
		"venus-1",
		"mars-1",
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&buf),
	)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(buf.Bytes(), &result))

	var clientResult2 hermes.ClientResult
	require.NoError(t, json.Unmarshal(result.Result, &clientResult2))
	fmt.Println(clientResult2)

	// create connection
	buf = bytes.Buffer{}
	result = hermes.Result{}
	err = h.CreateConnection(
		ctx,
		"mars-1",
		"07-tendermint-0",
		"07-tendermint-0",
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&buf),
	)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(buf.Bytes(), &result))

	var connection hermes.ConnectionResult
	require.NoError(t, json.Unmarshal(result.Result, &connection))
	fmt.Println(connection)

	// create and query channel
	buf = bytes.Buffer{}
	result = hermes.Result{}
	err = h.CreateChannel(
		ctx,
		"mars-1",
		"connection-0",
		"transfer",
		"transfer",
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&buf),
	)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(buf.Bytes(), &result))

	var channel hermes.ConnectionResult
	require.NoError(t, json.Unmarshal(result.Result, &channel))
	fmt.Println(channel)

	buf = bytes.Buffer{}
	result = hermes.Result{}
	err = h.QueryChannels(
		ctx,
		true,
		"mars-1",
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(&buf),
	)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(buf.Bytes(), &result))

	var channels []hermes.ChannelResult
	require.NoError(t, json.Unmarshal(result.Result, &channels))
	fmt.Println(channels)

	// start hermes
	err = h.Start(
		ctx,
		hermes.WithConfigFile(cfgPath),
		hermes.WithStdOut(os.Stdout),
	)
}
