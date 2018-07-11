package controller

import (
	"admin/app/models"
	"admin/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContextController struct{}

func (this *ContextController) Router(r *gin.Engine) {
	g := r.Group("/content")
	{
		g.POST("/add_link", this.AddFriendlyLink)
		g.GET("/link_list", this.GetFriendlyLink)
		g.GET("/article_list", this.GetArticleList)
		g.POST("/add_banner", this.AddBanner)
		g.GET("/banner_list", this.GetBannerList)
		g.POST("/add_article", this.AddArticle)
		g.GET("/article_type", this.GetArticleType)
		g.POST("/local_filetoali", this.LocalFileToAliCloud)
		g.GET("/get_filetoali", this.GetLocalFileToAliCloud)
	}
}

func (this *ContextController) GetLocalFileToAliCloud(c *gin.Context) {
	req := struct {
		FilePath  string `form:"file_path" json:"file_path" binding:"required"`
		ObjectKey string `form:"okey" json:"okey" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err})
		return
	}
	filepath, err := new(models.Article).GetLocalFileToAliCloud(req.ObjectKey, req.FilePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "", "path": filepath, "msg": "成功"})
	return
}

func (this *ContextController) LocalFileToAliCloud(c *gin.Context) {
	req := struct {
		FilePath string `form:"file_path" json:"file_path" binding:"required"`
		//ObjectKey string `form:"okey" json:"okey" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err})
		return
	}
	_, err = new(models.Article).LocalFileToAliCloud(req.FilePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	response := new(models.PolicyToken).Get_policy_token()
	// c.Request.Header.Set("Access-Control-Allow-Methods", "POST")
	// c.Request.Header.Set("Access-Control-Allow-Origin", "*")
	//io.WriteString(c, response)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": response, "msg": "成功"})
	//c.JSON(http.StatusOK, gin.H{"code": 0, "data": "", "msg": "成功"})
}

func (this *ContextController) GetArticleType(c *gin.Context) {
	result, err := new(models.ArticleType).GetArticleType()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": result, "msg": "成功"})
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
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}
	fmt.Println("..........................................")
	err = new(models.FriendlyLink).Add(1, req.LinkState, req.WebName, req.LinkName)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "", "msg": "成功"})
	return
}

func (this *ContextController) GetFriendlyLink(c *gin.Context) {

	req := struct {
		Page     int    `form:"page" json:"page" binding:"required"`
		Name     string `form:"name" json:"name" `
		Linkname string `form:"link_name" json:"link_name" `
		Count    int    `form:"count" json:"count" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		return
	}
	//operator db         GetFriendlyLinkList
	result, page, total, err := new(models.FriendlyLink).GetFriendlyLinkList(req.Page, req.Count, req.Name, req.Linkname)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "page": page, "total": total, "data": result, "msg": "成功"})
	return
}

func (this *ContextController) AddBanner(c *gin.Context) {
	req := struct {
		Order       int    `form:"order" json:"order" binding:"required"`
		PictureName string `form:"picture_n" json:"picture_n" binding:"required"`
		PicturePath string `form:"picture_p" json:"picture_p" binding:"required"`
		LinkAddr    string `form:"link_addr" json:"link_addr" binding:"required"`
		Status      int    `form:"status" json:"status" binding:"required"`
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}
	err = new(models.Banner).Add(req.Order, req.Status, req.PictureName, req.PicturePath, req.LinkAddr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "", "msg": "成功"})
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
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}
	list, page, total, err := new(models.Banner).GetBannerList(req.Page, req.Rows, req.Status, req.Start_t, req.End_t)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "page": page, "total": total, "data": list, "msg": "成功"})

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
		c.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err.Error()})
		return
	}
	reuslt, page, total, err := new(models.ArticleList).GetArticleList(req.Page, req.Rows, req.Type)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": reuslt, "page": page, "total": total, "msg": "成功"})
}

func (this *ContextController) AddArticle(c *gin.Context) {
	req := struct {
		Title         string `form:"title" json:"title" binding:"required"`
		Description   string `form:"desc" json:"desc" `
		Content       string `form:"path" json:"path" binding:"required"` //path
		Covers        string `form:"covers" json:"covers"`                //图片路径
		ContentImages string `form:"content_Image" json:"content_Image" ` //缩略图路径
		Type          int    `form:"type" json:"type" binding:"required"` //文章类型
		Weight        int    `form:"weight" json:"weight" `
		Astatus       int    `form:"status" json:"status" binding:"required"`
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorf(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
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
			c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
			return
		}
		req.ContentImages = path
	}

	if len(req.ContentImages) > 0 {
		path, err := new(models.Article).LocalFileToAliCloud(req.ContentImages)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
			return
		}
		req.ContentImages = path
	}

	//
	//
	err = new(models.Article).AddArticle(&models.Article{
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
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "", "msg": "成功"})
	return
}
