package work

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_that_NewPool_returns_a_Pool(t *testing.T) {
	p := NewPool()
	assert.NotNil(t, p)
	assert.NotNil(t, p.eg)
	assert.NotNil(t, p.ctx)
	assert.NotNil(t, p.cancel)
}

type MyContextKeyType string

const MyContextKey MyContextKeyType = "my_context_key"

func Test_that_NewPoolWithContext_returns_a_Pool(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, MyContextKey, "my_value")
	p := NewPoolWithContext(ctx)
	assert.NotNil(t, p)
	assert.NotNil(t, p.eg)
	assert.NotNil(t, p.cancel)
	assert.Equal(t, "my_value", p.ctx.Value(MyContextKey))
}

func Test_that_Pool_basically_works(t *testing.T) {
	p := NewPool()
	var called bool
	task := func(ctx context.Context) error {
		called = true
		return nil
	}

	p.Run(task)
	err := p.Wait()

	assert.NoError(t, err)
	assert.True(t, called)
}

func Test_that_Pool_runs_multiple_tasks_to_completion(t *testing.T) {
	p := NewPool()
	var called1, called2 bool
	task1 := func(ctx context.Context) error {
		called1 = true
		return nil
	}
	task2 := func(ctx context.Context) error {
		called2 = true
		return nil
	}

	p.Run(task1)
	p.Run(task2)
	err := p.Wait()

	assert.NoError(t, err)
	assert.True(t, called1)
	assert.True(t, called2)
}

func Test_that_Pool_can_be_cancelled(t *testing.T) {
	p := NewPool()
	var called bool
	task := func(ctx context.Context) error {
		called = true
		<-ctx.Done() // simulate work that can be cancelled
		return nil
	}

	p.Run(task)
	p.Cancel() // cancel the pool

	err := p.Wait()

	assert.NoError(t, err)
	assert.True(t, called)
}

func Test_that_Pool_can_be_cancelled_from_goroutine(t *testing.T) {
	p := NewPool()
	sync := make(chan struct{})
	defer close(sync)
	var called bool
	task := func(ctx context.Context) error {
		called = true
		sync <- struct{}{} // tell cancel goroutine we've started
		<-ctx.Done()       // simulate work that can be cancelled
		return nil
	}

	p.Run(task)

	go func() {
		<-sync     // wait for the task to start
		p.Cancel() // cancel the pool from a goroutine
	}()

	err := p.Wait()

	assert.NoError(t, err)
	assert.True(t, called)
}

func Test_that_Pool_wait_returns_error_when_task_returns_error(t *testing.T) {
	p := NewPool()
	task := func(ctx context.Context) error {
		return assert.AnError // simulate an error
	}

	p.Run(task)
	err := p.Wait()

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func Test_that_Pool_cancels_all_tasks_on_error(t *testing.T) {
	p := NewPool()
	var called1, called2 bool
	task1 := func(ctx context.Context) error {
		called1 = true
		return assert.AnError // simulate an error
	}
	task2 := func(ctx context.Context) error {
		called2 = true
		<-ctx.Done() // simulate work that can be cancelled
		return nil
	}

	p.Run(task1)
	p.Run(task2)
	err := p.Wait()

	assert.Error(t, err)
	assert.True(t, called1)
	assert.True(t, called2) // task2 should still be called even if task1 returns an error
}
