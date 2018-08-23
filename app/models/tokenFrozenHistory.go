package models

import (
	"admin/utils"
	"fmt"
)

//注册奖励表
//邀请奖励表
type FrozenHistory struct {
	BaseModel    `xorm:"-"`
	Uid        int64  `xorm:"comment('用户ID') BIGINT(20)" json:"uid"`
	Ukey       string `xorm:"comment('流水ID') unique(uni) VARCHAR(128)" json:"ukey"`
	Opt        int    `xorm:"comment('操作类型') INT(11)" json:"opt"`
	TokenId    int    `xorm:"comment('币类型') unique(uni) INT(11)" json:"token_id"`
	Num        int64  `xorm:"comment('数量') BIGINT(20)" json:"num"`
	Type       int    `xorm:"comment('业务使用类型') unique(uni) INT(11)" json:"type"`
	CreateTime int64  `xorm:"comment('创建时间') BIGINT(20)" json:"created_time"`
	Frozen     int64  `xorm:"comment('冻结余额') BIGINT(20)" json:"frozen"`
}

type FrozenHistoryGroup struct {
	FrozenHistory `xorm:"extends"`
	UserInfo    `xorm:"extends"`
	Balance int64	`xorm:"-",json:"balance"`
}
func (f*FrozenHistoryGroup) TableName()string{
	return "frozen_history"
}

func (f*FrozenHistory)GetFrozenHistory(page,rows,typ,tid,status int,bt,et uint64, search string)(*ModelList,error){
	engine :=utils.Engine_token
	query := engine.Table("frozen_history").Alias("fh")
	query = query.Join("LEFT", "g_common.user u ", "u.uid= fh.uid")
	query = query.Join("LEFT", "g_common.user_ex ex", "fh.uid=ex.uid")
	//opt,tid 是必须的所以不可能为空
	if typ!=0{
		query=query.Where("fh.type=?",typ)
	}
	if tid !=0{
		query=query.Where("fh.token_id=? ",tid)
	}
	if bt!=0{
		if et!=0{
			query = query.Where("fh.token_id=?  AND fh.created_time BETWEEN ? AND ? ", tid,bt ,et+86400)
		}else {
			query = query.Where("fh.token_id=?  AND fh.created_time BETWEEN ? AND ? ", tid,bt ,bt+86400)
		}

	}
	if search != `` {
		temp := fmt.Sprintf(" concat(IFNULL(u.`uid`,''),IFNULL(u.`phone`,''),IFNULL(ex.`nick_name`,''),IFNULL(u.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}
	//用户状态
	if status != 0 {
		query = query.Where("u.status=?", status)
	}
	queryCount:=*query
	count,err:=queryCount.Count(&FrozenHistoryGroup{})
	if err!=nil{
		return nil,err
	}
	offset,mList:=f.Paging(page,rows,int(count))
	list:=make([]FrozenHistoryGroup,0)
	err=query.Limit(mList.PageSize,offset).Find(&list)
	if err!=nil{
		return nil,err
	}
	mList.Items =list
	return mList,nil
}
