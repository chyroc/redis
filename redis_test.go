package redis_test

import (
	"github.com/Chyroc/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

func conn(t *testing.T) (*redis.Redis, *assert.Assertions) {
	as := assert.New(t)

	r, err := redis.Dial("127.0.0.1:6379")
	as.Nil(err)
	as.NotNil(r)

	as.Nil(r.FlushDB().Err())

	return r, as
}

func setbits(t *testing.T, r *redis.Redis, key string, index, result []int) {
	as := assert.New(t)
	as.Equal(len(index), len(result))

	for k := range index {
		c := false
		if result[k] == 1 {
			c = true
		}
		p := r.SetBit(key, index[k], c)
		as.Nil(p.Err())
	}
	getbits(t, r, key, index, result)
}

func getbits(t *testing.T, r *redis.Redis, key string, index, result []int) {
	as := assert.New(t)
	as.Equal(len(index), len(result))

	for k := range index {
		p := r.GetBit(key, index[k])
		as.Nil(p.Err())
		as.Equal(result[k], p.Integer())
	}
}
