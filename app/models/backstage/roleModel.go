package backstage

import (
	"strconv"
	"strings"

	"admin/app/models"
	"admin/errors"
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
		return nil, errors.NewSys(err)
	}

	// 获取分页
	offset, modelList := r.Paging(pageIndex, pageSize, int(count))

	// 获取列表数据
	var list []Role
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, errors.NewSys(err)
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

	// 判断用户组名称是否已存在
	engine := utils.Engine_backstage
	checkRole := new(Role)
	has, err := engine.Where("name=?", name).Get(checkRole)
	if err != nil {
		return 0, errors.NewSys(err)
	}
	if has && checkRole.Name == name {
		return 0, errors.NewNormal("名称已存在")
	}

	// 开始写入，事务
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return 0, errors.NewSys(err)
	}

	// 1. 用户组
	_, err = session.Insert(role)
	if err != nil {
		session.Rollback()
		return 0, errors.NewSys(err)
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
			return 0, errors.NewSys(err)
		}
	}

	err = session.Commit()
	if err != nil {
		return 0, errors.NewSys(err)
	}

	return roleId, nil
}
