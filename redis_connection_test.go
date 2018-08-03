package redis_test

import "testing"

func TestConnection(t *testing.T) {
	r, as := conn(t)

	p := r.Set("a", "b")
	as.Nil(p.Err())

	p = r.Get("a")
	as.Nil(p.Err())
	as.Equal("b", p.String())

	p = r.Select(2)
	as.Nil(p.Err())

	p = r.Get("a")
	as.Nil(p.Err())
	as.Equal("", p.String())
	as.True(p.Null())
}
