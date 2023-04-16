package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lightsaid/gotk/form"
	"github.com/lightsaid/gotk/mux"
	"github.com/lightsaid/short-net/models"
	"github.com/lightsaid/short-net/util"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

var (
	renderDataKey   = "rednder_data"
	authRequiredKey = "user_id"
)

type jsonResponse map[string]interface{}

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
		f.Errors.Add("email", "邮箱已被使用")
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

	// 在背后起协程发送邮件
	app.execInBackgorund(func() {
		// 生成 Token
		token, _, err := app.tokenMaker.GenToken(user.ID, app.env.MaxActivateTime)
		if err != nil {
			slog.Error("GenToken failed: "+err.Error(), "email", user.Email, "id", user.ID)
			http.Redirect(w, r, "/servererror", http.StatusSeeOther)
			return
		}

		link := fmt.Sprintf("%s:%d/activate/%s", app.env.HTTPServerHost, app.env.HTTPServerPort, token)

		subject := "激活邮件"
		content := `
			<h1>您好，欢迎注册 ShortNet</h1>
			<p>如果是你本人注册 ShortNet，请点击下面激活账户，若不是请忽略该邮件。</p>
		` + "<p><a href='" + link + "'>激活账户</a></p>"
		to := []string{user.Email}

		err = app.mailer.SendEmail(subject, content, to, nil, nil, nil)
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

// activateHandler 激活用户
func (app *application) activateHandler(w http.ResponseWriter, r *http.Request) {
	token := mux.Param(r, "token")

	payload, err := app.tokenMaker.VerifyToken(token)
	if err != nil {
		slog.Error("activateHandler", "error", err.Error(), "token", token)
		if errors.Is(err, jwt.ErrTokenExpired) {
			app.sessionMgr.Put(r.Context(), "error", "令牌过期，激活失败")
		} else {
			app.sessionMgr.Put(r.Context(), "error", "令牌无效")
		}

		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	if payload.UserID <= 0 {
		app.sessionMgr.Put(r.Context(), "error", "激活失败，用户ID不存在")
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	err = app.store.ActiveUserByID(payload.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		app.sessionMgr.Put(r.Context(), "error", "激活失败，用户不存在")
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	app.sessionMgr.Put(r.Context(), "message", "激活成功")
	http.Redirect(w, r, "/success", http.StatusSeeOther)
}

// loginHandler 登录
func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.sessionMgr.Put(r.Context(), "error", "解析数据失败，未知错误")
		app.renderTemplate(w, r, "sign.page.html", nil)
		return
	}

	f := form.New(r.PostForm)

	// 验证表单数据
	f.Required("email")
	f.Required("password")
	f.IsEmail("email", "请输入合法的邮箱地址")
	f.MinLength("password", 6, "密码长度必须大于或等于6")
	f.MaxLength("password", 24, "密码长度必须小于或等于24")

	data := renderData{
		Form: f,
	}

	// 验证不通过
	if !f.Valid() {
		data.Error = "验证不通过"
		app.sessionMgr.Put(r.Context(), renderDataKey, data)
		http.Redirect(w, r, "/sign", http.StatusSeeOther)
		return
	}

	// 查询用户
	user, err := app.store.GetUserByEmail(f.Get("email"))
	if err != nil {
		data.Error = "查询用户失败"
		if errors.Is(err, gorm.ErrRecordNotFound) {
			f.Errors.Add("email", "用户不存在")
			app.sessionMgr.Put(r.Context(), renderDataKey, data)
			http.Redirect(w, r, "/sign", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/servererror", http.StatusSeeOther)
		return
	}

	err = util.CheckHashedPassword(user.Password, f.Get("password"))
	if err != nil {
		data.Error = "密码不正确"
		f.Errors.Add("email", "邮箱或密码不匹配")
		app.sessionMgr.Put(r.Context(), renderDataKey, data)
		http.Redirect(w, r, "/sign", http.StatusSeeOther)
		return
	}

	if user.Active == 0 {
		f.Errors.Add("email", "用户未激活")
		app.sessionMgr.Put(r.Context(), renderDataKey, data)
		http.Redirect(w, r, "/sign", http.StatusSeeOther)
		return
	}

	// NOTE: 这里 github.com/alexedwards/scs/v2 有个bug，存 uint、int64  数据会丢失
	// fmt.Println(">>>>> user.ID: ", user.ID)
	app.sessionMgr.Put(r.Context(), authRequiredKey, int(user.ID))
	// tmp := app.sessionMgr.GetInt(r.Context(), authRequiredKey)
	// fmt.Println(">>>> tmp: ", tmp)

	app.sessionMgr.Put(r.Context(), "flash", "登录成功")

	http.Redirect(w, r, "/", http.StatusSeeOther)
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

func (app *application) errorHandler(w http.ResponseWriter, r *http.Request) {
	errMsg := app.sessionMgr.PopString(r.Context(), "error")
	if errMsg == "" {
		errMsg = "操作失败"
	}
	app.renderTemplate(w, r, "error.page.html", &renderData{Error: errMsg})
}

func (app *application) forgotHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "forgot.page.html", nil)
}

func (app *application) resetHandler(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "reset.page.html", nil)
}

// ================================ 以下是 Link Table Handler ==============================================

// createLinkHandler 创建 link 返回 json
func (app *application) createLinkHandler(w http.ResponseWriter, r *http.Request) {
	app.mutex.Lock()
	defer app.mutex.Unlock()

	rsp := make(jsonResponse)

	// 解析表单并验证
	r.ParseForm()
	f := form.New(r.PostForm)
	f.Required("long_url")
	if !f.Valid() {
		rsp["error"] = "请输入长网址"
		app.writeJSON(w, r, http.StatusBadRequest, rsp)
		return
	}

	app.shortID++
	shortHashed := util.EncodeBase62(app.shortID)
	userID, isLogin := app.IsLogin(r)
	if !isLogin || userID <= 0 {
		rsp["error"] = "请先登录"
		app.writeJSON(w, r, http.StatusBadRequest, rsp)
		return
	}

	fmt.Println("shortID: ", app.shortID, f.Get("long_url"))
	link := &models.Link{
		UserID:    userID,
		LongURL:   f.Get("long_url"),
		ShortHash: shortHashed,
		ExpiredAt: time.Now().Add(app.env.ShortDefaultExpire),
	}

	err := app.createLink(link)
	if err != nil {
		rsp["error"] = "服务内部错"
		app.writeJSON(w, r, http.StatusInternalServerError, rsp)
		return
	}

	rsp["data"] = link
	app.writeJSON(w, r, http.StatusOK, rsp)
}

func (app *application) updateLinkHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteLinkHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) redirectLinkHandler(w http.ResponseWriter, r *http.Request) {
	hash := mux.Param(r, "hash")
	link, err := app.store.GetLinkByHash(hash)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		r.Context().Done()
		return
	}
	http.Redirect(w, r, link.LongURL, http.StatusSeeOther)
}

func (app *application) listLinksHandler(w http.ResponseWriter, r *http.Request) {

}
