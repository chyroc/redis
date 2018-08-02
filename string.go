package redis

func (r *Redis) Set(key, value string) *Reply {
	if err := r.cmd("SET", key, value); err != nil {
		return &Reply{err: err}
	}

	return r.readToReply()
}

func (r *Redis) Get(key string) *Reply {
	if err := r.cmd("GET", key); err != nil {
		return &Reply{err: err}
	}

	return r.readToReply()
}
