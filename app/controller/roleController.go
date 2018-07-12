package controller

import (
	"unicode/utf8"

	"admin/app/models/backstage"
	"admin/constant"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	BaseController
}

func (r *RoleController) Router(e *gin.Engine) {
	group := e.Group("/role")
	{
		group.GET("/list", r.List)
		group.POST("/add", r.Add)
	}
}

// 用户组列表
func (r *RoleController) List(c *gin.Context) {
	// 获取参数
	page, err := r.GetInt(c, "page", 1)
	if err != nil {
		r.RespErr(c, constant.RESPONSE_CODE_ERROR, "参数page格式错误")
		return
	}

	rows, err := r.GetInt(c, "rows", 10)
	if err != nil {
		r.RespErr(c, constant.RESPONSE_CODE_ERROR, "参数rows格式错误")
		return
	}

	// 调用model
	list, err := new(backstage.Role).List(page, rows)
	if err != nil {
		r.RespErr(c, constant.RESPONSE_CODE_ERROR, "查询失败")
		return
	}

	// 设置返回数据
	r.Put(c, "list", list)

	// 返回
	r.RespOK(c, "查询成功")
	return
}

// 新增用户组
func (r *RoleController) Add(c *gin.Context) {
	// 获取参数
	name := r.GetString(c, "name")
	if strLen := utf8.RuneCountInString(name); strLen == 0 || strLen > 10 {
		r.RespErr(c, constant.RESPONSE_CODE_ERROR, "参数name格式错误")
		return
	}

	desc := r.GetString(c, "desc", "")
	nodeIds := r.GetString(c, "node_ids", "")

	// 调用model
	id, err := new(backstage.Role).Add(name, desc, nodeIds)
	if err != nil {
		r.RespErr(c, constant.RESPONSE_CODE_ERROR, "新增失败，请稍后重试")
		return
	}

	// 设置返回数据
	r.Put(c, "id", id)

	// 返回
	r.RespOK(c, "新增成功")
	return
}
