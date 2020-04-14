package movies

import (
	"net/http"
	"testing"
	"time"

	"github.com/DamianSkrzypczak/shift/apptest"
	"github.com/gavv/httpexpect"
)

func TestNoMoviesAtStart(t *testing.T) {
	dt := apptest.NewDomainTest("NoMovies", Initialize) // Create and run test server for domain
	e := httpexpect.New(t, dt.URL).GET("/").Expect()

	e.JSON().Array().Empty()
}

func TestAddingMovie(t *testing.T) {
	dt := apptest.NewDomainTest("AddingMovie", Initialize) // Create and run test server for domain

	now := JSONTime(time.Now())
	movie := Movie{
		Title:    "TestMovie",
		Genre:    "Absurdist",
		Released: &now,
	}

	e := httpexpect.New(t, dt.URL).POST("/").WithJSON(movie).Expect()
	e.JSON().Object().Value("movie").Equal(movie)
	e.Status(http.StatusCreated)
}

func TestGettingMovie(t *testing.T) {
	now := JSONTime(time.Now())
	movie := Movie{
		Title:    "TestMovie",
		Genre:    "Absurdist",
		Released: &now,
	}

	movies = []Movie{movie}

	dt := apptest.NewDomainTest("GettingMovie", Initialize) // Create and run test server for domain
	e := httpexpect.New(t, dt.URL).GET("/0").Expect()

	e.Status(http.StatusOK)
	e.JSON().Object().Equal(movie)
}

func TestGettingMovieOutOfRange(t *testing.T) {
	dt := apptest.NewDomainTest("GettingMovie", Initialize) // Create and run test server for domain
	e := httpexpect.New(t, dt.URL).GET("/99").Expect()

	e.Status(http.StatusNotFound)
	e.Body().Empty()
}

func TestGettingInvalidIDMovie(t *testing.T) {
	dt := apptest.NewDomainTest("GettingMovie", Initialize) // Create and run test server for domain
	e := httpexpect.New(t, dt.URL).GET("/aaa").Expect()

	e.Status(http.StatusNotFound)
	e.Body().Empty()
}

func TestRemovingExistingMovie(t *testing.T) {
	now := JSONTime(time.Now())
	movie := Movie{
		Title:    "TestMovie",
		Genre:    "Absurdist",
		Released: &now,
	}

	movies = []Movie{movie}

	dt := apptest.NewDomainTest("GettingMovie", Initialize) // Create and run test server for domain
	e := httpexpect.New(t, dt.URL).DELETE("/0").Expect()

	e.Status(http.StatusNoContent)
	e.Body().Empty()
}

func TestRemovingNonExistingMovie(t *testing.T) {
	dt := apptest.NewDomainTest("GettingMovie", Initialize) // Create and run test server for domain
	e := httpexpect.New(t, dt.URL).DELETE("/99").Expect()

	e.Status(http.StatusNotFound)
	e.Body().Empty()
}
