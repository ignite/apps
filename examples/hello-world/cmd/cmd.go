package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewHelloWorld creates a new hello-world command.
func NewHelloWorld() *cobra.Command {
	c := &cobra.Command{
		Use:           "hello-world",
		Short:         "Say hello to the world of ignite!",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Hello, world!")
			return nil
		},
	}

	return c
}
