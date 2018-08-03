package redis

type Reply struct {
	err     error
	null    bool
	str     string
	integer int64
	boo     bool
}

func (r *Reply) Err() error {
	return r.err
}

func (r *Reply) String() string {
	return r.str
}

func (r *Reply) Integer() int {
	return int(r.integer)
}

func (r *Reply) Integer64() int64 {
	return r.integer
}

func (r *Reply) Null() bool {
	return r.null
}

func (r *Reply) Bool() bool {
	return r.boo
}

func errToReply(err error) *Reply {
	if err != nil {
		return &Reply{err: err}
	}
	return nil
}

func intToReply(i int64) *Reply {
	return &Reply{integer: i}
}

func bytesToReply(bs []byte) *Reply {
	return &Reply{str: string(bs)}
}

func strToReply(s string) *Reply {
	return &Reply{str: s}
}

func nullReply() *Reply {
	return &Reply{null: true}
}
