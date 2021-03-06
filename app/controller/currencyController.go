package controller

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"admin/apis"
	"admin/app/models"
	"admin/constant"
	"admin/utils"
	"admin/utils/convert"

	"github.com/gin-gonic/gin"
)

type CurrencyController struct {
	BaseController
}

func (this *CurrencyController) Router(r *gin.Engine) {
	g := r.Group("/currency")
	{
		g.GET("/list", this.GetTradeList)                                  //p4-2-0法币挂单管理
		g.GET("/export_list", this.ExportTradeList)                        //p4-2-0法币挂单管理
		g.POST("/down_trade_order", this.DownTradeAds)                     //p4-2-0法币挂单管理 下架交易单
		g.GET("/tokens", this.GetTokensList)                               //获取 所有数据货币的名称及货币Id
		g.GET("/order_list", this.GetOderList)                             //p4-2-1法币成交管理
		g.GET("/export_order_list", this.ExportOderList)                   //p4-2-1法币成交管理
		g.GET("/total_balance", this.GetTotalCurrencyBalance)              //p2-3-1法币账户统计列表
		g.GET("/export_total_balance", this.ExportTotalCurrencyBalance)    //p2-3-1法币账户统计列表
		g.GET("/user_detail", this.GetUserDetailList)                      //p2-3-1-2法币账户资产展示
		g.GET("/export_user_detail", this.ExportUserDetailList)            //p2-3-1-2法币账户资产展示
		g.GET("/user_buysell", this.GetBuySellList)                        //p2-3-1-1查看统计买入_卖出_划转
		g.GET("/export_user_buysell", this.ExportBuySellList)              //p2-3-1-1查看统计买入_卖出_划转
		g.GET("/total", this.Total)                                        //p2-3-0总财产列表
		g.GET("/export_total", this.ExportTotal)                           //p2-3-0总财产列表
		g.GET("/currency_change", this.GetCurrencyChangeHistory)           //p2-3-3法币账户变更详情
		g.GET("/export_currency_change", this.ExportCurrencyChangeHistory) //p2-3-3法币账户变更详情

		// 统计平台币总数
		g.GET("/total_coin", this.TotalCoin)
		g.GET("/export_total_coin", this.ExportTotalCoin)

		//g.GET("/")                                               //p2-3-0-0币数统计列表
		//划转到币币账户货币数量日统计 注释接口没有实现
		g.GET("/layoff_list", this.GetLayOffList)
		g.GET("/export_layoff_list", this.GetLayOffList)
		//法币成交管理 放行 取消
		g.GET("/revoke_currency", this.SetRevokeCurrency)      //撤单
		g.GET("/verify_pass_currency", this.SetCurrencyToPass) //审核通过

		g.GET("/trade_trend", this.TradeTrend)
	}
}

//撤单
func (cu *CurrencyController) SetRevokeCurrency(c *gin.Context) {
	req := struct {
		Id int64 `form:"id" json:"id" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}
	err = new(apis.VendorApi).CurrencyRevoke(req.Id)
	if err != nil {
		cu.RespErr(c, err)
		return
	}
	cu.RespOK(c)
	return
}

func (cu *CurrencyController) SetCurrencyToPass(c *gin.Context) {
	req := struct {
		Id int64 `form:"id" json:"id" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}
	err = new(apis.VendorApi).CurrencyVerityPass(req.Id)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}
	cu.RespOK(c)
	return
}

func (cu *CurrencyController) GetLayOffList(c *gin.Context) {
	return
}

func (cu *CurrencyController) GetCurrencyChangeHistory(c *gin.Context) {
	cu.currencyChangeHistory(c)
	return
}

func (cu *CurrencyController) ExportCurrencyChangeHistory(c *gin.Context) {
	cu.currencyChangeHistory(c)
	return
}

func (cu *CurrencyController) currencyChangeHistory(c *gin.Context) {
	req := struct {
		Page   int    `form:"page" json:"page" binding:"required"`
		Rows   int    `form:"rows" json:"rows" `
		Bt     string `form:"bt" json:"bt"`
		Et     string `form:"et" json:"et"`
		Search string `form:"search" json:"search" `             //搜索的内容
		Status int    `form:"status" json:"status" `             //用户账号状态
		Tid    int    `form:"tid" json:"tid" binding:"required"` //货币id
		Chtype int    `form:"type" json:"type"`                  // 买入 卖出 提币 充币 划转
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
	ulist, err := new(models.UserCurrencyHistory).GetListForUid(req.Page, req.Rows, req.Tid, req.Status, req.Chtype, req.Bt, req.Et, req.Search)
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

	ulist.Items = list
	cu.Put(c, "list", ulist)
	cu.RespOK(c)
	return
}

//总财产统计列表
func (cu *CurrencyController) Total(c *gin.Context) {
	cu.total(c)
	return
}
func (cu *CurrencyController) ExportTotal(c *gin.Context) {
	cu.total(c)
	return
}

func (cu *CurrencyController) total(c *gin.Context) {
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

	result, err := new(models.UserGroup).GetAllUser1(req.Page, req.Rows, req.Status, req.Search)
	if err != nil {
		cu.RespErr(c, err)
	}

	uidList := make([]uint64, 0)
	value, OK := result.Items.([]models.Total)
	if !OK {
		cu.RespErr(c, errors.New("assert failed"))
		return
	}
	for _, value := range value {
		uidList = append(uidList, uint64(value.Uid))
	}
	fmt.Println("uid", uidList)

	//总资产折合
	//币币账户折合
	//tk :=make(chan int,2)
	tokenList, err := new(apis.VendorApi).GetCny(uidList, 1)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		cu.RespErr(c, err.Error())
		return
	}
	//资产总折合
	//法币账户折合
	currencyList, err := new(apis.VendorApi).GetCny(uidList, 2)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		cu.RespErr(c, err)
		return
	}
	var totalCurrencyInt int64
	var totalTokenInt int64
	fmt.Println(currencyList)
	for i, v := range value {
		for _, vt := range tokenList {
			if vt.Uid == uint64(v.Uid) {
				value[i].LockTokenCNY = vt.FrozenCny
				value[i].TotalTokenCNY = vt.BalanceCny
				//totalInt =vt.TotalCny
				totalTokenInt = vt.FrozenCnyInt + vt.BalanceCnyInt
				break
			}
		}
		for _, vc := range currencyList {
			if vc.Uid == uint64(v.Uid) {
				value[i].LockCurrentCNY = vc.FrozenCny
				value[i].TotalCurrentCNY = vc.BalanceCny
				//temp, _ := strconv.ParseFloat(vc.TotalCny, 64)
				totalCurrencyInt = vc.BalanceCnyInt + vc.FrozenCnyInt
				break
			}
		}
		value[i].TotalCNY = convert.Int64ToStringBy8Bit(totalTokenInt + totalCurrencyInt)
		totalCurrencyInt = 0
		totalTokenInt = 0
	}
	result.Items = value
	cu.Put(c, "list", result)
	cu.RespOK(c)
	return
}

/*
	total coin
*/
func (cu *CurrencyController) TotalCoin(c *gin.Context) {
	cu.totalCoin(c)
	return
}

func (cu *CurrencyController) ExportTotalCoin(c *gin.Context) {
	cu.totalCoin(c)
	return
}

func (cu *CurrencyController) totalCoin(c *gin.Context) {
	fmt.Println("total coin ...")
	req := struct {
		Page    int `form:"page"       json:"page" binding:"required"`
		Rows    int `form:"rows"       json:"rows" `
		TokenId int `form:"token_id"   json:"token_id"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}
	var tokenList []models.CommonTokens
	var total int64
	if req.Page <= 1 {
		req.Page = 1
	}
	if req.Rows <= 0 {
		req.Rows = 10
	}
	if req.TokenId > 0 {
		tokenList, total, err = new(models.CommonTokens).GetTokenPage(req.Page, req.Rows, int32(req.TokenId))
	} else {
		tokenList, total, err = new(models.CommonTokens).GetTokenPage(req.Page, req.Rows, 0)
	}

	if err != nil {
		fmt.Println(err)
	}
	var tokenIdList []int32
	for _, tk := range tokenList {
		tokenIdList = append(tokenIdList, int32(tk.Id))
	}

	type TotalCoin struct {
		TokenId    int    `json:"token_id"`
		TokenName  string `json:"token_name"`
		TotalUser  int64  `json:"total_user"`
		TotalNum   string `json:"total_num"`
		AverageNum string `json:"average_num"`
		TotalFree  string `json:"total_free"`
	}

	var totalcoinList []TotalCoin

	// 统计场外总币数和持有币总人数
	currencyBalanceList, currencyUserCoin, err := new(models.UserCurrency).GetAllCurrencyCoin(tokenIdList)
	// 统计币币交易总币数和对应币总持有人数
	tokenBalanceList, tokenUserCoin, err := new(models.UserToken).GetAllTokenCoin(tokenIdList)

	// 统计手续费
	addFreeList, delFreeList, err := new(models.TokenFreeHistory).GetFreeByTokenIds(tokenIdList)

	for _, tk := range tokenList {
		var totalFree string
		var totalnum string
		var totaluser int64
		var tmp TotalCoin
		tmp.TokenId = int(tk.Id)
		tmp.TokenName = tk.Mark

		// 手续费
		for _, addfree := range addFreeList {
			if addfree.TokenId == int32(tk.Id) {
				totalFree, _ = convert.StringAddString(totalFree, addfree.TotalAddFree)
			}
		}
		for _, delfree := range delFreeList {
			if delfree.TokenId == int32(tk.Id) {
				totalFree, _ = convert.StringSubString(totalFree, delfree.TotalDelFree)
			}
		}

		//

		for _, tokenBalance := range tokenBalanceList {
			if tokenBalance.TokenId == int32(tk.Id) {
				fmt.Println("bibi : ", tk.Mark, " you :", tokenBalance.TotalBalanceStr, tokenBalance.TotalFrozenStr)
				totalnum, _ = convert.StringAddString(totalnum, tokenBalance.TotalBalanceStr)
				totalnum, _ = convert.StringAddString(totalnum, tokenBalance.TotalFrozenStr)
			}
		}
		for _, currencyBalance := range currencyBalanceList {
			if currencyBalance.TokenId == int32(tk.Id) {
				fmt.Println("fabi : ", tk.Mark, " you :", currencyBalance.TotalBalance, currencyBalance.TotalFreeze)
				totalnum, _ = convert.StringAddString(totalnum, currencyBalance.TotalBalance)
				totalnum, _ = convert.StringAddString(totalnum, currencyBalance.TotalFreeze)
			}
		}

		if totalnum == `` {
			tmp.TotalNum = "0"
		} else {
			tmp.TotalNum, _ = convert.StringTo8Bit(totalnum)
		}

		if totalFree == `` {
			tmp.TotalFree = "0"
		} else {
			tmp.TotalFree, _ = convert.StringTo8Bit(totalFree)
		}

		if err != nil {
			fmt.Println(err)
		}
		for _, tkUser := range tokenUserCoin {
			if tkUser.TokenId == int32(tk.Id) {
				totaluser = totaluser + tkUser.TotalUser
			}
		}
		for _, cuUser := range currencyUserCoin {
			if cuUser.TokenId == int32(tk.Id) {
				totaluser = totaluser + cuUser.TotalUser
			}
		}
		fmt.Println("totalnum:", totalnum, totaluser)
		tmp.TotalUser = totaluser
		if totaluser <= 0 {
			tmp.AverageNum = "0"
		} else {
			fmt.Println(totalnum, totaluser, convert.Int64ToStringAdd8Bit(totaluser))

			tempStr, err := convert.StringDivString(totalnum, convert.Int64ToStringAdd8Bit(totaluser))
			if err != nil {
				tmp.AverageNum = "0"
			} else {
				fmt.Println(tk.Mark, " 人均持有=", tempStr)
				if tempStr == `` {
					tempStr = "0"
				}
				tmp.AverageNum = tempStr
			}

		}

		totalcoinList = append(totalcoinList, tmp)
	}
	respList := new(models.ModelList)
	respList.IsPage = true
	respList.Items = totalcoinList
	respList.Total = int(total)
	respList.PageIndex = req.Page
	respList.PageSize = req.Rows
	var pagecount int
	if int(total)%req.Rows == 0 {
		pagecount = int(total) / req.Rows
	} else {
		pagecount = (int(total) / req.Rows) + 1
	}
	respList.PageCount = pagecount

	cu.Put(c, "list", respList)
	cu.RespOK(c)
	return
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
	cu.buySellList(c)
	return
}
func (cu *CurrencyController) ExportBuySellList(c *gin.Context) {
	cu.buySellList(c)
	return
}

func (cu *CurrencyController) buySellList(c *gin.Context) {
	req := struct {
		Uid     int `form:"uid" json:"uid" binding:"required"`
		Page    int `form:"page" json:"page" binding:"required"`
		Rows    int `form:"rows" json:"rows" `
		TokenId int `form:"token_id" json:"token_id"`
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
	list, err := new(models.Order).GetOrderListOfUid(req.Page, req.Rows, req.Uid, req.TokenId)
	if err != nil {
		cu.RespErr(c, err)
		return
	}

	cu.Put(c, "list", list)
	cu.RespOK(c)
	return
}

func (cu *CurrencyController) GetUserDetailList(c *gin.Context) {
	cu.userDetailList(c)
	return
}

func (cu *CurrencyController) ExportUserDetailList(c *gin.Context) {
	cu.userDetailList(c)
	return
}
func (cu *CurrencyController) userDetailList(c *gin.Context) {
	req := struct {
		Uid     int `form:"uid" json:"uid" binding:"required"`
		Page    int `form:"page" json:"page" binding:"required"`
		Rows    int `form:"rows" json:"rows" `
		TokenId int `form:"token_id" json:"token_id"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}

	list, err := new(models.UserCurrency).GetCurrencyList(req.Page, req.Rows, req.Uid, req.TokenId)
	if err != nil {
		cu.RespErr(c, err)
		return
	}
	cu.Put(c, "list", list)
	cu.RespOK(c)
	return
}

//p2-3-1法币账户统计列表
func (cu *CurrencyController) GetTotalCurrencyBalance(c *gin.Context) {
	cu.totalCurrencyBalance(c)
	return
}

func (cu *CurrencyController) ExportTotalCurrencyBalance(c *gin.Context) {
	cu.totalCurrencyBalance(c)
	return
}

func (cu *CurrencyController) totalCurrencyBalance(c *gin.Context) {
	req := struct {
		Page    int    `form:"page" json:"page" binding:"required"`
		Rows    int    `form:"rows" json:"rows" `
		TokenId int    `form:"tid" json:"tid"`
		Search  string `form:"search" json:"search" ` //搜索的内容
		Status  int    `form:"status" json:"status" ` //用户账号状态
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}
	////result, err := new(models.UserGroup).GetAllUser(req.Page, req.Page_num, req.Status, req.Search)
	result, err := new(models.UserCurrency).CurrencyBalanceNew(req.Page, req.Rows, req.Status, req.TokenId, req.Search)
	if err != nil {
		cu.RespErr(c, err)
	}

	uidList := make([]uint64, 0)
	value, OK := result.Items.([]models.AmountToCny)
	if !OK {
		cu.RespErr(c, errors.New("assert failed"))
		return
	}
	for _, value := range value {
		uidList = append(uidList, uint64(value.Uid))
	}
	//fmt.Println("uid", uidList)

	tokenList, err := new(apis.VendorApi).GetCny(uidList, 2)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		cu.RespErr(c, err.Error())
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

	result.Items = value
	//法币账户折和没有计算
	cu.Put(c, "list", result)

	// 返回
	cu.RespOK(c)
	return
}

//法币挂单管理
func (cu *CurrencyController) GetTradeList(c *gin.Context) {
	cu.tradeList(c)
	return

}
func (cu *CurrencyController) ExportTradeList(c *gin.Context) {
	cu.tradeList(c)
	return

}

func (cu *CurrencyController) tradeList(c *gin.Context) {
	req := struct {
		Page    int    `form:"page" json:"page" binding:"required"`
		PageNum int    `form:"rows" json:"rows" `
		Ustatus int    `form:"status" json:"status" ` //用户登录状态
		Search  string `form:"search" json:"search" `
		Verify  int    `form:"verify" json:"verify" ` //实名认证 二级认证 google 验证  交易权限
		//Date    string `form:"date" json:"date" `         //挂单日期
		TokenId int `form:"token_id" json:"token_id" ` //货币名称
		TradeId int `form:"tid" json:"tid" `           //交易方向
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}
	list, err := new(models.Ads).GetAdsList(req.Page, req.PageNum, req.Ustatus, req.TokenId, req.TradeId, req.Verify, req.Search)
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
	cu.oderList(c)
	return
}
func (cu *CurrencyController) ExportOderList(c *gin.Context) {
	cu.oderList(c)
	return
}

func (cu *CurrencyController) oderList(c *gin.Context) {
	//参数一大堆
	req := struct {
		Page int `form:"page" json:"page" binding:"required"`
		Rows int `form:"rows" json:"rows" `
		//Start_t  string `form:"start_t" json:"start_t" `
		Search  string `form:"search" json:"search" `     //筛选
		Status  int    `form:"status" json:"status" `     //订单状态
		TokenId int    `form:"token_id" json:"token_id" ` //货币名称
		AdType  int    `form:"adtype" json:"adtype" `     //买卖方向
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		cu.RespErr(c, err)
		return
	}
	list, err := new(models.Order).GetOrderList(req.Page, req.Rows, req.AdType, req.Status, req.TokenId, req.Search)
	if err != nil {
		cu.RespErr(c, err)
		return
	}
	cu.Put(c, "list", list)
	cu.RespOK(c)
	return
}

// 法币交易走势
func (t *CurrencyController) TradeTrend(ctx *gin.Context) {
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
	list, err := new(models.CurrencyDailySheet).TradeTrendList(filter)
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
	loc, _ := time.LoadLocation("Local")
	for k, v := range list {
		datetime, _ := time.ParseInLocation(utils.LAYOUT_DATE_TIME, v.Date, loc)
		x[k] = datetime.Format("0102")
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
