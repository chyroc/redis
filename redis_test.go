package redis_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/chyroc/redis"
	"math/rand"
)

var (
	e                *redis.Redis
	zeroTimeDuration = time.Duration(0)
	zeroMap          = map[string]string{}
)

func NewRedis(t *testing.T) *redis.Redis {
	as := assert.New(t)
	e, err := redis.Dial("127.0.0.1:6379")
	as.Nil(err)
	as.NotNil(e)
	return e
}

func NewTest(t *testing.T) *testRedis {
	as := assert.New(t)
	e = NewRedis(t)
	as.Nil(e.FlushDB())
	return &testRedis{redis: e, t: t, as: as, map_string_string: make(map[string]string), map_string_int: make(map[string]int)}
}

type testRedis struct {
	redis *redis.Redis
	t     *testing.T
	as    *assert.Assertions

	err               error
	number            float64
	str               string
	boo               bool
	null              bool
	duration          *time.Duration
	results           []interface{}
	map_string_string map[string]string
	map_string_int    map[string]int
}

func (r *testRedis) run(fun interface{}, args ...interface{}) {
	ft := reflect.TypeOf(fun)
	fv := reflect.ValueOf(fun)
	r.as.Equal(reflect.Func, ft.Kind())
	r.as.Equal(reflect.Func, fv.Kind())

	atLeastCallInNumber := ft.NumIn()
	if ft.IsVariadic() {
		atLeastCallInNumber--
	}
	if atLeastCallInNumber > len(args) {
		r.as.Fail(fmt.Sprintf("expect at least %d arguments , but got %d: %#v", atLeastCallInNumber, len(args), args))
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
			case redis.SetOption, redis.MigrateOption, redis.ScanOption, redis.GeoLocation, time.Time:
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
			switch x := ithCallRealOut.Interface().(type) {
			case map[string]string:
				r.map_string_string = x
			case map[string]int:
				r.map_string_int = x
			default:
				panic(fmt.Sprintf("invalid map type: %#v\n", x))
			}
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
	r = &testRedis{redis: r.redis, t: r.t, as: r.as}
	r.run(fun, args...)
	return r
}

func (r *testRedis) Expect(expected ...interface{}) *testRedis {
	r.as.Nil(r.err)

	if len(expected) == 0 {
		r.as.Fail("expect at least 1 argument")
	}

	if len(r.results) > 0 {
		r.as.Len(r.results, len(expected))
		for k, v := range r.results {
			switch v.(type) {
			case int64:
				expected[k] = int64(expected[k].(int))
			default:
				break
			}
		}

		r.as.Equal(expected, r.results)
		return r
	}

	switch e := expected[0].(type) {
	case int:
		r.as.Equal(float64(e), r.number)
	case float64:
		r.as.Equal(e, r.number)
	case string:
		r.as.Equal(e, r.str)
	case bool:
		r.as.Equal(e, r.boo)
	case map[string]string:
		r.as.Equal(e, r.map_string_string)
	case map[string]int:
		r.as.Equal(e, r.map_string_int)
	case *time.Duration:
		r.as.Equal(e, r.duration)
	case time.Duration:
		r.as.Equal(e, *r.duration)
	case redis.KeyType:
		r.as.Equal(string(e), r.str)
	case redis.NullString:
		fmt.Printf("result : %#v", r.results)
		r.as.Equal(e.String, r.str)
		r.as.Equal(e.Valid, !r.null)
	case *redis.SortedSet:
		r.as.Equal(e.Score, r.str)
		r.as.Equal(e.Member, !r.null)
	case nil:
		r.as.Nil(r.duration)
	default:
		panic(fmt.Sprintf("invalid data type: %#v", e))
	}

	return r
}

func (r *testRedis) ExpectSuccess() {
	r.as.Nil(r.err)
}

func (r *testRedis) ExpectNull() {
	r.as.Nil(r.err)
	r.as.True(r.null)
	r.as.Empty(r.str)
}

func (r *testRedis) ExpectError(s string) {
	r.as.NotNil(r.err)
	r.as.Equal(s, r.err.Error())
}

func (r *testRedis) ExpectBigger(i int) {
	r.as.Nil(r.err)
	r.as.True(r.number > float64(i))
}

func (r *testRedis) ExpectLess(i interface{}) {
	r.as.Nil(r.err)
	switch v := i.(type) {
	case int:
		r.as.True(r.number <= float64(v))
	case time.Duration:
		r.as.True(*r.duration <= v)
	}
}

func (r *testRedis) ExpectBelong(s ...string) {
	r.as.Nil(r.err)
	for _, v := range s {
		if v == r.str {
			return
		}
	}
	r.as.Fail(fmt.Sprintf("expected %#v contain: %v", s, r.str))
}

func (r *testRedis) ExpectContains(s ...string) {
	r.as.Nil(r.err)
	stringContains(r.t, interfacesToStringSlice(r.results, 0), s)
}

func (r *testRedis) ExpectContainsBy(s ...string) {
	r.as.Nil(r.err)
	stringContains(r.t, s, interfacesToStringSlice(r.results, 0))
}

func (r *testRedis) ExpectSlice(s ...string) {
	r.as.Nil(r.err)

	m := make(map[string]int)
	for _, v := range r.results {
		m[v.(string)]++
	}
	for _, v := range s {
		if m[v] <= 0 {
			r.as.Fail(fmt.Sprintf("%#v != %#v", r.results, s))
		}
		m[v]--
	}
	for _, v := range m {
		if v > 0 {
			r.as.Fail(fmt.Sprintf("%#v != %#v", r.results, s))
		}
	}
}

func (r *testRedis) SetBits(key string, index, result []int) {
	r.as.Equal(len(index), len(result))
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
	r.as.Equal(len(index), len(result))
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
		case reflect.Int:
			var ints []int
			return reflect.ValueOf(&ints).Elem()
		case reflect.Struct:
			switch typ.Elem().Name() {
			case "MigrateOption":
				var ints []redis.MigrateOption
				return reflect.ValueOf(&ints).Elem()
			case "SetOption":
				var ints []redis.SetOption
				return reflect.ValueOf(&ints).Elem()
			case "LimitOption":
				var ints []redis.LimitOption
				return reflect.ValueOf(&ints).Elem()
			case "GeoLocation":
				var ints []redis.GeoLocation
				return reflect.ValueOf(&ints).Elem()
			default:
				panic(fmt.Sprintf("unsupport %v\n", typ.Elem().Name()))
			}
		case reflect.Interface:
			var ints []interface{}
			return reflect.ValueOf(&ints).Elem()
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
		if i+1 < len(args) {
			switch args[i+1].(type) {
			case string:
				var v = reflect.ValueOf(&[]interface{}{}).Elem()
				for k, s := range args[i:] {
					if k%2 == 0 {
						v = reflect.Append(v, reflect.ValueOf(s.(int)))
					} else {
						v = reflect.Append(v, reflect.ValueOf(s.(string)))
					}
				}
				return v
			}
		}

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
	case redis.MigrateOption:
		var v = reflect.ValueOf(&[]redis.MigrateOption{}).Elem()
		for _, s := range args[i:] {
			v = reflect.Append(v, reflect.ValueOf(s.(redis.MigrateOption)))
		}
		return v
	default:
		panic(fmt.Sprintf("expect slice, but got %#v\n", args[i:]))
	}
}

func toInterfaceSlice(typ reflect.Type, value reflect.Value) []interface{} {
	var slice []interface{}

	if typ.Elem().Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

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
		case "SortedSet":
			for _, v := range value.Convert(reflect.TypeOf([]*redis.SortedSet{})).Interface().([]*redis.SortedSet) {
				slice = append(slice, v)
			}
		case "GeoLocation":
			for _, v := range value.Convert(reflect.TypeOf([]*redis.GeoLocation{})).Interface().([]*redis.GeoLocation) {
				slice = append(slice, v)
			}
		default:
			panic(fmt.Sprintf("expect struct , but got %s: %#v", typ.Elem().Name(), value))
		}
	default:
		panic(fmt.Sprintf("expect slice , but got %s: %#v", typ.Elem().Kind(), value))
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
		r.as.Fail(fmt.Sprintf("timeout"))
	}

	r.as.Fail("unkone error")
}

// every ele in b in a slice
func stringContains(t *testing.T, a, b []string) {
	as := assert.New(t)

	m := make(map[string]bool)
	for _, v := range a {
		m[v] = true
	}
	for _, v := range b {
		if !m[v] {
			as.Fail(fmt.Sprintf("%#v should contain %#v", b, a))
		}
	}
}

func randString(n int) string {
	var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// RandString rand string
	rand.Seed(int64(time.Now().Nanosecond()))
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
