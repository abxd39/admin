package models

import (
	"admin/utils"
	"fmt"
)

type UserCurrencyHistory struct {
	BaseModel   `xorm:"-"`
	Id          int    `xorm:"not null pk autoincr comment('ID') INT(10)" json:"id"`
	Uid         int    `xorm:"not null default 0 INT(10)" json:"uid"`
	OrderId     string `xorm:"not null default '' comment('订单ID') VARCHAR(64)" json:"order_id"`
	TokenId     int    `xorm:"not null default 0 comment('货币类型') INT(10)" json:"token_id"`
	Num         int64  `xorm:"not null default 0 comment('数量') BIGINT(64)" json:"num"`
	Fee         int64  `xorm:"not null default 0 comment('手续费用') BIGINT(64)" json:"fee"`
	Surplus     int64  `xorm:"comment('账户余额') BIGINT(64)" json:"surplus"`
	Operator    int    `xorm:"not null default 0 comment('操作类型 1订单转入 2订单转出 3充币 4提币 5冻结') TINYINT(2)" json:"operator"`
	Address     string `xorm:"not null default '' comment('提币地址') VARCHAR(255)" json:"address"`
	States      int    `xorm:"not null default 0 comment('订单状态: 0删除 1待支付 2待放行(已支付) 3确认支付(已完成) 4取消') TINYINT(2)" json:"states"`
	CreatedTime string `xorm:"not null comment('创建时间') DATETIME" json:"created_time"`
	UpdatedTime string `xorm:"comment('修改间') DATETIME" json:"updated_time"`
}

func (u *UserCurrencyHistory) GetList(page, rows, ot int, date string) (*ModelList, error) {
	engine := utils.Engine_currency
	query := engine.Desc("id")
	if ot != 0 {
		query = query.Where("operator=?", ot)
	}
	if date != `` {
		sub := date[:11] + "23:59:59"
		temp := fmt.Sprintf("created_time BETWEEN ? AND ?", date, sub)
		query = query.Where(temp)
	}
	tempQuery := *query
	count, err := tempQuery.Count(&UserCurrencyHistory{})
	if err != nil {
		return nil, err
	}
	offset, modelList := u.Paging(page, rows, int(count))
	list := make([]UserCurrencyHistory, 0)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	modelList.Items = list
	return modelList, nil
}

func (u *UserCurrencyHistory) GetListForUid(page, rows int, uid []uint64) (*ModelList, error) {
	engine := utils.Engine_currency
	query := engine.Desc("id")
	query = query.In("uid", uid)
	tempQuery := *query
	count, err := tempQuery.Count(&UserCurrencyHistory{})
	if err != nil {
		return nil, err
	}
	offset, modelList := u.Paging(page, rows, int(count))
	list := make([]UserCurrencyHistory, 0)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	modelList.Items = list
	return modelList, nil

}
