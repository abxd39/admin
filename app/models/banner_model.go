package models

import (
	"admin/utils"
	"errors"
	"fmt"
	"time"
)

type Banner struct {
	Id          int    `xorm:"not null pk INT(11)"`
	Order       int    `xorm:"not null default 1 comment('排序') TINYINT(4)"`
	PictureName string `xorm:"not null default '' comment('图片名称') VARCHAR(255)"`
	UploadTime  string `xorm:"not null comment('展示开始日期') DATETIME"`
	LinkPath    string `xorm:"not null default '' comment('链接地址') VARCHAR(255)"`
	PicturePath string `xorm:"not null default '' comment('图片路径') VARCHAR(255)"`
	Status      int    `xorm:"not null default 1 comment('上架状态 1 上架 0下架') TINYINT(4)"`
}

func (b *Banner) Add(or, state int, picname, picp, linkaddr string) error {
	engine := utils.Engine_common
	current := time.Now().Format("2006-01-02 15:04:05")
	ban := &Banner{
		Order:       or,
		PictureName: picname,
		PicturePath: picp,
		UploadTime:  current,
		LinkPath:    linkaddr,
		Status:      state,
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

func (b *Banner) GetBannerList(page, rows, status int, start_t, end_t string) ([]Banner, int, int, error) {
	engine := utils.Engine_common
	limit := 0
	if rows <= 0 {
		rows = 100
	}
	if page <= 1 {
		page = 1
	} else {
		limit = (page - 1) * rows
	}

	list := make([]Banner, 0)
	fmt.Println("///////////////////////////////", rows, limit)
	query := engine.Desc("id")
	if status != 0 {
		query = query.Where("status=?", status)

	}
	if len(start_t) != 0 {
		query = query.Where("time_start>=?", end_t)
	}
	if len(end_t) != 0 {
		query = query.Where("time_end<=?", start_t)

	}
	Tquery := *query
	err := query.Limit(rows, limit).Find(&list)
	if err != nil {
		return nil, 0, 0, err
	}
	count, err := Tquery.Count(&Banner{})
	total_page := int(count) / rows
	return list, total_page, int(count), nil

}
