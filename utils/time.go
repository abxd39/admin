package utils

import (
	"time"
)

const (
	LAYOUT_DATE         = "2006-01-02"
	LAYOUT_DATE_TIME    = "2006-01-02 15:04:05"
	LAYOUT_DATE_TIME_12 = "2006-01-02 03:04:05"
)

// 日期转时间戳
func Date2Unix(date string, layout string) int64 {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(layout, date, loc)
	return theTime.Unix()
}

// 时间戳转日期
func Unix2Date(unix int64, layout string) string {
	return time.Unix(unix, 0).Format(layout)
}
