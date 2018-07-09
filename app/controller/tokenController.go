package controller

import (
	"admin/app/models"
	"admin/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenController struct{}

func (this *TokenController) Router(r *gin.Engine) {
	g := r.Group("/token")
	{
		g.GET("/list", this.GetTokenOderList)            //bibi 挂单信息
		g.GET("/record_list", this.GetRecordList)        //bibi 成交记录
		g.GET("/total_balance", this.GetTokenBalance)    //bibi 所有用户 总资产（币币总资产）
		g.GET("/user_token_detail", this.GetTokenDetail) //bibi账户资产展示
	}
}

func (this *TokenController) GetTokenDetail(c *gin.Context) {
	req := struct {
		Uid      int `form:"uid" json:"uid" binding:"required"`
		Token_id int `form:"token_id" json:"token_id"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}
	//bibi账户余额
	list, err := new(models.UserToken).GetTokenDetailOfUid(req.Uid, req.Token_id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
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
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}
	fmt.Printf("GetTokenBalance%#v\n", req)
	list, page, toal, err := new(models.PersonalProperty).TotalUserBalance(req.Page, req.Page_num, req.Status)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "page": page, "total": toal, "data": list, "msg": "成功"})
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
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}
	list, page, toal, err := new(models.EntrustDetail).GetTokenRecordList(req.Page, req.Page_num, req.Trade_id, req.Trade_duad, req.Ad_id, req.Start_t, req.End_t)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
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
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}
	list, page, toal, err := new(models.EntrustDetail).GetTokenOrderList(req.Page, req.Page_num, req.Trade_id, req.Trade_duad, req.Ad_id, req.Status, req.Start_t, req.End_t)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "page": page, "total": toal, "data": list, "msg": "成功"})
	return
}
