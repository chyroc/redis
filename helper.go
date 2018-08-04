package redis

import (
	"strconv"
	"time"
)

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func durationToMillisecond(t time.Duration) string {
	return strconv.Itoa(int(t / time.Millisecond))
}
