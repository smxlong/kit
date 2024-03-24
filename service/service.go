package service

import (
	"context"
	"os"
	"strconv"
	"syscall"

	"github.com/smxlong/kit/logger"
	"github.com/smxlong/kit/signalcontext"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	l, _ := logger.New(opts...)

	l.Debugw("service started")
	defer l.Debugw("service ended")
	return service(ctx, l)
}

// Runnable must be implemented by all services that are run with Main.
type Runnable interface {
	Run(ctx context.Context, l logger.Logger) error
}

// BindFlags can optionally be implemented by a service to bind command line
// flags to the service. BindFlags will be called before BindEnvironment and
// Run.
type BindFlags interface {
	BindFlags(flags *pflag.FlagSet)
}

// BindEnvironment can optionally be implemented by a service to bind
// environment variables to the service. BindEnvironment will be called after
// BindFlags and before Run.
type BindEnvironment interface {
	BindEnvironment() error
}

// Main runs your service, calling its BindFlags and BindEnvironment methods if
// it implements them, and then calling its Run method.
func Main(use, short string, s Runnable) {
	main(use, short, s, os.Exit)
}

func main(use, short string, s Runnable, exitFunc func(int)) {
	cmd := &cobra.Command{
		Use:          use,
		Short:        short,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if be, ok := s.(BindEnvironment); ok {
				if err := be.BindEnvironment(); err != nil {
					return err
				}
			}
			return Run(s.Run)
		},
	}
	if bf, ok := s.(BindFlags); ok {
		bf.BindFlags(cmd.Flags())
	}
	if err := cmd.Execute(); err != nil {
		exitFunc(1)
	}
}
