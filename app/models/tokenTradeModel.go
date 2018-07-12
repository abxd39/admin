package models

import (
	"admin/utils"
	"fmt"
)

type EntrustDetail struct {
	BaseModel   `xorm:"-"`
	EntrustId   string `xorm:"not null pk comment('委托记录表（委托管理）') VARCHAR(64)"`
	Uid         uint64 `xorm:"not null comment('用户id') INT(32)"`
	TokenId     int    `xorm:"not null comment('货币id') INT(32)"`
	AllNum      int64  `xorm:"not null comment('总数量') BIGINT(20)"`
	SurplusNum  int64  `xorm:"not null comment('剩余数量') BIGINT(20)"`
	Price       int64  `xorm:"not null comment('实际价格(卖出价格）') BIGINT(20)"`
	Opt         int    `xorm:"not null comment('类型 买入单1 卖出单2 ') TINYINT(4)"`
	Type        int    `xorm:"not null comment('类型 市价委托1 还是限价委托2') TINYINT(4)"`
	OnPrice     int64  `xorm:"not null comment('委托价格(挂单价格全价格 卖出价格是扣除手续费的）') BIGINT(20)"`
	Fee         int64  `xorm:"not null comment('手续费比例') BIGINT(20)"`
	States      int    `xorm:"not null comment('状态0正常1撤单2成交') TINYINT(4)"`
	CreatedTime int    `xorm:"not null comment('添加时间') created INT(10)"`
}

func (this *EntrustDetail) GetTokenRecordList(page, rows, trade_id, trade_duad, ad_id int, start_t, end_t string) ([]EntrustDetail, int, int, error) {
	engine := utils.Engine_token
	//
	if page <= 1 {
		page = 1
	}
	if rows <= 0 {
		rows = 100
	}
	var total int
	tok := new(EntrustDetail)
	count, err := engine.Count(tok)
	if err != nil {
		return nil, 0, 0, err
	}
	if int(count) < rows {
		total = 1
	} else {
		total = int(count) / rows
		v := int(count) % rows
		if v != 0 {
			total += 1
		}
	}
	list := make([]EntrustDetail, 0)
	fmt.Printf("$$$$$$$$$$$$$$$%#v\n", rows)
	err = engine.Where("states=0").Limit(rows, (page-1)*rows).Find(&list)
	if err != nil {
		return nil, 0, 0, err
	}
	return list, total, total * rows, nil
}

func (this *EntrustDetail) GetTokenOrderList(page, rows, trade_id, trade_duad, ad_id, status int, start_t, end_t string) ([]EntrustDetail, int, int, error) {
	engine := utils.Engine_token
	//
	if page <= 1 {
		page = 1
	}
	if rows <= 0 {
		rows = 100
	}
	var total int
	tok := new(EntrustDetail)
	count, err := engine.Count(tok)
	if err != nil {
		return nil, 0, 0, err
	}
	if int(count) < rows {
		total = 1
	} else {
		total = int(count) / rows
		v := int(count) % rows
		if v != 0 {
			total += 1
		}
	}
	list := make([]EntrustDetail, 0)
	fmt.Printf("@@@@@@@@@@%#v\n", rows)
	err = engine.Limit(rows, (page-1)*rows).Find(&list)
	if err != nil {
		return nil, 0, 0, err
	}
	return list, total, total * rows, nil
}
