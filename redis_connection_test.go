package redis_test

import (
	"github.com/Chyroc/redis"
	"testing"
)

func TestConnection(t *testing.T) {
	r, as := conn(t)

	// as.Nil(r.Set("a", "b").Err())
	r.RunTest(r.Set, "a", "b").Expect(true)
	r.RunTest(r.Get, "a").Expect("b")
	as.Nil(r.Select(2))
	r.RunTest(r.Get, "a").Expect(redis.NullString{})
}
