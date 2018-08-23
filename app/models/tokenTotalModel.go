package models

import (
	"admin/utils"
	"fmt"
	"admin/apis"
	"admin/errors"
)

type UserToken struct {
	BaseModel  `xorm:"-"`
	Id         int64
	Uid        uint64 `xorm:"unique(currency_uid) INT(11)"`
	TokenId    int    `xorm:"comment('币种') unique(currency_uid) INT(11)"`
	Balance    int64  `xorm:"comment('余额') BIGINT(20)"`
	Frozen     int64  `xorm:"comment('冻结余额') BIGINT(20)"`
	Version    int    `xorm:"version"`
	FrozenCny  int64  `xorm:"default 0 BIGINT(20)"`
	BalanceCny int64  `xorm:"default 0 BIGINT(20)"`
}

//资产
type PersonalProperty struct {
	UserToken  `xorm:"extends"`
	NickName   string  `json:"nick_name"`
	Phone      string  `json:"phone"`
	Email      string  `json:"email"`
	AmountTo   string `xorm:"-"json:"amount_to"`    //折合人民币
	//BalanceCny float64 `xorm:"-" json:"balance_cny"` // 这和人民币总数
	//FrozenCny  float64 `xorm:"-" json:"frozen_cny"`
	Status     int     `json:"status"` //账号状态
	Account    string  `json:"account"`
}

// var Total []PersonalProperty

// func NewTotal() {
// 	Total = make([]PersonalProperty, 0)
// }

type DetailToken struct {
	UserToken   `xorm:"extends"`
	Mark        string  `json:"mark" `
	AmountTo    string `xorm:"-"json:"amount_to"` //折合人民币
	BalanceTrue float64 `xorm:"-" json:"balance_true"`
	FrozenTrue  float64 `xorm:"-" json:"frozen_true"`
}

func (*UserToken) TableName() string {
	return "user_token"
}

func (u *DetailToken) GetTokenDetailOfUid(page, rows, uid, tokenId int) (*ModelList, error) {
	engine := utils.Engine_token
	query := engine.Alias("dt").Where("uid=?", uid)
	query = query.Join("INNER", "g_common.tokens t", "dt.token_id = t.id")
	if tokenId != 0 {
		query = query.Where("token_id=?", tokenId)
	}
	list := make([]DetailToken, 0)
	queryCount := *query
	count, err := queryCount.Count(&DetailToken{})
	if err != nil {
		return nil, err
	}
	offset, mList := u.Paging(page, rows, int(count))
	err = query.Limit(mList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	tidList:=make([]int,0)
	for _,v:=range list{
		tidList = append(tidList,v.TokenId)
	}
	fmt.Println(tidList)
	priceList,err:=new(apis.VendorApi).GetTokenCnyPriceList(tidList)
	if err!=nil{
		utils.AdminLog.Errorln(err.Error())
		return nil,errors.New(err.Error())
	}
	//注 折合为空是因为 还没有计算折合
	for i, v := range list {

		list[i].BalanceTrue = u.Int64ToFloat64By8Bit(v.Balance)
		list[i].FrozenTrue = u.Int64ToFloat64By8Bit(v.Frozen)
		for _,pv:=range priceList{
			if v.TokenId == pv.TokenId{
				//list[i].AmountTo = u.Int64ToFloat64By8Bit(v.BalanceCny) + u.Int64ToFloat64By8Bit(v.FrozenCny)
				fmt.Println("bibi",v.Balance)
				temp :=u.Int64MulInt64By8BitString(pv.CnyPriceInt,v.Balance)
				fmt.Println(temp)
				list[i].AmountTo = temp
			}
		}

	}
	mList.Items = list

	return mList, nil
}

//所有用户 的全部币币资产
//第一步get 所有用户
func (t *PersonalProperty) TotalUserBalance(page, rows, status int, search string) (*ModelList, error) {
	fmt.Println("this is run")
	engine := utils.Engine_token
	query := engine.Alias("ut")
	query = query.Join("LEFT", "g_common.user u", "u.uid=ut.uid")
	query = query.Join("LEFT", "g_common.user_ex ex", "ex.uid=ut.uid")
	if status != 0 {
		query = query.Where("u.status=?", status)
	}
	if search != `` {
		temp := fmt.Sprintf(" concat(IFNULL(ut.`uid`,''),IFNULL(u.`phone`,''),IFNULL(ex.`nick_name`,''),IFNULL(u.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)

	}

	countQuery := *query
	count, err := countQuery.Distinct("ut.uid").Count(&PersonalProperty{})
	if err != nil {
		return nil, err
	}
	offset, mList := t.Paging(page, rows, int(count))
	list := make([]PersonalProperty, 0)
	err = query.Desc("ut.uid").Select("ut.token_id,ut.uid uid, u.phone phone,u.email email,ex.nick_name nick_name,u.status status,u.account account").GroupBy("ut.uid").Limit(mList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	//折合rmb为空是因为 数据库上就空
	//h获取rmb 价格
	//tidList:=make([]int,0)
	//for _,v:=range list{
	//	tidList = append(tidList,v.TokenId)
	//}
	//fmt.Println("uid=",tidList)
	//
	//priceList,err:=new(apis.VendorApi).GetTokenCnyPriceList(tidList)
	//if err!=nil{
	//	utils.AdminLog.Errorln(err.Error())
	//	return nil,errors.New(err.Error())
	//}
	//
	//
	//for i, v := range list {
	//	for _,vt:=range priceList{
	//		if vt.TokenId == v.TokenId {
	//			list[i].AmountTo = convert.Int64MulInt64By8BitString( v.Balance + v.Frozen,vt.CnyPriceInt)
	//			break
	//		}
	//	}
	//	fmt.Println("bibi账户折合",v.Balance,v.Frozen)
	//
	//}
	mList.Items = list
	return mList, nil
}
