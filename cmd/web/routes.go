package main

import (
	"github.com/lightsaid/gotk/mux"
)

func (app *application) setupRoute() *mux.ServeMux {
	r := mux.NewServeMux()

	app.showpages(r)

	r.POST("/login", app.loginHandler)
	r.POST("/resister", app.registerHandler)

	// 静态资源
	r.Static("/static/", "./static")

	return r
}

func (app *application) showpages(r *mux.ServeMux) {
	r.GET("/", app.indexHandler)
	r.GET("/sign", app.signHandler)
	r.GET("/forgot", app.forgotHandler)
	r.GET("/reset", app.resetHandler)
	r.GET("/notfound", app.notFoundHandler)
	r.GET("/servererror", app.serverErrorHandler)
}
