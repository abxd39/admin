package controller

import (
	"admin/app/models"
	"admin/log"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContextController struct{}

func (cm *ContextController) Router(r *gin.Engine) {
	g := r.Group("/content")
	{
		g.POST("/addlink", cm.AddFriendlyLink)
		g.GET("/linklist", cm.GetFriendlyLink)
		g.GET("/article", cm.GetArticleList)
	}
}

func (cm *ContextController) AddFriendlyLink(c *gin.Context) {
	fmt.Println("..........................................")
	req := struct {
		WebName   string `form:"web_name" json:"web_name" binding:"required"`
		LinkName  string `form:"link_name" json:"link_name" binding:"required"`
		Aorder    int    `form:"order" json:"order" binding:"required"`
		LinkState int    `form:"link_state" json:"link_state" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.AdminLog.Errorf(err.Error())
		return
	}
	err = new(models.FriendlyLink).Add(req.Aorder, req.LinkState, req.WebName, req.LinkName)
	if err != nil {
		log.AdminLog.Errorf(err.Error())
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "", "msg": "成功"})
	return
}

func (cm *ContextController) GetFriendlyLink(c *gin.Context) {

	req := struct {
		Page  int `form:"page" json:"page" binding:"required"`
		Count int `form:"count" json:"count" `
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.AdminLog.Errorf(err.Error())
		return
	}
	//operator db         GetFriendlyLinkList
	result, err := new(models.FriendlyLink).GetFriendlyLinkList(req.Count, req.Page)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": result, "msg": "成功"})
	return
}

func (cm *ContextController) AddBanner(c *gin.Context) {
	req := struct {
		Order       int    `form:"order" json:"order" binding:"required"`
		PictureName string `form:"picture_n" json:"picture_n" binding:"required"`
		PicturePath string `form:"picture_p" json:"picture_p" binding:"required"`
		LinkAddr    string `form:"link_addr" json:"link_addr" binding:"required"`
		Start_t     string `form:"start_t" json:"start_t" binding:"required"`
		End_t       string `form:"end_t" json:"end_t" binding:"required"`
		State       int    `form:"state" json:"state" binding:"required"`
	}{}

	err := c.ShouldBind(&req)
	if err != nil {
		log.AdminLog.Errorf(err.Error())
		return
	}
	err = new(models.Banner).Add(req.Order, req.State, req.PictureName, req.PicturePath, req.LinkAddr, req.Start_t, req.End_t)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "", "msg": "成功"})
	return
}

func (cm *ContextController) GetArticleList(c *gin.Context) {
	req := struct {
		Page int `form:"page" json:"page" binding:"required"`
		Rows int `form:"rows" json:"rows" binding:"required"`
		Type int `form:"type" json:"type" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.AdminLog.Errorf(err.Error())
		return
	}
	reuslt, total, er := new(models.ArticleList).GetArticleList(req.Page, req.Rows, req.Type)
	if er != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": er.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": reuslt, "total": total, "msg": "成功"})
}

func (cm *ContextController) AddArticle(c *gin.Context) {
	req := struct {
		Title         string `form:"title" json:"title" binding:"required"`
		Description   string `form:"desc" json:"desc" `
		Content       string `form:"content" json:"content" binding:"required"`
		Covers        string `form:"covers" json:"covers"`
		ContentImages string `form:"content_Image" json:"content_Image" `
		Type          int    `form:"tpye" json:"type" binding:"required"`
		Author        string `form:"author" json:"author" binding:"required"`
		Weight        int    `form:"weight" json:"weight" binding:"required"`
		Astatus       int    `form:"status" json:"status" binding:"required"`
		//AdminId       int    `form:"admin_id" json:"admin_id" binding:"required"`
		//AdminNickname string `form:"admin_name" json:"admin_name" binding:"required"`
	}{}
	err := c.ShouldBind(&req)
	if err != nil {
		log.AdminLog.Errorf(err.Error())
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": "", "msg": "成功"})
	return
}
