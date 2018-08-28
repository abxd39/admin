package models

import (
	"admin/utils"
	"fmt"
)

type TokenFreeHistory struct {
	Id       int64     `xorm:"not null pk autoincr INT(18)"   json:"id"`
	TokenId  int64     `xorm:"INT(11)"       json:"token_id"`
	Opt      int32     `xorm:"INT(11)"       json:"opt"`
	Type     int32     `xorm:"INT(11)"       json:"type"`
	Num      int64     `xorm:"BIGINT(20)"    json:"num"`
	CreatedTime  int64  `xorm:"BIGINT(20)"   json:"created_time"`
	Ukey      string    `xorm:"VARCHAR(128)"  json:"ukey"`
	Uid      int64      `xorm:"BIGINT(20)"    json:"uid"`
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