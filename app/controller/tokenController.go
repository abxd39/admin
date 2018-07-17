package controller

import (
	"admin/app/models"
	"admin/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenController struct {
	BaseController
}

func (this *TokenController) Router(r *gin.Engine) {
	g := r.Group("/token")
	{
		g.GET("/list", this.GetTokenOderList) //bibi p4-1-0 币币委托管理 挂单信息

		g.GET("/record_list", this.GetRecordList)        //bibi p4-1-1 成交记录
		g.GET("/total_balance", this.GetTokenBalance)    //bibi p2-3-2币币账户统计列表
		g.GET("/user_token_detail", this.GetTokenDetail) //p2-3-2-1查看币币账户资产
		g.GET("/token_cash", this.GetTokenCashList)      //p4-1-2币兑管理
		// g.GET("/delete_cash", this.DeleteCash)           //删除币兑
		// g.GET("/modify_cash", this.ModifyCash)           //修改币兑
		// g.GET("add_cash", this.AddCash)                  //删除
		g.GET("/change_detail", this.ChangeDetail) //p2-3-4币币账户变更详情
	}
}

func (this *TokenController) ChangeDetail(c *gin.Context) {
	req := struct {
		Page    int    `form:"page" json:"page" binding:"required"`
		Rows    int    `form:"rows" json:"rows" `
		Start_t string `form:"start_t" json:"start_t" `
		Search  string `form:"search" json:"search" ` //刷选
		Ad_id   int    `form:"ad_id" json:"ad_id" `   //交易方向 买 卖 划转
		Status  int    `form:"status" json:"status" ` //用户状态
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	//没写完
}

func (this *TokenController) GetTokenCashList(c *gin.Context) {
	req := struct {
		Page    int `form:"page" json:"page" binding:"required"`
		Rows    int `form:"rows" json:"rows" `
		TokenId int `form:"token_id" json:"token_id" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.QuenesConfig).GetTokenCashList(req.Page, req.Rows, req.TokenId)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c, "成功")
	return
}

func (this *TokenController) GetTokenDetail(c *gin.Context) {
	req := struct {
		Uid      int `form:"uid" json:"uid" binding:"required"`
		Token_id int `form:"token_id" json:"token_id"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	//bibi账户余额
	list, err := new(models.UserToken).GetTokenDetailOfUid(req.Uid, req.Token_id)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": list, "msg": "成功"})
	return
}

//bibi 账户统计表
func (this *TokenController) GetTokenBalance(c *gin.Context) {
	req := struct {
		Page     int `form:"page" json:"page" binding:"required"`
		Page_num int `form:"rows" json:"rows" `
		Status   int `form:"status" json:"status" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	fmt.Printf("GetTokenBalance%#v\n", req)
	list, err := new(models.PersonalProperty).TotalUserBalance(req.Page, req.Page_num, req.Status)
	if err != nil {
		this.RespErr(c, err)
		return
	}

	this.Put(c, "list", list)
	this.RespOK(c, "成功")
	return
}

//bibi 成交记录
func (this *TokenController) GetRecordList(c *gin.Context) {
	req := struct {
		Page       int    `form:"page" json:"page" binding:"required"`
		Page_num   int    `form:"rows" json:"rows" `
		Trade_id   int    `form:"trade_id" json:"trade_id" ` //交易类型id 市价交易or 限价交易
		Start_t    string `form:"start_t" json:"start_t" `
		End_t      string `form:"end_t" json:"end_t" `
		Trade_duad int    `form:"trade_duad" json:"trade_duad" ` //交易对
		Ad_id      int    `form:"ad_id" json:"ad_id" `           //买卖方向
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, page, toal, err := new(models.EntrustDetail).GetTokenRecordList(req.Page, req.Page_num, req.Trade_id, req.Trade_duad, req.Ad_id, req.Start_t, req.End_t)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "page": page, "total": toal, "data": list, "msg": "成功"})
	return
}

//币币挂单列表
func (this *TokenController) GetTokenOderList(c *gin.Context) {
	req := struct {
		Page       int    `form:"page" json:"page" binding:"required"`
		Page_num   int    `form:"rows" json:"rows" `
		Trade_id   int    `form:"trade_id" json:"trade_id" ` //交易类型id 市价交易or 限价交易
		Start_t    string `form:"start_t" json:"start_t" `
		End_t      string `form:"end_t" json:"end_t" `
		Trade_duad int    `form:"trade_duad" json:"trade_duad" ` //交易对
		Ad_id      int    `form:"ad_id" json:"ad_id" `           //买卖方向
		Status     int    `form:"status" json:"staus" `          //订单状态
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, page, toal, err := new(models.EntrustDetail).GetTokenOrderList(req.Page, req.Page_num, req.Trade_id, req.Trade_duad, req.Ad_id, req.Status, req.Start_t, req.End_t)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "page": page, "total": toal, "data": list, "msg": "成功"})
	return
}
