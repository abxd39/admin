package models

import (
	"admin/utils"
)


// 货币类型表
type CommonTokens struct {
	Id   uint32 `xorm:"not null pk autoincr INT(10)" json:"id"`
	Name string `xorm:"VARBINARY(20)" json:"cn_name"` // 货币中文名
	Mark string `xorm:"VARBINARY(20)" json:"name"`    // 货币标识
}

func (t *CommonTokens) TableName() string {
	return "tokens"
}

//获取数字货币id及名称
func (t *CommonTokens) GetTokenList() ([]CommonTokens, error) {
	engine := utils.Engine_currency
	list := make([]CommonTokens, 0)
	err := engine.Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
