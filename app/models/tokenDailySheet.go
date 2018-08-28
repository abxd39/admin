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
}

// 手续费合计
type TokenFeeTotal struct {
	TodayTotal       string `xorm:"today_total"`         // 今日合计
	YesterdayTotal   string `xorm:"yesterday_total"`     // 上日合计
	LastWeekDayTotal string `xorm:"last_week_day_total"` // 上周同日合计
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

// 手续费合计
// 今日、上日、上周同日
func (this *TokenDailySheet) FeeTotal() (*TokenFeeTotal, error) {
	// 计算日期
	todayDate := time.Now().Format(utils.LAYOUT_DATE)
	datetime, _ := time.Parse(utils.LAYOUT_DATE_TIME, fmt.Sprintf("%s 00:00:00", todayDate))
	todayTime := datetime.Unix() // 零点
	yesterdayTime := todayTime - 24*60*60
	lastWeekDayTime := todayTime - 7*24*60*60

	// 开始合计
	//1. 今日
	feeTotal := &TokenFeeTotal{}
	session := utils.Engine_token.Where("1=1")
	_, err := session.
		Table(this).
		Select("IFNULL(sum(fee_buy_total+fee_sell_total), 0) today_total").
		And("date=?", todayTime).
		Get(feeTotal)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	//2. 上日
	yesFeeTotal := &TokenFeeTotal{}
	yesSession := utils.Engine_token.Where("1=1")
	_, err = yesSession.
		Table(this).
		Select("IFNULL(sum(fee_buy_total+fee_sell_total), 0) yesterday_total").
		And("date=?", yesterdayTime).
		Get(yesFeeTotal)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	//3. 上周同日
	lastWeekFeeTotal := &TokenFeeTotal{}
	lastWeekSession := utils.Engine_token.Where("1=1")
	_, err = lastWeekSession.
		Table(this).
		Select("IFNULL(sum(fee_buy_total+fee_sell_total), 0) last_week_day_total").
		And("date=?", lastWeekDayTime).
		Get(lastWeekFeeTotal)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	// 合并
	feeTotal.YesterdayTotal = yesFeeTotal.YesterdayTotal
	feeTotal.LastWeekDayTotal = lastWeekFeeTotal.LastWeekDayTotal

	return feeTotal, nil
}

// 交易走势
func (this *TokenDailySheet) TradeTrendList(filter map[string]interface{}) ([]*TokenTradeTrend, error) {
	// 时间区间，默认最近一周
	today, err := time.Parse(utils.LAYOUT_DATE_TIME, fmt.Sprintf("%s 00:00:00", time.Now().Format(utils.LAYOUT_DATE)))
	if err != nil {
		return nil, errors.NewSys(err)
	}
	todayTime := today.Unix()

	dateBegin := todayTime - 6*24*60*60
	dateEnd := todayTime

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
	sql := " SELECT id,date,SUM(fee_buy_cny) buy, SUM(fee_sell_cny) sell  FROM `token_daily_sheet` GROUP BY date ORDER BY `date` DESC "
	if date != 0 {
		sql = fmt.Sprintf(" SELECT id,date,SUM(fee_buy_cny) as buy, SUM(fee_sell_cny) as sell  FROM `token_daily_sheet` where date between %d and %d GROUP BY date ORDER BY `date` DESC", date, date+86400)
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
	}
	mList.Items = list
	tfd := new(TokenFeeDailySheetGroup)
	_, err = engine.SQL("SELECT COALESCE(sum(`fee_buy_cny`),0) AS total_buy, COALESCE(sum(`fee_sell_cny`),0) AS total_sell FROM `token_daily_sheet`").Get(tfd)
	if err != nil {
		return nil, nil, err
	}
	tfd.Total, _ = convert.StringAddString(tfd.TotalBuy, tfd.TotalSell)
	tfd.Total, _ = convert.StringTo8Bit(tfd.Total)
	tfd.TotalSell, _ = convert.StringTo8Bit(tfd.TotalSell)
	tfd.TotalBuy, _ = convert.StringTo8Bit(tfd.TotalBuy)
	return mList, tfd, nil
}

//李宇舶 写的
//func (tk *TokenDailySheet) TimingFunc(begin, end int64) {
//	//g:=make([]*Trade,0)
//	//buy
//	fmt.Println("定时任务开始--------------------------------------->")
//	fmt.Println(time.Now().Unix())
//	engine := utils.Engine_token
//	// 统计买的手续费
//	sql := fmt.Sprintf("select sum(num) as a,sum(fee) as b ,sum(fee_cny) as c ,sum(total_cny) as d,token_admission_id  from trade where deal_time>=%d and deal_time<%d  and opt=1 group by token_admission_id", begin, end)
//	r, err := engine.Query(sql)
//	if err != nil {
//		utils.AdminLog.Errorln(err.Error())
//		return
//	}
//
//	l := make(map[int]*TokenDailySheet)
//
//	for _, v := range r {
//		h := &TokenDailySheet{}
//		t, ok := v["token_admission_id"]
//		if !ok {
//			utils.AdminLog.Errorln("ok u")
//		}
//
//		a, ok := v["a"]
//		if !ok {
//			utils.AdminLog.Errorln("ok a")
//		}
//		b, ok := v["b"]
//		if !ok {
//			utils.AdminLog.Errorln("ok b")
//		}
//		c, ok := v["c"]
//		if !ok {
//			utils.AdminLog.Errorln("ok c")
//		}
//		d, ok := v["d"]
//		if !ok {
//			utils.AdminLog.Errorln("ok d")
//		}
//
//		h.TokenId = tk.BytesToIntAscii(t)
//		h.BuyTotal = tk.BytesToInt64Ascii(a)
//		h.FeeBuyTotal = tk.BytesToInt64Ascii(b)
//		h.FeeBuyCny = tk.BytesToInt64Ascii(c)
//		h.BuyTotalCny = tk.BytesToInt64Ascii(d)
//		l[h.TokenId] = h
//	}
//
//	// 统计卖的手续费
//	sql = fmt.Sprintf("select token_id, sum(num) as a,sum(fee) as b ,sum(fee_cny) as c ,sum(total_cny) as d,token_admission_id  from trade where deal_time>=%d and deal_time<%d  and opt=2 group by token_admission_id", begin, end)
//	r, err = engine.Query(sql)
//	if err != nil {
//		utils.AdminLog.Errorln(err.Error())
//		return
//	}
//
//	for _, v := range r {
//		h := &TokenDailySheet{}
//		t, _ := v["token_admission_id"]
//		a, _ := v["a"]
//		b, _ := v["b"]
//		c, _ := v["c"]
//		d, _ := v["d"]
//
//		t_ := convert.BytesToIntAscii(t)
//		m, ok := l[t_]
//		if !ok {
//			h.TokenId = convert.BytesToIntAscii(t)
//			h.SellTotal = convert.BytesToInt64Ascii(a)
//			h.FeeSellTotal = convert.BytesToInt64Ascii(b)
//			h.FeeSellCny = convert.BytesToInt64Ascii(c)
//			h.SellTotalCny = convert.BytesToInt64Ascii(d)
//			l[h.TokenId] = h
//		} else {
//			m.SellTotal = convert.BytesToInt64Ascii(a)
//			m.FeeSellTotal = convert.BytesToInt64Ascii(b)
//			m.FeeSellCny = convert.BytesToInt64Ascii(c)
//			m.SellTotalCny = convert.BytesToInt64Ascii(d)
//		}
//		h.Date = end
//	}
//
//	result := &struct {
//		Balance int64
//		Frozeen int64
//	}{}
//
//	fmt.Println("len l:", len(l))
//
//	for _, v := range l {
//		p := time.Unix(begin, 0).Format("2006-01-02 ")
//		utils.AdminLog.Printf("insert into token_id %d,time %s", v.TokenId, p)
//		v.Date = begin
//		sql := fmt.Sprintf("SELECT SUM(t.balance) AS balance ,SUM(t.frozen) AS frozeen FROM `user_token` t WHERE t.token_id=%d", v.TokenId)
//		_, err := engine.Table("user_token").SQL(sql).Get(result)
//		if err != nil {
//			continue
//		}
//		v.FrozenAll = result.Frozeen
//		v.BalanceAll = result.Balance
//
//		fmt.Println("================================= v ===============================")
//		fmt.Println("v:", v.TokenId, v.SellTotalCny, v.BuyTotalCny)
//		/*
//			_, err = engine.Table("token_daily_sheet").Cols("token_id", "fee_buy_cny", "fee_buy_total", "fee_sell_cny", "fee_sell_total", "buy_total", "sell_total_cny", "sell_total", "date", "balance_all", "frozen_all").InsertOne(v)
//			if err != nil {
//				utils.AdminLog.Errorln(err.Error())
//				return
//			}
//		*/
//
//		//  判断当前id这个date是否已经统计
//		tdsheet := TokenDailySheet{TokenId: v.TokenId, Date: v.Date}
//		isExistSql := " SELECT token_id, `date` FROM  g_token.`token_daily_sheet`   WHERE token_id=? AND `date`=?"
//		has, err := engine.Table("token_daily_sheet").SQL(isExistSql, v.TokenId, v.Date).Exist(&tdsheet)
//		if err != nil {
//			fmt.Println(err)
//			utils.AdminLog.Errorln(err)
//		}
//		if has {
//			fmt.Println("exists:", v.TokenId, v.Date)
//			utils.AdminLog.Infoln("exists:", v.TokenId, v.Date)
//			continue
//		}
//		// 不存在，则插入
//		newSql := "INSERT INTO `token_daily_sheet` (`token_id`,`fee_buy_cny`,`fee_buy_total`,`fee_sell_cny`,`fee_sell_total`,`buy_total`,`sell_total_cny`,`sell_total`,`balance_all`,`frozen_all`,`date`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
//		_, err = engine.Exec(newSql, v.TokenId, v.FeeBuyCny, v.FeeBuyTotal, v.FeeSellCny, v.FeeSellTotal, v.BuyTotal, v.SellTotalCny, v.SellTotal, v.BalanceAll, v.FrozenAll, v.Date)
//		if err != nil {
//			fmt.Println(err)
//			utils.AdminLog.Errorln(err)
//			return
//		}
//	}
//	fmt.Println("successful!!!!")
//}

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
		if v.Id==0{
			continue
		}
		buy, err := new(Trade).Get(v.Id, begin,end, 1)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		buyList = append(buyList, *buy)

		sell, err := new(Trade).Get(v.Id, begin,end, 2)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		sellList = append(sellList, *sell)

	}


	//插入数据
	for _, v := range buyList {
		fmt.Println("buy")
		tdsheet := TokenDailySheet{TokenId: v.TokenId, Date: v.Date}
		isExistSql := " SELECT token_id, `date` FROM  g_token.`token_daily_sheet`   WHERE token_id=? AND `date`=?"
		has, err := engine.Table("token_daily_sheet").SQL(isExistSql, v.TokenId, v.Date).Exist(&tdsheet)
		if err != nil {
			fmt.Println(err)
			utils.AdminLog.Errorln(err)
		}
		tdsheet.TokenId =v.TokenId
		tdsheet.FeeBuyTotal= v.FeeTotal
		tdsheet.FeeBuyCny=v.FeeTotalCny
		tdsheet.BuyTotal=v.Total
		tdsheet.BuyTotalCny=v.TotalCny
		tdsheet.Date = v.Date
		if has {

			_, err := engine.Cols("fee_buy_total", "fee_buy_cny", "buy_total", "buy_total_cny").Where("token_id=? and date=?",tdsheet.TokenId,tdsheet.Date).Update(&tdsheet)
			if err != nil {
				utils.AdminLog.Info("buyList定时任务执行更新数据失败", err.Error())
				fmt.Println(err.Error())

			}
		} else {
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
		has, err := engine.Table("token_daily_sheet").SQL(isExistSql, v.TokenId, v.Date).Exist(&tdsheet)
		if err != nil {
			fmt.Println(err)
			utils.AdminLog.Errorln(err)
		}
		tdsheet.TokenId= v.TokenId
		tdsheet.FeeSellCny = v.FeeTotalCny
		tdsheet.FeeSellTotal =v.FeeTotal
		tdsheet.SellTotal = v.Total
		tdsheet.SellTotalCny=v.TotalCny
		tdsheet.Date =v.Date
		if has {
			_, err := engine.Cols("fee_sell_total", "fee_sell_cny", "sell_total", "sell_total_cny").Where("token_id=? and date=?",tdsheet.TokenId,tdsheet.Date).Update(&tdsheet)
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
	for _,vu:=range tokenIdList{
		fmt.Println("balance_all")
		balance,free,err:=new(UserCurrency).GetUserCurrencyBalanceAndFree(vu.Id)
		if err!=nil{
			utils.AdminLog.Info("sellList定时任务执行更新数据失败", err.Error())
			fmt.Println(err.Error())
			continue
		}
		tdsheet := TokenDailySheet{TokenId: int64(vu.Id), Date: begin}
		isExistSql := " SELECT token_id, `date` FROM  g_token.`token_daily_sheet`   WHERE token_id=? AND `date`=?"
		has, err := engine.Table("token_daily_sheet").SQL(isExistSql, vu.Id, begin).Exist(&tdsheet)
		if err != nil {
			fmt.Println(err)
			utils.AdminLog.Errorln(err)
		}
		tdsheet.TokenId=int64(vu.Id)
		tdsheet.FrozenAll = free
		tdsheet.BalanceAll = balance
		tdsheet.Date = begin
		if has{
			_,err:=engine.Cols("balance_all","frozen_all").Where("token_id=? and date=?",tdsheet.TokenId,tdsheet.Date).Update(&tdsheet)
			if err!=nil{
				fmt.Println(err)
				utils.AdminLog.Errorln(err)
			}
		}else{
			_,err:=engine.Cols("token_id","balance_all","frozen_all","date").InsertOne(&tdsheet)
			if err != nil {
				fmt.Println(err)
				utils.AdminLog.Errorln(err)
			}
		}

	}
	//法币 冻结
	//更新 字段 BalanceAll FrozenAll
	for _,vu:=range tokenIdList{
		fmt.Println("balance_all")
		balance,free,err:=new(UserToken).GetUserTokenBalance(vu.Id)
		if err!=nil{
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
		tdsheet.TokenId=int64(vu.Id)
		fmt.Println("前",tdsheet.FrozenAll,tdsheet.BalanceAll)
		tdsheet.FrozenAll,_ = convert.Int64AddInt64(tdsheet.FrozenAll,free)
		tdsheet.BalanceAll, _= convert.Int64AddInt64(tdsheet.BalanceAll,balance)
		fmt.Println("后",tdsheet.FrozenAll,tdsheet.BalanceAll)
		tdsheet.Date = begin
		if has{
			_,err:=engine.Cols("balance_all","frozen_all").Where("token_id=? and date=?",tdsheet.TokenId,tdsheet.Date).Update(&tdsheet)
			if err!=nil{
				fmt.Println(err)
				utils.AdminLog.Errorln(err)
			}
		}else{
			_,err:=engine.Cols("token_id","balance_all","frozen_all","date").InsertOne(&tdsheet)
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
	i:=0
	for BeginTime < unix {
		i++
		fmt.Println("第",i,"次循环")
		t.TimingFuncNew(BeginTime, BeginTime+86400)//test
		BeginTime += 86400
	}

}
