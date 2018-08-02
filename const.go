package redis

import "errors"

const (
	LF byte = 10 // \n
	CR byte = 13 // \r
)

var CRLF = []byte{CR, LF}

var (
	UnSupportRespType = errors.New("unsupported redis resp type")
)
