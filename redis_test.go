package redis_test

import (
	"fmt"
	"github.com/Chyroc/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testRedis struct {
	*redis.Redis
	*testing.T
	*assert.Assertions

	err     error
	number  float64
	str     string
	boo     bool
	null    bool
	results []interface{}
}

func (r *testRedis) RunTest(fun interface{}, args ...interface{}) *testRedis {
	r.err = nil
	r.number = 0
	r.str = ""
	r.boo = false
	r.null = false
	r.results = nil

	switch f := fun.(type) {
	case func(string, ...int) (int, error):
		var ingeter int
		ingeter, r.err = f(args[0].(string), interfacesToIntSlice(args, 1)...)
		r.number = float64(ingeter)
	case func(string, string) (int, error):
		var ingeter int
		ingeter, r.err = f(args[0].(string), args[1].(string))
		r.number = float64(ingeter)
	case func(op redis.BitOp, destkey string, keys ...string) (int, error):
		var ingeter int
		ingeter, r.err = f(args[0].(redis.BitOp), args[1].(string), interfacesToStringSlice(args, 2)...)
		r.number = float64(ingeter)
	case func() ([]int64, error):
		ints, err := f()
		r.err = err
		for _, v := range ints {
			r.results = append(r.results, v)
		}
	case func(string) (int, error):
		var ingeter int
		ingeter, r.err = f(args[0].(string))
		r.number = float64(ingeter)
	case func(string, int) (int, error):
		var ingeter int
		ingeter, r.err = f(args[0].(string), args[1].(int))
		r.number = float64(ingeter)
	case func(string) (string, error):
		r.str, r.err = f(args[0].(string))
	case func(string, int, int) (string, error):
		r.str, r.err = f(args[0].(string), args[1].(int), args[2].(int))
	case func(string, string) (string, error):
		r.str, r.err = f(args[0].(string), args[1].(string))
	case func(string, float64) (float64, error):
		switch args[1].(type) {
		case int:
			r.number, r.err = f(args[0].(string), float64(args[1].(int)))
		default:
			r.number, r.err = f(args[0].(string), args[1].(float64))
		}
	case func(string, ...string) ([]redis.NullString, error):
		var ns []redis.NullString
		ns, r.err = f(args[0].(string), interfacesToStringSlice(args, 1)...)
		for _, v := range ns {
			r.results = append(r.results, v)
		}
	case func(key, value string, options ...redis.SetOption) (bool, error):
		if len(args) > 2 {
			r.boo, r.err = f(args[0].(string), args[1].(string), args[2].(redis.SetOption))
		} else {
			r.boo, r.err = f(args[0].(string), args[1].(string))
		}
	case func(key, value string, kvs ...string) (bool, error):
		r.boo, r.err = f(args[0].(string), args[1].(string), interfacesToStringSlice(args, 2)...)
	case func(key string, offset int, value string) (int, error):
		var integer int
		integer, r.err = f(args[0].(string), args[1].(int), args[2].(string))
		r.number = float64(integer)
	case func(key string, offset int, SetOrRemove bool) (int, error):
		var integer int
		integer, r.err = f(args[0].(string), args[1].(int), args[2].(bool))
		r.number = float64(integer)
	case func(key string) (bool, error):
		r.boo, r.err = f(args[0].(string))
	case func(key string, member ...string) (int, error):
		var integer int
		integer, r.err = f(args[0].(string), interfacesToStringSlice(args, 1)...)
		r.number = float64(integer)
	case func(key string) (redis.NullString, error):
		var ns redis.NullString
		ns, r.err = f(args[0].(string))
		r.str = ns.String
		r.null = !ns.Valid
	case func(key, value string) (redis.NullString, error):
		var ns redis.NullString
		ns, r.err = f(args[0].(string), args[1].(string))
		r.str = ns.String
		r.null = !ns.Valid
	case func(key string, seconds int) (bool, error):
		r.boo, r.err = f(args[0].(string), args[1].(int))
	case func(key string, t time.Time) (bool, error):
		r.boo, r.err = f(args[0].(string), args[1].(time.Time))
	case func(pattern string) ([]string, error):
		var s []string
		s, r.err = f(args[0].(string))
		for _, v := range s {
			r.results = append(r.results, v)
		}
	case func() (int, error):
		var integer int
		integer, r.err = f()
		r.number = float64(integer)
	case func() (string, error):
		r.str, r.err = f()
	default:
		panic("un support function")
	}

	return r
}

func (r *testRedis) equal(expected, actual interface{}) {
	switch actual.(type) {
	case int64:
		r.Equal(int64(expected.(int)), actual)
	default:
		r.Equal(expected, actual)
	}
}

func (r *testRedis) Expect(expected ...interface{}) {
	r.Nil(r.err)

	if len(r.results) > 0 {
		r.Len(r.results, len(expected))
		for k, v := range r.results {
			switch v.(type) {
			case int64:
				expected[k] = int64(expected[k].(int))
			default:
				break
			}
		}

		r.Equal(expected, r.results)

		return
	}

	switch e := expected[0].(type) {
	case int:
		r.Equal(float64(e), r.number)
	case float64:
		r.Equal(e, r.number)
	case string:
		r.Equal(e, r.str)
	case bool:
		r.Equal(e, r.boo)
	}
}

func (r *testRedis) ExpectNull() {
	r.Nil(r.err)
	r.True(r.null)
	r.Empty(r.str)
}

func (r *testRedis) ExpectError(s string) {
	r.Error(fmt.Errorf(s))
}

func (r *testRedis) ExpectBigger(i int) {
	r.Nil(r.err)
	r.True(r.number > float64(i))
}

func (r *testRedis) equalInt64s(ints []int64, expected ...int) {
	r.Equal(len(expected), len(ints))
	for k := range ints {
		r.Equal(int64(expected[k]), ints[k])
	}
}

func (r *testRedis) equalInts(ints []int, expected ...int) {
	r.Equal(len(expected), len(ints))
	for k := range ints {
		r.Equal(expected[k], ints[k])
	}
}

func conn(t *testing.T) (*testRedis, *assert.Assertions) {
	as := assert.New(t)

	r, err := redis.Dial("127.0.0.1:6379")
	as.Nil(err)
	as.NotNil(r)

	as.Nil(r.FlushDB())

	return &testRedis{Redis: r, T: t, Assertions: as}, as
}

func setbits(t *testing.T, r *testRedis, key string, index, result []int) {
	as := assert.New(t)
	as.Equal(len(index), len(result))

	for k := range index {
		c := false
		if result[k] == 1 {
			c = true
		}
		//r.RunTest(r.SetBit, index[k], c).Expect(1)
		_, err := r.SetBit(key, index[k], c)
		as.Nil(err)
	}
	getbits(t, r, key, index, result)
}

func getbits(t *testing.T, r *testRedis, key string, index, result []int) {
	as := assert.New(t)
	as.Equal(len(index), len(result))

	for k := range index {
		r.RunTest(r.GetBit, key, index[k]).Expect(result[k])
	}
}

func interfacesToIntSlice(args []interface{}, startIndex int) []int {
	var is []int
	for k, v := range args {
		if k < startIndex {
			continue
		}
		is = append(is, v.(int))
	}
	return is
}

// startIndex: 1 ~ len
func interfacesToStringSlice(args []interface{}, startIndex int) []string {
	var str []string
	for k, v := range args {
		if k < startIndex {
			continue
		}
		str = append(str, v.(string))
	}
	return str
}
