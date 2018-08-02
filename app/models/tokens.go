package models

import (
	"admin/utils"
	"admin/errors"
)

type Tokens struct {
	BaseModel `xorm:"-"`
	Id                    int     `xorm:"not null pk autoincr INT(11)" binding:"required" json:"id" `
	Name                  string  `xorm:"not null default '' comment('货币名称') VARCHAR(64)" binding:"required" json:"name"`
	Detail                string  `xorm:"not null default '' comment('法币交易说明') VARCHAR(255)" binding:"required" json:"detail"`
	Logo                  string  `xorm:"comment('货币logobase64') TINYTEXT" binding:"required" json:"logo"`
	Decimal               int     `xorm:"default 1 comment('精度 1个eos最小精度的10的18次方') INT(11)" binding:"required" json:"decimal"`
	Status                int     `xorm:"not null comment('是否可用 1 可用2 不可用') TINYINT(4)" binding:"required" json:"status"`
	InTokenMark           int     `xorm:"not null comment('') TINYINT(4)" binding:"required" json:"in_token_mark"`
	OutTokenMark          int     `xorm:"not null comment('') TINYINT(4)" binding:"required" json:"out_token_mark"`
	OutTokenLeastBalance int64   `xorm:"not null comment('') BIGINT(20)" binding:"required" json:"out_token_least_balance"`
	OutTokenFee           float32 `xorm:"not null comment('') FLOAT" binding:"required" json:"out_token_fee"`
}

//添加 删除 修改

func (t*Tokens)TokensSystemAdd(tokens Tokens)(error) {
	engine:=utils.Engine_common
	if tokens.Id!=0{
		_,err:=engine.Update(tokens)
		if err!=nil{
			return err
		}
		return nil
	}
	_,err:=engine.InsertOne(tokens)
	if err!=nil{
		return err
	}
	return nil
}

func (t*Tokens)GetSystemList(page,rows,status,in,out int,name string)(*ModelList,error)  {
	engine:=utils.Engine_common
	query:=engine.Desc("id")
	if status == -1{
		query=query.Where("status=?",0)
	}
	if status ==1{
		query =query.Where("status=?",status)
	}
	if in!=0{
		query =query.Where("in_token_mark=?",in)
	}
	if out!=0{
		query =query.Where("out_token_mark=?",out)
	}
	if name!=``{
		query =query.Where("mark=?",name)
	}
	countQuery:=*query
	count,err:=countQuery.Count(&Tokens{})
	if err!=nil{
		return nil,err
	}
	offset,mList:=t.Paging(page,rows,int(count))
	list:=make([]Tokens,0)
	err=query.Limit(mList.PageSize,offset).Find(&list)
	if err!=nil{
		return nil,err
	}
	mList.Items =list
	return mList,nil
}

func (t*Tokens)GetSystem(id int)(*Tokens,error)  {
	engine:=utils.Engine_common
	has,err:=engine.Where("id=?",id).Get(t)
	if err !=nil{
		return nil,err
	}
	if !has{
		return nil,errors.New("no exists")
	}
	return t,nil
}

func (t*Tokens)DeleteSystem(id int)(error)  {
	engine:=utils.Engine_common
	_,err:=engine.Where("id=?",id).Delete(&Tokens{})
	if err !=nil{
		return err
	}
	return nil
}
