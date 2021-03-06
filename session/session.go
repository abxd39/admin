package session

import (
	"admin/constant"
	"admin/errors"
	"admin/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

var Store sessions.Store

func init() {
	stype, err := utils.Cfg.GetValue("session", "type")
	if err != nil {
		panic("utils.session:" + err.Error())
	}

	switch stype {
	case "redis":
		addr, err := utils.Cfg.GetValue("redis", "addr")
		if err != nil {
			panic("utils.session:" + err.Error())
		}
		pwd, err := utils.Cfg.GetValue("redis", "pwd")
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

	// 设置属性
	//Store.Options(sessions.Options{MaxAge: 10 * 60}) // 超时时间，秒
}

// 当前登录的管理员的id
func GetUid(ctx *gin.Context) (int, error) {
	session := sessions.Default(ctx)
	uidInterface := session.Get("uid")
	if uidInterface == nil {
		return 0, errors.NewNormal(constant.RESPONSE_CODE_SESSION_INVALID)
	}

	return uidInterface.(int), nil
}

// 当前登录管理员是否超管
func IsSuper(ctx *gin.Context) (bool, error) {
	session := sessions.Default(ctx)
	isSuperInterface := session.Get("is_super")
	if isSuperInterface == nil {
		return false, errors.NewNormal(constant.RESPONSE_CODE_SESSION_INVALID)
	}

	return isSuperInterface.(bool), nil
}
