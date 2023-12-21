package signal

import (
	"context"
	"os"
	"os/signal"
)

// Context returns a context that is canceled when a signal is received.
func Context(signals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
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
