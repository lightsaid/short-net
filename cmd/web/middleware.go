package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/justinas/nosurf"
	"golang.org/x/exp/slog"
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

func (app *application) csrfMiddleware(next http.Handler) http.Handler {
	var secure bool
	if app.env.RunMode == "prod" {
		secure = true
	}

	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// recovererMiddleware 恐慌恢复
func (app *application) recovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil && r != http.ErrAbortHandler {
				slog.Error("PANIC", "recover", r)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// authRequired 身份认证
func (app *application) authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, isLogin := app.IsLogin(r); !isLogin {

			contentType := strings.ToLower(r.Header.Get("Content-Type"))

			// 判断请求是否由 fetch 发送
			if strings.Contains(contentType, "multipart/form-data") || strings.Contains(contentType, "application/json") {
				resp := make(jsonResponse)
				resp["error"] = "请先登录"
				resp["redirect"] = "/sign"
				app.writeJSON(w, r, http.StatusSeeOther, resp)
				return
			}

			app.sessionMgr.Put(r.Context(), "error", "请先登录")
			http.Redirect(w, r, "/sign", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
