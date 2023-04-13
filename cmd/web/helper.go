package main

import (
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
