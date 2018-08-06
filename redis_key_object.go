package redis

// Object ...
type Object struct {
	redis *Redis
	key   string
}

// RefCount <key> 返回给定 key 引用所储存的值的次数。此命令主要用于除错。
func (r *Object) RefCount() (int, error) {
	return r.redis.run("OBJECT", "REFCOUNT", r.key).int()
}

// Encoding <key> 返回给定 key 锁储存的值所使用的内部表示(representation)。
func (r *Object) Encoding() (string, error) {
	p := r.redis.run("OBJECT", "ENCODING", r.key)
	if p.errNotFromReply != nil {
		return "", p.errNotFromReply
	}
	return p.str, nil
}

// IdleTime <key> 返回给定 key 自储存以来的空闲时间(idle， 没有被读取也没有被写入)，以秒为单位。
func (r *Object) IdleTime() (int, error) {
	return r.redis.run("OBJECT", "IDLETIME", r.key).int()
}
