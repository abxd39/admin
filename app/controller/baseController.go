package controller

import (
	"net/http"
	"strconv"

	"admin/constant"

	"admin/errors"

	"github.com/gin-gonic/gin"
)

const (
	SAVE_DATA_KEY = "save_api_data_key"
)

// 返回的结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// base controller
type BaseController struct {
	resp Response
}

// 设置返回的数据，key-value
// 使用gin context的Keys保存
// gin context每个请求都会先reset
func (b *BaseController) Put(ctx *gin.Context, key string, value interface{}) {
	// lazy init
	if ctx.Keys[SAVE_DATA_KEY] == nil {
		ctx.Keys[SAVE_DATA_KEY] = make(map[string]interface{})
	}

	ctx.Keys[SAVE_DATA_KEY].(map[string]interface{})[key] = value
}

// 正确的响应
func (b *BaseController) RespOK(ctx *gin.Context, msg ...string) {
	b.resp.Code = constant.RESPONSE_CODE_OK
	b.resp.Msg = "成功"
	b.resp.Data = ctx.Keys[SAVE_DATA_KEY]

	// 没有数据时，让data字段的json值为[]而非null
	if b.resp.Data == nil {
		b.resp.Data = []int{}
	}

	ctx.JSON(http.StatusOK, b.resp)
}

// 错误的响应
func (b *BaseController) RespErr(ctx *gin.Context, options ...interface{}) {
	b.resp.Code = constant.RESPONSE_CODE_ERROR // 默认是常规错误
	b.resp.Msg = ""
	b.resp.Data = ctx.Keys[SAVE_DATA_KEY]

	// 继续确定code、msg
	for _, v := range options {
		switch opt := v.(type) {
		case int:
			b.resp.Code = opt // 当前指定code
		case string:
			b.resp.Msg = opt
		case errors.SysErrorInterface: // 系统错误
			b.resp.Code = constant.RESPONSE_CODE_SYSTEM // 设为系统错误code
			b.resp.Msg = opt.String()                   // todo 根据环境使用生产用opt.Error()，本地用opt.String()
		case errors.NormalErrorInterface: // 常规错误
			if opt.Status() != 0 { // 常规错误指定了code并且不为0
				b.resp.Code = opt.Status()
			}
			b.resp.Msg = opt.Error()
		case error: // go错误
			b.resp.Msg = opt.Error()
		}
	}

	// 没有数据时，让data字段的json值为[]而非null
	if b.resp.Data == nil {
		b.resp.Data = []int{}
	}

	ctx.JSON(http.StatusOK, b.resp)
}

// 获取get、post提交的参数
// 参数存在时第二个参数返回true，即使参数的值为空字符串
// 参数不存在时第二个参数返回false
func (b *BaseController) GetParam(ctx *gin.Context, key string) (string, bool) {
	param, ok := ctx.GetQuery(key)
	if len(param) == 0 { // get获取不到时，尝试post获取
		param, ok = ctx.GetPostForm(key)
	}

	return param, ok
}

// 获取get、post提交的string类型的参数
// def表示默认值，取第一个，多余的丢弃
func (b *BaseController) GetString(ctx *gin.Context, key string, def ...string) string {
	param, _ := b.GetParam(ctx, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0]
	}

	return param
}

// 获取get、post提交的int类型的参数
// def表示默认值，取第一个，多余的丢弃
func (b *BaseController) GetInt(ctx *gin.Context, key string, def ...int) (int, error) {
	param, _ := b.GetParam(ctx, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.Atoi(param)
}

// 获取get、post提交的int64类型的参数
// def表示默认值，取第一个，多余的丢弃
func (b *BaseController) GetInt64(ctx *gin.Context, key string, def ...int64) (int64, error) {
	param, _ := b.GetParam(ctx, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.ParseInt(param, 10, 64)
}

// 获取get、post提交的float64类型的参数
// def表示默认值，取第一个，多余的丢弃
func (b *BaseController) GetFloat64(ctx *gin.Context, key string, def ...float64) (float64, error) {
	param, _ := b.GetParam(ctx, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.ParseFloat(param, 64)
}

// 获取get、post提交的float64类型的参数
// def表示默认值，取第一个，多余的丢弃
func (b *BaseController) GetBool(ctx *gin.Context, key string, def ...bool) (bool, error) {
	param, _ := b.GetParam(ctx, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.ParseBool(param)
}
