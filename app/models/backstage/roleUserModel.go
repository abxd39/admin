package backstage

type RoleUser struct {
	Id     int `xorm:"id pk autoincr"`
	RoleId int `xorm:"role_id"`
	Uid    int `xorm:"uid"`
}

// 表名
func (*RoleUser) TableName() string {
	return "role_user"
}
