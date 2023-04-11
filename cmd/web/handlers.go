package main

import (
	"net/http"
)

func (app *application) indexHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "index.page.html")
}

func (app *application) signHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "sign.page.html")
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "sign.page.html")
}

func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "sign.page.html")
}

func (app *application) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "404.page.html")
}

func (app *application) serverErrorHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "500.page.html")
}

func (app *application) forgotHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "forgot.page.html")
}

func (app *application) resetHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "reset.page.html")
}
