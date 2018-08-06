package redis_test

import "testing"

func TestHyperLogLog(t *testing.T) {
	r := NewTest(t)

	// add
	r.RunTest(e.PFAdd, "a").Expect(true)
	r.RunTest(e.PFAdd, "a").Expect(false)
	r.RunTest(e.PFAdd, "a", "1", "2", "3").Expect(true)
	r.RunTest(e.PFAdd, "b", "1", "diff").Expect(true)

	// count
	r.RunTest(e.PFCount, "not-exist").Expect(0)
	r.RunTest(e.PFCount, "a").Expect(3)
	r.RunTest(e.PFCount, "b").Expect(2)
	r.RunTest(e.PFCount, "a", "b").Expect(4)

	// merge
	r.RunTest(e.PFMerge, "dest", "a", "b").ExpectSuccess()
	r.RunTest(e.PFCount, "dest").Expect(4)
}
