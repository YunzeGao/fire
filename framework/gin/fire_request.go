package gin

import (
	"github.com/spf13/cast"
	"mime/multipart"
)

type IRequest interface {
	DefaultQueryInt(key string, def int) (int, bool)
	DefaultQueryInt64(key string, def int64) (int64, bool)
	DefaultQueryFloat(key string, def float32) (float32, bool)
	DefaultQueryDouble(key string, def float64) (float64, bool)
	DefaultQueryBool(key string, def bool) (bool, bool)
	DefaultQueryString(key string, def string) (string, bool)
	DefaultQueryStringSlice(key string, def []string) ([]string, bool)
	DefaultParamInt(key string, def int) (int, bool)
	DefaultParamInt64(key string, def int64) (int64, bool)
	DefaultParamFloat(key string, def float32) (float32, bool)
	DefaultParamDouble(key string, def float64) (float64, bool)
	DefaultParamBool(key string, def bool) (bool, bool)
	DefaultParamString(key string, def string) (string, bool)
	Param(key string) string
	GetParam(key string) (string, bool)
	DefaultFormInt(key string, def int) (int, bool)
	DefaultFormInt64(key string, def int64) (int64, bool)
	DefaultFormFloat(key string, def float32) (float32, bool)
	DefaultFormDouble(key string, def float64) (float64, bool)
	DefaultFormBool(key string, def bool) (bool, bool)
	DefaultFormString(key string, def string) (string, bool)
	DefaultFormStringSlice(key string, def []string) ([]string, bool)
	DefaultFormFile(key string) (*multipart.FileHeader, error)
	DefaultForm(key string, def interface{}) interface{}
	DefaultGetRawData() ([]byte, error)
	DefaultUri() string
	DefaultMethod() string
	DefaultHost() string
	DefaultClientIp() string
	DefaultHeader(key string) string
	DefaultHeaders() map[string][]string
	Cookie(key string) (string, error)
	DefaultCookies() map[string]string
}

// #region context request function

// DefaultQueryAll 获取请求地址中所有参数
func (ctx *Context) DefaultQueryAll() map[string][]string {
	ctx.initQueryCache()
	return ctx.queryCache
}

func (ctx *Context) DefaultQueryInt(key string, def int) (int, bool) {
	if value, ok := ctx.GetQuery(key); ok {
		if v, err := cast.ToIntE(value); err == nil {
			return v, true
		}
	}
	return def, false
}

func (ctx *Context) DefaultQueryInt64(key string, def int64) (int64, bool) {
	if value, ok := ctx.GetQuery(key); ok {
		if v, err := cast.ToInt64E(value); err == nil {
			return v, true
		}
	}

	return def, false
}

func (ctx *Context) DefaultQueryFloat(key string, def float32) (float32, bool) {
	if value, ok := ctx.GetQuery(key); ok {
		if v, err := cast.ToFloat32E(value); err == nil {
			return v, true
		}
	}

	return def, false
}

func (ctx *Context) DefaultQueryDouble(key string, def float64) (float64, bool) {
	if value, ok := ctx.GetQuery(key); ok {
		if v, err := cast.ToFloat64E(value); err == nil {
			return v, true
		}
	}
	return def, false
}

func (ctx *Context) DefaultQueryBool(key string, def bool) (bool, bool) {
	if value, ok := ctx.GetQuery(key); ok {
		if v, err := cast.ToBoolE(value); err == nil {
			return v, true
		}
	}
	return def, false
}

func (ctx *Context) DefaultQueryString(key string, def string) (string, bool) {
	if value, ok := ctx.GetQuery(key); ok {
		return value, true
	}
	return def, false
}

func (ctx *Context) DefaultQueryStringSlice(key string, def []string) ([]string, bool) {
	params := ctx.DefaultQueryAll()
	if values, ok := params[key]; ok {
		return values, true
	}
	return def, false
}

func (ctx *Context) GetParam(key string) (string, bool) {
	if value := ctx.Param(key); value != "" {
		return value, true
	}
	return "", false
}

func (ctx *Context) DefaultParamInt(key string, def int) (int, bool) {
	if value, ok := ctx.GetParam(key); ok {
		if v, err := cast.ToIntE(value); err == nil {
			return v, true
		}
	}
	return def, false
}

func (ctx *Context) DefaultParamInt64(key string, def int64) (int64, bool) {
	if value, ok := ctx.GetParam(key); ok {
		if v, err := cast.ToInt64E(value); err == nil {
			return v, true
		}
	}
	return def, false
}

func (ctx *Context) DefaultParamFloat(key string, def float32) (float32, bool) {
	if value, ok := ctx.GetParam(key); ok {
		if v, err := cast.ToFloat32E(value); err == nil {
			return v, true
		}
	}
	return def, false
}

func (ctx *Context) DefaultParamDouble(key string, def float64) (float64, bool) {
	if value, ok := ctx.GetParam(key); ok {
		if v, err := cast.ToFloat64E(value); err == nil {
			return v, true
		}
	}
	return def, false
}

func (ctx *Context) DefaultParamBool(key string, def bool) (bool, bool) {
	if value, ok := ctx.GetParam(key); ok {
		if v, err := cast.ToBoolE(value); err == nil {
			return v, true
		}
	}
	return def, false
}

func (ctx *Context) DefaultParamString(key string, def string) (string, bool) {
	if value, ok := ctx.GetParam(key); ok {
		return value, true
	}
	return def, false
}

func (ctx *Context) DefaultFormInt(key string, def int) (int, bool) {
	if v, err := cast.ToIntE(ctx.DefaultForm(key, def)); err == nil {
		return v, true
	}
	return def, false
}

func (ctx *Context) DefaultFormInt64(key string, def int64) (int64, bool) {
	if v, err := cast.ToInt64E(ctx.DefaultForm(key, def)); err == nil {
		return v, true
	}
	return def, false
}

func (ctx *Context) DefaultFormFloat(key string, def float32) (float32, bool) {
	if v, err := cast.ToFloat32E(ctx.DefaultForm(key, def)); err == nil {
		return v, true
	}
	return def, false
}

func (ctx *Context) DefaultFormDouble(key string, def float64) (float64, bool) {
	if v, err := cast.ToFloat64E(ctx.DefaultForm(key, def)); err == nil {
		return v, true
	}
	return def, false
}

func (ctx *Context) DefaultFormBool(key string, def bool) (bool, bool) {
	if v, err := cast.ToBoolE(ctx.DefaultForm(key, def)); err == nil {
		return v, true
	}
	return def, false
}

func (ctx *Context) DefaultFormString(key string, def string) (string, bool) {
	if v, err := cast.ToStringE(ctx.DefaultForm(key, def)); err == nil {
		return v, true
	}
	return def, false
}

func (ctx *Context) DefaultFormStringSlice(key string, def []string) ([]string, bool) {
	params := ctx.formCache
	if values, ok := params[key]; ok {
		return values, true
	}
	return def, false
}

func (ctx *Context) DefaultFormFile(key string) (*multipart.FileHeader, error) {
	return ctx.FormFile(key)
}

func (ctx *Context) DefaultForm(key string, def interface{}) interface{} {
	if value, ok := ctx.GetPostForm(key); ok {
		return value
	}
	return def
}

func (ctx *Context) DefaultGetRawData() ([]byte, error) {
	return ctx.GetRawData()
}

func (ctx *Context) DefaultUri() string {
	return ctx.Request.RequestURI
}

func (ctx *Context) DefaultMethod() string {
	return ctx.Request.Method
}

func (ctx *Context) DefaultHost() string {
	return ctx.Request.URL.Host
}

func (ctx *Context) DefaultClientIp() string {
	return ctx.ClientIP()
}

func (ctx *Context) DefaultHeader(key string) string {
	return ctx.GetHeader(key)
}

func (ctx *Context) DefaultHeaders() map[string][]string {
	return ctx.Request.Header
}

func (ctx *Context) DefaultCookie(key string) (string, bool) {
	cookies := ctx.DefaultCookies()
	if val, ok := cookies[key]; ok {
		return val, true
	}
	return "", false
}

func (ctx *Context) DefaultCookies() map[string]string {
	cookies := ctx.Request.Cookies()
	ret := map[string]string{}
	for _, cookie := range cookies {
		ret[cookie.Name] = cookie.Value
	}
	return ret
}

//#endregion
