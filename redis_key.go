package redis

// Exists key
func (r *Redis) Exists(key string) (bool, error) {
	return r.run("EXISTS", key).fixBool()
}
