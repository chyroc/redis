package redis

import (
	"bytes"
	"fmt"
)

// Reply ...
type Reply struct {
	err     error
	null    bool
	str     string
	integer int64

	replys []*Reply
}

// String ...
func (p *Reply) String() string {
	if p.err != nil {
		return fmt.Sprintf("Err: %v", p.err)
	}
	if p.null {
		return "NULL"
	}
	if p.str != "" {
		return fmt.Sprintf("String: %v", p.str)
	}
	if p.integer != 0 {
		return fmt.Sprintf("Integet: %v", p.integer)
	}
	if len(p.replys) > 0 {
		buf := new(bytes.Buffer)
		buf.WriteString("List:")
		for _, v := range p.replys {
			buf.WriteString("  ")
			buf.WriteString(v.String())
		}
		return buf.String()
	}
	return ""
}

// Integer ...
func (p *Reply) int() (int, error) {
	if p.err != nil {
		return 0, p.err
	}
	return int(p.integer), nil // TODO int64?
}

func (p *Reply) string() (NullString, error) {
	if p.err != nil {
		return NullString{}, p.err
	}
	if p.null {
		return NullString{}, nil
	}

	return NullString{String: p.str, Valid: true}, nil
}

func (p *Reply) fixBool() (bool, error) {
	if p.err == nil {
		return p.integer == 1, nil
	}
	return false, p.err
}

func errToReply(err error) *Reply {
	if err != nil {
		return &Reply{err: err}
	}
	return nil
}

func bytesToReply(bs []byte) *Reply {
	return &Reply{str: string(bs)}
}
