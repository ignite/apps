package hermes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestHermes(t *testing.T) {
	ctx := context.Background()
	h, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer h.Cleanup()

	// Create the default config and add chains
	c := DefaultConfig()
	err = c.AddChain("mars-1", "http://localhost:26649", "http://localhost:9082")
	if err != nil {
		t.Fatal(err)
	}
	err = c.AddChain("venus-1", "http://localhost:26659", "http://localhost:9092")
	if err != nil {
		t.Fatal(err)
	}

	path, err := c.Save()
	if err != nil {
		t.Fatal(err)
	}

	// Add hermes keys
	var (
		buf    = bytes.Buffer{}
		result = Result{}
	)
	err = h.AddMnemonic(
		ctx,
		"mars-1",
		"letter column benefit acoustic evidence false trim cave jump pluck awesome lion",
		WithConfigFile(path),
		WithStdOut(&buf),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	buf = bytes.Buffer{}
	result = Result{}
	err = h.AddMnemonic(
		ctx,
		"venus-1",
		"jeans payment lock client result enemy bullet rug crush deny month salad",
		WithConfigFile(path),
		WithStdOut(&buf),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	// create clients
	buf = bytes.Buffer{}
	result = Result{}
	err = h.CreateClient(
		ctx,
		"mars-1",
		"venus-1",
		WithConfigFile(path),
		WithStdOut(&buf),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	buf = bytes.Buffer{}
	result = Result{}
	err = h.CreateClient(
		ctx,
		"venus-1",
		"mars-1",
		WithConfigFile(path),
		WithStdOut(&buf),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	// create connection
	buf = bytes.Buffer{}
	result = Result{}
	err = h.CreateConnection(
		ctx,
		"mars-1",
		"07-tendermint-0",
		"07-tendermint-0",
		WithConfigFile(path),
		WithStdOut(&buf),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	// create and query channel
	buf = bytes.Buffer{}
	result = Result{}
	err = h.CreateChannel(
		ctx,
		"mars-1",
		"connection-0",
		"transfer",
		"transfer",
		WithConfigFile(path),
		WithStdOut(&buf),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	buf = bytes.Buffer{}
	result = Result{}
	err = h.QueryChannels(
		ctx,
		true,
		"mars-1",
		WithConfigFile(path),
		WithStdOut(&buf),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	// start hermes
	buf = bytes.Buffer{}
	result = Result{}
	err = h.Start(
		ctx,
		WithConfigFile(path),
		WithStdOut(&buf),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}
