package backstage

import (
	"admin/app/models"
	"admin/errors"
	"admin/utils"
)

// 管理员登录日志
type UserLoginLog struct {
	models.BaseModel `xorm:"-"`
	Id               int   `xorm:"not null pk autoincr INT(11)" json:"id"`
	Uid              int   `xorm:"not null comment('管理员ID') INT(11)" json:"uid"`
	LoginTime        int64 `xorm:"not null comment('登录时间') INT(11)" json:"login_time"`
	States           int   `xorm:"not null comment('1登录成功 2登录失败') TINYINT(1)" json:"states"`
}

// 添加登录日志
func (l *UserLoginLog) Add(log *UserLoginLog) (int, error) {
	engine := utils.Engine_backstage
	_, err := engine.Insert(log)
	if err != nil {
		return 0, errors.NewSys(err)
	}

	return log.Id, nil
}
