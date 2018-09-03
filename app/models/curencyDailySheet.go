package models

import (
	"admin/errors"
	"admin/utils"
	"time"

)

//type CurencyFeeDailySheet struct {
//	BaseModel  `xorm:"-"`
//	Id         int   `xorm:"not null pk comment('自增id') TINYINT(4)"`
//	FeeBuyCny  int64 `xorm:"not null comment('法币手买续费折合cny') BIGINT(20)"`
//	FeeSellCny int64 `xorm:"not null comment('法币手卖续费折合cny') BIGINT(20)"`
//	BalanceCny int64 `xorm:"not null comment('法币交易总额折合cny') BIGINT(20)"`
//	Date       int64 `xorm:"not null comment('日期例如20180801') BIGINT(10)"`
//}

type CurrencyDailySheet struct {
	BaseModel       `xorm:"-"`
	Id              int    `xorm:"not null pk autoincr comment('自增id') TINYINT(4)"`
	TokenId         int    `xorm:"not null comment('币种ID') INT(11)"`
	SellTotal       int64  `xorm:"not null default 0 comment('法币卖出总数') BIGINT(20)"`
	SellCny         int64  `xorm:"not null comment('法币卖出总额折合cny') BIGINT(20)"`
	BuyTotal        int64  `xorm:"not null default 0 comment('法币买入总数') BIGINT(20)"`
	BuyCny          int64  `xorm:"not null default 0 comment('法币买入总额折合cny') BIGINT(20)"`
	FeeSellTotal    int64  `xorm:"not null default 0 comment('法币卖出手续费总数') BIGINT(20)"`
	FeeSellCny      int64  `xorm:"not null comment('法币卖出手续费折合cny') BIGINT(20)"`
	FeeBuyTotal     int64  `xorm:"not null default 0 comment('法币买入手续费总数') BIGINT(20)"`
	FeeBuyCny       int64  `xorm:"not null default 0 comment('法币买入手续费折合cny') BIGINT(20)"`
	BuyTotalAll     int64  `xorm:"not null comment('累计买入总额') BIGINT(20)"`
	BuyTotalAllCny  int64  `xorm:"not null comment('累计买入总额折合cny') BIGINT(20)"`
	SellTotalAll    int64  `xorm:"not null comment('累计卖出总额') BIGINT(20)"`
	SellTotalAllCny int64  `xorm:"not null comment('累计卖出总额折合') BIGINT(20)"`
	Total           int64  `xorm:"not null comment('总数') BIGINT(20)"`
	TotalCny        int64  `xorm:"not null comment('总数折合') BIGINT(20)"`
	Date            string `xorm:"not null comment('时间戳，精确到天') BIGINT(10)"`
}





// 走势返回string，内容是int64
// 如果用int64，数据太大时xorm sum会溢出报错
type CurrencyTradeTrend struct {
	BuyTotal     string `xorm:"buy_total"`
	BuyCny       string `xorm:"buy_cny"`
	SellTotal    string `xorm:"sell_total"`
	SellCny      string `xorm:"sell_cny"`
	FeeBuyTotal  string `xorm:"fee_buy_total"`
	FeeBuyCny    string `xorm:"fee_buy_cny"`
	FeeSellTotal string `xorm:"fee_sell_total"`
	FeeSellCny   string `xorm:"fee_sell_cny"`
	Date         string `xorm:"date"`
}

// 交易趋势
func (this *CurrencyDailySheet) TradeTrendList(filter map[string]interface{}) ([]*CurrencyTradeTrend, error) {
	// 时间区间，默认最近一周
	today := time.Now().Format(utils.LAYOUT_DATE)

	loc, err := time.LoadLocation("Local")
	if err != nil {
		return nil, errors.NewSys(err)
	}
	todayTime, err := time.ParseInLocation(utils.LAYOUT_DATE, today, loc)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	dateBegin := todayTime.AddDate(0, 0, -7).Format(utils.LAYOUT_DATE)
	dateEnd := today

	// 开始查询
	session := utils.Engine_currency.Where("1=1")

	// 筛选
	if v, ok := filter["date_begin"]; ok {
		dateBegin, _ = v.(string)
	}
	if v, ok := filter["date_end"]; ok {
		dateEnd, _ = v.(string)
	}
	if v, ok := filter["token_id"]; ok {
		session.And("token_id=?", v)
	}

	var list []*CurrencyTradeTrend
	err = session.
		Table(this).
		Select("date, sum(buy_total) as buy_total, sum(buy_cny) as buy_cny"+
			",sum(sell_total) as sell_total, sum(sell_cny) as sell_cny"+
			",sum(fee_buy_total) as fee_buy_total, sum(fee_buy_cny) as fee_buy_cny"+
			",sum(fee_sell_total) as fee_sell_total, sum(fee_sell_cny) as fee_sell_cny").
		And("date>=?", dateBegin+" 00:00:00").
		And("date<=?", dateEnd+" 00:00:00").
		GroupBy("date").
		OrderBy("date ASC	").
		Find(&list)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	return list, nil
}
