package models

import (
	"admin/utils"
	"fmt"
	"strconv"
	"time"
)

type CurencyFeeDailySheet struct {
	BaseModel  `xorm:"-"`
	Id         int   `xorm:"not null pk comment('自增id') TINYINT(4)"`
	FeeBuyCny  int64 `xorm:"not null comment('法币手买续费折合cny') BIGINT(20)"`
	FeeSellCny int64 `xorm:"not null comment('法币手卖续费折合cny') BIGINT(20)"`
	BalanceCny int64 `xorm:"not null comment('法币交易总额折合cny') BIGINT(20)"`
	Date       int64 `xorm:"not null comment('日期例如20180801') BIGINT(10)"`
}

//定时结算bibi 日交易报表表数据
func (this *CurencyFeeDailySheet) BoottimeTimingSettlement() {

	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		//Printf("定时结算Boottime表数据，结算完成: %v\n",time.Now())
		//以下为定时执行的操作
		current := now.Format("2006-01-02 15:04:05")
		//cunrrentUnixtime := now.Unix()
		type Result struct {
			Num   int64
			Fee   int64
			Price int64
		}
		engine := utils.Engine_currency
		//bibi 日报表统计
		sql := "SELECT m.num,m.fee fee,m.token_id,c.price price FROM (SELECT t.days,t.num,t.fee,t.token_id FROM (SELECT SUBSTRING(confirm_time,1,10) days,num,fee,token_id FROM g_currency.`order` WHERE "

		endSql := " GROUP BY t.token_id) m JOIN  g_token.`config_token_cny` c ON m.token_id= c.token_id"
		//sell ad_type=1
		sellSql := fmt.Sprintf("pay_status=3 AND ad_type=1) tWHERE t.days ='%s'", current[:10])
		listSell := make([]Result, 0)
		err := engine.SQL(sql + sellSql + endSql).Find(&listSell)
		if err != nil {
			utils.AdminLog.Println(err.Error())
			continue
		}
		var cfds CurencyFeeDailySheet
		date, _ := strconv.Atoi(current[:10])
		cfds.Date = int64(date)
		for _, v := range listSell {
			feeStr := this.Int64MulInt64By8BitString(v.Fee, v.Price)
			fResult, err := strconv.ParseFloat(feeStr, 64)
			if err != nil {
				utils.AdminLog.Println(err.Error())
				continue
			}
			cfds.FeeBuyCny += this.Float64ToInt64By8Bit(fResult)

			balanceStr := this.Int64MulInt64By8BitString(v.Num, v.Price)
			fResult, err = strconv.ParseFloat(balanceStr, 64)
			if err != nil {
				utils.AdminLog.Println(err.Error())
				continue
			}
			cfds.BalanceCny += this.Float64ToInt64By8Bit(fResult)
		}
		//buy ad_type=2
		buySql := fmt.Sprintf("pay_status=3 AND ad_type=2) tWHERE t.days ='%s' ", current[:10])
		listBuy := make([]Result, 0)
		err = engine.SQL(sql + buySql + endSql).Find(&listBuy)
		if err != nil {
			utils.AdminLog.Println(err.Error())
			continue
		}
		for _, v := range listSell {
			feeStr := this.Int64MulInt64By8BitString(v.Fee, v.Price)
			fResult, err := strconv.ParseFloat(feeStr, 64)
			if err != nil {
				utils.AdminLog.Println(err.Error())
				continue
			}
			cfds.FeeBuyCny += this.Float64ToInt64By8Bit(fResult)

			balanceStr := this.Int64MulInt64By8BitString(v.Num, v.Price)
			fResult, err = strconv.ParseFloat(balanceStr, 64)
			if err != nil {
				utils.AdminLog.Println(err.Error())
				continue
			}
			cfds.BalanceCny += this.Float64ToInt64By8Bit(fResult)
		}
		_, err = engine.AllCols().InsertOne(&cfds)
		if err != nil {
			utils.AdminLog.Println(err.Error())
			continue
		}
		fmt.Println("CurencyFeeDailySheet--->successfule")

	}
}
