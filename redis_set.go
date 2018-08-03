package redis

// SAdd key member [member ...]
func (r *Redis) SAdd(key string, member ...string) (int, error) {
	return r.run(append([]string{"SADD", key}, member...)...).int()
}
