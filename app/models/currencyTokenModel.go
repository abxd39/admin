package models

import (
	"admin/utils"
)


// 货币类型表
type CommonTokens struct {
	Id   uint32 `xorm:"not null pk autoincr INT(10)" json:"Id"`
	Name string `xorm:"VARBINARY(20)" json:"CnName"` // 货币中文名
	Mark string `xorm:"VARBINARY(20)" json:"Name"`    // 货币标识
	//Detail string `xorm:"VARBINARY(20)" json:"detail"`    // 货币标识
}

func (t *CommonTokens) TableName() string {
	return "tokens"
}

//获取数字货币id及名称
func (t *CommonTokens) GetTokenList() ([]CommonTokens, error) {
	engine := utils.Engine_common
	list := make([]CommonTokens, 0)
	err := engine.Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
