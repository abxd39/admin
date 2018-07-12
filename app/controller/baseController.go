package controller

import (
	"net/http"
	"strconv"

	"admin/constant"

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
func (b *BaseController) Put(c *gin.Context, key string, value interface{}) {
	// 使用gin context的Keys保存临时数据，保证每个请求之前都能reset
	c.Keys[SAVE_DATA_KEY] = map[string]interface{}{key: value}
}

// 正确的响应
func (b *BaseController) RespOK(c *gin.Context, msg string) {
	b.resp.Code = constant.RESPONSE_CODE_OK
	b.resp.Msg = msg

	if c.Keys[SAVE_DATA_KEY] != nil {
		b.resp.Data = c.Keys[SAVE_DATA_KEY]
	} else {
		b.resp.Data = []int{} // 没有数据时，让data字段的json值为[]而非null
	}

	c.JSON(http.StatusOK, b.resp)
}

// 错误的响应
func (b *BaseController) RespErr(c *gin.Context, code int, msg string) {
	b.resp.Code = code
	b.resp.Msg = msg

	if c.Keys[SAVE_DATA_KEY] != nil {
		b.resp.Data = c.Keys[SAVE_DATA_KEY]
	} else {
		b.resp.Data = []int{} // 没有数据时，让data字段的json值为[]而非null
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
