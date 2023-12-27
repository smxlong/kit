package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a logging interface. It supports some of the methods of the
// zap.SugaredLogger interface.
type Logger interface {
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	With(keysAndValues ...interface{}) Logger
	Sync() error
}

// logger is a wrapper around zap.SugaredLogger.
type logger struct {
	*zap.SugaredLogger
}

// New returns a new Logger.
func New(opts ...Option) (Logger, error) {
	options, err := optionsFrom(opts...)
	if err != nil {
		return nil, err
	}
	level, err := zap.ParseAtomicLevel(options.level)
	if err != nil {
		return nil, err
	}
	config := zap.NewProductionConfig()
	config.Level = level
	config.Encoding = options.format
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	l, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &logger{l.Sugar()}, nil
}

// With returns a new Logger with the given keys and values.
func (l *logger) With(keysAndValues ...interface{}) Logger {
	return &logger{l.SugaredLogger.With(keysAndValues...)}
}
