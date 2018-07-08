package controller

import (
	cur "admin/app/models/currency"
	models "admin/app/models/user"
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
		group.GET("/total_user", w.GetTotalUser)         //获取用户平台注册用户总数
		group.GET("/total_property", w.GetTotalProperty) //总资产统计列表
	}
}

func (w *WebUserManageController) GetTotalProperty(c *gin.Context) {
	req := struct {
		Page   int `form:"page" json:"page" binding:"required"`
		Rows   int `form:"rows" json:"rows" `
		Status int `form:"status" json:"status" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}
	//1所有用户数量
	result, page, total, err := new(models.WebUser).GetAllUser(req.Page, req.Rows, req.Status)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	userlist := make([]int, 0)
	for _, v := range result {
		userlist = append(userlist, int(v.Uid))
	}
	// //根据UID 获取广告id
	// adlist, err := new(cur.Ads).GetIdList(userlist)
	// if err != nil {
	// 	c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	// 	return
	// }
	orderlist, err := new(cur.Order).GetOrderId(userlist, req.Status)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	type Rsp struct {
		Uid           int
		NickName      string
		Phone         string
		Email         string
		Status        int
		TotalBalance  float32
		TotalToken    float32
		TotalCurrency float32
	}
	_ = orderlist
	// rsp := make([]Rsp, 0)
	// for _, v := range orderlist {

	// }
	//统计所有订单id与UID 的管理
	//var ad map[int][]int
	//orderId := make([]string, 0)

	//for _, uid := range userlist {
	// 	//找出相同uid的所有id
	// 	idlist := make([]int, 0)
	// 	for _, adid := range adlist {
	// 		if uid == int(adid.Uid) {
	// 			idlist = append(idlist, int(adid.Id))
	// 		}
	// 	}

	// }

	//2法币总资产
	//再此需要传递一个交易的状态作为参数 付状态: 1待支付 2待放行(已支付) 3确认支付(已完成)
	//new(cur.Order).GetOrderId(orderId, 3)
	//3币币总资产
	c.JSON(http.StatusOK, gin.H{"code": 0, "page": page, "total": total, "data": result, "msg": "成功"})
	return
}

func (w *WebUserManageController) GetTotalUser(c *gin.Context) {
	//fmt.Printf("list param %#v\n", req)
	total, err := new(models.WebUser).GetTotalUser()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "", "total": total, "msg": "成功"})
	return
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
	reuslt, page, total, err := new(models.WebUser).UserList(req.Page, req.Rows, req.Verify, req.Status, req.Uname, req.Phone, req.Email, req.Date, req.Uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "page": page, "data": reuslt, "total": total, "msg": "成功"})
	return
}
