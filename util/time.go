package util

import "time"

const defaultTimeZone = "Asia/Jakarta"

// ToDefaultTimezone force time to default timezone
func ToDefaultTimezone(format string, t *time.Time) (result time.Time, err error) {
	loc, err := time.LoadLocation(defaultTimeZone)
	if err != nil {
		return result, err
	}
	temp := t.Format(format)
	tz, err := time.Parse(format, temp)
	if err != nil {
		return result, err
	}

	result = tz.In(loc)
	return
}

// TimeToString convert pointer time based on format
func TimeToString(t *time.Time, format string) string {
	if t == nil {
		return ""
	}
	return t.Format(format)
}
