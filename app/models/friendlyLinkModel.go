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

func (f *FriendlyLink) DeleteFriendlyLink(id int) error {
	engine := utils.Engine_context
	has, err := engine.Id(id).Exist(&FriendlyLink{})
	if err != nil {
		return err
	}
	if !has {
		return errors.New("this friendlylink not exits!!")
	}
	_, err = engine.Id(id).Delete(&FriendlyLink{})
	if err != nil {
		return err
	}
	return nil
}

func (f *FriendlyLink) OPeratorFriendlyLink(id, status int) error {
	engine := utils.Engine_context
	has, err := engine.Id(id).Exist(&FriendlyLink{})
	if err != nil {
		return err
	}
	if !has {
		return errors.New("this friendlylink not exits!!")
	}
	_, err = engine.Id(id).Update(&FriendlyLink{
		LinkState: status,
	})
	if err != nil {
		return err
	}
	return nil
}

func (f *FriendlyLink) GetFriendlyLink(id int) (*FriendlyLink, error) {
	engine := utils.Engine_context
	fr := new(FriendlyLink)
	has, err := engine.Id(id).Get(fr)
	if err != nil {
		return nil, err
	}
	if has {
		return fr, nil
	}
	return nil, errors.New("this friendlylink not exits!!")
}

func (f *FriendlyLink) Add(id, order, state int, wn, ln string) error {
	engine := utils.Engine_context
	flink := &FriendlyLink{
		Aorder:    order,
		WebName:   wn,
		LinkName:  ln,
		LinkState: state,
	}
	has, err := engine.Id(id).Exist(&FriendlyLink{})
	if err != nil {
		return err
	}
	if has {
		_, err := engine.Id(id).Update(flink)
		if err != nil {
			return err
		}
		return nil
	}

	result, err := engine.Insert(flink)
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

func (f *FriendlyLink) GetFriendlyLinkList(page, rows, status int, search string) (*ModelList, error) {
	engine := utils.Engine_context
	fmt.Println("...................................")
	//query := engine.Desc("id")
	countsql := "select count(*) from friendly_link  "
	var sql string = "select  * from friendly_link "
	if len(search) > 0 && status != 0 {

		//sql = fmt.Sprintf("select * from friendly_link where concat(IFNULL(web_name,''),IFNULL(link_name,'')) LIKE '%%%s%%' ", search)
		sub := fmt.Sprintf("where concat(IFNULL(web_name,''),IFNULL(link_name,'')) LIKE '%%%s%%' AND link_state=%d ; ", search, status)
		sql += sub
		countsql += sub
		//fmt.Println("uname =", sql)
	} else if len(search) > 0 {
		sub := fmt.Sprintf("where concat(IFNULL(web_name,''),IFNULL(link_name,'')) LIKE '%%%s%%' ;", search)
		sql += sub
		countsql += sub
	} else if status != 0 {
		sub := fmt.Sprintf("where link_state=%d ;", status)
		sql += sub
		countsql += sub
	}

	count, err := engine.Sql(countsql).Count(&FriendlyLink{})
	if err != nil {
		utils.AdminLog.Errorln("统计所有记录失败")
		return nil, err
	}
	fmt.Println("count =", count)

	//获取分页
	offset, modelList := f.Paging(page, rows, int(count))
	fmt.Println("page=", modelList.PageSize, "offset=", offset)
	//查询数据
	friendlist := make([]FriendlyLink, 0)

	err = engine.Sql(sql).Limit(modelList.PageSize, offset).Find(&friendlist)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return nil, err
	}
	fmt.Println("friendlist-len=", len(friendlist))
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
