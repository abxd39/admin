package models

import (
	"admin/utils"
	"fmt"
)

//type TokenFreeHistory struct {
//	Id       int64     `xorm:"not null pk autoincr INT(18)"   json:"id"`
//	TokenId  int64     `xorm:"INT(11)"       json:"token_id"`
//	Opt      int32     `xorm:"INT(11)"       json:"opt"`
//	Type     int32     `xorm:"INT(11)"       json:"type"`
//	Num      int64     `xorm:"BIGINT(20)"    json:"num"`
//	CreatedTime  int64  `xorm:"BIGINT(20)"   json:"created_time"`
//	Ukey      string    `xorm:"VARCHAR(128)"  json:"ukey"`
//	Uid      int64      `xorm:"BIGINT(20)"    json:"uid"`
//}

type TokenFreeHistory struct {
	BaseModel  `xorm:"-"`
	Id          int64  `xorm:"pk autoincr BIGINT(18)" json:"id"`
	TokenId     int    `xorm:"comment('代币ID') unique(ckey) INT(11)" json:"token_id"`
	Opt         int    `xorm:"comment('操作方向1加2减') INT(11)" json:"opt"`
	Type        int    `xorm:"comment('流水类型1区块 2委托 3注册奖励 4邀请奖励 5撤销委托 6交易入账 7冻结退回 8系统自动退小额余额 9 交易确认扣减冻结数量 10划转到法币 11划转到币币 18-提币手续费') INT(11)" json:"type"`
	Num         int64  `xorm:"comment('数量') BIGINT(20)" json:"num"`
	CreatedTime int64  `xorm:"comment('操作时间') BIGINT(20)" json:"created_time"`
	Ukey        string `xorm:"comment('联合key') unique(ckey) VARCHAR(128)" json:"ukey"`
	Uid         int64  `xorm:"comment('用户uid') BIGINT(20)" json:"uid"`
	Uid1        int64  `xorm:"comment('用户uid') BIGINT(20)" json:"uid_1"`
}



type AddFreeHistory struct {
	TotalAddFree    string   `json:"total_add_free"`
	TokenId         int32    `json:"token_id"`
	Opt             int32    `json:"opt"`
}

type DelFreeHistory struct {
	TotalDelFree   string    `json:"total_del_free"`
	TokenId        int32     `json:"token_id"`
	Opt            int32     `json:"opt"`
}



func (this *TokenFreeHistory) GetFreeByTokenIds(tokenIdList []int32) (addFreeList []AddFreeHistory, delFreeList []DelFreeHistory, err error) {

	engine:= utils.Engine_token
	// 统计opt=1加的
	addSql := "SELECT token_id,opt,SUM(num) as total_add_free  FROM g_token.`token_free_history`  WHERE opt=1   GROUP BY token_id "
	err = engine.In("token_id", tokenIdList).SQL(addSql).Find(&addFreeList)
	if err != nil {
		fmt.Println(err)
	}

	// 统计opt=2减的
	delSql := "SELECT token_id,opt,SUM(num) as total_del_free FROM g_token.`token_free_history`  WHERE  opt=2   GROUP BY token_id"
	err = engine.In("token_id", tokenIdList).SQL(delSql).Find(&delFreeList)

	if err != nil {
		fmt.Println(err)
	}
	
	fmt.Println("err")
	return

}

