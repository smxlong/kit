package work

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testSignals implements signals and passes the given signal to the channel.
type testSignals struct {
	signal os.Signal
}

func (ts *testSignals) Notify(c chan<- os.Signal, sig ...os.Signal) {
	for _, s := range sig {
		if s == ts.signal {
			c <- s
			return
		}
	}
}

func (ts *testSignals) Stop(chan<- os.Signal) {
}

func Test_that_signalHandlerImplementation_works(t *testing.T) {
	ts := &testSignals{signal: os.Interrupt}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := signalHandlerImplementation(ctx, ts)
	assert.Error(t, err)
	assert.Equal(t, "signal received: interrupt", err.Error())
}

func Test_that_signalHandlerImplementation_returns_context_error(t *testing.T) {
	ts := &testSignals{} // no signal will be sent
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context to simulate a done state
	err := signalHandlerImplementation(ctx, ts)
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}
