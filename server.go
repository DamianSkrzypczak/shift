package shift

import "net/http"

type Server interface {
	Initialize(addr string, router *Router)
	ListenAndServe() error
	ListenAndServeTLS(certFile, keyFile string) error
}

type defaultServer struct {
	http.Server
}

func (s *defaultServer) Initialize(addr string, router *Router) {
	s.Addr = addr
	s.Handler = router
}
