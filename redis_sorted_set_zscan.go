package redis

import (
	"fmt"
	"strconv"
)

// ZScan ...
type ZScan struct {
	redis *Redis

	err  error
	key  string
	args []string

	cursor int
	result []*SortedSet

	eachCursor int
}

// Each ...
func (r *ZScan) Each(f func(k int, v *SortedSet) error) error {
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
func (r *ZScan) ALL() ([]*SortedSet, error) {
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
func (r *ZScan) Next() ([]*SortedSet, error) {
	if r.cursor == 0 {
		return nil, ErrIteratorEnd
	} else if r.cursor == -1 {
		r.cursor = 0
	}

	p := r.redis.run(append([]string{"ZSCAN", r.key, strconv.Itoa(r.cursor)}, r.args...)...)
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
	s, err := p.replys[1].fixSortedSetSliceWithScores()
	if err != nil {
		return nil, err
	}

	r.cursor = next
	r.result = append(r.result, s...)

	return s, nil
}
