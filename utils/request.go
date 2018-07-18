package utils

import (
	"github.com/gin-gonic/gin"
)

// 获取用户IP
// 如果通过nginx代理，需要设置nginx代理转发用户ip给go服务：proxy_set_header   X-Real-IP        $remote_addr;
// 如果没代理，直接通过Request获取
func GetRemoteAddr(ctx *gin.Context) string {
	ip := ctx.Request.Header.Get("X-Real-IP") // nginx转发ip
	if len(ip) == 0 {                         // 没有代理，直接去request取
		ip = ctx.Request.RemoteAddr
	}

	return ip
}
