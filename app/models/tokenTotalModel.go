package models

import (
	"admin/apis"
	"admin/errors"
	"admin/utils"
	"fmt"
	"strings"
"admin/utils/convert"
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
	UserToken `xorm:"extends"`
	NickName  string `json:"nick_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	AmountTo  string `xorm:"-"json:"amount_to"` //折合人民币
	BalanceTrue string `xorm:"-" json:"balance_true"` // 这和人民币总数
	FrozenTrue  string `xorm:"-" json:"frozen_true"`
	Status  int    `json:"status"` //账号状态
	Account string `json:"account"`
}

// var Total []PersonalProperty

// func NewTotal() {
// 	Total = make([]PersonalProperty, 0)
// }

type DetailToken struct {
	UserToken   `xorm:"extends"`
	Mark        string  `json:"mark" `
	AmountTo    string  `xorm:"-"json:"amount_to"` //折合人民币
	BalanceTrue float64 `xorm:"-" json:"balance_true"`
	FrozenTrue  float64 `xorm:"-" json:"frozen_true"`
}

func (*UserToken) TableName() string {
	return "user_token"
}

func (*UserToken) GetUserTokenBalance( tid int)(balance ,free int64,err error){
	engine := utils.Engine_token
	result, err := engine.Where("token_id=?", tid).SumsInt(&UserToken{}, "balance", "frozen")
	if err != nil {
		return 0, 0, err
	}
	return result[0], result[1], nil
}



//更具用戶UID获取账户余额
func  (t*UserToken)GetBalanceForUid( uid []int64)([]UserToken,error){
	engine := utils.Engine_token
	list:=make([]UserToken,0)
	err:=engine.In("uid",uid).Find(&list)
	if err!=nil{
		return nil,err
	}
	return list,nil
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
	tidList := make([]int, 0)
	for _, v := range list {
		tidList = append(tidList, v.TokenId)
	}
	fmt.Println(tidList)
	priceList, err := new(apis.VendorApi).GetTokenCnyPriceList(tidList)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return nil, errors.New(err.Error())
	}
	//注 折合为空是因为 还没有计算折合
	for i, v := range list {

		list[i].BalanceTrue = u.Int64ToFloat64By8Bit(v.Balance)
		list[i].FrozenTrue = u.Int64ToFloat64By8Bit(v.Frozen)
		for _, pv := range priceList {
			if v.TokenId == pv.TokenId {
				//list[i].AmountTo = u.Int64ToFloat64By8Bit(v.BalanceCny) + u.Int64ToFloat64By8Bit(v.FrozenCny)
				fmt.Println("bibi", v.Balance)
				temp := u.Int64MulInt64By8BitString(pv.CnyPriceInt, v.Balance)
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
func (t *PersonalProperty) TotalUserBalance(page, rows, status,tid int,Range, search string) (*ModelList, error) {
	fmt.Println("this is run")
	engine := utils.Engine_token
	query := engine.Alias("ut")
	query = query.Join("LEFT", "g_common.user u", "u.uid=ut.uid")
	query = query.Join("LEFT", "g_common.user_ex ex", "ex.uid=ut.uid")
	if Range!=``{
		query= query.Desc("ut.balance","ut.frozen")
	}else {
		query =query.Desc("ut.uid")
	}
	if tid !=0{
		query = query.Where("ut.token_id=?",tid)
	}
	fmt.Println("range===========",Range)
	value:=100000000
	Range =strings.Trim(Range," ")
	if strings.Compare(Range,"1~1000")==0{
		query =query.Where("(ut.frozen >? and ut.frozen<?)",1*value,1000*value)
	}else if strings.Compare(Range,"1001~10000")==0{
		query =query.Where("(ut.frozen >? and ut.frozen<?) ",1001*value,10000*value)
	}else if strings.Compare(Range,"10001~50000")==0{
		query =query.Where("(ut.frozen >? and ut.frozen<?) ",10001*value,50000*value)
	}else if strings.Compare(Range,"50万以上")==0{
		query =query.Where("ut.balance  >? or ut.frozen >?",500000*value,500000*value)
	}
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
	err = query.Select("ut.id,ut.token_id,ut.uid uid, ut.balance,ut.frozen, u.phone phone,u.email email,ex.nick_name nick_name,u.status status,u.account account ").GroupBy("ut.uid").Limit(mList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	for i,v:=range list{
		if tid !=0{
			list[i].BalanceTrue = convert.Int64ToStringBy8Bit(v.Balance)
			list[i].FrozenTrue =convert.Int64ToStringBy8Bit(v.Frozen)
		}else {
			list[i].BalanceTrue = "-"
			list[i].FrozenTrue ="-"
		}

	}
	mList.Items = list
	return mList, nil
}

/*
	统计平台所有币余额
*/
type TotalTokenCoin struct {
	//TotalBalance       int64  `xorm:"-" json:"total_balance"`
	//TotalFrozen        int64  `xorm:"-" json:"total_frozen"`
	TotalBalanceStr string `json:"total_balance_str"`
	TotalFrozenStr  string `json:"total_frozen_str"`
	TokenId         int32  `json:"token_id"`
	TokenName       string `json:"token_name"`
}
type TotalTokenCoinUser struct {
	TotalUser int64 `json:"total_user"`
	TokenId   int32 `json:"token_id"`
}

type CountToken struct {
	Total int64 `json:"total"`
}

func (this *UserToken) GetAllTokenCoin(tokenIdList []int32) (allbalanceList []TotalTokenCoin, allCoinUsers []TotalTokenCoinUser, err error) {
	engine := utils.Engine_token
	sql := "SELECT SUM(balance) AS total_balance_str, SUM(frozen) AS total_frozen_str, token_id, token_name FROM  g_token.`user_token` GROUP BY token_id"
	engine.In("token_id", tokenIdList).SQL(sql).Find(&allbalanceList)
	if err != nil {
		fmt.Println("token balance sum err:", err)
	}

	usersSql := "SELECT count(uid) as total_user, token_id  FROM g_token.`user_token` WHERE (balance > 0 OR frozen > 0 ) GROUP BY token_id"
	err = engine.In("token_id", tokenIdList).SQL(usersSql).Find(&allCoinUsers)
	if err != nil {
		fmt.Println(err)
	}

	return

}
