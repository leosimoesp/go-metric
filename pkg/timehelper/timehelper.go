package timehelper

import "time"

func CreateTimestamp() int64 {
	return time.Now().Unix()
}

func FromTimeToInt(t time.Time) int64 {
	return t.Unix()
}

func GetLastHourTimestamp() int64 {
	now := time.Now()

	lastHour := now.Add(time.Duration(-1) * time.Hour)

	return FromTimeToInt(lastHour)
}
