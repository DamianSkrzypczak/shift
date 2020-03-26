package shift

import (
	"context"
	"net/http"
)

type Middleware func(next Handler) Handler

func (m Middleware) ToHTTP() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		handler := FromHTTPHandler(h)
		return m(handler).ToHTTP()
	}
}

// FromHTTPMiddleware produces new Middleware from http-based middleware
func FromHTTPMiddleware(m func(h http.Handler) http.Handler) Middleware {
	return func(next Handler) Handler {
		return FromHTTPHandler(m(next.ToHTTP()))
	}
}

type Middlewares []Middleware

func (ms Middlewares) ToHTTP() []func(h http.Handler) http.Handler {
	httpMiddlewares := []func(h http.Handler) http.Handler{}
	for _, m := range ms {
		httpMiddlewares = append(httpMiddlewares, m.ToHTTP())
	}

	return httpMiddlewares
}

type ContextProvider func(*Request) context.Context

func (ctxProv ContextProvider) toMiddleware() Middleware {
	return func(next Handler) Handler {
		return func(requestContext RequestContext) {
			newCtx := ctxProv(requestContext.Request)
			requestContext.Request.Request = requestContext.Request.WithContext(newCtx)
			next(requestContext)
		}
	}
}

type Interceptor func() (before, after func(RequestContext))

func (i Interceptor) toMiddleware() Middleware {
	return func(h Handler) Handler {
		return func(requestContext RequestContext) {
			before, after := i()
			before(requestContext)
			h(requestContext)
			after(requestContext)
		}
	}
}
