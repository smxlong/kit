package work

import (
	"context"
	"time"
)

// Sleep the provided duration, or until the context is canceled.
func Sleep(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Timeout returns a function that calls Sleep with the provided duration.
func Timeout(d time.Duration) Task {
	return func(ctx context.Context) error {
		return Sleep(ctx, d)
	}
}
