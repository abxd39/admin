package models

import (
	"admin/utils"
	"errors"
	"fmt"
)

//二级认证结构
type UserSecondaryCertification struct {
	BaseModel             `xorm:"-"`
	Id                    int    `xorm:"not null pk autoincr comment('自增id') INT(10)"`
	Uid                   int    `xorm:"not null comment('用户uid') INT(64)"`
	VerifyCount           int    `xorm:"not null comment('认证次数') TINYINT(4)"`
	VerifyTime            int    `xorm:"not null comment('认证时间') VARCHAR(100)"`
	VideoRecordingDigital string `xorm:"not null comment('视频录制的数字') VARCHAR(100)"`
	//VerifyStatus          int    `xorm:"not null comment('认证状态，1通过认证，2认证失败') TINYINT(4)"`
	PositivePath    string `xorm:"not null comment('身份证正面图片路径') VARCHAR(100)"`
	ReverseSidePath string `xorm:"not null comment('身份证反面图片路径') VARCHAR(100)"`
	VideoPath       string `xorm:"not null comment('视频路径') VARCHAR(100)"`
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
	query = query.Cols("user_secondary_certification.uid", "user_secondary_certification.verify_count", "user_secondary_certification.verify_time", "user.security_auth", "user_secondary_certification.video_recording_digital", "user.email", "user.phone", "user.status", "user_ex.nick_name")
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
	if us.SecurityAuth&utils.AUTH_TWO == utils.AUTH_TWO{
		us.TwoVerifyMark = 1
	}
	return us, nil
}

//二级认证列表
func (u *UserSecondaryCertification) GetSecondaryCertificationList(page, rows, verify_status, user_status, time int, search string) (*ModelList, error) {

	engine := utils.Engine_common
	query := engine.Desc("id")
	query = query.Join("INNER", "user", "user.uid=user_secondary_certification.uid")
	query = query.Join("LEFT", "user_ex", "user_ex.uid = user.uid")
	query = query.Cols("user_secondary_certification.uid", "user_secondary_certification.verify_count", "user_secondary_certification.verify_time", "user.security_auth", "user_secondary_certification.video_recording_digital", "user.email", "user.phone", "user.status", "user_ex.nick_name")
	if verify_status ==-1 {//刷选未认证
		query = query.Where("user.security_auth & ? !=?", utils.AUTH_TWO,utils.AUTH_TWO)
	}
	if verify_status == utils.AUTH_TWO{
		query = query.Where("user.security_auth & ? =?", utils.AUTH_TWO,utils.AUTH_TWO)
	}
	if user_status != 0 {
		query = query.Where("status=?", user_status)
	}
	if time != 0 {
		query = query.Where("verify_time  BETWEEN ? AND ? ", time, time+86400)
	}
	if len(search) != 0 {
		temp := fmt.Sprintf(" concat(IFNULL(`user`.`uid`,''),IFNULL(`user`.`phone`,''),IFNULL(`user_ex`.`nick_name`,''),IFNULL(`user`.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}
	tempquery := *query
	count, err := tempquery.Count(&UserSecondaryCertificationGroup{})
	if err != nil {
		return nil, err
	}
	offset, modellist := u.Paging(page, rows, int(count))
	list := make([]UserSecondaryCertificationGroup, 0)

	err = query.Limit(modellist.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	for index, _ := range list {
		if list[index].SecurityAuth&utils.AUTH_TWO == utils.AUTH_TWO {
			list[index].TwoVerifyMark = 1
		}
	}
	modellist.Items = list
	return modellist, nil
}
