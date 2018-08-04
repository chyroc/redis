package redis_test

// every ele in b in a slice
func stringContains(a, b []string) bool {
	m := make(map[string]bool)
	for _, v := range a {
		m[v] = true
	}
	for _, v := range b {
		if !m[v] {
			return false
		}
	}
	return true
}
