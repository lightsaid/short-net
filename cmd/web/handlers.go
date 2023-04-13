package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/lightsaid/gotk/form"
	"github.com/lightsaid/short-net/models"
	"github.com/lightsaid/short-net/util"
	"golang.org/x/exp/slog"
)

var (
	renderDataKey = "rednder_data"
)

func (app *application) indexHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "index.page.html", nil)
}

func (app *application) signHandler(w http.ResponseWriter, r *http.Request) {
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
	app.renderTemplate(w, r, "sign.page.html", &data)
}

// registerHandler 注册
func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.renderTemplate(w, r, "sign.page.html", &renderData{Error: "解析表单数据失败"})
		return
	}

	f := form.New(r.PostForm)

	// 验证表单数据
	f.Required("name")
	f.Required("email")
	f.Required("password")
	f.Required("repassword")
	f.MinLength("name", 2, "用户名长度必须大于或等于2")
	f.MaxLength("name", 24, "用户名长度必须小于或等于24")
	f.IsEmail("email", "请输入合法的邮箱地址")
	f.MinLength("password", 6, "密码长度必须大于或等于6")
	f.MaxLength("password", 24, "密码长度必须小于或等于24")
	if f.Get("password") != f.Get("repassword") {
		f.Errors.Add("repassword", "两次密码不一致")
	}

	// 验证不通过
	if !f.Valid() {
		data := renderData{
			Error: "验证不通过",
			Form:  f,
		}
		app.sessionMgr.Put(r.Context(), renderDataKey, data)
		log.Println("设置了renderDataKey data ", "error", data.Error, "password_error", data.Form.Errors.Get("password"))
		http.Redirect(w, r, "/sign?t=1", http.StatusSeeOther)
		return
	}

	// 验证通过，注册流程，注册执行事务，等发送邮件成功，在提交事物
	hashedPwsd, err := util.GenHashedPassword(f.Get("password"))
	if err != nil {
		http.Redirect(w, r, "/servererror", http.StatusSeeOther)
		return
	}
	user := models.User{
		Name:     f.Get("name"),
		Email:    f.Get("email"),
		Password: hashedPwsd,
		Avatar:   fmt.Sprintf("%s%d%s", "/static/images/avatar", util.RandomInt(1, 2), ".png"),
	}

	qUser, _ := app.store.GetUserByEmail(user.Email)
	if qUser.ID > 0 {
		data := renderData{
			Error: "邮箱已被使用",
			Form:  f,
		}
		app.sessionMgr.Put(r.Context(), renderDataKey, data)
		http.Redirect(w, r, "/sign?t=1", http.StatusSeeOther)
		return
	}

	err = app.store.CreateUser(&user)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			data := renderData{
				Error: "邮箱已被使用",
				Form:  f,
			}
			app.sessionMgr.Put(r.Context(), renderDataKey, data)
			http.Redirect(w, r, "/sign?t=1", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/servererror", http.StatusSeeOther)
		return
	}

	// 发送邮件
	app.execInBackgorund(func() {
		subject := "激活邮件"
		content := `
			<h1>您好，欢迎注册 ShortNet</h1>
			<p>如果是你本人注册 ShortNet，请点击下面激活账户，若不是请忽略该邮件。</p>
			<p><a href="https://localhost:4000">激活账户</a></p>
		`
		to := []string{user.Email}

		err := app.mailer.SendEmail(subject, content, to, nil, nil, nil)
		if err != nil {
			slog.Error("sender register error: "+err.Error(), "email", user.Email)
		}
	})

	app.sessionMgr.Put(r.Context(), "message", "恭喜你注册成功，请到邮箱激活用户～")
	http.Redirect(w, r, "/success", http.StatusSeeOther)

	// NOTE: 这里只是一种尝试
	// txErr := app.store.TxRegister(&user, func(err error) {
	// 	if err != nil {
	// 		slog.Error("TxRegister: "+err.Error(), "email", user.Email)

	// 		var mysqlErr *mysql.MySQLError
	// 		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
	// 			data := renderData{
	// 				Error: "邮箱已被使用",
	// 				Form:  f,
	// 			}
	// 			app.sessionMgr.Put(r.Context(), renderDataKey, data)
	// 			http.Redirect(w, r, "/sign?t=1", http.StatusSeeOther)

	// 			return
	// 		} else {
	// 			http.Redirect(w, r, "/servererror", http.StatusSeeOther)
	// 			return
	// 		}
	// 	}

	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// })

	// if txErr != nil {
	// 	slog.Error("TxRegister final: ", txErr, "email", user.Email)
	// }
}

// loginHandler 登录
func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// TODO：错误提示
		app.renderTemplate(w, r, "sign.page.html", nil)
		return
	}

	// f := form.New(r.PostForm)
	// f.Required("")
	app.renderTemplate(w, r, "sign.page.html", nil)
}

func (app *application) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "404.page.html", nil)
}

func (app *application) serverErrorHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "500.page.html", nil)
}

func (app *application) operateSuccessfully(w http.ResponseWriter, r *http.Request) {
	message := app.sessionMgr.PopString(r.Context(), "message")
	if message == "" {
		slog.Info("message 丢失？")
		message = "操作成功"
	}
	data := renderData{
		StringMap: map[string]string{
			"message": message,
		},
	}
	app.renderTemplate(w, r, "success.page.html", &data)
}

func (app *application) forgotHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "forgot.page.html", nil)
}

func (app *application) resetHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "reset.page.html", nil)
}
