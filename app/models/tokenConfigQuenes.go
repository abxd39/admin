package models

import (
	"admin/utils"
	"errors"
	"fmt"
)

//交易对
type ConfigQuenes struct {
	BaseModel            `xorm:"-"`
	Id                   int64  `xorm:"pk autoincr BIGINT(20)" jons:"id"`
	TokenId              int    `xorm:"comment('交易币') unique(union_quene_id) INT(11)" jons:"token_id"`
	TokenTradeId         int    `xorm:"comment('实际交易币') unique(union_quene_id) INT(11)" jons:"token_trade_id"`
	BitCount             int    `xorm:"comment('小数位数') TINYINT(4)" json:"bit_count"`
	MinOrderUnm          int64  `xorm:"comment('最小交易量') BIGINT(20)" json:"min_order_unm"`
	Switch               int    `xorm:"comment('开关0关1开') TINYINT(4)" jons:"switch"`
	Price                int64  `xorm:"comment('初始价格') BIGINT(20)" jons:"price"`
	Scope                string `xorm:"comment('振幅') DECIMAL(6,2)" jons:"scope"`
	Name                 string `xorm:"comment('USDT/BTC') VARCHAR(32)" jons:"name"`
	Low                  int64  `xorm:"comment('最低价') BIGINT(20)" jons:"low"`
	High                 int64  `xorm:"comment('最高价') BIGINT(20)" jons:"high"`
	Amount               int64  `xorm:"comment('成交量') BIGINT(20)" jons:"amount"`
	SellPoundage         int64  `xorm:"comment('卖出手续费') BIGINT(20)" jons:"sell_poundage"`
	BuyPoundage          int64  `xorm:"comment('买入手续费') BIGINT(20)" jons:"buy_poundage"`
	BuyMinimumPrice      int64  `xorm:"comment('买入最小交易价') BIGINT(20)" jons:"buy_minimum_price"`
	BuyMaxmunPrice       int64  `xorm:"comment('买入最大交易价') BIGINT(20)" jons:"buy_maxmun_price"`
	SellMinimumPrice     int64  `xorm:"comment('卖出最小交易价') BIGINT(20)" jons:"sell_minimum_price"`
	SellMaxmumPrice      int64  `xorm:"comment('卖出最大交易价') BIGINT(20)" jons:"sell_maxmum_price"`
	MinimumTradingVolume int64  `xorm:"comment('最小交易额') BIGINT(20)" jons:"minimum_trading_volume"`
	MaxmumTradingVolume  int64  `xorm:"comment('最大交易额') BIGINT(20)" jons:"maxmum_trading_volume"`
	BeginTime            int    `xorm:"comment('开盘交易时间') INT(11)" jons:"begin_time"`
	EndTime              int    `xorm:"comment('闭盘交易时间') INT(11)" jons:"end_time"`
	SaturdaySwitch       int    `xorm:"comment('周六可交易') TINYINT(4)" jons:"saturday_switch"`
	SundaySwitch         int    `xorm:"comment('周日可交易') TINYINT(4)" jons:"sunday_switch"`
}

type ConfigQuenesEx struct {
	ConfigQuenes `xorm:"extends"`
	IsModifyMark bool `xorm:"-"`//是否允许修改状态
}

func (c *ConfigQuenesEx) TableName() string {
	return "config_quenes"
}

func (q *ConfigQuenes) GetTokenCashList(page, rows, id int) (*ModelList, error) {
	engine := utils.Engine_token

	query := engine.Desc("id")
	if id != 0 {
		query = query.Where("id=?", id)
	}
	tquery := *query
	count, err := tquery.Count(&ConfigQuenes{})
	if err != nil {
		return nil, err
	}
	offset, modelList := q.Paging(page, rows, int(count))
	query.Limit(modelList.PageSize, offset)

	list := make([]ConfigQuenes, 0)
	err = query.Find(&list)
	if err != nil {
		return nil, err
	}
	modelList.Items = list
	return modelList, nil
}

//删除 兑币
func (q *ConfigQuenes) DeleteCash(id int) error {
	engine := utils.Engine_token
	query := engine.Desc("id")
	query = query.Where("id=?", id)
	tempQuery := *query
	has, err := tempQuery.Exist(&ConfigQuenes{})
	if err != nil {
		return err
	}
	if !has {
		return errors.New(" 兑币对不存在！！")
	}
	_, err = query.Delete(&ConfigQuenes{})
	if err != nil {
		return err
	}
	return nil
}

//修改
func (q *ConfigQuenes) ModifyCash(id int) (*ConfigQuenesEx, error) {
	engine := utils.Engine_token
	query := engine.Desc("id")
	query = query.Where("id=?", id)
	tempQuery := *query
	c := new(ConfigQuenesEx)
	has, err := tempQuery.Get(c)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New(" 兑币对不存在！！")
	}
	//查数据库是否有挂单
	has, err = new(EntrustDetail).IsExist(c.Name)
	if err != nil {
		return nil, err
	}
	if has {
		c.IsModifyMark = has
	}
	return c, nil
}

//添加 和修改后的提交
func (q *ConfigQuenes) AddCash(c *ConfigQuenes) error {
	//
	engine := utils.Engine_token

	if c.Id != 0 {
		//修改 则需要判断是否存在
		has, err := engine.Id(c.Id).Exist(&ConfigQuenes{})
		if err != nil {

			return err
		}
		if !has {
			return errors.New("兑币不存在修改失败")
		}
		_, err = engine.Where("id=?", c.Id).Update(&ConfigQuenes{
			TokenId:              c.TokenId,
			TokenTradeId:         c.TokenTradeId,
			Switch:               c.Switch,
			Price:                c.Price,
			Scope:                c.Scope,
			Name:                 c.Name,
			Low:                  c.Low,
			High:                 c.High,
			Amount:               c.Amount,
			SellPoundage:         c.SellPoundage,
			BuyPoundage:          c.BuyPoundage,
			BuyMinimumPrice:      c.BuyMinimumPrice,
			BuyMaxmunPrice:       c.BuyMaxmunPrice,
			SellMinimumPrice:     c.SellMinimumPrice,
			SellMaxmumPrice:      c.SellMinimumPrice,
			MaxmumTradingVolume:  c.MaxmumTradingVolume,
			MinimumTradingVolume: c.MinimumTradingVolume,
			BeginTime:            c.BeginTime,
			EndTime:              c.EndTime,
			SaturdaySwitch:       c.SaturdaySwitch,
			SundaySwitch:         c.SundaySwitch,
		})
		if err != nil {
			fmt.Println("0000000000000001")
			return err
		}
		fmt.Println("11111111111111111")
		return nil
	} else {
		//新曾
		fmt.Println("2222222222222222")
		_, err := engine.InsertOne(c)
		if err != nil {
			return err
		}
		return nil

	}

}
