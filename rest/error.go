package rest

import "net/http"

// Error is an error from an endpoint.
type Error struct {
	message    string
	statusCode int
	cause      error
}

// NewError returns a new Error.
func NewError(message string, statusCode int) *Error {
	return &Error{
		message:    message,
		statusCode: statusCode,
	}
}

// Error implements the error interface.
func (e *Error) Error() string {
	return e.message
}

// Is implements the errors.Is interface.
func (e *Error) Is(target error) bool {
	if t, ok := target.(*Error); ok {
		return e.message == t.message && e.statusCode == t.statusCode
	}
	return false
}

// Unwrap implements the errors.Unwrap interface.
func (e *Error) Unwrap() error {
	return e.cause
}

// StatusCode implements the StatusCode interface.
func (e *Error) StatusCode() int {
	return e.statusCode
}

// Cause returns the cause of the error.
func (e *Error) Cause() error {
	return e.cause
}

// WithCause returns a new Error with the given cause.
func (e *Error) WithCause(cause error) *Error {
	return &Error{
		message:    e.message,
		statusCode: e.statusCode,
		cause:      cause,
	}
}

var (
	// ErrBadContentType is returned when a request has a bad content type.
	ErrBadContentType = NewError("bad content type", http.StatusUnsupportedMediaType)
	// ErrBadRequest is returned when a request is bad.
	ErrBadRequest = NewError("bad request", http.StatusBadRequest)
	// ErrEmptyBody is returned when a request has an empty body.
	ErrEmptyBody = NewError("empty body", http.StatusBadRequest)
	// ErrForbidden is returned when a request is forbidden.
	ErrForbidden = NewError("forbidden", http.StatusForbidden)
	// ErrInternal is returned when an internal error occurs.
	ErrInternal = NewError("internal error", http.StatusInternalServerError)
	// ErrNotFound is returned when a resource is not found.
	ErrNotFound = NewError("not found", http.StatusNotFound)
	// ErrNotSupported is returned when a method is not supported.
	ErrNotSupported = NewError("not supported", http.StatusMethodNotAllowed)
	// ErrUnauthorized is returned when a request is unauthorized.
	ErrUnauthorized = NewError("unauthorized", http.StatusUnauthorized)
	// ErrConflict is returned when a request causes a conflict.
	ErrConflict = NewError("conflict", http.StatusConflict)
)
