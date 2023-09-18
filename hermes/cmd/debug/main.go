package main

import (
	"fmt"
	"os"

	"github.com/ignite/plugins/relayer/cmd"
)

func main() {
	if err := cmd.NewRelayer().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
