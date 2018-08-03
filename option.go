package redis

import "time"

// SetOption ...
type SetOption struct {
	Expire time.Duration
	NX     bool // Only set the key if it does not already exist.
	XX     bool // Only set the key if it already exist.
}

// BitOp ...
type BitOp string

// BitOp 参数
const (
	BitOpAND BitOp = "AND"
	BitOpOR  BitOp = "OR"
	BitOpXOR BitOp = "XOR"
	BitOpNOT BitOp = "NOT"
)

// BitFieldOverflow type
type BitFieldOverflow string

// BitFieldOverflow type
const (
	BitFieldOverflowWrap BitFieldOverflow = "WRAP"
	BitFieldOverflowSat  BitFieldOverflow = "SAT"
	BitFieldOverflowFail BitFieldOverflow = "FAIL"
)
