package redis_test

import (
	"fmt"
	"github.com/Chyroc/redis"
	"os"
	"testing"
)

var e *redis.Redis

func TestConnection(t *testing.T) {
	r := conn(t)

	r.RunTest(e.Set, "a", "b").Expect(true)
	r.RunTest(e.Get, "a").Expect("b")
	r.RunTest(e.Select, 2).ExpectSuccess()
	r.RunTest(e.Get, "a").Expect(redis.NullString{})
}

func TestMultiRdisInstance(t *testing.T) {
	r := conn(t)

	_, ok := os.LookupEnv("TRAVIS")
	if !ok {
		t.SkipNow()
	}
	for _, v := range os.Environ() {
		fmt.Printf("%v\n", v)
	}

	r.RunTest(e.Set, "greeting", "Hello from 6379 instance").Expect(true)
	r.RunTest(e.Migrate, "127.0.0.1", 7777, "greeting", 0, 0).Expect(true)
	r.RunTest(e.Exists, "greeting").Expect(false)

	e2, err := redis.Dial("127.0.0.1:6379")
	r.Nil(err)
	x, err := e2.Get("greeting")
	r.Nil(err)
	r.Equal("Hello from 6379 instance", x)
}
