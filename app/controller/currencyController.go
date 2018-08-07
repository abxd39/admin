package controller

import (
	"admin/app/models"
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
		g.GET("/list", this.GetTradeList)                        //p4-2-0法币挂单管理
		g.POST("/down_trade_order", this.DownTradeAds)           //p4-2-0法币挂单管理 下架交易单
		g.GET("/tokens", this.GetTokensList)                     //获取 所有数据货币的名称及货币Id
		g.GET("/order_list", this.GetOderList)                   //p4-2-1法币成交管理
		g.GET("/total_balance", this.GetTotalCurrencyBalance)    //p2-3-1法币账户统计列表
		g.GET("/user_detail", this.GetUserDetailList)            //p2-3-1-2法币账户资产展示
		g.GET("/user_buysell", this.GetBuySellList)              //p2-3-1-1查看统计买入_卖出_划转
		g.GET("/total", this.Total)                              //p2-3-0总财产列表
		g.GET("/currency_change", this.GetCurrencyChangeHistory) //p2-3-3法币账户变更详情
		//g.GET("/")                                               //p2-3-0-0币数统计列表
		//划转到币币账户货币数量日统计
		g.GET("/layoff_list", this.GetLayOffList)
	}
}

func (cu *CurrencyController) GetLayOffList(c *gin.Context) {

}

func (cu *CurrencyController) GetCurrencyChangeHistory(c *gin.Context) {
	req := struct {
		Page   int    `form:"page" json:"page" binding:"required"`
		Rows   int    `form:"rows" json:"rows" `
		Search string `form:"search" json:"search" `               //搜索的内容
		Status int    `form:"status" json:"status" `               //用户账号状态
		Date   string `form:"date" json:"date" binding:"required"` //日期
		Chtype int    `form:"type" json:"type"`                    // 买入 卖出 提币 充币 划转
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}
	//把货币Id转换为货币名称
	tokenlist, err := new(models.CommonTokens).GetTokenList()
	if err != nil {
		cu.RespErr(c, err)
		return
	}
	ulist, err := new(models.UserCurrencyHistory).GetListForUid(req.Page, req.Rows, req.Status, req.Chtype, req.Search, req.Date)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}
	list, Ok := ulist.Items.([]models.UserCurrencyHistory)
	if !Ok {
		cu.RespErr(c, errors.New("assert type UserGroup failed!!"))
		return
	}
	for i, v := range list {
		for _, vt := range tokenlist {
			if vt.Id == uint32(v.TokenId) {
				list[i].TokenName = vt.Mark
				break
			}
		}
	}

	//if req.Search != `` || req.Status != 0 {
	//	// 先帅选用户资料
	//	list, err := new(models.UserGroup).GetAllUser(req.Page, req.Rows, req.Status, req.Search)
	//	if err != nil {
	//		cu.RespErr(c, err)
	//		return
	//	}
	//	value, Ok := list.Items.([]models.UserGroup)
	//	if !Ok {
	//		cu.RespErr(c, errors.New("assert type UserGroup failed!!"))
	//		return
	//	}
	//	uidlist := make([]int64, 0)
	//	for _, v := range value {
	//		uidlist = append(uidlist, v.Uid)
	//	}
	//	histroyList, err := new(models.UserCurrencyHistory).GetListForUid(req.Page, req.Rows, uidlist)
	//	if err != nil {
	//		cu.RespErr(c, err)
	//		return
	//	}
	//	histroyValue, ok := histroyList.Items.([]models.UserCurrencyHistory)
	//	if !ok {
	//		cu.RespErr(c, errors.New("assert type UserCurrencyHistory failed!!"))
	//		return
	//	}
	//
	//	for i, _ := range histroyValue {
	//		for _, v := range value {
	//			if histroyValue[i].Uid == int(v.Uid) {
	//				histroyValue[i].Email = v.Email
	//				histroyValue[i].NickName = v.NickName
	//				histroyValue[i].Phone = v.Phone
	//				histroyValue[i].Status = v.Status
	//				break
	//			}
	//		}
	//		for _, vt := range tokenlist {
	//			if int(vt.Id) == histroyValue[i].TokenId {
	//				histroyValue[i].TokenName = vt.Name
	//				break
	//			}
	//
	//		}
	//	}
	//	histroyList.Items = histroyValue
	//	cu.Put(c, "list", histroyList)
	//	cu.RespOK(c)
	//	return
	//
	//} else {
	//	list, err := new(models.UserCurrencyHistory).GetList(req.Page, req.Rows, req.Chtype, req.Date)
	//	if err != nil {
	//		cu.RespErr(c, err)
	//		return
	//	}
	//	uidList := make([]uint64, 0)
	//	Value, ok := list.Items.([]models.UserCurrencyHistory)
	//	if !ok {
	//		cu.RespErr(c, err)
	//		return
	//	}
	//	for _, v := range Value {
	//		uidList = append(uidList, uint64(v.Uid))
	//	}
	//	//
	//	ulist, err := new(models.UserGroup).GetUserListForUid(uidList)
	//	if err != nil {
	//		cu.RespErr(c, err)
	//		return
	//	}
	//
	//	for i, _ := range Value {
	//		for _, v := range ulist {
	//			if Value[i].Uid == int(v.Uid) {
	//				Value[i].NickName = v.NickName
	//				Value[i].Phone = v.Phone
	//				Value[i].Email = v.Email
	//				Value[i].Status = v.Status
	//				break
	//			}
	//		}
	//		for _, vt := range tokenlist {
	//			if int(vt.Id) == Value[i].TokenId {
	//				Value[i].TokenName = vt.Name
	//				break
	//			}
	//
	//		}
	//
	//	}
	//	list.Items = Value
	//	cu.Put(c, "list", list)
	//	cu.RespOK(c)
	//	return
	//}
	cu.Put(c, "list", list)
	cu.RespOK(c)
	return
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
	uidlist := make([]int64, 0)
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
			if value[i].Uid == int64(v.Uid) {
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
		cu.RespErr(c, err)
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

//p2-3-1法币账户统计列表
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
	//result, err := new(models.UserGroup).GetAllUser(req.Page, req.Page_num, req.Status, req.Search)
	result, err := new(models.AmountToCny).CurrencyBalance(req.Page, req.Page_num, req.Status, req.Search)
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
	list, err := new(models.CommonTokens).GetTokenList()
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
