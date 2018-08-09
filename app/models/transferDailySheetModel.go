package models

import (
	"admin/errors"
	"admin/utils"
	"time"
)

type TransferDailySheet struct {
	BaseModel `xorm:"-"`
	Id        int32 `xorm:"id"`
	TokenId   int32 `xorm:"id"`
	Type      int8  `xorm:"type"`
	Num       int64 `xrom:"num"`
	Date      int32 `xorm:"date"`
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
		time, err := time.Parse(utils.LAYOUT_DATE, v)
		if err != nil {
			return nil, nil, errors.NewSys(err)
		}
		session.And("date>=?", time.Unix())
	}
	if v, ok := filter["date_end"]; ok {
		time, err := time.Parse(utils.LAYOUT_DATE, v)
		if err != nil {
			return nil, nil, errors.NewSys(err)
		}
		session.And("date<=?", time.Unix())
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
	err = session.Select("*").OrderBy("date DESC, token_id ASC").Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}
	modelList.Items = list

	return modelList, list, nil
}
