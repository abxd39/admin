package utils

import (
	"github.com/Unknwon/goconfig"
)

var Cfg *goconfig.ConfigFile

func init() {
	var err error
	Cfg, err = goconfig.LoadConfigFile("conf/admin.ini")
	if err != nil {
		panic("load config err is " + err.Error())
	}
}
