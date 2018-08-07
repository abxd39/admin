package models

import (
	"admin/utils"
	"fmt"
	"strconv"
)

// 用户虚拟货币资产表
type UserCurrency struct {
	Id        uint64 `xorm:"not null pk autoincr INT(10)" json:"id"`
	Uid       uint64 `xorm:"INT(10)"     json:"uid"`                                          // 用户ID
	TokenId   uint32 `xorm:"INT(10)"     json:"token_id"`                                     // 虚拟货币类型
	TokenName string `xorm:"VARCHAR(36)" json:"token_name"`                                   // 虚拟货币名字
	Freeze    int64  `xorm:"BIGINT not null default 0"   json:"freeze"`                       // 冻结
	Balance   int64  `xorm:"not null default 0 comment('余额') BIGINT"   json:"balance"`        // 余额
	Address   string `xorm:"not null default '' comment('充值地址') VARCHAR(255)" json:"address"` // 充值地址
	Version   int64  `xorm:"version"`
}

//折合 rmb
type AmountToCny struct {
	BaseModel    `xorm:"-"`
	UserCurrency `xorm:"extends"`
	AmountTo     float64 `xorm:"-"json:"amount_to"` //折合人民币
	Email        string  `json:"email"`
	Phone        string  `json:"phone"`
	Status       int     `json:"status"`
	NickName     string  `json:"nick_name"`
	Account      string  `json:"account"`
}

func (this *AmountToCny) TableName() string {
	return "user_currency"
}

//func (this *UserCurrency)TableName()string  {
//	return "user_currency"
//}
//获取单个用户的所有法币资产
func (this *UserCurrency) GetCurrencyList(page, rows, uid, tokenid int) ([]UserCurrency, int, int, error) {
	engine := utils.Engine_currency
	data := new(UserCurrency)
	var beginrow int
	if rows <= 0 {
		rows = 100
	}
	if page <= 1 {
		beginrow = 0
	} else {
		beginrow = (page - 1) * rows
	}

	list := make([]UserCurrency, 0)
	//根据uid 和 token_id 查询
	var total int
	count, err := engine.Where("uid=?", uid).Count(data)
	if err != nil {
		return nil, 0, 0, err
	}
	if int(count) < rows {
		total = 1
	} else {
		total = int(count) / rows
		v := int(count) % rows
		if v != 0 {
			total = total + 1
		}
	}
	if tokenid != 0 {
		// count, err := engine.Where("uid=? AND token_id=?", uid, tokenid).Count(data)
		// if err != nil {
		// 	return nil, 0, 0, err
		// }
		err = engine.Where("uid=? AND token_id=?", uid, tokenid).Limit(rows, beginrow).Find(&list)
		if err != nil {
			return nil, 0, 0, err
		}
		return list, 1, 1, nil
	}
	err = engine.Where("uid=?", uid).Limit(rows, beginrow).Find(&list)
	if err != nil {
		return nil, 0, 0, err
	}
	return list, page, total, nil

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

func (this *AmountToCny) CurrencyBalance(page, rows, status int, search string) (*ModelList, error) {
	engine := utils.Engine_currency
	fmt.Println("------->amountToCny")
	query := engine.Alias("uc").Desc("uc.uid")
	query = query.Cols("uc.uid", "uc.token_id", "uc.freeze", "uc.balance")
	query = query.Cols("u.phone", "u.email", "u.status", "u.account")
	query = query.Cols("ex.nick_name")
	query = query.Join("LEFT", "g_common.user u", "u.uid=uc.uid")
	query = query.Join("LEFT", "g_common.user_ex ex", "ex.uid=uc.uid")
	//query =query.GroupBy("uc.token_id")
	if status != 0 {
		query = query.Where("u.status=?", status)
	}
	if len(search) != 0 {
		temp := fmt.Sprintf(" concat(IFNULL(uc.`uid`,''),IFNULL(u.`phone`,''),IFNULL(ex.`nick_name`,''),IFNULL(u.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}

	countQuery := *query
	count, err := countQuery.Count(&AmountToCny{})
	if err != nil {
		return nil, err
	}
	offset, mList := this.Paging(page, rows, int(count))
	list := make([]AmountToCny, 0)
	err = query.GroupBy("uc.token_id").Limit(mList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	mList.Items = list
	//抓取所有货币的id 和价格 2018-08-07 4：15  折合 无法计算 原因是 没有取价格的的表

	return mList, nil
}

func (this *AmountToCny)Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}