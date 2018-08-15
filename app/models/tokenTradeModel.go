package models

import (
	"admin/utils"
	"fmt"
	"strconv"
	"time"
)

//bibi 交易表
//type Trade struct {
//	BaseModel    `xorm:"-"`
//	TradeId      int    `xorm:"not null pk autoincr comment('交易表的id') INT(11)"`
//	TradeNo      string `xorm:"comment('订单号') unique(uni_reade_no) VARCHAR(32)"`
//	Uid          int64  `xorm:"comment('买家uid') index BIGINT(11)"`
//	TokenId      int    `xorm:"comment('主货币id') index INT(11)"`
//	TokenTradeId int    `xorm:"comment('交易币种') INT(11)"`
//	TokenName    string `xorm:"not null comment('交易对 名称 例如USDT/BTC') VARCHAR(10)"`
//	Price        int64  `xorm:"comment('价格') BIGINT(20)"`
//	Num          int64  `xorm:"comment('数量') BIGINT(20)"`
//	Fee          int64  `xorm:"comment('手续费') BIGINT(20)"`
//	Opt          int    `xorm:"comment(' buy  1或sell 2') index unique(uni_reade_no) TINYINT(4)"`
//	DealTime     int64  `xorm:"comment('成交时间') BIGINT(11)"`
//	States       int    `xorm:"comment('0是挂单，1是部分成交,2成交， -1撤销') INT(11)"`
//	FeeCny       int64  `xorm:"comment( '手续费折合CNY') BIGINT(20)"`
//	TotalCny     int64  `xorm:"comment( '总交易额折合CNY') BIGINT(20)"`
//}

type Trade struct {
	BaseModel    `xorm:"-"`
	TradeId      int    `xorm:"not null pk autoincr comment('交易表的id') INT(11)"`
	TradeNo      string `xorm:"comment('订单号') unique(uni_reade_no) VARCHAR(32)"`
	Uid          int64  `xorm:"comment('买家uid') index BIGINT(11)"`
	TokenId      int    `xorm:"comment('主货币id') index INT(11)"`
	TokenTradeId int    `xorm:"comment('交易币种') INT(11)"`
	Symbol       string `xorm:"not null default 'BTC' comment('交易对 名称 例如USDT/BTC') VARCHAR(16)" `
	Price        int64  `xorm:"comment('价格') BIGINT(20)"`
	Num          int64  `xorm:"comment('数量') BIGINT(20)"`
	Fee          int64  `xorm:"comment('手续费') BIGINT(20)"`
	Opt          int    `xorm:"comment(' buy  1或sell 2') index unique(uni_reade_no) TINYINT(4)"`
	DealTime     int64  `xorm:"comment('成交时间') index BIGINT(11)"`
	//States       int    `xorm:"comment('0是挂单，1是部分成交,2成交， -1撤销') INT(11)"`
	FeeCny    int64  `xorm:"comment('手续费折合CNY') BIGINT(20)"`
	TotalCny  int64  `xorm:"comment('总交易额折合CNY') BIGINT(20)"`
	EntrustId string `xorm:"VARCHAR(32)"`
}

type TradeReturn struct {
	Trade          `xorm:"extends"`
	AllNum         int64   `xorm:"not null comment('总数量') BIGINT(20)"`  //总数
	SurplusNum     int64   `xorm:"not null comment('剩余数量') BIGINT(20)"` //余数
	FinishCount    float64 `xorm:"-" json:"finish_count"` //已成
	FeeTrue        float64 `xorm:"-" json:"fee_true"`
	AllNumTrue     float64 `xorm:"-" json:"all_num_true"`
	SurplusNumTrue float64 `xorm:"-" json:"surplus_num_true"`
	PriceTrue      float64 `xorm:"-" json:"price_true"`
}

func (t *TradeReturn) TableName() string {
	return "trade"
}

type TradeEx struct {
	Trade          `xorm:"extends"`
	ConfigTokenCny `xorm:"extends"`
	TotalTrue       float64 //交易总额
	FeeTrue         float64 //交易手续费
}

type TotalTradeCNY struct {
	Date  int64  //日期
	Buy   uint64 //买入总额
	Sell  uint64 //卖出总额
	Total uint64 // 买卖总金额
}

func (t *TradeEx) TableName() string {
	return "trade"
}

func (t *TotalTradeCNY) TableName() string {
	return "trade"
}

func (this *Trade) TotalTotalTradeList(page, rows int, date uint64) (*ModelList, error) {

	engine := utils.Engine_token
	query := engine.Desc("deal_time")
	query = query.Join("left", "config_token_cny p", "trade.token_id = p.token_id")
	query = query.GroupBy("deal_time")
	if date != 0 {
		temp := date / 1000
		query = query.Where("left(deal_time,7)=?", temp)
	}
	tempQuery := *query
	buyQuery := *query
	sellQuery := *query
	count, err := tempQuery.Count(&Trade{})
	if err != nil {
		return nil, err
	}
	offset, mList := this.Paging(page, rows, int(count))
	//买入总额
	buyList := make([]TradeEx, 0)
	err = buyQuery.Where("opt=1").Limit(mList.PageSize, offset).Find(&buyList)
	if err != nil {
		return nil, err
	}
	//卖出总额
	sellList := make([]TradeEx, 0)
	err = sellQuery.Where("opt=2").Limit(mList.PageSize, offset).Find(&sellList)
	//var totalBuy uint64
	//var totalSell uint64
	//买卖总金额
	totalDateList := make([]map[int64]*TotalTradeCNY, 0)
	dateMap := make(map[int64]*TotalTradeCNY, 0)
	for _, v := range buyList {
		key := v.DealTime / 1000
		for i, _ := range totalDateList {
			if _, ok := totalDateList[i][key]; !ok {
				dateMap[key] = &TotalTradeCNY{Date: v.DealTime}
				totalDateList = append(totalDateList, dateMap)
			}
			strBuy := this.Int64MulInt64By8BitString(v.Num, v.ConfigTokenCny.Price)
			buy, err := strconv.ParseUint(strBuy, 10, 64)
			if err != nil {
				continue
			}
			totalDateList[i][key].Buy += buy

		}

	}

	for _, v := range sellList {
		key := v.DealTime / 1000
		for i, _ := range totalDateList {
			if _, ok := totalDateList[i][key]; !ok {
				dateMap[key] = &TotalTradeCNY{Date: v.DealTime}
				totalDateList = append(totalDateList, dateMap)
			}
			strSell := this.Int64MulInt64By8BitString(v.Num, v.ConfigTokenCny.Price)
			sell, err := strconv.ParseUint(strSell, 10, 64)
			if err != nil {
				continue
			}
			totalDateList[i][key].Sell += sell
		}

	}

	return nil, nil
}

func (this *Trade) GetTokenRecordList(page, rows, opt, uid int, bt,et uint64, name string) (*ModelList, error) {
	engine := utils.Engine_token
	query := engine.Desc("t.entrust_id")
	query = query.Alias("t").Join("left", "entrust_detail e", "e.entrust_id= t.entrust_id")
	query = query.Where("(e.states=2 or e.states=1) and e.symbol=?", name) //交易对
	if opt != 0 {
		query = query.Where("t.opt=?", opt) //交易方向
	}
	if uid != 0 {
		query = query.Where("t.uid=?", uid)
	}
	if bt != 0 {
		if et!=0{
			query = query.Where("t.deal_time BETWEEN ? AND ? ", bt, et+86400)
		}else {
			query = query.Where("t.deal_time BETWEEN ? AND ? ", bt, bt+86400)
		}

	}
	tempQuery := *query

	count, err := tempQuery.Count(&TradeReturn{})
	if err != nil {
		return nil, err
	}
	offset, modelList := this.Paging(page, rows, int(count))
	list := make([]TradeReturn, 0)
	//fmt.Printf("$$$$$$$$$$$$$$$%#v\n", rows)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	for i, v := range list {
		allNum, surPlusNUm := this.SubductionZeroMethodInt64(v.AllNum, v.SurplusNum)
		list[i].AllNumTrue = allNum
		list[i].SurplusNumTrue = surPlusNUm
		list[i].FinishCount = allNum - surPlusNUm
		list[i].FeeTrue = this.Int64ToFloat64By8Bit(v.Fee)
		list[i].PriceTrue = this.Int64ToFloat64By8Bit(v.Price)
	}
	modelList.Items = list
	return modelList, nil
}

//p5-1-0-1币币交易手续费明细
/********************************
* id 兑币id
* trade_type 交易方向 1 卖 2买
* search 筛选
 */
func (this *Trade) GetFeeInfoList(page, rows, uid, opt int, date uint64, name string) (*ModelList, error) {
	engine := utils.Engine_token
	query := engine.Desc("trade.token_id")
	query = query.Join("left", "config_token_cny p", "trade.token_id = p.token_id")

	if uid != 0 {
		query = query.Where("uid=?", uid)
	}
	if date != 0 {
		query = query.Where("deal_time BETWEEN ? AND ?", date, date+86400)
	}
	if opt != 0 {
		query = query.Where("opt=?", opt)
	}
	if name != `` {
		query = query.Where("token_name=?", name)
	}
	ValuQuery := *query
	count, err := query.Distinct("deal_time").Count(&Trade{})
	if err != nil {
		return nil, err
	}
	offset, mlist := this.Paging(page, rows, int(count))
	list := make([]TradeEx, 0)
	err = ValuQuery.Limit(mlist.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	//未完待续 折合成人民币
	//fmt.Println("len=",len(list))
	for i, v := range list {
		list[i].TotalTrue = this.Int64ToFloat64By8Bit(v.Num)
		list[i].FeeTrue = this.Int64ToFloat64By8Bit(v.Fee)
	}
	mlist.Items = list
	return mlist, nil
}

//获取单日bibi交易手续费
func (this *Trade) GetTodayFee() (float64, error) {
	engine := utils.Engine_token
	sql := "SELECT fee FROM (SELECT FROM_UNIXTIME(deal_time,'%Y-%m-%d')days,SUM(fee_cny) fee FROM g_token.trade where states=2  ) t WHERE "
	current := time.Now().Format("2006-01-02 15:04:05")
	current = fmt.Sprintf(" t.days='%s'", current[:10])
	fee := &struct {
		Fee float64
	}{}

	_, err := engine.SQL(sql + current).Get(fee)
	if err != nil {
		return 0, err
	}
	return fee.Fee, nil
}
