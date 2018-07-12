package backstage

import (
	"strconv"
	"strings"

	"admin/app/models"
	"admin/utils"
)

type Role struct {
	models.BaseModel `xorm:"-"`
	Id               int    `xorm:"not null pk autoincr INT(11)" json:"id"`
	Name             string `xorm:"not null default '' VARCHAR(36)" json:"name"`
	Desc             string `xorm:"not null VARCHAR(100)" json:"desc"`
	People           int    `xorm:"not null default 0 comment('人数') INT(6)" json:"people"`
	IsSuper          int    `xorm:"not null default 0 comment('是否超管 0 否 1是') TINYINT(1)" json:"is_super"`
}

// 用户组列表
func (r *Role) List(pageIndex, pageSize int) (modelList *models.ModelList, err error) {
	// 获取总数
	engine := utils.Engine_backstage
	query := engine.Desc("id")
	tempQuery := *query
	count, err := tempQuery.Count(&Role{})
	if err != nil {
		return nil, err
	}

	// 获取分页
	offset, modelList := r.Paging(pageIndex, pageSize, int(count))

	// 获取列表数据
	var list []Role
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return nil, err
	}
	modelList.Items = list

	return
}

// 新增用户组
func (r *Role) Add(name, desc, nodeIds string) (id int, err error) {
	// 用户组
	role := new(Role)
	role.Name = name
	role.Desc = desc

	// 关联的节点
	nodeIdArr := strings.Split(nodeIds, ",")

	// 开始写入，事务
	engine := utils.Engine_backstage
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return
	}

	// 1. 用户组
	_, err = session.Insert(role)
	if err != nil {
		session.Rollback()
		return
	}
	roleId := role.Id // 刚刚生成的id

	// 2. 用户组、节点关联
	for _, v := range nodeIdArr {
		nodeId, _ := strconv.Atoi(v)

		roleNode := &RoleNode{
			RoleId: roleId,
			NodeId: nodeId,
		}

		_, err = session.Insert(roleNode)
		if err != nil {
			session.Rollback()
			return
		}
	}

	err = session.Commit()
	if err != nil {
		return
	}

	return roleId, nil
}
