package controller

import (
	bk "admin/app/models/backstage"
	"admin/utils"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	google "github.com/mojocn/base64Captcha"
)

type AdminController struct {
}

func (this *AdminController) Router(r *gin.Engine) {
	group := r.Group("/admin")
	{
		group.GET("/code", this.Code)
		group.POST("/login", this.Login)
		group.GET("/logout", this.Logout)
		group.GET("/list", this.List)
		group.GET("/delete", this.Delete)
		group.GET("/update", this.Update)
	}
}

func (this *AdminController) Code(ctx *gin.Context) {
	fmt.Println("...............................................")
	var configD = google.ConfigDigit{
		Height:     40,
		Width:      120,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 4,
	}

	//ctx.Request.AddCookie()
	idKeyD, capD := google.GenerateCaptcha("", configD)
	// cook := &http.Cookie{
	// 	Name:  "idkey",
	// 	Path:  "/",
	// 	Value: idKeyD,
	// }
	//以base64编码
	base64stringD := google.CaptchaWriteToBase64Encoding(capD)
	//ctx.Request.AddCookie(cook)
	session := sessions.Default(ctx)
	session.Clear()
	session.Set("idkey", idKeyD)
	session.Save()
	ctx.Data(0, idKeyD, capD.BinaryEncodeing())
	fmt.Printf("idke=%s", idKeyD)
	fmt.Println(idKeyD, base64stringD)
	return
}

//verfiy

func verifyCaptcha(idkey, verifyValue string) bool {
	// return google.VerifyCaptcha(idkey, verifyValue)
	return false
}

func (this *AdminController) Login(ctx *gin.Context) {
	req := struct {
		Phone    string `form:"phone" json:"phone" binding:"required"`
		LoginPwd string `form:"pwd" json:"pwd" binding:"required"`
		//Verify   string `form:"verify" json:"verify" binding:"required"`
	}{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		utils.AdminLog.Errorln("param buind failed !!")
		ctx.JSON(http.StatusOK, gin.H{"code": 2, "data": "", "msg": err})
		return
	}
	//verify code
	fmt.Println("0000.0.0.0000.0.0.0.0.0.0.0.0.0.0.0000....0.0.00.0.0..0..")
	fmt.Println(req)
	session := sessions.Default(ctx) //因为每次获取验证码时都会清除一次。所以再次不应该调用clear函数
	// var idkey string
	// v := session.Get("idkey")
	// if v == nil {
	// 	utils.AdminLog.Errorln("sessions 中idkey 的值获取失败")
	// 	return
	// } else {
	// 	idkey = v.(string)
	// }
	// fmt.Println("idkey=", idkey, "verify=", req.Verify)
	// verifyResult := verifyCaptcha(idkey, req.Verify)
	// if !verifyResult {
	// 	//failed
	// 	fmt.Println("verify failed!!")
	// 	ctx.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": "验证错误"})
	// 	return
	// }
	//md5加密
	//var hanlen int
	//fmt.Println("vvvvvvvvvvvvvvvvvvvvv", req.LoginPwd)
	has := md5.New()
	has.Write([]byte(req.LoginPwd))
	hasvalue := has.Sum(nil)
	//查数据库 验证用户名和密码
	fmt.Println("hasvalue=", hex.EncodeToString(hasvalue))
	var uid int
	var name string
	name, uid, err = new(bk.User).Login(hex.EncodeToString(hasvalue), req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": "登录失败"})
		return
	}

	//success
	//fmt.Println("verify success")
	//添加cooke 用户名

	session.Set("uid", uid)
	session.Set("phone", req.Phone)
	session.Save()
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": "", "uid": uid, "name": name, "msg": "登录成功"})
	return
}

func (this *AdminController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": "", "msg": "成功"})
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
