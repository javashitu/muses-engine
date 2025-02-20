package util

import "time"

func CurMillionSeconds() int64 {
	now := time.Now()
	return now.UnixNano() / int64(time.Millisecond)
}

func CurDate() string {
	now := time.Now()
	return now.Format("2006-01-02")
}
