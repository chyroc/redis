package redis

import (
	"fmt"
	"strconv"
)

// SScan ...
type SScan struct {
	redis *Redis

	err  error
	key  string
	args []string

	cursor int
	result []string

	eachCursor int
}

// Each ...
func (r *SScan) Each(f func(k int, v string) error) error {
	r.eachCursor = -1
	for {
		if r.cursor == 0 {
			return nil // end
		}

		result, err := r.Next()
		if err != nil {
			return err
		}
		for _, v := range result {
			r.eachCursor++
			if err := f(r.eachCursor, v); err != nil {
				return err
			}
		}
	}
}

// ALL ...
func (r *SScan) ALL() ([]string, error) {
	for {
		_, err := r.Next()
		if err == ErrIteratorEnd {
			return r.result, nil
		} else if err != nil {
			return nil, err
		}
	}
}

// Next ...
func (r *SScan) Next() ([]string, error) {
	if r.cursor == 0 {
		return nil, ErrIteratorEnd
	} else if r.cursor == -1 {
		r.cursor = 0
	}

	p := r.redis.run(append([]string{"SSCAN", r.key, strconv.Itoa(r.cursor)}, r.args...)...)
	if p.errNotFromReply != nil {
		return nil, p.errNotFromReply
	}

	if len(p.replys) != 2 {
		return nil, fmt.Errorf("expect 2 return, bu got %d", len(p.replys))
	}

	// cursor
	if p.replys[0].errNotFromReply != nil {
		return nil, p.replys[0].errNotFromReply
	}
	next, err := strconv.Atoi(p.replys[0].str)
	if err != nil {
		return nil, err
	}

	// item
	var s []string
	for _, v := range p.replys[1].replys {
		s = append(s, v.str)
	}

	r.cursor = next
	r.result = append(r.result, s...)

	return s, nil
}
