package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"golang.org/x/exp/slog"
)

func (app *application) genTemplateCache() error {
	slog.Info("templates base path", "path", app.env.ViewPath+"/*.page.html")
	matches, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", app.env.ViewPath))
	if err != nil {
		slog.Error("match")
		return err
	}
	slog.Info("matches views", "views", matches)
	for _, match := range matches {
		t, err := template.ParseFiles(match)
		if err != nil {
			slog.Error("template.ParseFiles: "+err.Error(), "match", match)
			continue
		}

		filename := filepath.Base(match)
		if app.templateCache == nil {
			app.templateCache = make(map[string]*template.Template, len(matches))
		}
		app.templateCache[filename] = t
	}
	return nil
}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, tplname string) {
	t, exists := app.templateCache[tplname]
	if !exists {
		slog.Error("template not found", "tplname", tplname)
		// 错误处理
		return
	}

	err := t.Execute(w, nil)
	if err != nil {
		slog.Error("render template 'Execute' error: "+err.Error(), "tplname", tplname)
	}
}
