package signalcontext

import (
	"context"
	"os"
	"os/signal"
)

// WithSignals returns a context that is canceled when any of the given
// signals are received.
func WithSignals(ctx context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	sigch := make(chan os.Signal, 1)
	go func() {
		select {
		case <-ctx.Done():
		case <-sigch:
			cancel()
		}
	}()
	signal.Notify(sigch, signals...)
	return ctx, cancel
}
