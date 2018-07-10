package models

import (
	"admin/utils"
)

type QuenesConfig struct {
	Id           int64  `xorm:"pk autoincr BIGINT(20)"`
	TokenId      int    `xorm:"comment('交易币') unique(union_quene_id) INT(11)"`
	TokenTradeId int    `xorm:"comment('实际交易币') unique(union_quene_id) INT(11)"`
	Switch       int    `xorm:"comment('开关0关1开') TINYINT(4)"`
	Price        int64  `xorm:"comment('初始价格') BIGINT(20)"`
	Name         string `xorm:"comment('USDT/BTC') VARCHAR(32)"`
	Scope        string `xorm:"comment('振幅') DECIMAL(6,2)"`
}

func (q *QuenesConfig) GetTokenCashList(page, rows, token_id int) ([]QuenesConfig, int, int, error) {
	engine := utils.Engine_token
	limit := 0
	if rows <= 1 {
		rows = 50
	}
	if page <= 1 {
		page = 1
	} else {
		limit = (page - 1) * rows
	}
	query := engine.Desc("id")
	if token_id != 0 {
		query = query.Where("token_id=?", token_id)
	}

	query.Limit(rows, limit)
	tquery := *query
	list := make([]QuenesConfig, 0)
	err := query.Find(&list)
	if err != nil {
		return nil, 0, 0, err
	}
	count, err := tquery.Count(&QuenesConfig{})
	if err != nil {
		return nil, 0, 0, err
	}
	total_page := int(count) / rows
	return list, total_page, int(count), nil
}
