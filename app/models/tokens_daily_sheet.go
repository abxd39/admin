package models

import (
	"time"

	"admin/errors"
	"admin/utils"
	"fmt"
)

type TokensDailySheet struct {
	Id                  int32  `xorm:"id"`
	TokenId             int64  `xorm:"token_id"`
	TokenName           string `xorm:"token_name"`
	TokenTotal          int64  `xorm:"token_total"`
	TokenFrozenTotal    int64  `xorm:"token_frozen_total"`
	CurrencyTotal       int64  `xorm:"currency_total"`
	CurrencyFrozenTotal int64  `xorm:"currency_frozen_total"`
	Date                string `xorm:"date"`
}

func (*TokensDailySheet) TableName() string {
	return "g_common.tokens_daily_sheet"
}

// 走势返回string，内容是int64
// 如果用int64，数据太大时xorm sum会溢出报错
type TokensNumTrend struct {
	TokenTotal    string `xorm:"token_total"`
	CurrencyTotal string `xorm:"currency_total"`
	Total         string `xorm:"total"`
	Date          string `xorm:"date"`
}

//币种数量走势
func (t *TokensDailySheet) NumTrend(filter map[string]interface{}) ([]*TokensNumTrend, error) {
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
	session := utils.Engine_common.Where("1=1")

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

	var list []*TokensNumTrend
	err = session.
		Table(t).
		Select("date, sum(token_total+token_frozen_total) as token_total, sum(currency_total+currency_frozen_total) as currency_total").
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

// 币种数量汇总，基于g_token/user_token、g_currency/user_currency表
func (t *TokensDailySheet) DoDailySheet(today string) error {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		utils.AdminLog.Error("【币种数量日汇总】loc err：", err.Error())
		return errors.NewSys(err)
	}
	todayTime, err := time.ParseInLocation(utils.LAYOUT_DATE, today, loc)
	if err != nil {
		utils.AdminLog.Error("【币种数量日汇总】todayTime err：", err.Error())
		return errors.NewSys(err)
	}

	yesterdayTime := todayTime.AddDate(0, 0, -1)
	yesterdayDate := yesterdayTime.Format(utils.LAYOUT_DATE)

	// 开始汇总
	// 1.币币余额
	// 1.1.获取汇总数据
	type TokenSum struct {
		TokenId          int32  `xorm:"token_id"`
		TokenName        string `xorm:"token_name"`
		TokenTotal       int64  `xorm:"token_total"`
		TokenFrozenTotal int64  `xorm:"token_frozen_total"`
	}
	var tokenSumList []*TokenSum
	err = utils.Engine_token.SQL(fmt.Sprintf("SELECT token_id, token_name, SUM(balance) token_total, SUM(frozen) token_frozen_total FROM %s GROUP BY token_id",
		new(UserToken).TableName())).Find(&tokenSumList)
	if err != nil {
		return errors.NewSys(err)
	}

	// 1.2.写入汇总表
	for _, v := range tokenSumList {
		utils.Engine_common.Exec(fmt.Sprintf("INSERT INTO %s (token_id, token_name, token_total, token_frozen_total, date)"+
			" VALUES (%d, '%s', %d, %d, '%s') ON DUPLICATE KEY UPDATE token_total=%[4]d, token_frozen_total=%[5]d",
			t.TableName(), v.TokenId, v.TokenName, v.TokenTotal, v.TokenFrozenTotal, yesterdayDate))
	}

	// 2.代币余额
	// 2.1.获取汇总数据
	type CurrencySum struct {
		TokenId             int32  `xorm:"token_id"`
		TokenName           string `xorm:"token_name"`
		CurrencyTotal       int64  `xorm:"currency_total"`
		CurrencyFrozenTotal int64  `xorm:"currency_frozen_total"`
	}
	var currencySumList []*CurrencySum
	err = utils.Engine_currency.SQL(fmt.Sprintf("SELECT token_id, token_name, SUM(balance) currency_total, SUM(freeze) currency_frozen_total FROM %s GROUP BY token_id",
		new(UserCurrency).TableName())).Find(&currencySumList)

	// 2.2.写入汇总表
	for _, v := range currencySumList {
		utils.Engine_common.Exec(fmt.Sprintf("INSERT INTO %s (token_id, token_name, currency_total, currency_frozen_total, date)"+
			" VALUES (%d, '%s', %d, %d, '%s') ON DUPLICATE KEY UPDATE currency_total=%[4]d, currency_frozen_total=%[5]d",
			t.TableName(), v.TokenId, v.TokenName, v.CurrencyTotal, v.CurrencyFrozenTotal, yesterdayDate))
	}

	return nil
}
