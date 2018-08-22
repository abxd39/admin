package models

import (
	"admin/utils"
	"fmt"
	"strconv"
	"admin/utils/convert"
)

// 用户虚拟货币资产表
type UserCurrency struct {
	BaseModel  `xorm:"-"`
	Id         uint64 `xorm:"not null pk autoincr INT(10)" json:"id"`
	Uid        uint64 `xorm:"INT(10)"     json:"uid"`                                          // 用户ID
	TokenId    uint32 `xorm:"INT(10)"     json:"token_id"`                                     // 虚拟货币类型
	TokenName  string `xorm:"VARCHAR(36)" json:"token_name"`                                   // 虚拟货币名字
	Freeze     int64  `xorm:"BIGINT not null default 0"   json:"freeze"`                       // 冻结
	Balance    int64  `xorm:"not null default 0 comment('余额') BIGINT"   json:"balance"`        // 余额
	Address    string `xorm:"not null default '' comment('充值地址') VARCHAR(255)" json:"address"` // 充值地址
	Version    int64  `xorm:"version"`
	BalanceCny int64  `xorm:"default 0 comment('折合的余额人民币') BIGINT(20)" json:"balance_cny"`
	FreezeCny  int64  `xorm:"default 0 comment('冻结折合人民币') BIGINT(20)" json:"freeze_cny"`
}

//折合 rmb
type AmountToCny struct {
	UserCurrency `xorm:"extends"`
	AmountTo   string `xorm:"-"json:"amount_to"` //折合人民币
	Email      string  `json:"email"`
	Phone      string  `json:"phone"`
	Status     int     `json:"status"`
	NickName   string  `json:"nick_name"`
	Account    string  `json:"account"`
}

func (this *AmountToCny) TableName() string {
	return "user_currency"
}

type DetailCurrency struct {
	UserCurrency `xorm:"extends"`
	BalanceTrue float64 `xorm:"-" json:"balance_true"`
	FreezeTrue  float64 `xorm:"-" json:"freeze_true" `
	AmountTo    float64 `xorm:"-"json:"amount_to"` //折合人民币
}

func (this *DetailCurrency) TableName() string {
	return "user_currency"
}

//获取单个用户的所有法币资产
func (this *UserCurrency) GetCurrencyList(page, rows, uid, tokenId int) (*ModelList, error) {
	engine := utils.Engine_currency

	query := engine.Where("uid=?", uid)
	if tokenId != 0 {
		query = query.Where(" token_id=?", tokenId)

	}
	//根据uid 和 token_id 查询
	queryCount := *query
	count, err := queryCount.Count(&DetailCurrency{})
	if err != nil {
		return nil, err
	}
	offset, mList := this.Paging(page, rows, int(count))
	list := make([]DetailCurrency, 0)

	err = query.Limit(mList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	for i, v := range list {
		list[i].BalanceTrue = this.Int64ToFloat64By8Bit(v.Balance)
		list[i].FreezeTrue = this.Int64ToFloat64By8Bit(v.Freeze)
		list[i].AmountTo = this.Int64ToFloat64By8Bit(v.BalanceCny) + this.Int64ToFloat64By8Bit(v.FreezeCny)
		fmt.Println(v.BalanceCny)
		fmt.Println(v.FreezeCny)
	}
	mList.Items = list
	return mList, nil
}

func (this *UserCurrency) GetAll(uid []int64) ([]UserCurrency, error) {
	engine := utils.Engine_currency
	list := make([]UserCurrency, 0)
	err := engine.In("uid", uid).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (this *UserCurrency) GetBalance(uid, token_id int) (*UserCurrency, error) {
	engine := utils.Engine_currency
	data := new(UserCurrency)
	_, err := engine.Where("uid=? AND token_id=?", uid, token_id).Get(data)
	if err != nil {
		return nil, err
	}
	return data, nil

}

//p2-3-1法币账户统计列表

func (this *UserCurrency) CurrencyBalance(page, rows, status int, search string) (*ModelList, error) {
	engine := utils.Engine_currency
	query := engine.Alias("uc").Desc("uc.uid")
	//query = query.Cols()
	//query = query.Cols()
	//query = query.Cols("ex.nick_name")
	query = query.Join("LEFT", "g_common.user u", "u.uid=uc.uid")
	query = query.Join("LEFT", "g_common.user_ex ex", "ex.uid=uc.uid")
	//query =query.GroupBy("uc.uid")
	if status != 0 {
		query = query.Where("u.status=?", status)
	}
	if len(search) != 0 {
		temp := fmt.Sprintf(" concat(IFNULL(uc.`uid`,''),IFNULL(u.`phone`,''),IFNULL(ex.`nick_name`,''),IFNULL(u.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}

	countQuery := *query
	count, err := countQuery.Distinct("uc.uid").Count(&AmountToCny{})
	if err != nil {
		return nil, err
	}
	offset, mList := this.Paging(page, rows, int(count))
	list := make([]AmountToCny, 0)
	err = query.Select("sum(uc.freeze_cny) AS freeze_cny ,sum(uc.balance_cny) balance_cny,uc.uid ,uc.token_id,uc.freeze, u.phone,u.email,u.status,u.account").GroupBy("uc.uid").Limit(mList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}

	for i, v := range list {
		fmt.Println(v.BalanceCny)
		fmt.Println(v.FreezeCny)
		list[i].AmountTo = fmt.Sprintf("%.2f", convert.Int64ToFloat64By8Bit(v.BalanceCny +v.FreezeCny))
	}
	mList.Items = list

	return mList, nil
}

func (this *UserCurrency) Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
