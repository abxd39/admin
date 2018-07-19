package backstage

// 用户组、权限节点关联
type RoleNode struct {
	Id     int `xorm:"id pk autoincr"`
	RoleId int `xorm:"role_id"`
	NodeId int `xorm:"node_id"`
}

// 表名
func (*RoleNode) TableName() string {
	return "role_node"
}
