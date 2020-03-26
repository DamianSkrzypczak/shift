package main

import (
	"net/http"
	"simple/domains/users"

	"github.com/DamianSkrzypczak/shift"
)

func main() {
	app := shift.New("simpleApp", nil)

	app.Router.Get("/status", func(rc shift.RequestContext) {
		rc.Response.SetStatusCode(http.StatusOK)
		_ = rc.Response.SetBodyJSON(map[string]string{"status": "healthy"})
	})

	app.Domain("/users", users.Initialize)

	if err := app.Run(":8000"); err != nil {
		panic(err)
	}
}
