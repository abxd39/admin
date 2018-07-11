package controller

import (
	"admin/app/models/backstage"
	"admin/constant"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	BaseController
}

func (this *RoleController) Router(r *gin.Engine) {
	group := r.Group("/role")
	{
		group.GET("/list", this.List)
	}
}

// 用户组列表
func (this *RoleController) List(c *gin.Context) {
	// 获取参数
	req := struct {
		Page int `form:"page" json:"page" binding:"required"`
		Rows int `form:"rows" json:"rows"`
	}{}
	err := c.Bind(&req)
	if err != nil {
		this.RespErr(c, constant.RESPONSE_CODE_ERROR, "参数错误")
		return
	}

	// 调用model
	list, err := new(backstage.Role).List(req.Page, req.Rows)
	if err != nil {
		this.RespErr(c, constant.RESPONSE_CODE_SYSTEM, "系统错误")
		return
	}

	// 设置返回数据
	this.Put("list", list)

	// 返回
	this.RespOK(c, "成功")
	return
}
