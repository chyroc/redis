package redis

import "strconv"

// Select index
func (r *Redis) Select(index int) error {
	return r.run("SELECT", strconv.Itoa(index)).err
}
