package backstage

type NodeAPI struct {
	Id     int `xorm:"id pk autoincr"`
	NodeId int `xorm:"node_id"`
	Api    int `xorm:"api"`
}

// 表名
func (*NodeAPI) TableName() string {
	return "node_api"
}
