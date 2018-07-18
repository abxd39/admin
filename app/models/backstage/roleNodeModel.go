package backstage

// 用户组、权限节点关联
type RoleNode struct {
	Id     int `xorm:"not null pk autoincr INT(11)"`
	RoleId int `xorm:"not null default 0 comment('角色id') INT(11)"`
	NodeId int `xorm:"not null default 0 comment('节点id') INT(11)"`
}

// 表名
func (*RoleNode) TableName() string {
	return "role_node"
}
