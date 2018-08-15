package models

import (
	"fmt"
	"admin/utils"
	)


type TotalProperty struct {
	BaseModel `xorm:"-"`
	Uid       int64
	Phone     string
	NickName  string
	Status    int
	Rt        int64 //用户注册时间
	Email     string
	TokenId   uint32
	Ucb       int64 //法币总资产
	Ucf       int64 //法币冻结
	Utb       int64 //币币总资产
	Utc       int64 //币币冻结
}

func (*TotalProperty) TableName() string {
	return "user"
}

func (w *TotalProperty) GetTotalProperty(page, rows, status int, date uint64, search string) (*ModelList, error) {
	engine := utils.Engine_common

	countSql := "SELECT COUNT(u.`uid`) num "
	contentSql := " SELECT `u`.`uid`, `u`.`email`, `u`.`phone`, `ex`.`nick_name`, `ex`.`register_time` `rt`, `uc`.`balance`  `ucb`, `uc`.`freeze` `ucf`, `ut`.`balance` `utb`, `ut`.`frozen` `utf` "
	mid := "FROM `user` AS `u` LEFT JOIN g_common.user_ex ex ON u.uid=ex.uid LEFT JOIN g_currency.user_currency uc ON u.uid=uc.uid LEFT JOIN g_token.user_token ut ON u.uid=ut.uid "
	sql := fmt.Sprintf(" WHERE (ex.register_time BETWEEN %d AND %d)", date, date+86400)
	if status != 0 {
		appendSql := fmt.Sprintf("and  u.status=%d ", status)
		sql += appendSql
	}
	if search != `` {
		appendSql := fmt.Sprintf(" and concat(IFNULL(u.`uid`,''),IFNULL(u.`phone`,''),IFNULL(ex.`nick_name`,''),IFNULL(u.`email`,'')) LIKE '%%%s%%'  ", search)
		sql += appendSql
	}
	Count := &struct {
		Num int
	}{}
	_, err := engine.SQL(countSql + mid + sql).Get(Count)
	if err != nil {
		return nil, err
	}
	offset, mList := w.Paging(page, rows, int(Count.Num))
	list := make([]TotalProperty, 0)
	limitSql := fmt.Sprintf("limit %d offset %d", mList.PageSize, offset)
	err = engine.SQL(contentSql + mid + sql + limitSql).Find(&list)
	if err != nil {
		return nil, err
	}
	mList.Items = list
	return mList, nil
}
