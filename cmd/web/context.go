package main

// type contextKey string

/*
	NOTE: 以下代码废弃，本想在重定向前，往context中存数据，重定向后再取，但是此做法不是行不同的，go在重定向后，不会携带之前的context，
	解决方式：
		1. 使用 QueryString，需要将数据序列化成string，也有长度限制，在浏览器看到长长的一串地址，不美观
		2. 使用 Cookie, Cookie 本身也是只能存 string，因此需要将数据序列化，常用的序列化方式有 JSON、XML, 这两种对数据有要求，
		   因此在这个场景中也不适合，唯一的选择只能是 Gob, Gob 是 Go 自带的二进制序列化方式，它只能被 Go 语言所识别，
		   Gob 序列化后的数据更小，因为它是二进制数据，而 JSON 是文本数据。
		   Gob 支持 Go 中的所有数据类型。
*/

// const renderDataCtxkey contextKey = "renderdata"

// // contextSetRenderData 设置渲染的 renderData 值
// func (app *application) contextSetRenderData(r *http.Request, data *renderData) *http.Request {
// 	slog.Info(" renderDataCtxkey >>> ", "value", data.Form, "error", data.Error)
// 	ctx := context.WithValue(r.Context(), renderDataCtxkey, data)
// 	return r.WithContext(ctx)
// }

// // contextGetRenderData 获取 renderData 值
// func (app *application) contextGetRenderData(r *http.Request) *renderData {
// 	data, exists := r.Context().Value(renderDataCtxkey).(*renderData)
// 	if !exists {
// 		slog.Info("构造空 renderData", "url", r.URL.Path, "method", r.Method)
// 		data = app.newRenderData()
// 	}
// 	return data
// }
