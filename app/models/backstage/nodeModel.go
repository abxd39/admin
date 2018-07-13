package backstage

import (
	"admin/app/models"
	"admin/errors"
	"admin/utils"
)

// 权限节点
type Node struct {
	models.BaseModel `xorm:"-"`
	Id               int    `xorm:"not null pk autoincr INT(11)" json:"id"`
	Pid              int    `xorm:"not null default 0 comment('上级ID') INT(11)" json:"pid"`
	Title            string `xorm:"not null default '' comment('标题') VARCHAR(36)" json:"title"`
	Weight           int    `xorm:"not null default 0 comment('权重 逆序') INT(11)" json:"weight"`
	States           int    `xorm:"not null default 1 comment('1 正常  0 否') TINYINT(4)" json:"states"`
	Type             int    `xorm:"not null comment('权限类型，1 菜单，2 功能') TINYINT(1)" json:"type"`
	Depth            int    `xorm:"not null comment('深度，从1开') INT(11)" json:"depth"`
	MenuUrl          string `xorm:"not null default '' comment('菜单地址') VARCHAR(150)" json:"menu_url"`
	MenuIcon         string `xorm:"not null default '' comment('菜单图标') VARCHAR(50)" json:"menu_icon"`
	BelongSuper      int    `xorm:"not null default 0 comment('属于超管的权限，0 否，1 是') TINYINT(1)" json:"belong_super"`
	DependId         string `xorm:"not null default '' comment('依赖ID，逗号分隔') VARCHAR(500)" json:"depend_id"`
	FullId           string `xorm:"not null comment('全ID，逗号分隔，前后加逗号') VARCHAR(500)" json:"full_id"`
}

// 节点列表，all
func (n *Node) ListAll(pageIndex, pageSize int, filter map[string]string) (list []Node, err error) {
	// 获取总数
	engine := utils.Engine_backstage
	query := engine.Desc("weight")

	// 筛选
	if v, ok := filter["states"]; ok {
		query.Where("states=?", v)
	}

	// 获取列表数据
	err = query.Find(&list)
	if err != nil {
		return nil, errors.NewSys(err)
	}

	return
}
