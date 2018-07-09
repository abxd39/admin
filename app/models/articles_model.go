package models

import (
	"admin/utils"
	"fmt"
)

type Article struct {
	Id            int    `xorm:"not null pk autoincr comment('自增ID') INT(10)"`
	Title         string `xorm:"not null default '' comment('文章标题') VARCHAR(100)"`
	Description   string `xorm:"not null default '' comment('描述') VARCHAR(1000)"`
	Content       string `xorm:"not null comment('内容') TEXT"`
	Covers        string `xorm:"not null default '' comment('封面图片') VARCHAR(1000)"`
	ContentImages string `xorm:"not null comment('内容图片') TEXT"`
	Type          int    `xorm:"not null default 1 comment('类型 1 业界新闻 2 公告 3 帮助手册') TINYINT(4)"`
	TypeName      string `xorm:"not null default '' comment('类型名字') VARCHAR(50)"`
	Author        string `xorm:"not null default '' comment('作者名字') VARCHAR(150)"`
	Weight        int    `xorm:"not null default 0 comment('权重，排序字段') TINYINT(4)"`
	Shares        int    `xorm:"not null default 0 comment('分享数量') INT(11)"`
	Hits          int    `xorm:"not null default 0 comment('点击数量') INT(11)"`
	Comments      int    `xorm:"not null default 0 comment('评论数量') INT(11)"`
	Astatus       int    `xorm:"not null default 1 comment('1 显示 0 不显示') TINYINT(1)"`
	CreateTime    string `xorm:"not null default '' comment('创建时间') VARCHAR(36)"`
	UpdateTime    string `xorm:"not null VARCHAR(36)"`
	AdminId       int    `xorm:"not null INT(4)"`
	AdminNickname string `xorm:"not null default '' comment('管理员名字') VARCHAR(50)"`
}

type ArticleList struct {
	Weight     int    `xorm:"not null default 0 comment('权重，排序字段') TINYINT(4)"`
	Title      string `xorm:"not null default '' comment('文章标题') VARCHAR(100)"`
	Author     string `xorm:"not null default '' comment('作者名字') VARCHAR(150)"`
	Covers     string `xorm:"not null default '' comment('封面图片') VARCHAR(1000)"`
	CreateTime string `xorm:"not null default '' comment('创建时间') VARCHAR(36)"`
	Hits       int    `xorm:"not null default 0 comment('点击数量') INT(11)"`
	Astatus    int    `xorm:"not null default 1 comment('1 显示 0 不显示') TINYINT(1)"`
	Type       int    `xorm:"not null default 1 comment('类型 1 业界新闻 2 公告 3 帮助手册') TINYINT(4)"`
}

type ArticleType struct {
	Id       int    `xorm:"not null pk autoincr MEDIUMINT(6)"`
	TypeId   int    `xorm:"not null default 0 TINYINT(10)"`
	TypeName string `xorm:"not null default '' comment('类型名称 1关于我们，2媒体报道，3联系我们，4团队介绍，5数据资产介绍，6服务条款，7免责声明，8隐私保护9 业界新闻 10 公告 11 帮助手册 12 币种介绍') VARCHAR(100)"`
}

func (a *ArticleList) TableName() string {
	return "article"
}

func (a *ArticleType) GetArticleType() ([]ArticleType, error) {
	engine := utils.Engine_common
	list := make([]ArticleType, 0)
	err := engine.Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (a *ArticleList) GetArticleList(page, rows, tp int) ([]*ArticleList, int, int, error) {

	if page <= 0 {
		page = 1
	}

	if rows <= 0 {
		rows = 100
	}
	var start_rows int
	if page > 1 {
		start_rows = (page - 1) * rows
	}
	engine := utils.Engine_common
	fmt.Println("type=", tp, "page=", page, "起始行star_row=", start_rows, "page_num=", rows)
	u := make([]Article, 0)

	query := engine.Where("type=?", tp)
	TempQuery := query
	err := query.Limit(rows, start_rows).Find(&u)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return nil, 0, 0, err
	}
	list := make([]*ArticleList, 0)
	for _, v := range u {
		ret := ArticleList{
			Weight:     v.Weight,
			Title:      v.Title,
			Author:     v.Author,
			Covers:     v.Covers,
			CreateTime: v.CreateTime,
			Hits:       v.Hits,
			Astatus:    v.Astatus,
			Type:       v.Type,
		}
		list = append(list, &ret)
	}
	count, err := TempQuery.Count(&Article{})
	if err != nil {
		return nil, 0, 0, err
	}

	total_page := int(count) / rows
	return list, total_page, int(count), nil

}

func (a *Article) AddArticle(u *Article) error {
	engine := utils.Engine_common
	result, err := engine.InsertOne(u)
	if err != nil {
		return err
	}
	if result == 0 {
		utils.AdminLog.Errorln("article InsertOne failed ")
	}
	return nil
}
