package backstage

type RoleUser struct {
	Id     int `xorm:"not null pk autoincr INT(11)"`
	RoleId int `xorm:"not null INT(11)"`
	Uid    int `xorm:"not null INT(11)"`
}
