package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/lightsaid/gotk/form"
	"github.com/lightsaid/gotk/mux"
	"github.com/lightsaid/short-net/dbrepo"
	"github.com/lightsaid/short-net/models"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
	data := renderData{
		Data: map[string]any{},
	}
	defer app.renderTemplate(w, r, "book.page.html", &data)

	books, err := app.store.ListBooks(dbrepo.Filters{Page: 1, Size: 100})
	if err != nil {
		data.Error = "查询失败"
		return
	}

	if len(books) == 0 {
		data.Info = "暂无数据"
		return
	}

	data.Data["books"] = books
}

func (app *application) showOrderHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "order.page.html", app.newRenderData())
}

func (app *application) showCreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var data = *app.newRenderData()
	// 注册或者登录重定向过来
	if exists := app.sessionMgr.Exists(r.Context(), renderDataKey); exists {
		var ok bool
		data, ok = app.sessionMgr.Get(r.Context(), renderDataKey).(renderData)
		if !ok {
			slog.Error("signHandler", "renderDataKey", nil)
		}
		app.sessionMgr.Remove(r.Context(), renderDataKey)
	}
	app.renderTemplate(w, r, "createBook.page.html", &data)
}

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	f := form.New(r.PostForm)

	data := renderData{}

	f.Required("title", "price", "stock")
	title := f.Get("title")
	price, err1 := strconv.Atoi(f.Get("price"))
	stock, err2 := strconv.Atoi(f.Get("stock"))
	if err1 != nil || err2 != nil {
		if err1 != nil {
			f.Errors.Add("price", "价格必须是整数")
		}
		if err2 != nil {
			f.Errors.Add("stock", "库存必须是整数")
		}
		data.Form = f
		app.sessionMgr.Put(r.Context(), renderDataKey, data)

		http.Redirect(w, r, "/book/create", http.StatusSeeOther)
		return
	}

	filename, err := app.formUpload(w, r)
	if err != nil {
		slog.Error("create book upload error: " + err.Error())
		if errors.Is(err, http.ErrMissingFile) {
			f.Errors.Add("file", "请上传图书图片")
		} else {
			f.Errors.Add("file", "上传图片错误，文件大小不要超过 2 M")
		}
	}

	if !f.Valid() {
		data.Form = f
		app.sessionMgr.Put(r.Context(), renderDataKey, data)
		http.Redirect(w, r, "/book/create", http.StatusSeeOther)
		return
	}

	book := models.Book{
		Title:   title,
		Price:   uint(price),
		Stcok:   uint(stock),
		Picture: filename,
	}

	err = app.store.CreateBook(&book)
	if err != nil {
		slog.Error("createBook error: ", err.Error())
		data.Form = f
		app.sessionMgr.Put(r.Context(), "error", "服务内部错误")
		app.sessionMgr.Put(r.Context(), renderDataKey, data)
		http.Redirect(w, r, "/book/create", http.StatusSeeOther)
		return
	}
	app.sessionMgr.Put(r.Context(), "flash", "创建成功")
	http.Redirect(w, r, "/book/create", http.StatusSeeOther)
}

func (app *application) buyBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Param(r, "id"))
	if err != nil {
		app.sessionMgr.Put(r.Context(), "error", "发生异常, Id 不合法")
		http.Redirect(w, r, "/book/index", http.StatusSeeOther)
		return
	}

	userID, ok := app.IsLogin(r)
	if !ok || userID <= 0 {
		app.sessionMgr.Put(r.Context(), "error", "请先登录")
		http.Redirect(w, r, "/sign", http.StatusSeeOther)
		return
	}

	book, err := app.store.GetBook(uint(id))
	if err != nil {
		slog.Error("buy book error: "+err.Error(), "userId", userID, "bookid", book.ID)
		var msg = "服务错误"
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg = "book 不存在"
		}
		app.sessionMgr.Put(r.Context(), "error", msg)
		http.Redirect(w, r, "/book/index", http.StatusSeeOther)
		return
	}

	// NOTE: 目前一次仅仅对一本书下单
	order := models.Order{
		UserID:      userID,
		TotalAmount: book.Price,
		OrderDetails: []models.OrderDetail{
			{Qty: 1, Amount: book.Price, BookID: book.ID},
		},
	}
	// TODO: 扣减库存
	err = app.store.CreateOrder(&order)
	if err != nil {
		slog.Error("buy book CreateOrder error: "+err.Error(), "userId", userID, "bookid", book.ID)
		app.sessionMgr.Put(r.Context(), "error", "抢够失败")
		http.Redirect(w, r, "/book/index", http.StatusSeeOther)
		return
	}
	app.sessionMgr.Put(r.Context(), "flash", "抢够成功")
	http.Redirect(w, r, "/book/index", http.StatusSeeOther)
}
