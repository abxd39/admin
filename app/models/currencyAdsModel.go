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
	Uid         uint64 `xorm:"INT(10)" json:"uid"`              // 用户ID
	TypeId      uint32 `xorm:"TINYINT(1)" json:"type_id"`       // 类型:1出售 2购买
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
	Ads     `xorm:"extends"`
	Uname   string `xorm:"-"`
	Phone   string `xorm:"-"`
	Email   string `xorm:"-"`
	Ustatus uint32 `xorm:"-"`
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

//法币挂单管理
func (this *Ads) GetAdsList(cur Currency) (*ModelList, error) {

	if cur.PageNum <= 0 {
		cur.PageNum = 100
	}
	engine := utils.Engine_currency
	query := engine.Desc("id")

	//总页数
	// query := engine.Join("INNER", "user_currency", "ads.uid=user_currency.uid")
	// query = query.Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid")
	//挂单日期
	if cur.TradeId != 0 {
		query = query.Where("type_id", cur.TradeId)
	}
	if cur.TokenId != 0 {
		query = query.Where("token_id=?", cur.TokenId)
	}

	if len(cur.Date) != 0 {
		substr := cur.Date[:11] + "23:59:59"
		query = query.Where("created_time BETWEEN ? AND ? ", cur.Date, substr)
	}
	tempQuery := *query
	uidQuery := *query
	//分两部分查询Count
	count, err := tempQuery.Count(&AdsUserCurrencyCount{})
	if err != nil {
		return nil, err
	}

	fmt.Println("count=", count)
	offset, modelList := this.Paging(cur.Page, cur.PageNum, int(count))
	fmt.Println("offset=", offset, "modelList", modelList)
	list := make([]AdsUserCurrencyCount, 0)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}

	//去重
	uidList := make([]Ads, 0)
	err = uidQuery.Distinct("uid").Find(&uidList)
	if err != nil {
		return nil, err
	}
	//查询用户资料

	//无条件判断

	uid := make([]uint64, 0)
	for _, v := range uidList {
		uid = append(uid, v.Uid)
	}
	fmt.Println("uidList", len(uid))
	//跨库查询 用户资料
	ulist, err := new(UserGroup).GetCurreryList(uid, cur.Verify, cur.Search)
	if err != nil {
		return nil, err
	}
	fmt.Println("ulist=", len(ulist))
	for index, value := range list {
		for _, v := range ulist {
			if value.Uid == v.Uid {
				list[index].Phone = v.Phone
				list[index].Uname = v.NickName
				list[index].Ustatus = uint32(v.Status)
				list[index].Email = v.Email
				break
			}
		}
	}
	modelList.Items = list
	return modelList, nil
}
