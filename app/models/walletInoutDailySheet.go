package models

import (
	"admin/utils"
	"fmt"
	"time"
)

//数据库 g_wallet 日提币汇总表
type WalletInoutDailySheet struct {
	BaseModel       `xorm:"-"`
	Id              int    `xorm:"not null pk autoincr comment('自增id') TINYINT(4)"`
	TokenId         int    `xorm:"not null comment('货币id') TINYINT(4)"`
	TokenName       string `xorm:"not null comment('货币名称') VARCHAR(20)"`
	TotalDayBalance int64  `xorm:"not null comment('日提币总金额(cny)') BIGINT(20)"`
	TotalDayFee     int64  `xorm:"not null comment('日提币手续费总金额(cny)') BIGINT(20)"`
	TotalBalance    int64  `xorm:"not null comment('累计提币总金额(cny)') BIGINT(20)"`
	TotalFee        int64  `xorm:"not null comment('累计提币手续费总金额(cny)') BIGINT(20)"`
	Date            string `xorm:"not null default 'CURRENT_TIMESTAMP' comment('日期天') TIMESTAMP"`
}

func (this *WalletInoutDailySheet) TableName() string {
	return "token_inout_daily_sheet"
}

//定时结算bibi 日提币报表表数据
func (this *WalletInoutDailySheet) BoottimeTimingSettlement() {

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
		current = current[:10]
		//cunrrentUnixtime := now.Unix()
		//bibi 日报表统计
		engine := utils.Engine_wallet
		list := make([]WalletInoutDailySheet, 0)
		sql := fmt.Sprintf("SELECT SUM(fee_cny) total_day_fee,SUM(amount_cny) total_day_balance,t.days,tokenid token_id ,token_name FROM (SELECT SUBSTRING(created_time,1,10) days,tokenid,token_name, states,opt, fee_cny ,amount_cny FROM g_wallet.token_inout) t  WHERE t.states=2 AND t.opt=1 AND t.days='%s' GROUP BY t.tokenid", current)
		fmt.Println(sql)
		err := engine.SQL(sql).Find(&list)
		if err != nil {
			utils.AdminLog.Println("日提币统计失败！", err.Error(), current)
			//continue
		}
		fmt.Println("1->		", list)
		//根据id 抓取最后插入数据库的的数据
		for _, v := range list {
			if v.TokenId == 0 {
				fmt.Println("error......")
				continue
			}
			lastSql := fmt.Sprintf("SELECT total_balance,total_fee FROM g_wallet.token_inout_daily_sheet WHERE id= (SELECT MAX(id) FROM g_wallet.token_inout_daily_sheet where token_id=%d )", v.TokenId)
			fmt.Println(lastSql)
			_, err := engine.SQL(lastSql).Get(this)
			if err != nil {
				utils.AdminLog.Println(err.Error())
				continue
			}
			fmt.Println("2->		", this)
			this.Id = 0
			this.TokenId = v.TokenId
			this.TokenName = v.TokenName
			this.TotalDayFee = v.TotalDayFee
			this.TotalDayBalance = v.TotalDayBalance
			this.TotalBalance += this.TotalDayBalance
			this.TotalFee += this.TotalDayFee
			this.Date = now.Format("2006-01-02 15:04:05")

			fmt.Println("3->		", this)
			_, err = engine.Table("token_inout_daily_sheet").Cols("token_id","token_name","total_day_balance","total_day_fee","total_balance","total_fee","date").InsertOne(this)
			if err != nil {
				utils.AdminLog.Println(err.Error())
				continue
			}
			fmt.Println("successful")
		}

	}
}

//获取提币日报表信息
func (this *WalletInoutDailySheet) GetInOutDailSheetList(page, rows, tokenId int, date string) (*ModelList, error) {

	engine := utils.Engine_wallet
	query := engine.Desc("id")
	if tokenId != 0 {
		query = query.Where("token_id=?", tokenId)
	}
	if date != `` {
		subst := date[:11] + "23:59:59"
		fmt.Println(subst)
		sql := fmt.Sprintf("date between '%s' and '%s'", date, subst)
		query = query.Where(sql)
	}
	//query = query.Where("id>?", 0)
	countQuery := *query
	count, err := countQuery.Count(&WalletInoutDailySheet{})
	if err != nil {
		return nil, err
	}
	fmt.Println("已经到这里了")
	offset, mList := this.Paging(page, rows, int(count))
	list := make([]WalletInoutDailySheet, 0)
	err = query.Limit(mList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	mList.Items = list
	return mList, nil
}
