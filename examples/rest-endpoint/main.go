package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/smxlong/kit/rest"
	"github.com/smxlong/kit/signalcontext"
	"github.com/smxlong/kit/webserver"
)

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
	http.Handle("/api/example", &rest.Endpoint{
		Method: map[string]rest.Handler{
			"POST": {
				NewRequest: func() rest.Request {
					return &ExampleRequest{}
				},
				Handle: func(ctx context.Context, req rest.Request) rest.Response {
					r := req.(*ExampleRequest)
					return &ExampleResponse{
						Sum: r.A + r.B,
					}
				},
			},
		},
	})

	ctx, cancel := signalcontext.WithSignals(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	webserver.ListenAndServe(ctx, server)
	fmt.Println("Server stopped")
}
