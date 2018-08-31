package models

import (
	"admin/errors"
	"admin/utils"
	"admin/utils/convert"
	"fmt"
	"time"

	"log"

	"github.com/robfig/cron"
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
	TokenId      int64 `xorm:"not null comment('货币id') INT(11)" json:"token_id"`
	FeeBuyCny    int64 `xorm:"not null comment('买手续费折合cny') BIGINT(20)" json:"fee_buy_cny"`
	FeeBuyTotal  int64 `xorm:"not null comment('买手续费总额') BIGINT(20)" json:"fee_buy_total"`
	FeeSellCny   int64 `xorm:"not null comment('卖手续费折合cny') BIGINT(20)" json:"fee_sell_cny"`
	FeeSellTotal int64 `xorm:"not null comment('卖手续费总额') BIGINT(20)" json:"fee_sell_total"`
	BuyTotal     int64 `xorm:"not null comment('买总额') BIGINT(20)" json:"buy_total"`
	BuyTotalCny  int64 `xorm:"not null comment('买总额折合') BIGINT(20)" json:"but_total_cny"`
	SellTotalCny int64 `xorm:"not null comment('卖总额折合') BIGINT(20)" json:"sell_total_cny"`
	SellTotal    int64 `xorm:"not null comment('卖总额') BIGINT(20)" json:"sell_total"`
	BalanceAll   int64 `xorm:"comment('币总额') BIGINT(20)"`
	FrozenAll    int64 `xorm:"comment('冻结总额') BIGINT(20)"`
	Date         int64 `xorm:"not null comment('时间戳，精确到天') BIGINT(10)" json:"date"`
}

type TokenFeeDailySheetGroup struct {
	TotalBuy  string `json:"total_buy"`
	TotalSell string `json:"total_sell"`
	Total     string `json:"total"`
}

type total struct {
	Id    int
	Total string `xorm:"-" json:"total" `
	Buy   string `json:"buy"`
	Sell  string `json:"sell"`
	Date  int64  `json:"date"`
	DateStr string `xorm:"-" json:"date_str"`
 }

// 走势返回string，内容是int
// 如果用int64，数据太大时xorm sum会溢出报错
type TokenTradeTrend struct {
	BuyTotal     string `xorm:"buy_total"`
	BuyTotalCny  string `xorm:"but_total_cny"`
	SellTotal    string `xorm:"sell_total"`
	SellTotalCny string `xorm:"sell_total_cny"`
	FeeBuyTotal  string `xorm:"fee_buy_total"`
	FeeBuyCny    string `xorm:"fee_buy_cny"`
	FeeSellTotal string `xorm:"fee_sell_total"`
	FeeSellCny   string `xorm:"fee_sell_cny"`
	Date         int64  `xorm:"date"`
}

// 交易走势
func (this *TokenDailySheet) TradeTrendList(filter map[string]interface{}) ([]*TokenTradeTrend, error) {
	// 时间区间，默认最近一周
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return nil, errors.NewSys(err)
	}
	today, err := time.ParseInLocation(utils.LAYOUT_DATE, time.Now().Format(utils.LAYOUT_DATE), loc)
	if err != nil {
		return nil, errors.NewSys(err)
	}
	todayZeroUnix := today.Unix()

	dateBegin := todayZeroUnix - 7*24*60*60
	dateEnd := todayZeroUnix

	// 开始查询
	session := utils.Engine_token.Where("1=1")

	// 筛选
	if v, ok := filter["date_begin"]; ok {
		dateBegin, _ = v.(int64)
	}
	if v, ok := filter["date_end"]; ok {
		dateEnd, _ = v.(int64)
	}
	if v, ok := filter["token_id"]; ok {
		session.And("token_id=?", v)
	}

	var list []*TokenTradeTrend
	err = session.
		Table(this).
		Select("date, sum(buy_total) as buy_total, sum(buy_total_cny) as buy_total_cny"+
			",sum(sell_total) as sell_total, sum(sell_total_cny) as sell_total_cny"+
			",sum(fee_buy_total) as fee_buy_total, sum(fee_buy_cny) as fee_buy_cny"+
			",sum(fee_sell_total) as fee_sell_total, sum(fee_sell_cny) as fee_sell_cny").
		And("date>=?", dateBegin).
		And("date<=?", dateEnd).
		GroupBy("date").
		OrderBy("date ASC").
		Find(&list)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	return list, nil
}

//手续费报表 一天显示一条记录
func (this *TokenDailySheet) GetDailySheetList(page, rows int, date uint64) (*ModelList, *TokenFeeDailySheetGroup, error) {
	engine := utils.Engine_token
	fmt.Println("bibi 交易手续费汇总")
	//query := engine.Desc("date")
	//sql := " SELECT id,date,SUM(fee_buy_cny) buy, SUM(fee_sell_cny) sell  FROM `token_daily_sheet` GROUP BY date ORDER BY `date` DESC "
	//sql := " SELECT * FROM (SELECT id,DATE,SUM(fee_buy_cny) AS buy, SUM(fee_sell_cny) AS sell FROM `token_daily_sheet`  GROUP BY DATE ORDER BY `date`  DESC  )t WHERE t.buy !=0 OR t.sell !=0    "

	sql := "SELECT * FROM (SELECT id,DATE,SUM(fee_buy_total) AS buy, SUM(fee_sell_total) AS sell FROM `token_daily_sheet`  GROUP BY DATE ORDER BY `date`  DESC  )t WHERE t.buy !=0 OR t.sell !=0   "

	if date != 0 {
		//sql = fmt.Sprintf(" SELECT id,date,SUM(fee_buy_cny) as buy, SUM(fee_sell_cny) as sell  FROM `token_daily_sheet` where date between %d and %d GROUP BY date ORDER BY `date` DESC", date, date+86400)
		//sql = fmt.Sprintf("SELECT * FROM (SELECT id,DATE,SUM(fee_buy_cny) AS buy, SUM(fee_sell_cny) AS sell FROM `token_daily_sheet`  where date between %d and %d  GROUP BY DATE ORDER BY `date`  DESC  )t WHERE t.buy !=0 OR t.sell !=0   ", date, date+86400)
		sql = fmt.Sprintf("SELECT * FROM (SELECT id,DATE,SUM(fee_buy_total) AS buy, SUM(fee_sell_total) AS sell FROM `token_daily_sheet`  where date between %d and %d  GROUP BY DATE ORDER BY `date`  DESC  )t WHERE t.buy !=0 OR t.sell !=0   ", date, date+86400)
	}
	Count := &struct {
		Num int64
	}{}
	countSql := fmt.Sprintf("select  count(id) num from (%s) t", sql)
	_, err := engine.SQL(countSql).Get(Count)
	if err != nil {
		return nil, nil, err
	}
	offset, mList := this.Paging(page, rows, int(Count.Num))
	list := make([]total, 0)
	limitSql := fmt.Sprintf(" limit %d offset %d", mList.PageSize, offset)

	err = engine.Table("token_daily_sheet").SQL(sql + limitSql).Find(&list)
	if err != nil {
		return nil, nil, err
	}
	for i, v := range list {

		temp, _ := convert.StringAddString(v.Buy, v.Sell)
		list[i].Total, _ = convert.StringTo8Bit(temp)
		list[i].Buy, _ = convert.StringTo8Bit(v.Buy)
		list[i].Sell, _ = convert.StringTo8Bit(v.Sell)
		list[i].DateStr = time.Unix(v.Date,0).Format("2006-01-02 15:04:05")
	}
	mList.Items = list
	tfd := new(TokenFeeDailySheetGroup)
	//_, err = engine.SQL("SELECT COALESCE(sum(`fee_buy_cny`),0) AS total_buy, COALESCE(sum(`fee_sell_cny`),0) AS total_sell FROM `token_daily_sheet`").Get(tfd)
	_, err = engine.SQL("SELECT COALESCE(sum(`fee_buy_total`),0) AS total_buy, COALESCE(sum(`fee_sell_total`),0) AS total_sell FROM `token_daily_sheet`").Get(tfd)

	if err != nil {
		return nil, nil, err
	}
	tfd.Total, _ = convert.StringAddString(tfd.TotalBuy, tfd.TotalSell)
	tfd.Total, _ = convert.StringTo8Bit(tfd.Total)
	tfd.TotalSell, _ = convert.StringTo8Bit(tfd.TotalSell)
	tfd.TotalBuy, _ = convert.StringTo8Bit(tfd.TotalBuy)
	return mList, tfd, nil
}


//wyw 币币交易手续费定时任务呢
func (tk *TokenDailySheet) TimingFuncNew(begin, end int64) {
	//币币交易
	//第一步获取货币id
	engine := utils.Engine_token
	tokenIdList, err := new(Tokens).GetTokensList()
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return
	}
	buyList := make([]DayCount, 0)
	sellList := make([]DayCount, 0)
	//买入
	for _, v := range tokenIdList {
		if v.Id == 0 {
			continue
		}
		buy, err := new(Trade).Get(v.Id, begin, end, 1)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if buy.TokenId != 0 {
			buyList = append(buyList, *buy)
		}

		sell, err := new(Trade).Get(v.Id, begin, end, 2)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if sell.TokenId != 0 {
			sellList = append(sellList, *sell)
		}

	}

	//插入数据
	for _, v := range buyList {
		fmt.Println("buy")
		tdsheet := TokenDailySheet{TokenId: v.TokenId, Date: v.Date}
		isExistSql := " SELECT token_id, `date` FROM  g_token.`token_daily_sheet`   WHERE token_id=? AND `date`=?"
		has, err := engine.Table("token_daily_sheet").SQL(isExistSql, v.TokenId, v.Date).Get(&tdsheet)
		if err != nil {
			fmt.Println(err)
			utils.AdminLog.Errorln(err)
		}

		 feeTotal, err := convert.StringToInt64(v.FeeTotal)
		 if err!=nil{
			 utils.AdminLog.Info("buy", err.Error())
			 fmt.Println(err.Error())
		 }
		 feeTotalCny, err := convert.StringToInt64(v.FeeTotalCny)
		if err!=nil{
			utils.AdminLog.Info("buy", err.Error())
			fmt.Println(err.Error())
		}
		 buyTotal, err := convert.StringToInt64(v.Total)
		if err!=nil{
			utils.AdminLog.Info("buy", err.Error())
			fmt.Println(err.Error())
		}
		 buyTotalCny, err := convert.StringToInt64(v.TotalCny)
		if err!=nil{
			utils.AdminLog.Info("buy", err.Error())
			fmt.Println(err.Error())
		}
		fmt.Println("but", tdsheet)
		if has {
			tdsheet.BuyTotalCny,err =convert.Int64AddInt64(tdsheet.BuyTotalCny,buyTotalCny)
			if err!=nil{
				utils.AdminLog.Info("buy", err.Error())
				fmt.Println(err.Error())
			}
			tdsheet.BuyTotal,err=convert.Int64AddInt64(tdsheet.BuyTotal,buyTotal)
			if err!=nil{
				utils.AdminLog.Info("buy", err.Error())
				fmt.Println(err.Error())
			}
			tdsheet.FeeBuyCny,err=convert.Int64AddInt64(tdsheet.FeeBuyCny,feeTotalCny)
			if err!=nil{
				utils.AdminLog.Info("buy", err.Error())
				fmt.Println(err.Error())
			}
			tdsheet.FeeBuyTotal,err=convert.Int64AddInt64(tdsheet.FeeBuyTotal,feeTotal)
			if err!=nil{
				utils.AdminLog.Info("buy", err.Error())
				fmt.Println(err.Error())
			}
			_, err := engine.Cols("fee_buy_total", "fee_buy_cny", "buy_total", "buy_total_cny").Where("id=?", tdsheet.Id).Update(&tdsheet)
			if err != nil {
				utils.AdminLog.Info("buyList定时任务执行更新数据失败", err.Error())
				fmt.Println(err.Error())

			}
		} else {
			tdsheet.BuyTotalCny=buyTotalCny
			tdsheet.BuyTotal=buyTotal
			tdsheet.FeeBuyCny=feeTotalCny
			tdsheet.FeeBuyTotal=feeTotal
			_, err := engine.Cols("token_id", "fee_buy_total", "fee_buy_cny", "buy_total", "buy_total_cny", "date").InsertOne(&tdsheet)
			if err != nil {
				utils.AdminLog.Info("buyList定时任务执行更新数据失败", err.Error())
				fmt.Println(err.Error())

			}
		}

	}

	for _, v := range sellList {
		fmt.Println("sell")
		tdsheet := TokenDailySheet{TokenId: v.TokenId, Date: v.Date}
		isExistSql := " SELECT token_id, `date` FROM  g_token.`token_daily_sheet`   WHERE token_id=? AND `date`=?"
		has, err := engine.Table("token_daily_sheet").SQL(isExistSql, v.TokenId, v.Date).Get(&tdsheet)
		if err != nil {
			fmt.Println(err)
			utils.AdminLog.Errorln(err)
		}

		feeTotal, err := convert.StringToInt64(v.FeeTotal)
		if err != nil {
			utils.AdminLog.Info("sell", err.Error())
			fmt.Println(err.Error())
		}
		feeTotalCny, err := convert.StringToInt64(v.FeeTotalCny)
		if err != nil {
			utils.AdminLog.Info("sell", err.Error())
			fmt.Println(err.Error())
		}
		buyTotal, err := convert.StringToInt64(v.Total)
		if err != nil {
			utils.AdminLog.Info("sell", err.Error())
			fmt.Println(err.Error())
		}
		buyTotalCny, err := convert.StringToInt64(v.TotalCny)
		if err != nil {
			utils.AdminLog.Info("sell", err.Error())
			fmt.Println(err.Error())
		}

		fmt.Println("sell", tdsheet)
		if has {
			tdsheet.SellTotalCny,err =convert.Int64AddInt64(tdsheet.SellTotalCny,buyTotalCny)
			if err != nil {
				utils.AdminLog.Info("sell", err.Error())
				fmt.Println(err.Error())
			}
			tdsheet.SellTotal,err=convert.Int64AddInt64(tdsheet.SellTotal,buyTotal)
			if err != nil {
				utils.AdminLog.Info("sell", err.Error())
				fmt.Println(err.Error())
			}
			tdsheet.FeeSellCny,err=convert.Int64AddInt64(tdsheet.FeeSellCny,feeTotalCny)
			if err != nil {
				utils.AdminLog.Info("sell", err.Error())
				fmt.Println(err.Error())
			}
			tdsheet.FeeSellTotal,err=convert.Int64AddInt64(tdsheet.FeeSellTotal,feeTotal)
			if err != nil {
				utils.AdminLog.Info("sell", err.Error())
				fmt.Println(err.Error())
			}
			_, err := engine.Cols("fee_sell_total", "fee_sell_cny", "sell_total", "sell_total_cny").Where("id=?", tdsheet.Id).Update(&tdsheet)
			if err != nil {
				utils.AdminLog.Info("sellList定时任务执行更新数据失败", err.Error())
				fmt.Println(err.Error())
			}
		} else {
			_, err := engine.Cols("token_id", "fee_sell_total", "fee_sell_cny", "sell_total", "sell_total_cny", "date").InsertOne(&tdsheet)
			if err != nil {
				utils.AdminLog.Info("sellList定时任务执行更新数据失败", err.Error())
				fmt.Println(err.Error())
			}
		}

	}
	//法币 冻结
	//更新 字段 BalanceAll FrozenAll
	for _, vu := range tokenIdList {
		fmt.Println("balance_all")
		balance, free, err := new(UserCurrency).GetUserCurrencyBalanceAndFree(vu.Id)
		if err != nil {
			utils.AdminLog.Info("sellList定时任务执行更新数据失败", err.Error())
			fmt.Println(err.Error())
			continue
		}
		tdsheet := TokenDailySheet{TokenId: int64(vu.Id), Date: begin}
		isExistSql := " SELECT token_id, `date` FROM  g_token.`token_daily_sheet`   WHERE token_id=? AND `date`=?"
		has, err := engine.Table("token_daily_sheet").SQL(isExistSql, vu.Id, begin).Get(&tdsheet)
		if err != nil {
			fmt.Println(err)
			utils.AdminLog.Errorln(err)
		}

		fmt.Println("balance_all", tdsheet)
		if has {
			tdsheet.FrozenAll,err =convert.Int64AddInt64(tdsheet.FrozenAll,free)
			if err != nil {
				utils.AdminLog.Info("frozenAll", err.Error())
				fmt.Println(err.Error())
			}
			tdsheet.BalanceAll ,err =convert.Int64AddInt64(tdsheet.FrozenAll,balance)
			if err != nil {
				utils.AdminLog.Info("balanceAll", err.Error())
				fmt.Println(err.Error())
			}
			_, err := engine.Cols("balance_all", "frozen_all").Where("id=?", tdsheet.Id).Update(&tdsheet)
			if err != nil {
				fmt.Println(err)
				utils.AdminLog.Info("insert", err.Error())
			}
		} else {
			tdsheet.FrozenAll = free
			tdsheet.BalanceAll = balance
			_, err := engine.Cols("token_id", "balance_all", "frozen_all", "date").InsertOne(&tdsheet)
			if err != nil {
				fmt.Println(err)
				utils.AdminLog.Info("update", err.Error())
			}
		}

	}
	//法币 冻结
	//更新 字段 BalanceAll FrozenAll
	for _, vu := range tokenIdList {
		fmt.Println("balance_all")
		balance, free, err := new(UserToken).GetUserTokenBalance(vu.Id)
		if err != nil {
			utils.AdminLog.Info("sellList定时任务执行更新数据失败", err.Error())
			fmt.Println(err.Error())
			continue
		}
		tdsheet := TokenDailySheet{TokenId: int64(vu.Id), Date: begin}
		isExistSql := " SELECT token_id, `date` FROM  g_token.`token_daily_sheet`   WHERE token_id=? AND `date`=?"
		has, err := engine.Table("token_daily_sheet").SQL(isExistSql, vu.Id, begin).Get(&tdsheet)
		if err != nil {
			fmt.Println(err)
			utils.AdminLog.Errorln(err)
		}
		tdsheet.FrozenAll, _ = convert.Int64AddInt64(tdsheet.FrozenAll, free)
		tdsheet.BalanceAll, _ = convert.Int64AddInt64(tdsheet.BalanceAll, balance)
		if has {
			if tdsheet.FrozenAll == 0 && tdsheet.BalanceAll == 0 {
				continue
			}

			_, err := engine.Cols("balance_all", "frozen_all").Where("token_id=? and date=?", tdsheet.TokenId, tdsheet.Date).Update(&tdsheet)
			if err != nil {
				fmt.Println(err)
				utils.AdminLog.Errorln(err)
			}
		} else {

			_, err := engine.Cols("token_id", "balance_all", "frozen_all", "date").InsertOne(&tdsheet)
			if err != nil {
				fmt.Println(err)
				utils.AdminLog.Errorln(err)
			}
		}

	}

}

func (t *TokenDailySheet) Run() {
	toBeCharge := time.Now().Format("2006-01-02 ") + "00:00:00"
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
	unix := theTime.Unix()
	fmt.Println("当前时间戳", unix)
	t.TimingFuncNew(unix-86400, unix)
	//t.TimingFunc(1532448000, unix)
}

//启动
func DailyStart() {

	fmt.Println("daily count start ....")
	log.Println("daily count start ....")

	i := 0
	c := cron.New()

	//AddFunc
	spec := "0 0 1 * *" // every day ...
	c.AddFunc(spec, func() {
		i++
		log.Println("cron running:", i)
	})
	//AddJob方法
	c.AddJob(spec, &TokenDailySheet{})
	//启动计划任务
	c.Start()
	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer c.Stop()

	select {}
}

func DailyStart1() {
	new(TokenDailySheet).Run_tool()

}

func (t *TokenDailySheet) Run_tool() {
	toBeCharge := time.Now().Format("2006-01-02 ") + "00:00:00"
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
	unix := theTime.Unix()
	fmt.Println("当前时间戳", unix)
	begin := 1533657600
	BeginTime := int64(begin)
	i := 0
	for BeginTime < unix {
		i++
		fmt.Println("第", i, "次循环")
		t.TimingFuncNew(BeginTime, BeginTime+86400) //test
		BeginTime += 86400
	}

}
