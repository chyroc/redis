package redis

// FlushDB ...
func (r *Redis) FlushDB() error {
	return r.run("FLUSHDB").err
}
