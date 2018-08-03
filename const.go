package redis

import "errors"

// redis protocol split string
const (
	LF byte = 10 // \n
	CR byte = 13 // \r
)

// CRLF ...
var CRLF = []byte{CR, LF}

// error
var (
	ErrUnSupportRespType = errors.New("unsupported redis resp type")
	ErrEmptyCommand      = errors.New("empty command")
	ErrNull              = errors.New("empty string")
	ErrKeyNotExist       = errors.New("key not exist")
)
