package redis

func (r *Redis) Set(key, value string) error {
	if err := r.cmd("SET", key, value); err != nil {
		return err
	}

	_, err := r.readUntilCRCL()
	return err
}

func (r *Redis) Get(key string) error {
	if err := r.cmd("GET", key); err != nil {
		return err
	}

	_, err := r.readUntilCRCL()
	return err
}
