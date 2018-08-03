package redis_test

import (
	"github.com/Chyroc/redis"
	"testing"
	"time"
)

func TestDel(t *testing.T) {
	r, _ := conn(t)

	r.RunTest(r.Set, "a", "b").Expect(true)
	r.RunTest(r.Del, "a").Expect(1)

	r.RunTest(r.Del, "c").Expect(0)

	r.RunTest(r.Set, "a", "b").Expect(true)
	r.RunTest(r.Set, "b", "b").Expect(true)
	r.RunTest(r.Set, "c", "b").Expect(true)
	r.RunTest(r.Del, "a", "b", "c").Expect(3)
}

func TestDump(t *testing.T) {
	r, _ := conn(t)

	r.RunTest(r.Set, "a", "hello, dumping world!").Expect(true)
	r.RunTest(r.Dump, "a").Expect("\x00\x15hello, dumping world!\b\x00j`\u07bd\x84>wu")

	r.RunTest(r.Dump, "b").ExpectNull()
}

func TestExist(t *testing.T) {
	r, _ := conn(t)

	r.RunTest(r.Set, "a", "b").Expect(true)
	r.RunTest(r.Exists, "a").Expect(true)
}

func TestExpireTTL(t *testing.T) {
	//t.Parallel() // TODO

	r, _ := conn(t)

	r.RunTest(r.Set, "a", "b").Expect(true)
	r.RunTest(r.TTL, "a").Expect(-1)

	r.RunTest(r.Expire, "a", 2)
	r.RunTest(r.TTL, "a").Expect(2)
	time.Sleep(time.Second)
	r.RunTest(r.TTL, "a").Expect(1)
	time.Sleep(time.Second)
	r.RunTest(r.TTL, "a").ExpectError(redis.ErrKeyNotExist.Error())

	r.RunTest(r.Set, "c", "b").Expect(true)
	r.RunTest(r.ExpireAt, "c", time.Now().Add(time.Minute))
	r.RunTest(r.TTL, "c").Expect(60)
}

func TestKeys(t *testing.T) {
	r, as := conn(t)

	as.Nil(r.MSet("one", "1", "two", "2", "three", "3", "four", "4"))

	r.RunTest(r.Keys, "*o*").Expect("four", "one", "two")
	r.RunTest(r.Keys, "t??").Expect("two")
	r.RunTest(r.Keys, "t[w]*").Expect("two")
	r.RunTest(r.Keys, "*").Expect("three", "four", "one", "two")
}

func TestMove(t *testing.T) {
	r, as := conn(t)

	// 清空下面会用到的数据库 0 1
	as.Nil(r.Select(0))
	as.Nil(r.FlushDB())
	as.Nil(r.Select(1))
	as.Nil(r.FlushDB())

	// key exist
	as.Nil(r.Select(0))

	r.RunTest(r.Set, "a", "b").Expect(true)
	r.RunTest(r.Exists, "a").Expect(true)

	r.RunTest(r.Move, "a", 1).Expect(true)
	r.RunTest(r.Exists, "a").Expect(false)

	as.Nil(r.Select(1))
	r.RunTest(r.Exists, "a").Expect(true)

	// key not exist
	as.Nil(r.Select(0))

	r.RunTest(r.Exists, "b").Expect(false)
	r.RunTest(r.Move, "b", 1).Expect(false)

	as.Nil(r.Select(1))
	r.RunTest(r.Exists, "b").Expect(false)

	// 当源数据库和目标数据库有相同的 key 时
	as.Nil(r.Select(0))
	r.RunTest(r.Set, "c", "db 0").Expect(true)
	as.Nil(r.Select(1))
	r.RunTest(r.Set, "c", "db 1").Expect(true)
	as.Nil(r.Select(0))
	r.RunTest(r.Move, "c", 1).Expect(false)

	r.RunTest(r.Get, "c").Expect("db 0")
	as.Nil(r.Select(1))
	r.RunTest(r.Get, "c").Expect("db 1")
}

func TestObject(t *testing.T) {
	r, as := conn(t)
	as.Equal(1, 1)

	r.RunTest(r.Set, "a", "b").Expect(true)
	o := r.Object("a")
	r.RunTest(o.RefCount).Expect(1)
	r.RunTest(o.IdleTime).Expect(0)
	time.Sleep(time.Second * 2)
	r.RunTest(o.IdleTime).ExpectBigger(0)

	r.RunTest(r.Get, "a").Expect("b")
	r.RunTest(o.IdleTime).Expect(0)
	r.RunTest(o.Encoding).Expect("embstr")
	r.RunTest(r.Set, "b", "123456789012345678901234567890123456789012345").Expect(true)
	r.RunTest(r.Object("b").Encoding).Expect("raw")

	r.RunTest(r.Set, "c", "1").Expect(true)
	r.RunTest(r.Object("c").Encoding).Expect("int")
}
