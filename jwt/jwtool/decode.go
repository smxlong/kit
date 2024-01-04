package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

// decode is the decode command.
type decode struct {
	key        string // key to sign with
	token      string // token to decode
	noValidate bool   // don't validate the token
}

// Command returns the decode command.
func (d *decode) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decode",
		Short: "Decode a JWT.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return d.do(cmd, args)
		},
	}
	cmd.Flags().StringVarP(&d.key, "key", "k", "", "The key to sign with.")
	cmd.Flags().StringVarP(&d.token, "token", "t", "", "The token to decode. If not given, token will be read from stdin.")
	cmd.Flags().BoolVarP(&d.noValidate, "no-validate", "N", false, "Don't validate the token.")
	return cmd
}

func (d *decode) do(cmd *cobra.Command, args []string) error {
	if d.token == "" {
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return errors.New("no token given")
		}
		d.token = scanner.Text()
	}

	claims := gojwt.MapClaims{}
	tok, err := gojwt.ParseWithClaims(d.token, &claims, func(token *gojwt.Token) (interface{}, error) {
		if d.key == "" {
			return nil, errors.New("no key specified")
		}
		return []byte(d.key), nil
	})
	if err != nil {
		if tok != nil && d.noValidate {
			fmt.Printf("WARNING: %v\n", err)
		} else {
			return err
		}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(tok); err != nil {
		return err
	}
	return nil
}
