package gin

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
)

type IResponse interface {
	IJson(obj interface{}) IResponse
	IJsonp(obj interface{}) IResponse
	IXml(obj interface{}) IResponse
	IHtml(file string, obj interface{}) IResponse
	IText(format string, values ...interface{}) IResponse
	IRedirect(path string) IResponse
	ISetHeader(key string, val string) IResponse
	ISetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse
	ISetStatus(code int) IResponse
	ISetOkStatus() IResponse // 设置200状态
}

var _ IResponse = new(Context)

// #region context response function

func (ctx *Context) IJson(obj interface{}) IResponse {
	byt, err := json.Marshal(obj)
	if err != nil {
		return ctx.ISetStatus(http.StatusInternalServerError)
	}
	ctx.ISetHeader("Content-Type", "application/json")
	_, _ = ctx.Writer.Write(byt)
	return ctx
}

func (ctx *Context) IJsonp(obj interface{}) IResponse {
	// 获取请求参数callback
	callbackFunc := ctx.Query("callback")
	ctx.ISetHeader("Content-Type", "application/javascript")
	// 输出到前端页面的时候需要注意下进行字符过滤，否则有可能造成xss攻击
	callback := template.JSEscapeString(callbackFunc)
	// 输出函数名
	_, err := ctx.Writer.Write([]byte(callback))
	if err != nil {
		return ctx
	}
	// 输出左括号
	_, err = ctx.Writer.Write([]byte("("))
	if err != nil {
		return ctx
	}
	// 数据函数参数
	ret, err := json.Marshal(obj)
	if err != nil {
		return ctx
	}
	_, err = ctx.Writer.Write(ret)
	if err != nil {
		return ctx
	}
	// 输出右括号
	_, err = ctx.Writer.Write([]byte(")"))
	if err != nil {
		return ctx
	}
	return ctx
}

func (ctx *Context) IXml(obj interface{}) IResponse {
	byt, err := xml.Marshal(obj)
	if err != nil {
		return ctx.ISetStatus(http.StatusInternalServerError)
	}
	ctx.ISetHeader("Content-Type", "application/html")
	_, _ = ctx.Writer.Write(byt)
	return ctx
}

func (ctx *Context) IHtml(file string, obj interface{}) IResponse {
	// 读取模版文件，创建template实例
	t, err := template.New("output").ParseFiles(file)
	if err != nil {
		return ctx
	}
	// 执行Execute方法将obj和模版进行结合
	if err := t.Execute(ctx.Writer, obj); err != nil {
		return ctx
	}

	ctx.ISetHeader("Content-Type", "application/html")
	return ctx
}

func (ctx *Context) IText(format string, values ...interface{}) IResponse {
	out := fmt.Sprintf(format, values...)
	ctx.ISetHeader("Content-Type", "application/text")
	_, _ = ctx.Writer.Write([]byte(out))
	return ctx
}

func (ctx *Context) IRedirect(path string) IResponse {
	http.Redirect(ctx.Writer, ctx.Request, path, http.StatusMovedPermanently)
	return ctx
}

func (ctx *Context) ISetHeader(key string, val string) IResponse {
	ctx.Header(key, val)
	return ctx
}

func (ctx *Context) ISetCookie(key string, val string, maxAge int, path, domain string, secure, httpOnly bool) IResponse {
	ctx.SetCookie(key, val, maxAge, path, domain, secure, httpOnly)
	return ctx
}

func (ctx *Context) ISetStatus(code int) IResponse {
	ctx.Status(code)
	return ctx
}

func (ctx *Context) ISetOkStatus() IResponse {
	return ctx.ISetStatus(http.StatusOK)
}

//#endregion
