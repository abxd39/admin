package backstage

import (
	"strconv"
	"strings"

	"admin/app/models"
	"admin/errors"
	"admin/utils"
)

// 用户组
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

// 用户组详情
func (r *Role) Get(id int) (role *Role, err error) {
	engine := utils.Engine_backstage
	role = new(Role)
	has, err := engine.ID(id).Get(role)
	if err != nil {
		return nil, errors.NewSys(err)
	}
	if !has {
		return nil, errors.NewNormal("用户组不存在或已被删除")
	}

	return
}

// 用户组绑定的节点ID
func (r *Role) GetBindNodeIds(roleId int) (nodeIds []int, err error) {
	engine := utils.Engine_backstage
	err = engine.Table(new(RoleNode)).Where("role_id=?", roleId).Cols("node_id").Find(&nodeIds)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	return
}

// 新增用户组
func (r *Role) Add(role *Role, nodeIds string) (id int, err error) {
	// 判断用户组名称是否已存在
	engine := utils.Engine_backstage
	checkRole := new(Role)
	has, err := engine.Where("name=?", role.Name).Get(checkRole)
	if err != nil {
		return 0, errors.NewSys(err)
	}
	if has && checkRole.Name == role.Name {
		return 0, errors.NewNormal("名称已存在")
	}

	// 开始写入，事务
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return 0, errors.NewSys(err)
	}

	// 1. 新增用户组
	_, err = session.Insert(role)
	if err != nil {
		session.Rollback()
		return 0, errors.NewSys(err)
	}
	roleId := role.Id // 刚刚生成的id

	// 2. 新增用户组、节点关联
	nodeIdArr := strings.Split(nodeIds, ",") // 逗号分隔
	for _, v := range nodeIdArr {
		nodeId, err := strconv.Atoi(v)
		if err != nil || nodeId <= 0 {
			session.Rollback()
			return 0, errors.NewNormal("参数node_ids格式错误")
		}

		roleNodeMD := &RoleNode{
			RoleId: roleId,
			NodeId: nodeId,
		}

		_, err = session.Insert(roleNodeMD)
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

// 更新用户组
func (r *Role) Update(id int, name, desc, nodeIds string) error {
	// 验证用户组是否存在
	engine := utils.Engine_backstage
	has, err := engine.Id(id).Get(new(Role))
	if err != nil {
		return errors.NewSys(err)
	}
	if !has {
		return errors.NewNormal("用户组不存在或已删除")
	}

	// 判断用户组名称是否已存在
	has, err = engine.Where("name=?", name).And("id!=?", id).Get(new(Role))
	if err != nil {
		return errors.NewSys(err)
	}
	if has {
		return errors.NewNormal("名称已存在")
	}

	// 开始更新，事务
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return errors.NewSys(err)
	}

	// 1. 更新用户组
	roleMD := &Role{
		Name: name,
		Desc: desc,
	}
	_, err = session.ID(id).Update(roleMD)
	if err != nil {
		session.Rollback()
		return errors.NewSys(err)
	}

	// 2. 更新用户组、节点关联
	// 2.1 删除之前的关联
	_, err = session.Where("role_id=?", id).Delete(new(RoleNode))
	if err != nil {
		session.Rollback()
		return errors.NewSys(err)
	}

	// 2.2 新增关联
	nodeIdArr := strings.Split(nodeIds, ",") // 逗号分隔
	for _, v := range nodeIdArr {
		nodeId, err := strconv.Atoi(v)
		if err != nil || nodeId <= 0 {
			session.Rollback()
			return errors.NewNormal("参数node_ids格式错误")
		}

		roleNodeMD := &RoleNode{
			RoleId: id,
			NodeId: nodeId,
		}

		_, err = session.Insert(roleNodeMD)
		if err != nil {
			session.Rollback()
			return errors.NewSys(err)
		}
	}

	err = session.Commit()
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
}

// 删除用户组
func (r *Role) Delete(id int) error {
	// 验证用户组是否存在
	engine := utils.Engine_backstage
	has, err := engine.Id(id).Get(new(Role))
	if err != nil {
		return errors.NewSys(err)
	}
	if !has {
		return errors.NewNormal("用户组不存在或已删除")
	}

	// 开始删除，事务
	session := engine.NewSession()
	defer session.Close()

	// 1. 删除用户组
	_, err = session.ID(id).Delete(new(Role))
	if err != nil {
		session.Rollback()
		return errors.NewSys(err)
	}

	// 2. 删除用户组、节点关联
	_, err = session.Where("role_id=?", id).Delete(new(RoleNode))
	if err != nil {
		session.Rollback()
		return errors.NewSys(err)
	}

	err = session.Commit()
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
}
