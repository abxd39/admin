package log

import (
	cf "admin/conf"
	"os"

	"github.com/liudng/godump"
	"github.com/sirupsen/logrus"
)

var AdminLog *logrus.Logger

func InitLog() {
	AdminLog = logrus.New()

	filename := cf.Cfg.MustValue("log", "log_path")
	godump.Dump(filename)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		AdminLog.Out = file
	} else {
		AdminLog.Out = os.Stdout
		//panic("Failed to log to file, using default stderr")
	}
}
