package main

import (
	"github.com/lightsaid/gotk/mux"
)

func (app *application) setupRoute() *mux.ServeMux {
	r := mux.NewServeMux()

	r.GET("/", app.indexHandler)

	r.GET("/sign", app.signHandler)

	r.GET("/signin", app.signinHandler)

	r.GET("/signup", app.signupHandler)

	return r
}
