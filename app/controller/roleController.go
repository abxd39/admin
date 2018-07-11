package controller

import (
	"net/http"

	"admin/utils"

	"admin/app/models/backstage"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
}

func (this *RoleController) Router(r *gin.Engine) {
	group := r.Group("/role")
	{
		group.GET("/list", this.GetRoleList)
	}
}

// 获取用户组列表
func (this *RoleController) GetRoleList(c *gin.Context) {
	// 获取参数
	req := struct {
		Page int `form:"page" json:"page" binding:"required"`
		Rows int `form:"rows" json:"rows"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
	}

	// 调用model
	list, totalPage, total, err := new(backstage.Role).GetRoleList(req.Page, req.Rows)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "page": totalPage, "total": total, "data": list, "msg": "成功"})
}
