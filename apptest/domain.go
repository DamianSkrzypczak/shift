package apptest

import (
	"net/http"
	"net/http/httptest"

	"github.com/DamianSkrzypczak/shift"
)

type DomainTest struct {
	*shift.App
	URL string
}

func NewDomainTest(appName string, constructor func(d *shift.Domain)) *DomainTest {
	app := shift.New(appName, nil)
	app.Domain("/", constructor)

	server := &TestServer{}
	app.Server = server

	if err := app.Run(""); err != nil {
		panic(err)
	}

	return &DomainTest{
		App: app,
		URL: server.URL,
	}
}

type TestServer struct {
	*httptest.Server
}

func (ts *TestServer) Initialize(_ string, router *shift.Router) {
	ts.Server = httptest.NewServer(http.Handler(router))
}
func (ts *TestServer) ListenAndServe() error                            { return nil }
func (ts *TestServer) ListenAndServeTLS(certFile, keyFile string) error { return nil }
