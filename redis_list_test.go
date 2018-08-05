package redis_test

import (
	"github.com/Chyroc/redis"
	"testing"
	"time"
)

func TestListPushPop(t *testing.T) {
	r := NewTest(t)

	// len
	r.RunTest(e.LLen, "a").Expect(0)

	// left push
	r.RunTest(e.LPush, "a", "1").Expect(1)
	r.RunTest(e.LPush, "a", "1").Expect(2) // 重复
	r.RunTest(e.LPush, "a", "string").Expect(3)

	// len
	r.RunTest(e.LLen, "a").Expect(3)

	// left pop
	r.RunTest(e.LPop, "a").Expect("string")

	// right pop
	r.RunTest(e.RPop, "a").Expect("1")
	r.RunTest(e.RPop, "a").Expect("1")

	// pop nil
	r.RunTest(e.RPop, "a").Expect(redis.NullString{})
	r.RunTest(e.LPop, "a").Expect(redis.NullString{})
}

func TestListRange(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.RPush, "a", "1").Expect(1)
	r.RunTest(e.LRange, "a", 0, 0).Expect("1")

	r.RunTest(e.RPush, "a", "2").Expect(2)
	r.RunTest(e.LRange, "a", 0, 1).Expect("1", "2")
}

func TestListBlock(t *testing.T) {
	r := NewTest(t)

	// block pop

	// left
	r.RunTest(e.LPush, "a", "1").Expect(1)
	r.RunTest(e.LPush, "b", "2").Expect(1)

	r.RunTest(e.BLPop, zeroTimeDuration, "c", "a", "b").Expect(map[string]string{"a": "1"})
	r.RunTest(e.BLPop, zeroTimeDuration, "c", "a", "b").Expect(map[string]string{"b": "2"})

	r.RunTest(e.BLPop, time.Second, "c", "a", "b").Expect(zeroMap)

	// right
	r.RunTest(e.RPush, "a", "1").Expect(1)
	r.RunTest(e.RPush, "b", "2").Expect(1)

	r.RunTest(e.BRPop, zeroTimeDuration, "c", "a", "b").Expect(map[string]string{"a": "1"})
	r.RunTest(e.BRPop, zeroTimeDuration, "c", "a", "b").Expect(map[string]string{"b": "2"})

	r.RunTest(e.BRPop, time.Second, "c", "a", "b").Expect(zeroMap)
}

func TestListRPopLPush(t *testing.T) {
	r := NewTest(t)

	// push pop

	// source != dest
	{
		// init data
		r.RunTest(e.RPush, "a", "1").Expect(1)
		r.RunTest(e.RPush, "a", "2").Expect(2)
		r.RunTest(e.RPush, "a", "3").Expect(3)
		r.RunTest(e.RPush, "a", "4").Expect(4)
		r.RunTest(e.LRange, "a", 0, -1).Expect("1", "2", "3", "4")

		r.RunTest(e.RPopLPush, "a", "b").Expect("4")
		r.RunTest(e.LRange, "a", 0, -1).Expect("1", "2", "3")
		r.RunTest(e.LRange, "b", 0, -1).Expect("4")

		r.RunTest(e.RPopLPush, "a", "b").Expect("3")
		r.RunTest(e.LRange, "a", 0, -1).Expect("1", "2")
		r.RunTest(e.LRange, "b", 0, -1).Expect("3", "4")
	}

	// source == dest
	{
		r.as.Nil(e.FlushDB())

		// init data
		r.RunTest(e.RPush, "a", "1").Expect(1)
		r.RunTest(e.RPush, "a", "2").Expect(2)
		r.RunTest(e.RPush, "a", "3").Expect(3)
		r.RunTest(e.RPush, "a", "4").Expect(4)
		r.RunTest(e.LRange, "a", 0, -1).Expect("1", "2", "3", "4")

		r.RunTest(e.RPopLPush, "a", "a").Expect("4")
		r.RunTest(e.LRange, "a", 0, -1).Expect("4", "1", "2", "3")
		r.RunTest(e.RPopLPush, "a", "a").Expect("3")
		r.RunTest(e.LRange, "a", 0, -1).Expect("3", "4", "1", "2")
	}

	{
		r.as.Nil(e.FlushDB())
		// init data
		r.RunTest(e.RPush, "a", "1").Expect(1)
		r.RunTest(e.RPush, "a", "2").Expect(2)
		r.RunTest(e.LRange, "a", 0, -1).Expect("1", "2")

		r.RunTest(e.BRPopLPush, "a", "b", zeroTimeDuration).Expect("2")
		r.RunTest(e.BRPopLPush, "a", "b", zeroTimeDuration).Expect("1")
		r.RunTest(e.LRange, "a", 0, -1).Expect("")
		r.RunTest(e.LRange, "b", 0, -1).Expect("1", "2")

	}
}

func TestListIndexeSetInsert(t *testing.T) {
	r := NewTest(t)

	// index
	r.RunTest(e.LPush, "a", "1").Expect(1)
	r.RunTest(e.LPush, "a", "2").Expect(2)
	r.RunTest(e.LRange, "a", 0, -1).Expect("2", "1")

	r.RunTest(e.LIndex, "a", 0).Expect("2")
	r.RunTest(e.LIndex, "a", -1).Expect("1")
	r.RunTest(e.LIndex, "a", 3).Expect(redis.NullString{})

	// insert
	r.RunTest(e.LInsert, "a", true, "1", "pre-1").Expect(3)
	r.RunTest(e.LRange, "a", 0, -1).Expect("2", "pre-1", "1")

	r.RunTest(e.LInsert, "a", true, "not-exist", "2").Expect(-1)
	r.RunTest(e.LInsert, "not-exist-key", true, "not-exist-value", "2").Expect(0)

	// set
	r.RunTest(e.LRange, "a", 0, -1).Expect("2", "pre-1", "1")
	r.RunTest(e.LSet, "a", 1, "set-value").ExpectSuccess()
	r.RunTest(e.LRange, "a", 0, -1).Expect("2", "set-value", "1")
}

func TestListPushX(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.LPush, "a", "1").Expect(1)
	r.RunTest(e.LPushX, "a", "2").Expect(2)
	r.RunTest(e.LPushX, "b", "2").Expect(0)

	r.RunTest(e.RPush, "c", "1").Expect(1)
	r.RunTest(e.RPushX, "c", "2").Expect(2)
	r.RunTest(e.RPushX, "d", "2").Expect(0)
}

func TestListRem(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.RPush, "a", "morning").Expect(1)
	r.RunTest(e.RPush, "a", "hello").Expect(2)
	r.RunTest(e.RPush, "a", "morning").Expect(3)
	r.RunTest(e.RPush, "a", "hello").Expect(4)
	r.RunTest(e.RPush, "a", "morning").Expect(5)
	r.RunTest(e.LRange, "a", 0, -1).Expect("morning", "hello", "morning", "hello", "morning")

	r.RunTest(e.LRem, "a", 2, "morning").Expect(2)
	r.RunTest(e.LRange, "a", 0, -1).Expect("hello", "hello", "morning")

	r.RunTest(e.LRem, "a", -1, "morning").Expect(1)
	r.RunTest(e.LRange, "a", 0, -1).Expect("hello", "hello")

	r.RunTest(e.LRem, "a", 0, "hello").Expect(2)
	r.RunTest(e.LRange, "a", 0, -1).ExpectSuccess()
}

func TestListTrim(t *testing.T) {
	r := NewTest(t)

	{
		r.as.Nil(e.FlushDB())
		r.RunTest(e.RPush, "a", "1").Expect(1)
		r.RunTest(e.RPush, "a", "2").Expect(2)
		r.RunTest(e.RPush, "a", "3").Expect(3)
		r.RunTest(e.RPush, "a", "4").Expect(4)
		r.RunTest(e.RPush, "a", "5").Expect(5)
		r.RunTest(e.LRange, "a", 0, -1).Expect("1", "2", "3", "4", "5")

		// start 和 stop 都在列表的索引范围之内个
		r.RunTest(e.LTRIM, "a", 1, -1).ExpectSuccess()
		r.RunTest(e.LRange, "a", 0, -1).Expect("2", "3", "4", "5")

		// stop 比列表的最大下标还要大
		r.RunTest(e.LTRIM, "a", 1, 9090).ExpectSuccess()
		r.RunTest(e.LRange, "a", 0, -1).Expect("3", "4", "5")

		// start 和 stop 都比列表的最大下标要大，并且 start < stop
		r.RunTest(e.LTRIM, "a", 90, 9090).ExpectSuccess()
		r.RunTest(e.LRange, "a", 0, -1).Expect(nil)
	}

	{
		r.as.Nil(e.FlushDB())
		r.RunTest(e.RPush, "a", "1").Expect(1)
		r.RunTest(e.RPush, "a", "2").Expect(2)
		r.RunTest(e.RPush, "a", "3").Expect(3)
		r.RunTest(e.RPush, "a", "4").Expect(4)
		r.RunTest(e.RPush, "a", "5").Expect(5)
		r.RunTest(e.LRange, "a", 0, -1).Expect("1", "2", "3", "4", "5")

		// start 和 stop 都比列表的最大下标要大，并且 start > stop
		r.RunTest(e.LTRIM, "a", 9090, 90).ExpectSuccess()
		r.RunTest(e.LRange, "a", 0, -1).Expect(nil)
	}
}
