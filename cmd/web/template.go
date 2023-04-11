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
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", app.env.ViewPath))
	if err != nil {
		slog.Error("match")
		return err
	}
	slog.Info("matches views", "views", pages)
	for _, pages := range pages {
		// 解析模版
		t, err := template.ParseFiles(pages)
		if err != nil {
			slog.Error("template.ParseFiles: "+err.Error(), "pages", pages)
			return err
		}

		// 加载布局组件 *.layout.html
		t, err = t.ParseGlob(fmt.Sprintf("%s/*.layout.html", app.env.ViewPath))
		if err != nil {
			return err
		}

		// 加载其他组件 *.partial.html
		t, err = t.ParseGlob(fmt.Sprintf("%s/*.partial.html", app.env.ViewPath))
		if err != nil {
			return err
		}

		// 根据文件名缓存，如 index.page.html
		filename := filepath.Base(pages)
		if app.templateCache == nil {
			app.templateCache = make(map[string]*template.Template, len(pages))
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
