package hermes_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"relayer/pkg/hermes"
	"testing"
)

func TestHermes(t *testing.T) {
	cfgPath, err := hermes.ConfigPath()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Dir(cfgPath)); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	h, err := hermes.New()
	if err != nil {
		t.Fatal(err)
	}
	defer h.Cleanup()

	// Create the default config and add chains
	c := hermes.DefaultConfig()
	err = c.AddChain("mars-1", "http://localhost:26649", "http://localhost:9082")
	if err != nil {
		t.Fatal(err)
	}
	err = c.AddChain("venus-1", "http://localhost:26659", "http://localhost:9092")
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Save(); err != nil {
		t.Fatal(err)
	}

	// Add hermes keys
	result, err := h.AddMnemonic(
		ctx,
		"mars-1",
		"letter column benefit acoustic evidence false trim cave jump pluck awesome lion",
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	result, err = h.AddMnemonic(
		ctx,
		"venus-1",
		"jeans payment lock client result enemy bullet rug crush deny month salad",
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	// create clients
	result, err = h.CreateClient(ctx, "mars-1", "venus-1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	result, err = h.CreateClient(ctx, "venus-1", "mars-1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	// create connection
	result, err = h.CreateConnection(ctx, "mars-1", "07-tendermint-0", "07-tendermint-0")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	// create and query channel
	result, err = h.CreateChannel(ctx, "mars-1", "connection-0", "transfer", "transfer")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	result, err = h.QueryChannels(ctx, true, "mars-1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	// start hermes
	result, err = h.Start(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}
