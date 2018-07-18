package utils

import (
	"time"
)

// 日期转时间戳
func Date2Unix(date string) int64 {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", date, loc)
	return theTime.Unix()
}

// 时间戳转日期
func Unix2Date(unix int64) string {
	return time.Unix(unix, 0).Format("2006-01-02 15:04:05")
}
