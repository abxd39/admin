package models

import (
	"admin/utils"
)

// 用户虚拟货币资产表
type UserCurrency struct {
	Id        uint64 `xorm:"not null pk autoincr INT(10)" json:"id"`
	Uid       uint64 `xorm:"INT(10)"     json:"uid"`                                          // 用户ID
	TokenId   uint32 `xorm:"INT(10)"     json:"token_id"`                                     // 虚拟货币类型
	TokenName string `xorm:"VARCHAR(36)" json:"token_name"`                                   // 虚拟货币名字
	Freeze    int64  `xorm:"BIGINT not null default 0"   json:"freeze"`                       // 冻结
	Balance   int64  `xorm:"not null default 0 comment('余额') BIGINT"   json:"balance"`        // 余额
	Address   string `xorm:"not null default '' comment('充值地址') VARCHAR(255)" json:"address"` // 充值地址
	Version   int64  `xorm:"version"`
}

//获取单个用户的所有法币资产
func (this *UserCurrency) GetCurrencyList(page, rows, uid, tokenid int) ([]UserCurrency, int, int, error) {
	engine := utils.Engine_currency
	data := new(UserCurrency)
	var beginrow int
	if rows <= 0 {
		rows = 100
	}
	if page <= 1 {
		beginrow = 0
	} else {
		beginrow = (page - 1) * rows
	}

	list := make([]UserCurrency, 0)
	//根据uid 和 token_id 查询
	var total int
	count, err := engine.Where("uid=?", uid).Count(data)
	if err != nil {
		return nil, 0, 0, err
	}
	if int(count) < rows {
		total = 1
	} else {
		total = int(count) / rows
		v := int(count) % rows
		if v != 0 {
			total = total + 1
		}
	}
	if tokenid != 0 {
		// count, err := engine.Where("uid=? AND token_id=?", uid, tokenid).Count(data)
		// if err != nil {
		// 	return nil, 0, 0, err
		// }
		err = engine.Where("uid=? AND token_id=?", uid, tokenid).Limit(rows, beginrow).Find(&list)
		if err != nil {
			return nil, 0, 0, err
		}
		return list, 1, 1, nil
	}
	err = engine.Where("uid=?", uid).Limit(rows, beginrow).Find(&list)
	if err != nil {
		return nil, 0, 0, err
	}
	return list, page, total, nil

}

func (this *UserCurrency) GetAll(uid []uint64) ([]UserCurrency, error) {
	engine := utils.Engine_currency
	list := make([]UserCurrency, 0)
	err := engine.In("uid", uid).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (this *UserCurrency) GetBalance(uid, token_id int) (*UserCurrency, error) {
	engine := utils.Engine_currency
	data := new(UserCurrency)
	_, err := engine.Where("uid=? AND token_id=?", uid, token_id).Get(data)
	if err != nil {
		return nil, err
	}
	return data, nil

}
