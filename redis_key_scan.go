package redis

import (
	"fmt"
	"strconv"
)

// Scan ...
type Scan struct {
	redis *Redis

	err  error
	args []string

	cursor int
	result []string

	eachCursor int
}

// Each ...
func (r *Scan) Each(f func(k int, v string) error) error {
	r.eachCursor = 0
	for {
		if r.cursor == 0 {
			return nil // end
		}
		if len(r.result) <= r.eachCursor {
			if _, err := r.Next(); err != nil {
				return err
			}
		}
		if err := f(r.eachCursor, r.result[r.eachCursor]); err != nil {
			return err
		}
	}
}

// ALL ...
func (r *Scan) ALL() ([]string, error) {
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
func (r *Scan) Next() ([]string, error) {
	if r.cursor == 0 {
		return nil, ErrIteratorEnd
	} else if r.cursor == -1 {
		r.cursor = 0
	}

	p := r.redis.run(append([]string{"SCAN", strconv.Itoa(r.cursor)}, r.args...)...)
	if p.err != nil {
		return nil, p.err
	}

	if len(p.replys) != 2 {
		return nil, fmt.Errorf("expect 2 return, bu got %d", len(p.replys))
	}

	// cursor
	if p.replys[0].err != nil {
		return nil, p.replys[0].err
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
