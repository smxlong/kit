package service

import (
	"context"
	"errors"
	"testing"

	"github.com/smxlong/kit/logger"
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
