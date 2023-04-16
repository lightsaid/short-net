package main

import (
	"github.com/lightsaid/gotk/mux"
)

func (app *application) setupRoute() *mux.ServeMux {
	r := mux.NewServeMux()

	// NOTE: 使用 github.com/alexedwards/scs/v2 要执行 LoadAndSave()，
	// 否则会各种 panic
	r.Use(app.loadSessionAndSave)
	r.Use(app.csrfMiddleware)
	// r.Use(app.recovererMiddleware)
	r.Use(app.loggerMiddleware)

	app.showpages(r)
	app.signLogicHandler(r)
	app.userLogicHandler(r)
	app.shortLogicHandler(r)

	// 静态资源
	r.Static("/static/", "./static")

	return r
}

// showpages 页面访问
func (app *application) showpages(r *mux.ServeMux) {
	r.GET("/", app.indexHandler)
	r.GET("/sign", app.signHandler)
	r.GET("/forgot", app.forgotHandler)
	r.GET("/reset", app.resetHandler).Use(app.authRequired)
	r.GET("/notfound", app.notFoundHandler)
	r.GET("/servererror", app.serverErrorHandler)
	r.GET("/success", app.operateSuccessfully)
	r.GET("/error", app.errorHandler)
}

// 登录注册逻辑
func (app *application) signLogicHandler(r *mux.ServeMux) {
	r.POST("/login", app.loginHandler)
	r.POST("/register", app.registerHandler)
	r.GET("/activate/:token", app.activateHandler)
}

// userLogicHandler 用户handler
func (app *application) userLogicHandler(r *mux.ServeMux) {
	m := r.RouteGroup("")
	// TODO:
	m.POST("/forgot", app.indexHandler)
	m.POST("/reset", app.indexHandler)
	m.POST("/profile", app.indexHandler)
	m.PUT("/profile", app.indexHandler)
}

// shortLogicHandler 短网址逻辑相关
func (app *application) shortLogicHandler(r *mux.ServeMux) {
	r.GET("/:hash|^[a-zA-Z0-9]+$", app.redirectLinkHandler)

	s := r.RouteGroup("/short")
	s.Use(app.authRequired)
	s.POST("/create", app.createLinkHandler)
	s.POST("/update", app.updateLinkHandler)
	s.POST("/delete", app.deleteLinkHandler)
	s.POST("/list", app.listLinksHandler)
}
