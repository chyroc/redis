package redis_test

import (
	"github.com/Chyroc/redis"
	"testing"
)

func TestSortedSetAdd(t *testing.T) {
	r := NewTest(t)

	// add card
	r.RunTest(e.ZAdd, "a", 10, "1").Expect(1)
	r.RunTest(e.ZCard, "a").Expect(1)
	r.RunTest(e.ZAdd, "a", 9, "2", 8, "3").Expect(2)
	r.RunTest(e.ZCard, "a").Expect(3)

	// count
	r.RunTest(e.ZCount, "a", 0, 1).Expect(0)
	r.RunTest(e.ZCount, "a", 0, 8).Expect(1)
	r.RunTest(e.ZCount, "a", 0, 9).Expect(2)
	r.RunTest(e.ZCount, "a", 0, 10).Expect(3)

	// card
	r.RunTest(e.ZCard, "not-exist").Expect(0)

	// range revrange
	r.RunTest(e.ZRange, "a", 0, -1, true).Expect(&redis.SortedSet{"3", 8}, &redis.SortedSet{"2", 9}, &redis.SortedSet{"1", 10})
	r.RunTest(e.ZRange, "a", 0, -1, false).Expect(&redis.SortedSet{"3", 0}, &redis.SortedSet{"2", 0}, &redis.SortedSet{"1", 0})
	r.RunTest(e.ZRevRange, "a", 0, -1, true).Expect(&redis.SortedSet{"1", 10}, &redis.SortedSet{"2", 9}, &redis.SortedSet{"3", 8})
	r.RunTest(e.ZRevRange, "a", 0, -1, false).Expect(&redis.SortedSet{"1", 0}, &redis.SortedSet{"2", 0}, &redis.SortedSet{"3", 0})

	// change score
	r.RunTest(e.ZAdd, "a", 10, "1").Expect(0)
	r.RunTest(e.ZAdd, "a", 7, "1").Expect(0)
	r.RunTest(e.ZRange, "a", 0, -1, true).Expect(&redis.SortedSet{"1", 7}, &redis.SortedSet{"3", 8}, &redis.SortedSet{"2", 9})
	r.RunTest(e.ZRevRange, "a", 0, -1, true).Expect(&redis.SortedSet{"2", 9}, &redis.SortedSet{"3", 8}, &redis.SortedSet{"1", 7})

	// incrby
	r.RunTest(e.ZIncrBy, "a", -5, "2").Expect(4)
	r.RunTest(e.ZRange, "a", 0, -1, true).Expect(&redis.SortedSet{"2", 4}, &redis.SortedSet{"1", 7}, &redis.SortedSet{"3", 8})
	r.RunTest(e.ZIncrBy, "a", 5, "2").Expect(9)

	// rem
	r.RunTest(e.ZRem, "a", "1").Expect(1)
	r.RunTest(e.ZRem, "a", "2", "3").Expect(2)
	r.RunTest(e.ZRem, "a", "not-exist").Expect(0)
}

func TestSortedSetRangeByScore(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.ZAdd, "a", 10, "1", 9, "2", 8, "3").Expect(3)

	// range by score
	r.RunTest(e.ZRangeByScore, "a", "-inf", "+inf", true).Expect(&redis.SortedSet{"3", 8}, &redis.SortedSet{"2", 9}, &redis.SortedSet{"1", 10})
	r.RunTest(e.ZRangeByScore, "a", "-inf", "+inf", false).Expect(&redis.SortedSet{"3", 0}, &redis.SortedSet{"2", 0}, &redis.SortedSet{"1", 0})
	r.RunTest(e.ZRangeByScore, "a", "-inf", "9", true).Expect(&redis.SortedSet{"3", 8}, &redis.SortedSet{"2", 9})
	r.RunTest(e.ZRangeByScore, "a", "-inf", "(9", true).Expect(&redis.SortedSet{"3", 8})

	// rev range by score
	r.RunTest(e.ZRevRangeByScore, "a", "+inf", "-inf", true).Expect(&redis.SortedSet{"1", 10}, &redis.SortedSet{"2", 9}, &redis.SortedSet{"3", 8})
	r.RunTest(e.ZRevRangeByScore, "a", "+inf", "-inf", false).Expect(&redis.SortedSet{"1", 0}, &redis.SortedSet{"2", 0}, &redis.SortedSet{"3", 0})
	r.RunTest(e.ZRevRangeByScore, "a", "9", "-inf", true).Expect(&redis.SortedSet{"2", 9}, &redis.SortedSet{"3", 8})
	r.RunTest(e.ZRevRangeByScore, "a", "(9", "-inf", true).Expect(&redis.SortedSet{"3", 8})
}

func TestSortedSetRank(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.ZAdd, "a", 10, "1", 9, "2", 8, "3").Expect(3)

	// rank
	r.RunTest(e.ZRank, "a", "1").Expect(2)
	r.RunTest(e.ZRank, "a", "2").Expect(1)
	r.RunTest(e.ZRank, "a", "3").Expect(0)
	r.RunTest(e.ZRank, "a", "not-exist").ExpectError(redis.ErrKeyNotExist.Error())
	r.RunTest(e.ZRank, "not-exist", "not-exist").ExpectError(redis.ErrKeyNotExist.Error())

	// rev rank
	r.RunTest(e.ZRevRank, "a", "1").Expect(0)
	r.RunTest(e.ZRevRank, "a", "2").Expect(1)
	r.RunTest(e.ZRevRank, "a", "3").Expect(2)
	r.RunTest(e.ZRevRank, "a", "not-exist").ExpectError(redis.ErrKeyNotExist.Error())
	r.RunTest(e.ZRevRank, "not-exist", "not-exist").ExpectError(redis.ErrKeyNotExist.Error())

}

func TestSortedSetRemRange(t *testing.T) {
	r := NewTest(t)

	// rem range by rank
	r.RunTest(e.ZAdd, "a", 10, "1", 9, "2", 8, "3").Expect(3)
	r.RunTest(e.ZRemRangeByRank, "a", 0, 1).Expect(2)
	r.RunTest(e.ZRemRangeByRank, "a", 0, 1).Expect(1)
	r.RunTest(e.ZRemRangeByRank, "a", 0, 1).Expect(0)

	// rem range by score
	r.RunTest(e.ZAdd, "a", 10, "1", 9, "2", 8, "3").Expect(3)
	r.RunTest(e.ZRemRangeByScore, "a", "7", "8").Expect(1)
	r.RunTest(e.ZRemRangeByScore, "a", "-inf", "+inf").Expect(2)
	r.RunTest(e.ZRemRangeByScore, "a", "-inf", "+inf").Expect(0)
}

func TestSortedSetScore(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.ZAdd, "a", 10, "1", 9, "2", 8, "3").Expect(3)
	r.RunTest(e.ZRange, "a", 0, -1, true).Expect(&redis.SortedSet{"3", 8}, &redis.SortedSet{"2", 9}, &redis.SortedSet{"1", 10})

	r.RunTest(e.ZScore, "a", "1").Expect(10)
	r.RunTest(e.ZScore, "a", "2").Expect(9)
	r.RunTest(e.ZScore, "a", "3").Expect(8)
}

func TestSortedSetScan(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.ZAdd, "a", 10, "1", 9, "2", 8, "3").Expect(3)
	r.RunTest(e.ZRange, "a", 0, -1, true).Expect(&redis.SortedSet{"3", 8}, &redis.SortedSet{"2", 9}, &redis.SortedSet{"1", 10})

	r.RunTest(e.ZScan("a").ALL).Expect(&redis.SortedSet{"3", 8}, &redis.SortedSet{"2", 9}, &redis.SortedSet{"1", 10})

	var vv []*redis.SortedSet
	r.as.Nil(e.ZScan("a").Each(func(k int, v *redis.SortedSet) error {
		vv = append(vv, v)
		return nil
	}))
	r.as.Equal([]*redis.SortedSet{{"3", 8}, {"2", 9}, {"1", 10}}, vv)
}

func TestSortedSetRangeByLex(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.ZAdd, "a", 2, "a", 2, "b", 2, "c").Expect(3)

	// range by lex
	// lex count
	data := [][]interface{}{
		{"-", "+", 3, []string{"a", "b", "c"}},
		{"-", "[a", 1, []string{"a"}},
		{"-", "(a", 0, []string{}},
		{"-", "[b", 2, []string{"a", "b"}},
		{"-", "(b", 1, []string{"a"}},
		{"-", "[c", 3, []string{"a", "b", "c"}},
		{"-", "(c", 2, []string{"a", "b"}},
		{"-", "[d", 3, []string{"a", "b", "c"}},
		{"-", "(d", 3, []string{"a", "b", "c"}},
		{"(a", "[b", 1, []string{"b"}},
	}
	for _, vv := range data {
		r.RunTest(e.ZRangeByLex, "a", vv[0].(string), vv[1].(string)).ExpectSlice(vv[3].([]string)...)
		r.RunTest(e.ZLexCount, "a", vv[0].(string), vv[1].(string)).Expect(vv[2].(int))
	}

	// remove range by lex
	for _, vv := range data {
		r.RunTest(e.FlushDB).ExpectSuccess()
		r.RunTest(e.ZAdd, "a", 2, "a", 2, "b", 2, "c").Expect(3)
		r.RunTest(e.ZRemRangeByLex, "a", vv[0].(string), vv[1].(string)).Expect(vv[2].(int))
	}
}
