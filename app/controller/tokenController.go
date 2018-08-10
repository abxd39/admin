package controller

import (
	"admin/app/models"
	"admin/utils"
	"errors"
	"fmt"

	"time"

	"admin/constant"
	"digicon/common/convert"
	"github.com/gin-gonic/gin"
	"regexp"
)

type TokenController struct {
	BaseController
}

func (this *TokenController) Router(r *gin.Engine) {
	g := r.Group("/token")
	{
		g.GET("/list", this.GetTokenOderList)            //bibi p4-1-0 币币委托管理 挂单信息
		g.POST("/evacuate_order", this.EvacuateOder)     // p4-1-0 币币委托管理 挂单信息 撤单
		g.GET("/record_list", this.GetRecordList)        //bibi p4-1-1 成交记录
		g.GET("/total_balance", this.GetTokenBalance)    //bibi p2-3-2币币账户统计列表
		g.GET("/user_token_detail", this.GetTokenDetail) //p2-3-2-1查看币币账户资产
		g.GET("/token_cash_list", this.GetTokenCashList) //p4-1-2币兑管理
		g.POST("/delete_cash", this.DeleteCash)          //删除币兑
		g.GET("/modify_cash", this.ModifyCash)           //修改币兑
		g.POST("/add_cash", this.AddCash)                //添加币兑
		g.GET("/change_detail", this.ChangeDetail)       //p2-3-4币币账户变更详情
		g.GET("/fee_list", this.GetFeeInfoList)          //p5-1-0-1币币交易手续费明细
		g.GET("/add_take_list", this.GetAddTakeList)     //p5-1-1-1提币手续费明细
		g.GET("/total_trade", this.GetTradeTotalList)    //p5-1-0币币交易手续费汇总
		//提币 充币管理
		g.GET("/io_token_list", this.GetTokenInList) //
		g.POST("/opt_token", this.OptTakeToken)      //
		//日提币充币汇总
		g.GET("/total_token_list", this.GetTotalTokenList) //
		g.GET("/total_token_info_list", this.GetTotalTokenInfoList)

		//日划转汇总
		//币币交易手续费汇总
		g.GET("/token_order_fee_total", this.GetOderFeeTotalList)
		//提币手续费汇总表
		g.GET("/token_inout_daily_sheet", this.GetTokenInOutDailySheetList)
		//仪表盘
		//手续费走势图
		g.GET("/fee_trend_map", this.GetFeeTrendMap)
		g.GET("/test1", this.test1)

		g.GET("/list_transfer_daily_sheet", this.ListTransferDailySheet)

	}
}

func (this *TokenController) test1(c *gin.Context) {
	new(models.TokenFeeDailySheet).Test1()
	new(models.CurencyFeeDailySheet).BoottimeTimingSettlement()
	new(models.WalletInoutDailySheet).BoottimeTimingSettlement()
	new(models.TokenFeeDailySheet).BoottimeTimingSettlement()
	this.RespOK(c)
	return
}

func (this *TokenController) GetFeeTrendMap(c *gin.Context) {
	//手续费 注：当天
	//币币交易手续费
	tokenFee, err := new(models.Trade).GetTodayFee()
	if err != nil {
		this.RespErr(c, err)
		return
	}
	//fmt.Println("---------------->0")
	//法币交易手续费
	currencyTotalFee, err := new(models.Order).GetOrderDayFee()
	if err != nil {
		this.RespErr(c, err)
		return
	}
	//fmt.Println("---------------->1")
	//提币手续费
	outTokenFee, err := new(models.TokenInout).GetOutTokenFee()
	if err != nil {
		this.RespErr(c, err)
		return
	}
	//fmt.Println("---------------->2")
	//this.Put(c,"tFee",)
	this.Put(c, "tradeFee", currencyTotalFee+tokenFee)
	this.Put(c, "oFee", outTokenFee)
	date := fmt.Sprintf("%02d%02d", time.Now().Month(), time.Now().Day())
	this.Put(c, "date", date)
	this.RespOK(c)
	return
}

//提币手续费汇总表
func (this *TokenController) GetTokenInOutDailySheetList(c *gin.Context) {
	req := struct {
		Page    int    `form:"page" json:"page" binding:"required"`
		Rows    int    `form:"rows" json:"rows" `
		Date    string `form:"date",json:"date"` //日期
		TokenId int    `form:"tid",json:"tid"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.WalletInoutDailySheet).GetInOutDailSheetList(req.Page, req.Rows, req.TokenId, req.Date)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
	return
}

//币币交易手续费汇总
func (this *TokenController) GetOderFeeTotalList(c *gin.Context) {
	req := struct {
		Page int    `form:"page" json:"page" binding:"required"`
		Rows int    `form:"rows" json:"rows" `
		Date uint64 `form:"date",json:"date"` //日期
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, data, err := new(models.TokenFeeDailySheet).GetDailySheetList(req.Page, req.Rows, req.Date)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "data", data)
	this.Put(c, "list", list)
	this.RespOK(c)
	return
}

//冲提币汇总
func (this *TokenController) GetTotalTokenList(c *gin.Context) {
	req := struct {
		Opt     int    `form:"opt" json:"opt" binding:"required"` //1充 2 提币
		Page    int    `form:"page" json:"page" binding:"required"`
		Rows    int    `form:"rows" json:"rows" `
		TokenId int    `form:"tokenId" json:"tokenId" ` //货币id
		Date    string `form:"date",json:"date"`        //日期
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.TokenInout).GetTotalList(req.Page, req.Rows, req.TokenId, req.Opt, req.Date)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
}

func (this *TokenController) GetTotalTokenInfoList(c *gin.Context) {
	req := struct {
		Opt     int    `form:"opt" json:"opt" binding:"required"` //1充 2 提币
		Page    int    `form:"page" json:"page" binding:"required"`
		Rows    int    `form:"rows" json:"rows" `
		TokenId int    `form:"tokenId" json:"tokenId" `              //货币id
		Date    string `form:"date",json:"date"  binding:"required"` //日期
		Search  string `form:"search",json:"search"`                 //输入uid 搜索
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.TokenInout).GetTotalInfoList(req.Page, req.Rows, req.TokenId, req.Opt, req.Date, req.Search)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
}

//充提币管理
func (this *TokenController) GetTokenInList(c *gin.Context) {
	req := struct {
		Page    int    `form:"page" json:"page" binding:"required"`
		Rows    int    `form:"rows" json:"rows" `
		Ustatus int    `form:"ustatus" json:"ustatus" ` // 用户状态
		Status  int    `form:"status" json:"status" `   //提币状态
		TokenId int    `form:"tokenId" json:"tokenId" ` //货币id
		Opt     int    `form:"opt" json:"opt"  `        //操作方向
		Search  string `form:"search" json:"search" `   //筛选
		Date    string `form:"date",json:"date"`        //日期
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.TokenInout).GetTokenInList(req.Page, req.Rows, req.Ustatus, req.Status, req.TokenId, req.Opt, req.Search, req.Date)
	if err != nil {
		utils.AdminLog.Println(err.Error())
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
	return
}

func (this *TokenController) OptTakeToken(c *gin.Context) {
	req := struct {
		Id  int `form:"id" json:"id" binding:"required"`
		Uid int `form:"uid" json:"uid" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	err = new(models.TokenInout).OptTakeToken(req.Id, req.Uid)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.RespOK(c)
	return
}

func (this *TokenController) GetTradeTotalList(c *gin.Context) {
	req := struct {
		Page int    `form:"page" json:"page" binding:"required"`
		Rows int    `form:"rows" json:"rows" `
		Date uint64 `form:"date",json:"date"` //日期
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.Trade).TotalTotalTradeList(req.Page, req.Rows, req.Date)
	if err != nil {
		utils.AdminLog.Println(err.Error())
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
	return
}

//p5-1-1-1提币手续费明细
func (this *TokenController) GetAddTakeList(c *gin.Context) {
	req := struct {
		Page     int    `form:"page" json:"page" binding:"required"`
		Rows     int    `form:"rows" json:"rows" `
		Token_id int    `form:"token_id" json:"token_id" ` //币种
		Uid      int    `form:"uid" json:"uid" `
		Date     uint64 `form:"date",json:"date"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.TokenHistoryGroup).GetAddTakeList(req.Page, req.Rows, req.Token_id, req.Uid, req.Date)
	if err != nil {
		utils.AdminLog.Println(err.Error())
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
	return
}

//p5-1-0-1币币交易手续费明细
func (this *TokenController) GetFeeInfoList(c *gin.Context) {
	req := struct {
		Page int    `form:"page" json:"page" binding:"required"`
		Rows int    `form:"rows" json:"rows" `
		Opt  int    `form:"opt" json:"opt" ` //交易方向
		Uid  int    `form:"uid" json:"uid" `
		Date uint64 `form:"date",json:"date"`
		Name string `form:"name" json:"name"` //交易对名称 USDT/BTC
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.Trade).GetFeeInfoList(req.Page, req.Rows, req.Uid, req.Opt, req.Date, req.Name)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
}

func (this *TokenController) DeleteCash(c *gin.Context) {
	req := struct {
		Id int `form:"id" json:"id" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	err = new(models.ConfigQuenes).DeleteCash(req.Id)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.RespOK(c)
	return
}

func (this *TokenController) ModifyCash(c *gin.Context) {
	req := struct {
		Id int `form:"id" json:"id" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	result, err := new(models.ConfigQuenes).ModifyCash(req.Id)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "data", result)

	this.RespOK(c)
	return
}

func (this *TokenController) AddCash(c *gin.Context) {
	req := models.ConfigQuenes{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	err = new(models.ConfigQuenes).AddCash(&req)
	if err != nil {
		fmt.Println("0cccccccccccccccccc0", err)
		this.RespErr(c, err)
		return
	}
	this.RespOK(c)
	return
}

//订单操作 撤单处理

func (this *TokenController) EvacuateOder(c *gin.Context) {
	req := struct {
		Uid      int    `form:"uid" json:"uid" binding:"required"`
		OerderId string `form:"oid" json:"oid" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	//new(models.WalletInoutDailySheet).BoottimeTimingSettlement()
	err = new(models.EntrustDetail).EvacuateOder(req.Uid, req.OerderId)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.RespOK(c)
	return
}

func (this *TokenController) ChangeDetail(c *gin.Context) {
	req := struct {
		Page   int    `form:"page" json:"page" binding:"required"`
		Rows   int    `form:"rows" json:"rows" `
		Date   uint64 `form:"date" json:"date" binding:"required"`
		Search string `form:"search" json:"search" ` //刷选
		Type   int    `form:"type" json:"type" `     //交易方向 买 卖 划转
		Status int    `form:"status" json:"status" ` //用户状态
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	//把货币Id转换为货币名称
	tokenlist, err := new(models.CommonTokens).GetTokenList()
	if err != nil {
		this.RespErr(c, err)
		return
	}
	fmt.Println("------------------------>what funck you ", req.Status)
	mlist, err := new(models.MoneyRecord).GetMoneyListForDateOrType(req.Page, req.Rows, req.Type, req.Status, req.Date, req.Search)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	list, Ok := mlist.Items.([]models.MoneyRecord)
	if !Ok {
		this.RespErr(c, errors.New("assert type UserGroup failed!!"))
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
	//没写完
	//在此分道扬镳
	//if req.Search != `` || req.Status != 0 {
	//	list, err := new(models.UserGroup).GetAllUser(req.Page, req.Rows, req.Status, req.Search)
	//	if err != nil {
	//		this.RespErr(c, err)
	//		return
	//	}
	//	value, Ok := list.Items.([]models.UserGroup)
	//	if !Ok {
	//		this.RespErr(c, errors.New("assert type UserGroup failed!!"))
	//		return
	//	}
	//	uidlist := make([]int64, 0)
	//	for _, v := range value {
	//		uidlist = append(uidlist, v.Uid)
	//	}
	//	monerylist, err := new(models.MoneyRecord).GetMoneyList(req.Page, req.Rows, uidlist)
	//	if err != nil {
	//		this.RespErr(c, err)
	//		return
	//	}
	//	tokenValue, ok := monerylist.Items.([]models.MoneyRecord)
	//	if !ok {
	//		this.RespErr(c, errors.New("assert type MoneyRecord failed!!"))
	//		return
	//	}
	//	for i, _ := range tokenValue {
	//		for _, v := range value {
	//			if int(v.Uid) == tokenValue[i].Uid {
	//				tokenValue[i].NickName = v.NickName
	//				tokenValue[i].Email = v.Email
	//				tokenValue[i].Phone = v.Phone
	//				tokenValue[i].Status = v.Status
	//				break
	//			}
	//		}
	//		for _, tv := range tokenlist {
	//			if int(tv.Id) == tokenValue[i].TokenId {
	//				tokenValue[i].TokenName = tv.Name
	//				break
	//			}
	//		}
	//	}
	//	monerylist.Items = tokenValue
	//	this.Put(c, "list", monerylist)
	//	this.RespOK(c)
	//	return
	//} else {
	//	fmt.Println("1111111111111111111111111", err)
	//	list, err := new(models.MoneyRecord).GetMoneyListForDateOrType(req.Page, req.Rows, req.Type, req.Date)
	//	if err != nil {
	//		this.RespErr(c, err)
	//		return
	//	}
	//	uidList := make([]uint64, 0)
	//	Value, ok := list.Items.([]models.MoneyRecord)
	//	if !ok {
	//		this.RespErr(c, err)
	//		return
	//	}
	//	for _, v := range Value {
	//		uidList = append(uidList, uint64(v.Uid))
	//	}
	//	ulist, err := new(models.UserGroup).GetUserListForUid(uidList)
	//	if err != nil {
	//		this.RespErr(c, err)
	//		return
	//	}
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

	mlist.Items = list
	this.Put(c, "list", mlist)
	this.RespOK(c)
	return
}

func (this *TokenController) GetTokenCashList(c *gin.Context) {
	req := struct {
		Page int `form:"page" json:"page" binding:"required"`
		Rows int `form:"rows" json:"rows" `
		Id   int `form:"id" json:"id" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.ConfigQuenes).GetTokenCashList(req.Page, req.Rows, req.Id)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	fmt.Println("121212")
	this.Put(c, "list", list)
	this.RespOK(c, "成功")
	return
}

func (this *TokenController) GetTokenDetail(c *gin.Context) {
	req := struct {
		Uid      int `form:"uid" json:"uid" binding:"required"`
		Page     int `form:"page" json:"page" binding:"required"`
		Rows     int `form:"rows" json:"rows" `
		Token_id int `form:"token_id" json:"token_id"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	//bibi账户余额
	list, err := new(models.DetailToken).GetTokenDetailOfUid(req.Page, req.Rows, req.Uid, req.Token_id)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
	return
}

//bibi 账户统计表
func (this *TokenController) GetTokenBalance(c *gin.Context) {
	req := struct {
		Page   int    `form:"page" json:"page" binding:"required"`
		Rows   int    `form:"rows" json:"rows" `
		Status int    `form:"status" json:"status" `
		Search string `form:"search" json:"search" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	fmt.Printf("GetTokenBalance%#v\n", req)
	list, err := new(models.PersonalProperty).TotalUserBalance(req.Page, req.Rows, req.Status, req.Search)
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
		Page     int    `form:"page" json:"page" binding:"required"`
		Page_num int    `form:"rows" json:"rows" `
		Uid      int    `form:"uid" json:"uid" `
		Date     uint64 `form:"date" json:"date" `
		Name     string `form:"name" json:"name" binding:"required" ` //交易对
		Opt      int    `form:"opt" json:"opt" `                      //买卖方向
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	fmt.Println("start_t=", req.Date)
	list, err := new(models.Trade).GetTokenRecordList(req.Page, req.Page_num, req.Opt, req.Uid, req.Date, req.Name)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
	return
}

//币币挂单列表
func (this *TokenController) GetTokenOderList(c *gin.Context) {
	req := struct {
		Page      int    `form:"page" json:"page" binding:"required"`
		Rows      int    `form:"rows" json:"rows" `
		Uid       int    `form:"uid" json:"uid" `
		TradeId   string `form:"trade_id" json:"trade_id" ` //交易类型id 市价交易or 限价交易
		BeginTime int    `form:"bt" json:"bt"`
		EndTime   int    `form:"et" json:"et"`
		Symbol    string `form:"symbol" json:"symbol"  binding:"required" ` //交易对
		AdId      int    `form:"ad_id" json:"ad_id" `                       //买卖方向
		Status    int    `form:"status" json:"staus" `                      //订单状态
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.EntrustDetail).GetTokenOrderList(req.Page, req.Rows, req.AdId, req.Status, req.BeginTime, req.EndTime, req.Uid, req.Symbol, req.TradeId)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
	return
}

// 日划转汇总列表
func (t *TokenController) ListTransferDailySheet(ctx *gin.Context) {
	// 获取参数
	page, err := t.GetInt(ctx, "page", 1)
	if err != nil {
		t.RespErr(ctx, "参数page格式错误")
		return
	}
	rows, err := t.GetInt(ctx, "rows", 10)
	if err != nil {
		t.RespErr(ctx, "参数rows格式错误")
		return
	}

	// 筛选
	filter := make(map[string]string)
	if v := t.GetString(ctx, "type"); v != "" {
		filter["type"] = v
	}
	if v := t.GetString(ctx, "token_id"); v != "" {
		filter["token_id"] = v
	}
	if v := t.GetString(ctx, "date_begin"); v != "" {
		if matched, err := regexp.MatchString(constant.REGE_PATTERN_DATE, v); err != nil || !matched {
			t.RespErr(ctx, "参数date_begin格式错误")
			return
		}

		filter["date_begin"] = v
	}
	if v := t.GetString(ctx, "date_end"); v != "" {
		if matched, err := regexp.MatchString(constant.REGE_PATTERN_DATE, v); err != nil || !matched {
			t.RespErr(ctx, "参数date_end格式错误")
			return
		}

		filter["date_end"] = v
	}

	// 调用model
	modelList, list, err := new(models.TransferDailySheet).List(page, rows, filter)
	if err != nil {
		t.RespErr(ctx, err)
		return
	}

	// 整理数据
	type newItem struct {
		Id        int32   `json:"id"`
		TokenId   int32   `json:"token_id"`
		TokenName string  `json:"token_name"`
		Type      int8    `json:"type"`
		Num       float64 `json:"num"`
		Date      string  `json:"date"`
	}
	newItems := make([]newItem, len(list))
	for k, v := range list {
		newItems[k] = newItem{
			Id:        v.Id,
			TokenId:   v.TokenId,
			TokenName: v.TokenName,
			Type:      v.Type,
			Num:       convert.Int64ToFloat64By8Bit(v.Num),
			Date:      v.Date,
		}
	}
	modelList.Items = newItems

	// 设置返回数据
	t.Put(ctx, "list", modelList)

	// 返回
	t.RespOK(ctx)
	return
}
