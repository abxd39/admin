package models

import (
	"admin/apis"
	"admin/errors"
	"admin/utils"
	"admin/utils/convert"
	"fmt"
	"time"
	"encoding/json"
	"strings"
)

//冲 提 币明细流水表

//type TokenInout struct {
//	BaseModel   `xorm:"-"`
//	Id          int    `xorm:"not null pk autoincr comment('自增id') INT(11)"`
//	Uid         int    `xorm:"not null comment('用户id') INT(11)"`
//	Opt         int    `xorm:"not null comment('操作方向 1 充币 2 提币') TINYINT(4)"`
//	Txhash      string `xorm:"not null comment('交易hash') VARCHAR(200)"`
//	From        string `xorm:"not null comment('打款方') VARCHAR(42)"`
//	To          string `xorm:"not null comment('收款方') VARCHAR(42)"`
//	Amount      int64  `xorm:"not null comment('金额(数量)') BIGINT(20)"`
//	Fee         int64  `xorm:"not null comment('提币手续费(数量)') BIGINT(20)"`
//	AmountCny   int64  `xorm:"not null comment('提币数量折合cny') BIGINT(20)"`
//	FeeCny      int64  `xorm:"not null comment('手续费折合cny') BIGINT(20)"`
//	Value       string `xorm:"not null comment('原始16进制转账数据') VARCHAR(32)"`
//	Chainid     int    `xorm:"not null comment('链id') INT(11)"`
//	Contract    string `xorm:"not null default '' comment('合约地址') VARCHAR(42)"`
//	Tokenid     int    `xorm:"not null comment('币种id') INT(11)"`
//	States      int    `xorm:"not null comment('人充提币状态 1正在提币，2 已完成，3提币已取消，4提币失败') TINYINT(1)"`
//	TokenName   string `xorm:"not null comment('币种名称') VARCHAR(10)"`
//	CreatedTime string `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间 提币创建时间') TIMESTAMP"`
//	DoneTime    string `xorm:"not null default '0000-00-00 00:00:00' comment('充币到账时间') TIMESTAMP"`
//	Remarks     string `xorm:"not null comment('备注信息') VARCHAR(100)"`
//}


type TokenInout struct {
	BaseModel   `xorm:"-"`
	Id          int       `xorm:"not null pk autoincr comment('自增id') index INT(11)"`
	Uid         int       `xorm:"not null comment('用户id') INT(11)"`
	Opt         int       `xorm:"not null comment('操作方向 1 充币 2 提币') TINYINT(4)"`
	Txhash      string    `xorm:"not null comment('交易hash') VARCHAR(191)"`
	From        string    `xorm:"not null comment('打款方') VARCHAR(42)"`
	To          string    `xorm:"not null comment('收款方') VARCHAR(42)"`
	Amount      int64     `xorm:"not null comment('金额(数量)') BIGINT(20)"`
	Fee         int64     `xorm:"not null comment('提币手续费(数量)') BIGINT(20)"`
	AmountCny   int64     `xorm:"not null comment('提币数量折合cny') BIGINT(20)"`
	FeeCny      int64     `xorm:"not null comment('手续费折合cny') BIGINT(20)"`
	Value       string    `xorm:"not null comment('原始16进制转账数据') VARCHAR(32)"`
	Chainid     int       `xorm:"not null comment('链id') INT(11)"`
	Contract    string    `xorm:"not null default '' comment('合约地址') VARCHAR(42)"`
	Tokenid     int       `xorm:"not null comment('币种id') INT(11)"`
	States      int       `xorm:"not null comment('人充提币状态 1正在提币，2 已完成，3提币已取消，4提币失败') TINYINT(1)"`
	TokenName   string    `xorm:"not null comment('币种名称') VARCHAR(10)"`
	CreatedTime string `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间 提币创建时间') TIMESTAMP"`
	DoneTime    string `xorm:"not null default '0000-00-00 00:00:00' comment('充币到账时间') TIMESTAMP"`
	Remarks     string    `xorm:"not null comment('备注信息') VARCHAR(100)"`
	Gas         int64     `xorm:"comment('gas数量') BIGINT(20)"`
	GasPrice    int64     `xorm:"comment('gas价格,单位：wei') BIGINT(20)"`
	RealFee     int64     `xorm:"comment('实际消耗手续费') BIGINT(20)"`
}


type TokenInoutGroup struct {
	TokenInout `xorm:"extends"`
	NickName   string  `json:"nick_name"`
	Phone      string  `json:"phone"`
	Email      string  `json:"email"`
	Status     int     `json:"status"`
	AmountTrue float64 `xorm:"-" json:"amount_ture"`
	FeeTrue    float64 `xorm:"-" json:"fee_true"`
	OutCount   float64 `xorm:"-" json:"out_count"` //提币数量
}

func (t *TokenInoutGroup) TableName() string {
	return "token_inout"
}

type TokenFeeHistoryGroup struct {
	//TokenInout `xorm:"-"`
	Id          int `json:"id"`
	Uid         int `json:"uid"`
	Amount      int64 `json:"amount"`
	Fee         int64 `json:"fee"`
	Mark         string `json:"mark"`
	NumTrue      string `json:"num_true"`
	FeeTrue      string `json:"fee_true"`
	CreatedTime string `json:"created_time"`
}

func (this *TokenFeeHistoryGroup) TableName() string {
	return "token_inout"
}

//p5-1-1-1提币手续费明细
func (this *TokenInout) GetAddTakeList(page, rows, tid, uid int) (*ModelList, error) {
	fmt.Println("p5-1-1-1提币手续费明细")
	engine := utils.Engine_wallet
	query := engine.Alias("ti").Desc("ti.id")
	query = query.Join("LEFT", "g_common.tokens t", "t.id=ti.tokenid")
	query = query.Where("ti.states=2")//已完成
	if tid != 0 {
		query = query.Where("ti.tokenid=?", tid)
	}
	if uid != 0 {
		query = query.Where("ti.uid=?", uid)
	}
	//if date != 0 {
	//	query = query.Where("check_time BETWEEN ? AND ?", date, date+864000)
	//}
	countQuery := *query
	count, err := countQuery.Count(&TokenInout{})
	if err != nil {
		return nil, err
	}

	offset, mlist := this.Paging(page, rows, int(count))
	list := make([]TokenFeeHistoryGroup, 0)
	err = query.Limit(mlist.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	for i, v := range list {
		list[i].FeeTrue = convert.Int64ToStringBy8Bit(v.Fee)
		list[i].NumTrue = convert.Int64AddInt64Float64Percent(v.Amount,v.Fee)
	}
	mlist.Items = list
	return mlist, nil
}


//日提币 每个用户提币信息
func (t *TokenInout) GetTotalInfoList(page, rows, tid, opt int, search string) (*ModelList, error) {
	enginge := utils.Engine_wallet
	//SELECT t.time,t.token_name,t.total,t.uid
	sql1 := " FROM (SELECT DATE_FORMAT(created_time,'%Y-%m-%d %H:%i:%s') DAY,created_time time ,opt,SUM(amount) amount ,tokenid,token_name name,uid FROM token_inout "
	sql := fmt.Sprintf("WHERE opt=%d", opt)
	if tid != 0 {
		tmp := fmt.Sprintf(" AND tokenid=%d", tid)
		sql += tmp
	}

	//刷选
	if search != `` {
		tmp := fmt.Sprintf(" AND uid=%s", search)
		sql += tmp
	}
	sql += " GROUP BY DAY, uid)t  WHERE t.amount!=0 "
	sql = sql1 + sql
	//if date != `` {
	//	sub := date[:8]
	//	sql = sql + sub
	//}

	type Count struct {
		Count int
	}
	count := new(Count)
	sqlCount := "select count(*) count " + sql
	fmt.Println(sqlCount)
	_, err := enginge.SQL(sqlCount).Get(count)
	if err != nil {
		return nil, err
	}
	offset, mList := t.Paging(page, rows, int(count.Count))
	type Return struct {
		Day    string
		Uid    int
		Amount int64
		Total  string `xorm:"-"` //提币总数
		Name   string //货币名称
	}
	limitSql := fmt.Sprintf(" limit %d offset %d ", mList.PageSize, offset)
	list := make([]Return, 0)
	contentSql := "SELECT t.time,t.name,t.amount,t.uid ,t.day" + sql + limitSql
	fmt.Println(contentSql)
	err = enginge.SQL(contentSql).Find(&list)
	if err != nil {
		return nil, err
	}
	for i, v := range list {
		list[i].Total = convert.Int64ToStringBy8Bit(v.Amount)
	}
	mList.Items = list
	return mList, nil
}

//日提币 充币 汇总
func (t *TokenInout) GetTotalList(page, rows, tokenId, opt int, date string) (*ModelList, error) {
	engine := utils.Engine_wallet
	sql1 := "FROM (SELECT DATE_FORMAT(created_time,'%Y%m%d') DAY,id,opt,SUM(amount) count,token_name name,tokenid tid FROM token_inout WHERE "
	sql := fmt.Sprintf("opt= %d GROUP BY DAY, tokenid) t WHERE t.amount!=0 ", opt)
	sql = sql1 + sql
	limitSql := " limit %d offset %d"
	search := "where t.id>0"
	if tokenId != 0 {
		temp := fmt.Sprintf(" AND t.tid=%d", tokenId)
		search += temp
	}
	if date != `` {
		sub := date[:8]
		fmt.Println("date=", sub, "len(date) =", len(sub))
		temp := " AND t.day=" + sub
		search += temp
	}
	type Count struct {
		Count int
	}
	count := new(Count)
	query := "SELECT COUNT(*) count  " + sql + search
	//fmt.Println("query=",query)
	_, err := engine.SQL(query).Get(count)
	if err != nil {
		return nil, err
	}
	fmt.Println("count=", count.Count)
	type Return struct {
		Day   int     //日期
		Count int64   //总数
		Total float64 //提币总量
		Name  string  // 货币名称
		Tid   int     //货币id
	}
	offset, mList := t.Paging(page, rows, int(count.Count))
	limitSql = fmt.Sprintf(limitSql, mList.PageSize, offset)
	list := make([]Return, 0)
	sql += search + limitSql
	queryContent := "SELECT * " + sql
	//fmt.Println(queryContent)
	engine.SQL(queryContent).Find(&list)
	for i, v := range list {
		list[i].Total = t.Int64ToFloat64By8Bit(v.Count)
	}
	mList.Items = list
	return mList, nil
}

//提币 充币 p3-1-0 充币 提币管理
func (t *TokenInoutGroup) GetTokenInList(page, rows, uStatus, status, tokenId, opt int, search string) (*ModelList, error) {
	engine := utils.Engine_wallet
	query := engine.Alias("t").Desc("t.uid")
	query = query.Join("LEFT", "g_common.user u", "u.uid=t.uid")
	query = query.Join("LEFT", "g_common.user_ex ex", "ex.uid=t.uid")
	if tokenId != 0 {
		query = query.Where("t.tokenid=?", tokenId)
	}
	if status != 0 {
		query = query.Where("t.states=?", status)
	}
	if opt != 0 {
		query = query.Where("t.opt=?", opt)
	}
	//if date != `` {
	//	subst := date[:11] + "23:59:59"
	//	fmt.Println(subst)
	//	sql := fmt.Sprintf("t.create_time  BETWEEN '%s' AND '%s' ", date, subst)
	//	query = query.Where(sql)
	//}
	if uStatus != 0 {
		query = query.Where("u.status=?", uStatus)
	}
	if len(search) != 0 {
		temp := fmt.Sprintf(" concat(IFNULL(u.`uid`,''),IFNULL(u.`phone`,''),IFNULL(ex.`nick_name`,''),IFNULL(u.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}
	queryCount := *query
	count, err := queryCount.Count(t)
	if err != nil {
		return nil, err
	}
	offset, mList := t.Paging(page, rows, int(count))
	list := make([]TokenInoutGroup, 0)
	err = query.Limit(mList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	for i, v := range list {
		list[i].FeeTrue = t.Int64ToFloat64By8Bit(v.Fee)
		list[i].AmountTrue = t.Int64ToFloat64By8Bit(v.Amount)
		list[i].OutCount = list[i].AmountTrue + list[i].FeeTrue
	}
	mList.Items = list

	return mList, nil
}

//提币管理
func (t *TokenInout) OptTakeToken(id, status int) error {
	engine := utils.Engine_wallet
	//t:=new(TokenInout)
	has, err := engine.Where("id=? ", id).Get(t)
	if err != nil {
		return err
	}
	if !has {
		strErr := fmt.Sprintf("订单不存在!!! id=%d", id)
		return errors.New(strErr)
	}
	engineToken := utils.Engine_common
	token := new(Tokens)
	has, err = engineToken.Table("tokens").Where("id=?", t.Tokenid).Get(token)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if !has {
		strErr := fmt.Sprintf("数字货币不存在!!!token_id=%d", t.Tokenid)
		return errors.New(strErr)
	}
	sess := engine.NewSession()
	if err = sess.Begin(); err != nil {
		return err
	}
	defer sess.Close()

	_, err = sess.Where("id=?", id).Update(&TokenInout{
		States: status,
	})
	if err != nil {
		sess.Rollback()
		return err
	}
	//审核通过
	fmt.Println("id=", id, "status=", status)
	if status == utils.VERIFY_OUT_TOKEN_MARK {
		fmt.Println("审核通过")
		mount := t.Int64ToFloat64By8Bit(t.Amount)
		////fmt.Println("num=",)
		strMount := fmt.Sprintf("%.10f", mount)
		Name := strings.ToUpper(token.Signature)
		Name =strings.Trim(Name," ")
		if Name == "EIP155" || Name == "EIP" { //ERC20
			fmt.Sprintf(strMount)
			fmt.Println("获取签名")
			sign, err := new(apis.VendorApi).GetTradeSigntx(t.Uid, t.Tokenid, t.To, strMount)
			if err != nil {
				sess.Rollback()
				return err
			}
			fmt.Println("发送提币申请")
			err = new(apis.VendorApi).PostOutToken(t.Uid, t.Tokenid, t.Id, sign)
			if err != nil {
				sess.Rollback()
				return err
			}
		}
		if Name == "BTC" { //btc
			fmt.Println("btc 提币申请")
			err = new(apis.VendorApi).PostOutTokenBtc(t.Uid, t.Tokenid, t.Id, t.To, strMount)
			if err != nil {
				sess.Rollback()
				return err
			}
		}
		if Name =="OMNI"{ //表示提USDT
			fmt.Println("usdt 提币申请")
			params := make(map[string]interface{})
			params["uid"]=t.Uid
			params["token_id"]=t.Tokenid
			params["apply_id"]=t.Id
			params["from_address"]=t.From
			params["to_address"]= t.To
			params["protertyid"]=1 //测试环境 为1 正式环境为=31
			params["amount"]= convert.Int64ToStringBy8Bit(t.Amount)
			result,_:=json.Marshal(params)
			err =new(apis.VendorApi).PostOutTokenUsdt(string(result))
			if err!=nil{
				utils.AdminLog.Infof("usdt 提币失败",err.Error())
				fmt.Println(err.Error())
				sess.Rollback()
				return err
			}
		}

	}
	//审核撤销
	if status == utils.VERIFY_REVOKE_TOKEN_MARK {
		//需要王炳雨提供接口
		err := new(apis.VendorApi).RevokeOutToken(int64(t.Uid), int64(t.Tokenid), t.Amount+t.Fee)
		if err != nil {
			sess.Rollback()
			utils.AdminLog.Error(err.Error())
			return err
		}
	}
	sess.Commit()
	return nil

}

// 交易合计
type InOutTradeTotal struct {
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
func (this *TokenInout) TradeTotal() (*InOutTradeTotal, error) {
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
	feeTotal := &InOutTradeTotal{}
	_, err = utils.Engine_wallet.
		Table(this).
		Select("COUNT(id) total_time, IFNULL(SUM(amount+fee), 0) total_num, IFNULL(SUM(fee), 0) total").
		Where("opt=2").
		Get(feeTotal)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	//2. 今日
	todayFeeTotal := &InOutTradeTotal{}
	_, err = utils.Engine_wallet.
		Table(this).
		Select("COUNT(id) today_total_time, IFNULL(SUM(amount+fee), 0) today_total_num, IFNULL(SUM(fee), 0) today_total").
		Where("opt=2").
		And("created_time>=?", todayDate).
		Get(todayFeeTotal)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	//3. 上日
	yesFeeTotal := &InOutTradeTotal{}
	_, err = utils.Engine_wallet.
		Table(this).
		Select("COUNT(id) yesterday_total_time, IFNULL(SUM(amount+fee), 0) yesterday_total_num, IFNULL(SUM(fee), 0) yesterday_total").
		Where("opt=2").
		And("created_time>=?", yesterdayDateBegin).
		And("created_time<=?", yesterdayDateEnd).
		Get(yesFeeTotal)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	//4. 上周同日
	lastWeekFeeTotal := &InOutTradeTotal{}
	_, err = utils.Engine_wallet.
		Table(this).
		Select("COUNT(id) last_week_day_total_time, IFNULL(SUM(amount+fee), 0) last_week_day_total_num, IFNULL(SUM(fee), 0) last_week_day_total").
		Where("opt=2").
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
