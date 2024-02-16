package main

import (
	"fmt"
	"os"

	"github.com/ignite/apps/wasm/cmd"
)

func main() {
	if err := cmd.NewWasm().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
