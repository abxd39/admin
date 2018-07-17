package controller

import (
	"admin/app/models"
	"admin/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebUserManageController struct {
	BaseController
}

func (w *WebUserManageController) Router(r *gin.Engine) {
	g := r.Group("/webuser")
	{
		g.GET("/list", w.GetWebUserList)                              //用户管理
		g.GET("/total_user", w.GetTotalUser)                          //获取用户平台注册用户总数
		g.GET("/total_property", w.GetTotalProperty)                  //总资产统计列表
		g.GET("/login_log", w.GetLoginList)                           //用户登录日志
		g.GET("/seconde_certification", w.GetSecodeCertificationList) //获取二级认证列表
		g.POST("/modeify_user_status", w.ModifyUserStatus)            //修改用户状态
		g.POST("/addwhite_list", w.AddWhiteList)                      //增加删除白名单
		g.GET("/user_whitelist", w.WhiteUserList)                     //白名单用户列表
		g.GET("/terminal_list", w.GetTerminalTypeList)                //登录终端类型
		g.GET("/get_second_detail", w.GetSecondDetail)                //二级实名详情
		g.GET("/get_first_datail", w.GetFirstDetail)                  //一级实名认证详情
		g.GET("/get_first_list", w.GetFirstList)                      //p2-4一级实名认证列表
		g.POST("/certification_affirm", w.CertificationAffirm)        //审核用户一级认证
	}
}
func (w *WebUserManageController) CertificationAffirm(c *gin.Context) {
	req := struct {
		Uid int `form:"uid" json:"uid" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	err = new(models.WebUser).CertificationAffirmLimit(req.Uid)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.RespOK(c)
	return
}

func (w *WebUserManageController) GetFirstList(c *gin.Context) {
	req := struct {
		Page    int    `form:"page" json:"page" binding:"required"`
		Rows    int    `form:"rows" json:"rows" `
		Cstatus int    `form:"cstatus" json:"cstatus" ` //认证状态
		Date    uint64 `form:"time" json:"time" `       //日期
		Status  int    `form:"status" json:"status" `   //用户状态
		Search  string `form:"search" json:"search" `   //刷选
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	list, err := new(models.WebUser).GetFirstList(req.Page, req.Rows, req.Status, req.Cstatus, req.Date, req.Search)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.Put(c, "list", list)
	w.RespOK(c)
}

func (w *WebUserManageController) GetFirstDetail(c *gin.Context) {
	req := struct {
		Uid int `form:"uid" json:"uid" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	result, err := new(models.FirstDetail).GetFirstDetail(req.Uid)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.Put(c, "list", result)
	w.RespOK(c)
	return
}

func (w *WebUserManageController) GetSecondDetail(c *gin.Context) {
	req := struct {
		Uid int `form:"uid" json:"uid" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	user, err := new(models.UserSecondaryCertificationGroup).GetSecondaryCertificationOfUid(req.Uid)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.Put(c, "list", user)
	w.RespOK(c)
	return
}

func (w *WebUserManageController) GetTerminalTypeList(c *gin.Context) {
	list, err := new(models.UserLoginTerminalType).GetTerminalTypeList()
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.Put(c, "list", list)
	w.RespOK(c)
	return
}

func (w *WebUserManageController) WhiteUserList(c *gin.Context) {
	req := struct {
		Page   int    `form:"page" json:"page" binding:"required"`
		Rows   int    `form:"rows" json:"rows" `
		Date   uint64 `form:"time" json:"time" `
		Status int    `form:"status" json:"status" `
		Search string `form:"search" json:"search" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	list, err := new(models.WebUser).GetWhiteList(req.Page, req.Rows, req.Status, req.Date, req.Search)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.Put(c, "list", list)
	w.RespOK(c)
	return
}

func (w *WebUserManageController) AddWhiteList(c *gin.Context) {
	req := struct {
		Uid     int `form:"uid" json:"uid" binding:"required"`
		WStatus int `form:"wstatus" json:"wstatus" binding:"required"` //黑白名单状态
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	err = new(models.WebUser).ModifyWhiteStatus(req.Uid, req.WStatus)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.RespOK(c)
	return
}

func (w *WebUserManageController) ModifyUserStatus(c *gin.Context) {
	req := struct {
		Uid    int `form:"uid" json:"uid" binding:"required"`
		Status int `form:"status" json:"status" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	err = new(models.WebUser).ModifyUserStatus(req.Uid, req.Status)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.RespOK(c)
	return
}

func (w *WebUserManageController) GetSecodeCertificationList(c *gin.Context) {
	req := struct {
		Page         int    `form:"page" json:"page" binding:"required"`
		Rows         int    `form:"rows" json:"rows" `
		VerifyStatus int    `form:"verify_status" json:"verify_status" `
		Status       int    `form:"user_status" json:"user_status" `
		VerifyTime   int    `form:"verify_time" json:"verify_time" `
		Search       string `form:"search" json:"search" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	list, err := new(models.UserSecondaryCertificationGroup).GetSecondaryCertificationList(req.Page, req.Rows, req.VerifyStatus, req.Status, req.VerifyTime, req.Search)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.Put(c, "list", list)
	w.RespOK(c)
	return
}

func (w *WebUserManageController) GetLoginList(c *gin.Context) {
	req := struct {
		Page         int    `form:"page" json:"page" binding:"required"`
		Rows         int    `form:"rows" json:"rows" `
		LoginTime    uint64 `form:"login_time" json:"login_time" `
		TerminalType int    `form:"t_type" json:"t_type" `
		Status       int    `form:"status" json:"status" `
		Search       string `form:"search" json:"search" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	list, err := new(models.UserLogInLogGroup).GetUserLoginLogList(req.Page, req.Rows, req.TerminalType, req.Status, req.LoginTime, req.Search)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.Put(c, "list", list)
	w.RespOK(c)
	return
}

////总资产统计列表
func (w *WebUserManageController) GetTotalProperty(c *gin.Context) {
	req := struct {
		Page   int    `form:"page" json:"page" binding:"required"`
		Rows   int    `form:"rows" json:"rows" `
		Status int    `form:"status" json:"status" `
		Search string `form:"search" json:"search" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	//1所有用户数量
	result, err := new(models.WebUser).GetAllUser(req.Page, req.Rows, req.Status, req.Search)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	userlist := make([]int, 0)
	list, Ok := result.Items.([]models.WebUser)
	if !Ok {
		w.RespErr(c, err)
		return
	}
	for _, v := range list {
		userlist = append(userlist, int(v.Uid))
	}
	// //根据UID 获取广告id
	// adlist, err := new(cur.Ads).GetIdList(userlist)
	// if err != nil {
	// 	c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	// 	return
	// }
	orderlist, err := new(models.Order).GetOrderId(userlist, req.Status)
	if err != nil {
		w.RespErr(c, err)
		return
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
	//result.Items =
	w.Put(c, "list", list)
	w.RespOK(c, "成功")
	return
}

func (w *WebUserManageController) GetTotalUser(c *gin.Context) {
	//fmt.Printf("list param %#v\n", req)
	total, err := new(models.WebUser).GetTotalUser()
	if err != nil {
		w.RespErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "", "total": total, "msg": "成功"})
	return
}

func (w *WebUserManageController) GetWebUserList(c *gin.Context) {
	req := struct {
		Page           int    `form:"page" json:"page" binding:"required"`
		Rows           int    `form:"rows" json:"rows" `
		Search_content string `form:"search" json:"search" `
		Date           int64  `form:"date" json:"date" `
		Verify         int    `form:"verify" json:"verify" `
		Status         int    `form:"status" json:"status" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	reuslt, err := new(models.WebUser).UserList(req.Page, req.Rows, req.Verify, req.Status, req.Search_content, req.Date)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	fmt.Println("0.00.0.0.0.0.000000000000000", reuslt)
	list, Ok := reuslt.Items.([]models.UserGroup)
	fmt.Println("GetWebUserList-1")
	if !Ok {
		fmt.Println("GetWebUserList-2")
		w.RespErr(c, err)
		return
	}
	for index, _ := range list {

		if list[index].SecurityAuth == 28 {
			list[index].GoogleVerifyMark = 1
			list[index].RealNameVerifyMark = 1
			list[index].TWOVerifyMark = 1
		} else if list[index].SecurityAuth == 4 {
			list[index].TWOVerifyMark = 1

		} else if list[index].SecurityAuth == 8 {
			list[index].GoogleVerifyMark = 1
		} else if list[index].SecurityAuth == 12 {
			list[index].GoogleVerifyMark = 1
			list[index].TWOVerifyMark = 1
		} else if list[index].SecurityAuth == 16 {
			list[index].RealNameVerifyMark = 1
		} else if list[index].SecurityAuth == 20 {
			list[index].RealNameVerifyMark = 1
			list[index].TWOVerifyMark = 1
		} else if list[index].SecurityAuth == 24 {
			list[index].RealNameVerifyMark = 1
			list[index].GoogleVerifyMark = 1
		}

	}
	if err != nil {
		w.RespErr(c, err)
		return
	}
	// 设置返回数据
	w.Put(c, "list", reuslt)
	fmt.Println("GetWebUserList-3")
	// 返回
	w.RespOK(c)
	//c.JSON(http.StatusOK, gin.H{"code": 0, "page": page, "data": reuslt, "total": total, "msg": "成功"})
	return
}
