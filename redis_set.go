package redis

// SAdd key member [member ...]
func (r *Redis) SAdd(key string, member ...string) *Reply {
	return r.run(append([]string{"SADD", key}, member...)...)
}
