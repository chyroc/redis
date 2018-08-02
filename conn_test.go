package redis_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Chyroc/redis"
)

func TestConn(t *testing.T) {
	as := assert.New(t)

	r, err := redis.Dial("127.0.0.1:6379")
	as.Nil(err)
	as.NotNil(r)
}
