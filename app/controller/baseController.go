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
func (b *BaseController) Put(c *gin.Context, key string, value interface{}) {
	// lazy init
	if c.Keys[SAVE_DATA_KEY] == nil {
		c.Keys[SAVE_DATA_KEY] = make(map[string]interface{})
	}

	c.Keys[SAVE_DATA_KEY].(map[string]interface{})[key] = value
}

// 正确的响应
func (b *BaseController) RespOK(c *gin.Context, msg ...string) {
	b.resp.Code = constant.RESPONSE_CODE_OK
	b.resp.Msg = "成功"
	b.resp.Data = c.Keys[SAVE_DATA_KEY]

	// 没有数据时，让data字段的json值为[]而非null
	if b.resp.Data == nil {
		b.resp.Data = []int{}
	}

	c.JSON(http.StatusOK, b.resp)
}

// 错误的响应
func (b *BaseController) RespErr(c *gin.Context, options ...interface{}) {
	b.resp.Code = constant.RESPONSE_CODE_ERROR // 默认是常规错误
	b.resp.Msg = ""
	b.resp.Data = c.Keys[SAVE_DATA_KEY]

	// 继续确定code、msg
	for _, v := range options {
		switch opt := v.(type) {
		case int:
			b.resp.Code = opt
		case string:
			b.resp.Msg = opt
		case errors.SysErrorInterface: // 系统错误
			b.resp.Code = constant.RESPONSE_CODE_SYSTEM
			b.resp.Msg = opt.String() // todo opt.Error()
		case errors.NormalErrorInterface: // 常规错误
			if opt.Status() != 0 {
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

	c.JSON(http.StatusOK, b.resp)
}

// 获取get、post提交的参数
func (b *BaseController) GetParam(c *gin.Context, key string) string {
	param := c.Query(key)
	if len(param) == 0 { // get获取不到时，尝试post获取
		param = c.PostForm(key)
	}

	return param
}

// 获取get、post提交的string类型的参数
// def表示默认值，取第一个，多余的丢弃
func (b *BaseController) GetString(c *gin.Context, key string, def ...string) string {
	param := b.GetParam(c, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0]
	}

	return param
}

// 获取get、post提交的int类型的参数
// def表示默认值，取第一个，多余的丢弃
func (b *BaseController) GetInt(c *gin.Context, key string, def ...int) (int, error) {
	param := b.GetParam(c, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.Atoi(param)
}

// 获取get、post提交的int64类型的参数
// def表示默认值，取第一个，多余的丢弃
func (b *BaseController) GetInt64(c *gin.Context, key string, def ...int64) (int64, error) {
	param := b.GetParam(c, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.ParseInt(param, 10, 64)
}

// 获取get、post提交的float64类型的参数
// def表示默认值，取第一个，多余的丢弃
func (b *BaseController) GetFloat64(c *gin.Context, key string, def ...float64) (float64, error) {
	param := b.GetParam(c, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.ParseFloat(param, 64)
}

// 获取get、post提交的float64类型的参数
// def表示默认值，取第一个，多余的丢弃
func (b *BaseController) GetBool(c *gin.Context, key string, def ...bool) (bool, error) {
	param := b.GetParam(c, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.ParseBool(param)
}
