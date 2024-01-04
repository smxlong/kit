package jwt

import (
	"context"
	"errors"
	"net/http"
	"strings"

	gojwt "github.com/golang-jwt/jwt/v5"
)

// Middleware is a middleware that extracts and validates a JWT from the
// Authorization header.
type Middleware struct {
	publicKey interface{}
	newClaims func() gojwt.Claims
	audience  string
	issuer    string
}

// Option is an option for NewMiddleware.
type Option func(*Middleware)

// WithKey sets the key to use for JWT validation. This can be the key itself, a
// keyfunc (see https://godoc.org/github.com/golang-jwt/jwt#Keyfunc), or a
// no-arg function that returns the key.
//
// If the key is a keyfunc or a no-arg function, it will be called for each
// request.
func WithKey(publicKey interface{}) Option {
	return func(m *Middleware) {
		m.publicKey = publicKey
	}
}

// WithNewClaims sets the function that returns a new claims instance. This
// allows the use of custom claims.
//
// Custom claims types must implement jwt.Claims. The simplest way to do this
// is to embed jwt.RegisteredClaims.
func WithNewClaims(newClaims func() gojwt.Claims) Option {
	return func(m *Middleware) {
		m.newClaims = newClaims
	}
}

// WithAudience sets the expected audience. If the audience is not empty, the
// token must contain an audience claim that matches.
func WithAudience(audience string) Option {
	return func(m *Middleware) {
		m.audience = audience
	}
}

// WithIssuer sets the expected issuer. If the issuer is not empty, the token
// must contain an issuer claim that matches.
func WithIssuer(issuer string) Option {
	return func(m *Middleware) {
		m.issuer = issuer
	}
}

// NewMiddleware returns a middleware. The middleware extracts and validates a
// JWT from the Authorization header, and stores the parsed claims in the
// request context under the key ContextKeyClaims.
func NewMiddleware(opts ...Option) *Middleware {
	m := &Middleware{
		newClaims: func() gojwt.Claims {
			return &gojwt.RegisteredClaims{}
		},
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Wrap the handler with middleware that extracts and validates a JWT from the
// Authorization header, and stores the parsed claims in the request context
// under the key ContextKeyClaims.
func (m *Middleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token, err := m.getToken(r); err == nil {
			if claims, err := m.parseClaims(token); err == nil {
				ctx := context.WithValue(r.Context(), ContextKeyClaims, claims)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// getToken extracts the token from the Authorization header.
func (m *Middleware) getToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return auth[7:], nil
	}
	return "", errors.New("no token")
}

// parseClaims parses the claims from the token. using the key.
func (m *Middleware) parseClaims(token string) (gojwt.Claims, error) {
	opts := []gojwt.ParserOption{}
	if m.audience != "" {
		opts = append(opts, gojwt.WithAudience(m.audience))
	}
	if m.issuer != "" {
		opts = append(opts, gojwt.WithIssuer(m.issuer))
	}

	claims := m.newClaims()
	if _, err := gojwt.ParseWithClaims(token, claims, m.keyFunc, opts...); err != nil {
		return nil, err
	}

	return claims, nil
}

// keyFunc returns the key to use for JWT validation.
func (m *Middleware) keyFunc(token *gojwt.Token) (interface{}, error) {
	if keyfunc, ok := m.publicKey.(func(*gojwt.Token) (interface{}, error)); ok {
		return keyfunc(token)
	}

	if keyfunc, ok := m.publicKey.(func() (interface{}, error)); ok {
		return keyfunc()
	}

	return m.publicKey, nil
}
