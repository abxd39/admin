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
	Ads           `xorm:"extends"`
	SubductionZero `xorm:"-"`
	TWOVerifyMark int    `xorm:"-"`
	Uname         string `xorm:"-"`
	Phone         string `xorm:"-"`
	Email         string `xorm:"-"`
	Ustatus       uint32 `xorm:"-"`
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
	_, err = engine.Exec("UPDATE g_currency.`ads` a SET a.`states`=0 WHERE a.`id`=? AND a.`uid` =?",id,uid)
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
func (this *Ads) GetAdsList(page, rows, status, tokenid, tradeid, verify int, search, date string) (*ModelList, error) {

	engine := utils.Engine_currency
	query := engine.Desc("id")

	//总页数
	// query := engine.Join("INNER", "user_currency", "ads.uid=user_currency.uid")
	// query = query.Join("LEFT", "user_currency_count", "ads.uid=user_currency_count.uid")
	// 分叉
	//挂单日期
	if verify != 0 || status != 0 || search != `` {
		ulist, err := new(UserGroup).UserList(page, rows, verify, status, search, 0)
		if err!=nil{
			return nil,err
		}
		fmt.Println("认证刷选",ulist)
		uid := make([]int64, 0)
		value, ok := ulist.Items.([]UserGroup)
		if !ok {
			return nil, errors.New("assert failed!!!")
		}
		for _, v := range value {
			uid = append(uid, v.Uid)
		}
		fmt.Println("uidlist", uid)
		list := make([]AdsUserCurrencyCount, 0)
		query = query.In("uid", uid)
		tempQuery := *query
		count, err := tempQuery.Count(&AdsUserCurrencyCount{})
		offset, modelList := this.Paging(page, rows, int(count))
		err = query.Limit(modelList.PageSize, offset).Find(&list)
		if err != nil {
			return nil, err
		}
		for index, vads := range list {
			for _, v := range value {
				if vads.Uid == v.Uid {
					list[index].TWOVerifyMark = v.TWOVerifyMark
					list[index].Phone = v.Phone
					list[index].Uname = v.NickName
					list[index].Ustatus = uint32(v.Status)
					list[index].Email = v.Email
					break
				}
			}
		}
		//去掉零
		for i,v:=range list{
			num,price :=this.SubductionZeroMethod(v.Num,v.Price)
			list[i].NumberTrue =num
			list[i].PriceTrue = price
		}
		modelList.Items = list
		return modelList, nil
	}

	if tradeid != 0 {
		query = query.Where("type_id=?", tradeid)
	}
	if tokenid != 0 {
		query = query.Where("token_id=?", tokenid)
	}
	if date != `` {
		subst := date[:11] + "23:59:59"
		temp := fmt.Sprintf("created_time  BETWEEN '%s' AND '%s' ", date, subst)
		query = query.Where(temp)
	}
	tempQuery := *query
	uidQuery := *query
	//分两部分查询Count
	count, err := tempQuery.Count(&AdsUserCurrencyCount{})
	if err != nil {
		return nil, err
	}

	fmt.Println("count=", count)
	offset, modelList := this.Paging(page, rows, int(count))
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

	uid := make([]int64, 0)
	for _, v := range uidList {
		uid = append(uid, v.Uid)
	}
	fmt.Println("uidList", len(uid))
	//跨库查询 用户资料
	ulist, err := new(UserGroup).GetCurrencyList(uid, verify, search)
	if err != nil {
		return nil, err
	}
	fmt.Println("ulist=", len(ulist))
	for index, value := range list {
		for _, v := range ulist {
			if value.Uid == v.Uid {
				list[index].TWOVerifyMark = v.TWOVerifyMark
				list[index].Phone = v.Phone
				list[index].Uname = v.NickName
				list[index].Ustatus = uint32(v.Status)
				list[index].Email = v.Email
				break
			}
		}
	}
	//去掉零
	for i,v:=range list{
		num,price :=this.SubductionZeroMethod(v.Num,v.Price)
		list[i].NumberTrue =num
		list[i].PriceTrue = price
	}
	modelList.Items = list
	return modelList, nil
}
