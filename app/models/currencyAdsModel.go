package models

import (
	"admin/utils"
	"errors"
	"fmt"
)

// 买卖(广告)表
type Ads struct {
	BaseModel   `xorm:"-"`
	Id          uint64 `xorm:"not null pk autoincr INT(10)" json:"id"`
	Uid         int64  `xorm:"INT(10)" json:"uid"`              // 用户ID
	TypeId      uint32 `xorm:"TINYINT(1)" json:"type_id"`       // 类型:2出售 1购买
	TokenId     uint32 `xorm:"INT(10)" json:"token_id"`         // 货币类型
	TokenName   string `xorm:"VARBINARY(36)" json:"token_name"` // 货币名称
	Price       uint64 `xorm:"BIGINT(20)" json:"price"`         // 单价
	Num         uint64 `xorm:"BIGINT(20)" json:"num"`           // 数量
	Premium     int32  `xorm:"INT(10)" json:"premium"`          // 溢价
	AcceptPrice uint64 `xorm:"BIGINT(20)" json:"accept_price"`  // 可接受最低[高]单价
	MinLimit    uint32 `xorm:"INT(10)" json:"min_limit"`        // 最小限额
	MaxLimit    uint32 `xorm:"INT(10)" json:"max_limit"`        // 最大限额
	IsTwolevel  uint32 `xorm:"TINYINT(1)" json:"is_twolevel"`   // 是否要通过二级认证:0不通过 1通过
	Pays        string `xorm:"VARBINARY(50)" json:"pays"`       // 支付方式:以 , 分隔: 1,2,3
	Remarks     string `xorm:"VARBINARY(512)" json:"remarks"`   // 交易备注
	Reply       string `xorm:"VARBINARY(512)" json:"reply"`     // 自动回复问候语
	IsUsd       uint32 `xorm:"TINYINT(1)" json:"is_usd"`        // 是否美元支付:0否 1是
	States      uint32 `xorm:"TINYINT(1)" json:"states"`        // 状态:2下架 1上架
	CreatedTime string `xorm:"DATETIME" json:"created_time"`    // 创建时间
	UpdatedTime string `xorm:"DATETIME" json:"updated_time"`    // 修改时间
	IsDel       uint32 `xorm:"TINYINT(1)" json:"is_del"`        // 是否删除:0不删除 1删除
}

// 法币交易列表 - 用户虚拟币-订单统计 - 用户虚拟货币资产
type AdsUserCurrencyCount struct {
	Ads            `xorm:"extends"`
	SubductionZero `xorm:"-"`
	TWOVerifyMark  int    `xorm:"-"`
	Uname          string `xorm:"-"`
	Phone          string `xorm:"-"`
	Email          string `xorm:"-"`
	Ustatus        uint32 `xorm:"-"`
}

type AdsUserExUser struct {
	Ads             `xorm:"extends"`
	SubductionZero  `xorm:"-"`
	TWOVerifyMark   int `xorm:"-"`
	NickName        string
	Phone           string
	Email           string
	Status          int32
	SecurityAuth    int
	PremiumTrue     float64 `xorm:"-" json:"premium_true"`
	AcceptPriceTrue float64 `xorm:"-" json:"accept_price_true"`
	Range           string  `xorm:"-" json:"range"`
}

func (a *AdsUserExUser) TableName() string {
	return "ads"
}

//g_currency
func (AdsUserCurrencyCount) TableName() string {
	return "ads"
}

//法币挂单管理

//更具UID 获取所有广告id

type AdIdOfUid struct {
	Id     uint64 `xorm:"not null pk autoincr INT(10)" json:"id"`
	Uid    uint64 `xorm:"INT(10)" json:"uid"`        // 用户ID
	TypeId uint32 `xorm:"TINYINT(1)" json:"type_id"` // 类型:1出售 2购买
}

func (a *AdIdOfUid) TableName() string {
	return "ads"
}

func (a *Ads) TableName() string {
	return "ads"
}

//根据UID提取广告id
func (this *Ads) GetIdList(uid []int) ([]AdIdOfUid, error) {
	if len(uid) <= 0 {
		return nil, errors.New("uid list is empty !!")
	}
	engine := utils.Engine_currency
	adlist := make([]AdIdOfUid, 0)
	err := engine.In("uid", uid).Cols("id", "uid", "type_id").Find(&adlist)
	if err != nil {
		return nil, err
	}
	return adlist, nil
}

//法币挂单 下架订单
func (this *Ads) DownTradeAds(id, uid int) error {
	engine := utils.Engine_currency
	query := engine.Desc("id")
	query = query.Where("id=? AND uid=?", id, uid)
	has, err := query.Exist(&Ads{})
	if err != nil {
		//fmt.Println("0.0.0.0.0.0.0.0.00.0.0.")
		return err
	}
	if !has {
		return errors.New(" 订单不存在！！")
	}
	//sql:=fmt.Sprintf()
	_, err = engine.Exec("UPDATE g_currency.`ads` a SET a.`states`=0 WHERE a.`id`=? AND a.`uid` =?", id, uid)
	if err != nil {
		fmt.Println("what fuck you ")
		return err
	}
	return nil
}

//法币挂单管理
//交易 权限
//交易方向
// 刷选
func (this *Ads) GetAdsList(page, rows, status, tokenid, tradeid, verify int, search string) (*ModelList, error) {

	engine := utils.Engine_currency
	query := engine.Desc("ads.id")

	//总页数
	query = engine.Join("LEFT", "g_common.user u ", "u.uid= ads.uid")
	query = query.Join("LEFT", "g_common.user_ex ex", "ads.uid=ex.uid")
	// 分叉

	if verify != 0 {
		query = query.Where("ads.is_twolevel=?", verify)
	} //||
	if status != 0 {
		query = query.Where("u.status=?", status)
	}
	if search != `` {
		if len(search) > 0 {
			temp := fmt.Sprintf(" concat(IFNULL(u.`uid`,''),IFNULL(u.`phone`,''),IFNULL(ex.`nick_name`,''),IFNULL(u.`email`,'')) LIKE '%%%s%%'  ", search)
			query = query.Where(temp)
		}
	}
	if tradeid != 0 {
		query = query.Where("type_id=?", tradeid)
	}
	if tokenid != 0 {
		query = query.Where("token_id=?", tokenid)
	}
	//if date != `` {
	//	subst := date[:11] + "23:59:59"
	//	temp := fmt.Sprintf("created_time  BETWEEN '%s' AND '%s' ", date, subst)
	//	query = query.Where(temp)
	//}
	tempQuery := *query
	//uidQuery := *query
	//分两部分查询Count
	count, err := tempQuery.Count(&AdsUserExUser{})
	if err != nil {
		return nil, err
	}
	offset, mList := this.Paging(page, rows, int(count))
	list := make([]AdsUserExUser, 0)
	err = query.Limit(mList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}

	for i, v := range list {
		num, price := this.SubductionZeroMethod(v.Num, v.Price)
		list[i].NumberTrue = num
		list[i].PriceTrue = price
		if v.SecurityAuth&4 == 4 {
			list[i].TWOVerifyMark = 1
		} else {
			list[i].TWOVerifyMark = 0
		}
		list[i].PremiumTrue = this.Int64ToFloat64By8Bit(int64(v.Premium))
		list[i].AcceptPriceTrue = this.Int64ToFloat64By8Bit(int64(v.AcceptPrice))
		list[i].Range = fmt.Sprintf("%d-%d", v.MinLimit, v.MaxLimit)
	}
	mList.Items = list
	return mList, nil
}
