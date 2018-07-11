package controller

import (
	"admin/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 返回的结构
type Response struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

// base controller
type BaseController struct {
	resp Response
}

// 设置返回的数据，key-value
func (b *BaseController) Put(key string, value interface{}) {
	if b.resp.Data == nil {
		b.resp.Data = make(map[string]interface{})
	}

	b.resp.Data[key] = value
}

// 正确的响应
func (b *BaseController) RespOK(c *gin.Context, msg string) {
	b.resp.Code = constant.RESPONSE_CODE_OK
	b.resp.Msg = msg
	c.JSON(http.StatusOK, b.resp)
}

// 错误的响应
func (b *BaseController) RespErr(c *gin.Context, code int, msg string) {
	b.resp.Code = code
	b.resp.Msg = msg

	c.JSON(http.StatusOK, b.resp)
}
