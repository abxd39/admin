package middleware

import (
	"strings"

	"admin/app/controller"
	"admin/app/models/backstage"
	"admin/constant"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 验证登录、权限
func CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取api地址
		uri := ctx.Request.RequestURI                            // 得到：/admin/login?name=xxx&pwd=xxx
		api := strings.TrimLeft(strings.Split(uri, "?")[0], "/") // 得到：admin/login

		// 1. 验证登录
		// 无需登录的api
		noNeedLoginAPIs := map[string]bool{
			"admin/code":   true, // 值用不到
			"admin/login":  true,
			"admin/logout": true,
		}

		var uid int
		if _, ok := noNeedLoginAPIs[api]; !ok { // !ok
			session := sessions.Default(ctx)
			uidInterface := session.Get("uid")
			if uidInterface == nil { // 不存在，未登录
				new(controller.BaseController).RespErr(ctx, constant.RESPONSE_CODE_SESSION_INVALID)
				ctx.Abort()
				return
			}

			uid = uidInterface.(int)
		}

		// 2. 验证权限
		// 无需验证权限的api
		noNeedAuthAPIs := map[string]bool{
			"admin/my_left_menu":  true,
			"admin/my_right_menu": true,
		}

		// 合并无需登录的接口，无需登录的接口肯定也无需验证权限
		for k, v := range noNeedLoginAPIs {
			noNeedAuthAPIs[k] = v
		}

		if _, ok := noNeedAuthAPIs[api]; !ok { // !ok
			has, err := new(backstage.User).CheckPermission(ctx, uid, api)
			if err != nil {
				new(controller.BaseController).RespErr(ctx, err)
				ctx.Abort()
				return
			}
			if !has {
				new(controller.BaseController).RespErr(ctx, constant.RESPONSE_CODE_NO_API_PERMISSION)
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
