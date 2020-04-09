package movies

import (
	"net/http"
	"testing"
	"time"

	"github.com/DamianSkrzypczak/shift/apptest"
)

func setupClient(t *testing.T) *apptest.DomainTest {
	client := apptest.NewDomainTest(t) // Create and run test server
	client.Domain(Initialize)          // Register domain

	return client
}

func TestNoMoviesAtStart(t *testing.T) {
	c := setupClient(t)
	c.GET("/").Expect().JSON().Array().Empty()
}

func TestAddingMovie(t *testing.T) {
	c := setupClient(t)

	now := JSONTime(time.Now())
	movie := Movie{
		Title:    "TestMovie",
		Genre:    "Absurdist",
		Released: &now,
	}

	c.POST("/").WithJSON(movie).Expect().Status(http.StatusCreated)
	c.POST("/").WithJSON(movie).Expect().JSON().Object().Value("movie").Equal(movie)
}
