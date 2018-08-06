package redis_test

import (
	"testing"
	"fmt"
)

// fixme read lock and write lock
func TestPubSubSubUnSub(t *testing.T) {
	t.SkipNow()
	r := NewTest(t)

	e.Subscribe("a")

	r.RunTest(e.UnSubscribe, "a").Expect(1)

}

func TestPubSubPSubscribe(t *testing.T) {
	r := NewTest(t)
	e2 := NewRedis(t)

	c, err := e.PSubscribe("chan*")
	r.as.Nil(err)

	var expected []string
	go func() {
		for i := 0; i < 100; i++ {
			msg := fmt.Sprintf("msg-%d", i)
			expected = append(expected, msg)
			r.RunTest(e2.Publish, "chan", msg).Expect(1)
		}
	}()

	var s []string
	for i := 0; i < 100; i++ {
		x := <-c
		r.as.Nil(x.Err)
		s = append(s, x.Message)
	}

	r.as.Equal(expected, s)
}

func TestPubSubSubscribe(t *testing.T) {
	r := NewTest(t)
	e2 := NewRedis(t)

	c, err := e.Subscribe("chan")
	r.as.Nil(err)

	var expected []string
	go func() {
		for i := 0; i < 100; i++ {
			msg := fmt.Sprintf("msg-%d", i)
			expected = append(expected, msg)
			r.RunTest(e2.Publish, "chan", msg).Expect(1)
		}
	}()

	var s []string
	for i := 0; i < 100; i++ {
		x := <-c
		r.as.Nil(x.Err)
		s = append(s, x.Message)
	}

	r.as.Equal(expected, s)
}

func TestPubSubNumSubscribe(t *testing.T) {
	r := NewTest(t)

	k := randString(10)
	r.RunTest(e.Pubsub().NumSubscribe, k).Expect(map[string]int{k: 0})
	k1, k2 := randString(10), randString(32)
	r.RunTest(e.Pubsub().NumSubscribe, k1, k2).Expect(map[string]int{k1: 0, k2: 0})

	k4, k5 := randString(10), randString(10)
	_, err := NewRedis(t).Subscribe(k4)
	r.as.Nil(err)
	r.RunTest(e.Pubsub().NumSubscribe, k4, k5).Expect(map[string]int{k4: 1, k5: 0})
	_, err = NewRedis(t).Subscribe(k5)
	r.as.Nil(err)
	r.RunTest(e.Pubsub().NumSubscribe, k4, k5).Expect(map[string]int{k4: 1, k5: 1})
}

func TestPubSubNumPattern(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.Pubsub().NumPattern).Expect(0)
	NewRedis(t).PSubscribe("a*")
	r.RunTest(e.Pubsub().NumPattern).Expect(1)
	NewRedis(t).PSubscribe("b*")
	r.RunTest(e.Pubsub().NumPattern).Expect(2)
	NewRedis(t).PSubscribe("c*")
	r.RunTest(e.Pubsub().NumPattern).Expect(3)
}
