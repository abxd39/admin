package models

import (
	"admin/utils"
	"fmt"
)

type TokenHistory struct {
	BaseModel  `xorm:"-"`
	Id         int    `xorm:"comment('操作序号') INT(255)"`
	Uid        int    `xorm:"comment('用户id') INT(11)"`
	TokenId    int    `xorm:"comment(' 货币类型') INT(11)"`
	Num        int64  `xorm:"comment('提现数量') BIGINT(20)"`
	Fee        int64  `xorm:"comment('手续费') BIGINT(20)"`
	Address    string `xorm:"comment('提现地址') VARCHAR(255)"`
	RecordTime int    `xorm:"comment('提交时间') INT(11)"`
	CheckTime  int    `xorm:"comment('审核时间') INT(11)"`
	AdminId    int    `xorm:"comment('审核人') INT(11)"`
	Status     int    `xorm:"comment('状态0审核中，1拒绝，2成功') INT(11)"`
	Operator   int    `xorm:"comment('操作类型0充币，1提币') INT(11)"`
}
type TokenHistoryGroup struct {
	TokenHistory `xorm:"extends"`
	Name string
}

func (this*TokenHistoryGroup)TableName()string  {
		return "token_history"
}
//p5-1-1-1提币手续费明细
func (this *TokenHistory) GetAddTakeList(page, rows, tid, uid int, date uint64) (*ModelList, error) {
	fmt.Println("p5-1-1-1提币手续费明细")
	engine := utils.Engine_token
	query := engine.Desc("token_history.id")
	query = query.Join("LEFT","tokens","tokens.id=token_history.token_id")
	if tid != 0 {
		query = query.Where("token_id=?", tid)
	}
	if uid != 0 {
		query = query.Where("uid=?", uid)
	}
	if date != 0 {
		query = query.Where("check_time BETWEEN ? AND ?", date, date+864000)
	}
	countQuery:=*query
	count,err:=countQuery.Count(&TokenHistory{})
	if err!=nil{
		return nil,err
	}

	offset,mlist:=this.Paging(page,rows,int(count))
	list:=make([]TokenHistoryGroup,offset)
	err=query.Limit(mlist.PageSize,offset).Find(&list)
	if err!=nil{
		return nil,err
	}
	mlist.Items = list
	return  mlist,nil
}
