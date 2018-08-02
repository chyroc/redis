package redis

// FLUSHDB
func (r *Redis) FlushDB() *Reply {
	if err := r.cmd("FLUSHDB"); err != nil {
		return &Reply{err: err}
	}

	return r.readToReply()
}
