package main

import (
	"fmt"
	"net/http"

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

func (app *application) IsLogin(r *http.Request) (int, bool) {
	if ok := app.sessionMgr.Exists(r.Context(), authRequiredKey); ok {
		userID := app.sessionMgr.GetInt(r.Context(), authRequiredKey)
		fmt.Println(">>>>>>>> ID: ", userID)
		if userID > 0 {
			return userID, true
		}
	}

	return 0, false
}
