package backstage

import (
	"admin/utils"

	"admin/app/models"
)

type Role struct {
	models.BaseModel `xorm:"-"`
	Id               int    `xorm:"not null pk autoincr INT(11)"`
	Name             string `xorm:"not null default '' VARCHAR(36)"`
	Desc             string `xorm:"not null VARCHAR(100)"`
	People           int    `xorm:"not null default 0 comment('人数') INT(6)"`
	IsSuper          int    `xorm:"not null default 0 comment('是否超管 0 否 1是') TINYINT(1)"`
}

// 用户组列表
func (r *Role) GetRoleList(pageIndex, pageSize int) (list []Role, totalPage, total int, err error) {
	// 获取总数
	engine := utils.Engine_backstage
	query := engine.Desc("id")
	tempQuery := *query
	count, err := tempQuery.Count(&Role{})
	if err != nil {
		return nil, 0, 0, err
	}
	total = int(count)

	// 获取分页
	offset, totalPage := r.Paging(pageIndex, pageSize, total)

	// 获取列表数据
	err = query.Limit(pageSize, offset).Find(&list)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return nil, 0, 0, err
	}

	return list, totalPage, total, nil
}
