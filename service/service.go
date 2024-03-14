package service

import (
	"context"
	"os"
	"strconv"
	"syscall"

	"github.com/smxlong/kit/logger"
	"github.com/smxlong/kit/signalcontext"
)

// Run a service. A context is created which will be canceled when the service
// receives a SIGINT or SIGTERM signal. A logger is created. These are passed
// to the service function. The logger is created with a level of DEBUG if the
// DEBUG environment variable is set to true.
func Run(service func(ctx context.Context, l logger.Logger) error) error {
	ctx, cancel := signalcontext.WithSignals(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	debug := false
	if dbg, ok := os.LookupEnv("DEBUG"); ok {
		if v, err := strconv.ParseBool(dbg); err == nil {
			debug = v
		}
	}
	opts := []logger.Option{}
	if debug {
		opts = append(opts, logger.WithLevel("DEBUG"))
	}
	l, err := logger.New(opts...)
	if err != nil {
		return err
	}

	l.Debugw("service started")
	defer l.Debugw("service ended")
	return service(ctx, l)
}
