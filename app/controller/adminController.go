package controller

import (
	"admin/app/models"
	"admin/utils"
	"crypto/md5"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

type AdminController struct {
}

func (this *AdminController) Router(r *gin.Engine) {
	group := r.Group("/admin")
	{
		group.GET("/code", this.Code)
		group.GET("/login", this.Login)
		group.GET("/loginout", this.Loginout)
		group.GET("/list", this.List)
		group.GET("/delete", this.Delete)
		group.GET("/update", this.Update)
	}
}

func (this *AdminController) Code(ctx *gin.Context) {
	fmt.Println("...............................................")
	var configD = base64Captcha.ConfigDigit{
		Height:     80,
		Width:      240,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 4,
	}

	//ctx.Request.AddCookie()
	idKeyD, capD := base64Captcha.GenerateCaptcha("", configD)
	cook := &http.Cookie{
		Name:  "idkey",
		Path:  "/",
		Value: idKeyD,
	}
	//以base64编码
	base64stringD := base64Captcha.CaptchaWriteToBase64Encoding(capD)
	ctx.Data(0, idKeyD, capD.BinaryEncodeing())
	ctx.Request.AddCookie(cook)
	fmt.Println(idKeyD, base64stringD)
	return
}

//verfiy

func verifyCaptcha(idkey, verifyValue string) bool {
	return base64Captcha.VerifyCaptcha(idkey, verifyValue)

}

func (this *AdminController) Login(ctx *gin.Context) {
	req := struct {
		Phone    string `form:"phone" json:"phone" binding:"required"`
		LoginPwd string `form:"pwd" json:"pwd" binding:"required"`
		Verify   string `form:"verify" json:"verify" binding:"required"`
	}{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		log.Fatalln("param buind failed !!")
		return
	}
	//verify code
	var idkey string
	verifyResult := verifyCaptcha(idkey, req.Verify)
	if !verifyResult {
		//failed
		fmt.Println("verify failed!!")
		ctx.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": "登录失败"})
		return
	}
	//md5加密
	//var hanlen int
	has := md5.New()
	_, err = has.Write([]byte(req.LoginPwd))
	if err != nil {
		utils.AdminLog.Errorln("md5 加密失败!!")
		ctx.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": "登录失败"})
		return
	}
	hasvalue := has.Sum(nil)
	//查数据库 验证用户名和密码
	var uid int
	uid, err = new(models.User).Login(string(hasvalue), req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": "登录失败"})
		return
	}

	//success
	fmt.Println("verify success")
	//添加cooke 用户名
	session := sessions.Default(ctx)
	session.Clear()
	session.Set("idkey", idkey)
	session.Set("uid", uid)
	session.Set("phone", req.Phone)
	session.Save()
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": "url", "msg": "登录成功"})
	return
}

func (this *AdminController) Loginout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	return
}
func (this *AdminController) List(ctx *gin.Context) {

	return
}
func (this *AdminController) Delete(ctx *gin.Context) {

	return
}
func (this *AdminController) Update(ctx *gin.Context) {

	return
}
