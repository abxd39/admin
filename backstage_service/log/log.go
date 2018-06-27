package log

import (
	cf "admin/backstage_service/conf"
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLog() {
	Log = logrus.New()
	filename := cf.Inif.MustValue("log", "log_path")

	_, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		//Log.Out = file
		Log.Out = os.Stdout
	} else {
		Log.Out = os.Stdout
		//panic("Failed to log to file, using default stderr")
	}
}
