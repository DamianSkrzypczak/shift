package apptest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DamianSkrzypczak/shift"
	"github.com/gavv/httpexpect"
)

type DomainTest struct {
	*shift.App
	*httpexpect.Expect
}

func (at *DomainTest) Run() {}

func (at *DomainTest) Domain(constructor func(d *shift.Domain)) {
	at.App.Domain("/", constructor)
}

func NewDomainTest(t *testing.T) *DomainTest {
	app := shift.New("TestApp", nil)
	server := &TestServer{}
	app.Server = server

	if err := app.Run(""); err != nil {
		panic(err)
	}

	return &DomainTest{
		App:    app,
		Expect: httpexpect.New(t, server.URL),
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
