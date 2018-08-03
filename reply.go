package redis

type Reply struct {
	err     error
	null    bool
	str     string
	integer int64
	boo     bool
}

func (p *Reply) Err() error {
	return p.err
}

func (p *Reply) String() string {
	return p.str
}

func (p *Reply) Integer() int {
	return int(p.integer)
}

func (p *Reply) Integer64() int64 {
	return p.integer
}

func (p *Reply) Null() bool {
	return p.null
}

func (p *Reply) Bool() bool {
	return p.boo
}

func (p *Reply) fixBool() {
	if p.err == nil {
		p.boo = p.integer == 1
	}
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
