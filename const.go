package redis

import "errors"

// error
var (
	ErrUnSupportRespType = errors.New("unsupported redis resp type")
	ErrEmptyCommand      = errors.New("empty command")
	ErrKeyNotExist       = errors.New("key not exist")
	ErrIteratorEnd       = errors.New("iterator end")
	ErrInvalidBitOp      = errors.New("invalid operation, should be one of: or, and, xor and not")
)
