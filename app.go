package shift

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type App struct {
	Server     Server
	Logger     *logrus.Logger
	rootDomain *Domain
	Router     *Router
	Name       string
	Config     *AppConfig
}

func New(name string, config *AppConfig) *App {
	if config == nil {
		config = NewAppConfig()
	}

	rootDomain := newDomain("/", nil)
	config.Router.apply(rootDomain.Router)

	app := App{
		Name:       name,
		Server:     &defaultServer{},
		Logger:     logrus.New(),
		rootDomain: rootDomain,
		Router:     rootDomain.Router,
		Config:     config,
	}

	return &app
}

func (app *App) Domain(path string, constructor func(d *Domain)) {
	constructor(newDomain(path, app.rootDomain))
}

func (app *App) Run(addr string) error {
	app.Logger.SetLevel(logrus.DebugLevel)
	app.Logger.WithField("Address", addr).Infof("Serving shift app \"%s\"", app.Name)

	return http.ListenAndServe(addr, app.Router)
}

func (app *App) RunTLS(addr, certFile, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, app.Router)
}
