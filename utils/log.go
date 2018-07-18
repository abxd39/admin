package utils

import (
	"os"

	"github.com/liudng/godump"
	"github.com/sirupsen/logrus"
)

var AdminLog *logrus.Logger

func init() {
	AdminLog = logrus.New()

	filename := Cfg.MustValue("log", "path")
	godump.Dump(filename)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		AdminLog.Out = file
	} else {
		AdminLog.Out = os.Stdout
		//panic("Failed to log to file, using default stderr")
	}
}
