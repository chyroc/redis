package redis_test

import (
	"github.com/Chyroc/redis"
	"testing"
)

func TestHash(t *testing.T) {
	r := NewTest(t)

	// hash set
	r.RunTest(e.HSet, "addr", "a", "1").Expect(true)
	r.RunTest(e.HSet, "addr", "b", "1").Expect(true)
	r.RunTest(e.HSet, "addr", "c", "1").Expect(true)
	r.RunTest(e.HSet, "addr", "d", "1").Expect(true)
	r.RunTest(e.HSet, "addr", "d", "1").Expect(false) // exist field

	// hash getall
	r.RunTest(e.HGetALL, "addr").Expect(map[string]string{"a": "1", "b": "1", "c": "1", "d": "1"})
	r.RunTest(e.HGetALL, "not-exist").Expect(zeroMap)

	// hash exists
	r.RunTest(e.HExists, "not-exist", "not-exist").Expect(false)
	r.RunTest(e.HExists, "addr", "not-exist").Expect(false)
	r.RunTest(e.HExists, "addr", "a").Expect(true)

	// hash get
	r.RunTest(e.HGet, "not-exist", "not-exist").Expect(redis.NullString{})
	r.RunTest(e.HGet, "addr", "not-exist").Expect(redis.NullString{})
	r.RunTest(e.HGet, "addr", "a").Expect(redis.NullString{String: "1", Valid: true})

	// hash del
	r.RunTest(e.HDel, "addr", "a").Expect(1)
	r.RunTest(e.HDel, "addr", "not-exist").Expect(0)
	r.RunTest(e.HDel, "addr", "b", "c", "d").Expect(3)
	r.RunTest(e.HGetALL, "addr").Expect(zeroMap)

	// hash incr by
	r.RunTest(e.HSet, "incr", "a", "1").Expect(true)
	r.RunTest(e.HIncrBy, "incr", "a", 10).Expect(11)
	r.RunTest(e.HIncrBy, "incr", "a", -5).Expect(6)

	r.RunTest(e.HSet, "incr", "b", "string").Expect(true)
	r.RunTest(e.HIncrBy, "incr", "b", 10).ExpectError("ERR hash value is not an integer")
	r.RunTest(e.HIncrBy, "incr", "not-exist", 10).Expect(false)

	// hash incr by float
	r.RunTest(e.HSet, "key", "a", "10.5").Expect(true)
	r.RunTest(e.HIncrByFloat, "key", "a", 0.1).Expect(10.6)
	r.RunTest(e.HSet, "key", "b", "5.0e3").Expect(true) // 指数符号
	r.RunTest(e.HIncrByFloat, "key", "b", 2.0e2).Expect(5200)

	// not exist key
	r.RunTest(e.HIncrByFloat, "not-exist", "c", 2.0e2).Expect(200)
	r.RunTest(e.HGetALL, "not-exist").Expect(map[string]string{"c": "200"})

	// not exist field
	r.RunTest(e.HIncrByFloat, "not-exist", "not-exist-field", -1).Expect(-1)
	r.RunTest(e.HGetALL, "not-exist").Expect(map[string]string{"c": "200", "not-exist-field": "-1"})

	// hash keys
	r.RunTest(e.HKeys, "not-exist").Expect("c", "not-exist-field")
	r.RunTest(e.HKeys, "not-exist real").Expect("") // todo fix test

	// hash len
	r.RunTest(e.HLen, "not-exist").Expect(2)
	r.RunTest(e.HLen, "not-exist real").Expect(0)

	// hash multi set and get
	r.RunTest(e.HMSet, "k1", "field1", "v1", "field2", "v2", "field3", "v3").ExpectSuccess()
	r.RunTest(e.HMSet, "k1", "field1", "v1-fix", "field2", "v2-fix", "field4", "v4").ExpectSuccess()
	r.RunTest(e.HMGet, "k1", "field1", "field-not-exist", "field4", "field-not-exist").Expect(redis.NullString{"v1-fix", true}, redis.NullString{}, redis.NullString{"v4", true}, redis.NullString{})

	// hash set nx
	r.RunTest(e.HSetNX, "not-exist-key-nx", "a", "1").Expect(true)
	r.RunTest(e.HSetNX, "not-exist-key-nx", "a", "1").Expect(false)
	r.RunTest(e.HSetNX, "not-exist-key-nx", "b", "1").Expect(true)
	r.RunTest(e.HSetNX, "not-exist-key-nx", "b", "1").Expect(false)

	// hash values
	r.RunTest(e.HVals, "not-exist").Expect("200", "-1")

	// hash string length
	r.RunTest(e.HStrLen, "not-exist", "c").Expect(3)               // 200
	r.RunTest(e.HStrLen, "not-exist", "not-exist-field").Expect(2) // -1

	r.Nil(e.FlushDB())

	// hash scan
	r.RunTest(e.HSet, "addr", "a", "1").Expect(true)
	r.RunTest(e.HSet, "addr", "b", "1").Expect(true)
	r.RunTest(e.HSet, "addr", "c", "1").Expect(true)
	r.RunTest(e.HSet, "addr", "d", "1").Expect(true)
	r.RunTest(e.HSet, "addr", "d1w", "1").Expect(true)
	r.RunTest(e.HSet, "addr", "d2w", "1").Expect(true)

	r.RunTest(e.HScan("addr").ALL).Expect(map[string]string{"a": "1", "b": "1", "c": "1", "d": "1", "d1w": "1", "d2w": "1"})
	r.RunTest(e.HScan("addr", redis.ScanOption{Match: "d?*"}).ALL).Expect(map[string]string{"d1w": "1", "d2w": "1"})

	var kk []int
	var ff []string
	var vv []string
	r.Nil(e.HScan("addr").Each(func(k int, field, value string) error {
		kk = append(kk, k)
		ff = append(ff, field)
		vv = append(vv, value)
		return nil
	}))
	r.Equal([]int{0, 1, 2, 3, 4, 5}, kk)
	stringContains(t, []string{"d", "d1w", "d2w", "a", "b", "c"}, ff)
	r.Equal([]string{"1", "1", "1", "1", "1", "1"}, vv)
}
