package redis_test

import (
	"github.com/chyroc/redis"
	"os"
	"sync"
	"testing"
	"time"
)

func TestConnection(t *testing.T) {
	r := NewTest(t)

	// select
	r.RunTest(e.Set, "a", "b").Expect(true)
	r.RunTest(e.Get, "a").Expect("b")
	r.RunTest(e.Select, 2).ExpectSuccess()
	r.RunTest(e.Get, "a").Expect(redis.NullString{})

	// ping echo
	r.RunTest(e.Ping).Expect("PONG")
	r.RunTest(e.Echo, "message").Expect("message")

	// quit
	r.RunTest(e.Quit).ExpectSuccess()
	r.RunTest(e.Get, "a").ExpectError("EOF")
}

func TestMultiRdisInstance(t *testing.T) {
	r := NewTest(t)

	if _, ok := os.LookupEnv("TRAVIS"); !ok {
		t.SkipNow()
	}

	r.RunTest(e.Set, "greeting", "Hello from 6379 instance").Expect(true)
	r.RunTest(e.Migrate, "127.0.0.1", 7777, "greeting", 0, zeroTimeDuration).ExpectSuccess()
	r.RunTest(e.Exists, "greeting").Expect(false)

	e2, err := redis.Dial("127.0.0.1:7777")
	r.as.Nil(err)
	x, err := e2.Get("greeting")
	r.as.Nil(err)
	r.as.Equal(redis.NullString{String: "Hello from 6379 instance", Valid: true}, x)
}

func TestLock(t *testing.T) {
	r := NewTest(t)
	r.TestTimeout(func() {
		wg := new(sync.WaitGroup)
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				r.RunTest(e.Set, "k2", "2").Expect(true)
			}()
		}
		wg.Wait()
	}, time.Second*2)
}
