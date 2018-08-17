package models

import (
	"admin/utils"
	"fmt"
	"log"
	"time"
	"github.com/robfig/cron"
	"digicon/common/convert"
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
	BalanceAll   int64 `xorm:"comment('币总额') BIGINT(20)"`
	FrozenAll    int64 `xorm:"comment('冻结总额') BIGINT(20)"`
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

//李宇舶 写的
func (tk*TokenDailySheet) TimingFunc(begin, end int64) {
	//g:=make([]*Trade,0)
	//buy
	fmt.Println("定时任务开始--------------------------------------->")
	fmt.Println(time.Now().Unix())
	engine :=utils.Engine_token
	sql := fmt.Sprintf("select sum(num) as a,sum(fee) as b ,sum(fee_cny) as c ,sum(total_cny) as d,token_admission_id  from trade where deal_time>=%d and deal_time<%d  and opt=1 group by token_admission_id", begin, end)
	r, err := engine.Query(sql)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	l := make(map[int]*TokenDailySheet)

	if len(r) > 0 {
		for _, v := range r {
			h:=&TokenDailySheet{}
			t, ok := v["token_admission_id"]
			if !ok {
				log.Fatal("ok u")
			}

			a, ok := v["a"]
			if !ok {
				log.Fatal("ok a")
			}
			b, ok := v["b"]
			if !ok {
				log.Fatal("ok b")
			}
			c, ok := v["c"]
			if !ok {
				log.Fatal("ok c")
			}
			d, ok := v["d"]
			if !ok {
				log.Fatal("ok d")
			}

			h.TokenId  = tk.BytesToIntAscii(t)
			h.BuyTotal= tk.BytesToInt64Ascii(a)
			h.FeeBuyTotal= tk.BytesToInt64Ascii(b)
			h.FeeBuyCny = tk.BytesToInt64Ascii(c)
			h.BuyTotalCny= tk.BytesToInt64Ascii(d)
			l[h.TokenId] = h
		}

	}


	sql = fmt.Sprintf("select sum(num) as a,sum(fee) as b ,sum(fee_cny) as c ,sum(total_cny) as d,token_admission_id  from trade where deal_time>=%d and deal_time<%d  and opt=2 group by token_admission_id", begin, end)
	r, err = engine.Query(sql)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return
	}

	if len(r) > 0 {
		for _, v := range r {
			h:=&TokenDailySheet{}
			t, _ := v["token_admission_id"]
			a, _ := v["a"]
			b, _ := v["b"]
			c, _ := v["c"]
			d, _ := v["d"]

			t_ := convert.BytesToIntAscii(t)
			m, ok := l[t_]
			if !ok {
				h.TokenId  = convert.BytesToIntAscii(t)
				h.SellTotal= convert.BytesToInt64Ascii(a)
				h.FeeSellTotal = convert.BytesToInt64Ascii(b)
				h.FeeSellCny = convert.BytesToInt64Ascii(c)
				h.SellTotalCny = convert.BytesToInt64Ascii(d)
				l[h.TokenId] = h
			} else {
				m.SellTotal  = convert.BytesToInt64Ascii(a)
				m.FeeSellTotal = convert.BytesToInt64Ascii(b)
				m.FeeSellCny = convert.BytesToInt64Ascii(c)
				m.SellTotalCny = convert.BytesToInt64Ascii(d)
			}
		}
	}

	result:=&struct {
		Balance int64
		Frozeen int64
	}{}
	for _,v:=range l  {
		p:=time.Unix(begin, 0).Format("2006-01-02 ")
		utils.AdminLog.Printf("insert into token_id %d,time %s",v.TokenId,p)
		v.Date=begin
		sql :=fmt.Sprintf("SELECT SUM(t.balance) AS balance ,SUM(t.frozen) AS frozen FROM `user_token` t WHERE t.token_id=%d",v.TokenId)
		_,err:=engine.Table("user_token").SQL(sql).Get(result)
		if err!=nil{
			continue
		}
		v.FrozenAll = result.Frozeen
		v.BalanceAll =result.Balance
		_,err = engine.Cols("token_id","fee_buy_cny","fee_buy_total","fee_sell_cny","fee_sell_total","buy_total","sell_total_cny","sell_total","date","balance_all","frozen_all").InsertOne(v)
		if err != nil {
			utils.AdminLog.Errorln(err.Error())
			return
		}
	}
	//sql := fmt.Sprintf("insert into TokenDailySheet (`token_id`,`FeeBuyCny`,`FeeBuyTotal`,`FeeSellCny`,`FeeSellTotal`,`BuyTotal`,`BuyTotalCny`,`SellTotalCny`,`SellTotal`)  values(20001,0,1) on  DUPLICATE key update num=num+values(num)")
	/*
		_,err = DB.GetMysqlConn().Insert(l)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
	*/

	//be:=begin+86400
	//if be>time.Now().Unix() {
	//	return
	//}
	//tk.TimingFunc(begin+86400,end+86400)
	fmt.Println("successful!!!!")
}


func (t *TokenDailySheet ) Run(){
	toBeCharge := time.Now().Format("2006-01-02 ") + "00:00:00"
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
	unix:= theTime.Unix()

	t.TimingFunc(unix-86400,unix)
}


//启动
func DailyStart() {
	fmt.Println("daily count start ....")
	log.Println("daily count start ....")

	i := 0
	c := cron.New()

	//AddFunc
	spec := "0 0 1 * *" // every day ...
	c.AddFunc(spec, func() {
		i++
		log.Println("cron running:", i)
	})
	//AddJob方法
	c.AddJob(spec,&TokenDailySheet{})
	//启动计划任务
	c.Start()
	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer c.Stop()

	select {}
}
