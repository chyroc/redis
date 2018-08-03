package redis

import (
	"fmt"
	"strconv"
)

// DataType ...
type DataType interface {
	Err() error
	String() string
}

type integer struct {
	err    error
	signed bool
	length int
}

func (r integer) Err() error {
	if r.length <= 0 || (r.signed && r.length > 64) || (!r.signed && r.length > 63) {
		return fmt.Errorf("invalid type. use something like i16 u8. note that u64 is not supported but i64 is")
	}
	return nil
}

func (r integer) String() string {
	if r.signed {
		return "i" + strconv.Itoa(r.length)
	}
	return "u" + strconv.Itoa(r.length)
}

// UnSignedInt ...
func UnSignedInt(length int) DataType {
	return integer{signed: false, length: length}
}

// SignedInt ...
func SignedInt(length int) DataType {
	return integer{signed: true, length: length}
}
