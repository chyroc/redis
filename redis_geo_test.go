package redis_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Chyroc/redis"
)

func diffLessThan(t *testing.T, a, b, diff float64) {
	if a > b {
		if a-b > diff {
			assert.Fail(t, fmt.Sprintf("%f - %f should less than %f", a, b, diff))
		}
	} else {
		if b-a > diff {
			assert.Fail(t, fmt.Sprintf("%f - %f should less than %f", b, a, diff))
		}
	}
}

func expectGeoLocation(t *testing.T, geo *redis.GeoLocation, expectLongitude, expectLatitude float64) {
	assert.NotNil(t, geo)
	diffLessThan(t, geo.Longitude, expectLongitude, 1)
	diffLessThan(t, geo.Latitude, expectLatitude, 1)
}

func TestGeo(t *testing.T) {
	r := NewTest(t)

	// add
	r.RunTest(e.GeoAdd, "key", redis.GeoLocation{0, 0, "a"}).Expect(1)
	r.RunTest(e.GeoAdd, "key", redis.GeoLocation{180, 80, "b"}).Expect(1)
	r.RunTest(e.GeoAdd, "key", redis.GeoLocation{180, -80, "c"}).Expect(1)

	// pos
	g, err := e.GeoPos("key", "a")
	r.as.Nil(err)
	r.as.Len(g, 1)
	expectGeoLocation(t, g[0], 0, 0)

	g, err = e.GeoPos("key", "b", "c")
	r.as.Nil(err)
	r.as.Len(g, 2)
	expectGeoLocation(t, g[0], 180, 80)
	expectGeoLocation(t, g[1], 180, -80)

	// dist
	r.RunTest(e.GeoDist, "key", "a", "b").ExpectSuccess()
	r.RunTest(e.GeoDist, "key", "a", "b", "km").ExpectSuccess()

	// hash
	r.RunTest(e.GeoHash, "key", "a").Expect("s0000000000")
	r.RunTest(e.GeoHash, "key", "a", "b", "c").Expect("s0000000000", "s0000000000", "s0000000000")
}
