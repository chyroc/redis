package redis

import "strconv"

// BitField ...
type BitField struct {
	r       *Redis
	key     string
	err     error
	actions []string
}

// Get type offset
//
//   返回指定的二进制位范围
func (b *BitField) Get(typ DataType, offset int) *BitField {
	if b.err != nil {
		return b
	}
	if err := typ.Err(); err != nil {
		b.err = err
		return b
	}
	b.actions = append(b.actions, "GET", typ.String(), strconv.Itoa(offset))
	return b
}

// Set type offset, value
//
//   对指定的二进制位范围进行设置，并返回它的旧值
func (b *BitField) Set(typ DataType, offset, value int) *BitField {
	if b.err != nil {
		return b
	}
	if err := typ.Err(); err != nil {
		b.err = err
		return b
	}
	b.actions = append(b.actions, "SET", typ.String(), strconv.Itoa(offset), strconv.Itoa(value))
	return b
}

// IncrBy type offset increment
//
//   对指定的二进制位范围执行加法操作，并返回它的旧值
func (b *BitField) IncrBy(typ DataType, offset, increment int) *BitField {
	if b.err != nil {
		return b
	}
	if err := typ.Err(); err != nil {
		b.err = err
		return b
	}
	b.actions = append(b.actions, "INCRBY", typ.String(), strconv.Itoa(offset), strconv.Itoa(increment))
	return b
}

// Overflow WARP|SAT|FAIL
//
//   用户可以通过 OVERFLOW 命令以及以下展示的三个参数， 指定 BITFIELD 命令在执行自增或者自减操作时， 碰上向上溢出（overflow）或者向下溢出（underflow）情况时的行为：
//
//     WRAP ： 使用回绕（wrap around）方法处理有符号整数和无符号整数的溢出情况。
//     对于无符号整数来说， 回绕就像使用数值本身与能够被储存的最大无符号整数执行取模计算， 这也是 C 语言的标准行为。
//     对于有符号整数来说， 上溢将导致数字重新从最小的负数开始计算， 而下溢将导致数字重新从最大的正数开始计算。
//     比如说， 如果我们对一个值为 127 的 i8 整数执行加一操作， 那么将得到结果 -128 。
//     SAT ： 使用饱和计算（saturation arithmetic）方法处理溢出， 也即是说， 下溢计算的结果为最小的整数值，
//     而上溢计算的结果为最大的整数值。 举个例子， 如果我们对一个值为 120 的 i8 整数执行加 10 计算，
//     那么命令的结果将为 i8 类型所能储存的最大整数值 127 。 与此相反， 如果一个针对 i8 值的计算造成了下溢，
//     那么这个 i8 值将被设置为 -127 。
//     FAIL ： 在这一模式下， 命令将拒绝执行那些会导致上溢或者下溢情况出现的计算， 并向用户返回空值表示计算未被执行。
//
//   需要注意的是， OVERFLOW 子命令只会对紧随着它之后被执行的 INCRBY 命令产生效果，
//   这一效果将一直持续到与它一同被执行的下一个 OVERFLOW 命令为止。 在默认情况下， INCRBY 命令使用 WRAP 方式来处理溢出计算。
func (b *BitField) Overflow(f BitFieldOverflow) *BitField {
	if b.err != nil {
		return b
	}
	b.actions = append(b.actions, "OVERFLOW", string(f))
	return b
}

// Run ...
func (b *BitField) Run() ([]int64, error) {
	if b.err != nil {
		return nil, b.err
	}
	p := b.r.run(append([]string{"BITFIELD", b.key}, b.actions...)...)
	if p.err != nil {
		return nil, p.err
	}

	var is []int64
	for _, v := range p.replys {
		if v.err != nil {
			return nil, v.err // TODO 检查这是不是真的有err
		}
		is = append(is, v.integer)
	}
	return is, nil
}
