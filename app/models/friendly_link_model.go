package models

import (
	"admin/utils"
	"errors"
	"fmt"
)

type FriendlyLink struct {
	Id        int    `xorm:"not null pk autoincr comment('自增id') INT(10)"`
	Aorder    int    `xorm:"not null comment('排序') INT(10)"`
	WebName   string `xorm:"not null default '' comment('网址名称') VARCHAR(100)"`
	LinkName  string `xorm:"not null default '' comment('网站链接') VARCHAR(100)"`
	LinkState int    `xorm:"not null comment('1,上架2，下架') INT(10)"`
}

func (f *FriendlyLink) Add(order, state int, wn, ln string) error {

	flink := &FriendlyLink{
		Aorder:    order,
		WebName:   wn,
		LinkName:  ln,
		LinkState: state,
	}
	result, err := utils.Engine_context.Insert(flink)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return nil
	}
	if result == 0 {
		err = errors.New("friendly link insert fail!!")
		utils.AdminLog.Errorf(err.Error())
		return err
	}

	return nil
}

func (f *FriendlyLink) GetFriendlyLinkList(page, count int) ([]*FriendlyLink, int, int, error) {
	engine := utils.Engine_context
	//page !=0
	if 1 >= page {
		page = 1
	}
	var limit int
	if 1 >= page {
		limit = 0
	} else {
		limit = (page - 1) * count
	}

	defa := count
	if 0 == defa {
		count = 100
	}
	u := &FriendlyLink{}

	friendlist := make([]FriendlyLink, 0)
	fmt.Println("count=", count, "limit=", limit)
	err := engine.Limit(count, limit).Find(&friendlist)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return nil, 0, 0, err
	}
	fmt.Println("00000000000000000000000000")
	link := make([]*FriendlyLink, 0)
	for _, frd := range friendlist {
		ret := FriendlyLink{
			Id:        frd.Id,
			Aorder:    frd.Aorder,
			WebName:   frd.WebName,
			LinkName:  frd.LinkName,
			LinkState: frd.LinkState,
		}
		link = append(link, &ret)
	}
	fmt.Println("1111111111111111111111111111111", link)

	total, err := engine.Count(u)
	if err != nil {
		utils.AdminLog.Errorln("统计所有记录失败")
		return nil, 0, 0, err
	}

	page = int(total) / count

	return link, page, int(total), nil
}
