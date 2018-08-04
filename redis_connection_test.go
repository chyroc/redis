package redis_test

import (
	"github.com/Chyroc/redis"
	"testing"
	"fmt"
	"os"
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
	//r := conn(t)

	for _, v := range os.Environ() {
		fmt.Printf("%v\n", v)
	}
}
