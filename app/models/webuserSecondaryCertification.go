package models

import (
	"admin/utils"
	"errors"
	"fmt"
	"time"
)

//二级认证结构

type UserSecondaryCertification struct {
	BaseModel             `xorm:"-"`
	Id                    int    `xorm:"not null pk autoincr comment('自增id') INT(10)"`
	Uid                   int    `xorm:"not null comment('用户uid') INT(64)"`
	VerifyCount           int    `xorm:"not null default 0 comment('认证次数') TINYINT(4)"`
	VerifyTime            int    `xorm:"not null comment('认证时间戳') INT(11)"`
	VideoRecordingDigital string `xorm:"not null comment('视频录制的数字') VARCHAR(100)"`
	PositivePath          string `xorm:"not null comment('身份证正面图片路径') VARCHAR(100)"`
	ReverseSidePath       string `xorm:"not null comment('身份证反面图片路径') VARCHAR(100)"`
	VideoPath             string `xorm:"not null comment('视频路径') VARCHAR(100)"`
	InHandPicturePath     string `xorm:"not null comment('手持身份证照片路径') VARCHAR(100)"`
}

type UserSecondaryCertificationGroup struct {
	UserSecondaryCertification `xorm:"extends"`
	NickName                   string `xorm:"not null default '' comment('用户昵称') VARCHAR(64)"`
	SecurityAuth               int    `xorm:"comment('认证状态1110') TINYINT(8)"`
	TwoVerifyMark              int
	Phone                      string `xorm:"comment('手机') unique VARCHAR(64)"`
	Email                      string `xorm:"comment('邮箱') unique VARCHAR(128)"`
	Account                    string `xorm:"comment('账号') unique VARCHAR(64)"`
	Status                     int    `xorm:"comment('用户账号状态') TINYINT(4)"`
}

func (u *UserSecondaryCertificationGroup) TableName() string {
	return "user_secondary_certification"
}

//二级认证详情
func (u *UserSecondaryCertificationGroup) GetSecondaryCertificationOfUid(uid int) (*UserSecondaryCertificationGroup, error) {
	engine := utils.Engine_common
	//
	query := engine.Desc("id")
	query = query.Join("INNER", "user", "user.uid=user_secondary_certification.uid")
	query = query.Join("LEFT", "user_ex", "user_ex.uid = user_secondary_certification.uid")
	//query = query.Cols("user_secondary_certification.uid", "user_secondary_certification.verify_count", "user_secondary_certification.verify_time", "user.security_auth", "user_secondary_certification.video_recording_digital", "user.email", "user.phone", "user.status", "user_ex.nick_name")
	query = query.Where("user_secondary_certification.uid=?", uid)
	tempQuery := *query
	has, err := tempQuery.Exist(&UserSecondaryCertification{})
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("没有二次认证信息！！")
	}
	us := new(UserSecondaryCertificationGroup)
	_, err = query.Get(us)
	if err != nil {
		return nil, err
	}
	if us.SecurityAuth&utils.AUTH_TWO == utils.AUTH_TWO {
		us.TwoVerifyMark = 1
	}

	return us, nil
}

//二级认证列表
func (u *UserSecondaryCertification) GetSecondaryCertificationList(page, rows, verify_status, user_status, tm int, search string) (*ModelList, error) {

	engine := utils.Engine_common
	//query := engine.Desc("id")
	//query = query.Join("INNER", "user", "user.uid=user_secondary_certification.uid")
	//query = query.Join("LEFT", "user_ex", "user_ex.uid = user.uid AND user.set_tarde_mark & 4=4")
	//query = query.Cols("user_secondary_certification.uid", "user_secondary_certification.verify_count", "user_secondary_certification.verify_time", "user.security_auth", "user_secondary_certification.video_recording_digital", "user.email", "user.phone", "user.status", "user_ex.nick_name")
	//if verify_status ==-1 {//刷选未认证
	//	query = query.Where("user.security_auth & ? !=?", utils.AUTH_TWO,utils.AUTH_TWO)
	//}
	//if verify_status == utils.AUTH_TWO{
	//	query = query.Where("user.security_auth & ? =?", utils.AUTH_TWO,utils.AUTH_TWO)
	//}
	//if user_status != 0 {
	//	query = query.Where("status=?", user_status)
	//}
	//if time != 0 {
	//	query = query.Where("verify_time  BETWEEN ? AND ? ", time, time+86400)
	//}
	//if len(search) != 0 {
	//	temp := fmt.Sprintf(" concat(IFNULL(`user`.`uid`,''),IFNULL(`user`.`phone`,''),IFNULL(`user_ex`.`nick_name`,''),IFNULL(`user`.`email`,'')) LIKE '%%%s%%'  ", search)
	//	query = query.Where(temp)
	//}

	sql := "SELECT t.phone,t.email,t.status,t.uid,t.set_tarde_mark,t.security_auth,us.verify_time,us.verify_count,us.video_recording_digital,ex.nick_name "
	countSql := "SELECT COUNT(*) num "

	Value := "FROM  g_common.`user` t JOIN g_common.`user_secondary_certification` us ON us.uid =t.uid  JOIN g_common.`user_ex` ex ON us.uid= ex.uid "
	var condition string
	if verify_status == -1 {
		temp := "WHERE ( t.set_tarde_mark &4=4 OR  t.set_tarde_mark&8=8 ) "
		condition = temp
	} else {
		temp := "WHERE (t.security_auth&4=4 OR t.set_tarde_mark &4=4 OR  t.set_tarde_mark&8=8 ) "
		condition = temp
	}
	if user_status != 0 {
		temp := fmt.Sprintf("AND t.status=%d ", user_status)
		condition += temp
	}
	if tm != 0 {
		temp := fmt.Sprintf(" AND us.verify_time BETWEEN %d AND %d ", tm, tm+86400)
		condition += temp
	}
	if search != `` {
		temp := fmt.Sprintf(" and concat(IFNULL(t.`uid`,''),IFNULL(t.`phone`,''),IFNULL(ex.`nick_name`,''),IFNULL(t.`email`,'')) LIKE '%%%s%%'  ", search)
		condition +=temp
	}

	count := &struct {
		Num int
	}{}

	str:=countSql + Value + condition
	fmt.Println(str)
	_, err := engine.SQL(str).Get(count)
	if err != nil {
		return nil, err
	}
	fmt.Println("-------------->num=", count.Num)
	offset, modellist := u.Paging(page, rows, int(count.Num))
	type UserCer struct {
		NickName              string
		SecurityAuth          int
		TwoVerifyMark         int
		Phone                 string
		Email                 string
		Account               string
		Status                int
		Uid                   int
		VerifyCount           int
		VerifyTime            int64
		VideoRecordingDigital string
		VerifyTimeStr string `xorm:"-" json:"verify_time_str"`
	}
	list := make([]UserCer, 0)
	limitSql := fmt.Sprintf(" ORDER BY us.`uid`  DESC LIMIT %d OFFSET %d", modellist.PageSize, offset)
	str=sql + Value + condition + limitSql
	fmt.Println(str)
	err = engine.SQL(str).Find(&list)
	if err != nil {
		return nil, err
	}
	for index, _ := range list {
		if list[index].SecurityAuth&utils.AUTH_TWO == utils.AUTH_TWO {
			fmt.Println("---> securityauth=", list[index].SecurityAuth)
			list[index].TwoVerifyMark = 1
		}
	}
	for i,v:=range  list{
		list[i].VerifyTimeStr = time.Unix(v.VerifyTime,0).Format("2006-01-02 15:04:05")
	}
	modellist.Items = list
	return modellist, nil
}
