package utils

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

var Store sessions.Store

func init() {
	stype, _ := Cfg.GetValue("session", "type")

	switch stype {
	case "redis":
		addr, err := Cfg.GetValue("session", "addr")
		if err != nil {
			panic("utils.session:" + err.Error())
		}
		Store, err = redis.NewStore(10, "tcp", addr, "ailaiduokeji657@@@", []byte("secret"))
		if err != nil {
			panic("utils.session:" + err.Error())
		}
	default:
		panic("unkown session.type")
	}

}

/* 使用示例
r.GET("/incr", func(c *gin.Context) {
	session := sessions.Default(c)
	var count int
	v := session.Get("count")
	if v == nil {
		count = 0
	} else {
		count = v.(int)
		count++
	}
	session.Set("count", count)
	session.Save()
	c.JSON(200, gin.H{"count": count})
})

*/
