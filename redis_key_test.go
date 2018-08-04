package redis_test

import (
	"github.com/Chyroc/redis"
	"github.com/stretchr/testify/assert"
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

func TestDumpRestore(t *testing.T) {
	r, _ := conn(t)

	r.RunTest(r.Set, "a", "hello, dumping world!").Expect(true)
	r.RunTest(r.Dump, "a").Expect("\x00\x15hello, dumping world!\b\x00j`\u07bd\x84>wu")
	r.RunTest(r.Dump, "b").ExpectNull()

	r.RunTest(r.Restore, "a-2", 0, "\x00\x15hello, dumping world!\b\x00j`\u07bd\x84>wu", false).ExpectSuccess()
	r.RunTest(r.Get, "a-2").Expect("hello, dumping world!")

	r.RunTest(r.Restore, "a-2", 0, "\x00\x15hello, dumping world!\b\x00j`\u07bd\x84>wu", false).ExpectError("BUSYKEY Target key name already exists.")
	r.RunTest(r.Restore, "a-2", 0, "\x00\x15hello, dumping world!\b\x00j`\u07bd\x84>wu", true).ExpectSuccess()

	r.RunTest(r.Restore, "a-3", 0, "invalid dump data", false).ExpectError("ERR DUMP payload version or checksum are wrong")
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
	r.RunTest(r.TTL, "a").Expect(nil)

	r.RunTest(r.Expire, "a", time.Second*2)
	r.RunTest(r.TTL, "a").ExpectLess(time.Second * 2)
	time.Sleep(time.Second)
	r.RunTest(r.TTL, "a").ExpectLess(time.Second)
	time.Sleep(time.Second)
	r.RunTest(r.TTL, "a").ExpectError(redis.ErrKeyNotExist.Error())

	r.RunTest(r.Set, "c", "b").Expect(true)
	r.RunTest(r.ExpireAt, "c", time.Now().Add(time.Minute))
	r.RunTest(r.TTL, "c").ExpectLess(time.Second * 60)
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
	r, _ := conn(t)

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

func TestPersist(t *testing.T) {
	r, _ := conn(t)

	r.RunTest(r.Set, "a", "b").Expect(true)
	r.RunTest(r.Expire, "a", time.Second*10).True(true)
	r.RunTest(r.TTL, "a").Expect(time.Second * 10)
	r.RunTest(r.Persist, "a").Expect(true)
	r.RunTest(r.Persist, "a").Expect(false)
}

func TestRandomKey(t *testing.T) {
	r, as := conn(t)

	r.RunTest(r.MSet, "fruit", "apple", "drink", "beer", "food", "cookies").ExpectSuccess()
	r.RunTest(r.RandomKey).ExpectBelong("fruit", "drink", "food")
	r.RunTest(r.RandomKey).ExpectBelong("fruit", "drink", "food")
	r.RunTest(r.RandomKey).ExpectBelong("fruit", "drink", "food")
	r.RunTest(r.Keys, "*") // TODO list contain list

	as.Nil(r.FlushDB())
	r.RunTest(r.RandomKey).ExpectNull()
}

func TestRenameRenameNx(t *testing.T) {
	r, as := conn(t)

	// rename
	{
		r.RunTest(r.Set, "a", "b").Expect(true)
		r.RunTest(r.Rename, "a", "new-name").ExpectSuccess()
		r.RunTest(r.Exists, "a").Expect(false)
		r.RunTest(r.Exists, "new-name").Expect(true)

		// not exist key
		r.RunTest(r.Rename, "b", "new-name").ExpectError("ERR no such key")

		// 覆盖
		r.RunTest(r.Set, "c", "1").Expect(true)
		r.RunTest(r.Set, "d", "2").Expect(true)
		r.RunTest(r.Rename, "c", "d").ExpectSuccess()
		r.RunTest(r.Exists, "c").Expect(false)
		r.RunTest(r.Exists, "d").Expect(true)
		r.RunTest(r.Get, "d").Expect("1")
	}

	as.Nil(r.FlushDB())

	// renamenx
	{
		r.RunTest(r.Set, "a", "1").Expect(true)
		r.RunTest(r.RenameNX, "a", "b").Expect(true)
		r.RunTest(r.Set, "c", "1").Expect(true)
		r.RunTest(r.RenameNX, "b", "c").Expect(false)
	}
}

func TestType(t *testing.T) {
	r, _ := conn(t)

	r.RunTest(r.Set, "a", "1").Expect(true)
	r.RunTest(r.Type, "a").Expect(redis.KeyTypeString)

	// TODO LPUSH book_list "programming in scala" list

	r.RunTest(r.SAdd, "b", "1").Expect(1)
	r.RunTest(r.Type, "b").Expect(redis.KeyTypeSet)
}

func TestScan(t *testing.T) {
	r, _ := conn(t)

	// all
	r.RunTest(r.Set, "a", "1").Expect(true)
	r.RunTest(r.Set, "b", "1").Expect(true)
	r.RunTest(r.Set, "c", "1").Expect(true)
	r.RunTest(r.Set, "d", "1").Expect(true)

	r.RunTest(r.Scan().ALL).ExpectContains("a", "b", "c", "d")

	r.RunTest(r.Set, "a1a", "1").Expect(true)
	r.RunTest(r.Set, "a1b", "1").Expect(true)
	r.RunTest(r.Set, "a1c", "1").Expect(true)

	r.RunTest(r.Scan(redis.ScanOption{Match: "a?*"}).ALL).ExpectContains("a1a", "a1c", "a1b")

	// each
	var s []string
	assert.Nil(t, r.Scan().Each(func(k int, v string) error {
		s = append(s, v)
		return nil
	}))
	stringContains(s, []string{"a", "b", "c", "d"})
}
