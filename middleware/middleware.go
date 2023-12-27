package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// Chain chains the given middleware together, in order.
func Chain(middlewares ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		return h
	}
}

// ChainFunc chains the given middleware functions together, in order.
func ChainFunc(middlewares ...func(http.Handler) http.Handler) Middleware {
	return func(h http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		return h
	}
}
