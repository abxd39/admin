package models

import (
	"admin/utils"
	"fmt"
)

type UserLoginLog struct {
	BaseModel    `xorm:"-"`
	Id           int    `xorm:"not null pk autoincr INT(10)"`
	Uid          int    `xorm:"not null comment('用户uid') INT(64)"`
	TerminalType int    `xorm:"comment('终端类型') TINYINT(4)"`
	TerminalName string `xorm:"not null comment('登录的终端名称') VARCHAR(100)"`
	LoginIp      string `xorm:"not null comment('登录IP') VARCHAR(15)"`
	LoginTime    int    `xorm:"not null comment('登录日期时间') VARCHAR(10)"`
}

type UserLogInLogGroup struct {
	//UserLoginLog `xorm:"extends"`
	NickName     string `xorm:"not null default '' comment('用户昵称') VARCHAR(64)"`
	Phone        string `xorm:"comment('手机') unique VARCHAR(64)"`
	Email        string `xorm:"comment('邮箱') unique VARCHAR(128)"`
	Status       int    `xorm:"default 0 comment('用户状态，1正常，2冻结') INT(11)"`
}

func (u *UserLogInLogGroup) TableName() string {
	return "user_login_log"
}


//获取用户登录日志
func (u *UserLoginLog) GetUserLoginLogList(page, rows, terminal_type, status int, login_time uint64, search string) (*ModelList, error) {
	engine := utils.Engine_common

	query := engine.Desc("id")
	query.Join("INNER", "user", "user.uid=user_login_log.uid")
	query.Join("LEFT", "user_ex", "user_ex.uid=user.uid")
	if status != 0 {
		query = query.Where("status=?", status)
	}
	if terminal_type != 0 {
		query = query.Where("terminal_type=?", terminal_type)
	}
	if login_time != 0 {
		query = query.Where("`user_login_log`.`login_time`  BETWEEN ? AND ? ", login_time, login_time+86400) //86400 为一天的秒数
	}
	if len(search) != 0 {
		temp := fmt.Sprintf(" concat(IFNULL(`user`.`uid`,''),IFNULL(`user`.`phone`,''),IFNULL(`user_ex`.`nick_name`,''),IFNULL(`user`.`email`,'')) LIKE '%%%s%%'  ", search)
		query = query.Where(temp)
	}
	tquery := *query
	count, err := tquery.Count(&UserLogInLogGroup{})
	if err != nil {
		return nil, err
	}
	offset, modelist := u.Paging(page, rows, int(count))

	list := make([]UserLogInLogGroup, 0)
	err = query.Limit(modelist.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}
	modelist.Items = list
	return modelist, nil
}
