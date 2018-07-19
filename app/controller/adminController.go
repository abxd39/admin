package controller

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"

	bk "admin/app/models/backstage"
	"admin/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	google "github.com/mojocn/base64Captcha"
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
		group.GET("/list_login_log", a.ListLoginLog)
		group.GET("/my_left_menu", a.MyLeftMenu)
		group.GET("/my_right_menu", a.MyRightMenu)
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
	//查数据库 验证用户名和密码
	isSuper, nickName, uid, err := new(bk.User).Login(req.Phone, req.LoginPwd)
	if err != nil {
		// @@@写入登录日志，登录失败@@@
		if _, err := new(bk.UserLoginLog).Add(&bk.UserLoginLog{
			Uid:       uid,
			NickName:  req.Phone,
			States:    2,
			LoginIp:   utils.GetRemoteAddr(ctx),
			LoginTime: time.Now().Unix(),
		}); err != nil {
			utils.AdminLog.Error("记录管理员登录日志失败", err)
		}

		// 返回
		ctx.JSON(http.StatusOK, gin.H{"code": 1, "data": "", "msg": "登录失败"})
		return
	}

	//success
	//fmt.Println("verify success")
	//添加cooke 用户名

	session.Set("uid", uid)
	session.Set("name", req.Phone)
	session.Set("is_super", isSuper)
	session.Save()

	// @@@写入登录日志，登录成功@@@
	if _, err := new(bk.UserLoginLog).Add(&bk.UserLoginLog{
		Uid:       uid,
		NickName:  nickName,
		States:    1,
		LoginIp:   utils.GetRemoteAddr(ctx),
		LoginTime: time.Now().Unix(),
	}); err != nil {
		utils.AdminLog.Error("记录管理员登录日志失败", err)
	}

	// 返回
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "data": "", "uid": uid, "name": nickName, "msg": "登录成功"})
	return
}

// 登出
func (a *AdminController) Logout(ctx *gin.Context) {
	// 清空session
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()

	// 返回
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

	// 获取绑定的节点ID
	roleIds, err := userMD.GetBindRoleIds(uid)
	if err != nil {
		a.RespErr(ctx, err)
		return
	}

	// 设置返回数据
	a.Put(ctx, "user", user)
	a.Put(ctx, "role_ids", roleIds)

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
	params := make(map[string]interface{})

	// 获取参数
	uid, err := a.GetInt(ctx, "uid")
	if err != nil || uid < 1 {
		a.RespErr(ctx, "参数uid格式错误")
		return
	}

	if name, ok := a.GetParam(ctx, "name"); ok {
		if matched, err := regexp.MatchString(`^[a-z0-9_\-]{3,15}$`, name); err != nil || !matched {
			a.RespErr(ctx, "参数name格式错误，3-15个字符")
			return
		}

		params["name"] = name
	}

	if nickname, ok := a.GetParam(ctx, "nickname"); ok {
		if strLen := utf8.RuneCountInString(nickname); strLen < 2 || strLen > 15 {
			a.RespErr(ctx, "参数nickname格式错误，2-15个字符")
			return
		}

		params["nick_name"] = nickname
	}

	if pwd, ok := a.GetParam(ctx, "pwd"); ok {
		if matched, err := regexp.MatchString(`^[a-zA-Z0-9~!@#$%^&*_\-=+:;|,.?]{6,20}$`, pwd); err != nil || !matched {
			a.RespErr(ctx, "密码格式错误，6-20个字符")
			return
		}

		rePwd, _ := a.GetParam(ctx, "re_pwd")
		if rePwd != pwd {
			a.RespErr(ctx, "两次输入的密码不一致")
			return
		}

		params["pwd"] = pwd
	}

	if roleIds, ok := a.GetParam(ctx, "role_ids"); ok { // 逗号分隔
		params["role_ids"] = roleIds
	}

	// 调用model
	err = (new(bk.User)).Update(uid, params)
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

// 删除管理员
func (a *AdminController) Delete(ctx *gin.Context) {
	// 获取参数
	uid, err := a.GetInt(ctx, "uid")
	if err != nil || uid < 1 {
		a.RespErr(ctx, "参数uid格式错误")
		return
	}

	// 调用model
	err = new(bk.User).Delete(uid)
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

// 管理员登录日志列表
func (a *AdminController) ListLoginLog(ctx *gin.Context) {
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

	// 筛选参数
	filter := make(map[string]string)
	if v, ok := a.GetParam(ctx, "login_date"); ok {
		filter["login_date"] = v
	}

	// 调用model
	list, err := new(bk.UserLoginLog).List(page, rows, filter)
	if err != nil {
		a.RespErr(ctx, err)
		return
	}

	// 重新组装数据
	if list.Total > 0 {
		if items, ok := list.Items.([]bk.UserLoginLog); ok {
			type NewItem struct {
				Id        int    `json:"id"`
				Uid       int    `json:"uid"`
				NickName  string `json:"nick_name"`
				LoginIp   string `json:"login_ip"`
				LoginTime int64  `json:"login_time"`
				Desc      string `json:"desc"`
			}

			newItems := make([]NewItem, len(items))
			for k, v := range items {
				state := "成功"
				if v.States == 2 {
					state = "失败"
				}

				newItems[k] = NewItem{
					Id:        v.Id,
					Uid:       v.Uid,
					NickName:  v.NickName,
					LoginIp:   v.LoginIp,
					LoginTime: v.LoginTime,
					Desc:      fmt.Sprintf("%s在%s时，登录%s。", v.NickName, utils.Unix2Date(v.LoginTime), state),
				}
			}

			// 用新的items
			list.Items = newItems
		} else {
			// 类型断言错误
			a.RespErr(ctx, "类型错误")
			return
		}
	}

	// 设置返回数据
	a.Put(ctx, "list", list)

	// 返回
	a.RespOK(ctx)
	return
}

// 获取左侧菜单
func (a *AdminController) MyLeftMenu(ctx *gin.Context) {
	// 调用model
	list, err := new(bk.User).MyLeftMenu(ctx)
	if err != nil {
		a.RespErr(ctx, err)
		return
	}

	// 重新组装数据
	type NewItem struct {
		Id       int    `json:"id"`
		Pid      int    `json:"pid"`
		Weight   int    `json:"weight"`
		Title    string `json:"title"`
		Depth    int    `json:"depth"`
		MenuUrl  string `json:"menu_url"`
		MenuIcon string `json:"menu_icon"`
		FullId   string `json:"full_id"`
	}

	newList := make([]NewItem, len(list))
	for k, v := range list {
		newList[k] = NewItem{
			Id:       v.Id,
			Pid:      v.Pid,
			Weight:   v.Weight,
			Title:    v.Title,
			Depth:    v.Depth,
			MenuUrl:  v.MenuUrl,
			MenuIcon: v.MenuIcon,
			FullId:   v.FullId,
		}
	}

	// 设置返回数据
	a.Put(ctx, "list", newList)

	// 返回
	a.RespOK(ctx)
	return
}

// 获取右侧菜单
func (a *AdminController) MyRightMenu(ctx *gin.Context) {
	// 获取参数
	pid, err := a.GetInt(ctx, "pid")
	if err != nil {
		a.RespErr(ctx, "参数pid格式错误")
		return
	}

	// 调用model
	list, err := new(bk.User).MyRightMenu(ctx, pid)
	if err != nil {
		a.RespErr(ctx, err)
		return
	}

	// 重新组装数据
	type NewItem struct {
		Id       int    `json:"id"`
		Pid      int    `json:"pid"`
		Weight   int    `json:"weight"`
		Title    string `json:"title"`
		Depth    int    `json:"depth"`
		MenuUrl  string `json:"menu_url"`
		MenuIcon string `json:"menu_icon"`
		FullId   string `json:"full_id"`
	}

	newList := make([]NewItem, len(list))
	for k, v := range list {
		newList[k] = NewItem{
			Id:       v.Id,
			Pid:      v.Pid,
			Weight:   v.Weight,
			Title:    v.Title,
			Depth:    v.Depth,
			MenuUrl:  v.MenuUrl,
			MenuIcon: v.MenuIcon,
			FullId:   v.FullId,
		}
	}

	// 设置返回数据
	a.Put(ctx, "list", newList)

	// 返回
	a.RespOK(ctx)
	return
}
