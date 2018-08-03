package redis

// Reply ...
type Reply struct {
	err     error
	null    bool
	str     string
	integer int64
	boo     bool

	replys []*Reply
}

// Err ...
func (p *Reply) Err() error {
	return p.err
}

// String ...
func (p *Reply) String() string {
	return p.str
}

// Integer ...
func (p *Reply) Integer() int {
	return int(p.integer)
}

// Integer64 ...
func (p *Reply) Integer64() int64 {
	return p.integer
}

// Null ...
func (p *Reply) Null() bool {
	return p.null
}

// Bool ...
func (p *Reply) Bool() bool {
	return p.boo
}

// Replys ...
func (p *Reply) Replys() []*Reply {
	return p.replys
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

func nullReply() *Reply {
	return &Reply{null: true}
}
