package models

import (
	"admin/errors"
	"admin/utils"
	"fmt"
	"time"
)

type TransferDailySheet struct {
	BaseModel `xorm:"-"`
	Id        int32  `xorm:"id" json:"id"`
	TokenId   int32  `xorm:"token_id" json:"token_id"`
	Type      int8   `xorm:"type" json:"type"` // 1-划转到币币 2-划转到法币
	Num       int64  `xrom:"num" json:"num"`
	Date      string `xorm:"date" json:"date"`
}

func (TransferDailySheet) TableName() string {
	return "transfer_daily_sheet"
}

// 列表
func (t *TransferDailySheet) List(pageIndex, pageSize int, filter map[string]string) (*ModelList, []*TransferDailySheet, error) {
	session := utils.Engine_token.Where("1=1")

	// 筛选
	if v, ok := filter["type"]; ok {
		session.And("type=?", v)
	}
	if v, ok := filter["token_id"]; ok {
		session.And("token_id=?", v)
	}
	if v, ok := filter["date_begin"]; ok {
		session.And("date>=?", v)
	}
	if v, ok := filter["date_end"]; ok {
		session.And("date<=?", v)
	}

	//计算分页
	countSession := session.Clone()
	count, err := countSession.Count(t)
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}
	offset, modelList := t.Paging(pageIndex, pageSize, int(count))

	// 获取列表
	var list []*TransferDailySheet
	err = session.Select("*").OrderBy("date DESC, type ASC, token_id ASC").Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}
	modelList.Items = list

	return modelList, list, nil
}

// 汇总，基于g_currency/transfer_record、g_token/transfer_record表
// 汇总指定日期上一天
func (t *TransferDailySheet) DoDailySheet(today string) error {
	// 获取昨天
	loc, err := time.LoadLocation("Local")
	if err != nil {
		utils.AdminLog.Error("【划转日汇总】loc err：", err.Error())
		return errors.NewSys(err)
	}

	todayTime, err := time.ParseInLocation(utils.LAYOUT_DATE, today, loc)
	if err != nil {
		utils.AdminLog.Error("【划转日汇总】todayTime err：", err.Error())
		return errors.NewSys(err)
	}

	yesterdayDate := todayTime.AddDate(0, 0, -1).Format(utils.LAYOUT_DATE)
	yesterdayBeginTime, err := time.ParseInLocation(utils.LAYOUT_DATE_TIME, fmt.Sprintf("%s 00:00:00", yesterdayDate), loc)
	if err != nil {
		utils.AdminLog.Error("【划转日汇总】yesterdayBeginTime err：", err.Error())
		return errors.NewSys(err)
	}

	yesterdayEndTime, err := time.ParseInLocation(utils.LAYOUT_DATE_TIME, fmt.Sprintf("%s 23:59:59", yesterdayDate), loc)
	if err != nil {
		utils.AdminLog.Error("【划转日汇总】yesterdayEndTime err：", err.Error())
		return errors.NewSys(err)
	}

	yesterdayBeginUnix := yesterdayBeginTime.Unix()
	yesterdayEndUnix := yesterdayEndTime.Unix()

	// 开始汇总
	// 1.币币划转到法币
	// 1.1.获取汇总数据
	type TokenToCurrencySum struct {
		TokenId   int32  `xorm:"token_id"`
		TokenName string `xorm:"token_name"`
		TotalNum  int64  `xorm:"total_num"`
	}
	var tokenToCurrencySumList []*TokenToCurrencySum
	err = utils.Engine_token.SQL(fmt.Sprintf("SELECT token_id, token_name, SUM(num) AS total_num FROM %s WHERE create_time>=%d AND create_time<=%d GROUP BY token_id",
		new(TokenTransferRecordModel).TableName(), yesterdayBeginUnix, yesterdayEndUnix)).Find(&tokenToCurrencySumList)
	if err != nil {
		return errors.NewSys(err)
	}

	// 1.2.写入汇总表
	for _, v := range tokenToCurrencySumList {
		utils.Engine_token.Exec(fmt.Sprintf("INSERT INTO %s (token_id, token_name, type, num, date) VALUES (%d, '%s', 2, %d, '%s') ON DUPLICATE KEY UPDATE num=%[4]d",
			t.TableName(), v.TokenId, v.TokenName, v.TotalNum, yesterdayDate))
	}

	// 2.法币划转到币币
	// 2.1.获取汇总数据
	type CurrencyToTokenSum struct {
		TokenId   int32  `xorm:"token_id"`
		TokenName string `xorm:"token_name"`
		TotalNum  int64  `xorm:"total_num"`
	}
	var currencyToTokenSumList []*CurrencyToTokenSum
	err = utils.Engine_currency.SQL(fmt.Sprintf("SELECT token_id, token_name, SUM(num) AS total_num FROM %s WHERE create_time>=%d AND create_time<=%d GROUP BY token_id",
		new(CurrencyTransferRecordModel).TableName(), yesterdayBeginUnix, yesterdayEndUnix)).Find(&currencyToTokenSumList)

	// 2.2.写入汇总表
	for _, v := range currencyToTokenSumList {
		utils.Engine_token.Exec(fmt.Sprintf("INSERT INTO %s (token_id, token_name, type, num, date) VALUES (%d, '%s', 1, %d, '%s') ON DUPLICATE KEY UPDATE num=%[4]d",
			t.TableName(), v.TokenId, v.TokenName, v.TotalNum, yesterdayDate))
	}

	return nil
}
