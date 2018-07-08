package backstage

type Role struct {
	Id     int    `xorm:"not null pk autoincr INT(11)"`
	Name   string `xorm:"not null default '' VARCHAR(36)"`
	Level  int    `xorm:"not null default 1 comment('等级 1 超级管理员  2 管理员  3 运营  4 审计  5 财务') INT(6)"`
	Descr  string `xorm:"not null VARCHAR(100)"`
	People int    `xorm:"not null default 0 comment('人数') INT(6)"`
}
