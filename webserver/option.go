package webserver

import "time"

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
