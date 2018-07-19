package controller

import (
	"admin/app/models/backstage"
	"admin/errors"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"unicode/utf8"
)

// 后台配置
type ConfigController struct {
	BaseController
}

func (c *ConfigController) Router(e *gin.Engine) {
	group := e.Group("/config")
	{
		group.GET("/get_site", c.GetSite)
		group.POST("/set_site", c.SetSite)
	}
}

// 获取基础配置
func (c *ConfigController) GetSite(ctx *gin.Context) {
	// 调用model
	configMD := new(backstage.Config)
	config, err := configMD.Get(backstage.CONFIG_SITE)
	if err != nil {
		c.RespErr(ctx, err)
		return
	}

	// json反解析
	value := make(map[string]interface{})
	err = json.Unmarshal([]byte(config.Value), &value)
	if err != nil {
		c.RespErr(ctx, errors.NewSys(err))
		return
	}

	// 设置返回数据
	c.Put(ctx, "config", value)

	// 返回
	c.RespOK(ctx)
	return
}

// 设置基础配置
func (c *ConfigController) SetSite(ctx *gin.Context) {
	// 获取参数
	name := c.GetString(ctx, "name")
	if strLen := utf8.RuneCountInString(name); strLen < 1 {
		c.RespErr(ctx, "参数name格式错误")
		return
	}

	englishName := c.GetString(ctx, "english_name")
	if strLen := len(englishName); strLen < 1 {
		c.RespErr(ctx, "参数english_name格式错误")
		return
	}

	title := c.GetString(ctx, "title")
	if strLen := utf8.RuneCountInString(title); strLen < 1 {
		c.RespErr(ctx, "参数title格式错误")
		return
	}

	englishTitle := c.GetString(ctx, "english_title")
	if strLen := len(englishTitle); strLen < 1 {
		c.RespErr(ctx, "参数english_title格式错误")
		return
	}

	logo := c.GetString(ctx, "logo")
	if strLen := len(logo); strLen < 1 {
		c.RespErr(ctx, "参数logo格式错误")
		return
	}

	keyword := c.GetString(ctx, "keyword")
	if strLen := utf8.RuneCountInString(keyword); strLen < 1 {
		c.RespErr(ctx, "参数keyword格式错误")
		return
	}

	englishKeyword := c.GetString(ctx, "english_keyword")
	if strLen := len(englishKeyword); strLen < 1 {
		c.RespErr(ctx, "参数english_keyword格式错误")
		return
	}

	desc := c.GetString(ctx, "desc")
	if strLen := utf8.RuneCountInString(desc); strLen < 1 {
		c.RespErr(ctx, "参数desc格式错误")
		return
	}

	englishDesc := c.GetString(ctx, "english_desc")
	if strLen := len(englishDesc); strLen < 1 {
		c.RespErr(ctx, "参数english_desc格式错误")
		return
	}

	beian := c.GetString(ctx, "beian")
	if strLen := len(beian); strLen < 1 {
		c.RespErr(ctx, "参数beian格式错误")
		return
	}

	statisticScript := c.GetString(ctx, "statistic_script")
	if strLen := len(statisticScript); strLen < 1 {
		c.RespErr(ctx, "参数statistic_script格式错误")
		return
	}

	// 序列化成json
	valueInterface, err := json.Marshal(map[string]interface{}{
		"name":             name,
		"english_name":     englishName,
		"title":            title,
		"english_title":    englishTitle,
		"logo":             logo,
		"keyword":          keyword,
		"english_keyword":  englishKeyword,
		"desc":             desc,
		"english_desc":     englishDesc,
		"beian":            beian,
		"statistic_script": statisticScript,
	})
	if err != nil {
		c.RespErr(ctx, errors.NewSys(err))
		return
	}

	// 组织参数
	config := &backstage.Config{
		Name:  backstage.CONFIG_SITE,
		Value: string(valueInterface),
	}

	// 调用model
	err = new(backstage.Config).Set(config)
	if err != nil {
		c.RespErr(ctx, err)
		return
	}

	// 返回
	c.RespOK(ctx)
	return
}
