package models

import (
	"admin/errors"
	"admin/utils"
	"fmt"
	"time"
)

type MoneyRecord struct {
	BaseModel   `xorm:"-"`
	UserInfo    `xorm:"extends"`
	Id          int64  `xorm:"pk autoincr BIGINT(20)" json:"id"`
	Uid         int    `xorm:"comment('用户ID') unique(hash_index) INT(11)" json:"uid"`
	TokenId     int    `xorm:"comment('代币ID') INT(11)" json:"token_id"`
	Ukey        string `xorm:"comment('联合key') unique(hash_index) VARCHAR(128)" json:"ukey"`
	Type        int    `xorm:"comment('流水类型1区块2委托') INT(11)" json:"type"`
	Opt         int    `xorm:"comment('操作方向1加2减') unique(hash_index) TINYINT(4)" json:"opt"`
	Num         int64  `xorm:"comment('数量') BIGINT(20)" json:"num"`
	Surplus     int64  `xorm:"comment('余额') BIGINT(20)" json:"surplus"`
	CreatedTime int64  `xorm:"comment('操作时间') BIGINT(20)" json:"created_time"`
}

type MoneyRecordWithToken struct {
	MoneyRecord `xorm:"extends"`
	TokenName   string `xorm:"token_name" json:"token_name"`
}

func (m *MoneyRecord) TableName() string {
	return "money_record"
}

func (m *MoneyRecord) GetMoneyList(page, rows int, uid []int64) (*ModelList, error) {
	engine := utils.Engine_token
	query := engine.Desc("id")
	query = query.In("uid", uid)
	tquery := *query
	count, err := tquery.Count(&MoneyRecord{})
	if err != nil {
		return nil, err
	}
	offset, modelList := m.Paging(page, rows, int(count))
	query.Limit(modelList.PageSize, offset)

	list := make([]MoneyRecord, 0)
	err = query.Find(&list)
	if err != nil {
		return nil, err
	}
	modelList.Items = list
	return modelList, nil
}

func (m *MoneyRecord) GetMoneyListForDateOrType(page, rows, ty, status int, date uint64, search string) (*ModelList, error) {
	engine := utils.Engine_token
	query := engine.Alias("uch").Desc("id")
	query = query.Join("LEFT", "g_common.user u ", "u.uid= uch.uid")
	query = query.Join("LEFT", "g_common.user_ex ex", "uch.uid=ex.uid")
	query = query.Where("uch.created_time between ? and ?", date, date+86400)

	if search != `` {
		temp := fmt.Sprintf(" concat(IFNULL(u.`uid`,''),IFNULL(u.`phone`,''),IFNULL(ex.`nick_name`,''),IFNULL(u.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}
	if ty != 0 {
		query = query.Where("uch.opt=?", ty)
	}
	if status != 0 {
		query = query.Where("u.status=?", status)
	}

	tempQuery := *query
	count, err := tempQuery.Count(&MoneyRecord{})
	if err != nil {
		return nil, err
	}
	offset, modelList := m.Paging(page, rows, int(count))
	list := make([]MoneyRecord, 0)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	for i, v := range list {
		list[i].NumTrue = m.Int64ToFloat64By8Bit(v.Num)
		list[i].SurplusTrue = m.Int64ToFloat64By8Bit(v.Surplus)
	}
	modelList.Items = list
	return modelList, nil

}

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
		dayBeginTime, _ := time.Parse(utils.LAYOUT_DATE_TIME, v.(string)+" 00:00:00")
		dayEndTime, _ := time.Parse(utils.LAYOUT_DATE_TIME, v.(string)+" 23:59:59")

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
