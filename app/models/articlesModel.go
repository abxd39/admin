package models

import (
	"admin/utils"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
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
	engine := utils.Engine_context
	list := make([]ArticleType, 0)
	err := engine.Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (a *ArticleList) GetArticleList(page, rows, tp, status int, title, st string) (*ModelList, error) {
	fmt.Println("GetArticleList")
	engine := utils.Engine_context
	query := engine.Desc("id")
	if tp != 0 {
		query = query.Where("type=?", tp)
	}
	if status != 0 {
		query = query.Where("astatus=?", status)
	}
	if len(title) != 0 {
		temp := fmt.Sprintf(" concat(IFNULL(title,'')) LIKE '%%%s%%'  ", title)
		query = query.Where(temp)
	}
	if st != `` {
		substr := st[:11] + "23:59:59"
		temp := fmt.Sprintf("create_time BETWEEN '%s' AND '%s' ", st, substr)
		query = query.Where(temp)
	}

	TempQuery := *query
	count, err := TempQuery.Count(&Article{})
	if err != nil {
		return nil, err
	}
	offset, modelList := a.Paging(page, rows, int(count))
	fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", offset, modelList.PageSize)
	u := make([]Article, count)
	err = query.Limit(modelList.PageSize, offset).Find(&u)
	if err != nil {
		utils.AdminLog.Errorln(err.Error())
		return nil, err
	}
	fmt.Println("111111111111111111111111111111111111", len(u))
	list := make([]ArticleList, 0)
	for _, v := range u {
		if v.Id == 0 {
			continue
		}
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
		list = append(list, ret)
	}

	modelList.Items = list
	return modelList, nil

}

//上下架 文章
func (a *Article) UpArticle(id, status int) error {
	engine := utils.Engine_context
	art := new(Article)
	has, _ := engine.Exist(art)
	if !has {
		return errors.New("文章不存在！！")
	}
	current := time.Now().Format("2006-01-02 15:04:05")
	//下架
	_, err := engine.Id(id).Update(&Article{
		Astatus:    status,
		UpdateTime: current,
	})
	if err != nil {
		return err
	}

	return nil
}

//获取文章
func (a *Article) GetArticle(id int) (*Article, error) {
	engine := utils.Engine_context
	art := new(Article)
	_, err := engine.Id(id).Get(art)
	if err != nil {
		return nil, err
	}
	return art, nil
}

//删除文章
func (a *Article) DeleteArticle(id int) error {
	engine := utils.Engine_context
	has, _ := engine.Id(id).Exist(&Article{})
	if !has {
		return errors.New("文章不存在")
	}
	_, err := engine.Id(id).Delete(&Article{})
	if err != nil {
		return err
	}
	return nil
}

func (a *Article) AddArticle(id int, u *Article) error {
	engine := utils.Engine_context
	if id != 0 {
		fmt.Println("cccccc")
		art := new(Article)
		has, err := engine.Id(id).Get(art)
		if err != nil {
			return err
		}
		if has {

			if v := strings.Compare(art.Covers, u.Covers); v != 0 {
				a.DeletFileToAliCloud(art.Covers)
			}
			_, err = engine.Id(id).Update(u)
			if err != nil {
				return err
			}
			return nil
		}
		return errors.New("Article is not exist !!!")
	}
	current := time.Now().Format("2006-01-02 15:04:05")
	u.CreateTime = current
	u.UpdateTime = current
	result, err := engine.InsertOne(u)
	if err != nil {
		return err
	}
	if result == 0 {
		utils.AdminLog.Errorln("article InsertOne failed ")
	}
	return nil
}

//删除oss对象
func (a *Article) DeletFileToAliCloud(filepath string) error {
	client := utils.AliClient
	bucket, err := client.Bucket("sdun")
	if err != nil {
		return err
	}
	index := strings.LastIndex(filepath, "//")
	if index < 0 {
		return errors.New("oss object delete failed!!")
	}
	substr := filepath[index+1:]
	isExist, err := bucket.IsObjectExist(substr)
	if err != nil {
		return err
	}
	if isExist {
		err := bucket.DeleteObject(substr)
		if err != nil {
			return nil
		}
	}

	return nil
}

//上传Ali coud
func (a *Article) LocalFileToAliCloud(filePath string) (string, error) {
	client := utils.AliClient
	bucket, err := client.Bucket("sdun")
	if err != nil {
		return "", err
	}
	//查找base64
	fmt.Println("base34-1")
	base := strings.Index(filePath, ";base64,")
	if base < 0 {
		fmt.Println("base34-3")
		// 是远程的oss 文件路径
		return filePath, nil
	}
	fmt.Println("base34-2")
	fmt.Println(filePath)
	// if len(remotePath) != 0 {
	// 	index := strings.LastIndex(remotePath, "//")
	// 	if index <= 0 {
	// 		return "", errors.New("oss okject no exits!!")
	// 	}
	// 	substr := remotePath[index+1:]
	// 	isExist, err := bucket.IsObjectExist(substr)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	if isExist {

	// 	}
	// }
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	subm := strings.IndexByte(filePath, ',')
	if subm < 0 {
		return "", errors.New("find failed!!")
	}
	substr := filePath[:subm]
	subb := strings.IndexByte(substr, '/')
	sube := strings.IndexByte(substr, ';')
	if subb < 0 || sube < 0 {
		return "", errors.New("find fail!!")
	}
	fmt.Println(subb, sube, subm)
	fSuffix := substr[subb+1 : sube]
	value := filePath[subm+1:]
	h := md5.New()
	tempValue := value
	tempValue += timestamp
	h.Write([]byte(tempValue)) // 需要加密的字符串为 123456
	cipherStr := h.Sum(nil)
	okey := hex.EncodeToString(cipherStr)
	fmt.Println(okey)
	okey += "."
	okey += fSuffix
	fmt.Printf("%#v\n", okey)
	ddd, _ := base64.StdEncoding.DecodeString(value)
	err = bucket.PutObject(okey, bytes.NewReader(ddd))
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
