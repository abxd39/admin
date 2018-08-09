package models

import "fmt"
import "admin/utils"
import (
	"errors"
	"strconv"
)

//bibi委托表
type EntrustDetail struct {
	BaseModel           `xorm:"-"`
	EntrustId           string `xorm:"not null pk comment('委托记录表（委托管理）') VARCHAR(64)"`
	Uid                 int64  `xorm:"not null comment('用户id') BIGINT(32)"`
	Symbol              string `xorm:"comment('队列') VARCHAR(64)"`
	TokenId             int    `xorm:"not null comment('货币id') INT(32)"`
	AllNum              int64  `xorm:"not null comment('总数量') BIGINT(20)"`
	SurplusNum          int64  `xorm:"not null comment('剩余数量') BIGINT(20)"`
	Price               int64  `xorm:"not null comment('实际平均价格(卖出价格）') BIGINT(20)"`
	Opt                 int    `xorm:"not null comment('类型 卖出单1 还是买入单0') TINYINT(4)"`
	Type                int    `xorm:"comment('交易类型') TINYINT(4)"`
	OnPrice             int64  `xorm:"not null comment('委托价格(挂单价格全价格 卖出价格是扣除手续费的）') BIGINT(20)"`
	FeePercent          float64  `xorm:"not null comment('手续费比例') BIGINT(20)"`//手续费比例 0.1%=0.001
	States              int    `xorm:"not null comment('0是挂单，1是部分成交,2成交， 3撤销') TINYINT(4)" json:"states"`
	CreatedTime         int64  `xorm:"not null comment('添加时间') BIGINT(10)"`
	Sum               	int64  `xorm:"comment('委托总额') BIGINT(20)"`
}


type ReturnValueOperator struct {
	EntrustDetail `xorm:"extends"`
	AllNumTrue     float64 `json:"all_num_true"`
	SurplusNumTrue float64 `json:"surplus_num_true"`//剩余
	PriceTrue      float64 `json:"price_true"`
	OnPriceTrue    float64 `json:"on_price_true"` //
	FeeTrue        float64 `json:"fee_true"` //手续费
	//MountTrue      float64 `json:"mount_true"`
	FinishCount    float64 `json:"finish_count"` //已成
	Fee int64	`json:"fee"`//手续费
}

func (this*ReturnValueOperator)TableName()string  {
	return "entrust_detail"
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

func (this *EntrustDetail) GetTokenOrderList(page, rows, adId, status, bt,et, uid int, symbol, tradeId string) (*ModelList, error) {
	engine := utils.Engine_token
	query := engine.Desc("en.entrust_id")
	query  =query.Alias("en")
	query =query.Desc("en.entrust_id ")
	query = query.Join("left","trade t","t.entrust_id = en.entrust_id ")
	query = query.Where("en.symbol=?",symbol)
	if tradeId != `` {
		query = query.Where("en.entrust_id=?", tradeId)
	}

	if adId != 0 {
		query = query.Where("en.opt=?", adId)
	}
	if status != 0 {
		query = query.Where("en.states=?", status)
	}
	if bt != 0 {
		if et!=0{
			query = query.Where("en.created_time  BETWEEN ? AND ? ", bt, et+86400)
		}else {
			query = query.Where("en.created_time  BETWEEN ? AND ? ", bt, bt+86400)
		}
	}
	if uid != 0 {
		query = query.Where("en.uid=?", uid)
	}
	fmt.Println("debug------>1110")
	tempQuery := *query
	count, err := tempQuery.Count(&ReturnValueOperator{})
	if err != nil {
		return nil, err
	}

	offset, modelList := this.Paging(page, rows, int(count))
	fmt.Println("---------------------->count=",count,"modelList.PageSize=", modelList.PageSize, "offset=?", offset)
	list := make([]ReturnValueOperator, 0)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}

	for i, v := range list {
		list[i].PriceTrue = this.Int64ToFloat64By8Bit(v.Price)
		list[i].FeeTrue = this.Decimal(this.Int64ToFloat64By8Bit(v.Fee))
		list[i].AllNumTrue = this.Int64ToFloat64By8Bit(v.AllNum)
		list[i].OnPriceTrue = this.Int64ToFloat64By8Bit(v.OnPrice)
		list[i].SurplusNumTrue = this.Int64ToFloat64By8Bit(v.SurplusNum)
		//list[i].MountTrue = this.Int64ToFloat64By8Bit(v.Mount)
		list[i].FinishCount = list[i].AllNumTrue - list[i].SurplusNumTrue
	}
	modelList.Items = list
	return modelList, nil
}

func (this *EntrustDetail)Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", value), 64)
	return value
}