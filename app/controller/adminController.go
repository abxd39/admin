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
	"regexp"
	"unicode/utf8"
)

// 管理员
type AdminController struct {
	BaseController
}

func (a *AdminController) Router(e *gin.Engine) {
	group := e.Group("/admin")
	{
		group.GET("/code", a.Code)
		group.POST("/login", a.Login)
		group.GET("/logout", a.Logout)
		group.GET("/list", a.List)
		group.GET("/get", a.Get)
		group.POST("/add", a.Add)
		group.POST("/update", a.Update)
		group.POST("/delete", a.Delete)

	}
}

// 登录验证码
func (a *AdminController) Code(ctx *gin.Context) {
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

// 登录
func (a *AdminController) Login(ctx *gin.Context) {
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

// 登出
func (a *AdminController) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": "", "msg": "成功"})
	return
}

// 管理员列表
func (a *AdminController) List(ctx *gin.Context) {
	// 获取参数
	page, err := a.GetInt(ctx, "page", 1)
	if err != nil {
		a.RespErr(ctx, "参数page格式错误")
		return
	}

	rows, err := a.GetInt(ctx, "rows", 10)
	if err != nil {
		a.RespErr(ctx, "参数rows格式错误")
		return
	}

	// 调用model
	list, err := new(bk.User).List(page, rows, nil)
	if err != nil {
		a.RespErr(ctx, err)
		return
	}

	// 设置返回数据
	a.Put(ctx, "list", list)

	// 返回
	a.RespOK(ctx)
	return

	return
}

// 管理员详情
func (a *AdminController) Get(ctx *gin.Context) {
	// 获取参数
	uid, err := a.GetInt(ctx, "uid")
	if err != nil || uid < 1 {
		a.RespErr(ctx, "参数uid格式错误")
		return
	}

	// 调用model
	userMD := new(bk.User)
	user, err := userMD.Get(uid)
	if err != nil {
		a.RespErr(ctx, err)
		return
	}

	// 设置返回数据
	a.Put(ctx, "user", user)

	// 返回
	a.RespOK(ctx)
	return
}

// 新增管理员
func (a *AdminController) Add(ctx *gin.Context) {
	// 获取参数
	name := a.GetString(ctx, "name")

	if matched, err := regexp.MatchString(`^[a-z0-9_\-]{3,15}$`, name); err != nil || !matched {
		a.RespErr(ctx, "参数name格式错误，只能为3-15位小写字母a-z、数字0-9、中划线-或下划线_")
		return
	}

	nickname := a.GetString(ctx, "nickname")
	if strLen := utf8.RuneCountInString(nickname); strLen < 2 || strLen > 15 {
		a.RespErr(ctx, "参数nickname格式错误，2-15个字符")
		return
	}

	pwd := a.GetString(ctx, "pwd")
	if matched, err := regexp.MatchString(`^[a-zA-Z0-9~!@#$%^&*_\-=+:;|,.?]{6,20}$`, pwd); err != nil || !matched {
		a.RespErr(ctx, "密码格式错误，6-20个字符")
		return
	}

	rePwd := a.GetString(ctx, "re_pwd")
	if rePwd != pwd {
		a.RespErr(ctx, "两次输入的密码不一致")
		return
	}

	roleIds := a.GetString(ctx, "role_ids") // 逗号分隔

	// 调用model
	user := &bk.User{
		Name:     name,
		NickName: nickname,
		Pwd:      pwd,
	}

	uid, err := (new(bk.User)).Add(user, roleIds)
	if err != nil {
		a.RespErr(ctx, err)
		return
	}

	// 设置返回数据
	a.Put(ctx, "uid", uid)

	// 返回
	a.RespOK(ctx)
	return
}

// 更新管理员
func (a *AdminController) Update(ctx *gin.Context) {

	return
}

// 删除管理员
func (a *AdminController) Delete(ctx *gin.Context) {

	return
}
