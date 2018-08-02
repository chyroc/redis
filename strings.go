package redis

import (
	"fmt"
	"strconv"
	"time"
)

// https://redis.io/commands#string
// http://redisdoc.com/string/index.html

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

	return r.run(args...)
}

// GET key
func (r *Redis) Get(key string) *Reply {
	return r.run("GET", key)
}

// INCR key
// Available since 1.0.0.
// Time complexity: O(1)
func (r *Redis) Incr(key string) *Reply {
	return r.run("INCR", key)
}
