package models

import (
	"fmt"
	"time"

	"admin/errors"
	"admin/utils"
)

//数据库 g_wallet 日提币汇总表
//type WalletInoutDailySheet struct {
//	BaseModel       `xorm:"-"`
//	Id              int    `xorm:"not null pk autoincr comment('自增id') TINYINT(4)"`
//	TokenId         int    `xorm:"not null comment('货币id') TINYINT(4)"`
//	TokenName       string `xorm:"not null comment('货币名称') VARCHAR(20)"`
//	TotalCny int64  `xorm:"not null comment('日提币总金额(cny)') BIGINT(20)"`
//	TotalDayFee     int64  `xorm:"not null comment('日提币手续费总金额(cny)') BIGINT(20)"`
//	TotalBalance    int64  `xorm:"not null comment('累计提币总金额(cny)') BIGINT(20)"`
//	TotalFee        int64  `xorm:"not null comment('累计提币手续费总金额(cny)') BIGINT(20)"`
//	Date            string `xorm:"not null default 'CURRENT_TIMESTAMP' comment('日期天') TIMESTAMP"`
//}

type TokenInoutDailySheet struct {
	BaseModel      `xorm:"-"`
	Id             int    `xorm:"not null pk autoincr comment('自增id') TINYINT(4)" `
	TokenId        int    `xorm:"not null comment('货币id') TINYINT(4)" json:"token_id"`
	TokenName      string `xorm:"not null comment('货币名称') VARCHAR(20)" json:"token_name"`
	TotalDayNum    int64  `xorm:"not null comment('日提币总量') BIGINT(20)" json:"total_day_num"`
	TotalDayCny    int64  `xorm:"not null comment('日提币总数折合') BIGINT(20)" json:"total_day_cny"`
	Total          int64  `xorm:"not null comment('提币累计总金额') BIGINT(20)" json:"total"`
	TotalDayNumFee int64  `xorm:"not null comment('日提币手续费数量') BIGINT(20)" json:"total_day_num_fee"`
	TotalFee       int64  `xorm:"not null comment('提币手续费累计总金额') BIGINT(20)" json:"total_fee"`
	TotalDayFeeCny int64  `xorm:"not null comment('日提币手续费总数折合') BIGINT(20)" json:"total_day_fee_cny"`
	TotalPut       int64  `xorm:"not null comment('充币累计总额') BIGINT(20)" json:"total_put"`
	TotalDayPut    int64  `xorm:"not null comment('日充币总额') BIGINT(20)" json:"total_day_put"`
	TotalDayPutCny int64  `xorm:"not null default 0 comment('日充币折合') BIGINT(20)" json:"total_day_put_cny"`
	Date           string `xorm:"not null comment('时间戳') datetime" json:"date"`
}

type FeeTotalSheet struct {
	TokenInoutDailySheet `xorm:"extends"`
	TotalDayNumTrue      float64 `xorm:"-" json:"total_num_true"`
	TotalTrue            float64 `xorm:"-" json:"total_true"`
	TotalFeeTrue         float64 `xorm:"-" json:"total_fee_true"`
	TotalDayNumFeeTrue   float64 `xorm:"-" json:"total_day_num_fee_true"`
}

// 走势返回string，内容是int64
// 如果用int64，数据太大时xorm sum会溢出报错
type InOutTrend struct {
	InTotal  string `xorm:"in_total"`
	OutTotal string `xorm:"out_total"`
	FeeTotal string `xorm:"fee_total"`
	Date     string `xorm:"date"`
}

func (this *FeeTotalSheet) TableName() string {
	return "token_inout_daily_sheet"
}
func (this *TokenInoutDailySheet) TableName() string {
	return "token_inout_daily_sheet"
}

//定时结算bibi 日提币报表表数据
//func (this *TokenInoutDailySheet) BoottimeTimingSettlement() {
//
//	for {
//		now := time.Now()
//		// 计算下一个零点
//		next := now.Add(time.Hour * 24)
//		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
//		t := time.NewTimer(next.Sub(now))
//		<-t.C
//		//Printf("定时结算Boottime表数据，结算完成: %v\n",time.Now())
//		//以下为定时执行的操作
//		current := now.Format("2006-01-02 15:04:05")
//		current = current[:10]
//		//cunrrentUnixtime := now.Unix()
//		//bibi 日报表统计
//		engine := utils.Engine_wallet
//		list := make([]TokenInoutDailySheet, 0)
//		sql := fmt.Sprintf("SELECT SUM(fee_cny) total_day_fee,SUM(amount_cny) total_day_balance,t.days,tokenid token_id ,token_name FROM (SELECT SUBSTRING(created_time,1,10) days,tokenid,token_name, states,opt, fee_cny ,amount_cny FROM g_wallet.token_inout) t  WHERE t.states=2 AND t.opt=1 AND t.days='%s' GROUP BY t.tokenid", current)
//		fmt.Println(sql)
//		err := engine.SQL(sql).Find(&list)
//		if err != nil {
//			utils.AdminLog.Println("日提币统计失败！", err.Error(), current)
//			//continue
//		}
//		fmt.Println("1->		", list)
//		//根据id 抓取最后插入数据库的的数据
//		for _, v := range list {
//			if v.TokenId == 0 {
//				fmt.Println("error......")
//				continue
//			}
//			lastSql := fmt.Sprintf("SELECT total_balance,total_fee FROM g_wallet.token_inout_daily_sheet WHERE id= (SELECT MAX(id) FROM g_wallet.token_inout_daily_sheet where token_id=%d )", v.TokenId)
//			fmt.Println(lastSql)
//			_, err := engine.SQL(lastSql).Get(this)
//			if err != nil {
//				utils.AdminLog.Println(err.Error())
//				continue
//			}
//			fmt.Println("2->		", this)
//			this.Id = 0
//			this.TokenId = v.TokenId
//			this.TokenName = v.TokenName
//			this.TotalFee = v.TotalDayFee
//			this.TotalDayBalance = v.TotalDayBalance
//			this.TotalBalance += this.TotalDayBalance
//			this.TotalFee += this.TotalDayFee
//			this.Date = now.Unix()
//			fmt.Println("3->		", this)
//			_, err = engine.Table("token_inout_daily_sheet").AllCols().InsertOne(this)
//			if err != nil {
//				utils.AdminLog.Println(err.Error())
//				continue
//			}
//			fmt.Println("successful")
//		}
//	}
//}

//提币手续费汇总表
func (this *TokenInoutDailySheet) GetInOutDailySheetList(page, rows, tokenId int, bt, et string) (*ModelList, error) {

	engine := utils.Engine_wallet
	query := engine.Desc("id")
	query = query.Where("total !=0 or total_day_num !=0 or total_fee !=0 or total_day_num_fee !=0 ")
	if tokenId != 0 {
		query = query.Where("token_id=?", tokenId)
	}
	//substr := st[:11] + "23:59:59"
	if bt != `` {
		if et != `` {
			query = query.Where("date between ? and ?", bt, et[:11]+"23:59:59")
		} else {
			query = query.Where("date between ? and ?", bt, bt[:11]+"23:59:59")
		}
	}
	//query = query.Where("id>?", 0)
	countQuery := *query
	count, err := countQuery.Count(&TokenInoutDailySheet{})
	if err != nil {
		return nil, err
	}
	fmt.Println("已经到这里了")
	offset, mList := this.Paging(page, rows, int(count))

	list := make([]FeeTotalSheet, 0)
	err = query.Limit(mList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	for i, v := range list {
		list[i].TotalTrue = this.Int64ToFloat64By8Bit(v.Total)
		list[i].TotalDayNumTrue = this.Int64ToFloat64By8Bit(v.TotalDayNum)
		list[i].TotalFeeTrue = this.Int64ToFloat64By8Bit(v.TotalFee)
		list[i].TotalDayNumFeeTrue = this.Int64ToFloat64By8Bit(v.TotalDayNumFee)
	}
	mList.Items = list
	return mList, nil
}

//日冲币汇总表
func (t *TokenInoutDailySheet) DayPutDailySheet(page, rows, tid int, bt, et string) (*ModelList, error) {
	engine := utils.Engine_wallet
	query := engine.Desc("id")
	query = query.Where("total_day_put !=0 or total_put !=0")
	if tid != 0 {
		query = query.Where("token_id=?", tid)
	}

	if bt != `` {
		if et != `` {
			query = query.Where("date between ? and ?", bt, et[:11]+"23:59:59")
		} else {
			query = query.Where("date between ? and ?", bt, bt[:11]+"23:59:59")
		}
	}

	countCount := *query
	count, err := countCount.Count(t)
	if err != nil {
		return nil, err
	}
	offset, mList := t.Paging(page, rows, int(count))

	type temp struct {
		TotalPut        int64   `json:"total_put"`
		TotalDayPut     int64   `json:"total_day_put"`
		TokenName       string  `json:"token_name"`
		TokenId         int     `json:"token_id"`
		TotalPutTrue    float64 `xorm:"-" json:"total_true"`
		TotalDayPutTrue float64 `xorm:"-" json:"total_day_true"`
		Date            string  ` json:"date"`
	}
	list := make([]temp, 0)

	err = query.Table("token_inout_daily_sheet").Limit(mList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	for i, v := range list {
		list[i].TotalPutTrue = t.Int64ToFloat64By8Bit(v.TotalPut)
		list[i].TotalDayPutTrue = t.Int64ToFloat64By8Bit(v.TotalDayPut)
	}
	//result:=make([]temp,0)
	//for _,v:=range list {
	//	if v.TotalDayPutTrue !=0 || v.TotalPutTrue !=0{
	//		result = append(result,v)
	//	}
	//}
	//_,mList = t.Paging(page,rows,len(result))
	mList.Items = list
	return mList, nil
}

//日提币
func (t *TokenInoutDailySheet) DayOutDailySheet(page, rows, tid int, bt, et string) (*ModelList, error) {
	engine := utils.Engine_wallet
	query := engine.Desc("id")
	query = query.Where("total !=0 or total_day_num !=0")
	if tid != 0 {
		query = query.Where("token_id=?", tid)
	}

	if bt != `` {
		if et != `` {
			query = query.Where("date between ? and ?", bt, et[:11]+"23:59:59")
		} else {
			query = query.Where("date between ? and ?", bt, bt[:11]+"23:59:59")
		}
	}

	countCount := *query
	count, err := countCount.Count(t)
	if err != nil {
		return nil, err
	}

	offset, mList := t.Paging(page, rows, int(count))
	type temp struct {
		Total        int64   `json:"total"`
		TotalDayNum  int64   `json:"total_num"`
		TokenName    string  `json:"token_name"`
		TokenId      int     `json:"token_id"`
		TotalTrue    float64 `xorm:"-" json:"total_true"`
		TotalNumTrue float64 `xorm:"-" json:"total_day_true"`
		Date         string  ` json:"date"`
	}
	list := make([]temp, 0)
	err = query.Table("token_inout_daily_sheet").Limit(mList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	for i, v := range list {
		list[i].TotalTrue = t.Int64ToFloat64By8Bit(v.Total)
		list[i].TotalNumTrue = t.Int64ToFloat64By8Bit(v.TotalDayNum)
	}

	mList.Items = list
	return mList, nil
}

// 充币提币走势
func (this *TokenInoutDailySheet) InOutTrendList(filter map[string]interface{}) ([]*InOutTrend, error) {
	// 时间区间，默认最近一周
	today := time.Now().Format(utils.LAYOUT_DATE)

	loc, err := time.LoadLocation("Local")
	if err != nil {
		return nil, errors.NewSys(err)
	}
	todayTime, err := time.ParseInLocation(utils.LAYOUT_DATE, today, loc)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	dateBegin := todayTime.AddDate(0, 0, -7).Format(utils.LAYOUT_DATE)
	dateEnd := today

	// 开始查询
	session := utils.Engine_wallet.Where("1=1")

	// 筛选
	if v, ok := filter["date_begin"]; ok {
		dateBegin, _ = v.(string)
	}
	if v, ok := filter["date_end"]; ok {
		dateEnd, _ = v.(string)
	}
	if v, ok := filter["token_id"]; ok {
		session.And("token_id=?", v)
	}

	var list []*InOutTrend
	err = session.Table(this).
		Select("date, sum(total_day_num) as in_total, sum(total_day_put) as out_total, sum(total_day_num_fee) as fee_total").
		And("date>=?", dateBegin+" 00:00:00").
		And("date<=?", dateEnd+" 00:00:00").
		GroupBy("date").
		OrderBy("date ASC	").
		Find(&list)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	return list, nil
}
