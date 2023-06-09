package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/lightsaid/gotk/form"
	"github.com/lightsaid/short-net/models"

	"golang.org/x/exp/slog"
)

// renderData 渲染模版需要用的数据
type renderData struct {
	Form      *form.Form        // 表单数据，如在表单提供验证不通过时，通过自此字段返回错误信息
	Flash     string            // 操作成功通过
	Error     string            // 操作错误通知
	Warning   string            // 警告通知
	Info      string            // 普通通知
	StringMap map[string]string // string map
	Data      map[string]any
	CSRFToken string // csrf token
	IsLogin   int    // 是否已经登录
}

func (app *application) newRenderData() *renderData {
	return &renderData{
		Form: form.New(nil),
	}
}

// genTemplateCache 生成模版缓存
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

// addDefaultData 添加一些默认数据
func (app *application) addDefaultData(r *http.Request, data *renderData) *renderData {
	if data == nil {
		data = app.newRenderData()
	}

	data.CSRFToken = nosurf.Token(r)

	data.Error = app.sessionMgr.PopString(r.Context(), "error")
	data.Flash = app.sessionMgr.PopString(r.Context(), "flash")
	data.Warning = app.sessionMgr.PopString(r.Context(), "warning")
	data.Info = app.sessionMgr.PopString(r.Context(), "info")

	userInfo, ok := app.sessionMgr.Get(r.Context(), UserInfoKey).(models.User)
	if ok && userInfo.ID > 0 {
		if data.Data == nil {
			data.Data = make(map[string]any)
		}
		data.Data["userInfo"] = userInfo
	}

	// 是否登录
	if _, ok := app.IsLogin(r); ok {
		data.IsLogin = 1
	}

	return data
}

// renderTemplate 根据 tplname 渲染模板
func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, tplname string, data *renderData) {
	t, exists := app.templateCache[tplname]
	if !exists {
		slog.Error("template not found", "tplname", tplname)

		// 错误处理
		http.Redirect(w, r, "/notfound", http.StatusSeeOther)
		return
	}

	data = app.addDefaultData(r, data)

	buf := new(bytes.Buffer)
	err := t.Execute(buf, data)
	if err != nil {
		slog.Error("执行模版错误: "+err.Error(), "tplname", tplname, "data", data)
		http.Redirect(w, r, "/servererror", http.StatusSeeOther)
		return
	}

	// err = t.Execute(w, data)
	_, err = buf.WriteTo(w)
	if err != nil {
		slog.Error("render template 'Execute' error: "+err.Error(), "tplname", tplname)
	}
}
