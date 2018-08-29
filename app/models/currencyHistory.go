package models

import (
	"fmt"
	"time"

	"admin/errors"
	"admin/utils"
)

type UserCurrencyHistory struct {
	UserInfo    `xorm:"extends"`
	BaseModel   `xorm:"-"`
	Id          int    `xorm:"not null pk autoincr comment('ID') INT(10)" json:"id"`
	Uid         int    `xorm:"not null default 0 INT(10)" json:"uid"`
	OrderId     string `xorm:"not null default '' comment('订单ID') VARCHAR(64)" json:"order_id"`
	TokenId     int    `xorm:"not null default 0 comment('货币类型') INT(10)" json:"token_id"`
	Num         int64  `xorm:"not null default 0 comment('数量') BIGINT(64)" json:"num"`
	Fee         int64  `xorm:"not null default 0 comment('手续费用') BIGINT(64)" json:"fee"`
	Surplus     int64  `xorm:"comment('账户余额') BIGINT(64)" json:"surplus"`
	Operator    int    `xorm:"not null default 0 comment('操作类型 操作类型 1订单转入 2订单转出 3币币划转到法币 4法币划转到币币 5冻结') TINYINT(2)" json:"operator"`
	Address     string `xorm:"not null default '' comment('提币地址') VARCHAR(255)" json:"address"`
	States      int    `xorm:"not null default 0 comment('订单状态: 0删除 1待支付 2待放行(已支付) 3确认支付(已完成) 4取消') TINYINT(2)" json:"states"`
	CreatedTime string `xorm:"not null comment('创建时间') DATETIME" json:"created_time"`
	UpdatedTime string `xorm:"comment('修改间') DATETIME" json:"updated_time"`
}

func (u *UserCurrencyHistory) TableName() string {
	return "user_currency_history"
}

func (u *UserCurrencyHistory) GetList(page, rows, ot int, date string) (*ModelList, error) {
	engine := utils.Engine_currency
	query := engine.Desc("id")
	//fmt.Println("000000000000000000000000000000000", ot)
	if ot != 0 {

		query = query.Where("operator=?", ot)
	}
	if date != `` {
		sub := date[:11] + "23:59:59"
		temp := fmt.Sprintf("created_time BETWEEN '%s' AND '%s'", date, sub)
		//fmt.Println(temp)
		query = query.Where(temp)
	}
	tempQuery := *query
	count, err := tempQuery.Count(&UserCurrencyHistory{})
	if err != nil {
		return nil, err
	}
	offset, modelList := u.Paging(page, rows, int(count))
	list := make([]UserCurrencyHistory, 0)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	modelList.Items = list
	return modelList, nil
}

//p2-3-3法币账户变更详情
func (u *UserCurrencyHistory) GetListForUid(page, rows, tid, status, chType int, bt, et string, search string) (*ModelList, error) {
	engine := utils.Engine_currency
	fmt.Println("------------------------>")
	query := engine.Alias("uch").Desc("u.uid")
	query = query.Join("LEFT", "g_common.user u ", "u.uid= uch.uid")
	query = query.Join("LEFT", "g_common.user_ex ex", "uch.uid=ex.uid")
	//substr := date[:11] + "23:59:59"
	//temp:= fmt.Sprintf("create_time BETWEEN '%s' AND '%s' ", st, substr)
	//query = query.Where(temp)
	//query =query.Where("uch.created_time between ? and ?", date,substr)

	query = query.Where("uch.token_id=?", tid)
	if chType != 0 {
		query = query.Where("uch.operator=?", chType)
	}
	if status != 0 {
		query = query.Where("u.status=?", status)
	}

	if bt != `` {
		if et != `` {
			subst := et[:11] + "23:59:59"
			query = query.Where("uch.created_time BETWEEN ? AND ? ", bt, subst)
		} else {
			subst := bt[:11] + "23:59:59"
			query = query.Where("uch.created_time BETWEEN ? AND ? ", bt, subst)
		}

	}
	if search != `` {
		temp := fmt.Sprintf(" concat(IFNULL(u.`uid`,''),IFNULL(u.`phone`,''),IFNULL(ex.`nick_name`,''),IFNULL(u.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}
	tempQuery := *query
	count, err := tempQuery.Count(&UserCurrencyHistory{})
	if err != nil {
		return nil, err
	}
	offset, modelList := u.Paging(page, rows, int(count))
	list := make([]UserCurrencyHistory, 0)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	for i, v := range list {
		list[i].NumTrue = u.Int64ToFloat64By8Bit(v.Num)
		list[i].SurplusTrue = u.Int64ToFloat64By8Bit(v.Surplus)
	}
	for i, v := range list {
		if v.Operator == 3 || v.Operator == 4 {
			list[i].UpdatedTime = v.CreatedTime
		}
	}
	modelList.Items = list
	return modelList, nil

}

// 法币交易合计
type CurrencyTradeTotal struct { // 交易次数
	TotalTime            int64 `xorm:"total_time"`               // 交易总次数
	TodayTotalTime       int64 `xorm:"today_total_time"`         // 今日交易次数
	YesterdayTotalTime   int64 `xorm:"yesterday_total_time"`     // 上日交易次数
	LastWeekDayTotalTime int64 `xorm:"last_week_day_total_time"` // 上周同日交易次数

	// 交易量
	TotalNum            string `xorm:"total_num"`               // 总计交易量
	TodayTotalNum       string `xorm:"today_total_num"`         // 今日交易量
	YesterdayTotalNum   string `xorm:"yesterday_total_num"`     // 上日交易量
	LastWeekDayTotalNum string `xorm:"last_week_day_total_num"` // 上周同日交易量

	// 交易手续费
	TotalFee            string `xorm:"total_fee"`               // 手续费总计
	TodayTotalFee       string `xorm:"today_total_fee"`         // 今日合计手续费
	YesterdayTotalFee   string `xorm:"yesterday_total_fee"`     // 上日合计手续费
	LastWeekDayTotalFee string `xorm:"last_week_day_total_fee"` // 上周同日合计手续费

}

// 交易次数、数量、手续费合计
// 今日、上日、上周同日
func (this *UserCurrencyHistory) TradeTotal() (*CurrencyTradeTotal, error) {
	// 计算日期
	todayDate := time.Now().Format(utils.LAYOUT_DATE)

	loc, err := time.LoadLocation("Local")
	if err != nil {
		return nil, errors.NewSys(err)
	}
	todayTime, err := time.ParseInLocation(utils.LAYOUT_DATE, todayDate, loc)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	yesterdayTime := todayTime.AddDate(0, 0, -1)
	lastWeekDayTime := todayTime.AddDate(0, 0, -7)

	todayDate = fmt.Sprintf("%s 00:00:00", todayDate)
	yesterdayDateBegin := fmt.Sprintf("%s 00:00:00", yesterdayTime.Format(utils.LAYOUT_DATE))
	yesterdayDateEnd := fmt.Sprintf("%s 23:59:59", yesterdayTime.Format(utils.LAYOUT_DATE))
	lastWeekDayDateBegin := fmt.Sprintf("%s 00:00:00", lastWeekDayTime.Format(utils.LAYOUT_DATE))
	lastWeekDayDateEnd := fmt.Sprintf("%s 23:59:59", lastWeekDayTime.Format(utils.LAYOUT_DATE))

	// 开始合计
	//1. 合计
	feeTotal := &CurrencyTradeTotal{}
	_, err = utils.Engine_currency.
		Table(this).
		Select("COUNT(id) total_time, IFNULL(SUM(num+fee), 0) total_num, IFNULL(SUM(fee), 0) total_fee").
		Where("operator IN (1,2)").
		Get(feeTotal)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	//2. 今日
	todayFeeTotal := &CurrencyTradeTotal{}
	_, err = utils.Engine_currency.
		Table(this).
		Select("COUNT(id) today_total_time, IFNULL(SUM(num+fee), 0) today_total_num, IFNULL(SUM(fee), 0) today_total_fee").
		Where("operator IN (1,2)").
		And("created_time>=?", todayDate).
		Get(todayFeeTotal)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	//3. 上日
	yesFeeTotal := &CurrencyTradeTotal{}
	_, err = utils.Engine_currency.
		Table(this).
		Select("COUNT(id) yesterday_total_time, IFNULL(SUM(num+fee), 0) yesterday_total_num, IFNULL(SUM(fee), 0) yesterday_total_fee").
		Where("operator IN (1,2)").
		And("created_time>=?", yesterdayDateBegin).
		And("created_time<=?", yesterdayDateEnd).
		Get(yesFeeTotal)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	//4. 上周同日
	lastWeekFeeTotal := &CurrencyTradeTotal{}
	_, err = utils.Engine_currency.
		Table(this).
		Select("COUNT(id) last_week_day_total_time, IFNULL(SUM(num+fee), 0) last_week_day_total_num, IFNULL(SUM(fee), 0) last_week_day_total_fee").
		Where("operator IN (1,2)").
		And("created_time>=?", lastWeekDayDateBegin).
		And("created_time<=?", lastWeekDayDateEnd).
		Get(lastWeekFeeTotal)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	// 合并
	feeTotal.TodayTotalTime = todayFeeTotal.TodayTotalTime
	feeTotal.TodayTotalNum = todayFeeTotal.TodayTotalNum
	feeTotal.TodayTotalFee = todayFeeTotal.TodayTotalFee

	feeTotal.YesterdayTotalTime = yesFeeTotal.YesterdayTotalTime
	feeTotal.YesterdayTotalNum = yesFeeTotal.YesterdayTotalNum
	feeTotal.YesterdayTotalFee = yesFeeTotal.YesterdayTotalFee

	feeTotal.LastWeekDayTotalTime = lastWeekFeeTotal.LastWeekDayTotalTime
	feeTotal.LastWeekDayTotalNum = lastWeekFeeTotal.LastWeekDayTotalNum
	feeTotal.LastWeekDayTotalFee = lastWeekFeeTotal.LastWeekDayTotalFee

	return feeTotal, nil
}
