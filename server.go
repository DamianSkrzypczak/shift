package shift

import "net/http"

type Server interface {
	ListenAndServe(addr string, router *Router) error
	ListenAndServeTLS(addr, certFile, keyFile string, router *Router) error
}

type defaultServer struct{}

func (s *defaultServer) ListenAndServe(addr string, router *Router) error {
	return http.ListenAndServe(addr, router)
}

func (s *defaultServer) ListenAndServeTLS(addr, certFile, keyFile string, router *Router) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, router)
}
