package models

import (
	"admin/utils"
	"time"
	"strconv"
	"admin/errors"
)

type ConfigureTradeRule struct {
	Id          int   `xorm:"not null pk autoincr comment('自增id') TINYINT(4)" json:"id"`
	Cuid        int   `xorm:"not null comment('创建者uid') INT(4)" json:"cuid"`
	Muid        int   `xorm:"not null comment('修改者uid') INT(4)" json:"muid"`
	OneTradeMax int64 `xorm:"comment('一级认证单笔最大交易额限定') BIGINT(20)" json:"one_trade_max"`
	OneTotal    int64 `xorm:"comment('一级认证总交易额上限限定') BIGINT(20)" json:"one_total"`
	TwoTradeMax int64 `xorm:"comment('二级认证单笔最大交易额限定') BIGINT(20)" json:"two_trade_max"`
	TwoTotal    int64 `xorm:"comment('二级认证总交易额上限限定') BIGINT(20)" json:"two_total"`
	CreateDate  int   `xorm:"comment('创建时间') INT(11)" json:"create_date"`
	UpdateDate  int   `xorm:"comment('修改时间') INT(11)" json:"update_date"`
}

func (this *ConfigureTradeRule) AddTradeRule(c ConfigureTradeRule) error {
	engine := utils.Engine_common
	con := new(ConfigureTradeRule)
	t:=time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	value,err:=strconv.Atoi(timestamp)
	if err!=nil{
		return err
	}
	query :=engine.Where("id=1")
	has, err := query.Get(con)
	if err != nil {
		return err
	}
	if !has {
		query.InsertOne(&ConfigureTradeRule{
			Cuid:c.Cuid,
			Muid:c.Muid,
			OneTradeMax:c.OneTradeMax,
			OneTotal:c.OneTotal,
			TwoTradeMax:c.TwoTradeMax,
			TwoTotal:c.TwoTotal,
			CreateDate:value,
		})
	}else {
		_,err:=query.Update(&ConfigureTradeRule{
			Cuid:c.Cuid,
			Muid:c.Muid,
			OneTradeMax:c.OneTradeMax,
			OneTotal:c.OneTotal,
			TwoTradeMax:c.TwoTradeMax,
			TwoTotal:c.TwoTotal,
			UpdateDate:value,
		})
		if err!=nil{
			return err
		}
	}
	return nil
}

func (this*ConfigureTradeRule)GetTradeRule()(*ConfigureTradeRule,error)  {
	engine := utils.Engine_common
	con := new(ConfigureTradeRule)
	has,err:=engine.Where("id=1").Get(con)
	if err!=nil{
		return nil,err
	}
	if !has{
		return nil,errors.New("交易规则不存在，请添加！！")
	}
	return con,nil
}
