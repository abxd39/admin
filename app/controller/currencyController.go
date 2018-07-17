package controller

import (
	"admin/app/models"
	"admin/constant"
	"admin/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CurrencyController struct {
	BaseController
}

func (this *CurrencyController) Router(r *gin.Engine) {
	g := r.Group("/currency")
	{
		g.GET("/list", this.GetTradeList)                     //p4-2-0法币挂单管理
		g.POST("/down_trade_order", this.DownTradeAds)        //p4-2-0法币挂单管理 下架交易单
		g.GET("/tokens", this.GetTokensList)                  //获取 所有数据货币的名称及货币Id
		g.GET("/order_list", this.GetOderList)                //p4-2-1法币成交管理
		g.GET("/total_balance", this.GetTotalCurrencyBalance) //p2-3-1法币账户统计列表
		g.GET("/user_detail", this.GetUserDetailList)         //p2-3-1-2法币账户资产展示
		g.GET("/user_buysell", this.GetBuySellList)           //p2-3-1-1查看统计买入_卖出_划转
		g.GET("total", this.Total)                            //p2-3-0总财产列表
		//g.GET("/")                                            //p2-3-3法币账户变更详情
		//g.GET("/")                                            //p2-3-0-0币数统计列表
	}
}

//总财产统计列表
func (cu *CurrencyController) Total(c *gin.Context) {
	req := struct {
		Page   int    `form:"page" json:"page" binding:"required"`
		Rows   int    `form:"rows" json:"rows" `
		Search string `form:"search" json:"search" ` //搜索的内容
		Status int    `form:"status" json:"status" ` //用户账号状态
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}
	result, err := new(models.UserGroup).GetAllUser(req.Page, req.Rows, req.Status, req.Search)
	if err != nil {
		cu.RespErr(c, err)
	}
	uidlist := make([]uint64, 0)
	value, OK := result.Items.([]models.UserGroup)
	if !OK {
		cu.RespErr(c, errors.New("assert [] userGroup failed!!"))
		return
	}
	for _, value := range value {
		uidlist = append(uidlist, value.Uid)
	}
	//资产总折合
	//法币账户折合
	currencylist, err := new(models.UserCurrency).GetAll(uidlist)
	for i, _ := range value {
		for _, v := range currencylist {
			if value[i].Uid == v.Uid {
				value[i].TotalCurrentCNY = v.Balance
				value[i].LockCurrentCNY = v.Freeze
				break
			}
		}
	}
	//币币账户折合
	cu.Put(c, "list", result)
	cu.RespOK(c)
}

func (cu *CurrencyController) DownTradeAds(c *gin.Context) {
	req := struct {
		Id  int `form:"id" json:"id" binding:"required"`   //订单id
		Uid int `form:"uid" json:"uid" binding:"required"` //用户uid
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}
	err = new(models.Ads).DownTradeAds(req.Id, req.Uid)
	if err != nil {
		cu.RespErr(c, err)
		return

	}
	cu.RespOK(c)
}

//查看法币统计买入_卖出_划转
func (cu *CurrencyController) GetBuySellList(c *gin.Context) {
	req := struct {
		Uid      int `form:"uid" json:"uid" binding:"required"`
		Page     int `form:"page" json:"page" binding:"required"`
		Rows     int `form:"rows" json:"rows" `
		Token_id int `form:"token_id" json:"token_id"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}
	uid := make([]int, 0)
	uid = append(uid, req.Uid)
	fmt.Printf("GetBuySellList%#v\n", uid)
	list, err := new(models.Order).GetOrderListOfUid(req.Page, req.Rows, req.Uid, req.Token_id)
	if err != nil {
		cu.RespErr(c, err)
		return
	}

	cu.Put(c, "list", list)
	cu.RespOK(c)
	return
}

func (cu *CurrencyController) GetUserDetailList(c *gin.Context) {
	req := struct {
		Uid      int `form:"uid" json:"uid" binding:"required"`
		Page     int `form:"page" json:"page" binding:"required"`
		Rows     int `form:"rows" json:"rows" `
		Token_id int `form:"token_id" json:"token_id"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, constant.RESPONSE_CODE_ERROR, "参数错误")
		return
	}

	result, page, total, err := new(models.UserCurrency).GetCurrencyList(req.Page, req.Rows, req.Uid, req.Token_id)
	if err != nil {
		cu.RespErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "page": page, "total": total, "data": result, "msg": "成功"})
	return
}

func (cu *CurrencyController) GetTotalCurrencyBalance(c *gin.Context) {
	req := struct {
		Page     int    `form:"page" json:"page" binding:"required"`
		Page_num int    `form:"rows" json:"rows" `
		Search   string `form:"search" json:"search" ` //搜索的内容
		Status   int    `form:"status" json:"status" ` //用户账号状态
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}
	fmt.Println(".0.................0.0.0.0.0.0.0.0.......")
	result, err := new(models.UserGroup).GetAllUser(req.Page, req.Page_num, req.Status, req.Search)
	if err != nil {
		cu.RespErr(c, err)
	}

	//法币账户折和没有计算
	cu.Put(c, "list", result)

	// 返回
	cu.RespOK(c, "成功")
	return
}

//法币挂单管理
func (cu *CurrencyController) GetTradeList(c *gin.Context) {
	req := struct {
		Page    int    `form:"page" json:"page" binding:"required"`
		PageNum int    `form:"rows" json:"rows" `
		Ustatus int    `form:"status" json:"status" ` //用户登录状态
		Search  string `form:"search" json:"search" `
		Verify  int    `form:"verify" json:"verify" `     //实名认证 二级认证 google 验证  交易权限
		Date    string `form:"date" json:"date" `         //挂单日期
		TokenId int    `form:"token_id" json:"token_id" ` //货币名称
		TradeId int    `form:"tid" json:"tid" `           //交易方向
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}
	list, err := new(models.Ads).GetAdsList(req.Page, req.PageNum, req.Ustatus, req.TokenId, req.TradeId, req.Verify, req.Search, req.Date)
	if err != nil {
		cu.RespErr(c, err)
		return
	}
	cu.Put(c, "list", list)
	cu.RespOK(c)
	return

}

func (cu *CurrencyController) GetTokensList(c *gin.Context) {
	fmt.Println("tttttttttttttttttttttttttttttttttttttttt")
	list, err := new(models.Tokens).GetTokenList()
	if err != nil {
		cu.RespErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": list, "msg": "成功"})
	return
}

//法币成交列表
func (cu *CurrencyController) GetOderList(c *gin.Context) {
	//参数一大堆
	req := struct {
		Page     int    `form:"page" json:"page" binding:"required"`
		Page_num int    `form:"rows" json:"rows" `
		Start_t  string `form:"start_t" json:"start_t" `
		Search   string `form:"search" json:"search" `     //筛选
		Status   int    `form:"status" json:"status" `     //订单状态
		Token_id int    `form:"token_id" json:"token_id" ` //货币名称
		Ad_type  int    `form:"adtype" json:"adtype" `     //买卖方向
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}
	list, err := new(models.Order).GetOrderList(req.Page, req.Page_num, req.Ad_type, req.Status, req.Token_id, req.Start_t, req.Search)
	if err != nil {
		cu.RespErr(c, err)
		return
	}
	cu.Put(c, "list", list)
	cu.RespOK(c)
	return
}
