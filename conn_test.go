package redis_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Chyroc/redis"
)

func conn(t *testing.T) *redis.Redis {
	as := assert.New(t)

	r, err := redis.Dial("127.0.0.1:6379")
	as.Nil(err)
	as.NotNil(r)

	return r
}

func TestConn(t *testing.T) {
	as := assert.New(t)

	t.Run("", func(t *testing.T) {
		r := conn(t)
		as.Nil(r.Set("k", "v").Err())
		getReply := r.Get("k")
		as.Nil(getReply.Err())
		as.Equal("v", getReply.String())
	})
}
