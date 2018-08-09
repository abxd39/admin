package backstage

import (
	"admin/app/models"
	"admin/errors"
	"admin/utils"
)

type NodeAPI struct {
	models.BaseModel `xorm:"-"`
	Id               int    `xorm:"id pk autoincr"`
	NodeId           int    `xorm:"node_id"`
	Api              string `xorm:"api"`
}

// 表名
func (*NodeAPI) TableName() string {
	return "node_api"
}

// 节点关联的api列表
func (n *NodeAPI) ListAll(nodeId int, filter map[string]string) (*models.ModelList, []*NodeAPI, error) {
	engine := utils.Engine_backstage
	query := engine.Where("node_id=?", nodeId)

	// 获取列表数据
	var list []*NodeAPI
	err := query.Find(&list)
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}

	return n.NoPaging(len(list), list), list, nil
}

// 新增
func (n *NodeAPI) Add(nodeAPI *NodeAPI) (int, error) {
	engine := utils.Engine_backstage

	_, err := engine.Insert(nodeAPI)
	if err != nil {
		return 0, errors.NewSys(err)
	}

	return nodeAPI.Id, nil
}

// 删除
func (n *NodeAPI) Delete(id int) error {
	engine := utils.Engine_backstage
	_, err := engine.ID(id).Delete(n)
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
}
