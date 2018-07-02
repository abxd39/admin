package models

import (
	"admin/utils"
	"errors"
	_ "time"
)

type Banner struct {
	Id          int    `xorm:"not null pk INT(11)"`
	Order       int    `xorm:"not null default 1 comment('排序') TINYINT(4)"`
	PictureName string `xorm:"not null default '' comment('图片名称') VARCHAR(255)"`
	TimeStart   string `xorm:"not null comment('展示开始日期') DATETIME"`
	TimeEnd     string `xorm:"not null comment('展示结束日期') DATETIME"`
	LinkPath    string `xorm:"not null default '' comment('链接地址') VARCHAR(255)"`
	PicturePath string `xorm:"not null default '' comment('图片路径') VARCHAR(255)"`
	State       int    `xorm:"not null default 1 comment('上架状态 1 上架 0下架') TINYINT(4)"`
}

func (b *Banner) Add(or, state int, picname, picp, linkaddr, st, et string) error {
	engine := utils.Engine_common
	//current := time.Now().Format("2006-01-02 15:04:05")
	ban := &Banner{
		Order:       or,
		PictureName: picname,
		PicturePath: picp,
		TimeStart:   st,
		TimeEnd:     et,
		LinkPath:    linkaddr,
		State:       state,
	}
	result, err := engine.InsertOne(ban)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		return err
	}
	if 0 == result {
		err = errors.New("Unkown error")
		utils.AdminLog.Errorf(err.Error())
		return err
	}
	return nil
}
