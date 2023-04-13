package main

import (
	"log"
	"net/http"
	"time"
)

func (app *application) loadSessionAndSave(next http.Handler) http.Handler {
	return app.sessionMgr.LoadAndSave(next)
}

func (app *application) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s \n", r.Method, r.RequestURI, time.Since(t))
	})
}
