package controller

import (
	"admin/app/models"
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
		group.GET("/get_sms", c.GetSms)
		group.POST("/set_sms", c.SetSms)
		group.GET("/get_kefu", c.GetKefu)
		group.POST("/set_kefu", c.SetKefu)
	}
}

// 获取基础配置
func (c *ConfigController) GetSite(ctx *gin.Context) {
	// 调用model
	configMD := new(models.Config)
	config, err := configMD.Get(models.CONFIG_SITE)
	if err != nil {
		c.RespErr(ctx, err)
		return
	}

	// json反解析
	value := models.SiteConfig{}
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
	valueBytes, err := json.Marshal(models.SiteConfig{
		Name:            name,
		EnglishName:     englishName,
		Title:           title,
		EnglishTitle:    englishTitle,
		Logo:            logo,
		Keyword:         keyword,
		EnglishKeyword:  englishKeyword,
		Desc:            desc,
		EnglishDesc:     englishDesc,
		Beian:           beian,
		StatisticScript: statisticScript,
	})
	if err != nil {
		c.RespErr(ctx, errors.NewSys(err))
		return
	}

	// 组织参数
	config := &models.Config{
		Name:  models.CONFIG_SITE,
		Value: string(valueBytes),
	}

	// 调用model
	err = new(models.Config).Set(config)
	if err != nil {
		c.RespErr(ctx, err)
		return
	}

	// 返回
	c.RespOK(ctx)
	return
}

// 获取短信配置
func (c *ConfigController) GetSms(ctx *gin.Context) {
	// 调用model
	configMD := new(models.Config)
	config, err := configMD.Get(models.CONFIG_SMS)
	if err != nil {
		c.RespErr(ctx, err)
		return
	}

	// json反解析
	value := models.SmsConfig{}
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

// 设置短信配置
func (c *ConfigController) SetSms(ctx *gin.Context) {
	// 获取参数
	url := c.GetString(ctx, "url")
	if strLen := utf8.RuneCountInString(url); strLen < 1 {
		c.RespErr(ctx, "参数url格式错误")
		return
	}

	userName := c.GetString(ctx, "user_name")
	if strLen := utf8.RuneCountInString(userName); strLen < 1 {
		c.RespErr(ctx, "参数user_name格式错误")
		return
	}

	password := c.GetString(ctx, "password")
	if strLen := utf8.RuneCountInString(password); strLen < 1 {
		c.RespErr(ctx, "参数password格式错误")
		return
	}

	// 序列化成json
	valueBytes, err := json.Marshal(models.SmsConfig{
		Url:      url,
		UserName: userName,
		Password: password,
	})
	if err != nil {
		c.RespErr(ctx, errors.NewSys(err))
		return
	}

	// 组织参数
	config := &models.Config{
		Name:  models.CONFIG_SMS,
		Value: string(valueBytes),
	}

	// 调用model
	err = new(models.Config).Set(config)
	if err != nil {
		c.RespErr(ctx, err)
		return
	}

	// 返回
	c.RespOK(ctx)
	return
}

// 获取客服配置
func (c *ConfigController) GetKefu(ctx *gin.Context) {
	// 调用model
	configMD := new(models.Config)
	config, err := configMD.Get(models.CONFIG_KEFU)
	if err != nil {
		c.RespErr(ctx, err)
		return
	}

	// json反解析
	value := models.KefuConfig{}
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

// 设置短信配置
func (c *ConfigController) SetKefu(ctx *gin.Context) {
	// 获取参数
	phone := c.GetString(ctx, "phone")
	if strLen := utf8.RuneCountInString(phone); strLen < 1 {
		c.RespErr(ctx, "参数url格式错误")
		return
	}

	email := c.GetString(ctx, "email")
	if strLen := utf8.RuneCountInString(email); strLen < 1 {
		c.RespErr(ctx, "参数email格式错误")
		return
	}

	address := c.GetString(ctx, "address")
	if strLen := utf8.RuneCountInString(address); strLen < 1 {
		c.RespErr(ctx, "参数address格式错误")
		return
	}

	// 序列化成json
	valueBytes, err := json.Marshal(models.KefuConfig{
		Phone:   phone,
		Email:   email,
		Address: address,
	})
	if err != nil {
		c.RespErr(ctx, errors.NewSys(err))
		return
	}

	// 组织参数
	config := &models.Config{
		Name:  models.CONFIG_KEFU,
		Value: string(valueBytes),
	}

	// 调用model
	err = new(models.Config).Set(config)
	if err != nil {
		c.RespErr(ctx, err)
		return
	}

	// 返回
	c.RespOK(ctx)
	return
}
