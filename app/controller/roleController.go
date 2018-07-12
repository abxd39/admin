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
	page, err := this.GetInt(c, "page", 1)
	if err != nil {
		this.RespErr(c, constant.RESPONSE_CODE_ERROR, "参数page格式错误")
		return
	}

	rows, err := this.GetInt(c, "rows", 10)
	if err != nil {
		this.RespErr(c, constant.RESPONSE_CODE_ERROR, "参数rows格式错误")
		return
	}

	// 调用model
	list, err := new(backstage.Role).List(page, rows)
	if err != nil {
		this.RespErr(c, constant.RESPONSE_CODE_ERROR, "查询失败")
		return
	}

	// 设置返回数据
	this.Put("list", list)

	// 返回
	this.RespOK(c, "成功")
	return
}
