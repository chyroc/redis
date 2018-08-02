package redis

import (
	"fmt"
	"strconv"
	"time"
)

// https://redis.io/commands#string
// http://redisdoc.com/string/set.html

type SetOption struct {
	Expire time.Duration
	NX     bool // Only set the key if it does not already exist.
	XX     bool // Only set the key if it already exist.
}

// SET key value [expiration EX seconds|PX milliseconds] [NX|XX]
func (r *Redis) Set(key, value string, options ...SetOption) *Reply {
	if len(options) > 1 {
		return &Reply{err: fmt.Errorf("must have 1 option")}
	}

	args := []string{"SET", key, value}

	if len(options) > 0 {
		option := options[0]
		if option.Expire >= time.Millisecond {
			args = append(args, "PX", strconv.Itoa(int(option.Expire/time.Millisecond)))
		}
		if option.NX && option.XX {
			return &Reply{err: fmt.Errorf("cannot set NX and XX option at the same time")}
		} else if option.NX {
			args = append(args, "NX")
		} else if option.XX {
			args = append(args, "XX")
		}

	}

	if err := r.cmd(args...); err != nil {
		return &Reply{err: err}
	}

	return r.readToReply()
}

// GET key
func (r *Redis) Get(key string) *Reply {
	if err := r.cmd("GET", key); err != nil {
		return &Reply{err: err}
	}

	return r.readToReply()
}
