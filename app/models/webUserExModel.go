package models

import (
	"admin/utils"
	"errors"
	"fmt"
)

type UserEx struct {
	BaseModel     `xorm:"-"`
	Uid           int64  `xorm:"not null pk comment(' 用户ID') BIGINT(11)"`
	NickName      string `xorm:"not null default '' comment('用户昵称') VARCHAR(64)"`
	HeadSculpture string `xorm:"not null default '' comment('头像图片路径') VARCHAR(100)"`
	RegisterTime  int64  `xorm:"comment('注册时间') BIGINT(20)"`
	AffirmTime    int64  `xorm:"comment('实名认证时间') BIGINT(20)"`
	InviteCode    string `xorm:"comment('邀请码') VARCHAR(64)"`
	RealName      string `xorm:"comment(' 真名') VARCHAR(32)"`
	IdentifyCard  string `xorm:"comment('身份证号') VARCHAR(64)"`
	InviteId      int64  `xorm:"comment('邀请者id') BIGINT(11)"`
	Invites       int    `xorm:"default 0 comment('邀请人数') INT(11)"`
	AffirmCount   int    `xorm:"default 0 comment('实名认证的次数') TINYINT(4)"`
	ChannelName   string `xorm:"not null default '' comment('邀请的渠道名称') VARCHAR(100)" json:"channel_name"`
}

type FirstDetail struct {
	UserEx       `xorm:"extends"`
	Account      string `xorm:"comment('账号') unique VARCHAR(64)"`
	SecurityAuth int    `xorm:"comment('认证状态1110') TINYINT(8)"`
	VerifyMark   int    //一级实名认证状态
}

//func (f *FirstDetail) TableName() string {
//	return "user_ex"
//}

type InviteGroup struct {
	UserEx      `xorm:"extends"`
	Account     string `xorm:"comment('账号') unique VARCHAR(64)"`
	Email       string `xorm:"comment('邮箱') unique VARCHAR(128)"`
	Phone       string `xorm:"comment('手机') unique VARCHAR(64)"`
	Status      int    `xorm:"default 0 comment('用户状态，0正常，1冻结') INT(11)"`
	InviteCount int
}

//func (w *UserEx) TableName() string {
//	return "user_ex"
//}

//邀请人统计表—账号：18888888888
func (w *UserEx) GetInviteInfoList(uid, page, rows int, date uint64, name, account string) (*ModelList, error) {
	engine := utils.Engine_common
	query := engine.Desc("user_ex.uid")
	query = query.Join("LEFT", "user", "user.uid=user_ex.uid")
	query = query.Cols("user_ex.uid", "user.account", "user_ex.register_time", "user_ex.channel_name")
	query = query.Where("`user_ex`.`invite_id`=?", uid)
	if name != `` {
		temp := fmt.Sprintf("channer_name=%s", name)
		query = query.Where(temp)
	}
	if account != `` {
		temp := fmt.Sprintf("user.account=%s", account)
		query = query.Where(temp)
	}
	fmt.Println("刷选时间=", date)
	if date != 0 {
		query = query.Where("`user_ex`.`register_time` BETWEEN ? AND ? ", date, date+86400)
	}
	tempQuery := *query
	count, err := tempQuery.Count(&UserEx{})
	if err != nil {
		return nil, err
	}
	offset, modelList := w.Paging(page, rows, int(count))
	list := make([]InviteGroup, 0)
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	modelList.Items = list
	return modelList, nil
}

//p2-5好友邀请 ___ 有邀请用户注册的列表

func (w *UserEx) GetInViteList(page, rows int, search string) (*ModelList, error) {
	engine := utils.Engine_common
	countSql := "SELECT COUNT(*) cnt FROM (%s) newtable"
	sql := "SELECT u.uid,u.`email`,u.`phone`,ue.`real_name`,ue.nick_name,u.account,ue.`invite_id`,u.`status`,tmp.cnt invite_count FROM (SELECT invite_id,COUNT(invite_id) cnt FROM user_ex GROUP BY invite_id) tmp JOIN user_ex ue ON ue.uid=tmp.invite_id JOIN `user` u ON u.uid=tmp.invite_id WHERE tmp.invite_id!=0"
	limitSql := " LIMIT %d OFFSET %d "

	if search != `` {
		temp := fmt.Sprintf(" AND CONCAT(IFNULL(`ue`.`uid`,''),IFNULL(`u`.`phone`,''),IFNULL(`ue`.`nick_name`,''),IFNULL(`u`.`email`,'')) LIKE '%%%s%%' ", search)
		sql += temp
	}
	//query := engine.SQL(contentSql + sql)

	type test struct {
		Cnt int
	}
	var temp test
	_, err := engine.SQL(fmt.Sprintf(countSql, sql)).Get(&temp)

	fmt.Println("cnt=", temp.Cnt)
	fmt.Println("page=", page, "rows=", rows, "cnt=", temp.Cnt)
	offset, modelList := w.Paging(page, rows, temp.Cnt)
	list := make([]InviteGroup, 0)
	limitSql = fmt.Sprintf(limitSql, modelList.PageSize, offset)
	fmt.Println("sql2=", sql+limitSql)
	err = engine.SQL(sql + limitSql).Find(&list)
	if err != nil {
		return nil, err
	}

	//fmt.Println("resultList=",list)
	modelList.Items = list
	return modelList, nil

}

//单个用户的认证详情
func (w *UserEx) GetFirstDetail(uid int) (*FirstDetail, error) {
	engine := utils.Engine_common
	query := engine.Desc("user_ex.uid")
	query = query.Join("INNER", "user", "user.uid=user_ex.uid")
	query = query.Where("user_ex.uid=?", uid)
	//query = query.Cols("user_ex.register_time", "user_ex.uid", "user_ex.real_name", "user_ex.identify_card", "user_ex.affirm_time", "user_ex.affirm_count", "user.account", "user_ex.nick_name", "user.security_auth")
	temp := *query
	has, err := temp.Exist(&FirstDetail{})
	if err != nil {
		fmt.Println("bukexue ")
		return nil, err
	}
	if !has {
		return nil, errors.New("用户不存在！！")
	}
	u := new(FirstDetail)
	_, err = query.Get(u)
	if err != nil {
		return nil, err
	}
	if u.SecurityAuth&utils.AUTH_FIRST == utils.AUTH_FIRST {
		u.VerifyMark = 1
	}
	return u, nil
}
