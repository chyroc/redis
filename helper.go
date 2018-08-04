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

var strapp *stringAppend

func init() {
	strapp = &stringAppend{}
}

type stringAppend struct {
	s []string
}

func (r *stringAppend) reset() *stringAppend {
	r.s = nil
	return r
}

func (r *stringAppend) adds(s ...string) *stringAppend {
	r.s = append(r.s, s...)
	return r
}

func (r *stringAppend) build() []string {
	return r.s
}
