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
	Id               int    `xorm:"id pk autoincr" json:"id"`
	Name             string `xorm:"name" json:"name"`
	Desc             string `xorm:"desc" json:"desc"`
	People           int    `xorm:"people" json:"people"`
}

// 表名
func (*Role) TableName() string {
	return "role"
}

// 用户组列表
func (r *Role) List(pageIndex, pageSize int) (modelList *models.ModelList, list []Role, err error) {
	// 获取总数
	engine := utils.Engine_backstage
	query := engine.Desc("id")
	tempQuery := *query
	count, err := tempQuery.Count(&Role{})
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}

	// 获取分页
	offset, modelList := r.Paging(pageIndex, pageSize, int(count))

	// 获取列表数据
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}
	modelList.Items = list

	return
}

// 用户组列表，不分页
func (r *Role) ListAll() (*models.ModelList, []Role, error) {
	engine := utils.Engine_backstage
	query := engine.Desc("id")

	var list []Role
	err := query.Find(&list)
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}

	return r.NoPaging(len(list), list), list, nil
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
	has, err := engine.Where("name=?", role.Name).Exist(new(Role))
	if err != nil {
		return 0, errors.NewSys(err)
	}
	if has {
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
	if nodeIds != "" {
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
	}

	err = session.Commit()
	if err != nil {
		return 0, errors.NewSys(err)
	}

	return roleId, nil
}

// 更新用户组
func (r *Role) Update(id int, params map[string]interface{}) error {
	// 验证用户组是否存在
	engine := utils.Engine_backstage
	has, err := engine.Id(id).Exist(new(Role))
	if err != nil {
		return errors.NewSys(err)
	}
	if !has {
		return errors.NewNormal("用户组不存在或已删除")
	}

	// 验证参数、设置更新的字段
	roleData := make(map[string]interface{})
	if v, ok := params["name"]; ok {
		// 判断用户组名称是否已存在
		has, err = engine.Where("name=?", v).And("id!=?", id).Exist(new(Role))
		if err != nil {
			return errors.NewSys(err)
		}
		if has {
			return errors.NewNormal("名称已存在")
		}

		roleData["name"] = v
	}
	if v, ok := params["desc"]; ok {
		roleData["desc"] = v
	}

	// 开始更新，事务
	session := engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return errors.NewSys(err)
	}

	// 1. 更新用户组
	_, err = session.Table(r).ID(id).Update(roleData)
	if err != nil {
		session.Rollback()
		return errors.NewSys(err)
	}

	// 2. 更新用户组、节点关联
	if v, ok := params["node_ids"]; ok {
		nodeIds := v.(string)
		if nodeIds != "" {
			// 2.1 删除之前的关联
			_, err = session.Where("role_id=?", id).Delete(new(RoleNode))
			if err != nil {
				session.Rollback()
				return errors.NewSys(err)
			}

			// 2.2 更新关联
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
	has, err := engine.Id(id).Exist(new(Role))
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
