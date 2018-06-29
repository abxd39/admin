package models

type Node struct {
	Id     int    `xorm:"not null pk autoincr INT(11)"`
	Pid    int    `xorm:"not null default 0 comment('上级ID') INT(11)"`
	Name   string `xorm:"not null default '' comment('名字') VARCHAR(36)"`
	Title  string `xorm:"not null default '' comment('标题') VARCHAR(36)"`
	Weight int    `xorm:"not null default 0 comment('权重 逆序') INT(11)"`
	States int    `xorm:"not null default 1 comment('1 正常  0 否') TINYINT(4)"`
}
