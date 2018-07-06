package models

import (
	"admin/utils"
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
	PayId       uint64 `xorm:"not null default 0 comment('支付类型') INT(10)"       json:"pay_id"`
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

//
//func (this *Order) GetList

//列出订单
func (this *Order) GetOrderList(Page, PageNum, AdType, States, TokenId int, StartTime, EndTime string) (*[]Order, int, error) {

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
		return nil, 0, err
	}
	return &list, int(total) / PageNum, nil

}
