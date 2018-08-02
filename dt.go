package redis

type nullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}
