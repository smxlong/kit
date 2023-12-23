package logger

// Option is an option for the logger.
type Option func(*options) error

// options is a struct that holds the options for the logger.
type options struct {
	level  string
	format string
}

// optionsFrom returns an options from the given options.
func optionsFrom(opts ...Option) (*options, error) {
	options := &options{
		level:  "info",
		format: "json",
	}
	for _, opt := range opts {
		if err := opt(options); err != nil {
			return nil, err
		}
	}
	return options, nil
}

// WithLevel sets the level of the logger.
func WithLevel(level string) Option {
	return func(o *options) error {
		o.level = level
		return nil
	}
}

// WithFormat sets the format of the logger.
func WithFormat(format string) Option {
	return func(o *options) error {
		o.format = format
		return nil
	}
}
