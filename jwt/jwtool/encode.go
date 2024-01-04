package main

import (
	"fmt"
	"strings"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

// encode is the encode command.
type encode struct {
	key       string   // key to sign with
	claims    []string // key=value
	expiresIn time.Duration
	expiresAt time.Time
	notBefore time.Time
	subject   string
	audience  []string
	issuer    string
}

// Command returns the encode command.
func (e *encode) Command() *cobra.Command {
	var expiresAt, notBefore string
	cmd := &cobra.Command{
		Use:   "encode",
		Short: "Encode a JWT.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if expiresAt != "" {
				e.expiresAt, err = time.Parse(time.RFC3339, expiresAt)
				if err != nil {
					return err
				}
			}
			if notBefore != "" {
				e.notBefore, err = time.Parse(time.RFC3339, notBefore)
				if err != nil {
					return err
				}
			}
			return e.do(cmd, args)
		},
	}
	cmd.Flags().StringVarP(&e.key, "key", "k", "", "The key to sign with.")
	cmd.Flags().StringArrayVarP(&e.claims, "claim", "c", nil, "A \"key=value\" claim to add to the JWT. Can be specified multiple times.")
	cmd.Flags().DurationVarP(&e.expiresIn, "expires-in", "e", 24*time.Hour, "The duration until the JWT expires.")
	cmd.Flags().StringVarP(&expiresAt, "expires-at", "E", "", "The time at which the JWT expires.")
	cmd.Flags().StringVarP(&notBefore, "not-before", "n", "", "The time before which the JWT is not valid.")
	cmd.Flags().StringVarP(&e.subject, "subject", "s", "jdoe@example.com", "The subject of the JWT.")
	cmd.Flags().StringArrayVarP(&e.audience, "audience", "a", []string{"https://example.com"}, "The audience of the JWT. Can be specified multiple times.")
	cmd.Flags().StringVarP(&e.issuer, "issuer", "i", "https://example.com", "The issuer of the JWT.")
	return cmd
}

func (e *encode) do(cmd *cobra.Command, args []string) error {
	claims := map[string]interface{}{}
	for _, kv := range e.claims {
		kvPair := strings.SplitN(kv, "=", 2)
		if len(kvPair) != 2 {
			return cmd.Help()
		}
		claims[kvPair[0]] = kvPair[1]
	}
	if e.expiresIn != 0 {
		claims["exp"] = time.Now().Add(e.expiresIn).Unix()
	}
	if !e.expiresAt.IsZero() {
		claims["exp"] = e.expiresAt.Unix()
	}
	if !e.notBefore.IsZero() {
		claims["nbf"] = e.notBefore.Unix()
	}
	if e.subject != "" {
		claims["sub"] = e.subject
	}
	if len(e.audience) > 0 {
		claims["aud"] = e.audience
	}
	if e.issuer != "" {
		claims["iss"] = e.issuer
	}
	tok, err := gojwt.NewWithClaims(gojwt.SigningMethodHS256, gojwt.MapClaims(claims)).SignedString([]byte(e.key))
	if err != nil {
		return err
	}
	fmt.Println(tok)
	return nil
}
