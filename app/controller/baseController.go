package controller

import (
	"admin/constant"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

// 获取get提交的int类型的参数
// 如果类型不是int，err不为nil
// def表示默认值，取第一个，多余的丢弃
func (b *BaseController) GetInt(c *gin.Context, key string, def ...int) (int, error) {
	param := c.Query(key)
	if len(param) == 0 && len(def) > 0 {
		return def[0], nil
	}

	return strconv.Atoi(param)
}
