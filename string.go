package redis

func (r *Redis) Set(key, value string) error {
	if err := r.cmd("SET", key, value); err != nil {
		return err
	}

	_, err := r.read()
	return err
}

func (r *Redis) Get(key string) (string, error) {
	if err := r.cmd("GET", key); err != nil {
		return "", err
	}

	bs, err := r.read()
	return string(bs.([]byte)), err
}
