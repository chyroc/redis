package redis_test

import (
	"testing"

	"github.com/Chyroc/redis"
	"time"
)

func TestStringGetSetNxXxExpire(t *testing.T) {
	r, _ := conn(t)

	{
		// get and set
		r.RunTest(r.Set, "key", "hello").Expect(true)
		r.RunTest(r.Set, "key", "hello")
		r.RunTest(r.Get, "key").Expect("hello")
	}

	{
		// expire
		r.RunTest(r.Set, "key-with-expire-time", "hello", redis.SetOption{Expire: time.Second * 2}).Expect(true)
		r.RunTest(r.Get, "key-with-expire-time").Expect("hello")

		time.Sleep(time.Second*2 + time.Millisecond*100)

		r.RunTest(r.Get, "key-with-expire-time").Expect(redis.NullString{})
	}

	{
		// nx set id not exist
		r.RunTest(r.Set, "not-exists-key", "hello", redis.SetOption{NX: true}).Expect(true)
		r.RunTest(r.Set, "not-exists-key", "new hello", redis.SetOption{NX: true}).Expect(false)

		r.RunTest(r.Get, "not-exists-key").Expect("hello")
	}

	{
		// xx only set when already exist
		r.RunTest(r.Set, "exists-key", "hello", redis.SetOption{XX: true}).Expect(false)
		r.RunTest(r.Set, "exists-key", "value").Expect(true)
		r.RunTest(r.Set, "exists-key", "new hello", redis.SetOption{XX: true}).Expect(true)

		r.RunTest(r.Get, "exists-key").Expect("new hello")
	}

	// TODO test lock with expire + nx xx
}

func TestStringAppend(t *testing.T) {
	r, _ := conn(t)

	r.RunTest(r.Exists, "key").Expect(false)

	r.RunTest(r.Append, "key", "value1").Expect(6)
	r.RunTest(r.Append, "key", " - vl2").Expect(12)

	r.RunTest(r.Get, "key").Expect("value1 - vl2")
}

func TestStringBit(t *testing.T) {
	r, as := conn(t)
	as.Equal(1, 1)

	// bitcount
	r.RunTest(r.BitCount, "bits", 0)

	// setbit
	r.RunTest(r.SetBit, "bits", 0, true).Expect(0)
	r.RunTest(r.BitCount, "bits", 1)
	r.RunTest(r.SetBit, "bits", 3, true).Expect(0)
	r.RunTest(r.BitCount, "bits", 2)

	// getbits
	r.RunTest(r.GetBit, "bits", 0).Expect(1)
	r.RunTest(r.GetBit, "bits", 3).Expect(1)
}

func TestBiTop(t *testing.T) {
	r, _ := conn(t)

	{
		setbits(t, r, "bits-1", []int{0, 1, 2, 3}, []int{1, 0, 0, 1}) // bits-1 1001
		setbits(t, r, "bits-2", []int{0, 1, 2, 3}, []int{1, 0, 1, 1}) // bits-1 1011

		// 1001 & 1011 = 1001
		r.RunTest(r.BitOp, redis.BitOpAND, "and-result", "bits-1", "bits-2").Expect(1)
		getbits(t, r, "and-result", []int{0, 1, 2, 3}, []int{1, 0, 0, 1})

		// 1001 | 1011 = 1011
		r.RunTest(r.BitOp, redis.BitOpOR, "or-result", "bits-1", "bits-2").Expect(1)
		getbits(t, r, "or-result", []int{0, 1, 2, 3}, []int{1, 0, 1, 1})

		// 1001 ^ 1011 = 0010
		r.RunTest(r.BitOp, redis.BitOpXOR, "xor-result", "bits-1", "bits-2").Expect(1)
		getbits(t, r, "xor-result", []int{0, 1, 2, 3}, []int{0, 0, 1, 0})

		// ^1001  = 0110
		r.RunTest(r.BitOp, redis.BitOpNOT, "not-result-1", "bits-1").Expect(1)
		getbits(t, r, "not-result-1", []int{0, 1, 2, 3}, []int{0, 1, 1, 0})

		// ^1011  = 0100
		r.RunTest(r.BitOp, redis.BitOpNOT, "not-result-2", "bits-2").Expect(1)
		getbits(t, r, "not-result-2", []int{0, 1, 2, 3}, []int{0, 1, 0, 0})
	}

	{
		r.RunTest(r.Set, "key1", "foobar").Expect(true)
		r.RunTest(r.Set, "key2", "abcdef").Expect(true)
		r.RunTest(r.BitOp, redis.BitOpAND, "dest", "key1", "key2").Expect(6)
		r.RunTest(r.Get, "dest").Expect("`bc`ab")
	}
}

func TestStringBitField(t *testing.T) {
	r, _ := conn(t)

	datatype := redis.SignedInt(4) // -8 ~ 7

	// incrby
	r.RunTest(r.BitField("mykey").IncrBy(datatype, 10, 1).Run).Expect(1)

	// incrby -> incrby -> get
	r.RunTest(r.BitField("mykey").
		IncrBy(datatype, 10, 1).
		IncrBy(datatype, 10, 1).
		Get(datatype, 10).Run).Expect(2, 3, 3)

	// todo overflow test
	// todo set test
}

func TestStringDecrIncr(t *testing.T) {
	r, _ := conn(t)

	// exist key
	r.RunTest(r.Set, "k1", "10").Expect(true)

	r.RunTest(r.Decr, "k1").Expect(9)        // decr
	r.RunTest(r.IncrBy, "k1", 10).Expect(19) // incrby
	r.RunTest(r.DecrBy, "k1", 5).Expect(14)  // decrby
	r.RunTest(r.Incr, "k1").Expect(15)       // incr

	// not exist key decr
	r.RunTest(r.Decr, "k2").Expect(-1)
	r.RunTest(r.IncrByFloat, "k2", 23.1).Expect(22.1)
	r.RunTest(r.Get, "k2").Expect("22.1")

	// invalid data type
	r.RunTest(r.Set, "k3", "string").Expect(true)
	r.RunTest(r.Decr, "k3").ExpectError("ERR value is not an integer or out of range")
	r.RunTest(r.DecrBy, "k3", 10).ExpectError("ERR value is not an integer or out of range")
	r.RunTest(r.Incr, "k3").ExpectError("ERR value is not an integer or out of range")
	r.RunTest(r.IncrBy, "k3", 10).ExpectError("ERR value is not an integer or out of range")

	r.RunTest(r.Set, "k3", "1.1").Expect(true)
	r.RunTest(r.IncrByFloat, "k3", 2.3).Expect(3.4)
	r.RunTest(r.Get, "k3").Expect("3.4")

	r.RunTest(r.Set, "k4", "314e-2").Expect(true)
	r.RunTest(r.Get, "k4").Expect("314e-2")
	r.RunTest(r.IncrByFloat, "k4", 0).Expect(3.14)
	r.RunTest(r.Get, "k4").Expect("3.14") // 执行 INCRBYFLOAT 之后格式会被改成非指数符号

	r.RunTest(r.Set, "k5", "2.0000").Expect(true)
	r.RunTest(r.Get, "k5").Expect("2.0000")
	r.RunTest(r.IncrByFloat, "k5", 0).Expect(2)
	r.RunTest(r.Get, "k5").Expect("2") // INCRBYFLOAT 会将无用的 0 忽略掉
}

func TestStringRange(t *testing.T) {
	r, _ := conn(t)

	// set
	r.RunTest(r.Set, "k", "hello, my friend").Expect(true)

	// get-range
	r.RunTest(r.GetRange, "k", 0, 4).Expect("hello")
	r.RunTest(r.GetRange, "k", -1, -5).Expect("") // 不支持回绕操作

	r.RunTest(r.GetRange, "k", -3, -1).Expect("end")
	r.RunTest(r.GetRange, "k", 0, -1).Expect("hello, my friend")
	r.RunTest(r.GetRange, "k", 0, 9090).Expect("hello, my friend")

	// set-range
	r.RunTest(r.SetRange, "k", 1, "----").Expect(16)
	r.RunTest(r.Get, "k").Expect("h----, my friend")
	r.RunTest(r.SetRange, "k", 20, "====").Expect(24)
	r.RunTest(r.Get, "k").Expect("h----, my friend\x00\x00\x00\x00====") // 空白处被"\x00"填充

	// strlen
	r.RunTest(r.StrLen, "k").Expect(24)
}

func TestStringGetSet(t *testing.T) {
	r, _ := conn(t)

	r.RunTest(r.GetSet, "a", "b").ExpectNull()

	r.RunTest(r.Get, "a").Expect("b")
	r.RunTest(r.GetSet, "a", "c").Expect("b")
	r.RunTest(r.Get, "a").Expect("c")

	r.RunTest(r.SAdd, "b", "m").Expect(1)
	r.RunTest(r.SAdd, "b", "k", "l").Expect(2)
	r.RunTest(r.GetSet, "b", "m").ExpectError("WRONGTYPE Operation against a key holding the wrong kind of value")
}

func TestMultiGetSet(t *testing.T) {
	r, as := conn(t)

	// mset
	as.Nil(r.MSet("a", "av"))
	as.Nil(r.MSet("b", "bv", "c", "cv"))
	as.Equal("key value pair, but got 3 arguments", r.MSet("1", "2", "3").Error())

	// mget
	r.RunTest(r.MGet, "a", "b", "c").Expect(redis.NullString{"av", true}, redis.NullString{"bv", true}, redis.NullString{"cv", true})

	r.RunTest(r.MGet, "c", "not-exist").Expect(redis.NullString{"cv", true}, redis.NullString{})

	// msetnx
	r.RunTest(r.MSetNX, "a", "1", "not-exist-1", "1").Expect(false)
	r.RunTest(r.Get, "not-exist-1").Expect(redis.NullString{})
	r.RunTest(r.MSetNX, "not-exist-1", "1", "not-exist-2", "1").Expect(true)
	r.RunTest(r.Get, "not-exist-1").Expect("1")
}
