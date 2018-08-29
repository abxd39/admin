package models

import (
	"admin/errors"
	"admin/utils"
	"fmt"
	"time"
)

type MoneyRecord struct {
	BaseModel    `xorm:"-"`
	Id           int64  `xorm:"pk autoincr BIGINT(20)" json:"id"`
	Uid          int    `xorm:"comment('用户ID') unique(hash_index) INT(11)" json:"uid"`
	TokenId      int    `xorm:"comment('代币ID') INT(11)" json:"token_id"`
	Ukey         string `xorm:"comment('联合key') unique(hash_index) VARCHAR(128)" json:"ukey"`
	Type         int    `xorm:"comment('流水类型1区块2委托') INT(11)" json:"type"`
	Opt          int    `xorm:"comment('操作方向1加2减') unique(hash_index) TINYINT(4)" json:"opt"`
	Num          int64  `xorm:"comment('数量') BIGINT(20)" json:"num"`
	Balance      int64  `xorm:"comment('余额') BIGINT(20)" json:"surplus"`
	CreatedTime  int64  `xorm:"comment('操作时间') BIGINT(20)" json:"created_time"`
	TransferTime int64  `xorm:"transfer_time"`
	Comment      string `xorm:"comment('备注') varchar(255)" json:"comment"`
}

type MoneyRecordWithToken struct {
	MoneyRecord `xorm:"extends"`
	TokenName   string `xorm:"token_name" json:"token_name"`
}

func (m *MoneyRecord) TableName() string {
	return "money_record"
}

type MoneyRecordGroup struct {
	MoneyRecord `xorm:"extends"`
	UserInfo    `xorm:"extends"`
}

func (m *MoneyRecord) BackstagePut(count float64, uid, tid int, comment string) error {
	opt := 1 //加
	tp := 14 //后台充值
	ukey := fmt.Sprintf("%d_backstage", time.Now().UnixNano())
	num := m.Float64ToInt64By8Bit(count)
	engine := utils.Engine_token
	sess := engine.NewSession()
	if err := sess.Begin(); err != nil {
		return err
	}
	defer sess.Close()
	ut := new(UserToken)
	has, err := sess.Table("user_token").Where("uid=? and token_id=?", uid, tid).Get(ut)
	if err != nil {
		sess.Rollback()
		return err
	}
	if !has {
		sess.Rollback()
		return errors.New("用户账户不存在该资产！！")
	}

	_, err = sess.Table("money_record").InsertOne(&MoneyRecord{
		Uid:         uid,
		TokenId:     tid,
		Ukey:        ukey,
		Type:        tp,
		Opt:         opt,
		Num:         num,
		Balance:     ut.Balance + num,
		CreatedTime: time.Now().Unix(),
		Comment:     comment,
	})
	if err != nil {
		sess.Rollback()
		return err
	}
	ut.Balance += num
	cnt, err := sess.Table("user_token").Cols("balance").Where("uid=? and token_id=?", uid, tid).Update(ut)
	if err != nil {
		sess.Rollback()
		return err
	}
	if cnt == 0 {
		sess.Rollback()
		return errors.New("充值失败!!!")
	}
	sess.Commit()
	return nil
}

func (m *MoneyRecord) GetMoneyList(page, rows int, uid []int64) (*ModelList, error) {
	engine := utils.Engine_token
	query := engine.Desc("id")
	query = query.In("uid", uid)
	tquery := *query
	count, err := tquery.Count(&MoneyRecordGroup{})
	if err != nil {
		return nil, err
	}
	offset, modelList := m.Paging(page, rows, int(count))
	query.Limit(modelList.PageSize, offset)

	list := make([]MoneyRecordGroup, 0)
	err = query.Find(&list)
	if err != nil {
		return nil, err
	}
	modelList.Items = list
	return modelList, nil
}

func (m *MoneyRecord) GetMoneyListForDateOrType(page, rows, ty, status int, tid int, bt, et uint64, search string) (*ModelList, error) {
	engine := utils.Engine_token
	query := engine.Alias("uch").Desc("u.uid")
	query = query.Join("LEFT", "g_common.user u ", "u.uid= uch.uid")
	query = query.Join("LEFT", "g_common.user_ex ex", "uch.uid=ex.uid")
	query = query.Where("uch.num/10000000 !=0 or uch.balance/100000000 ")
	subDate := time.Now().Unix()
	if bt != 0 {
		if et != 0 {
			query = query.Where("uch.token_id=?  AND uch.created_time BETWEEN ? AND ? ", tid, bt, et+86400)
		} else {
			query = query.Where("uch.token_id=?  AND uch.created_time BETWEEN ? AND ? ", tid, bt, bt+86400)
		}
	} else {
		if tid == 1 {
			query = query.Where("uch.token_id=?  AND uch.created_time BETWEEN ? AND ? ", tid, subDate-3600, subDate)
		} else {
			query = query.Where("uch.token_id=?  AND uch.created_time BETWEEN ? AND ? ", tid, subDate-86400, subDate)
		}

	}

	if search != `` {
		temp := fmt.Sprintf(" concat(IFNULL(u.`uid`,''),IFNULL(u.`phone`,''),IFNULL(ex.`nick_name`,''),IFNULL(u.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}
	if ty != 0 {
		query = query.Where("uch.type=?", ty)
	}
	if status != 0 {
		query = query.Where("u.status=?", status)
	}

	tempQuery := *query
	count, err := tempQuery.Count(&MoneyRecordGroup{})
	if err != nil {
		return nil, err
	}
	offset, modelList := m.Paging(page, rows, int(count))
	list := make([]MoneyRecordGroup, 0)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	for i, v := range list {
		list[i].NumTrue = m.Int64ToFloat64By8Bit(v.Num)
		list[i].SurplusTrue = m.Int64ToFloat64By8Bit(v.Balance)
	}
	modelList.Items = list
	return modelList, nil

}

func (m *MoneyRecord) GetMoneyListForUId(uid []int64) ([]MoneyRecord, error) {
	engine := utils.Engine_token
	query := engine.Desc("uid")
	list := make([]MoneyRecord, 0)
	err := query.In("uid", uid).Find(&list)
	if err != nil {
		return nil, err
	}
	//for i, v := range list {
	//	list[i].NumTrue = m.Int64ToFloat64By8Bit(v.Num)
	//	list[i].SurplusTrue = m.Int64ToFloat64By8Bit(v.Balance)
	//}
	//modelList.Items = list
	return list, nil

}

//u优化

//流水列表
func (s *MoneyRecord) List(pageIndex, pageSize int, filter map[string]interface{}) (*ModelList, []*MoneyRecordWithToken, error) {
	query := utils.Engine_token.Alias("mr").Join("LEFT", []string{new(UserToken).TableName(), "ut"}, "ut.token_id=mr.token_id AND ut.uid=mr.uid").Where("1=1")

	//筛选
	orderBy := "mr.id DESC"
	if _, ok := filter["uid"]; ok {
		query.And("mr.uid=?", filter["uid"])
	}
	if _, ok := filter["transfer"]; ok { //划转流水
		query.And("mr.type IN (?,?)", 10, 11)
		orderBy = "mr.transfer_time DESC, mr.id DESC"
	}
	if v, ok := filter["type"]; ok {
		query.And("mr.type=?", v)
	}
	if v, ok := filter["transfer_date"]; ok {
		loc, _ := time.LoadLocation("Local")
		dayBeginTime, _ := time.ParseInLocation(utils.LAYOUT_DATE_TIME, v.(string)+" 00:00:00", loc)
		dayEndTime, _ := time.ParseInLocation(utils.LAYOUT_DATE_TIME, v.(string)+" 23:59:59", loc)

		query.And("mr.transfer_time>=?", dayBeginTime.Unix()).And("mr.transfer_time<=?", dayEndTime.Unix())
	}
	if v, ok := filter["transfer_time_begin"]; ok {
		query.And("mr.transfer_time>=?", v)
	}
	if v, ok := filter["transfer_time_end"]; ok {
		query.And("mr.transfer_time<=?", v)
	}

	//分页
	tmpQuery := *query
	total, err := tmpQuery.Count(s)
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}
	offset, modelList := s.Paging(pageIndex, pageSize, int(total))

	var list []*MoneyRecordWithToken
	err = query.Select("mr.*, ut.token_name").OrderBy(orderBy).Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}
	modelList.Items = list

	return modelList, list, nil
}

//拉去平台内充值的所有记录 平台当天的单个币的充币数量
func (m *MoneyRecord) GetPlatformAll(page, rows int, tid, bt, et uint64) (*ModelList, error) {
	engine := utils.Engine_token
	sql := "SELECT FROM_UNIXTIME(created_time,'%Y%m%d') day,created_time `time` ,TYPE,SUM(num) total ,token_id tid,uid FROM g_token.money_record WHERE TYPE=14 "
	var condition string
	if bt != 0 {
		if et != 0 {
			condition = fmt.Sprintf(" WHERE TYPE=6 and created_time BETWEEN %d AND %d ", bt, et+86400)
		} else {
			condition = fmt.Sprintf(" WHERE TYPE=6 and created_time BETWEEN %d AND %d ", bt, bt+86400)
		}
	}
	if tid != 0 {
		condition += fmt.Sprintf("  AND token_id=%d ", tid)
	}

	condition += " GROUP BY day , token_id  order by day desc "
	count := fmt.Sprintf("SELECT COUNT(*) COUNT FROM (%s)t", sql+condition)
	num := &struct {
		Count int
	}{}
	_, err := engine.SQL(count).Get(num)
	if err != nil {
		return nil, err
	}
	offset, mList := m.Paging(page, rows, int(num.Count))
	limit := fmt.Sprintf(" limit %d offset %d ", mList.PageSize, offset)
	type temp struct {
		Time      int64   `json:"day"`
		Total     int64   `json:"total"`
		TotalTrue float64 `xorm:"-" json:"total_true"`
		Name      string  ` xorm:"-" json:"name"`
		Tid       int     `json:"tid"`
	}
	list := make([]temp, 0)
	err = engine.SQL(sql + condition + limit).Find(&list)
	if err != nil {
		return nil, err
	}
	tl, err := new(Tokens).GetTokensList()
	if err != nil {
		return nil, err
	}
	for i, v := range list {
		for _, tv := range tl {
			if v.Tid == tv.Id {
				list[i].Name = tv.Mark
				break
			}
		}
		list[i].TotalTrue = m.Int64ToFloat64By8Bit(v.Total)
	}
	mList.Items = list
	return mList, nil
}

//平台内充值明细
func (m *MoneyRecord) GetPlatForTokenOfDay(page, rows, uid, tid int, date uint64) (*ModelList, error) {
	engine := utils.Engine_token
	sql := "SELECT FROM_UNIXTIME(created_time,'%Y-%m-%d %H:%i:%s') day, SUM(num) total ,uid,comment FROM g_token.money_record  "

	var condition string
	//condition = fmt.Sprintf("WHERE TYPE=14 AND token_id=%d AND created_time  BETWEEN %d AND %d ", tid, date, date+86400)
	condition = fmt.Sprintf("WHERE TYPE=14 AND token_id=%d AND created_time  BETWEEN %d AND %d ", tid, date, date+86400)

	if uid != 0 {
		condition += fmt.Sprintf(" AND uid=%d ", uid)
	}
	sql += condition
	sql += " GROUP BY uid , token_id "
	countSql := fmt.Sprintf("select count(*) count from (%s) t", sql)
	num := &struct {
		Count int
	}{}
	_, err := engine.SQL(countSql).Get(num)
	if err != nil {
		return nil, err
	}
	fmt.Println(num.Count)
	type tmep struct {
		Day       string  `json:"day"`
		Total     int64   `json:"total"`
		TotalTrue float64 `xorm:"-" json:"total_true"`
		Uid       int64   `json:"uid"`
		Comment   string  `json:"comment"`
	}
	list := make([]tmep, 0)
	offset, mList := m.Paging(page, rows, int(num.Count))
	limit := fmt.Sprintf("limit %d offset %d", mList.PageSize, offset)
	err = engine.SQL(sql + limit).Find(&list)
	if err != nil {
		return nil, err
	}
	for i, v := range list {
		list[i].TotalTrue = m.Int64ToFloat64By8Bit(v.Total)
	}
	fmt.Println(len(list))
	mList.Items = list
	return mList, nil
}
