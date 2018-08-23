package models

import (
	"time"

	"admin/errors"
	"admin/utils"
)

type TokensDailySheet struct {
	Id            int32  `xorm:"id"`
	TokenId       int64  `xorm:"token_id"`
	TokenTotal    int64  `xorm:"token_total"`
	CurrencyTotal int64  `xorm:"currency_total"`
	Total         int64  `xorm:"total"`
	Date          string `xorm:"date"`
}

func (*TokensDailySheet) TableName() string {
	return "g_common.tokens_daily_sheet"
}

//币种数量走势
func (t *TokensDailySheet) NumTrend(filter map[string]interface{}) ([]*TokensDailySheet, error) {
	// 时间区间，默认最近一周
	today := time.Now().Format(utils.LAYOUT_DATE)
	todayTime, _ := time.Parse(utils.LAYOUT_DATE, today)

	dateBegin := todayTime.AddDate(0, 0, -6).Format(utils.LAYOUT_DATE)
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

	var list []*TokensDailySheet
	err := session.
		Select("date, sum(token_total) as token_total, sum(currency_total) as currency_total, sum(total) as total").
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
