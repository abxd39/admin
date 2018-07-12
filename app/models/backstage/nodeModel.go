package backstage

type Node struct {
	Id          int    `xorm:"not null pk autoincr INT(11)"`
	Pid         int    `xorm:"not null default 0 comment('上级ID') INT(11)"`
	Title       string `xorm:"not null default '' comment('标题') VARCHAR(36)"`
	Weight      int    `xorm:"not null default 0 comment('权重 逆序') INT(11)"`
	States      int    `xorm:"not null default 1 comment('1 正常  0 否') TINYINT(4)"`
	Type        int    `xorm:"not null comment('权限类型，1 菜单，2 功能') TINYINT(1)"`
	Depth       int    `xorm:"not null comment('深度，从1开') INT(11)"`
	MenuUrl     string `xorm:"not null default '' comment('菜单地址') VARCHAR(150)"`
	MenuIcon    string `xorm:"not null default '' comment('菜单图标') VARCHAR(50)"`
	BelongSuper int    `xorm:"not null default 0 comment('属于超管的权限，0 否，1 是') TINYINT(1)"`
	DependId    string `xorm:"not null default '' comment('依赖ID，逗号分隔') VARCHAR(500)"`
	FullId      string `xorm:"not null comment('全ID，逗号分隔，前后加逗号') VARCHAR(500)"`
}
