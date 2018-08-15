package models

import (
	"admin/utils"
)

//type TokenFeeDailySheet struct {
//	BaseModel    `xorm:"-"`
//	Id           int   `xorm:"not null pk autoincr comment('自增id') TINYINT(4)"`
//	Date         int64 `xorm:"not null comment('日期 精确到天 20180727 年保持4位月两位日两位') BIGINT(20)" json:"date"`
//	BuyBalance   int64 `xorm:"not null comment('买手续费日总额单位为人民币') BIGINT(20)" json:"buy_balance"`
//	SellBalance  int64 `xorm:"not null comment('卖手续费日总额单位为人民币') BIGINT(20)" json:"sell_balance"`
//	TotalBalance int64 `xorm:"not null comment('买卖总额单位为人民币') BIGINT(20)" json:"total_balance"`
//}

type TokenDailySheet struct {
	BaseModel    `xorm:"-"`
	Id           int   `xorm:"not null pk autoincr comment('自增id') TINYINT(4)"`
	TokenId      int   `xorm:"not null comment('货币id') INT(11)" json:"token_id"`
	FeeBuyCny    int64 `xorm:"not null comment('买手续费折合cny') BIGINT(20)" json:"fee_buy_cny"`
	FeeBuyTotal  int64 `xorm:"not null comment('买手续费总额') BIGINT(20)" json:"fee_buy_total"`
	FeeSellCny   int64 `xorm:"not null comment('卖手续费折合cny') BIGINT(20)" json:"fee_sell_cny"`
	FeeSellTotal int64 `xorm:"not null comment('卖手续费总额') BIGINT(20)" json:"fee_sell_total"`
	BuyTotal     int64 `xorm:"not null comment('买总额') BIGINT(20)" json:"buy_total"`
	ButTotalCny  int64 `xorm:"not null comment('买总额折合') BIGINT(20)" json:"but_total_cny"`
	SellTotalCny int64 `xorm:"not null comment('卖总额折合') BIGINT(20)" json:"sell_total_cny"`
	SellTotal    int64 `xorm:"not null comment('卖总额') BIGINT(20)" json:"sell_total"`
	Date         int64 `xorm:"not null comment('时间戳，精确到天') BIGINT(10)" json:"date"`
}

type TokenFeeDailySheetGroup struct {
	TotalBuy  float64 `json:"total_buy"`
	TotalSell float64 `json:"total_sell"`
	Total     float64 `json:"total"`
}

//func (this *TokenDailySheet)Test1(){
//	now := time.Now()
//	current := now.Format("2006-01-02 15:04:05")
//	cunrrentUnixtime := now.Unix()
//	//bibi 日报表统计
//	type statisticsDayTrade struct {
//		TotalBuyFeeCny  int64 //买
//		TotalSellFeeCny int64 //卖
//		TotalBalanceCny int64 //总
//		Total           int64
//		Date            int64
//	}
//
//	engine := utils.Engine_token
//	sql := "FROM (SELECT  FROM_UNIXTIME(deal_time,'%Y-%m-%d') days, states,opt, fee_cny FROM g_token.trade) t "
//	hearSql := fmt.Sprintf("SELECT SUM(fee_cny)  %s ", "total_buy_fee_cny ")
//	temp := fmt.Sprintf(" WHERE t.states=2 AND t.opt=%d AND t.days='%s'", 1, current)
//	sql1 := hearSql + sql + temp
//	statistics := new(statisticsDayTrade)
//	_, err := engine.SQL(sql1).Get(statistics)
//	if err != nil {
//		fmt.Println(err)
//		utils.AdminLog.Println("定时任务执行失败")
//		return
//	}
//	fmt.Println(statistics)
//	hearSql = fmt.Sprintf("SELECT SUM(fee_cny)  %s ", "total_sell_fee_cny ")
//	sqlSell := fmt.Sprintf(" WHERE t.states=2 AND t.opt=%d AND t.days='%s'", 2, current)
//	_, err = engine.SQL(hearSql + sql + sqlSell).Get(statistics)
//	if err != nil {
//		utils.AdminLog.Println("定时任务执行失败")
//		return
//	}
//	statistics.TotalBalanceCny = statistics.TotalBuyFeeCny + statistics.TotalSellFeeCny
//	statistics.Date = cunrrentUnixtime
//	_, err = engine.InsertOne(&TokenDailySheet{
//		TotalBalance: statistics.TotalBalanceCny,
//		BuyBalance:   statistics.TotalBuyFeeCny,
//		SellBalance:  statistics.TotalSellFeeCny,
//		Date:         statistics.Date,
//	})
//	fmt.Println(statistics)
//	if err != nil {
//		utils.AdminLog.Println("定时任务执行失败")
//		return
//	}
//	fmt.Println("successful")
//}
//
//定时结算bibi 日交易报表表数据
//func (this *TokenDailySheet) BoottimeTimingSettlement() {
//
//	for {
//		now := time.Now()
//		// 计算下一个零点
//		next := now.Add(time.Hour * 24)
//		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
//		t := time.NewTimer(next.Sub(now))
//		<-t.C
//		//Printf("定时结算Boottime表数据，结算完成: %v\n",time.Now())
//		//以下为定时执行的操作
//		current := now.Format("2006-01-02 15:04:05")
//		cunrrentUnixtime := now.Unix()
//		//bibi 日报表统计
//		type statisticsDayTrade struct {
//			TotalBuyFeeCny  int64 //买
//			TotalSellFeeCny int64 //卖
//			TotalBalanceCny int64 //总
//			Total           int64
//			Date            int64
//		}
//
//		engine := utils.Engine_token
//		sql := "FROM (SELECT  FROM_UNIXTIME(deal_time,'%Y-%m-%d') days, states,opt, fee_cny FROM g_token.trade) t "
//		hearSql := fmt.Sprintf("SELECT SUM(fee_cny)  %s ", "total_buy_fee_cny ")
//		temp := fmt.Sprintf(" WHERE t.states=2 AND t.opt=%d AND t.days='%s'", 1, current)
//		sql1 := hearSql + sql + temp
//		statistics := new(statisticsDayTrade)
//		_, err := engine.SQL(sql1).Get(statistics)
//		if err != nil {
//			fmt.Println(err)
//			utils.AdminLog.Println("定时任务执行失败")
//			continue
//		}
//		fmt.Println(statistics)
//		hearSql = fmt.Sprintf("SELECT SUM(fee_cny)  %s ", "total_sell_fee_cny ")
//		sqlSell := fmt.Sprintf(" WHERE t.states=2 AND t.opt=%d AND t.days='%s'", 2, current)
//		_, err = engine.SQL(hearSql + sql + sqlSell).Get(statistics)
//		if err != nil {
//			utils.AdminLog.Println("定时任务执行失败")
//			continue
//		}
//		statistics.TotalBalanceCny = statistics.TotalBuyFeeCny + statistics.TotalSellFeeCny
//		statistics.Date = cunrrentUnixtime
//		_, err = engine.AllCols().InsertOne(&TokenDailySheet{
//			TotalBalance: statistics.TotalBalanceCny,
//			BuyBalance:   statistics.TotalBuyFeeCny,
//			SellBalance:  statistics.TotalSellFeeCny,
//			Date:         statistics.Date,
//		})
//		fmt.Println(statistics)
//		if err != nil {
//			utils.AdminLog.Println("定时任务执行失败")
//			continue
//		}
//		fmt.Println("successful")
//	}
//}

//获取历史交易记录
func (this *TokenDailySheet) GetDailySheetList(page, rows int, date uint64) (*ModelList, *TokenFeeDailySheetGroup, error) {
	engine := utils.Engine_token
	query := engine.Desc("id")
	if date != 0 {
		query = query.Where("date between ? and ?", date, date+86400)
	}
	countQuery := *query
	count, err := countQuery.Count(&TokenDailySheet{})
	if err != nil {
		return nil, nil, err
	}
	offset, mList := this.Paging(page, rows, int(count))
	list := make([]TokenDailySheet, 0)
	err = query.Limit(mList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, nil, err
	}
	mList.Items = list
	result, err := engine.Sums(this, "buy_balance", "sell_balance", "total_balance")
	if err != nil {
		return nil, nil, err
	}
	total := result[2]
	totalBuy := result[1]
	totalSell := result[0]
	return mList, &TokenFeeDailySheetGroup{
		Total:     total,
		TotalBuy:  totalBuy,
		TotalSell: totalSell,
	}, nil
}
