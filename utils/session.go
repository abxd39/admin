package utils

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

var Store sessions.Store

func init() {
	stype, err := Cfg.GetValue("session", "type")
	if err != nil {
		panic("utils.session:" + err.Error())
	}

	switch stype {
	case "redis":
		addr, err := Cfg.GetValue("session", "addr")
		if err != nil {
			panic("utils.session:" + err.Error())
		}
		pwd, err := Cfg.GetValue("session", "pwd")
		if err != nil {
			panic("utils.session:" + err.Error())
		}

		Store, err = redis.NewStore(10, "tcp", addr, pwd, []byte("secret"))
		if err != nil {
			panic("utils.session:" + err.Error())
		}
	default:
		panic("unkown session.type")
	}

}

// 当前登录的管理员的id
func GetUid(ctx *gin.Context) int {
	session := sessions.Default(ctx)
	return session.Get("uid").(int)
}

// 当前登录管理员是否超管
func IsSuper(ctx *gin.Context) bool {
	session := sessions.Default(ctx)
	return session.Get("is_super").(bool)
}
