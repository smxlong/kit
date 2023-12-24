package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

// Request is a REST request.
type Request interface{}

// Response is a REST response.
type Response interface{}

// ErrorResponse is a response that contains an error.
type ErrorResponse struct {
	// Error is the error.
	Error string `json:"error"`
}

// Handler is a REST handler.
type Handler struct {
	// NewRequest returns a new request for the handler.
	NewRequest func() Request
	// Handle handles the request.
	Handle func(context.Context, Request) Response
}

// Endpoint is the specification of a REST endpoint.
type Endpoint struct {
	// Method is the HTTP method of the endpoint.
	Method map[string]Handler
}

// Validate is implemented by Requests that can be validated.
type Validate interface {
	// Validate validates the request.
	Validate() error
}

// StatusCode is implemented by Responses that have a status code.
type StatusCode interface {
	// StatusCode returns the status code of the response.
	StatusCode() int
}

// ServeHTTP implements the http.Handler interface.
func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := e.Method[r.Method]
	if !ok {
		errorResponse(w, ErrNotSupported)
		return
	}
	req := handler.NewRequest()
	if err := decode(r, req); err != nil {
		errorResponse(w, err)
		return
	}
	if v, ok := req.(Validate); ok {
		if err := v.Validate(); err != nil {
			if _, ok := err.(StatusCode); !ok {
				err = ErrBadRequest.WithCause(err)
			}
			errorResponse(w, err)
			return
		}
	}
	ctx := r.Context()
	res := handler.Handle(ctx, req)
	if err, ok := res.(error); ok {
		errorResponse(w, err)
		return
	}
	statusCode := statusCodeOrDefault(http.StatusOK, res)
	encode(w, res, statusCode)
}

// errorResponse sends an error response.
func errorResponse(w http.ResponseWriter, err error) {
	statusCode := statusCodeOrDefault(http.StatusInternalServerError, err)
	e := &ErrorResponse{err.Error()}
	encode(w, e, statusCode)
}

// decode decodes the given http.Request to the given Request.
func decode(r *http.Request, req Request) error {
	if err := validateHeaders(r); err != nil {
		return err
	}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		if err == io.EOF {
			return ErrEmptyBody
		}
		return ErrBadRequest.WithCause(err)
	}
	return nil
}

// validateHeaders validates the headers of the given http.Request.
func validateHeaders(r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return ErrBadContentType
	}
	return nil
}

// encode encodes the given response to the given http.ResponseWriter.
func encode(w http.ResponseWriter, resp Response, status int) {
	setHeaders(w)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

// setHeaders sets the headers of the given http.ResponseWriter.
func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// statusCodeOrDefault returns the status code of the given response or the
// given status code if the response does not have a status code.
func statusCodeOrDefault(statusCode int, res interface{}) int {
	if v, ok := res.(StatusCode); ok {
		return v.StatusCode()
	}
	return statusCode
}
