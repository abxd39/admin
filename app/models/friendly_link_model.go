package models

import (
	"admin/utils"
	"errors"
	"fmt"
)

type FriendlyLink struct {
	BaseModel `xorm:"-"`
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

func (f *FriendlyLink) GetFriendlyLinkList(page, rows, status int, name, link_name string) (*ModelList, error) {
	engine := utils.Engine_context
	//page !=0

	u := &FriendlyLink{}
	query := engine.Desc("id")

	if len(name) > 0 {
		query = query.Where("web_name=?", name)
	}
	if len(link_name) > 0 {
		query = query.Where("link_name=?", link_name)
	}
	if status != 0 {
		query = query.Where("link_state=?", status)
	}
	tempquery := *query
	count, err := tempquery.Count(u)
	if err != nil {
		utils.AdminLog.Errorln("统计所有记录失败")
		return nil, err
	}
	//获取分页
	offset, modelList := f.Paging(page, rows, int(count))
	friendlist := make([]FriendlyLink, 0)
	err = engine.Limit(modelList.PageSize, offset).Find(&friendlist)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return nil, err
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

	modelList.Items = link

	return modelList, nil
}
