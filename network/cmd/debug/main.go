package main

import (
	"fmt"
	"os"

	"github.com/ignite/apps/network/cmd"
)

func main() {
	if err := cmd.NewNetwork().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
