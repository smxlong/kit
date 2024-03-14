package boolean

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// ContextPredicate is a boolean function of a single argument and a context.
type ContextPredicate[T any] func(context.Context, T) bool

// ContextNot returns the negation of the given predicate.
func ContextNot[T any](p ContextPredicate[T]) ContextPredicate[T] {
	return func(ctx context.Context, x T) bool {
		return !p(ctx, x)
	}
}

// ContextAnd returns a predicate that is true if and only if all of the given
// predicates are true.
func ContextAnd[T any](predicates ...ContextPredicate[T]) ContextPredicate[T] {
	return func(ctx context.Context, x T) bool {
		for _, p := range predicates {
			if !p(ctx, x) {
				return false
			}
		}
		return true
	}
}

// ContextOr returns a predicate that is true if and only if at least one of the
// given predicates is true.
func ContextOr[T any](predicates ...ContextPredicate[T]) ContextPredicate[T] {
	return func(ctx context.Context, x T) bool {
		for _, p := range predicates {
			if p(ctx, x) {
				return true
			}
		}
		return false
	}
}

// Predicate is a boolean function of a single argument.
type Predicate[T any] func(T) bool

// Not returns the negation of the given predicate.
func Not[T any](p Predicate[T]) Predicate[T] {
	return func(x T) bool {
		return !p(x)
	}
}

// And returns a predicate that is true if and only if all of the given
// predicates are true.
func And[T any](predicates ...Predicate[T]) Predicate[T] {
	return func(x T) bool {
		for _, p := range predicates {
			if !p(x) {
				return false
			}
		}
		return true
	}
}

// Or returns a predicate that is true if and only if at least one of the given
// predicates is true.
func Or[T any](predicates ...Predicate[T]) Predicate[T] {
	return func(x T) bool {
		for _, p := range predicates {
			if p(x) {
				return true
			}
		}
		return false
	}
}

// ParallelAnd returns a predicate that is true if and only if all of the given
// predicates are true. The predicates are evaluated in parallel.
func ParallelAnd[T any](predicates ...ContextPredicate[T]) ContextPredicate[T] {
	return func(ctx context.Context, x T) bool {
		v := make([]bool, len(predicates))
		g, ctx := errgroup.WithContext(ctx)
		for i, p := range predicates {
			i, p := i, p
			g.Go(func() error {
				v[i] = p(ctx, x)
				return nil
			})
		}
		_ = g.Wait()
		for _, b := range v {
			if !b {
				return false
			}
		}
		return true
	}
}

// ParallelOr returns a predicate that is true if and only if at least one of the
// given predicates is true. The predicates are evaluated in parallel.
func ParallelOr[T any](predicates ...ContextPredicate[T]) ContextPredicate[T] {
	return func(ctx context.Context, x T) bool {
		v := make([]bool, len(predicates))
		g, ctx := errgroup.WithContext(ctx)
		for i, p := range predicates {
			i, p := i, p
			g.Go(func() error {
				v[i] = p(ctx, x)
				return nil
			})
		}
		_ = g.Wait()
		for _, b := range v {
			if b {
				return true
			}
		}
		return false
	}
}
