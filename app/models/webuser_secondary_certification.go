package models

import (
	"admin/utils"
)

//二级认证结构
type UserSecondaryCertification struct {
	BaseModel             `xorm:"-"`
	Id                    int    `xorm:"not null pk autoincr comment('自增id') INT(10)"`
	Uid                   int    `xorm:"not null comment('用户uid') INT(64)"`
	VerifyCount           int    `xorm:"not null comment('认证次数') TINYINT(4)"`
	VerifyTime            string `xorm:"not null comment('认证时间') VARCHAR(100)"`
	VideoRecordingDigital string `xorm:"not null comment('视频录制的数字') VARCHAR(100)"`
	VerifyStatus          int    `xorm:"not null comment('认证状态，1通过认证，2认证失败') TINYINT(4)"`
	PositivePath          string `xorm:"not null comment('身份证正面图片路径') VARCHAR(100)"`
	ReverseSidePath       string `xorm:"not null comment('身份证反面图片路径') VARCHAR(100)"`
	VideoPath             string `xorm:"not null comment('视频路径') VARCHAR(100)"`
}

type UserSecondaryCertificationGroup struct {
	UserSecondaryCertification `xorm:"extends"`
	NickName                   string `xorm:"not null default '' comment('用户昵称') VARCHAR(64)"`
	Phone                      string `xorm:"comment('手机') unique VARCHAR(64)"`
	Email                      string `xorm:"comment('邮箱') unique VARCHAR(128)"`
	Account                    string `xorm:"comment('账号') unique VARCHAR(64)"`
	Status                     int    `xorm:"comment('用户账号状态') TINYINT(4)"`
}

func (u *UserSecondaryCertificationGroup) TableName() string {
	return "user_secondary_certification"
}

//二级认证详情
func (u *UserSecondaryCertificationGroup) GetSecondaryCertificationOfUid(uid int) (UserSecondaryCertificationGroup, error) {
	return UserSecondaryCertificationGroup{}, nil
}

//二级认证列表
func (u *UserSecondaryCertification) GetSecondaryCertificationList(page, rows, verify_status, user_status int, certification_time string) ([]UserSecondaryCertificationGroup, int, int, error) {
	limit := 0
	if rows <= 0 {
		rows = 50
	}
	if page <= 1 {
		page = 1
	} else {
		limit = (page - 1) * rows
	}
	engine := utils.Engine_common
	query := engine.Desc("id")
	query = query.Cols("user_secondary_certification.uid", "verify_count", "verify_time", "verify_status", "video_recording_digital", "email", "phone", "status", "nick_name")
	query = query.Join("INNER", "user", "user.uid=user_secondary_certification.uid")
	query = query.Join("LEFT", "user_ex", "user_ex.uid=user.uid")
	if verify_status != 0 {
		query = query.Where("verify_status=?", verify_status)
	}
	if user_status != 0 {
		query = query.Where("status=?", user_status)
	}
	if len(certification_time) != 0 {
		query = query.Where("verify_time=?", certification_time)
	}
	tempquery := *query
	list := make([]UserSecondaryCertificationGroup, 0)
	err := query.Limit(rows, limit).Find(&list)
	if err != nil {
		return nil, 0, 0, err
	}
	count, err := tempquery.Count(&UserSecondaryCertificationGroup{})
	if err != nil {
		return nil, 0, 0, err
	}
	total_page := int(count) / rows
	return list, total_page, int(count), nil
}
