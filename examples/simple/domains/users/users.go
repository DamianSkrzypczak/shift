package users

import (
	"net/http"

	"github.com/DamianSkrzypczak/shift"
)

// Initialize users domain.
// - Setup routing
func Initialize(d *shift.Domain) {
	d.Router.Get("/{name}", getUser)
	d.Router.Post("/{name}", newUser)
}

// model representing single userModel entity.
type userModel struct {
	Avatar string `json:"avatar"`
}

// users storage.
// <nickname>:<user>.
var db = map[string]userModel{}

// User retrieval handler.
func getUser(rc shift.RequestContext) {
	user, ok := db[rc.Request.URLParam("name")]
	if !ok {
		rc.Response.SetStatusCode(http.StatusNotFound)
		return
	}

	if err := rc.Response.WithJSON(user, http.StatusOK); err != nil {
		rc.Response.SetStatusCode(http.StatusInternalServerError)
		return
	}
}

// New user creation handler.
func newUser(rc shift.RequestContext) {
	name := rc.Request.URLParam("name")

	newUser := userModel{}
	if err := rc.Request.JSON(&newUser); err != nil {
		rc.Response.SetStatusCode(http.StatusBadRequest)
		return
	}

	db[name] = newUser

	rc.Response.SetStatusCode(http.StatusAccepted)
}
