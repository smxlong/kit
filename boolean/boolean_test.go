package boolean

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_ContextNot tests the ContextNot function.
func Test_ContextNot(t *testing.T) {
	ctx := context.Background()
	var gotCtx context.Context
	// Create a context predicate that is always true.
	p := func(ctx context.Context, v int) bool {
		gotCtx = ctx
		return true
	}
	// Create a context predicate that is always false.
	np := ContextNot(p)
	require.False(t, np(ctx, 0))
	require.Equal(t, ctx, gotCtx)
}

// Test_ContextAnd tests the ContextAnd function.
func Test_ContextAnd(t *testing.T) {
	// Eight test cases with three input values and one expected result.
	tc := [][]int{
		{0, 0, 0, 0},
		{0, 0, 1, 0},
		{0, 1, 0, 0},
		{0, 1, 1, 0},
		{1, 0, 0, 0},
		{1, 0, 1, 0},
		{1, 1, 0, 0},
		{1, 1, 1, 1},
	}

	for _, c := range tc {
		ctx := context.Background()
		var gotCtx context.Context

		p := []ContextPredicate[int]{
			func(ctx context.Context, _ int) bool {
				gotCtx = ctx
				return false
			},
			func(ctx context.Context, _ int) bool {
				gotCtx = ctx
				return true
			},
		}

		pp := []ContextPredicate[int]{}
		for _, v := range c {
			pp = append(pp, p[v])
		}
		require.Equal(t, c[3] == 1, ContextAnd(pp...)(ctx, 0))
		require.Equal(t, ctx, gotCtx)
	}
}

// Test_ContextOr tests the ContextOr function.
func Test_ContextOr(t *testing.T) {
	// Eight test cases with three input values and one expected result.
	tc := [][]int{
		{0, 0, 0, 0},
		{0, 0, 1, 1},
		{0, 1, 0, 1},
		{0, 1, 1, 1},
		{1, 0, 0, 1},
		{1, 0, 1, 1},
		{1, 1, 0, 1},
		{1, 1, 1, 1},
	}

	for _, c := range tc {
		ctx := context.Background()
		var gotCtx context.Context

		p := []ContextPredicate[int]{
			func(ctx context.Context, _ int) bool {
				gotCtx = ctx
				return false
			},
			func(ctx context.Context, _ int) bool {
				gotCtx = ctx
				return true
			},
		}

		pp := []ContextPredicate[int]{}
		for _, v := range c {
			pp = append(pp, p[v])
		}
		require.Equal(t, c[3] == 1, ContextOr(pp...)(ctx, 0))
		require.Equal(t, ctx, gotCtx)
	}
}

// Test_Not tests the Not function.
func Test_Not(t *testing.T) {
	fa := func(int) bool { return false }
	tr := func(int) bool { return true }
	require.True(t, Not(fa)(0))
	require.False(t, Not(tr)(0))
}

// Test_And tests the And function.
func Test_And(t *testing.T) {
	// Eight test cases with three input values and one expected result.
	tc := [][]int{
		{0, 0, 0, 0},
		{0, 0, 1, 0},
		{0, 1, 0, 0},
		{0, 1, 1, 0},
		{1, 0, 0, 0},
		{1, 0, 1, 0},
		{1, 1, 0, 0},
		{1, 1, 1, 1},
	}

	for _, c := range tc {
		called := []bool{false, false, false}

		pp := []Predicate[int]{}
		for i, v := range c[:3] {
			i, v := i, v
			pp = append(pp, func(int) bool {
				called[i] = true
				return v == 1
			})
		}
		require.Equal(t, c[3] == 1, And(pp...)(0), "tc=%v", c)

		// should be called up until the first false
		x := true
		for i := range c[:3] {
			require.Equal(t, x, called[i], "tc=%v", c)
			if c[i] == 0 {
				x = false
			}
		}
	}
}

// Test_Or tests the Or function.
func Test_Or(t *testing.T) {
	// Eight test cases with three input values and one expected result.
	tc := [][]int{
		{0, 0, 0, 0},
		{0, 0, 1, 1},
		{0, 1, 0, 1},
		{0, 1, 1, 1},
		{1, 0, 0, 1},
		{1, 0, 1, 1},
		{1, 1, 0, 1},
		{1, 1, 1, 1},
	}

	for _, c := range tc {
		called := []bool{false, false, false}

		pp := []Predicate[int]{}
		for i, v := range c[:3] {
			i, v := i, v
			pp = append(pp, func(int) bool {
				called[i] = true
				return v == 1
			})
		}
		require.Equal(t, c[3] == 1, Or(pp...)(0), "tc=%v", c)

		// should be called up until the first true
		x := true
		for i := range c[:3] {
			require.Equal(t, x, called[i], "tc=%v", c)
			if c[i] == 1 {
				x = false
			}
		}
	}
}

// Test_ParallelAnd tests the ParallelAnd function.
func Test_ParallelAnd(t *testing.T) {
	// Eight test cases with three input values and one expected result.
	tc := [][]int{
		{0, 0, 0, 0},
		{0, 0, 1, 0},
		{0, 1, 0, 0},
		{0, 1, 1, 0},
		{1, 0, 0, 0},
		{1, 0, 1, 0},
		{1, 1, 0, 0},
		{1, 1, 1, 1},
	}

	for _, c := range tc {
		called := []bool{false, false, false}

		pp := []ContextPredicate[int]{}
		for i, v := range c[:3] {
			i, v := i, v
			pp = append(pp, func(ctx context.Context, _ int) bool {
				called[i] = true
				return v == 1
			})
		}

		ctx := context.Background()
		require.Equal(t, c[3] == 1, ParallelAnd(pp...)(ctx, 0), "tc=%v", c)

		for i := range c[:3] {
			require.True(t, called[i], "tc=%v", c)
		}
	}
}

// Test_ParallelOr tests the ParallelOr function.
func Test_ParallelOr(t *testing.T) {
	// Eight test cases with three input values and one expected result.
	tc := [][]int{
		{0, 0, 0, 0},
		{0, 0, 1, 1},
		{0, 1, 0, 1},
		{0, 1, 1, 1},
		{1, 0, 0, 1},
		{1, 0, 1, 1},
		{1, 1, 0, 1},
		{1, 1, 1, 1},
	}

	for _, c := range tc {
		called := []bool{false, false, false}

		pp := []ContextPredicate[int]{}
		for i, v := range c[:3] {
			i, v := i, v
			pp = append(pp, func(ctx context.Context, _ int) bool {
				called[i] = true
				return v == 1
			})
		}

		ctx := context.Background()
		require.Equal(t, c[3] == 1, ParallelOr(pp...)(ctx, 0), "tc=%v", c)

		for i := range c[:3] {
			require.True(t, called[i], "tc=%v", c)
		}
	}
}
