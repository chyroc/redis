package redis_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Chyroc/redis"
)

var (
	zeroTimeDuration = time.Duration(0)
	zeroMap          = map[string]string{}
)

func NewTest(t *testing.T) *testRedis {
	as := assert.New(t)

	var err error
	e, err = redis.Dial("127.0.0.1:6379")
	as.Nil(err)
	as.NotNil(e)

	as.Nil(e.FlushDB())

	return &testRedis{redis: e, t: t, Assertions: as, mapp: make(map[string]string)}
}

type testRedis struct {
	redis *redis.Redis
	t     *testing.T
	*assert.Assertions

	err      error
	number   float64
	str      string
	boo      bool
	null     bool
	duration *time.Duration
	results  []interface{}
	mapp     map[string]string
}

func (r *testRedis) run(fun interface{}, args ...interface{}) {
	ft := reflect.TypeOf(fun)
	fv := reflect.ValueOf(fun)
	r.Equal(reflect.Func, ft.Kind())
	r.Equal(reflect.Func, fv.Kind())

	atLeastCallInNumber := ft.NumIn()
	if ft.IsVariadic() {
		atLeastCallInNumber--
	}
	if atLeastCallInNumber > len(args) {
		r.Fail(fmt.Sprintf("expect at least %d arguments , but got %d: %#v", atLeastCallInNumber, len(args), args))
	}

	var in []reflect.Value
	for i := 0; i < ft.NumIn(); i++ {
		ithCallInType := ft.In(i)

		switch ithCallInType.Kind() {
		case reflect.String:
			switch args[i].(type) {
			case redis.BitOp:
				in = append(in, reflect.ValueOf(args[i].(redis.BitOp)))
			default:
				in = append(in, reflect.ValueOf(args[i].(string)))
			}
		case reflect.Bool:
			in = append(in, reflect.ValueOf(args[i].(bool)))
		case reflect.Int:
			in = append(in, reflect.ValueOf(args[i].(int)))
		case reflect.Int64:
			switch l := args[i].(type) {
			case int, int64, time.Duration:
				in = append(in, reflect.ValueOf(l))
			default:
				panic(fmt.Sprintf("unsupprt %v\n", l))
			}
		case reflect.Float64:
			switch l := args[i].(type) {
			case int:
				in = append(in, reflect.ValueOf(float64(l)))
			default:
				in = append(in, reflect.ValueOf(args[i].(float64)))
			}
		case reflect.Struct:
			switch l := args[i].(type) {
			case redis.SetOption, redis.MigrateOption, redis.ScanOption, time.Time:
				in = append(in, reflect.ValueOf(l))
			default:
				panic(fmt.Sprintf("unsupport %v: %#v\n", ithCallInType.Kind(), args[i]))
			}
		case reflect.Slice:
			in = append(in, interfaceSliceToReflectValue(ithCallInType, args, i))
			break
		default:
			panic(fmt.Sprintf("unsupport %v\n", ithCallInType.Kind()))
		}
	}

	var out []reflect.Value
	if ft.IsVariadic() {
		if len(in) < ft.NumIn() {
			in = append(in, reflect.ValueOf(nil))
		}
		out = fv.CallSlice(in)
	} else {
		out = fv.Call(in)
	}

	for i := 0; i < ft.NumOut(); i++ {
		ithCallRealOut := out[i] // 比input多的
		ithCallOutType := ft.Out(i)

		switch ithCallRealOut.Kind() {
		case reflect.Int:
			r.number = float64(ithCallRealOut.Int())
		case reflect.Float64:
			r.number = ithCallRealOut.Float()
		case reflect.Bool:
			r.boo = ithCallRealOut.Bool()
		case reflect.String:
			r.str = ithCallRealOut.String()
		case reflect.Interface:
			if ithCallRealOut.IsNil() {
				continue
			}
			switch l := ithCallRealOut.Interface().(type) {
			case error:
				r.err = l
			default:
				panic(fmt.Sprintf("unsupport interafce: [%v: %#v]", ithCallRealOut.Elem().Kind(), ithCallRealOut))
			}
		case reflect.Struct:
			switch l := ithCallRealOut.Interface().(type) {
			case redis.NullString:
				r.str = l.String
				r.null = !l.Valid
			default:
				panic(fmt.Sprintf("unsupport %v: %#v\n", ithCallRealOut.Kind(), ithCallRealOut.Interface()))
			}
		case reflect.Slice:
			r.results = toInterfaceSlice(ithCallOutType, ithCallRealOut)
		case reflect.Map:
			r.mapp = ithCallRealOut.Interface().(map[string]string)
		case reflect.Ptr:
			if !ithCallRealOut.IsNil() {
				t := ithCallRealOut.Elem().Interface().(time.Duration)
				r.duration = &t
			}
		default:
			panic(fmt.Sprintf("unsupport %v\n", ithCallRealOut.Kind()))
		}
	}
}

func (r *testRedis) RunTest(fun interface{}, args ...interface{}) *testRedis {
	r = &testRedis{redis: r.redis, t: r.t, Assertions: r.Assertions}
	r.run(fun, args...)
	return r
}

func (r *testRedis) Expect(expected ...interface{}) *testRedis {
	r.Nil(r.err)

	if len(expected) == 0 {
		r.Fail("expect at least 1 argument")
	}

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

		return r
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
	case map[string]string:
		r.Equal(e, r.mapp)
	case *time.Duration:
		r.Equal(e, r.duration)
	case time.Duration:
		r.Equal(e, *r.duration)
	case redis.KeyType:
		r.Equal(string(e), r.str)
	case redis.NullString:
		r.Equal(e.String, r.str)
		r.Equal(e.Valid, !r.null)
	case nil:
		r.Nil(r.duration)
	default:
		panic(fmt.Sprintf("invalid data type: %#v", e))
	}

	return r
}

func (r *testRedis) ExpectSuccess() {
	r.Nil(r.err)
}

func (r *testRedis) ExpectNull() {
	r.Nil(r.err)
	r.True(r.null)
	r.Empty(r.str)
}

func (r *testRedis) ExpectError(s string) {
	r.NotNil(r.err)
	r.Equal(s, r.err.Error())
}

func (r *testRedis) ExpectBigger(i int) {
	r.Nil(r.err)
	r.True(r.number > float64(i))
}

func (r *testRedis) ExpectLess(i interface{}) {
	r.Nil(r.err)
	switch v := i.(type) {
	case int:
		r.True(r.number <= float64(v))
	case time.Duration:
		r.True(*r.duration <= v)
	}
}

func (r *testRedis) ExpectBelong(s ...string) {
	r.Nil(r.err)
	for _, v := range s {
		if v == r.str {
			return
		}
	}
	r.Fail(fmt.Sprintf("expected %#v contain: %v", s, r.str))
}

func (r *testRedis) ExpectContains(s ...string) {
	r.Nil(r.err)
	stringContains(r.t, interfacesToStringSlice(r.results, 0), s)
}

func (r *testRedis) SetBits(key string, index, result []int) {
	r.Equal(len(index), len(result))
	for k := range index {
		c := false
		if result[k] == 1 {
			c = true
		}
		r.RunTest(e.SetBit, key, index[k], c).ExpectSuccess()
	}
	//r.GetBits(key, index, result)
}

func (r *testRedis) GetBits(key string, index, result []int) {
	r.Equal(len(index), len(result))
	for k := range index {
		r.RunTest(e.GetBit, key, index[k]).Expect(result[k])
	}
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

func interfaceSliceToReflectValue(typ reflect.Type, args []interface{}, i int) reflect.Value {
	if i >= len(args) {
		switch typ.Elem().Kind() {
		case reflect.String:
			var str []string
			return reflect.ValueOf(&str).Elem()
		default:
			if typ.Elem().ConvertibleTo(reflect.TypeOf(redis.SetOption{})) {
				var str []redis.SetOption
				return reflect.ValueOf(&str).Elem()
			}
		}
		panic(fmt.Sprintf("unsupport zero value: %v %v", typ.Kind(), typ.Elem().Kind()))
	}

	switch args[i].(type) {
	case int:
		var v = reflect.ValueOf(&[]int{}).Elem()
		for _, s := range args[i:] {
			v = reflect.Append(v, reflect.ValueOf(s.(int)))
		}
		return v
	case string:
		var v = reflect.ValueOf(&[]string{}).Elem()
		for _, s := range args[i:] {
			v = reflect.Append(v, reflect.ValueOf(s.(string)))
		}
		return v
	case redis.SetOption:
		var v = reflect.ValueOf(&[]redis.SetOption{}).Elem()
		for _, s := range args[i:] {
			v = reflect.Append(v, reflect.ValueOf(s.(redis.SetOption)))
		}
		return v
	default:
		panic(fmt.Sprintf("expect slice, but got %#v\n", args[i:]))
	}
}

func toInterfaceSlice(typ reflect.Type, value reflect.Value) []interface{} {
	var slice []interface{}

	switch typ.Elem().Kind() {
	case reflect.String:
		for _, v := range value.Convert(reflect.TypeOf([]string{})).Interface().([]string) {
			slice = append(slice, v)
		}
	case reflect.Int64:
		for _, v := range value.Convert(reflect.TypeOf([]int64{})).Interface().([]int64) {
			slice = append(slice, v)
		}
	case reflect.Int:
		for _, v := range value.Convert(reflect.TypeOf([]int{})).Interface().([]int) {
			slice = append(slice, v)
		}
	case reflect.Struct:
		switch typ.Elem().Name() {
		case "NullString":
			for _, v := range value.Convert(reflect.TypeOf([]redis.NullString{})).Interface().([]redis.NullString) {
				slice = append(slice, v)
			}
		default:
			panic(fmt.Sprintf("expect slice , but got %#v", value))
		}
	default:
		panic(fmt.Sprintf("expect slice , but got %#v", value))
	}

	return slice
}

func (r *testRedis) TestTimeout(f func(), timeout time.Duration) {
	done := make(chan bool)
	go func() {
		defer func() {
			done <- true
		}()

		f()
	}()

	select {
	case <-done:
		return
	case <-time.After(timeout):
		r.Fail(fmt.Sprintf("timeout"))
	}

	r.Fail("unkone error")
}
