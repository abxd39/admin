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
	Id               int    `xorm:"id pk autoincr" json:"id"`
	Uid              int    `xorm:"uid" json:"uid"`
	NickName         string `xorm:"nick_name" json:"nick_name"`
	States           int    `xorm:"states" json:"states"`
	LoginIp          string `xorm:"login_ip" json:"login_ip"`
	LoginTime        int64  `xorm:"login_time" json:"login_time"`
}

// 表名
func (*UserLoginLog) TableName() string {
	return "user_login_log"
}

// 登录日志列表
func (l *UserLoginLog) List(pageIndex, pageSize int, filter map[string]string) (modelList *models.ModelList, list []UserLoginLog, err error) {
	// 获取总数
	engine := utils.Engine_backstage
	query := engine.Desc("id")

	// 筛选
	query.Where("1=1")
	if v, ok := filter["login_date_start"]; ok {
		query.And("login_time>=?", utils.Date2Unix(v, utils.LAYOUT_DATE))
	}
	if v, ok := filter["login_date_end"]; ok {
		query.And("login_time<=?", utils.Date2Unix(fmt.Sprint(v, " 23:59:59"), utils.LAYOUT_DATE_TIME))
	}

	tempQuery := *query
	count, err := tempQuery.Count(&UserLoginLog{})
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}

	// 获取分页
	offset, modelList := l.Paging(pageIndex, pageSize, int(count))

	// 获取列表数据
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, nil, errors.NewSys(err)
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
