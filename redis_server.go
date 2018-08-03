package redis

// FlushDB ...
func (r *Redis) FlushDB() *Reply {
	return r.run("FLUSHDB")
}
