package models

import (
	"admin/errors"
	"admin/utils"
)

type Tokens struct {
	BaseModel `xorm:"-"`
	Id        int    `xorm:"not null pk autoincr INT(11)" json:"id" `
	Mark      string `xorm:"not null default '' comment('货币名称') VARCHAR(10)" json:"mark" `
	Detail    string `xorm:"not null default '' comment('法币交易说明') VARCHAR(255)" json:"detail"`
	Logo      string `xorm:"comment('货币logobase64') TINYTEXT"  json:"logo" `
	//Decimal               int     `xorm:"default 1 comment('精度 1个eos最小精度的10的18次方') INT(11)" form:"decimal"  binding:"required" json:"decimal"`
	Status               int     `xorm:"not null comment('是否可用 1 可用2 不可用') TINYINT(4)"  json:"status" `
	InTokenMark          int     `xorm:"not null comment('') TINYINT(4)"  json:"in_mark"`
	InTokenLeastBalance  int64   `xorm:"not null comment('') BIGINT(20)" json:"in_least_balance"`
	OutTokenMark         int     `xorm:"not null comment('') TINYINT(4)"  json:"out_mark"`
	OutTokenLeastBalance int64   `xorm:"not null comment('') BIGINT(20)"  json:"out_least_balance"`
	OutTokenFee          float32 `xorm:"not null comment('') FLOAT" json:"out_fee"`
	InRemarks            string  `xorm:"not null default '' comment('法币交易说明') VARCHAR(255)"  json:"in_remarks"`
	OutRemarks           string  `xorm:"not null default '' comment('法币交易说明') VARCHAR(255)"  json:"out_remarks"`
}

type TokensGroup struct {
	Tokens `xorm:"extends"`
	InLeast float64 `xorm:"-"`
	OutLeast float64 `xorm:"-"`
}

func(t*Tokens)TableName()string{
	return "tokens"
}
//添加 删除 修改

func (t *Tokens) TokensSystemAdd(tokens Tokens) error {
	engine := utils.Engine_common
	if tokens.Id != 0 {
		has, err := engine.Where("id=?", tokens.Id).Get(t)
		if err != nil {
			return err
		}
		if !has {
			return errors.New("no exists!!")
		}
		_, err = engine.Where("id=?", tokens.Id).Update(&Tokens{
			Mark:                 tokens.Mark,
			Detail:               tokens.Detail,
			Logo:                 tokens.Logo,
			Status:               tokens.Status,
			InTokenMark:          tokens.InTokenMark,
			InTokenLeastBalance:  tokens.InTokenLeastBalance,
			OutTokenMark:         tokens.OutTokenMark,
			OutTokenLeastBalance: tokens.OutTokenLeastBalance,
			OutTokenFee:          tokens.OutTokenFee,
		})
		if err != nil {
			return err
		}
		return nil
	}

	_, err := engine.InsertOne(tokens)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tokens) GetSystemList(page, rows, status, in, out int, name string) (*ModelList, error) {
	engine := utils.Engine_common
	query := engine.Desc("id")

	if status != 0 {
		query = query.Where("status=?", status)
	}
	if in != 0 {
		query = query.Where("in_token_mark=?", in)
	}
	if out != 0 {
		query = query.Where("out_token_mark=?", out)
	}
	if name != `` {
		query = query.Where("mark=?", name)
	}
	countQuery := *query
	count, err := countQuery.Count(&Tokens{})
	if err != nil {
		return nil, err
	}
	offset, mList := t.Paging(page, rows, int(count))
	list := make([]TokensGroup, 0)
	err = query.Limit(mList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	for i,v:=range list{
		list[i].InLeast = t.Int64ToFloat64By8Bit(v.InTokenLeastBalance)
		list[i].OutLeast =t.Int64ToFloat64By8Bit(v.OutTokenLeastBalance)
	}
	mList.Items = list
	return mList, nil
}

func (t *Tokens) GetSystem(id int) (*TokensGroup, error) {
	engine := utils.Engine_common
	tg:=new(TokensGroup)
	has, err := engine.Where("id=?", id).Get(tg)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("no exists")
	}
	tg.InLeast = t.Int64ToFloat64By8Bit(tg.InTokenLeastBalance)
	tg.OutLeast =t.Int64ToFloat64By8Bit(tg.OutTokenLeastBalance)
	return tg, nil
}

func (t *Tokens) DeleteSystem(id int) error {
	engine := utils.Engine_common
	has, err := engine.Where("id=?", id).Exist(&Tokens{})
	if err != nil {
		return err
	}
	if !has {
		return errors.New("not exists!!")
	}
	_, err = engine.Where("id=?", id).Delete(&Tokens{})
	if err != nil {
		return err
	}
	return nil
}
