package shift

import (
	"net/http"
)

type Handler func(RequestContext)

// ToHTTTP adapts Handler interface to http.HandlerFunc interface
func (h Handler) ToHTTP() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		h(contextConstructor(NewRequestContext).FromHTTP(resp, req))
	}
}

// FromHTTP produces new Handler from http.HandlerFunc instance
func FromHTTPHandler(h http.Handler) Handler {
	return func(ctx RequestContext) {
		h.ServeHTTP(ctx.Response, ctx.Request.Request)
	}
}

// FromHTTP produces new Handler from http.HandlerFunc instance
func FromHTTPHandlerFunc(h http.HandlerFunc) Handler {
	return func(ctx RequestContext) {
		h(ctx.Response, ctx.Request.Request)
	}
}
