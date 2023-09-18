package main

import (
	"fmt"
	"os"

	"github.com/ignite/apps/hermes/cmd"
)

func main() {
	cobraCmd, err := cmd.NewRelayer()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := cobraCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
