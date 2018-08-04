package redis

import (
	"fmt"
	"strconv"
)

// HScan ...
type HScan struct {
	redis *Redis

	err  error
	key  string
	args []string

	cursor int
	result map[string]string

	eachCursor int
}

// Each ...
func (r *HScan) Each(f func(k int, field, value string) error) error {
	r.eachCursor = -1
	for {
		if r.cursor == 0 {
			return nil // end
		}
		m, err := r.Next()
		if err != nil {
			return err
		}
		for k, v := range m {
			r.eachCursor++
			if err := f(r.eachCursor, k, v); err != nil {
				return err
			}
		}
	}
}

// ALL ...
func (r *HScan) ALL() (map[string]string, error) {
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
func (r *HScan) Next() (map[string]string, error) {
	if r.cursor == 0 {
		return nil, ErrIteratorEnd
	} else if r.cursor == -1 {
		r.cursor = 0
	}

	p := r.redis.run(append([]string{"HSCAN", r.key, strconv.Itoa(r.cursor)}, r.args...)...)
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
	m, err := p.replys[1].fixMap()
	if err != nil {
		return nil, err
	}

	r.cursor = next
	for k, v := range m {
		r.result[k] = v
	}

	return m, nil
}
