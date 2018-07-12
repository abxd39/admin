package models

import (
	"admin/utils"
)

type UserLoginLog struct {
	BaseModel    `xorm:"-"`
	Id           int    `xorm:"not null pk autoincr INT(10)"`
	Uid          int    `xorm:"not null comment('用户uid') INT(64)"`
	TerminalType int    `xorm:"comment('终端类型') TINYINT(4)"`
	TerminalName string `xorm:"not null comment('登录的终端名称') VARCHAR(100)"`
	LoginIp      string `xorm:"not null comment('登录IP') VARCHAR(15)"`
	LoginTime    string `xorm:"not null comment('登录日期时间') VARCHAR(10)"`
	Status       int    `xorm:"comment('用户账号状态') TINYINT(4)"`
}

type UserLogInLogGroup struct {
	UserLoginLog `xorm:"extends"`
	NickName     string `xorm:"not null default '' comment('用户昵称') VARCHAR(64)"`
	Phone        string `xorm:"comment('手机') unique VARCHAR(64)"`
	Email        string `xorm:"comment('邮箱') unique VARCHAR(128)"`
}

func (u *UserLogInLogGroup) TableName() string {
	return "user_login_log"
}

func (u *UserLogInLogGroup) GetUserLoginLogList(page, rows, terminal_type, status int, login_time string) ([]UserLogInLogGroup, int, int, error) {
	engine := utils.Engine_common
	limit := 0
	if page <= 1 {
		page = 1
	} else {
		limit = (page - 1) * rows
	}
	if rows <= 0 {
		rows = 100
	}

	query := engine.Desc("id")
	query.Join("INNER", "user", "user.uid=user_login_log.uid")
	query.Join("LEFT", "user_ex", "user_ex.uid=user.uid")
	if status != 0 {
		query = query.Where("status=?", status)
	}
	if terminal_type != 0 {
		query = query.Where("terminal_type=?", terminal_type)
	}
	if len(login_time) != 0 {
		query = query.Where("login_time=?", login_time)
	}
	tquery := query
	list := make([]UserLogInLogGroup, 0)
	err := query.Limit(rows, limit).Find(&list)
	if err != nil {
		return nil, 0, 0, nil
	}
	count, err := tquery.Count(&UserLogInLogGroup{})
	if err != nil {
		return nil, 0, 0, nil
	}
	total_page := int(count) / rows
	return list, total_page, int(count), nil
}
