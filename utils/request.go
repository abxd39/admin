package utils

import (
	"net"

	"github.com/gin-gonic/gin"
)

// 获取用户IP
// 如果通过nginx代理，需要设置nginx代理转发用户ip给go服务：proxy_set_header   X-Real-IP        $remote_addr;
// 如果没代理，直接通过Request获取
func GetRemoteAddr(ctx *gin.Context) string {
	remoteAddr := ctx.Request.RemoteAddr
	if ip := ctx.Request.Header.Get("X-Real-IP"); ip != "" { // nginx转发ip
		remoteAddr = ip
	} else if ip = ctx.Request.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}
