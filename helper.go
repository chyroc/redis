package redis

func boolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
