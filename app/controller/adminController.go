package controller

import (
	"fmt"
	"log"
	"net/http"

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
		UserName string `form:"uname" json:"uname" binding:"required"`
		LoginPwd string `form:"pwd" json:"pwd" binding:"required"`
		Idkey    string `form:"idkey" json:"idkey" binding:"required"`
		Verify   string `form:"verify" json:"verify" binding:"required"`
	}{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		log.Fatalln("param buind failed !!")
		return
	}
	//verify code
	verifyResult := verifyCaptcha(req.Idkey, req.Verify)
	if verifyResult {
		//success
		fmt.Println("verify success")
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": "url", "msg": "登录成功"})

	} else {
		//failed
		fmt.Println("verify failed!!")
		ctx.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": "登录失败"})
		return
	}
	return
}

func (this *AdminController) Loginout(ctx *gin.Context) {

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
