package redis

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
