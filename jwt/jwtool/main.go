package main

// jwtool is a tool for generating, validating, and manipulating JWTs.

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Commands:
//
// - encode
// - decode
// - validate

func main() {
	j := &jwtool{}
	if err := j.Command().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// jwtool is the jwtool command.
type jwtool struct {
	encode encode
	decode decode
}

func (j *jwtool) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jwtool",
		Short: "A tool for encoding, decoding, validating, and manipulating JWTs.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		j.encode.Command(),
		j.decode.Command(),
	)
	return cmd
}
