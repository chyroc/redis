package redis

import "strconv"

// Select index
func (r *Redis) Select(index int) *Reply {
	return r.run("SELECT", strconv.Itoa(index))
}
