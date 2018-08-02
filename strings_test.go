package redis_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Chyroc/redis"
	"time"
)

func conn(t *testing.T) *redis.Redis {
	as := assert.New(t)

	r, err := redis.Dial("127.0.0.1:6379")
	as.Nil(err)
	as.NotNil(r)

	return r
}

func TestStringGetSet(t *testing.T) {
	as := assert.New(t)
	r := conn(t)
	as.Nil(r.FlushDB().Err())

	{
		// get and set
		as.Nil(r.Set("key", "hello").Err())
		p := r.Get("key")
		as.Nil(p.Err())
		as.False(p.Null())
		as.Equal("hello", p.String())
	}

	{
		// expire
		as.Nil(r.Set("key-with-expire-time", "hello", redis.SetOption{Expire: time.Millisecond * 1000}).Err())
		p := r.Get("key-with-expire-time")
		as.Nil(p.Err())
		as.False(p.Null())
		as.Equal("hello", p.String())

		time.Sleep(time.Millisecond * 1000)

		p = r.Get("key-with-expire-time")
		as.Nil(p.Err())
		as.True(p.Null())
	}

	{
		// nx set id not exist
		p := r.Set("not-exists-key", "hello", redis.SetOption{NX: true})
		as.Nil(p.Err())
		as.False(p.Null())
		as.Equal("OK", p.String())

		p = r.Set("not-exists-key", "new hello", redis.SetOption{NX: true})
		as.Nil(p.Err())
		as.True(p.Null())

		p = r.Get("not-exists-key")
		as.Nil(p.Err())
		as.Equal("hello", p.String())
	}

	{
		// xx only set when already exist
		p := r.Set("exists-key", "hello", redis.SetOption{XX: true})
		as.Nil(p.Err())
		as.Empty(p.String())
		as.True(p.Null())

		p = r.Set("exists-key", "value")
		as.Nil(p.Err())

		p = r.Set("exists-key", "new hello", redis.SetOption{XX: true})
		as.Nil(p.Err())
		as.Equal("OK", p.String())
		as.False(p.Null())

		p = r.Get("exists-key")
		as.Nil(p.Err())
		as.Equal("new hello", p.String())
	}

	// TODO test lock with expire + nx xx
}
