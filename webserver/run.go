package webserver

import (
	"context"
	"net/http"
)

// run runs the http.Server until the context is canceled or an error occurs.
func run(ctx context.Context, server *http.Server, sfunc func() error, opts ...Option) error {
	o, err := optionsFrom(opts...)
	if err != nil {
		return err
	}
	errCh := make(chan error)
	go func() {
		errCh <- sfunc()
	}()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), o.shutdownTimeout)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}
