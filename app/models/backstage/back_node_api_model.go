package backstage

type NodeAPI struct {
	Id     int `xorm:"not null pk autoincr INT(11)"`
	NodeId int `xorm:"not null comment('节点ID') INT(11)"`
	Api    int `xorm:"not null comment('API接口地址') VARCHAR(200)"`
}
