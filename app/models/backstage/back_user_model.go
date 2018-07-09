package backstage

import (
	"admin/utils"
	"errors"
	"fmt"
)

type User struct {
	Id          int    `xorm:"not null pk autoincr INT(11)"`
	Uid         int    `xorm:"not null INT(11)"`
	Phone       string `xorm:"not null default '' comment('电话') VARCHAR(20)"`
	NickName    string `xorm:"not null default '' comment('昵称') VARCHAR(60)"`
	Pwd         string `xorm:"not null default '' comment('用户登录密码') VARCHAR(255)"`
	States      int    `xorm:"not null default 1 comment('1正常  0 锁定') TINYINT(4)"`
	Remark      string `xorm:"not null default '' comment('备注') VARCHAR(36)"`
	CreatedTime string `xorm:"not null comment('创建时间') DATETIME"`
	UpdatedTime string `xorm:"not null comment('修改时间') DATETIME"`
}

func (u *User) Login(pwd, phone string) (string, int, error) {
	engine := utils.Engine_backstage
	fmt.Println("login")
	use := &User{}
	_, err := engine.Where("phone=?", phone).Get(use)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		fmt.Println("login", err.Error())
		return "", 0, err
	}
	//find 如果不存在 数据库是否会返回 一个错误给我
	if use.States == 0 {
		return "", 0, errors.New("该用户已锁定")
	}
	fmt.Printf("数据库中的has值为%s", use.Pwd)
	if pwd != use.Pwd {
		return "", 0, errors.New("密码不对！！")
	}
	//
	return use.NickName, use.Uid, nil
}