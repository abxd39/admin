package backstage

import (
	"admin/app/models"
	"admin/errors"
	"admin/utils"
	"fmt"
)

// 管理员登录日志
type UserLoginLog struct {
	models.BaseModel `xorm:"-"`
	Id               int    `xorm:"not null pk autoincr INT(11)" json:"id"`
	Uid              int    `xorm:"not null comment('管理员ID') INT(11)" json:"uid"`
	NickName         string `xorm:"not null comment('管理员昵称') VARCHAR(60)" json:"nick_name"`
	States           int    `xorm:"not null comment('1登录成功 2登录失败') TINYINT(1)" json:"states"`
	LoginIp          string `xorm:"not null comment('登录IP') VARCHAR(30)" json:"login_ip"`
	LoginTime        int64  `xorm:"not null comment('登录时间') INT(11)" json:"login_time"`
}

// 表名
func (*UserLoginLog) TableName() string {
	return "user_login_log"
}

// 登录日志列表
func (l *UserLoginLog) List(pageIndex, pageSize int, filter map[string]string) (modelList *models.ModelList, err error) {
	// 获取总数
	engine := utils.Engine_backstage
	query := engine.Desc("id")

	// 筛选
	query.Where("1=1")
	if v, ok := filter["login_date"]; ok {
		timeStart := utils.Date2Unix(fmt.Sprint(v, " 00:00:00"))
		timeEnd := utils.Date2Unix(fmt.Sprint(v, " 23:59:59"))

		query.And("login_time>=?", timeStart).And("login_time<=?", timeEnd)
	}

	tempQuery := *query
	count, err := tempQuery.Count(&UserLoginLog{})
	if err != nil {
		return nil, errors.NewSys(err)
	}

	// 获取分页
	offset, modelList := l.Paging(pageIndex, pageSize, int(count))

	// 获取列表数据
	var list []UserLoginLog
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, errors.NewSys(err)
	}
	modelList.Items = list

	return
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
