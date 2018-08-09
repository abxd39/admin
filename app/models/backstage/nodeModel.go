package backstage

import (
	"admin/app/models"
	"admin/errors"
	"admin/utils"
	"fmt"
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
	MenuType         int    `xorm:"menu_type" json:"menu_type"`
	BelongSuper      int    `xorm:"belong_super" json:"belong_super"`
	DependId         string `xorm:"depend_id" json:"depend_id"`
	FullId           string `xorm:"full_id" json:"full_id"`
}

// 表名
func (*Node) TableName() string {
	return "node"
}

// 节点列表，all
func (n *Node) ListAll(filter map[string]string) (modelList *models.ModelList, list []Node, err error) {
	engine := utils.Engine_backstage
	query := engine.Desc("weight")

	// 筛选
	if v, ok := filter["states"]; ok {
		query.Where("states=?", v)
	}

	// 获取列表数据
	err = query.Find(&list)
	if err != nil {
		return nil, nil, errors.NewSys(err)
	}

	return n.NoPaging(len(list), list), list, nil
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

// 新增节点
func (n *Node) Add(node *Node) (int, error) {
	depth := 1
	var parent *Node
	var err error
	if node.Pid > 0 { // 获取上级信息
		parent, err = n.Get(node.Pid)
		if err != nil {
			return 0, err
		}

		depth = parent.Depth + 1
	}

	node.Depth = depth
	node.States = 1

	engine := utils.Engine_backstage
	session := engine.NewSession()
	err = session.Begin()
	if err != nil {
		return 0, errors.NewSys(err)
	}

	//1.新增节点
	_, err = session.Insert(node)
	if err != nil {
		session.Rollback()
		return 0, errors.NewSys(err)
	}

	//2.更新full_id
	fullId := fmt.Sprintf(",%d,", node.Id) //前后加逗号
	if node.Pid > 0 {                      // 获取上级信息
		fullId = fmt.Sprintf("%s%d,", parent.FullId, node.Id)
	}
	_, err = session.Table(n).ID(node.Id).Update(map[string]interface{}{"full_id": fullId})
	if err != nil {
		session.Rollback()
		return 0, errors.NewSys(err)
	}

	err = session.Commit()
	if err != nil {
		return 0, errors.NewSys(err)
	}

	return node.Id, nil
}

// 更新节点
func (n *Node) Update(id int, params map[string]interface{}) error {
	// 验证用户组是否存在
	engine := utils.Engine_backstage
	has, err := engine.Id(id).Exist(new(Node))
	if err != nil {
		return errors.NewSys(err)
	}
	if !has {
		return errors.NewNormal("节点不存在或已删除")
	}

	//设置更新的字段
	nodeData := make(map[string]interface{})
	if v, ok := params["title"]; ok {
		nodeData["title"] = v
	}
	if v, ok := params["type"]; ok {
		nodeData["type"] = v
	}
	if v, ok := params["depth"]; ok {
		nodeData["depth"] = v
	}
	if v, ok := params["menu_url"]; ok {
		nodeData["menu_url"] = v
	}
	if v, ok := params["menu_icon"]; ok {
		nodeData["menu_icon"] = v
	}
	if v, ok := params["menu_type"]; ok {
		nodeData["menu_type"] = v
	}
	if v, ok := params["belong_super"]; ok {
		nodeData["belong_super"] = v
	}
	if v, ok := params["states"]; ok {
		nodeData["states"] = v
	}
	if v, ok := params["weight"]; ok {
		nodeData["weight"] = v
	}

	//更新
	_, err = engine.Table(n).Id(id).Update(nodeData)
	if err != nil {
		return errors.NewSys(err)
	}

	return nil
}
