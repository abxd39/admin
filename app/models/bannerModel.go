package models

import (
	"admin/utils"
	"errors"
	"time"
)

type Banner struct {
	BaseModel   `xorm:"-"`
	Id          int    `xorm:"not null pk INT(11)"`
	Order       int    `xorm:"not null default 1 comment('排序') TINYINT(4)"`
	PictureName string `xorm:"not null default '' comment('图片名称') VARCHAR(255)"`
	UploadTime  string `xorm:"not null comment('展示开始日期') DATETIME"`
	LinkPath    string `xorm:"not null default '' comment('链接地址') VARCHAR(255)"`
	PicturePath string `xorm:"not null default '' comment('图片路径') VARCHAR(255)"`
	Status      int    `xorm:"not null default 1 comment('上架状态 1 上架 0下架') TINYINT(4)"`
}

func (b *Banner) DeleteBanner(id int) error {
	engine := utils.Engine_common
	has, _ := engine.Id(id).Exist(&Banner{})
	if !has {
		return errors.New("banner is not exist!!")
	}
	ba := new(Banner)
	engine.Id(id).Get(ba)
	new(Article).DeletFileToAliCloud(ba.PicturePath)

	_, err := engine.Id(id).Delete(&Banner{})
	if err != nil {
		return err
	}
	return nil
}

func (b *Banner) GetBanner(id int) (*Banner, error) {
	engine := utils.Engine_common
	ba := new(Banner)
	_, err := engine.Id(id).Get(ba)
	if err != nil {
		return nil, err
	}
	new(Article).DeletFileToAliCloud(ba.PicturePath)

	return ba, nil
}

func (b *Banner) OperatorUp(id, mark int) error {
	engine := utils.Engine_common
	current := time.Now().Format("2006-01-02 15:04:05")
	_, err := engine.Id(id).Update(&Banner{
		Status:     mark,
		UploadTime: current,
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *Banner) Add(id, or, state int, picname, picp, linkaddr string) error {
	engine := utils.Engine_common

	current := time.Now().Format("2006-01-02 15:04:05")
	ban := &Banner{
		Id:          id,
		Order:       or,
		PictureName: picname,
		PicturePath: picp,
		UploadTime:  current,
		LinkPath:    linkaddr,
		Status:      state,
	}
	if id != 0 {
		_, err := engine.Id(id).Update(ban)
		if err != nil {
			return err
		}
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

func (b *Banner) GetBannerList(page, rows, status int, start_t, end_t string) (*ModelList, error) {
	engine := utils.Engine_common
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
	count, err := Tquery.Count(&Banner{})
	if err != nil {
		return nil, err
	}
	//获取分页
	offset, modelList := b.Paging(page, rows, int(count))

	list := make([]Banner, 0)
	//断言
	//if _,ok:= modelist.Items.()
	err = query.Limit(modelList.PageSize, offset).Find(&list)
	if err != nil {
		return nil, err
	}

	modelList.Items = list
	return modelList, nil

}
