package work

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// ErrSignalReceived is the error returned when a signal is received by the SignalHandler.
var ErrSignalReceived = fmt.Errorf("signal received")

// SignalHandler is a Task that waits for a signal (os.Interrupt or
// syscall.SIGTERM) and returns an error when one is received.
func SignalHandler(ctx context.Context) error {
	return signalHandlerImplementation(ctx, &defaultSignals{})
}

// signals is an interface to signal.Notify and signal.Stop, allowing for
// easier testing and mocking of the signal handling behavior.
type signals interface {
	Notify(c chan<- os.Signal, sig ...os.Signal)
	Stop(c chan<- os.Signal)
}

// defaultSignals is the default implementation of the signals interface, using
// the standard library's signal package.
type defaultSignals struct{}

func (*defaultSignals) Notify(c chan<- os.Signal, sig ...os.Signal) {
	signal.Notify(c, sig...)
}

func (*defaultSignals) Stop(c chan<- os.Signal) {
	signal.Stop(c)
}

// signalHandlerImplementation is the implementation of SignalHandler, abstracting
// the signal.Notify and signal.Stop calls for testing purposes.
func signalHandlerImplementation(ctx context.Context, sigs signals) error {
	sigChan := make(chan os.Signal, 1)
	sigs.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer sigs.Stop(sigChan)
	defer close(sigChan)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case sig := <-sigChan:
		return fmt.Errorf("%w: %s", ErrSignalReceived, sig.String())
	}
}
