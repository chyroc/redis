package redis_test

import (
	"github.com/Chyroc/redis"
	"os"
	"testing"
)

var e *redis.Redis

func TestConnection(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.Set, "a", "b").Expect(true)
	r.RunTest(e.Get, "a").Expect("b")
	r.RunTest(e.Select, 2).ExpectSuccess()
	r.RunTest(e.Get, "a").Expect(redis.NullString{})
}

func TestMultiRdisInstance(t *testing.T) {
	r := NewTest(t)

	if _, ok := os.LookupEnv("TRAVIS"); !ok {
		t.SkipNow()
	}

	r.RunTest(e.Set, "greeting", "Hello from 6379 instance").Expect(true)
	r.RunTest(e.Migrate, "127.0.0.1", 7777, "greeting", 0, 0).ExpectSuccess()
	r.RunTest(e.Exists, "greeting").Expect(false)

	e2, err := redis.Dial("127.0.0.1:7777")
	r.Nil(err)
	x, err := e2.Get("greeting")
	r.Nil(err)
	r.Equal(redis.NullString{String: "Hello from 6379 instance", Valid: true}, x)
}
