package controller

import (
	"net/http"
	"os"
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

// controller
type Controller struct {
	resp Response
}

// 设置返回的数据，key-value
// 使用gin context的Keys保存
// gin context每个请求都会先reset
func (c *Controller) Put(ctx *gin.Context, key string, value interface{}) {
	// lazy init
	if ctx.Keys[SAVE_DATA_KEY] == nil {
		ctx.Keys[SAVE_DATA_KEY] = make(map[string]interface{})
	}

	ctx.Keys[SAVE_DATA_KEY].(map[string]interface{})[key] = value
}

// 正确的响应
func (c *Controller) RespOK(ctx *gin.Context, msg ...string) {
	c.resp.Code = constant.RESPONSE_CODE_OK
	c.resp.Msg = "成功"
	c.resp.Data = ctx.Keys[SAVE_DATA_KEY]

	// 没有数据时，让data字段的json值为[]而非null
	/*if c.resp.Data == nil {
		c.resp.Data = []int{}
	}*/

	ctx.JSON(http.StatusOK, c.resp)
}

// 错误的响应
func (c *Controller) RespErr(ctx *gin.Context, options ...interface{}) {
	c.resp.Code = constant.RESPONSE_CODE_ERROR // 默认是常规错误
	c.resp.Msg = ""
	c.resp.Data = ctx.Keys[SAVE_DATA_KEY]

	// 继续确定code、msg
	for _, v := range options {
		switch opt := v.(type) {
		case int:
			c.resp.Code = opt // 当前指定code
		case string:
		case errors.SysErrorInterface: // 系统错误
			c.resp.Code = constant.RESPONSE_CODE_SYSTEM // 设为系统错误code

			if os.Getenv("API_ENV") == "prod" { // 生产环境不显示错误细节
				c.resp.Msg = opt.String()
			} else { // 开发环境显示错误细节
				c.resp.Msg = opt.Error()
			}
		case errors.NormalErrorInterface: // 常规错误
			if opt.Status() != 0 { // 常规错误指定了code并且不为0
				c.resp.Code = opt.Status()
			}
			c.resp.Msg = opt.Error()
		case error: // go错误
			c.resp.Msg = opt.Error()
		}
	}

	// 优先使用系统指定msg
	c.resp.Msg = constant.GetResponseMsg(c.resp.Code)

	// 没有数据时，让data字段的json值为[]而非null
	/*if c.resp.Data == nil {
		c.resp.Data = []int{}
	}*/

	ctx.JSON(http.StatusOK, c.resp)
}

// 获取get、post提交的参数
// 参数存在时第二个参数返回true，即使参数的值为空字符串
// 参数不存在时第二个参数返回false
func (c *Controller) GetParam(ctx *gin.Context, key string) (string, bool) {
	param, ok := ctx.GetQuery(key)
	if len(param) == 0 { // get获取不到时，尝试post获取
		param, ok = ctx.GetPostForm(key)
	}

	return param, ok
}

// 获取get、post提交的string类型的参数
// def表示默认值，取第一个，多余的丢弃
func (c *Controller) GetString(ctx *gin.Context, key string, def ...string) string {
	param, _ := c.GetParam(ctx, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0]
	}

	return param
}

// 获取get、post提交的int类型的参数
// def表示默认值，取第一个，多余的丢弃
func (c *Controller) GetInt(ctx *gin.Context, key string, def ...int) (int, error) {
	param, _ := c.GetParam(ctx, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.Atoi(param)
}

// 获取get、post提交的int64类型的参数
// def表示默认值，取第一个，多余的丢弃
func (c *Controller) GetInt64(ctx *gin.Context, key string, def ...int64) (int64, error) {
	param, _ := c.GetParam(ctx, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.ParseInt(param, 10, 64)
}

// 获取get、post提交的float64类型的参数
// def表示默认值，取第一个，多余的丢弃
func (c *Controller) GetFloat64(ctx *gin.Context, key string, def ...float64) (float64, error) {
	param, _ := c.GetParam(ctx, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.ParseFloat(param, 64)
}

// 获取get、post提交的float64类型的参数
// def表示默认值，取第一个，多余的丢弃
func (c *Controller) GetBool(ctx *gin.Context, key string, def ...bool) (bool, error) {
	param, _ := c.GetParam(ctx, key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.ParseBool(param)
}
