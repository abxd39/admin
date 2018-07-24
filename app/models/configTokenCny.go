package models

import (
	"admin/utils"
	"admin/errors"
)

//汇率
type ConfigTokenCny struct {
	TokenId int   `xorm:"not null pk comment(' 币类型') INT(10)"`
	Price   int64 `xorm:"comment('人民币价格') BIGINT(20)"`
}

func (c*ConfigTokenCny)GetPrice(id int)(int64,error){
	engine:=utils.Engine_token
	query:=engine.Desc("token_id")
	if id!=0{
		query = query.Where("token_id=?",id)
	}
	has,err:=query.Get(c)
	if err!=nil{
		return 0,err
	}
	if !has{
		return 0,errors.New("token id not exits !!")
	}
	return  c.Price,nil
}
