package main

import (
	"github.com/DamianSkrzypczak/shift"
)

func main() {
	app := shift.New("minimalApp", nil)

	app.Router.Get("/", func(rc shift.RequestContext) {
		rc.Response.SetBodyString("HELLO WORLD!")
	})

	if err := app.Run(":8000"); err != nil {
		panic(err)
	}
}
