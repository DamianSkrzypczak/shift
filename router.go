package shift

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Router struct {
	_router chi.Router
}

func newRouter() *Router {
	return &Router{
		_router: chi.NewRouter(),
	}
}

func (r *Router) mount(pattern string, router *Router) {
	r._router.Mount(pattern, router._router)
}

func (r *Router) Connect(path string, h Handler) { r._router.Connect(path, h.ToHTTP()) }
func (r *Router) Delete(path string, h Handler)  { r._router.Delete(path, h.ToHTTP()) }
func (r *Router) Get(path string, h Handler)     { r._router.Get(path, h.ToHTTP()) }
func (r *Router) Head(path string, h Handler)    { r._router.Head(path, h.ToHTTP()) }
func (r *Router) Options(path string, h Handler) { r._router.Options(path, h.ToHTTP()) }
func (r *Router) Patch(path string, h Handler)   { r._router.Patch(path, h.ToHTTP()) }
func (r *Router) Post(path string, h Handler)    { r._router.Post(path, h.ToHTTP()) }
func (r *Router) Put(path string, h Handler)     { r._router.Put(path, h.ToHTTP()) }
func (r *Router) Trace(path string, h Handler)   { r._router.Trace(path, h.ToHTTP()) }

func (r *Router) Method(method, path string, h Handler) { r._router.Method(method, path, h.ToHTTP()) }

func (r *Router) NotFound(h Handler)         { r._router.NotFound(h.ToHTTP()) }
func (r *Router) MethodNotAllowed(h Handler) { r._router.MethodNotAllowed(h.ToHTTP()) }

func (r *Router) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	r._router.ServeHTTP(responseWriter, request)
}

// Use pushes given middleware onto router stack
func (r *Router) Use(middlewares ...Middleware) { r._router.Use(Middlewares(middlewares).ToHTTP()...) }

// Use pushes inlined middleware onto router stack.
func (r *Router) With(middlewares ...Middleware) *Router {
	return &Router{r._router.With(Middlewares(middlewares).ToHTTP()...)}
}

// UseContext works like Use but provides interface for contextProvider middleware
func (r *Router) UseContext(contextProvider ContextProvider) { r.Use(contextProvider.toMiddleware()) }

// WithContext works like With but provides interface for contextProvider middleware
func (r *Router) WithContext(contextProvider ContextProvider) *Router {
	return r.With(contextProvider.toMiddleware())
}

// UseContext works like Use but provides interface for Interceptor middleware
func (r *Router) UseInterceptor(interceptor Interceptor) { r.Use(interceptor.toMiddleware()) }

// WithContext works like With but provides interface for Interceptor middleware
func (r *Router) WithInterceptor(interceptor Interceptor) *Router {
	middleware := interceptor.toMiddleware()
	return r.With(middleware)
}
