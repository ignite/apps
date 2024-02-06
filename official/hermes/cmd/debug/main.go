package main

import (
	"fmt"
	"os"

	"github.com/ignite/apps/official/hermes/cmd"
)

func main() {
	if err := cmd.NewHermes().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
