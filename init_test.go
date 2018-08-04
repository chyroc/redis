package redis_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
