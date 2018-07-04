package controller

import (
	"admin/app/models"
	"admin/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebUserManageController struct{}

func (w *WebUserManageController) Router(r *gin.Engine) {
	group := r.Group("/webuser")
	{
		group.GET("/list", w.GetWebUserList)
	}
}
func (w *WebUserManageController) GetWebUserList(c *gin.Context) {
	req := struct {
		Page   int    `form:"page" json:"page" binding:"required"`
		Rows   int    `form:"rows" json:"rows" `
		Uid    int64  `form:"uid" json:"uid" `
		Uname  string `form:"uname" json:"uname" `
		Phone  string `form:"phone" json:"phone" `
		Email  string `form:"email" json:"email" `
		Date   int64  `form:"date" json:"date" `
		Verify int    `form:"verify" json:"verify" `
		Status int    `form:"status" json:"status" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}
	fmt.Printf("list param %#v\n", req)
	reuslt, total, erro := new(models.WebUser).UserList(req.Page, req.Rows, req.Verify, req.Status, req.Uname, req.Phone, req.Email, req.Date, req.Uid)
	if erro != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": erro.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": reuslt, "total": total, "msg": "成功"})
	return
}
