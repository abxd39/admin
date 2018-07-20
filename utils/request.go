package utils

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

// 获取用户IP
func GetRemoteAddr(ctx *gin.Context) string {
	remoteAddr, _, _ := net.SplitHostPort(ctx.Request.RemoteAddr)

	// 如果有代理，上面取的addr不准确
	if ip := ctx.Request.Header.Get("X-Forwarded-For"); ip != "" { // 最准确
		remoteAddr = strings.Split(ip, ",")[0] // 多层代理，取第一个
	} else if ip = ctx.Request.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}
