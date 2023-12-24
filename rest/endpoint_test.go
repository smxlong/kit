package rest

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TEmptyRequest is a test request.
type TEmptyRequest struct{}

// TEmptyResponse is a test response.
type TEmptyResponse struct{}

// TRequest is a test request.
type TRequest struct {
	Message         string `json:"message"`
	Number          int    `json:"number"`
	validationError error
}

// Validate implements the Validate interface.
func (r *TRequest) Validate() error {
	if r.validationError != nil {
		return r.validationError
	}
	// Number must be even and > 0
	if r.Number%2 != 0 || r.Number <= 0 {
		return ErrBadRequest.WithCause(errors.New("number must be even and > 0"))
	}
	// Message must have a 'q' in it
	if !strings.Contains(r.Message, "q") {
		return ErrBadRequest.WithCause(errors.New("message must contain a 'q'"))
	}
	return nil
}

// TResponse is a test response.
type TResponse struct {
	ResponseMessage string `json:"response_message"`
	ResponseNumber  int    `json:"response_number"`
}

// TResponseWithStatusCodeAndError is a test response.
type TResponseWithStatusCodeAndError struct {
	ResponseMessage string `json:"response_message"`
	ResponseNumber  int    `json:"response_number"`
	status          int
	err             error
}

// StatusCode implements the StatusCode interface.
func (r *TResponseWithStatusCodeAndError) StatusCode() int {
	return r.status
}

// Error implements the error interface.
func (r *TResponseWithStatusCodeAndError) Error() string {
	return r.err.Error()
}

// Is implements the errors.Is interface.
func (r *TResponseWithStatusCodeAndError) Is(target error) bool {
	return errors.Is(r.err, target)
}

// TResponseWithStatusCode is a test response.
type TResponseWithStatusCode struct {
	ResponseMessage string `json:"response_message"`
	ResponseNumber  int    `json:"response_number"`
	status          int
}

// StatusCode implements the StatusCode interface.
func (r *TResponseWithStatusCode) StatusCode() int {
	return r.status
}

//////////////////////////////////////////////////////////////////////////////
// errorResponse tests

func Test_That_errorResponse_Returns_A_JSON_Response(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	errorResponse(rec, NewError("test", 400))
	assert.Equal(t, 400, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "{\"error\":\"test\"}\n", rec.Body.String())
	// Test with an error that lacks a statuscode
	rec = httptest.NewRecorder()
	errorResponse(rec, errors.New("test"))
	assert.Equal(t, 500, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "{\"error\":\"test\"}\n", rec.Body.String())
}

//////////////////////////////////////////////////////////////////////////////
// decode tests

func Test_That_decode_Returns_Error_For_Empty_Request_Body(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	assert.Equal(t, ErrEmptyBody, decode(req, &TEmptyRequest{}))
}

func Test_That_decode_Returns_Error_For_Incorrect_Content_Type(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "/test", nil)
	assert.Equal(t, ErrBadContentType, decode(req, &TEmptyRequest{}))
}

func Test_That_decode_Returns_Error_For_Invalid_Request(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBufferString(`{`))
	assert.Equal(t, ErrBadRequest.WithCause(io.ErrUnexpectedEOF), decode(req, &TEmptyRequest{}))
}

func Test_That_decode_Returns_No_Error_For_Correct_Content_Type(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBufferString(`{"message":"test","number":1}`))
	out := &TRequest{}
	assert.NoError(t, decode(req, out))
	assert.Equal(t, "test", out.Message)
	assert.Equal(t, 1, out.Number)
}

//////////////////////////////////////////////////////////////////////////////
// validateHeaders tests

func Test_That_validateHeaders_Returns_Error_For_Incorrect_Content_Type(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "/test", nil)
	assert.Equal(t, ErrBadContentType, validateHeaders(req))
}

func Test_That_validateHeaders_Returns_No_Error_For_Correct_Content_Type(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, validateHeaders(req))
}

//////////////////////////////////////////////////////////////////////////////
// Encode tests

func Test_That_Encode_Sets_The_Correct_Content_Type(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	encode(rec, &TEmptyResponse{}, 200)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

func Test_That_Encode_Sets_The_Correct_Status_Code(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	encode(rec, &TEmptyResponse{}, 201)
	assert.Equal(t, 201, rec.Code)
}

func Test_That_Encode_Produces_JSON(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	encode(rec, &TResponse{"test", 1}, 200)
	assert.Equal(t, "{\"response_message\":\"test\",\"response_number\":1}\n", rec.Body.String())
}

//////////////////////////////////////////////////////////////////////////////
// setHeaders tests

func Test_That_setHeaders_Sets_The_Correct_Content_Type(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	setHeaders(rec)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
}

//////////////////////////////////////////////////////////////////////////////
// statusCodeOrDefault tests

func Test_That_statusCodeOrDefault_Returns_The_Status_Code_Of_The_Response(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 200, statusCodeOrDefault(201, &TResponseWithStatusCode{"test", 1, 200}))
}

func Test_That_statusCodeOrDefault_Returns_The_Default_Status_Code(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 201, statusCodeOrDefault(201, &TResponse{"test", 1}))
}

//////////////////////////////////////////////////////////////////////////////
// ServeHTTP tests

func Test_That_ServeHTTP_Returns_Error_For_Unsupported_Method(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	ep := &Endpoint{
		Method: map[string]Handler{
			"POST": {
				NewRequest: func() Request {
					return &TEmptyRequest{}
				},
				Handle: func(ctx context.Context, req Request) Response {
					return &TEmptyResponse{}
				},
			},
		},
	}
	ep.ServeHTTP(rec, req)
	assert.Equal(t, 405, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "{\"error\":\"not supported\"}\n", rec.Body.String())
}

func Test_That_ServeHTTP_Returns_Error_For_Invalid_Request(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBufferString(`{`))
	ep := &Endpoint{
		Method: map[string]Handler{
			"POST": {
				NewRequest: func() Request {
					return &TEmptyRequest{}
				},
				Handle: func(ctx context.Context, req Request) Response {
					return &TEmptyResponse{}
				},
			},
		},
	}
	ep.ServeHTTP(rec, req)
	assert.Equal(t, 400, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "{\"error\":\"bad request\"}\n", rec.Body.String())
}

func Test_That_ServeHTTP_Returns_Error_For_Invalid_Request_With_StatusCode(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBufferString(`{"message":"test","number":1}`))
	ep := &Endpoint{
		Method: map[string]Handler{
			"POST": {
				NewRequest: func() Request {
					return &TRequest{}
				},
				Handle: func(ctx context.Context, req Request) Response {
					return &TResponse{}
				},
			},
		},
	}
	ep.ServeHTTP(rec, req)
	assert.Equal(t, 400, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "{\"error\":\"bad request\"}\n", rec.Body.String())
}

func Test_That_ServeHTTP_Returns_Error_For_Invalid_Request_Without_StatusCode(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBufferString(`{"message":"test","number":1}`))
	ep := &Endpoint{
		Method: map[string]Handler{
			"POST": {
				NewRequest: func() Request {
					return &TRequest{
						validationError: errors.New("test"),
					}
				},
				Handle: func(ctx context.Context, req Request) Response {
					return &TEmptyResponse{}
				},
			},
		},
	}
	ep.ServeHTTP(rec, req)
	assert.Equal(t, 400, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "{\"error\":\"bad request\"}\n", rec.Body.String())
}

func Test_That_ServeHTTP_Returns_Response_For_Valid_Request(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBufferString(`{"message":"testq","number":2}`))
	ep := &Endpoint{
		Method: map[string]Handler{
			"POST": {
				NewRequest: func() Request {
					return &TRequest{}
				},
				Handle: func(ctx context.Context, req Request) Response {
					return &TResponse{req.(*TRequest).Message[:4], req.(*TRequest).Number - 1}
				},
			},
		},
	}
	ep.ServeHTTP(rec, req)
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "{\"response_message\":\"test\",\"response_number\":1}\n", rec.Body.String())
}

func Test_That_ServeHTTP_Returns_Response_For_Valid_Request_With_StatusCode(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBufferString(`{"message":"testq","number":2}`))
	ep := &Endpoint{
		Method: map[string]Handler{
			"POST": {
				NewRequest: func() Request {
					return &TRequest{}
				},
				Handle: func(ctx context.Context, req Request) Response {
					return &TResponseWithStatusCode{req.(*TRequest).Message[:4], req.(*TRequest).Number - 1, 201}
				},
			},
		},
	}
	ep.ServeHTTP(rec, req)
	assert.Equal(t, 201, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "{\"response_message\":\"test\",\"response_number\":1}\n", rec.Body.String())
}

func Test_That_ServeHTTP_Returns_Response_For_Valid_Request_With_StatusCode_With_Error(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBufferString(`{"message":"testq","number":2}`))
	ep := &Endpoint{
		Method: map[string]Handler{
			"POST": {
				NewRequest: func() Request {
					return &TRequest{}
				},
				Handle: func(ctx context.Context, req Request) Response {
					return &TResponseWithStatusCodeAndError{req.(*TRequest).Message[:4], req.(*TRequest).Number - 1, 421, errors.New("test")}
				},
			},
		},
	}
	ep.ServeHTTP(rec, req)
	assert.Equal(t, 421, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "{\"error\":\"test\"}\n", rec.Body.String())
}

func Test_That_ServeHTTP_Returns_Response_For_Valid_Request_With_Error(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(bytes.NewBufferString(`{"message":"testq","number":2}`))
	ep := &Endpoint{
		Method: map[string]Handler{
			"POST": {
				NewRequest: func() Request {
					return &TRequest{}
				},
				Handle: func(ctx context.Context, req Request) Response {
					return errors.New("test")
				},
			},
		},
	}
	ep.ServeHTTP(rec, req)
	assert.Equal(t, 500, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "{\"error\":\"test\"}\n", rec.Body.String())
}
