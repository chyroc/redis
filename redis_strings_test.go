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
		as.Nil(r.Set("key-with-expire-time", "hello", redis.SetOption{Expire: time.Second * 2}).Err())
		p := r.Get("key-with-expire-time")
		as.Nil(p.Err())
		as.False(p.Null())
		as.Equal("hello", p.String())

		time.Sleep(time.Second*2 + time.Millisecond*100)

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

func TestStringBitField(t *testing.T) {
	r, as := conn(t)

	datatype := redis.SignedInt(4) // -8 ~ 7
	// incrby
	p := r.BitField("mykey").Incrby(datatype, 10, 1).Run()
	as.Nil(p.Err())
	as.Len(p.Replys(), 1)
	as.Equal(1, p.Replys()[0].Integer())

	// incrby -> incrby -> get
	p = r.BitField("mykey").Incrby(datatype, 10, 1).
		Incrby(datatype, 10, 1).
		Get(datatype, 10).Run()
	as.Nil(p.Err())
	as.Len(p.Replys(), 3)
	as.Equal(2, p.Replys()[0].Integer())
	as.Equal(3, p.Replys()[1].Integer())
	as.Equal(3, p.Replys()[2].Integer())

	// todo overflow test
	// todo set test
}

func TestStringDecrIncr(t *testing.T) {
	r, as := conn(t)

	// exist key
	as.Nil(r.Set("k1", "10").Err())

	as.Equal(9, r.Decr("k1").Integer())        // decr
	as.Equal(19, r.IncrBy("k1", 10).Integer()) // incrby
	as.Equal(14, r.DecrBy("k1", 5).Integer())  // decrby
	as.Equal(15, r.Incr("k1").Integer())       // incr

	// not exist key decr
	as.Equal(-1, r.Decr("k2").Integer())

	// invalid data type
	as.Nil(r.Set("k3", "string").Err())
	as.Equal("ERR value is not an integer or out of range", r.Decr("k3").Err().Error())
	as.Equal("ERR value is not an integer or out of range", r.DecrBy("k3", 10).Err().Error())
	as.Equal("ERR value is not an integer or out of range", r.Incr("k3").Err().Error())
	as.Equal("ERR value is not an integer or out of range", r.IncrBy("k3", 10).Err().Error())

}
