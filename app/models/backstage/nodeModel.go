package backstage

import (
	"admin/app/models"
	"admin/errors"
	"admin/utils"
)

// 权限节点
type Node struct {
	models.BaseModel `xorm:"-"`
	Id               int    `xorm:"id pk autoincr" json:"id"`
	Pid              int    `xorm:"pid" json:"pid"`
	Title            string `xorm:"title" json:"title"`
	Weight           int    `xorm:"weight" json:"weight"`
	States           int    `xorm:"states" json:"states"`
	Type             int    `xorm:"type" json:"type"`
	Depth            int    `xorm:"depth" json:"depth"`
	MenuUrl          string `xorm:"menu_url" json:"menu_url"`
	MenuIcon         string `xorm:"menu_icon" json:"menu_icon"`
	BelongSuper      int    `xorm:"belong_super" json:"belong_super"`
	DependId         string `xorm:"depend_id" json:"depend_id"`
	FullId           string `xorm:"full_id" json:"full_id"`
}

// 表名
func (*Node) TableName() string {
	return "node"
}

// 节点列表，all
func (n *Node) ListAll(filter map[string]string) (modelList *models.ModelList, err error) {
	// 获取总数
	engine := utils.Engine_backstage
	query := engine.Desc("weight")

	// 筛选
	if v, ok := filter["states"]; ok {
		query.Where("states=?", v)
	}

	// 获取列表数据
	var list []Node
	err = query.Find(&list)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	return n.NoPaging(len(list), list), nil
}

// 节点详情
func (n *Node) Get(id int) (node *Node, err error) {
	engine := utils.Engine_backstage
	node = new(Node)
	has, err := engine.Id(id).Get(node)
	if err != nil {
		return nil, errors.NewSys(err)
	}
	if !has {
		return nil, errors.NewNormal("权限节点不存在")
	}

	return
}
