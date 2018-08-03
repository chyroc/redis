package redis

// FLUSHDB
func (r *Redis) FlushDB() *Reply {
	return r.run("FLUSHDB")
}
