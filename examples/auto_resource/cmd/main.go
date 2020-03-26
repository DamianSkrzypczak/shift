package main

import (
	"app/domains/movies"

	"github.com/DamianSkrzypczak/shift"
)

func main() {
	app := shift.New("MoviesAPI", nil)
	app.Domain("/movies", movies.Initialize)

	if err := app.Run(":8000"); err != nil {
		panic(err)
	}
}
