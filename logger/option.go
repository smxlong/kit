package logger

// Option is an option for the logger.
type Option func(*options) error

// options is a struct that holds the options for the logger.
type options struct {
	level      string
	format     string
	stacktrace bool
}

// optionsFrom returns an options from the given options.
func optionsFrom(opts ...Option) (*options, error) {
	options := &options{
		level:      "info",
		format:     "json",
		stacktrace: false,
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

// WithStacktrace sets the stacktrace of the logger.
func WithStacktrace(stacktrace bool) Option {
	return func(o *options) error {
		o.stacktrace = stacktrace
		return nil
	}
}
