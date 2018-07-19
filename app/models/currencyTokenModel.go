package models

import (
	"admin/utils"
)

type Tokens struct {
	Id     int    `xorm:"not null pk autoincr INT(10)"`
	Name   string `xorm:"not null default '' comment('虚拟货币名称') VARCHAR(20)"`
	CnName string `xorm:"not null default '' comment('虚拟货币中文名称') VARCHAR(20)"`
}

func (t *Tokens) TableName() string {
	return "tokens"
}

//获取数字货币id及名称
func (t *Tokens) GetTokenList() ([]Tokens, error) {
	engine := utils.Engine_currency
	list := make([]Tokens, 0)
	err := engine.Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
