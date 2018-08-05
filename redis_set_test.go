package redis_test

import "testing"

func TestSetAddMember(t *testing.T) {
	r := NewTest(t)

	// add
	r.RunTest(e.SAdd, "a", "1").Expect(1)
	r.RunTest(e.SAdd, "a", "1").Expect(0)
	r.RunTest(e.SAdd, "a", "2", "3", "4").Expect(3)
	r.RunTest(e.SAdd, "b", "b-2", "3", "4").Expect(3)
	r.RunTest(e.SAdd, "c", "c-2", "c-3", "4").Expect(3)
	r.RunTest(e.SAdd, "d", "d").Expect(1)

	// card
	r.RunTest(e.SCard, "a").Expect(4)

	// member
	r.RunTest(e.SMembers, "a").ExpectSlice("1", "2", "3", "4")
	r.RunTest(e.SMembers, "b").ExpectSlice("b-2", "3", "4")
	r.RunTest(e.SMembers, "c").ExpectSlice("c-2", "c-3", "4")
	r.RunTest(e.SMembers, "d").ExpectSlice("d")

	// diff
	r.RunTest(e.SDiff, "a").ExpectSlice("1", "2", "3", "4")
	r.RunTest(e.SDiff, "a", "b").ExpectSlice("1", "2")
	r.RunTest(e.SDiff, "a", "c").ExpectSlice("1", "2", "3")
	r.RunTest(e.SDiff, "b", "c").ExpectSlice("3", "b-2")

	// diff store
	r.RunTest(e.SDiffStore, "a-b", "a", "b").Expect(2)
	r.RunTest(e.SDiffStore, "a-c", "a", "c").Expect(3)
	r.RunTest(e.SDiffStore, "b-c", "b", "c").Expect(2)
	r.RunTest(e.SMembers, "a-b").ExpectSlice("1", "2")
	r.RunTest(e.SMembers, "a-c").ExpectSlice("1", "2", "3")
	r.RunTest(e.SMembers, "b-c").ExpectSlice("3", "b-2")

	// inter
	r.RunTest(e.SInter, "a").ExpectSlice("1", "2", "3", "4")
	r.RunTest(e.SInter, "a", "b").ExpectSlice("3", "4")
	r.RunTest(e.SInter, "a", "c").ExpectSlice("4")
	r.RunTest(e.SInter, "b", "c").ExpectSlice("4")
	r.RunTest(e.SInter, "b", "c", "d").ExpectSlice()

	// inter store
	r.RunTest(e.SInterStore, "a-b", "a", "b").Expect(2)
	r.RunTest(e.SInterStore, "a-c", "a", "c").Expect(1)
	r.RunTest(e.SInterStore, "b-c", "b", "c").Expect(1)
	r.RunTest(e.SInterStore, "b-c-d", "b", "c", "d").ExpectSlice()
	r.RunTest(e.SMembers, "a-b").ExpectSlice("3", "4")
	r.RunTest(e.SMembers, "a-c").ExpectSlice("4")
	r.RunTest(e.SMembers, "b-c").ExpectSlice("4")
	r.RunTest(e.SMembers, "b-c-d").ExpectSlice()

	// is member
	r.RunTest(e.SIsMember, "a", "1").Expect(true)
	r.RunTest(e.SIsMember, "a", "not-exist").Expect(false)
	r.RunTest(e.SIsMember, "not-exist", "not-exist").Expect(false)
}

func TestSetMove(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.SAdd, "a", "1", "2").Expect(2)
	r.RunTest(e.SMembers, "a").ExpectSlice("1", "2")

	r.RunTest(e.SMove, "a", "dest", "1").Expect(true)
	r.RunTest(e.SMembers, "a").ExpectSlice("2")
	r.RunTest(e.SMembers, "dest").ExpectSlice("1")

	r.RunTest(e.SMove, "not-exist", "dest", "1").Expect(false)
}

func TestSetPopRandMember(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.SAdd, "a", "1", "2", "3").Expect(3)
	r.RunTest(e.SRandMember, "a").ExpectContainsBy("1", "2", "3")
	r.RunTest(e.SRandMember, "a").ExpectContainsBy("1", "2", "3")
	r.RunTest(e.SRandMember, "a", 2).ExpectContainsBy("1", "2", "3")
	r.RunTest(e.SRandMember, "a", -2).ExpectContainsBy("1", "2", "3")
	r.RunTest(e.SRandMember, "a", 3).ExpectSlice("1", "2", "3")
	r.RunTest(e.SRandMember, "a", 9090).ExpectSlice("1", "2", "3")

	r.as.Contains([]string{"1", "2", "3"}, r.RunTest(e.SPop, "a").str)
	r.as.Len(r.RunTest(e.SMembers, "a").results, 2)
}

func TestSetRandMember(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.SAdd, "a", "1", "2", "3").Expect(3)
	r.as.Contains([]string{"1", "2", "3"}, r.RunTest(e.SPop, "a").str)
	r.as.Len(r.RunTest(e.SMembers, "a").results, 2)
}

/*
# 添加元素

redis> SADD fruit apple banana cherry
(integer) 3

# 只给定 key 参数，返回一个随机元素

redis> SRANDMEMBER fruit
"cherry"

redis> SRANDMEMBER fruit
"apple"

# 给定 3 为 count 参数，返回 3 个随机元素
# 每个随机元素都不相同

redis> SRANDMEMBER fruit 3
1) "apple"
2) "banana"
3) "cherry"

# 给定 -3 为 count 参数，返回 3 个随机元素
# 元素可能会重复出现多次

redis> SRANDMEMBER fruit -3
1) "banana"
2) "cherry"
3) "apple"

redis> SRANDMEMBER fruit -3
1) "apple"
2) "apple"
3) "cherry"

# 如果 count 是整数，且大于等于集合基数，那么返回整个集合

redis> SRANDMEMBER fruit 10
1) "apple"
2) "banana"
3) "cherry"

# 如果 count 是负数，且 count 的绝对值大于集合的基数
# 那么返回的数组的长度为 count 的绝对值

redis> SRANDMEMBER fruit -10
1) "banana"
2) "apple"
3) "banana"
4) "cherry"
5) "apple"
6) "apple"
7) "cherry"
8) "apple"
9) "apple"
10) "banana"

# SRANDMEMBER 并不会修改集合内容

redis> SMEMBERS fruit
1) "apple"
2) "cherry"
3) "banana"

# 集合为空时返回 nil 或者空数组

redis> SRANDMEMBER not-exists
(nil)

redis> SRANDMEMBER not-eixsts 10
(empty list or set)
*/

func TestSetRem(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.SAdd, "a", "1", "2", "3").Expect(3)
	r.RunTest(e.SRem, "a", "1").Expect(1)
	r.RunTest(e.SRem, "a", "not-exist").Expect(0)
	r.RunTest(e.SRem, "a", "2", "3", "not-exist").Expect(2)
}

func TestSetUnion(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.SAdd, "a", "1", "2", "3").Expect(3)
	r.RunTest(e.SAdd, "b", "b-1").Expect(1)
	r.RunTest(e.SAdd, "c", "1", "c-2").Expect(2)

	r.RunTest(e.SUnion, "a", "b").ExpectSlice("1", "2", "3", "b-1")
	r.RunTest(e.SUnion, "a", "b", "c").ExpectSlice("1", "2", "3", "b-1", "c-2")

	r.RunTest(e.SUnionStore, "a-b", "a", "b").Expect(4)
	r.RunTest(e.SUnionStore, "a-b-c", "a", "b", "c").Expect(5)
	r.RunTest(e.SMembers, "a-b").ExpectSlice("1", "2", "3", "b-1")
	r.RunTest(e.SMembers, "a-b-c").ExpectSlice("1", "2", "3", "b-1", "c-2")
}

func TestSetScan(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.SAdd, "a", "1", "2", "3").Expect(3)

	r.RunTest(e.SScan("a").ALL).ExpectSlice("1", "2", "3")
	var vv []string
	r.as.Nil(e.SScan("a").Each(func(k int, v string) error {
		vv = append(vv, v)
		return nil
	}))
	r.as.Equal([]string{"1", "2", "3"}, vv)
}
