package models

import (
	"admin/utils"
)

type MoneyRecord struct {
	BaseModel   `xorm:"-"`
	UserInfo    `xorm:"-"`
	Id          int64  `xorm:"pk autoincr BIGINT(20)" json:"id"`
	Uid         int    `xorm:"comment('用户ID') unique(hash_index) INT(11)" json:"uid"`
	TokenId     int    `xorm:"comment('代币ID') INT(11)" json:"token_id"`
	Ukey        string `xorm:"comment('联合key') unique(hash_index) VARCHAR(128)" json:"ukey"`
	Type        int    `xorm:"comment('流水类型1区块2委托') INT(11)" json:"type"`
	Opt         int    `xorm:"comment('操作方向1加2减') unique(hash_index) TINYINT(4)" json:"opt"`
	Num         int64  `xorm:"comment('数量') BIGINT(20)" json:"num"`
	Surplus     int64  `xorm:"comment('余额') BIGINT(20)" json:"surplus"`
	CreatedTime int64  `xorm:"comment('操作时间') BIGINT(20)" json:"created_time"`
}

func (m *MoneyRecord) TableName() string {
	return "money_record"
}

func (m *MoneyRecord) GetMoneyList(page, rows int, uid []uint64) (*ModelList, error) {
	engine := utils.Engine_token
	query := engine.Desc("id")
	query = query.In("uid", uid)
	tquery := *query
	count, err := tquery.Count(&MoneyRecord{})
	if err != nil {
		return nil, err
	}
	offset, modelList := m.Paging(page, rows, int(count))
	query.Limit(modelList.PageSize, offset)

	list := make([]MoneyRecord, 0)
	err = query.Find(&list)
	if err != nil {
		return nil, err
	}
	modelList.Items = list
	return modelList, nil
}

func (m *MoneyRecord) GetMoneyListForDateOrType(page, rows, ty int, date uint64) (*ModelList, error) {
	engine := utils.Engine_token
	query := engine.Desc("id")
	if ty != 0 {
		query = query.Where("opt=?", ty)
	}
	if date != 0 {
		query = query.Where("created_time BETWEEN ? AND ?", date, date+86400)
	}
	tempQuery := *query
	count, err := tempQuery.Count(&MoneyRecord{})
	if err != nil {
		return nil, err
	}
	offset, modelList := m.Paging(page, rows, int(count))
	list := make([]MoneyRecord, 0)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	modelList.Items = list
	return modelList, nil

}
