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

func float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 10, 64)
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

func buildSlice1(s string, ss []string) []string {
	return append([]string{s}, ss...)
}

func buildSlice2(a, b string, ss []string) []string {
	return append([]string{a, b}, ss...)
}

func buildSlice3(a, b, c string, ss []string) []string {
	return append([]string{a, b, c}, ss...)
}
