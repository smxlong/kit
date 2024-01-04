package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

	gojwt "github.com/golang-jwt/jwt/v5"

	"github.com/smxlong/kit/jwt"
	"github.com/smxlong/kit/rest"
	"github.com/smxlong/kit/signalcontext"
	"github.com/smxlong/kit/webserver"
)

// This example is based on rest-endpoint. Compare and contrast to understand
// how the JWT middleware is used.

type ExampleRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

type ExampleResponse struct {
	Sum int `json:"sum"`
}

func main() {
	server := &http.Server{
		Addr: ":8080",
	}
	jwtMiddleware := jwt.NewMiddleware(
		jwt.WithKey([]byte("secret")),
		jwt.WithAudience("https://example.com"),
		jwt.WithIssuer("https://example.com"),
	)
	http.Handle("/api/example", jwtMiddleware.Wrap(&rest.Endpoint{
		Method: map[string]rest.Handler{
			"POST": {
				NewRequest: func() rest.Request {
					return &ExampleRequest{}
				},
				Handle: func(ctx context.Context, req rest.Request) rest.Response {
					if claims, ok := ctx.Value(jwt.ContextKeyClaims).(gojwt.Claims); ok {
						fmt.Println("Claims:", claims)
					}
					r := req.(*ExampleRequest)
					return &ExampleResponse{
						Sum: r.A + r.B,
					}
				},
			},
		},
	}))

	ctx, cancel := signalcontext.WithSignals(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	webserver.ListenAndServe(ctx, server)
	fmt.Println("Server stopped")
}
