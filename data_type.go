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

// NullString ...
type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

// KeyType ...
type KeyType string

// KeyType ...
const (
	KeyTypeNone   KeyType = "none"   // key不存在
	KeyTypeString KeyType = "string" // 字符串
	KeyTypeList   KeyType = "list"   // 列表
	KeyTypeSet    KeyType = "set"    // 集合
	KeyTypeZSet   KeyType = "zset"   // 有序集
	KeyTypeHash   KeyType = "hash"   // 哈希表
)

// SortedSet ...
type SortedSet struct {
	Member string
	Score  int
}

// GeoLocation ...
type GeoLocation struct {
	Longitude float64 // 经度
	Latitude  float64 // 纬度
	Member    string
}
