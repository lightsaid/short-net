package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/lightsaid/short-net/models"
	"github.com/lightsaid/short-net/util"
	"golang.org/x/exp/slog"
)

// execInBackgorund 在背后执行一个异步操作
func (app *application) execInBackgorund(fn func()) {
	app.wg.Add(1)

	go func() {
		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				slog.Error("execInBackgorund panic: ", "error", err)
			}
		}()

		fn()
	}()
}

func (app *application) IsLogin(r *http.Request) (uint, bool) {
	if ok := app.sessionMgr.Exists(r.Context(), authRequiredKey); ok {
		userID := app.sessionMgr.GetInt(r.Context(), authRequiredKey)
		if userID > 0 {
			return uint(userID), true
		}
	}

	return 0, false
}

func (app *application) createLink(link *models.Link) error {
	var times = 0
	err := app.store.CreateLink(link)
	if err != nil {
		// 检查是否是唯一约束错误
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			times++
			if times >= 10 {
				// TODO: 发送邮件通知管理员
				slog.Error(
					"创建 link hash 唯一约束错误次数过多，快来看看吧",
					"error",
					err.Error(),
					"userId",
					link.UserID,
					"ShortHash",
					link.ShortHash,
				)

				return errors.New("ShortHash 重复严重")
			}
			app.shortID++
			link.ShortHash = util.EncodeBase62(app.shortID)
			app.createLink(link)
		} else {
			slog.Error(
				"创建 Link 错误",
				"error", err.Error(),
				"userID", link.UserID,
				"shortID", app.shortID,
				"hashe", link.ShortHash,
			)
			return err
		}
	}
	return nil
}

func (app *application) writeJSON(w http.ResponseWriter, r *http.Request, status int, data jsonResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	data["status"] = status
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Redirect(w, r, "/servererror", http.StatusSeeOther)
	}
}
