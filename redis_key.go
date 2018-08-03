package redis

// Exists key
func (r *Redis) Exists(key string) *Reply {
	reply := r.run("EXISTS", key)
	reply.fixBool()
	return reply
}
