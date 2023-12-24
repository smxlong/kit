package rest

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_That_NewError_Returns_An_Error_With_The_Given_Message_And_Status_Code(t *testing.T) {
	t.Parallel()
	err := NewError("test", 1)
	assert.Equal(t, "test", err.Error())
	assert.Equal(t, 1, err.StatusCode())
}

func Test_That_Error_Implements_The_Error_Interface(t *testing.T) {
	t.Parallel()
	var err error = NewError("test", 1)
	assert.Equal(t, "test", err.Error())
}

func Test_That_Error_Implements_The_Is_Interface(t *testing.T) {
	t.Parallel()
	err := NewError("test", 1)
	assert.True(t, errors.Is(err, NewError("test", 1)))
	assert.False(t, errors.Is(err, NewError("test", 2)))
	assert.False(t, errors.Is(err, NewError("test2", 1)))
	assert.False(t, errors.Is(err, errors.New("test")))
}

func Test_That_Error_Is_Returns_False_When_The_Target_Is_Not_Equal_To_The_Error(t *testing.T) {
	t.Parallel()
	err := NewError("test", 1)
	assert.False(t, err.Is(NewError("test", 2)))
}

func Test_That_Error_Implements_The_Unwrap_Interface(t *testing.T) {
	t.Parallel()
	cause := errors.New("cause")
	err := NewError("test", 1).WithCause(cause)
	assert.Equal(t, cause, errors.Unwrap(err))
}

func Test_That_Error_Implements_The_StatusCode_Interface(t *testing.T) {
	t.Parallel()
	var err StatusCode = NewError("test", 1)
	assert.Equal(t, 1, err.StatusCode())
}

func Test_That_Error_Implements_The_Cause_Method(t *testing.T) {
	t.Parallel()
	cause := errors.New("cause")
	err := NewError("test", 1).WithCause(cause)
	assert.Equal(t, cause, err.Cause())
}
