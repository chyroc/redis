package redis

// EXISTS key
func (r *Redis) Exists(key string) *Reply {
	reply := r.run("EXISTS", key)
	if reply.err != nil {
		return reply
	}
	reply.boo = reply.integer == 1
	return reply
}
