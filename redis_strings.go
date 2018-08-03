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

type BitOpOption struct {
	AND bool
	OR  bool
	XOR bool
	NOT bool
}

// BITOP operation destkey key [key ...]
//
// Available since 2.6.0.
// Time complexity: O(N)
//
// 对一个或多个保存二进制位的字符串 key 进行位元操作，并将结果保存到 destkey 上。
//
// operation 可以是 AND 、 OR 、 NOT 、 XOR 这四种操作中的任意一种：
//
//   BITOP AND destkey key [key ...] ，对一个或多个 key 求逻辑并，并将结果保存到 destkey 。
//   BITOP OR destkey key [key ...] ，对一个或多个 key 求逻辑或，并将结果保存到 destkey 。
//   BITOP XOR destkey key [key ...] ，对一个或多个 key 求逻辑异或，并将结果保存到 destkey 。
//   BITOP NOT destkey key ，对给定 key 求逻辑非，并将结果保存到 destkey 。
//   除了 NOT 操作之外，其他操作都可以接受一个或多个 key 作为输入。
//
//  处理不同长度的字符串
//
//  当 BITOP 处理不同长度的字符串时，较短的那个字符串所缺少的部分会被看作 0 。
//
//  空的 key 也被看作是包含 0 的字符串序列。
//
// 返回值：
//   保存到 destkey 的字符串的长度，和输入 key 中最长的字符串长度相等。
func (r *Redis) BitOp(option BitOpOption, destkey string, keys ...string) *Reply {
	if len(keys) == 0 {
		return errToReply(fmt.Errorf("need at least 1 key argument, but got 0"))
	}
	if option.AND {
		if option.OR || option.XOR || option.NOT {
			return errToReply(fmt.Errorf("only support 1 operation, already set: and"))
		}
		return r.run(append([]string{"BITOP", "AND", destkey}, keys...)...)
	} else if option.OR {
		if option.AND || option.XOR || option.NOT {
			return errToReply(fmt.Errorf("only support 1 operation, already set: or"))
		}
		return r.run(append([]string{"BITOP", "OR", destkey}, keys...)...)
	} else if option.XOR {
		if option.OR || option.AND || option.NOT {
			return errToReply(fmt.Errorf("only support 1 operation, already set: xor"))
		}
		return r.run(append([]string{"BITOP", "XOR", destkey}, keys...)...)
	} else if option.NOT {
		if option.OR || option.XOR || option.AND {
			return errToReply(fmt.Errorf("only support 1 operation, already set: not"))
		}
		return r.run("BITOP", "NOT", destkey, keys[0])
	}

	return errToReply(fmt.Errorf("invalid operation, should be one of: or, and, xor and not"))
}

type BitField struct {
	r       *Redis
	key     string
	err     error
	actions []string
}

type BitFieldOverflow string

const (
	BitFieldOverflowWrap BitFieldOverflow = "WRAP"
	BitFieldOverflowSat  BitFieldOverflow = "SAT"
	BitFieldOverflowFail BitFieldOverflow = "FAIL"
)

func (b *BitField) Get(typ DataType, offset int) *BitField {
	if b.err != nil {
		return b
	}
	if err := typ.Err(); err != nil {
		b.err = err
		return b
	}
	b.actions = append(b.actions, "GET", typ.String(), strconv.Itoa(offset))
	return b
}

func (b *BitField) Set(typ DataType, offset, value int) *BitField {
	if b.err != nil {
		return b
	}
	if err := typ.Err(); err != nil {
		b.err = err
		return b
	}
	b.actions = append(b.actions, "SET", typ.String(), strconv.Itoa(offset), strconv.Itoa(value))
	return b
}

func (b *BitField) Incrby(typ DataType, offset, increment int) *BitField {
	if b.err != nil {
		return b
	}
	if err := typ.Err(); err != nil {
		b.err = err
		return b
	}
	b.actions = append(b.actions, "INCRBY", typ.String(), strconv.Itoa(offset), strconv.Itoa(increment))
	return b
}

func (b *BitField) Overflow(f BitFieldOverflow) *BitField {
	if b.err != nil {
		return b
	}
	b.actions = append(b.actions, "OVERFLOW", string(f))
	return b
}

func (b *BitField) Run() *Reply {
	if b.err != nil {
		return errToReply(b.err)
	}
	return b.r.run(append([]string{"BITFIELD", b.key}, b.actions...)...)
}

// BITFIELD key [GET type offset] [SET type offset value] [INCRBY type offset increment] [OVERFLOW WRAP|SAT|FAIL]
//
//   Available since 3.2.0.
//   Time complexity: O(1) for each subcommand specified
func (r *Redis) BitField(key string) *BitField {
	return &BitField{r: r, key: key}
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
