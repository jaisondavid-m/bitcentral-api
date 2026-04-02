package utils

import "time"

func TsToString(ts int64) string {
	if ts == 0 {
		return ""
	}
	t := time.Unix(ts/1000, (ts%1000)*int64(time.Millisecond))
	return t.Format(time.RFC3339)
}