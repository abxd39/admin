package models

import (
	"admin/apis"
	"admin/errors"
	"admin/utils"
	"fmt"
	"time"
	"admin/utils/convert"
)

//冲 提 币明细流水表

type TokenInout struct {
	BaseModel   `xorm:"-"`
	Id          int    `xorm:"not null pk autoincr comment('自增id') INT(11)"`
	Uid         int    `xorm:"not null comment('用户id') INT(11)"`
	Opt         int    `xorm:"not null comment('操作方向 1 充币 2 提币') TINYINT(4)"`
	Txhash      string `xorm:"not null comment('交易hash') VARCHAR(200)"`
	From        string `xorm:"not null comment('打款方') VARCHAR(42)"`
	To          string `xorm:"not null comment('收款方') VARCHAR(42)"`
	Amount      int64  `xorm:"not null comment('金额(数量)') BIGINT(20)"`
	Fee         int64  `xorm:"not null comment('提币手续费(数量)') BIGINT(20)"`
	AmountCny   int64  `xorm:"not null comment('提币数量折合cny') BIGINT(20)"`
	FeeCny      int64  `xorm:"not null comment('手续费折合cny') BIGINT(20)"`
	Value       string `xorm:"not null comment('原始16进制转账数据') VARCHAR(32)"`
	Chainid     int    `xorm:"not null comment('链id') INT(11)"`
	Contract    string `xorm:"not null default '' comment('合约地址') VARCHAR(42)"`
	Tokenid     int    `xorm:"not null comment('币种id') INT(11)"`
	States      int    `xorm:"not null comment('人充提币状态 1正在提币，2 已完成，3提币已取消，4提币失败') TINYINT(1)"`
	TokenName   string `xorm:"not null comment('币种名称') VARCHAR(10)"`
	CreatedTime string `xorm:"not null default 'CURRENT_TIMESTAMP' comment('创建时间 提币创建时间') TIMESTAMP"`
	DoneTime    string `xorm:"not null default '0000-00-00 00:00:00' comment('充币到账时间') TIMESTAMP"`
	Remarks     string `xorm:"not null comment('备注信息') VARCHAR(100)"`
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

//仪表盘 日提币手续费
func (t *TokenInout) GetOutTokenFee() (float64, error) {
	engine := utils.Engine_wallet
	current := time.Now().Format("2006-01-02 15:04:05")
	sql := fmt.Sprintf("SELECT SUM(t.fee)FROM (SELECT SUBSTRING(done_time,1,10)days,fee_cny fee FROM g_wallet.`token_inout` WHERE states=2)t WHERE t.days='%s'", current[:10])
	fee := &struct {
		Fee float64
	}{}
	_, err := engine.SQL(sql).Get(fee)
	if err != nil {
		utils.AdminLog.Println(err.Error())
		return 0, err
	}
	return fee.Fee, nil
}

//日提币 每个用户提币信息
func (t *TokenInout) GetTotalInfoList(page, rows, tid, opt int, search string) (*ModelList, error) {
	enginge := utils.Engine_wallet
	//SELECT t.time,t.token_name,t.total,t.uid
	sql1 := " FROM (SELECT DATE_FORMAT(created_time,'%Y%m%d') DAY,created_time time ,opt,SUM(amount) amount ,tokenid,token_name name,uid FROM token_inout "
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
	sql += " GROUP BY DAY, uid)t "
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
		Day   int
		Uid   int
		Amount int64
		Total string `xorm:"-"` //提币总数
		Name  string //货币名称
	}
	limitSql := fmt.Sprintf(" limit %d offset %d ", mList.PageSize, offset)
	list := make([]Return, 0)
	contentSql := "SELECT t.time,t.name,t.amount,t.uid ,t.day" + sql + limitSql
	fmt.Println(contentSql)
	err = enginge.SQL(contentSql).Find(&list)
	if err != nil {
		return nil, err
	}
	for i,v:=range list{
		list[i].Total = convert.Int64ToStringBy8Bit(v.Amount)
	}
	mList.Items = list
	return mList, nil
}

//日提币 充币 汇总
func (t *TokenInout) GetTotalList(page, rows, tokenId, opt int, date string) (*ModelList, error) {
	engine := utils.Engine_wallet
	sql1 := "FROM (SELECT DATE_FORMAT(created_time,'%Y%m%d') DAY,id,opt,SUM(amount) count,token_name name,tokenid tid FROM token_inout WHERE "
	sql := fmt.Sprintf("opt= %d GROUP BY DAY, tokenid) t ", opt)
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
	if status == utils.VERIFY_OUT_TOKEN_MARK {
		fmt.Println("审核通过")
		mount := t.Int64ToFloat64By8Bit(t.Amount)
		////fmt.Println("num=",)
		strMount := fmt.Sprintf("%.10f", mount)
		if token.Signature == "eip155" || token.Signature == "eip" { //ERC20
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
		if token.Signature == "btc" { //btc
			fmt.Println("btc 提币申请")
			err = new(apis.VendorApi).PostOutTokenBtc(t.Uid, t.Tokenid, t.Id, t.To, strMount)
			if err != nil {
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
