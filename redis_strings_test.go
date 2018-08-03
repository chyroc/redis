package redis_test

import (
	"testing"

	"github.com/Chyroc/redis"
	"time"
)

func TestStringGetSet(t *testing.T) {
	r, as := conn(t)

	{
		// get and set
		as.Nil(r.Set("key", "hello").Err())
		p := r.Get("key")
		as.Nil(p.Err())
		as.False(p.Null())
		as.Equal("hello", p.String())
	}

	{
		// expire
		as.Nil(r.Set("key-with-expire-time", "hello", redis.SetOption{Expire: time.Millisecond * 1000}).Err())
		p := r.Get("key-with-expire-time")
		as.Nil(p.Err())
		as.False(p.Null())
		as.Equal("hello", p.String())

		time.Sleep(time.Millisecond * 1000)

		p = r.Get("key-with-expire-time")
		as.Nil(p.Err())
		as.True(p.Null())
	}

	{
		// nx set id not exist
		p := r.Set("not-exists-key", "hello", redis.SetOption{NX: true})
		as.Nil(p.Err())
		as.False(p.Null())
		as.Equal("OK", p.String())

		p = r.Set("not-exists-key", "new hello", redis.SetOption{NX: true})
		as.Nil(p.Err())
		as.True(p.Null())

		p = r.Get("not-exists-key")
		as.Nil(p.Err())
		as.Equal("hello", p.String())
	}

	{
		// xx only set when already exist
		p := r.Set("exists-key", "hello", redis.SetOption{XX: true})
		as.Nil(p.Err())
		as.Empty(p.String())
		as.True(p.Null())

		p = r.Set("exists-key", "value")
		as.Nil(p.Err())

		p = r.Set("exists-key", "new hello", redis.SetOption{XX: true})
		as.Nil(p.Err())
		as.Equal("OK", p.String())
		as.False(p.Null())

		p = r.Get("exists-key")
		as.Nil(p.Err())
		as.Equal("new hello", p.String())
	}

	// TODO test lock with expire + nx xx
}

func TestStringIncr(t *testing.T) {
	r, as := conn(t)

	p := r.Incr("page_view")
	as.Nil(p.Err())
	as.Equal(1, p.Integer())

	p = r.Incr("page_view")
	as.Nil(p.Err())
	as.Equal(2, p.Integer())
}

func TestStringAppend(t *testing.T) {
	r, as := conn(t)

	p := r.Exists("key")
	as.Nil(p.Err())
	as.False(p.Bool())

	p = r.Append("key", "value1")
	as.Nil(p.Err())
	as.Equal(6, p.Integer())

	p = r.Append("key", " - vl2")
	as.Nil(p.Err())
	as.Equal(12, p.Integer())

	p = r.Get("key")
	as.Nil(p.Err())
	as.Equal("value1 - vl2", p.String())
}

func TestStringBit(t *testing.T) {
	r, as := conn(t)

	// bitcount
	p := r.BitCount("bits")
	as.Nil(p.Err())
	as.Equal(0, p.Integer())

	// setbit
	p = r.SetBit("bits", 0, true)
	as.Nil(p.Err())

	p = r.BitCount("bits")
	as.Nil(p.Err())
	as.Equal(1, p.Integer())

	p = r.SetBit("bits", 3, true)
	as.Nil(p.Err())

	p = r.BitCount("bits")
	as.Nil(p.Err())
	as.Equal(2, p.Integer())

	// getbits
	p = r.GetBit("bits", 0)
	as.Nil(p.Err())
	as.Equal(1, p.Integer())

	p = r.GetBit("bits", 3)
	as.Nil(p.Err())
	as.Equal(1, p.Integer())
}

func TestBiTop(t *testing.T) {
	r, as := conn(t)

	{
		// bits-1 1001
		setbits(t, r, "bits-1", []int{0, 1, 2, 3}, []int{1, 0, 0, 1})

		// bits-1 1011
		setbits(t, r, "bits-2", []int{0, 1, 2, 3}, []int{1, 0, 1, 1})

		// 1001 & 1011 = 1001
		p := r.BitOp(redis.BitOpOption{AND: true}, "and-result", "bits-1", "bits-2")
		as.Nil(p.Err())
		as.Equal(1, p.Integer())
		getbits(t, r, "and-result", []int{0, 1, 2, 3}, []int{1, 0, 0, 1})

		// 1001 | 1011 = 1011
		p = r.BitOp(redis.BitOpOption{OR: true}, "or-result", "bits-1", "bits-2")
		as.Nil(p.Err())
		as.Equal(1, p.Integer())
		getbits(t, r, "or-result", []int{0, 1, 2, 3}, []int{1, 0, 1, 1})

		// 1001 ^ 1011 = 0010
		p = r.BitOp(redis.BitOpOption{XOR: true}, "xor-result", "bits-1", "bits-2")
		as.Nil(p.Err())
		as.Equal(1, p.Integer())
		getbits(t, r, "xor-result", []int{0, 1, 2, 3}, []int{0, 0, 1, 0})

		// ^1001  = 0110
		p = r.BitOp(redis.BitOpOption{NOT: true}, "not-result-1", "bits-1")
		as.Nil(p.Err())
		as.Equal(1, p.Integer())
		getbits(t, r, "not-result-1", []int{0, 1, 2, 3}, []int{0, 1, 1, 0})

		// ^1011  = 0100
		p = r.BitOp(redis.BitOpOption{NOT: true}, "not-result-2", "bits-2")
		as.Nil(p.Err())
		as.Equal(1, p.Integer())
		getbits(t, r, "not-result-2", []int{0, 1, 2, 3}, []int{0, 1, 0, 0})
	}

	{
		as.Nil(r.Set("key1", "foobar").Err())
		as.Nil(r.Set("key2", "abcdef").Err())
		as.Nil(r.BitOp(redis.BitOpOption{AND: true}, "dest", "key1", "key2").Err())
		as.Equal("`bc`ab", r.Get("dest").String())
	}
}
