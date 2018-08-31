package controller

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"admin/app/models"
	"admin/constant"
	"admin/utils"
	"admin/utils/convert"

	"admin/apis"

	"github.com/gin-gonic/gin"
)

type TokenController struct {
	BaseController
}

func (this *TokenController) Router(r *gin.Engine) {
	g := r.Group("/token")
	{
		g.GET("/list", this.GetTokenOderList)               //bibi p4-1-0 币币委托管理 挂单信息
		g.POST("/evacuate_order", this.EvacuateOder)        // p4-1-0 币币委托管理 挂单信息 撤单
		g.GET("/record_list", this.GetRecordList)           //bibi p4-1-1 成交记录
		g.GET("/export_record_list", this.ExportRecordList) //bibi p4-1-1 成交记录 导出
		g.GET("/total_balance", this.GetTokenBalance)       //bibi p2-3-2币币账户统计列表
		g.GET("/user_token_detail", this.GetTokenDetail)    //p2-3-2-1查看币币账户资产
		g.GET("/token_cash_list", this.GetTokenCashList)    //p4-1-2币兑管理
		g.POST("/delete_cash", this.DeleteCash)             //删除币兑
		g.GET("/modify_cash", this.ModifyCash)              //修改币兑
		g.POST("/add_cash", this.AddCash)                   //添加币兑
		g.GET("/change_detail", this.ChangeDetail)          //p2-3-4币币账户变更详情
		g.GET("/fee_list", this.GetFeeInfoList)             //p5-1-0-1币币交易手续费明细
		g.GET("/add_take_list", this.GetAddTakeList)        //p5-1-1-1提币手续费明细
		//g.GET("/total_trade", this.GetTradeTotalList)       //p5-1-0-1币币交易手续费明細 此函数已经重写
		//提币 充币管理
		g.GET("/in_token_list", this.GetTokenInList)         //充币
		g.GET("/out_token_list", this.GetTokenOutList)       //提币
		g.POST("/opt_token_pass", this.OptTakeTokenPass)     // 审核通过
		g.POST("/opt_token_failed", this.OptTakeTokenFailed) // 审核撤销
		//日提币充币汇总
		g.GET("/total_token_list", this.GetTotalTokenList) //
		g.GET("/total_token_info_list", this.GetTotalTokenInfoList)


		//币币交易手续费汇总
		g.GET("/token_order_fee_total", this.GetOderFeeTotalList)
		//提币手续费汇总表
		g.GET("/token_inout_daily_sheet", this.GetTokenInOutDailySheetList)

		//划转
		g.GET("/list_transfer_daily_sheet", this.ListTransferDailySheet)
		g.GET("/list_transfer", this.ListTransfer)
		//后台充值
		g.POST("/backstage_put", this.BackstagePut)
		//平台充值总表
		g.GET("/platform_all", this.PlatformAllSheet)
		//每天平台内充币详细信息
		g.GET("/platform_day", this.PlatformDay)

		//仪表盘
		g.GET("/trade_count", this.TradeCount)
		g.GET("/trade_total", this.TradeTotal)
		g.GET("/fee_total", this.FeeTotal)
		g.GET("/fee_trend", this.FeeTrend)
		g.GET("/trade_trend", this.TradeTrend)
		g.GET("/num_trend", this.NumTrend)
	}
}

func (this *TokenController) PlatformDay(c *gin.Context) {
	req := struct {
		Page int    `form:"page" json:"page" binding:"required"`
		Rows int    `form:"rows" json:"rows" `
		Date uint64 `form:"date" json:"date" binding:"required" `
		Tid  int    `form:"tid" json:"tid" binding:"required"`
		Uid  int    `form:"uid" json:"uid"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.MoneyRecord).GetPlatForTokenOfDay(req.Page, req.Rows, req.Uid, req.Tid, req.Date)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
	return
}

func (this *TokenController) PlatformAllSheet(c *gin.Context) {
	req := struct {
		Page int    `form:"page" json:"page" binding:"required"`
		Rows int    `form:"rows" json:"rows" `
		Bt   uint64 `form:"bt" json:"bt"` //日期
		Et   uint64 `form:"et" json:"et"`
		Tid  uint64 `form:"tid" json:"tid"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.MoneyRecord).GetPlatformAll(req.Page, req.Rows, req.Tid, req.Bt, req.Et)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
	return
}

func (this *TokenController) BackstagePut(c *gin.Context) {
	req := struct {
		Comment string  `form:"comment" json:"comment" binding:"required"`
		Count   float64 `form:"count" json:"count" binding:"required" `
		TokenId int     `form:"tid" json:"tid" binding:"required"`
		Uid     int     `form:"uid" json:"uid" binding:"required" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	err = new(models.MoneyRecord).BackstagePut(req.Count, req.Uid, req.TokenId, req.Comment)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.RespOK(c)
	return
}

//提币手续费汇总表
func (this *TokenController) GetTokenInOutDailySheetList(c *gin.Context) {
	req := struct {
		Page    int    `form:"page" json:"page" binding:"required"`
		Rows    int    `form:"rows" json:"rows" `
		Bt      string `form:"bt" json:"bt"` //日期
		Et      string `form:"et" json:"et"`
		TokenId int    `form:"tid" json:"tid"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.TokenInoutDailySheet).GetInOutDailySheetList(req.Page, req.Rows, req.TokenId, req.Bt, req.Et)
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
		Date uint64 `form:"date" json:"date"` //日期
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, data, err := new(models.TokenDailySheet).GetDailySheetList(req.Page, req.Rows, req.Date)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "data", data)
	this.Put(c, "list", list)

	this.RespOK(c)
	return
}

//提 充 币汇总
func (this *TokenController) GetTotalTokenList(c *gin.Context) {
	req := struct {
		Opt     int    `form:"opt" json:"opt" binding:"required"` //1充 2 提币
		Page    int    `form:"page" json:"page" binding:"required"`
		Rows    int    `form:"rows" json:"rows" `
		TokenId int    `form:"tokenId" json:"tokenId" ` //货币id
		Bt      string `form:"bt" json:"bt"`            //筛选开始日期
		Et      string `form:"et" json:"et"`            //筛选结束日期
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	//list, err := new(models.TokenInout).GetTotalList(req.Page, req.Rows, req.TokenId, req.Opt, req.Date)
	if req.Opt == 2 { //提币
		list, err := new(models.TokenInoutDailySheet).DayOutDailySheet(req.Page, req.Rows, req.TokenId, req.Bt, req.Et)
		if err != nil {
			this.RespErr(c, err)
			return
		}
		this.Put(c, "list", list)
		this.RespOK(c)
		return
	}
	if req.Opt == 1 { //充币
		list, err := new(models.TokenInoutDailySheet).DayPutDailySheet(req.Page, req.Rows, req.TokenId, req.Bt, req.Et)
		if err != nil {
			this.RespErr(c, err)
			return
		}
		this.Put(c, "list", list)
		this.RespOK(c)
		return
	}
	return
}

func (this *TokenController) GetTotalTokenInfoList(c *gin.Context) {
	req := struct {
		Opt     int    `form:"opt" json:"opt" binding:"required"` //1充 2 提币
		Page    int    `form:"page" json:"page" binding:"required"`
		Rows    int    `form:"rows" json:"rows" `
		TokenId int    `form:"tokenId" json:"tokenId" ` //货币id
		Search  string `form:"search" json:"search"`    //输入uid 搜索
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.TokenInout).GetTotalInfoList(req.Page, req.Rows, req.TokenId, req.Opt, req.Search)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
}

//充提币管理列表
func (this *TokenController) GetTokenInList(c *gin.Context) {
	this.getTokenList(c)
	return
}
func (this *TokenController) GetTokenOutList(c *gin.Context) {
	this.getTokenList(c)
	return
}

func (this *TokenController) getTokenList(c *gin.Context) {
	req := struct {
		Page    int    `form:"page" json:"page" binding:"required"`
		Rows    int    `form:"rows" json:"rows" `
		Ustatus int    `form:"ustatus" json:"ustatus" ` // 用户状态
		Status  int    `form:"status" json:"status" `   //提币状态
		TokenId int    `form:"tokenId" json:"tokenId" ` //货币id
		Opt     int    `form:"opt" json:"opt"  `        //操作方向
		Search  string `form:"search" json:"search" `   //筛选
		//Date    string `form:"date",json:"date"`        //日期
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.TokenInoutGroup).GetTokenInList(req.Page, req.Rows, req.Ustatus, req.Status, req.TokenId, req.Opt, req.Search)
	if err != nil {
		utils.AdminLog.Println(err.Error())
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", list)
	this.RespOK(c)
	return
}

// 提币审核
func (this *TokenController) OptTakeTokenPass(c *gin.Context) {
	this.optTakeToken(c)
	return
}
func (this *TokenController) OptTakeTokenFailed(c *gin.Context) {
	this.optTakeToken(c)
	return
}

func (this *TokenController) optTakeToken(c *gin.Context) {
	req := struct {
		Id     int `form:"id" json:"id" binding:"required"`
		Status int `form:"status" json:"status" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	err = new(models.TokenInout).OptTakeToken(req.Id, req.Status)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.RespOK(c)
	return
}

//func (this *TokenController) GetTradeTotalList(c *gin.Context) {
//	req := struct {
//		Page int    `form:"page" json:"page" binding:"required"`
//		Rows int    `form:"rows" json:"rows" `
//		Date uint64 `form:"date",json:"date"` //日期
//	}{}
//	err := c.ShouldBind(&req)
//	if err != nil {
//		utils.AdminLog.Errorf(err.Error())
//		this.RespErr(c, err)
//		return
//	}
//	list, err := new(models.Trade).TotalTotalTradeList(req.Page, req.Rows, req.Date)
//	if err != nil {
//		utils.AdminLog.Println(err.Error())
//		this.RespErr(c, err)
//		return
//	}
//	this.Put(c, "list", list)
//	this.RespOK(c)
//	return
//}

//p5-1-1-1提币手续费明细
func (this *TokenController) GetAddTakeList(c *gin.Context) {
	req := struct {
		Page    int `form:"page" json:"page" binding:"required"`
		Rows    int `form:"rows" json:"rows" `
		TokenId int `form:"token_id" json:"token_id" binding:"required"` //币种
		Uid     int `form:"uid" json:"uid" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.TokenInout).GetAddTakeList(req.Page, req.Rows, req.TokenId, req.Uid)
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
		Tid    int    `form:"tid" json:"tid" binding:"required"`
		Bt     uint64 `form:"bt" json:"bt"`
		Et     uint64 `form:"et" json:"et"`
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
	//这里分开
	//type =3 注册奖励
	//type=4 邀请奖励
	//因为表不同所以在此分开查询
	if req.Type == 3 || req.Type == 4 {
		fmt.Println("hehe")
		list, err := new(models.FrozenHistory).GetFrozenHistory(req.Page, req.Rows, req.Type, req.Tid, req.Status, req.Bt, req.Et, req.Search)
		if err != nil {
			this.RespErr(c, err)
			return
		}
		uidList := make([]int64, 0)
		value, ok := list.Items.([]models.FrozenHistoryGroup)
		if !ok {
			this.RespErr(c, errors.New("assert failed"))
			return
		}
		for _, v := range value {
			uidList = append(uidList, v.Uid)
		}
		balanceList, err := new(models.UserToken).GetBalanceForUid(uidList)
		if err != nil {
			this.RespErr(c, err)
			return
		}
		fmt.Println("value=", len(value), "uid=", len(balanceList))
		for i, v := range value {
			for _, bv := range balanceList {
				if v.Uid == int64(bv.Uid) {
					value[i].NumTrue = convert.Int64ToFloat64By8Bit(v.Num)
					//fmt.Println("前=",v.Num,"后=",value[i].NumTrue)
					value[i].SurplusTrue = convert.Int64ToFloat64By8Bit(bv.Balance)
					break
				}
			}
			for _, vt := range tokenlist {
				if vt.Id == uint32(v.TokenId) {
					value[i].TokenName = vt.Mark
					break
				}
			}

		}
		list.Items = value
		this.Put(c, "list", list)
		this.RespOK(c)
		return

	} else {
		mlist, err := new(models.MoneyRecord).GetMoneyListForDateOrType(req.Page, req.Rows, req.Type, req.Status, req.Tid, req.Bt, req.Et, req.Search)
		if err != nil {
			this.RespErr(c, err)
			return
		}

		list, Ok := mlist.Items.([]models.MoneyRecordGroup)
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

		mlist.Items = list
		this.Put(c, "list", mlist)
		this.RespOK(c)
		return
	}

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
	this.RespOK(c)
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
		Tid    int    `form:"tid" json:"tid"`
		Range  string `form:"range" json:"range" `
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
	list, err := new(models.PersonalProperty).TotalUserBalance(req.Page, req.Rows, req.Status, req.Tid, req.Range, req.Search)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	uidList := make([]uint64, 0)
	value, OK := list.Items.([]models.PersonalProperty)
	if !OK {
		this.RespErr(c, errors.New("assert failed"))
		return
	}
	for _, value := range value {
		uidList = append(uidList, uint64(value.Uid))
	}
	//fmt.Println("uid", uidList)

	tokenList, err := new(apis.VendorApi).GetCny(uidList, 1)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		this.RespErr(c, err.Error())
		return
	}
	for i, v := range value {
		for _, vc := range tokenList {
			if vc.Uid == uint64(v.Uid) {
				value[i].AmountTo = vc.TotalCny
				break
			}
		}
	}

	this.Put(c, "list", list)
	this.RespOK(c)
	return
}

//bibi 成交记录
func (this *TokenController) GetRecordList(c *gin.Context) {
	this.getRecordList(c)
}

//bibi 成交记录
func (this *TokenController) ExportRecordList(c *gin.Context) {
	this.getRecordList(c)
}
func (this *TokenController) getRecordList(c *gin.Context) {
	req := struct {
		Page int    `form:"page" json:"page" binding:"required"`
		Rows int    `form:"rows" json:"rows" `
		Uid  int    `form:"uid" json:"uid" `
		Bt   uint64 `form:"bt" json:"bt" `
		Et   uint64 `form:"et" json:"et" `
		Name string `form:"name" json:"name" ` //交易对
		Opt  int    `form:"opt" json:"opt" `   //买卖方向
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
	list, err := new(models.Trade).GetTokenRecordList(req.Page, req.Rows, req.Opt, req.Uid, req.Bt, req.Et, req.Name)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	value, Ok := list.Items.([]models.TradeReturn)
	if !Ok {
		this.RespErr(c, errors.New("assert type UserGroup failed!!"))
		return
	}
	for i, v := range value {
		for _, vt := range tokenlist {
			if vt.Id == uint32(v.TokenAdmissionId) {
				value[i].TokenName = vt.Mark
				break
			}
		}
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
		Tp        int    `form:"ty" json:"ty" ` //交易类型id 市价交易or 限价交易
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
	list, err := new(models.EntrustDetail).GetTokenOrderList(req.Page, req.Rows, req.AdId, req.Status, req.BeginTime, req.EndTime, req.Uid, req.Tp, req.Symbol)
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

// 划转明细
func (t *TokenController) ListTransfer(ctx *gin.Context) {
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
	filter := map[string]interface{}{
		"transfer": true,
	}
	if v := t.GetString(ctx, "uid"); v != "" {
		filter["uid"] = v
	}
	if v := t.GetString(ctx, "token_id"); v != "" {
		filter["token_id"] = v
	}
	if v := t.GetString(ctx, "transfer_date"); v != "" {
		filter["transfer_date"] = v
	}

	// 调用model
	modelList, list, err := new(models.MoneyRecord).List(page, rows, filter)

	// 重组data
	type item struct {
		Id           int64  `json:"id"`
		Uid          int    `json:"uid"`
		TokenId      int32  `json:"token_id"`
		TokenName    string `json:"token_name"`
		Type         int8   `json:"type"`
		Num          string `json:"num"`
		TransferTime int64  `json:"transfer_time"`
	}

	newItems := make([]*item, len(list))
	for k, v := range list {
		newItems[k] = &item{
			Id:           v.Id,
			Uid:          v.Uid,
			TokenId:      int32(v.TokenId),
			TokenName:    v.TokenName,
			Type:         int8(v.Type),
			Num:          convert.Int64ToStringBy8Bit(v.Num),
			TransferTime: v.TransferTime,
		}
	}
	modelList.Items = newItems

	// 设置返回数据
	t.Put(ctx, "list", modelList)

	// 返回
	t.RespOK(ctx)
	return
}

// 交易次数合计
func (t *TokenController) TradeCount(ctx *gin.Context) {
	// 调用model
	//1. 币币交易
	tokenTradeTotal, err := new(models.Trade).TradeTotal()
	if err != nil {
		t.RespErr(ctx, err)
		return
	}

	//2. 法币交易
	currencyTradeTotal, err := new(models.UserCurrencyHistory).TradeTotal()
	if err != nil {
		t.RespErr(ctx, err)
		return
	}

	//3. 提币交易
	inOutTradeTotal, err := new(models.TokenInout).TradeTotal()
	if err != nil {
		t.RespErr(ctx, err)
		return
	}

	// 整理数据
	total := tokenTradeTotal.TotalTime + currencyTradeTotal.TotalTime + inOutTradeTotal.TotalTime

	// 相较昨日、上周今日
	todayTotal := tokenTradeTotal.TodayTotalTime + currencyTradeTotal.TodayTotalTime + inOutTradeTotal.TodayTotalTime
	yesterdayTotal := tokenTradeTotal.YesterdayTotalTime + currencyTradeTotal.YesterdayTotalTime + inOutTradeTotal.YesterdayTotalTime
	lastWeekDayTotal := tokenTradeTotal.LastWeekDayTotalTime + currencyTradeTotal.LastWeekDayTotalTime + inOutTradeTotal.LastWeekDayTotalTime

	// 昨日
	upDay := todayTotal - yesterdayTotal

	// 上周今日
	upWeek := todayTotal - lastWeekDayTotal

	// 设置返回数据
	t.Put(ctx, "total_count", total)
	t.Put(ctx, "upday", upDay)
	t.Put(ctx, "upweek", upWeek)

	// 返回
	t.RespOK(ctx)
	return
}

// 交易数量合计
func (t *TokenController) TradeTotal(ctx *gin.Context) {
	// 调用model
	//1. 币币交易
	tokenTradeTotal, err := new(models.Trade).TradeTotal()
	if err != nil {
		t.RespErr(ctx, err)
		return
	}

	//2. 法币交易
	currencyTradeTotal, err := new(models.UserCurrencyHistory).TradeTotal()
	if err != nil {
		t.RespErr(ctx, err)
		return
	}

	//3. 提币交易
	inOutTradeTotal, err := new(models.TokenInout).TradeTotal()
	if err != nil {
		t.RespErr(ctx, err)
		return
	}

	// 整理数据
	total, _ := convert.StringAddString(tokenTradeTotal.TotalNum, currencyTradeTotal.TotalNum, inOutTradeTotal.TotalNum)
	total, _ = convert.StringTo8Bit(total)

	// 相较昨日、上周今日
	todayTotal, _ := convert.StringAddString(tokenTradeTotal.TodayTotalNum, currencyTradeTotal.TodayTotalNum, inOutTradeTotal.TodayTotalNum)
	yesterdayTotal, _ := convert.StringAddString(tokenTradeTotal.YesterdayTotalNum, currencyTradeTotal.YesterdayTotalNum, inOutTradeTotal.YesterdayTotalNum)
	lastWeekDayTotal, _ := convert.StringAddString(tokenTradeTotal.LastWeekDayTotalNum, currencyTradeTotal.LastWeekDayTotalNum, inOutTradeTotal.LastWeekDayTotalNum)

	// 昨日
	upDay, _ := convert.StringSubString(todayTotal, yesterdayTotal)
	upDay, _ = convert.StringTo8Bit(upDay)

	// 上周今日
	upWeek, _ := convert.StringSubString(todayTotal, lastWeekDayTotal)
	upWeek, _ = convert.StringTo8Bit(upWeek)

	// 设置返回数据
	t.Put(ctx, "total_num", total)
	t.Put(ctx, "upday", upDay)
	t.Put(ctx, "upweek", upWeek)

	// 返回
	t.RespOK(ctx)
	return
}

// 手续费合计
func (t *TokenController) FeeTotal(ctx *gin.Context) {
	// 调用model
	//1. 币币交易
	tokenTradeTotal, err := new(models.Trade).TradeTotal()
	if err != nil {
		t.RespErr(ctx, err)
		return
	}

	//2. 法币交易
	currencyTradeTotal, err := new(models.UserCurrencyHistory).TradeTotal()
	if err != nil {
		t.RespErr(ctx, err)
		return
	}

	//3. 提币交易
	inOutTradeTotal, err := new(models.TokenInout).TradeTotal()
	if err != nil {
		t.RespErr(ctx, err)
		return
	}

	// 整理数据
	total, _ := convert.StringAddString(tokenTradeTotal.TotalFee, currencyTradeTotal.TotalFee, inOutTradeTotal.TotalFee)
	total, _ = convert.StringTo8Bit(total)

	// 相较昨日、上周今日
	todayTotal, _ := convert.StringAddString(tokenTradeTotal.TodayTotalFee, currencyTradeTotal.TodayTotalFee, inOutTradeTotal.TodayTotalFee)
	yesterdayTotal, _ := convert.StringAddString(tokenTradeTotal.YesterdayTotalFee, currencyTradeTotal.YesterdayTotalFee, inOutTradeTotal.YesterdayTotalFee)
	lastWeekDayTotal, _ := convert.StringAddString(tokenTradeTotal.LastWeekDayTotalFee, currencyTradeTotal.LastWeekDayTotalFee, inOutTradeTotal.LastWeekDayTotalFee)

	// 昨日
	upDay, _ := convert.StringSubString(todayTotal, yesterdayTotal)
	upDay, _ = convert.StringTo8Bit(upDay)

	// 上周今日
	upWeek, _ := convert.StringSubString(todayTotal, lastWeekDayTotal)
	upWeek, _ = convert.StringTo8Bit(upWeek)

	// 设置返回数据
	t.Put(ctx, "total_fee", total)
	t.Put(ctx, "upday", upDay)
	t.Put(ctx, "upweek", upWeek)

	// 返回
	t.RespOK(ctx)
	return
}

// 手续费走势
func (t *TokenController) FeeTrend(ctx *gin.Context) {
	loc, _ := time.LoadLocation("Local")

	// 筛选
	tokenFilter := make(map[string]interface{})
	currencyFilter := make(map[string]interface{})
	inOutFilter := make(map[string]interface{})
	if v := t.GetString(ctx, "token_id"); v != "" {
		tokenFilter["token_id"] = v
		currencyFilter["token_id"] = v
		inOutFilter["token_id"] = v
	}
	if v := t.GetString(ctx, "date_begin"); v != "" {
		if matched, err := regexp.Match(constant.REGE_PATTERN_DATE, []byte(v)); err != nil || !matched {
			t.RespErr(ctx, "参数date_begin格式错误")
			return
		}
		dateTime, _ := time.ParseInLocation(utils.LAYOUT_DATE, v, loc)

		tokenFilter["date_begin"] = dateTime.Unix()
		currencyFilter["date_begin"] = v
		inOutFilter["date_begin"] = v
	}
	if v := t.GetString(ctx, "date_end"); v != "" {
		if matched, err := regexp.Match(constant.REGE_PATTERN_DATE, []byte(v)); err != nil || !matched {
			t.RespErr(ctx, "参数date_end格式错误")
			return
		}
		dateTime, _ := time.ParseInLocation(utils.LAYOUT_DATE, v, loc)

		tokenFilter["date_end"] = dateTime.Unix()
		currencyFilter["date_end"] = v
		inOutFilter["date_end"] = v
	}

	// 调用model
	//1. 币币买入、卖出手续费
	tokenList, err := new(models.TokenDailySheet).TradeTrendList(tokenFilter)
	if err != nil {
		t.RespErr(ctx, err)
		return
	}

	//2. 法币买入、卖出手续费
	currencyList, err := new(models.CurrencyDailySheet).TradeTrendList(currencyFilter)
	if err != nil {
		t.RespErr(ctx, err)
		return
	}

	//3. 提币手续费
	inOutList, err := new(models.TokenInoutDailySheet).InOutTrendList(inOutFilter)
	if err != nil {
		t.RespErr(ctx, err)
		return
	}

	// 转换成以日期作为key的map
	tokenListMap := make(map[string]*models.TokenTradeTrend)
	for _, v := range tokenList {
		tokenListMap[time.Unix(v.Date, 0).Format("0102")] = v
	}

	currencyListMap := make(map[string]*models.CurrencyTradeTrend)
	for _, v := range currencyList {
		datetime, _ := time.ParseInLocation(utils.LAYOUT_DATE_TIME, v.Date, loc)
		currencyListMap[datetime.Format("0102")] = v
	}

	inOutListMap := make(map[string]*models.InOutTrend)
	for _, v := range inOutList {
		datetime, _ := time.ParseInLocation(utils.LAYOUT_DATE_TIME, v.Date, loc)
		inOutListMap[datetime.Format("0102")] = v
	}

	// 组装数据
	listLen := len(tokenList)
	x := make([]string, listLen)
	yBuy := make([]string, listLen)
	ySell := make([]string, listLen)
	yOut := make([]string, listLen)

	var allBuyTotal = "0"
	var allSellTotal = "0"
	var allOutTotal = "0"
	var allTotal = "0"
	for k, v := range tokenList {
		date := time.Unix(v.Date, 0).Format("0102")

		var buy = "0"
		var sell = "0"
		var out = "0"
		if md, ok := tokenListMap[date]; ok {
			buy, _ = convert.StringAddString(buy, md.FeeBuyTotal)
			sell, _ = convert.StringAddString(sell, md.FeeSellTotal)
		}
		if md, ok := currencyListMap[date]; ok {
			buy, _ = convert.StringAddString(buy, md.FeeBuyTotal)
			sell, _ = convert.StringAddString(sell, md.FeeSellTotal)
		}
		if md, ok := inOutListMap[date]; ok {
			out, _ = convert.StringAddString(out, md.FeeTotal)
		}

		x[k] = date

		yBuy[k], _ = convert.StringTo8Bit(buy)
		ySell[k], _ = convert.StringTo8Bit(sell)
		yOut[k], _ = convert.StringTo8Bit(out)

		allBuyTotal, _ = convert.StringAddString(allBuyTotal, buy)
		allSellTotal, _ = convert.StringAddString(allSellTotal, sell)
		allOutTotal, _ = convert.StringAddString(allOutTotal, out)
		allTotal, _ = convert.StringAddString(allTotal, buy, sell, out)
	}

	allBuyTotalFloat, _ := convert.StringTo8Bit(allBuyTotal) // 转成float
	allSellTotalFloat, _ := convert.StringTo8Bit(allSellTotal)
	allOutTotalFloat, _ := convert.StringTo8Bit(allOutTotal)
	allTotalFloat, _ := convert.StringTo8Bit(allTotal)

	// 设置返回数据
	t.Put(ctx, "x", x)
	t.Put(ctx, "y_buy", yBuy)
	t.Put(ctx, "y_sell", ySell)
	t.Put(ctx, "y_out", yOut)
	t.Put(ctx, "all_buy_total", allBuyTotalFloat)
	t.Put(ctx, "all_sell_total", allSellTotalFloat)
	t.Put(ctx, "all_out_total", allOutTotalFloat)
	t.Put(ctx, "all_total", allTotalFloat)

	// 返回
	t.RespOK(ctx)
	return
}

// 币币交易走势
func (t *TokenController) TradeTrend(ctx *gin.Context) {
	loc, _ := time.LoadLocation("Local")

	// 筛选
	filter := make(map[string]interface{})
	if v := t.GetString(ctx, "token_id"); v != "" {
		filter["token_id"] = v
	}
	if v := t.GetString(ctx, "date_begin"); v != "" {
		if matched, err := regexp.Match(constant.REGE_PATTERN_DATE, []byte(v)); err != nil || !matched {
			t.RespErr(ctx, "参数date_begin格式错误")
			return
		}
		dateTime, _ := time.ParseInLocation(utils.LAYOUT_DATE, v, loc)

		filter["date_begin"] = dateTime.Unix()
	}
	if v := t.GetString(ctx, "date_end"); v != "" {
		if matched, err := regexp.Match(constant.REGE_PATTERN_DATE, []byte(v)); err != nil || !matched {
			t.RespErr(ctx, "参数date_end格式错误")
			return
		}
		dateTime, _ := time.ParseInLocation(utils.LAYOUT_DATE, v, loc)

		filter["date_end"] = dateTime.Unix()
	}

	// 调用model
	list, err := new(models.TokenDailySheet).TradeTrendList(filter)
	if err != nil {
		t.RespErr(ctx, err)
		return
	}

	// 组装数据
	listLen := len(list)
	x := make([]string, listLen)
	yBuy := make([]string, listLen)
	ySell := make([]string, listLen)

	allBuyTotal := "0"  // 买入总计
	allSellTotal := "0" // 卖出总计
	for k, v := range list {
		x[k] = time.Unix(v.Date, 0).Format("0102")
		yBuy[k], _ = convert.StringTo8Bit(v.BuyTotal)
		ySell[k], _ = convert.StringTo8Bit(v.SellTotal)

		allBuyTotal, _ = convert.StringAddString(allBuyTotal, v.BuyTotal)
		allSellTotal, _ = convert.StringAddString(allSellTotal, v.SellTotal)
	}
	allBuyTotalFloat, _ := convert.StringTo8Bit(allBuyTotal)   // 转成float
	allSellTotalFloat, _ := convert.StringTo8Bit(allSellTotal) // 转成float

	// 设置返回数据
	t.Put(ctx, "x", x)
	t.Put(ctx, "y_buy", yBuy)
	t.Put(ctx, "y_sell", ySell)
	t.Put(ctx, "all_buy_total", allBuyTotalFloat)
	t.Put(ctx, "all_sell_total", allSellTotalFloat)

	// 返回
	t.RespOK(ctx)
	return
}

// 币种数量走势
func (t *TokenController) NumTrend(ctx *gin.Context) {
	// 筛选
	filter := make(map[string]interface{})

	tokenId, err := t.GetInt(ctx, "token_id")
	if err != nil || tokenId < 0 {
		t.RespErr(ctx, "参数token_id格式错误")
		return
	}
	filter["token_id"] = tokenId

	if v := t.GetString(ctx, "date_begin"); v != "" {
		if matched, err := regexp.Match(constant.REGE_PATTERN_DATE, []byte(v)); err != nil || !matched {
			t.RespErr(ctx, "参数date_begin格式错误")
			return
		}

		filter["date_begin"] = v
	}
	if v := t.GetString(ctx, "date_end"); v != "" {
		if matched, err := regexp.Match(constant.REGE_PATTERN_DATE, []byte(v)); err != nil || !matched {
			t.RespErr(ctx, "参数date_end格式错误")
			return
		}

		filter["date_end"] = v
	}

	// 调用model
	list, err := new(models.TokensDailySheet).NumTrend(filter)
	if err != nil {
		t.RespErr(ctx, err)
		return
	}

	// 组装数据
	listLen := len(list)
	x := make([]string, listLen)
	y := make([]string, listLen)

	var allTotal = "0"
	loc, _ := time.LoadLocation("Local")
	for k, v := range list {
		datetime, _ := time.ParseInLocation(utils.LAYOUT_DATE, v.Date, loc)
		total, _ := convert.StringAddString(v.TokenTotal, v.CurrencyTotal)

		x[k] = datetime.Format("0102")
		y[k], _ = convert.StringTo8Bit(total)

		allTotal, _ = convert.StringAddString(allTotal, v.TokenTotal, v.CurrencyTotal)
	}
	allTotalFloat, _ := convert.StringTo8Bit(allTotal)

	// 设置返回数据
	t.Put(ctx, "x", x)
	t.Put(ctx, "y", y)
	t.Put(ctx, "all_total", allTotalFloat)

	// 返回
	t.RespOK(ctx)
	return
}
