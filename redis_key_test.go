package redis_test

import (
	"github.com/Chyroc/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDel(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.Set, "a", "b").Expect(true)
	r.RunTest(e.Del, "a").Expect(1)

	r.RunTest(e.Del, "c").Expect(0)

	r.RunTest(e.Set, "a", "b").Expect(true)
	r.RunTest(e.Set, "b", "b").Expect(true)
	r.RunTest(e.Set, "c", "b").Expect(true)
	r.RunTest(e.Del, "a", "b", "c").Expect(3)
}

func TestDumpRestore(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.Set, "a", "hello, dumping world!").Expect(true)
	r.RunTest(e.Dump, "a").Expect("\x00\x15hello, dumping world!\b\x00j`\u07bd\x84>wu")
	r.RunTest(e.Dump, "b").ExpectNull()

	r.RunTest(e.Restore, "a-2", zeroTimeDuration, "\x00\x15hello, dumping world!\b\x00j`\u07bd\x84>wu", false).ExpectSuccess()
	r.RunTest(e.Get, "a-2").Expect("hello, dumping world!")

	r.RunTest(e.Restore, "a-2", zeroTimeDuration, "\x00\x15hello, dumping world!\b\x00j`\u07bd\x84>wu", false).ExpectError("BUSYKEY Target key name already exists.")
	r.RunTest(e.Restore, "a-2", zeroTimeDuration, "\x00\x15hello, dumping world!\b\x00j`\u07bd\x84>wu", true).ExpectSuccess()

	r.RunTest(e.Restore, "a-3", zeroTimeDuration, "invalid dump data", false).ExpectError("ERR DUMP payload version or checksum are wrong")
}

func TestExist(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.Set, "a", "b").Expect(true)
	r.RunTest(e.Exists, "a").Expect(true)
}

func TestExpireTTL(t *testing.T) {
	//t.Parallel() // TODO

	r := NewTest(t)

	r.RunTest(e.Set, "a", "b").Expect(true)
	r.RunTest(e.TTL, "a").Expect(nil)

	r.RunTest(e.Expire, "a", time.Second*2)
	r.RunTest(e.TTL, "a").ExpectLess(time.Second * 2)
	time.Sleep(time.Second)
	r.RunTest(e.TTL, "a").ExpectLess(time.Second)
	time.Sleep(time.Second)
	r.RunTest(e.TTL, "a").ExpectError(redis.ErrKeyNotExist.Error())

	r.RunTest(e.Set, "c", "b").Expect(true)
	r.RunTest(e.ExpireAt, "c", time.Now().Add(time.Minute))
	r.RunTest(e.TTL, "c").ExpectLess(time.Second * 60)
}

func TestKeys(t *testing.T) {
	r := NewTest(t)

	r.Nil(e.MSet("one", "1", "two", "2", "three", "3", "four", "4"))

	r.RunTest(e.Keys, "*o*").ExpectContains("four", "one", "two")
	r.RunTest(e.Keys, "t??").ExpectContains("two")
	r.RunTest(e.Keys, "t[w]*").ExpectContains("two")
	r.RunTest(e.Keys, "*").ExpectContains("three", "four", "one", "two")
}

func TestMove(t *testing.T) {
	r := NewTest(t)

	// 清空下面会用到的数据库 0 1
	r.Nil(e.Select(0))
	r.Nil(e.FlushDB())
	r.Nil(e.Select(1))
	r.Nil(e.FlushDB())

	// key exist
	r.Nil(e.Select(0))

	r.RunTest(e.Set, "a", "b").Expect(true)
	r.RunTest(e.Exists, "a").Expect(true)

	r.RunTest(e.Move, "a", 1).Expect(true)
	r.RunTest(e.Exists, "a").Expect(false)

	r.Nil(e.Select(1))
	r.RunTest(e.Exists, "a").Expect(true)

	// key not exist
	r.Nil(e.Select(0))

	r.RunTest(e.Exists, "b").Expect(false)
	r.RunTest(e.Move, "b", 1).Expect(false)

	r.Nil(e.Select(1))
	r.RunTest(e.Exists, "b").Expect(false)

	// 当源数据库和目标数据库有相同的 key 时
	r.Nil(e.Select(0))
	r.RunTest(e.Set, "c", "db 0").Expect(true)
	r.Nil(e.Select(1))
	r.RunTest(e.Set, "c", "db 1").Expect(true)
	r.Nil(e.Select(0))
	r.RunTest(e.Move, "c", 1).Expect(false)

	r.RunTest(e.Get, "c").Expect("db 0")
	r.Nil(e.Select(1))
	r.RunTest(e.Get, "c").Expect("db 1")
}

func TestObject(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.Set, "a", "b").Expect(true)
	o := e.Object("a")
	r.RunTest(o.RefCount).Expect(1)
	r.RunTest(o.IdleTime).Expect(0)
	time.Sleep(time.Second * 2)
	r.RunTest(o.IdleTime).ExpectBigger(0)

	r.RunTest(e.Get, "a").Expect("b")
	r.RunTest(o.IdleTime).Expect(0)
	r.RunTest(o.Encoding).Expect("embstr")
	r.RunTest(e.Set, "b", "123456789012345678901234567890123456789012345").Expect(true)
	r.RunTest(e.Object("b").Encoding).Expect("raw")

	r.RunTest(e.Set, "c", "1").Expect(true)
	r.RunTest(e.Object("c").Encoding).Expect("int")
}

func TestPersist(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.Set, "a", "b").Expect(true)
	r.RunTest(e.Expire, "a", time.Second*10).True(true)
	r.RunTest(e.TTL, "a").ExpectLess(time.Second * 10)
	r.RunTest(e.Persist, "a").Expect(true)
	r.RunTest(e.Persist, "a").Expect(false)
}

func TestRandomKey(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.MSet, "fruit", "apple", "drink", "beer", "food", "cookies").ExpectSuccess()
	r.RunTest(e.RandomKey).ExpectBelong("fruit", "drink", "food")
	r.RunTest(e.RandomKey).ExpectBelong("fruit", "drink", "food")
	r.RunTest(e.RandomKey).ExpectBelong("fruit", "drink", "food")
	r.RunTest(e.Keys, "*") // TODO list contain list

	r.Nil(e.FlushDB())
	r.RunTest(e.RandomKey).ExpectNull()
}

func TestRenameRenameNx(t *testing.T) {
	r := NewTest(t)

	// rename
	{
		r.RunTest(e.Set, "a", "b").Expect(true)
		r.RunTest(e.Rename, "a", "new-name").ExpectSuccess()
		r.RunTest(e.Exists, "a").Expect(false)
		r.RunTest(e.Exists, "new-name").Expect(true)

		// not exist key
		r.RunTest(e.Rename, "b", "new-name").ExpectError("ERR no such key")

		// 覆盖
		r.RunTest(e.Set, "c", "1").Expect(true)
		r.RunTest(e.Set, "d", "2").Expect(true)
		r.RunTest(e.Rename, "c", "d").ExpectSuccess()
		r.RunTest(e.Exists, "c").Expect(false)
		r.RunTest(e.Exists, "d").Expect(true)
		r.RunTest(e.Get, "d").Expect("1")
	}

	r.Nil(e.FlushDB())

	// renamenx
	{
		r.RunTest(e.Set, "a", "1").Expect(true)
		r.RunTest(e.RenameNX, "a", "b").Expect(true)
		r.RunTest(e.Set, "c", "1").Expect(true)
		r.RunTest(e.RenameNX, "b", "c").Expect(false)
	}
}

func TestType(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.Set, "a", "1").Expect(true)
	r.RunTest(e.Type, "a").Expect(redis.KeyTypeString)

	// TODO LPush book_list "programming in scala" list

	r.RunTest(e.SAdd, "b", "1").Expect(1)
	r.RunTest(e.Type, "b").Expect(redis.KeyTypeSet)
}

func TestScan(t *testing.T) {
	r := NewTest(t)

	// all
	r.RunTest(e.Set, "a", "1").Expect(true)
	r.RunTest(e.Set, "b", "1").Expect(true)
	r.RunTest(e.Set, "c", "1").Expect(true)
	r.RunTest(e.Set, "d", "1").Expect(true)

	r.RunTest(e.Scan().ALL).ExpectContains("a", "b", "c", "d")

	r.RunTest(e.Set, "a1a", "1").Expect(true)
	r.RunTest(e.Set, "a1b", "1").Expect(true)
	r.RunTest(e.Set, "a1c", "1").Expect(true)

	r.RunTest(e.Scan(redis.ScanOption{Match: "a?*"}).ALL).ExpectContains("a1a", "a1c", "a1b")

	// each
	var vv []string
	var kk []int
	assert.Nil(t, e.Scan().Each(func(k int, v string) error {
		vv = append(vv, v)
		kk = append(kk, k)
		return nil
	}))
	stringContains(r.t, vv, []string{"a", "b", "c", "d"})
	r.Equal([]int{0, 1, 2, 3, 4, 5, 6}, kk)
}
