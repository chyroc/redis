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

// APPEND key value
func (r *Redis) Append(key, value string, options ...SetOption) *Reply {
	return r.run("APPEND", key, value)
}

// BITCOUNT key [start] [end]
func (r *Redis) BitCount(key string, startEnd ...int) *Reply {
	args := []string{"BITCOUNT", key}
	switch len(startEnd) {
	case 0:
	case 1:
		args = append(args, strconv.Itoa(startEnd[0]))
	case 2:
		args = append(args, strconv.Itoa(startEnd[0]), strconv.Itoa(startEnd[1]))
	default:
		return errToReply(fmt.Errorf("expect get 0, 1 or 2 arguments, but got %v", startEnd))
	}
	return r.run(args...)
}

// GET key
func (r *Redis) Get(key string) *Reply {
	return r.run("GET", key)
}

// GETBIT key offset
func (r *Redis) GetBit(key string, offset int) *Reply {
	return r.run("GETBIT", key, strconv.Itoa(offset))
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

// INCR key
// Available since 1.0.0.
// Time complexity: O(1)
func (r *Redis) Incr(key string) *Reply {
	return r.run("INCR", key)
}

// SETBIT key offset value
func (r *Redis) SetBit(key string, offset int, SetOrRemove bool) *Reply {
	p := r.run("SETBIT", key, strconv.Itoa(offset), boolToString(SetOrRemove))
	p.fixBool()
	return p
}
