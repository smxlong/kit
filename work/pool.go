package work

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Task is a function that can be executed in a work pool.
type Task func(ctx context.Context) error

// Pool is a work pool.
type Pool struct {
	eg     *errgroup.Group
	ctx    context.Context
	cancel context.CancelFunc
}

// NewPool creates a new work pool.
func NewPool() *Pool {
	return NewPoolWithContext(context.Background())
}

// NewPoolWithContext creates a new work pool with the given context.
func NewPoolWithContext(ctx context.Context) *Pool {
	ctx, cancel := context.WithCancel(ctx)
	eg, ctx := errgroup.WithContext(ctx)
	return &Pool{
		eg:     eg,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Run adds a task to the work pool.
func (p *Pool) Run(task Task) {
	p.eg.Go(func() error {
		return task(p.ctx)
	})
}

// Wait waits for all tasks to complete and returns any error encountered.
func (p *Pool) Wait() error {
	return p.eg.Wait()
}

// Cancel cancels the work pool, stopping any running tasks. Cancel returns
// immediately and does not wait for tasks to finish.
func (p *Pool) Cancel() {
	p.cancel()
}
