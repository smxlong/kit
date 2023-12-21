package webserver

import (
	"context"
	"net/http"
	"time"
)

// Option is a function that modifies the behavior of ListenAndServe and Serve.
type Option func(*options) error

// options is a struct that holds the options for ListenAndServe and Serve.
type options struct {
	shutdownTimeout time.Duration
}

// optionsFrom returns a runOptions from the given options.
func optionsFrom(opts ...Option) (*options, error) {
	options := &options{
		shutdownTimeout: 5 * time.Second,
	}
	for _, opt := range opts {
		if err := opt(options); err != nil {
			return nil, err
		}
	}
	return options, nil
}

// WithShutdownTimeout sets the timeout for graceful shutdown.
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(o *options) error {
		o.shutdownTimeout = timeout
		return nil
	}
}

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
