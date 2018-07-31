package models

import "fmt"
import "admin/utils"
import "errors"

//bibi委托表
type EntrustDetail struct {
	BaseModel   `xorm:"-"`
	EntrustId   string `xorm:"not null pk comment('委托记录表（委托管理）') VARCHAR(64)"`
	Uid         int64  `xorm:"not null comment('用户id') BIGINT(32)"`
	Symbol      string `xorm:"comment('队列') VARCHAR(64)"`
	TokenId     int    `xorm:"not null comment('货币id') INT(32)"`
	AllNum      int64  `xorm:"not null comment('总数量') BIGINT(20)"`
	SurplusNum  int64  `xorm:"not null comment('剩余数量') BIGINT(20)"`
	Price       int64  `xorm:"not null comment('实际平均价格(卖出价格）') BIGINT(20)"`
	Opt         int    `xorm:"not null comment('类型 卖出单1 还是买入单0') TINYINT(4)"`
	Type        int    `xorm:"comment('交易类型') TINYINT(4)"`
	OnPrice     int64  `xorm:"not null comment('委托价格(挂单价格全价格 卖出价格是扣除手续费的）') BIGINT(20)"`
	Fee         int64  `xorm:"not null comment('手续费比例') BIGINT(20)"`
	States      int    `xorm:"not null comment('0是挂单，1是部分成交,2成交， 3撤销') TINYINT(4)"`
	CreatedTime int64  `xorm:"not null comment('添加时间') BIGINT(10)"`
	Mount       int64  `xorm:"comment('总金额') BIGINT(20)"`
}


func (this *EntrustDetail) IsExist(symbol string) (bool, error) {
	engine := utils.Engine_token
	query := engine.Desc("entrust_id")
	return query.Where("symbol=?", symbol).Exist(&EntrustDetail{})
}

func (this *EntrustDetail) EvacuateOder(uid int, odid string) error {
	engine := utils.Engine_token
	//query := engine.Desc("")
	temp := fmt.Sprintf("uid=%d AND entrust_id =%s", uid, odid)
	query := engine.Where(temp)
	TempQuery := *query
	has, err := TempQuery.Exist(&EntrustDetail{})
	if err != nil {
		return err
	}
	if !has {
		return errors.New("订单不存在！！")
	}
	_, err = query.Update(&EntrustDetail{
		States: -1,
	})
	if err != nil {
		return err
	}
	return nil

}

func (this *EntrustDetail) GetTokenOrderList(page, rows, ad_id, status, start_t, uid int, symbo, trade_id string) (*ModelList, error) {
	engine := utils.Engine_token
	//
	query := engine.Desc("entrust_id")
	if trade_id != `` {
		query = query.Where("entrust_id=?", trade_id)
	}
	if symbo != `` {
		query = query.Where("symbol=?", symbo)
	}
	if ad_id != 0 {
		query = query.Where("opt=?", ad_id)
	}
	if status !=-1 && status!=0 {
		query = query.Where("states=?", status)
	}else {
		//-1 标识刷选未成交的订单
		query =query.Where("states=?", 0)
	}

	if start_t != 0 {
		query = query.Where("created_time  BETWEEN ? AND ? ", start_t, start_t+86400)
	}
	if uid != 0 {
		query = query.Where("uid=?", uid)
	}
	tempQuery := *query
	count, err := tempQuery.Count(&EntrustDetail{})
	if err != nil {
		return nil, err
	}
	offset, modelList := this.Paging(page, rows, int(count))
	list := make([]EntrustDetail, offset)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	modelList.Items = list
	return modelList, nil
}
