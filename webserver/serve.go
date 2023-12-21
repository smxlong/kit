package webserver

import (
	"context"
	"net"
	"net/http"
)

// Serve runs the http.Server until the context is canceled or an error occurs.
func Serve(ctx context.Context, server *http.Server, l net.Listener, opts ...Option) error {
	return run(ctx, server, func() error { return server.Serve(l) }, opts...)
}

// ListenAndServe runs the http.Server until the context is canceled or an error
// occurs.
func ListenAndServe(ctx context.Context, server *http.Server, opts ...Option) error {
	return run(ctx, server, server.ListenAndServe, opts...)
}
