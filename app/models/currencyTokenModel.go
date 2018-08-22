package models

import (
	"admin/utils"
	"fmt"
)

// 货币类型表
type CommonTokens struct {
	Id   uint32 `xorm:"not null pk autoincr INT(10)" json:"Id"`
	Name string `xorm:"VARBINARY(20)" json:"CnName"` // 货币中文名
	Mark string `xorm:"VARBINARY(20)" json:"Name"`   // 货币标识
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



func (t *CommonTokens) GetTokenPage(page, rows int, tokenid int32) (ctoks []CommonTokens,total int64,err error) {
	engine := utils.Engine_common
	query := engine.Asc("id")
	ctk := new(CommonTokens)
	if tokenid != 0{
		query = query.Where("id=?", tokenid)
	}
	tmpQuery := *query
	countQuery := &tmpQuery
	err = query.Limit(int(rows), (int(page)-1)*int(rows)).Find(&ctoks)
	if err != nil {
		fmt.Println(err)
		return
	}

	total, _ = countQuery.Count(ctk)
	return
}


func (t *CommonTokens) GetTokenByTokenIds(tokenList []int32)(ctoks []CommonTokens, err error) {
	engine := utils.Engine_common
	err = engine.In("id", tokenList).Find(&ctoks)
	return
}