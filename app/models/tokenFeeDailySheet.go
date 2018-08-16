package models

import (
	"admin/utils"
)

//type TokenFeeDailySheet struct {
//	BaseModel    `xorm:"-"`
//	Id           int   `xorm:"not null pk autoincr comment('自增id') TINYINT(4)"`
//	Date         int64 `xorm:"not null comment('日期 精确到天 20180727 年保持4位月两位日两位') BIGINT(20)" json:"date"`
//	BuyBalance   int64 `xorm:"not null comment('买手续费日总额单位为人民币') BIGINT(20)" json:"buy_balance"`
//	SellBalance  int64 `xorm:"not null comment('卖手续费日总额单位为人民币') BIGINT(20)" json:"sell_balance"`
//	TotalBalance int64 `xorm:"not null comment('买卖总额单位为人民币') BIGINT(20)" json:"total_balance"`
//}

type TokenDailySheet struct {
	BaseModel    `xorm:"-"`
	Id           int   `xorm:"not null pk autoincr comment('自增id') TINYINT(4)"`
	TokenId      int   `xorm:"not null comment('货币id') INT(11)" json:"token_id"`
	FeeBuyCny    int64 `xorm:"not null comment('买手续费折合cny') BIGINT(20)" json:"fee_buy_cny"`
	FeeBuyTotal  int64 `xorm:"not null comment('买手续费总额') BIGINT(20)" json:"fee_buy_total"`
	FeeSellCny   int64 `xorm:"not null comment('卖手续费折合cny') BIGINT(20)" json:"fee_sell_cny"`
	FeeSellTotal int64 `xorm:"not null comment('卖手续费总额') BIGINT(20)" json:"fee_sell_total"`
	BuyTotal     int64 `xorm:"not null comment('买总额') BIGINT(20)" json:"buy_total"`
	BuyTotalCny  int64 `xorm:"not null comment('买总额折合') BIGINT(20)" json:"but_total_cny"`
	SellTotalCny int64 `xorm:"not null comment('卖总额折合') BIGINT(20)" json:"sell_total_cny"`
	SellTotal    int64 `xorm:"not null comment('卖总额') BIGINT(20)" json:"sell_total"`
	Date         int64 `xorm:"not null comment('时间戳，精确到天') BIGINT(10)" json:"date"`
}

type TokenFeeDailySheetGroup struct {
	TotalBuy  float64 `json:"total_buy"`
	TotalSell float64 `json:"total_sell"`
	Total     float64 `json:"total"`
}

type total struct {
	TokenDailySheet `xorm:"extends"`
	Total float64 `xorm:"-" json:"total" `
}

//获取历史交易记录
func (this *TokenDailySheet) GetDailySheetList(page, rows int, date uint64) (*ModelList, *TokenFeeDailySheetGroup, error) {
	engine := utils.Engine_token
	query := engine.Desc("id")
	if date != 0 {
		query = query.Where("date between ? and ?", date, date+86400)
	}
	countQuery := *query
	count, err := countQuery.Count(&TokenDailySheet{})
	if err != nil {
		return nil, nil, err
	}
	offset, mList := this.Paging(page, rows, int(count))
	list := make([]total, 0)
	err = query.Table("token_daily_sheet").Limit(mList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, nil, err
	}
	for i,v:=range list{
		list[i].Total = this.Int64ToFloat64By8Bit(v.BuyTotalCny +v.SellTotalCny)
	}
	mList.Items = list
	result, err := engine.SumsInt(this, "buy_total", "sell_total")
	if err != nil {
		return nil, nil, err
	}
	totalBuy := result[1]
	totalSell := result[0]
	return mList, &TokenFeeDailySheetGroup{
		Total:     this.Int64ToFloat64By8Bit(totalBuy+totalSell),
		TotalBuy:  this.Int64ToFloat64By8Bit(totalBuy),
		TotalSell: this.Int64ToFloat64By8Bit(totalSell),
	}, nil
}
