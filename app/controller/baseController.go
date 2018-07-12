package controller

import (
	"admin/constant"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 返回的结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`

	Result map[string]interface{} `json:"-"`
}

// base controller
type BaseController struct {
	resp Response
}

// 设置返回的数据，key-value
func (b *BaseController) Put(key string, value interface{}) {
	if b.resp.Result == nil {
		b.resp.Result = make(map[string]interface{})
	}

	b.resp.Result[key] = value
}

// 正确的响应
func (b *BaseController) RespOK(c *gin.Context, msg string) {
	b.resp.Code = constant.RESPONSE_CODE_OK
	b.resp.Msg = msg

	if b.resp.Result != nil {
		b.resp.Data = b.resp.Result
	} else {
		b.resp.Data = []int{} // 没有数据时，让data字段的json值为[]而非null
	}

	c.JSON(http.StatusOK, b.resp)
}

// 错误的响应
func (b *BaseController) RespErr(c *gin.Context, code int, msg string) {
	b.resp.Code = code
	b.resp.Msg = msg

	if b.resp.Result != nil {
		b.resp.Data = b.resp.Result
	} else {
		b.resp.Data = []int{} // 没有数据时，让data字段的json值为[]而非null
	}

	c.JSON(http.StatusOK, b.resp)
}

// 获取get提交的int类型的参数，如果类型不是int，err不为nil
// def表示默认值
func (b *BaseController) GetInt(c *gin.Context, key string, def ...int) (int, error) {
	param := c.Query(key)
	if len(param) == 0 && len(def) == 1 {
		return def[0], nil
	}

	return strconv.Atoi(param)
}
