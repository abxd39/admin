package models

import (
	"admin/utils"
	"errors"
	"fmt"
)

// 订单表
type Order struct {
	Id          uint64 `xorm:"not null pk autoincr comment('ID')  INT(10)"  json:"id"`
	OrderId     string `xorm:"not null pk comment('订单ID') INT(10)"   json:"order_id"` // hash( type_id, 6( user_id, + 时间秒）
	AdId        uint64 `xorm:"not null default 0 comment('广告ID') index INT(10)"  json:"ad_id"`
	AdType      uint32 `xorm:"not null default 0 comment('广告类型:1出售 2购买') TINYINT(1)"  json:"ad_type"`
	Price       int64  `xorm:"not null default 0 comment('价格') BIGINT(64)"   json:"price"`
	Num         int64  `xorm:"not null default 0 comment('数量') BIGINT(64)"   json:"num"`
	TokenId     uint64 `xorm:"not null default 0 comment('货币类型') INT(10)"       json:"token_id"`
	PayId       string `xorm:"not null default 0 comment('支付类型') VARCHAR(64)"       json:"pay_id"`
	SellId      uint64 `xorm:"not null default 0 comment('卖家id') INT(10)"         json:"sell_id"`
	SellName    string `xorm:"not null default '' comment('卖家昵称') VARCHAR(64)"  json:"sell_name"`
	BuyId       uint64 `xorm:"not null default 0 comment('买家id') INT(10)"    json:"buy_id"`
	BuyName     string `xorm:"not null default '' comment('买家昵称') VARCHAR(64)"   json:"buy_name"`
	Fee         int64  `xorm:"not null default 0 comment('手续费用') BIGINT(64)"  json:"fee"`
	States      uint32 `xorm:"not null default 0 comment('订单状态: 0删除 1待支付 2待放行(已支付) 3确认支付(已完成) 4取消') TINYINT(1)"   json:"states"`
	PayStatus   uint32 `xorm:"not null default 0 comment('支付状态: 1待支付 2待放行(已支付) 3确认支付(已完成)') TINYINT(1)"  json:"pay_status"`
	CancelType  uint32 `xorm:"not null default 0 comment('取消类型: 1卖方 2 买方') TINYINT(1)"   json:"cancel_type"`
	CreatedTime string `xorm:"not null comment('创建时间') DATETIME"  json:"created_time"`
	UpdatedTime string `xorm:"comment('修改时间') DATETIME"    json:"updated_time"`
}

type OrderGroup struct {
	Order `xorm:"extends"`
	Uid   uint64 `xorm:"INT(10)"     json:"uid"`
}

func (o *Order) TableName() string {
	return "order"
}

//查询个人的所有数据货币的交易记录
func (this *Order) GetOrderListOfUid(page, rows, uid, token_id int) ([]OrderGroup, int, int, error) {
	if page <= 1 {
		page = 1
	}
	if rows < 1 {
		rows = 50
	}
	begin := 0
	if page > 1 {
		begin = (page - 1) * rows
	}
	list := make([]OrderGroup, 0)
	engine := utils.Engine_currency

	query := engine.Desc("order.id")
	query = query.Join("INNER", "ads", "order.ad_id=ads.id")
	if token_id != 0 {
		query.Where("token_id=? and pay_status=3", token_id)
	}
	query = query.Where("uid=?", uid)
	query = query.Limit(rows, begin)
	//query.GroupBy
	tempquery := query
	err := query.Find(&list)
	if err != nil {
		return nil, 0, 0, err
	}
	count, err := tempquery.Distinct("token_id").Count(&Order{})
	//engine.Query("")
	fmt.Printf("%#v\n", count)
	if err != nil {
		return nil, 0, 0, err
	}

	total := int(count) / rows
	fmt.Println("000000000000000000")
	return list, total, int(count), nil
}

//
//根据 uid  及交易状态 多表查询拉取 所有相关订单的交易记录
func (this *Order) GetOrderId(uid []int, status int) ([]OrderGroup, error) {
	if len(uid) <= 0 {
		return nil, errors.New("uid [] is empty!!")
	}
	fmt.Println("GetOrderId", uid, status)
	list := make([]OrderGroup, 0)
	engine := utils.Engine_currency
	query := engine.Desc("order.id")
	query = query.Join("INNER", "ads", "order.ad_id=ads.id")
	query = query.In("uid", uid)
	query = query.Where("pay_status=?", status)
	err := query.Find(&list)

	//err := engine.In("uid", orderId).Where("status=?", status).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

//列出订单
func (this *Order) GetOrderList(Page, PageNum, AdType, States, TokenId int, StartTime, EndTime string) (*[]Order, int, int, error) {

	engine := utils.Engine_currency
	if Page <= 1 {
		Page = 1
	}
	if PageNum <= 0 {
		PageNum = 100
	}

	query := engine.Desc("id")
	orderModel := new(Order)

	// if States != 0 { // 状态为0，表示已经删除
	// 	query = query.Where("states = 0")
	// } else {
	// 	query = query.Where("states = ?", States)
	// }

	// if Id != 0 {
	// 	query = query.Where("id = ?", Id)
	// }
	// if AdType != 0 {
	// 	query = query.Where("ad_type = ?", AdType)
	// }
	// if TokenId != 0 {
	// 	query = query.Where("token_id = ?", TokenId)
	// }
	//fmt.Println(StartTime, EndTime)
	if StartTime != `` {
		query = query.Where("created_time >= ?", StartTime)
	}
	if EndTime != `` {
		query = query.Where("created_time <= ?", EndTime)
	}
	list := make([]Order, 0)
	tmpQuery := *query
	countQuery := &tmpQuery
	//查询符合要求数据
	err := query.Limit(PageNum, (Page-1)*PageNum).Find(&list)
	//所有符合要求的数据的函数
	total, _ := countQuery.Count(orderModel)

	if err != nil {
		return nil, 0, 0, err
	}
	page := int(total) / PageNum
	return &list, page, int(total), nil

}
