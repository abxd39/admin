package controller

import (
	"admin/app/models"
	"admin/utils"
	"fmt"

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
		g.GET("/seconde_certification", w.GetSecondCertificationList) //获取二级认证列表
		g.POST("/modeify_user_status", w.ModifyUserStatus)            //修改用户状态
		g.POST("/addwhite_list", w.AddWhiteList)                      //增加删除白名单
		g.GET("/user_whitelist", w.WhiteUserList)                     //白名单用户列表
		g.GET("/terminal_list", w.GetTerminalTypeList)                //登录终端类型
		g.GET("/get_second_detail", w.GetSecondDetail)                //二级实名详情
		g.GET("/get_first_datail", w.GetFirstDetail)                  //一级实名认证详情
		g.GET("/get_first_list", w.GetFirstList)                      //p2-4一级实名认证列表
		g.POST("/first_affirm", w.FirstAffirm)                        //审核用户实名认证
		g.POST("/second_affirm", w.SecondAffirm)                      //审核二级实名认证
		g.POST("/trade_rule", w.SetTradeRule)                         //设置交易规则
		g.GET("/get_trade_rule", w.GetTradeRule)                      //获取交易规则
		g.GET("/get_invite_list", w.GetInviteList)                    //获取 p2-5好友邀列表 被邀请人列表 邀请人—账号：18888888888
		g.GET("/get_invite_info", w.GetInviteInfoList)                //p2-5-1邀请人统计列表

		//用户系统设置
		g.POST("/token_system_add", w.TokenSystemAdd)
		g.POST("/delete_system", w.DeleteSystem)
		g.GET("/get_system", w.GetSystem)
		g.GET("/get_system_list", w.GetSystemList)
	}
}

func (w *WebUserManageController) DeleteSystem(c *gin.Context) {
	req := struct {
		Id int `form:"id" json:"id" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}

	err = new(models.Tokens).DeleteSystem(req.Id)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.RespOK(c)
}

// 获取单条记录
func (w *WebUserManageController) GetSystem(c *gin.Context) {
	req := struct {
		Id int `form:"id" json:"id" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	result, err := new(models.Tokens).GetSystem(req.Id)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.Put(c, "data", result)
	w.RespOK(c)
}

//系统设置 获取列表
func (w *WebUserManageController) GetSystemList(c *gin.Context) {
	req := struct {
		Page   int    `form:"page" json:"page" binding:"required"`
		Rows   int    `form:"rows" json:"rows" `
		Status int    `form:"status" json:"status"` //是否可交易状态
		In     int    `form:"in" json:"in"`
		Out    int    `form:"out" json:"out"`
		Name   string `form:"name" json:"name"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	fmt.Println("---------------->", req.Name)
	list, err := new(models.Tokens).GetSystemList(req.Page, req.Rows, req.Status, req.In, req.Out, req.Name)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.Put(c, "list", list)
	w.RespOK(c)
	return
}

//系统设置 币种配置
func (w *WebUserManageController) TokenSystemAdd(c *gin.Context) {
	// 获取参数
	req := struct {
		Id                   int     `form:"id"  json:"id" `
		Mark                 string  `form:"mark" json:"mark" binding:"required" `
		Detail               string  `form:"detail" json:"detail" binding:"required" `
		Logo                 string  `form:"logo" json:"logo" binding:"required" json:"logo"`
		Status               int     `form:"status"  json:"status" binding:"required" json:"status"`
		InTokenMark          int     `form:"in_mark" json:"in_mark" binding:"required" json:"in_mark"`
		InTokenLeastBalance  float64 `form:"in_least_balance" json:"in_least_balance" binding:"required" json:"in_least_balance"`
		OutTokenMark         int     `form:"out_mark" json:"out_mark" binding:"required" json:"out_mark"`
		OutTokenLeastBalance float64 `form:"out_least_balance" json:"out_least_balance" binding:"required" json:"out_least_balance"`
		OutTokenFee          float32 `form:"out_fee" json:"out_fee" binding:"required" json:"out_fee"`
		InRemarks            string  `form:"in_remarks"   json:"in_remarks"`
		OutRemarks           string  `form:"out_remarks" json:"out_remarks"`
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		w.RespErr(c, err)
		return
	}
	err = new(models.Tokens).TokensSystemAdd(models.Tokens{
		Id:                   req.Id,
		Mark:                 req.Mark,
		Detail:               req.Detail,
		Logo:                 req.Logo,
		Status:               req.Status,
		InTokenLeastBalance:  new(models.Tokens).Float64ToInt64By8Bit(req.InTokenLeastBalance),
		OutTokenLeastBalance: new(models.Tokens).Float64ToInt64By8Bit(req.OutTokenLeastBalance),
		InTokenMark:          req.InTokenMark,
		OutTokenMark:         req.OutTokenMark,
		OutTokenFee:          req.OutTokenFee,
		InRemarks:            req.InRemarks,
		OutRemarks:           req.OutRemarks,
	})
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.RespOK(c)
	return
}

//邀请人统计表—账号：18888888888
func (w *WebUserManageController) GetInviteInfoList(c *gin.Context) {
	req := struct {
		Uid     int    `form:"uid" json:"uid" binding:"required"`
		Page    int    `form:"page" json:"page" binding:"required"`
		Rows    int    `form:"rows" json:"rows" `
		Date    uint64 `form:"date" json:"date" `       //日期
		Name    string `form:"name" json:"name" `       //渠道名称
		Account string `form:"account" json:"account" ` //刷选
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	list, err := new(models.UserEx).GetInviteInfoList(req.Uid, req.Page, req.Rows, req.Date, req.Name, req.Account)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.Put(c, "list", list)
	w.RespOK(c)
	return
}

//被邀请人列表
func (w *WebUserManageController) GetInviteList(c *gin.Context) {
	req := struct {
		Page   int    `form:"page" json:"page" binding:"required"`
		Rows   int    `form:"rows" json:"rows" `
		Date   uint64 `form:"time" json:"time" `     //日期
		Status int    `form:"status" json:"status" ` //用户状态
		Search string `form:"search" json:"search" ` //刷选
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	list, err := new(models.UserEx).GetInViteList(req.Page, req.Rows, req.Search)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.Put(c, "list", list)
	w.RespOK(c)
	return
}

func (w *WebUserManageController) GetTradeRule(c *gin.Context) {
	result, err := new(models.ConfigureTradeRule).GetTradeRule()
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.Put(c, "list", result)
	w.RespOK(c)
}

func (w *WebUserManageController) SetTradeRule(c *gin.Context) {

	req := struct {
		Cuid        int   `form:"cuid" json:"cuid" `
		Muid        int   `form:"muid" json:"muid" `
		OneTradeMax int64 `form:"one_trade_max" json:"one_trade_max" `
		OneTotal    int64 `form:"one_total" json:"one_total" `
		TwoTotal    int64 `form:"two_total" json:"two_total" `
		TwoTradeMax int64 `form:"two_trade_max" json:"two_trade_max" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	fmt.Println("two_total=", req.TwoTotal)
	param := models.ConfigureTradeRule{
		Cuid:        req.Cuid,
		Muid:        req.Muid,
		OneTradeMax: req.OneTradeMax,
		OneTotal:    req.OneTotal,
		TwoTotal:    req.TwoTotal,
		TwoTradeMax: req.TwoTradeMax,
	}
	err = new(models.ConfigureTradeRule).AddTradeRule(param)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.RespOK(c)
	return
}

func (w *WebUserManageController) SecondAffirm(c *gin.Context) {
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
	err = new(models.WebUser).SecondAffirmLimit(req.Uid, req.Status)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.RespOK(c)
	return
}

func (w *WebUserManageController) FirstAffirm(c *gin.Context) {
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
	err = new(models.WebUser).FirstAffirmLimit(req.Uid, req.Status)
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
	result, err := new(models.UserEx).GetFirstDetail(req.Uid)
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

func (w *WebUserManageController) GetSecondCertificationList(c *gin.Context) {
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
	list, err := new(models.UserLoginLog).GetUserLoginLogList(req.Page, req.Rows, req.TerminalType, req.Status, req.LoginTime, req.Search)
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
		Date   uint64 `form:"date",json:"date" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	//
	list, err := new(models.TotalProperty).GetTotalProperty(req.Page, req.Rows, req.Status, req.Date, req.Search)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		w.RespErr(c, err)
		return
	}

	w.Put(c, "list", list)
	w.RespOK(c, "成功")
	return
}

func (w *WebUserManageController) GetTotalUser(c *gin.Context) {
	//fmt.Printf("list param %#v\n", req)
	total, upday, upweek, err := new(models.WebUser).GetTotalUser()
	if err != nil {
		w.RespErr(c, err)
		return
	}
	w.Put(c, "upweek", upweek)
	w.Put(c, "upday", upday)
	w.Put(c, "total", total)
	w.RespOK(c)
	return
}

func (w *WebUserManageController) GetWebUserList(c *gin.Context) {
	req := struct {
		Page   int    `form:"page" json:"page" binding:"required"`
		Rows   int    `form:"rows" json:"rows" `
		Search string `form:"search" json:"search" `
		Date   int64  `form:"date" json:"date" `
		Verify int    `form:"verify" json:"verify" `
		Status int    `form:"status" json:"status" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		w.RespErr(c, err)
		return
	}
	result, err := new(models.WebUser).UserList(req.Page, req.Rows, req.Verify, req.Status, req.Search, req.Date)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	fmt.Println("0.00.0.0.0.0.000000000000000", result)
	list, Ok := result.Items.([]models.UserGroup)

	fmt.Println("GetWebUserList-1")
	if !Ok {
		fmt.Println("GetWebUserList-2")
		w.RespErr(c, err)
		return
	}
	list = w.VerifyOperator(list)
	if err != nil {
		w.RespErr(c, err)
		return
	}
	// 设置返回数据
	w.Put(c, "list", result)
	fmt.Println("GetWebUserList-3")
	// 返回
	w.RespOK(c)
	//c.JSON(http.StatusOK, gin.H{"code": 0, "page": page, "data": reuslt, "total": total, "msg": "成功"})
	return
}

func (w *WebUserManageController) VerifyOperator(list []models.UserGroup) []models.UserGroup {
	for index, _ := range list {

		if list[index].SecurityAuth&utils.AUTH_GOOGLE == utils.AUTH_GOOGLE {
			list[index].GoogleVerifyMark = 1
		}
		if list[index].SecurityAuth&utils.AUTH_TWO == utils.AUTH_TWO {
			list[index].TWOVerifyMark = 1
		}
		if list[index].SecurityAuth&utils.AUTH_FIRST == utils.AUTH_FIRST {
			list[index].RealNameVerifyMark = 1
		}

	}
	return nil
}
