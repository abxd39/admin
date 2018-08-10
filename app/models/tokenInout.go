package models

import (
	"admin/errors"
	"admin/utils"
	"fmt"
	"time"
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
	NickName   string `json:"nick_name"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Status     int    `json:"status"`
	AmountTrue float64 `xorm:"-" json:"amount_ture"`
	FeeTrue float64 `xorm:"-" json:"fee_true"`
	ToCount float64 `xorm:"-" json:"to_count"`
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
func (t *TokenInout) GetTotalInfoList(page, rows, tid, opt int, date, search string) (*ModelList, error) {
	enginge := utils.Engine_wallet
	//SELECT t.time,t.token_name,t.total,t.uid
	sql1 := " FROM (SELECT DATE_FORMAT(created_time,'%Y%m%d') DAY,created_time time ,opt,SUM(amount) total ,tokenid,token_name name,uid FROM token_inout "
	sql := fmt.Sprintf("WHERE %d", opt)
	if tid != 0 {
		tmp := fmt.Sprintf(" AND tokenid=%d", tid)
		sql += tmp
	}
	//sql:=fmt.Sprintf("  opt= %d AND tokenid=%d GROUP BY DAY, uid)t WHERE t.day=",opt,tid)
	//刷选
	if search != `` {
		tmp := fmt.Sprintf(" AND uid=%s", search)
		sql += tmp
	}
	sql += " GROUP BY DAY, uid)t WHERE t.day="
	sql = sql1 + sql
	if date != `` {
		sub := date[:8]
		sql = sql + sub
	}

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
		Total uint64 //提币总数
		Name  string //货币名称
	}
	limitSql := fmt.Sprintf(" limit %d offset %d ", mList.PageSize, offset)
	list := make([]Return, 0)
	contentSql := "SELECT t.time,t.name,t.total,t.uid ,t.day" + sql + limitSql
	fmt.Println(contentSql)
	err = enginge.SQL(contentSql).Find(&list)
	if err != nil {
		return nil, err
	}
	mList.Items = list
	return mList, nil
}

//日提币汇总
func (t *TokenInout) GetTotalList(page, rows, tokenId, opt int, date string) (*ModelList, error) {
	engine := utils.Engine_wallet
	sql1 := "FROM (SELECT DATE_FORMAT(created_time,'%Y%m%d') DAY,id,opt,SUM(amount) total,token_name name,tokenid tid FROM token_inout WHERE "
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
		Day   int    //日期
		Total uint64 //提币总量
		Name  string // 货币名称
		Tid   int    //货币id
	}
	offset, mList := t.Paging(page, rows, int(count.Count))
	limitSql = fmt.Sprintf(limitSql, mList.PageSize, offset)
	list := make([]Return, 0)
	sql += search + limitSql
	queryContent := "SELECT * " + sql
	//fmt.Println(queryContent)
	engine.SQL(queryContent).Find(&list)
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
		query = query.Where("t.token_id=?", tokenId)
	}
	if status != 0 {
		query = query.Where("u.states=?", status)
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
	if status != 0 {
		query = query.Where("u.status=?", status)
	}
	if len(search) != 0 {
		temp := fmt.Sprintf(" concat(IFNULL(u.`uid`,''),IFNULL(u.`phone`,''),IFNULL(ex.`nick_name`,''),IFNULL(u.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}
	queryCount:=*query
	count,err:=queryCount.Count(t)
	if err!=nil{
		return nil, err
	}
	offset,mList:=t.Paging(page,rows,int(count))
	list:=make([]TokenInoutGroup,0)
	err=query.Limit(mList.PageSize,offset).Find(&list)
	if err!=nil{
		return nil, err
	}
	for i,v:=range list{
		list[i].FeeTrue = t.Int64ToFloat64By8Bit(v.Fee)
		list[i].AmountTrue =t.Int64ToFloat64By8Bit(v.Amount)
		list[i].ToCount = list[i].AmountTrue-list[i].FeeTrue
	}
	mList.Items = list

	//两个方向  用户信息库和 钱包库
	//if uStatus != 0 || search != `` {
	//	mList, err := new(WebUser).GetAllUser(page, rows, uStatus, search)
	//	if err != nil {
	//		return nil, err
	//	}
	//	value, ok := mList.Items.([]UserGroup)
	//	if !ok {
	//		return nil, errors.New("assert []webUser type failed!!")
	//	}
	//	uidList := make([]int64, 0)
	//	for _, v := range value {
	//		uidList = append(uidList, v.Uid)
	//	}
	//	if len(uidList) < 1 {
	//		//没有匹配刷选条件的用户
	//		return nil, nil
	//	}
	//	query = query.In("uid", uidList)
	//	countQuery := *query
	//	count, err := countQuery.Count(&TokenInout{})
	//	offset, modelList := t.Paging(page, rows, int(count))
	//	list := make([]TokenInoutGroup, 0)
	//	err = query.Limit(modelList.PageSize, offset).Find(&list)
	//	if err != nil {
	//		return nil, err
	//	}
	//	//
	//	for i, _ := range list {
	//		for _, v := range value {
	//			if list[i].Uid == int(v.Uid) {
	//				list[i].NickName = v.NickName
	//				list[i].Phone = v.Phone
	//				list[i].Email = v.Email
	//				list[i].Status = v.Status
	//				break
	//			}
	//		}
	//	}
	//	modelList.Items = list
	//	return modelList, nil
	//} else {
	//	if tokenId != 0 {
	//		query = query.Where("token_id=?", tokenId)
	//	}
	//	if status != 0 {
	//		query = query.Where("states=?", status)
	//	}
	//	if opt != 0 {
	//		query = query.Where("opt=?", opt)
	//	}
	//	if date != `` {
	//		subst := date[:11] + "23:59:59"
	//		fmt.Println(subst)
	//		sql := fmt.Sprintf("create_time  BETWEEN '%s' AND '%s' ", date, subst)
	//		query = query.Where(sql)
	//	}
	//	countQuery := *query
	//	count, err := countQuery.Count(&TokenInout{})
	//	if err != nil {
	//		return nil, err
	//	}
	//	offset, modelList := t.Paging(page, rows, int(count))
	//	list := make([]TokenInoutGroup, 0)
	//	err = query.Limit(modelList.PageSize, offset).Find(&list)
	//	if err != nil {
	//		return nil, err
	//	}
	//	uidList := make([]uint64, 0)
	//	for _, v := range list {
	//		uidList = append(uidList, uint64(v.Uid))
	//	}
	//	if len(uidList) < 1 {
	//		//没有匹配刷选条件的用户
	//		return nil, nil
	//	}
	//	uList, err := new(UserGroup).GetUserListForUid(uidList)
	//	if err != nil {
	//		return nil, err
	//	}
	//	for i, _ := range list {
	//		for _, v := range uList {
	//			if int(v.Uid) == list[i].Uid {
	//				list[i].TokenName = v.NickName
	//				list[i].Phone = v.Phone
	//				list[i].Email = v.Email
	//				list[i].Status = v.Status
	//				break
	//			}
	//		}
	//	}
	//	modelList.Items = list
	//	return modelList, nil
	//} //else
	return mList, nil
}

//提币管理
func (t *TokenInout) OptTakeToken(id, uid, status int) error {
	engine := utils.Engine_wallet
	//t:=new(TokenInout)
	has, err := engine.Where("id=? and uid=?", id, uid).Get(t)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("rescind failed !!")
	}

	sess := engine.NewSession()
	if err = sess.Begin(); err != nil {
		return err
	}
	defer sess.Close()

	_, err = sess.Where("id=? and uid=?", id, uid).Update(&TokenInout{
		States: status,
	})
	if err != nil {
		sess.Rollback()
		return err
	}

	sess.Commit()
	return nil

}
