package models

import (
	"admin/errors"
	"admin/utils"
	"fmt"
	"time"
)

type MoneyRecord struct {
	BaseModel   `xorm:"-"`
	Id          int64  `xorm:"pk autoincr BIGINT(20)" json:"id"`
	Uid         int    `xorm:"comment('用户ID') unique(hash_index) INT(11)" json:"uid"`
	TokenId     int    `xorm:"comment('代币ID') INT(11)" json:"token_id"`
	Ukey        string `xorm:"comment('联合key') unique(hash_index) VARCHAR(128)" json:"ukey"`
	Type        int    `xorm:"comment('流水类型1区块2委托') INT(11)" json:"type"`
	Opt         int    `xorm:"comment('操作方向1加2减') unique(hash_index) TINYINT(4)" json:"opt"`
	Num         int64  `xorm:"comment('数量') BIGINT(20)" json:"num"`
	Balance     int64  `xorm:"comment('余额') BIGINT(20)" json:"surplus"`
	CreatedTime int64  `xorm:"comment('操作时间') BIGINT(20)" json:"created_time"`
	Comment     string `xorm:"comment('备注') varchar(255)" json:"comment"`
}

func (m *MoneyRecord) TableName() string {
	return "money_record"
}

type MoneyRecordGroup struct {
	MoneyRecord `xorm:"extends"`
	UserInfo    `xorm:"extends"`
}

func (m *MoneyRecordGroup) TableName() string {
	return "money_record"
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
		CreatedTime: time.Now().UnixNano(),
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
	if cnt ==0{
		sess.Rollback()
		return errors.New("充值失败!!!")
	}
	sess.Commit()
	return nil
}

func (m *MoneyRecordGroup) GetMoneyList(page, rows int, uid []int64) (*ModelList, error) {
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

func (m *MoneyRecordGroup) GetMoneyListForDateOrType(page, rows, ty, status int, tid int, search string) (*ModelList, error) {
	engine := utils.Engine_token
	query := engine.Alias("uch").Desc("id")
	query = query.Join("LEFT", "g_common.user u ", "u.uid= uch.uid")
	query = query.Join("LEFT", "g_common.user_ex ex", "uch.uid=ex.uid")
	query = query.Where("uch.token_id=?", tid)

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
