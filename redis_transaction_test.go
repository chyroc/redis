package redis_test

import "testing"

func TestTransaction(t *testing.T) {
	r := NewTest(t)

	r.RunTest(e.Multi).ExpectSuccess()
	r.RunTest(e.Set, "a", "1").ExpectSuccess()
	r.RunTest(e.Set, "b", "2").ExpectSuccess()
	r.RunTest(e.Incr, "a").ExpectSuccess()
	r.RunTest(e.Scan().ALL)
	r.RunTest(e.HSet, "hash", "ha", "1").ExpectSuccess()
	r.RunTest(e.HSet, "hash", "hb", "2").ExpectSuccess()
	r.RunTest(e.HScan("hash").ALL)
	r.RunTest(e.Exec).ExpectSuccess()
}
