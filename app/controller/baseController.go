package controller

import (
	"admin/app/models/backstage"
	"admin/constant"
	"admin/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strings"
)

type BaseController struct {
	Controller
}

// 验证登录、权限
func (b *BaseController) CheckLogin(ctx *gin.Context) error {
	// 1. 验证登录
	session := sessions.Default(ctx)
	uid := session.Get("uid")
	if uid == nil { // 不存在，未登录
		return errors.NewNormal(constant.RESPONSE_CODE_SESSION_INVALID)
	}

	// 2. 验证权限
	// 无需验证权限的api
	noNeedAuthAPIs := map[string]bool{
		"admin/logout": true,
	}

	// 获取api地址
	uri := ctx.Request.RequestURI                            // 得到：/admin/login?name=xxx&pwd=xxx
	api := strings.TrimLeft(strings.Split(uri, "?")[0], "/") // 得到：admin/login
	if _, ok := noNeedAuthAPIs[api]; !ok {                   // !ok
		uidInt, _ := uid.(int)
		result, err := new(backstage.User).CheckPermission(ctx, uidInt, api)
		if err != nil {
			b.RespErr(ctx, err)
			return err
		}
		if !result {
			return errors.NewNormal(constant.RESPONSE_CODE_NO_API_PERMISSION)
		}
	}

	return nil
}

func (b *BaseController) GetUid() int {
	return 0
}
