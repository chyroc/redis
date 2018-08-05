package redis_test

import (
	"testing"

	"github.com/Chyroc/redis"
	"time"
)

func TestStringGetSetNxXxExpire(t *testing.T) {
	r := NewTest(t)

	{
		// get and set
		r.RunTest(e.Set, "key", "hello").Expect(true)
		r.RunTest(e.Set, "key", "hello")
		r.RunTest(e.Get, "key").Expect("hello")
	}

	{
		// expire
		r.RunTest(e.Set, "key-with-expire-time", "hello", redis.SetOption{Expire: time.Second * 2}).Expect(true)
		r.RunTest(e.Get, "key-with-expire-time").Expect("hello")

		time.Sleep(time.Second*2 + time.Millisecond*100)

		r.RunTest(e.Get, "key-with-expire-time").Expect(redis.NullString{})
	}

	{
		// nx set id not exist
		r.RunTest(e.Set, "not-exists-key", "hello", redis.SetOption{NX: true}).Expect(true)
		r.RunTest(e.Set, "not-exists-key", "new hello", redis.SetOption{NX: true}).Expect(false)

		r.RunTest(e.Get, "not-exists-key").Expect("hello")
	}

	{
		// xx only set when already exist
		r.RunTest(e.Set, "exists-key", "hello", redis.SetOption{XX: true}).Expect(false)
		r.RunTest(e.Set, "exists-key", "value").Expect(true)
		r.RunTest(e.Set, "exists-key", "new hello", redis.SetOption{XX: true}).Expect(true)

		r.RunTest(e.Get, "exists-key").Expect("new hello")
	}

	// TODO test lock with expire + nx xx
}

func TestStringAppend(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.Exists, "key").Expect(false)

	r.RunTest(e.Append, "key", "value1").Expect(6)
	r.RunTest(e.Append, "key", " - vl2").Expect(12)

	r.RunTest(e.Get, "key").Expect("value1 - vl2")
}

func TestStringBit(t *testing.T) {
	r := NewTest(t)
	r.as.Equal(1, 1)

	// bitcount
	r.RunTest(e.BitCount, "bits", 0)

	// setbit
	r.RunTest(e.SetBit, "bits", 0, true).Expect(0)
	r.RunTest(e.BitCount, "bits", 1)
	r.RunTest(e.SetBit, "bits", 3, true).Expect(0)
	r.RunTest(e.BitCount, "bits", 2)

	// getbits
	r.RunTest(e.GetBit, "bits", 0).Expect(1)
	r.RunTest(e.GetBit, "bits", 3).Expect(1)
}

func TestBiTop(t *testing.T) {
	r := NewTest(t)

	{
		r.SetBits("bits-1", []int{0, 1, 2, 3}, []int{1, 0, 0, 1}) // bits-1 1001
		r.SetBits("bits-2", []int{0, 1, 2, 3}, []int{1, 0, 1, 1}) // bits-1 1011

		// 1001 & 1011 = 1001
		r.RunTest(e.BitOp, redis.BitOpAND, "and-result", "bits-1", "bits-2").Expect(1)
		r.GetBits("and-result", []int{0, 1, 2, 3}, []int{1, 0, 0, 1})

		// 1001 | 1011 = 1011
		r.RunTest(e.BitOp, redis.BitOpOR, "or-result", "bits-1", "bits-2").Expect(1)
		r.GetBits("or-result", []int{0, 1, 2, 3}, []int{1, 0, 1, 1})

		// 1001 ^ 1011 = 0010
		r.RunTest(e.BitOp, redis.BitOpXOR, "xor-result", "bits-1", "bits-2").Expect(1)
		r.GetBits("xor-result", []int{0, 1, 2, 3}, []int{0, 0, 1, 0})

		// ^1001  = 0110
		r.RunTest(e.BitOp, redis.BitOpNOT, "not-result-1", "bits-1").Expect(1)
		r.GetBits("not-result-1", []int{0, 1, 2, 3}, []int{0, 1, 1, 0})

		// ^1011  = 0100
		r.RunTest(e.BitOp, redis.BitOpNOT, "not-result-2", "bits-2").Expect(1)
		r.GetBits("not-result-2", []int{0, 1, 2, 3}, []int{0, 1, 0, 0})
	}

	{
		r.RunTest(e.Set, "key1", "foobar").Expect(true)
		r.RunTest(e.Set, "key2", "abcdef").Expect(true)
		r.RunTest(e.BitOp, redis.BitOpAND, "dest", "key1", "key2").Expect(6)
		r.RunTest(e.Get, "dest").Expect("`bc`ab")
	}
}

func TestStringBitField(t *testing.T) {
	r := NewTest(t)

	datatype := redis.SignedInt(4) // -8 ~ 7

	// incrby
	r.RunTest(e.BitField("mykey").IncrBy(datatype, 10, 1).Run).Expect(1)

	// incrby -> incrby -> get
	r.RunTest(e.BitField("mykey").
		IncrBy(datatype, 10, 1).
		IncrBy(datatype, 10, 1).
		Get(datatype, 10).Run).Expect(2, 3, 3)

	// todo overflow test
	// todo set test
}

func TestStringDecrIncr(t *testing.T) {
	r := NewTest(t)

	// exist key
	r.RunTest(e.Set, "k1", "10").Expect(true)

	r.RunTest(e.Decr, "k1").Expect(9)        // decr
	r.RunTest(e.IncrBy, "k1", 10).Expect(19) // incrby
	r.RunTest(e.DecrBy, "k1", 5).Expect(14)  // decrby
	r.RunTest(e.Incr, "k1").Expect(15)       // incr

	// not exist key decr
	r.RunTest(e.Decr, "k2").Expect(-1)
	r.RunTest(e.IncrByFloat, "k2", 23.1).Expect(22.1)
	r.RunTest(e.Get, "k2").Expect("22.1")

	// invalid data type
	r.RunTest(e.Set, "k3", "string").Expect(true)
	r.RunTest(e.Decr, "k3").ExpectError("ERR value is not an integer or out of range")
	r.RunTest(e.DecrBy, "k3", 10).ExpectError("ERR value is not an integer or out of range")
	r.RunTest(e.Incr, "k3").ExpectError("ERR value is not an integer or out of range")
	r.RunTest(e.IncrBy, "k3", 10).ExpectError("ERR value is not an integer or out of range")

	r.RunTest(e.Set, "k3", "1.1").Expect(true)
	r.RunTest(e.IncrByFloat, "k3", 2.3).Expect(3.4)
	r.RunTest(e.Get, "k3").Expect("3.4")

	r.RunTest(e.Set, "k4", "314e-2").Expect(true)
	r.RunTest(e.Get, "k4").Expect("314e-2")
	r.RunTest(e.IncrByFloat, "k4", 0).Expect(3.14)
	r.RunTest(e.Get, "k4").Expect("3.14") // 执行 INCRBYFLOAT 之后格式会被改成非指数符号

	r.RunTest(e.Set, "k5", "2.0000").Expect(true)
	r.RunTest(e.Get, "k5").Expect("2.0000")
	r.RunTest(e.IncrByFloat, "k5", 0).Expect(2)
	r.RunTest(e.Get, "k5").Expect("2") // INCRBYFLOAT 会将无用的 0 忽略掉
}

func TestStringRange(t *testing.T) {
	r := NewTest(t)

	// set
	r.RunTest(e.Set, "k", "hello, my friend").Expect(true)

	// get-range
	r.RunTest(e.GetRange, "k", 0, 4).Expect("hello")
	r.RunTest(e.GetRange, "k", -1, -5).Expect("") // 不支持回绕操作

	r.RunTest(e.GetRange, "k", -3, -1).Expect("end")
	r.RunTest(e.GetRange, "k", 0, -1).Expect("hello, my friend")
	r.RunTest(e.GetRange, "k", 0, 9090).Expect("hello, my friend")

	// set-range
	r.RunTest(e.SetRange, "k", 1, "----").Expect(16)
	r.RunTest(e.Get, "k").Expect("h----, my friend")
	r.RunTest(e.SetRange, "k", 20, "====").Expect(24)
	r.RunTest(e.Get, "k").Expect("h----, my friend\x00\x00\x00\x00====") // 空白处被"\x00"填充

	// strlen
	r.RunTest(e.StrLen, "k").Expect(24)
}

func TestStringGetSet(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.GetSet, "a", "b").ExpectNull()

	r.RunTest(e.Get, "a").Expect("b")
	r.RunTest(e.GetSet, "a", "c").Expect("b")
	r.RunTest(e.Get, "a").Expect("c")

	r.RunTest(e.SAdd, "b", "m").Expect(1)
	r.RunTest(e.SAdd, "b", "k", "l").Expect(2)
	r.RunTest(e.GetSet, "b", "m").ExpectError("WRONGTYPE Operation against a key holding the wrong kind of value")
}

func TestMultiGetSet(t *testing.T) {
	r := NewTest(t)

	// mset
	r.as.Nil(e.MSet("a", "av"))
	r.as.Nil(e.MSet("b", "bv", "c", "cv"))
	r.RunTest(e.MSet, "1", "2", "3").ExpectError("key value pair, but got 3 arguments")

	// mget
	r.RunTest(e.MGet, "a", "b", "c").Expect(redis.NullString{"av", true}, redis.NullString{"bv", true}, redis.NullString{"cv", true})
	r.RunTest(e.MGet, "c", "not-exist").Expect(redis.NullString{"cv", true}, redis.NullString{})

	// msetnx
	r.RunTest(e.MSetNX, "a", "1", "not-exist-1", "1").Expect(false)
	r.RunTest(e.Get, "not-exist-1").Expect(redis.NullString{})
	r.RunTest(e.MSetNX, "not-exist-1", "1", "not-exist-2", "1").Expect(true)
	r.RunTest(e.Get, "not-exist-1").Expect("1")
}
