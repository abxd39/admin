package models

import (
	"admin/utils"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var remoteurl string = "https://sdun.oss-cn-shenzhen.aliyuncs.com/"

type Article struct {
	BaseModel     `xorm:"-"`
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
	BaseModel  `xorm:"-"`
	Id         int    `xorm:"not null pk autoincr comment('自增ID') INT(10)"`
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
	BaseModel `xorm:"-"`
	Id        int    `xorm:"not null pk autoincr MEDIUMINT(6)"`
	TypeId    int    `xorm:"not null default 0 TINYINT(10)"`
	TypeName  string `xorm:"not null default '' comment('类型名称 1关于我们，2媒体报道，3联系我们，4团队介绍，5数据资产介绍，6服务条款，7免责声明，8隐私保护9 业界新闻 10 公告 11 帮助手册 12 币种介绍') VARCHAR(100)"`
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

func (a *ArticleList) GetArticleList(page, rows, tp, status int, st, et string) (*ModelList, error) {
	engine := utils.Engine_common
	query := engine.Desc("id")
	query = query.Where("type=?", tp)
	TempQuery := *query
	count, err := TempQuery.Count(&Article{})
	if err != nil {
		return nil, err
	}
	offset, modelList := a.Paging(page, rows, int(count))

	u := make([]Article, 0)
	err = query.Limit(modelList.PageSize, offset).Find(&u)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return nil, err
	}
	list := make([]*ArticleList, 0)
	for _, v := range u {
		ret := ArticleList{
			Id:         v.Id,
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

	modelList.Items = list
	return modelList, nil

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

func (a *Article) LocalFileToAliCloud(filePath string) (string, error) {
	client := utils.AliClient
	bucket, err := client.Bucket("sdun")
	if err != nil {
		return "", err
	}
	fmt.Println("111111111111111", filePath)
	//读取内容做has mde5
	// fd, err := os.OpenFile(filePath, os.O_RDONLY, 0660)
	// if err != nil {
	// 	// HandleError(err)
	// 	return "", err
	// }

	// body, err := ioutil.ReadAll(fd)
	// if err != nil {
	// 	fmt.Println("ReadAll", err)
	// 	return "", err
	// }
	// fd.Close()
	h := md5.New()
	h.Write([]byte(filePath)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	okey := hex.EncodeToString(cipherStr)
	fmt.Println(okey)
	fSuffix := ".png" //path.Ext(filePath)
	okey += fSuffix
	fmt.Printf("%#v\n", okey)
	ddd, _ := base64.StdEncoding.DecodeString(filePath) //成图片文件并把文件写入到buffer

	fmt.Println("111111111111111", ddd)

	err = ioutil.WriteFile("./output133.png", []byte(filePath), 0666)

	err = bucket.PutObject(okey, strings.NewReader(filePath))

	//err = bucket.PutObjectFromFile(okey, filePath)
	if err != nil {
		fmt.Println(filePath)
		return "", err
	}
	fmt.Println(remoteurl + okey)
	return remoteurl + okey, nil
}

func (a *Article) GetLocalFileToAliCloud(object_key, filepath string) (string, error) {
	client := utils.AliClient
	bucket, err := client.Bucket("sdun")
	if err != nil {
		return "", err
	}
	lsRes, err := bucket.ListObjects()
	if err != nil {
		// HandleError(err)
		return "", nil
	}

	for _, object := range lsRes.Objects {
		fmt.Println("Objects:", object.Key)
	}
	body, err := bucket.GetObject(object_key)
	if err != nil {
		// HandleError(err)
		return "", err
	}
	fd, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		// HandleError(err)
		return "", nil
	}
	defer fd.Close()

	io.Copy(fd, body)
	return filepath, nil
}
