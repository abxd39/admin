package controller

import (
	"admin/app/models"
	"admin/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ContextController struct {
	BaseController
}

func (this *ContextController) Router(r *gin.Engine) {
	g := r.Group("/content")
	{
		g.POST("/add_link", this.AddFriendlyLink)
		g.GET("/link_list", this.GetFriendlyLink)
		g.GET("/article_list", this.GetArticleList)          //获取文章列表
		g.POST("/add_banner", this.AddBanner)                //添加banner
		g.POST("/upordown_banner", this.UpBanner)            //修改banner的上下架状态
		g.GET("/get_banner", this.GetBanner)                 //获取bannner
		g.POST("/delete_banner", this.DeleteBanner)          //上传banner  同时删除 oss 上的图片对象
		g.GET("/banner_list", this.GetBannerList)            //获取 banner 列表
		g.POST("/add_article", this.AddArticle)              // 添加文章
		g.GET("/article_type", this.GetArticleType)          //获取文章类型
		g.POST("/local_filetoali", this.LocalFileToAliCloud) //上传图片到 oss
		g.POST("/delete_article", this.DeleteArticle)        //删除文章
		g.GET("/get_article", this.GetArticle)               //获取文章
		g.POST("/upordown_article", this.UpArticle)          //上下架文章
	}
}

func (this *ContextController) UpArticle(c *gin.Context) {
	req := struct {
		Id   int `form:"id" json:"id" binding:"required"` //文章Id
		Mark int `form:"op" json:"op" binding:"required"` //mark 1 下架 2 删除

	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		this.RespErr(c, err)
		return
	}
	err = new(models.Article).UpArticle(req.Id, req.Mark)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.RespOK(c)
	return
}

func (this *ContextController) GetArticle(c *gin.Context) {
	req := struct {
		Id int `form:"id" json:"id" binding:"required"` //文章Id
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		this.RespErr(c, err)
		return
	}
	result, err := new(models.Article).GetArticle(req.Id)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "data", result)
	this.RespOK(c)
	return
}

func (this *ContextController) DeleteArticle(c *gin.Context) {
	req := struct {
		Id int `form:"id" json:"id" binding:"required"` //文章Id
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		this.RespErr(c, err)
		return
	}
	err = new(models.Article).DeleteArticle(req.Id)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.RespOK(c)
	return
}

func (this *ContextController) DeleteBanner(c *gin.Context) {
	req := struct {
		Id int `form:"id" json:"id" binding:"required"` //文章Id
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		this.RespErr(c, err)
		return
	}
	err = new(models.Banner).DeleteBanner(req.Id)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.RespOK(c, "成功")
	return
}

func (this *ContextController) GetBanner(c *gin.Context) {
	req := struct {
		Id int `form:"id" json:"id" binding:"required"` //文章Id
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		this.RespErr(c, err)
		return
	}
	result, err := new(models.Banner).GetBanner(req.Id)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "data", result)
	this.RespOK(c, "成功")
	return
}

func (this *ContextController) UpBanner(c *gin.Context) {
	req := struct {
		Id   int `form:"id" json:"id" binding:"required"` //文章Id
		Mark int `form:"op" json:"op" binding:"required"` //mark 1 下架 2 删除

	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		this.RespErr(c, err)
		return
	}
	err = new(models.Banner).OperatorUp(req.Id, req.Mark)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		this.RespErr(c, err)
		return
	}
	this.RespOK(c, "成功")
	return
}

func (this *ContextController) LocalFileToAliCloud(c *gin.Context) {
	req := struct {
		FilePath string `form:"file_path" json:"file_path" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		this.RespErr(c, err)
		return
	}
	_, err = new(models.Article).LocalFileToAliCloud(req.FilePath)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	response := new(models.PolicyToken).Get_policy_token()
	// c.Request.Header.Set("Access-Control-Allow-Methods", "POST")
	// c.Request.Header.Set("Access-Control-Allow-Origin", "*")
	//io.WriteString(c, response)
	this.Put(c, "data", response)
	this.RespOK(c, "成功")
	return
}

//
func (this *ContextController) GetArticleType(c *gin.Context) {
	result, err := new(models.ArticleType).GetArticleType()
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.Put(c, "list", result)
	this.RespOK(c, "成功")
	return
}
func (this *ContextController) AddFriendlyLink(c *gin.Context) {

	req := struct {
		Order     string `form:"order" json:"order" binding:"required"`
		WebName   string `form:"web_name" json:"web_name" binding:"required"`
		LinkName  string `form:"link_name" json:"link_name" binding:"required"`
		LinkState int    `form:"link_state" json:"link_state" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	fmt.Println("..........................................")
	err = new(models.FriendlyLink).Add(1, req.LinkState, req.WebName, req.LinkName)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	//response := new(models.PolicyToken).Get_policy_token()
	// c.Request.Header.Set("Access-Control-Allow-Methods", "POST")
	// c.Request.Header.Set("Access-Control-Allow-Origin", "*")
	this.RespOK(c, "成功")
	return
}

func (this *ContextController) GetFriendlyLink(c *gin.Context) {

	req := struct {
		Page     int    `form:"page" json:"page" binding:"required"`
		Rows     int    `form:"rows" json:"rows" `
		Name     string `form:"name" json:"name" `
		Linkname string `form:"link_name" json:"link_name" `
		Status   int    `form:"status" json:"status" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	//operator db         GetFriendlyLinkList
	list, err := new(models.FriendlyLink).GetFriendlyLinkList(req.Page, req.Rows, req.Status, req.Name, req.Linkname)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	// 设置返回数据
	this.Put(c, "list", list)

	// 返回
	this.RespOK(c, "成功")
	return
}

func (this *ContextController) AddBanner(c *gin.Context) {
	req := struct {
		Id          int    `form:"id" json:"id"`
		Order       int    `form:"order" json:"order" binding:"required"`
		PictureName string `form:"picture_n" json:"picture_n" binding:"required"`
		PicturePath string `form:"picture_p" json:"picture_p" binding:"required"`
		LinkAddr    string `form:"link_addr" json:"link_addr" binding:"required"`
		Status      int    `form:"status" json:"status" binding:"required"`
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	path, err := new(models.Article).LocalFileToAliCloud(req.PicturePath)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	err = new(models.Banner).Add(req.Id, req.Order, req.Status, req.PictureName, path, req.LinkAddr)
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.RespOK(c, "成功")
	return
}

func (this *ContextController) GetBannerList(c *gin.Context) {
	req := struct {
		Page    int    `form:"page" json:"page" binding:"required"`
		Rows    int    `form:"rows" json:"rows" `
		Start_t string `form:"start_t" json:"start_t" `
		End_t   string `form:"end_t" json:"end_t" `
		Status  int    `form:"status" json:"status" `
	}{}
	fmt.Printf("GetBannerList%#v\n", req)
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.Banner).GetBannerList(req.Page, req.Rows, req.Status, req.Start_t, req.End_t)
	if err != nil {
		this.RespErr(c, err)
		return
	}

	// 设置返回数据
	this.Put(c, "list", list)

	// 返回
	this.RespOK(c, "成功")
	return
}

func (this *ContextController) GetArticleList(c *gin.Context) {
	req := struct {
		Page    int    `form:"page" json:"page" binding:"required"`
		Rows    int    `form:"rows" json:"rows" `
		Type    int    `form:"type" json:"type" binding:"required"`
		Start_t string `form:"start_t" json:"start_t" `
		End_t   string `form:"end_t" json:"end_t" `
		Status  int    `form:"status" json:"status" `
	}{}
	fmt.Println("获取文章列表")
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}
	list, err := new(models.ArticleList).GetArticleList(req.Page, req.Rows, req.Type, req.Status, req.Start_t, req.End_t)
	if err != nil {
		this.RespErr(c, err)
		return
	}

	// 设置返回数据
	this.Put(c, "list", list)

	// 返回
	this.RespOK(c, "成功")
	return
}

func (this *ContextController) AddArticle(c *gin.Context) {
	req := struct {
		Id            int    `form:"id" json:"id" `
		Title         string `form:"title" json:"title" binding:"required"`
		Description   string `form:"desc" json:"desc" `
		ContentText   string `form:"text" json:"text" `                   //文本内容
		Content       string `form:"path" json:"path" `                   //path
		Covers        string `form:"covers" json:"covers"`                //图片路径
		ContentImages string `form:"content_Image" json:"content_Image" ` //缩略图路径
		Type          int    `form:"type" json:"type" binding:"required"` //文章类型
		Weight        int    `form:"weight" json:"weight" `
		Astatus       int    `form:"status" json:"status" binding:"required"`
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		this.RespErr(c, err)
		return
	}

	//获取作者 用户名
	//获取管理员ID
	//管理员 用户名
	//获取Cooke 用户登录ID
	//截取文件名作为objecetde 对象
	if len(req.Content) > 0 {
		path, err := new(models.Article).LocalFileToAliCloud(req.Covers)
		if err != nil {
			this.RespErr(c, err)
			return
		}
		req.ContentImages = path
	}

	if len(req.ContentImages) > 0 {
		path, err := new(models.Article).LocalFileToAliCloud(req.ContentImages)
		if err != nil {
			this.RespErr(c, err)
			return
		}
		req.ContentImages = path
	}

	//
	//
	if len(req.ContentText) > 0 {
		req.Content = req.ContentText
	}
	err = new(models.Article).AddArticle(&models.Article{
		Id:            req.Id,
		Title:         req.Title,
		Description:   req.Description,
		Content:       req.Content,
		Covers:        req.Covers,
		ContentImages: req.ContentImages,
		Type:          req.Type,
		Weight:        req.Weight,
		Astatus:       req.Astatus,
	})
	if err != nil {
		this.RespErr(c, err)
		return
	}
	this.RespOK(c, "成功")
	return
}
