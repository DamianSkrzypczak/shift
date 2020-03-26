package shift

import chimiddleware "github.com/go-chi/chi/middleware"

const (
	TrailingSlashRedirect = "redirect" // Default trailing slash strategy
	TrailingSlashStrip    = "strip"
)

type AppConfig struct {
	Router RouterConfig
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		Router: RouterConfig{
			TrailingSlashStrategy: TrailingSlashRedirect,
		},
	}
}

type RouterConfig struct {
	TrailingSlashStrategy string
}

func (rconf RouterConfig) apply(r *Router) {
	switch rconf.TrailingSlashStrategy {
	case TrailingSlashRedirect:
		r.Use(FromHTTPMiddleware(chimiddleware.RedirectSlashes))
	case TrailingSlashStrip:
		r.Use(FromHTTPMiddleware(chimiddleware.StripSlashes))
	}
}
