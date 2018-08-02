package redis

type Reply struct {
	err error
	str nullString
}

func (r *Reply) Err() error {
	return r.err
}

func (r *Reply) String() string {
	return r.str.String
}

func (r *Reply) StringNull() bool {
	return !r.str.Valid
}
