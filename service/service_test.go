package service

import (
	"context"
	"errors"
	"testing"

	"github.com/smxlong/kit/logger"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

// Test_Run_Error_Propagates tests that an error returned from the service
// function is propagated to the caller of Run.
func Test_Run_Error_Propagates(t *testing.T) {
	errTest := errors.New("test")
	err := Run(func(ctx context.Context, log logger.Logger) error {
		return errTest
	})
	require.Equal(t, errTest, err)
}

func Test_Run_Has_Context_And_Logger(t *testing.T) {
	err := Run(func(ctx context.Context, l logger.Logger) error {
		require.NotNil(t, ctx)
		require.NotNil(t, l)
		return nil
	})
	require.NoError(t, err)
}

func Test_Run_Debug(t *testing.T) {
	t.Setenv("DEBUG", "true")
	err := Run(func(ctx context.Context, l logger.Logger) error {
		return nil
	})
	require.NoError(t, err)
}

type testServiceFunctions struct {
	run             func(ctx context.Context, l logger.Logger) error
	bindFlags       func(flags *pflag.FlagSet)
	bindEnvironment func() error
}

// A service that only implements the Run method.
type testServiceRunOnly struct {
	testServiceFunctions
}

func (ts *testServiceRunOnly) Run(ctx context.Context, l logger.Logger) error {
	return ts.run(ctx, l)
}

// A service that implements the Run, BindFlags, and BindEnvironment methods.
type testServiceAll struct {
	testServiceFunctions
}

func (ts *testServiceAll) Run(ctx context.Context, l logger.Logger) error {
	return ts.run(ctx, l)
}

func (ts *testServiceAll) BindFlags(flags *pflag.FlagSet) {
	ts.bindFlags(flags)
}

func (ts *testServiceAll) BindEnvironment() error {
	return ts.bindEnvironment()
}

// A service that implements the Run and BindFlags methods.
type testServiceRunAndBindFlags struct {
	testServiceFunctions
}

func (ts *testServiceRunAndBindFlags) Run(ctx context.Context, l logger.Logger) error {
	return ts.run(ctx, l)
}

func (ts *testServiceRunAndBindFlags) BindFlags(flags *pflag.FlagSet) {
	ts.bindFlags(flags)
}

// A service that implements the Run and BindEnvironment methods.
type testServiceRunAndBindEnvironment struct {
	testServiceFunctions
}

func (ts *testServiceRunAndBindEnvironment) Run(ctx context.Context, l logger.Logger) error {
	return ts.run(ctx, l)
}

func (ts *testServiceRunAndBindEnvironment) BindEnvironment() error {
	return ts.bindEnvironment()
}

func Test_main_Function_Run_Works(t *testing.T) {
	var called bool
	s := &testServiceRunOnly{
		testServiceFunctions{
			run: func(ctx context.Context, l logger.Logger) error {
				called = true
				return nil
			},
		},
	}
	main("test", "test", s, func(int) {})
	require.True(t, called)
}

func Test_main_Function_BindFlags_Works(t *testing.T) {
	var called bool
	var bfcalled bool
	s := &testServiceRunAndBindFlags{
		testServiceFunctions{
			run: func(ctx context.Context, l logger.Logger) error {
				called = true
				return nil
			},
			bindFlags: func(flags *pflag.FlagSet) {
				bfcalled = true
			},
		},
	}
	main("test", "test", s, func(int) {})
	require.True(t, bfcalled)
	require.True(t, called)
}

func Test_main_Function_BindEnvironment_Works(t *testing.T) {
	var called bool
	var becalled bool
	s := &testServiceRunAndBindEnvironment{
		testServiceFunctions{
			run: func(ctx context.Context, l logger.Logger) error {
				called = true
				return nil
			},
			bindEnvironment: func() error {
				becalled = true
				return nil
			},
		},
	}
	main("test", "test", s, func(int) {})
	require.True(t, becalled)
	require.True(t, called)
}

func Test_main_Function_All_Works(t *testing.T) {
	var called bool
	var bfcalled bool
	var becalled bool
	var bfcalledFirst bool
	s := &testServiceAll{
		testServiceFunctions{
			run: func(ctx context.Context, l logger.Logger) error {
				called = true
				return nil
			},
			bindFlags: func(flags *pflag.FlagSet) {
				if !becalled {
					bfcalledFirst = true
				}
				bfcalled = true
			},
			bindEnvironment: func() error {
				becalled = true
				return nil
			},
		},
	}
	main("test", "test", s, func(int) {})
	require.True(t, bfcalled)
	require.True(t, becalled)
	require.True(t, bfcalledFirst)
	require.True(t, called)
}

func Test_main_Propagates_Error_From_BindEnvironment(t *testing.T) {
	errTest := errors.New("test")
	s := &testServiceRunAndBindEnvironment{
		testServiceFunctions{
			bindEnvironment: func() error {
				return errTest
			},
			run: func(ctx context.Context, l logger.Logger) error {
				return nil
			},
		},
	}
	var code int
	main("test", "test", s, func(c int) {
		code = c
	})
	require.Equal(t, 1, code)
}

func Test_main_Propagates_Error_From_Run(t *testing.T) {
	errTest := errors.New("test")
	s := &testServiceRunOnly{
		testServiceFunctions{
			run: func(ctx context.Context, l logger.Logger) error {
				return errTest
			},
		},
	}
	var code int
	main("test", "test", s, func(c int) {
		code = c
	})
	require.Equal(t, 1, code)
}

func Test_Main_Success(t *testing.T) {
	s := &testServiceRunOnly{
		testServiceFunctions{
			run: func(ctx context.Context, l logger.Logger) error {
				return nil
			},
		},
	}
	Main("test", "test", s)
}
