package redis

import (
	"fmt"
	"strconv"
	"time"
)

// https://redis.io/commands#string
// http://redisdoc.com/string/index.html

// Append key value
//
//   如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
//   如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
//   返回值：
//     追加 value 之后， key 中字符串的长度。
func (r *Redis) Append(key, value string) (int, error) {
	return r.run("APPEND", key, value).int()
}

// BitCount key [start] [end]
//
//   计算给定字符串中，被设置为 1 的比特位的数量。
//   一般情况下，给定的整个字符串都会被进行计数，通过指定额外的 start 或 end 参数，可以让计数只在特定的位上进行。
//   start 和 end 参数的设置和 GETRANGE 命令类似，都可以使用负数值： 比如 -1 表示最后一个字节， -2 表示倒数第二个字节，以此类推。
//   不存在的 key 被当成是空字符串来处理，因此对一个不存在的 key 进行 BITCOUNT 操作，结果为 0 。
//
//   时间复杂度： O(N)
//   返回值：
//     被设置为 1 的位的数量。
func (r *Redis) BitCount(key string, startEnd ...int) (int, error) {
	args := []string{"BITCOUNT", key}
	switch len(startEnd) {
	case 0:
	case 1:
		args = append(args, strconv.Itoa(startEnd[0]))
	case 2:
		args = append(args, strconv.Itoa(startEnd[0]), strconv.Itoa(startEnd[1]))
	default:
		return 0, fmt.Errorf("expect get 0, 1 or 2 arguments, but got %v", startEnd)
	}
	return r.run(args...).int()
}

// BitOp operation destkey key [key ...]
//
//   Available since 2.6.0.
//   Time complexity: O(N)
//
//   对一个或多个保存二进制位的字符串 key 进行位元操作，并将结果保存到 destkey 上。
//
//   operation 可以是 AND 、 OR 、 NOT 、 XOR 这四种操作中的任意一种：
//
//     BITOP AND destkey key [key ...] ，对一个或多个 key 求逻辑并，并将结果保存到 destkey 。
//     BITOP OR destkey key [key ...] ，对一个或多个 key 求逻辑或，并将结果保存到 destkey 。
//     BITOP XOR destkey key [key ...] ，对一个或多个 key 求逻辑异或，并将结果保存到 destkey 。
//     BITOP NOT destkey key ，对给定 key 求逻辑非，并将结果保存到 destkey 。
//     除了 NOT 操作之外，其他操作都可以接受一个或多个 key 作为输入。
//
//   处理不同长度的字符串
//
//    当 BITOP 处理不同长度的字符串时，较短的那个字符串所缺少的部分会被看作 0 。
//    空的 key 也被看作是包含 0 的字符串序列。
//
//   返回值：
//     保存到 destkey 的字符串的长度，和输入 key 中最长的字符串长度相等。
func (r *Redis) BitOp(op BitOp, destkey string, keys ...string) (int, error) {
	if len(keys) == 0 {
		return 0, fmt.Errorf("need at least 1 key argument, but got 0")
	}
	switch op {
	case BitOpAND, BitOpXOR, BitOpOR:
		return r.run(append([]string{"BITOP", string(op), destkey}, keys...)...).int()
	case BitOpNOT:
		return r.run("BITOP", "NOT", destkey, keys[0]).int()
	}

	return 0, fmt.Errorf("invalid operation, should be one of: or, and, xor and not")
}

// BitField key [GET type offset] [SET type offset value] [INCRBY type offset increment] [OVERFLOW WRAP|SAT|FAIL]
//
//   Available since 3.2.0.
//   Time complexity: O(1) for each subcommand specified
//
//   BITFIELD 命令可以将一个 Redis 字符串看作是一个由二进制位组成的数组， 并对这个数组中储存的长度不同的整数进行访问 （被储存的整数无需进行对齐）。
//   换句话说， 通过这个命令， 用户可以执行诸如 “对偏移量 1234 上的 5 位长有符号整数进行设置”、 “获取偏移量 4567 上的 31 位长无符号整数”等操作。
//   此外， BITFIELD 命令还可以对指定的整数执行加法操作和减法操作， 并且这些操作可以通过设置妥善地处理计算时出现的溢出情况。
//
//   BITFIELD 命令可以在一次调用中同时对多个位范围进行操作： 它接受一系列待执行的操作作为参数， 并返回一个数组作为回复， 数组中的每个元素就是对应操作的执行结果。
//
//   注意：
//
//     使用 GET 子命令对超出字符串当前范围的二进制位进行访问（包括键不存在的情况）， 超出部分的二进制位的值将被当做是 0 。
//     使用 SET 子命令或者 INCRBY 子命令对超出字符串当前范围的二进制位进行访问将导致字符串被扩大， 被扩大的部分会使用值为 0 的二进制位进行填充。
//     在对字符串进行扩展时， 命令会根据字符串目前已有的最远端二进制位， 计算出执行操作所需的最小长度。
//
//   以下是 BITFIELD 命令支持的子命令：
//
//     GET <type> <offset> —— 返回指定的二进制位范围。
//     SET <type> <offset> <value> —— 对指定的二进制位范围进行设置，并返回它的旧值。
//     INCRBY <type> <offset> <increment> —— 对指定的二进制位范围执行加法操作，并返回它的旧值。用户可以通过向 increment 参数传入负值来实现相应的减法操作。
//
//   除了以上三个子命令之外， 还有一个子命令， 它可以改变之后执行的 INCRBY 子命令在发生溢出情况时的行为：
//
//     OVERFLOW [WRAP|SAT|FAIL]
//
//   当被设置的二进制位范围值为整数时， 用户可以在类型参数的前面添加 i 来表示有符号整数， 或者使用 u 来表示无符号整数。
//   比如说， 我们可以使用 u8 来表示 8 位长的无符号整数， 也可以使用 i16 来表示 16 位长的有符号整数。
//
//   BITFIELD 命令最大支持 64 位长的有符号整数以及 63 位长的无符号整数， 其中无符号整数的 63 位长度限制是由于 Redis 协议目前还无法返回 64 位长的无符号整数而导致的。
func (r *Redis) BitField(key string) *BitField {
	return &BitField{r: r, key: key}
}

// Decr key
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(1)
//
//   将 key 中储存的数字值减一。
//   如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECR 操作。
//   如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
//   本操作的值限制在 64 位(bit)有符号数字表示之内。
//
//   返回值：
//     执行 DECR 命令之后 key 的值。
func (r *Redis) Decr(key string) (int, error) {
	return r.run("DECR", key).int()
}

// DecrBy key decrement
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(1)
//
//   将 key 所储存的值减去减量 decrement 。
//   如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 DECRBY 操作。
//   如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
//   本操作的值限制在 64 位(bit)有符号数字表示之内。
//   关于更多递增(increment) / 递减(decrement)操作的更多信息，请参见 INCR 命令。
//
//   返回值：
//     减去 decrement 之后， key 的值。
func (r *Redis) DecrBy(key string, decrement int) (int, error) {
	return r.run("DECRBY", key, strconv.Itoa(decrement)).int()
}

// Get key
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(1)
//
//   返回 key 所关联的字符串值。
//   如果 key 不存在那么返回特殊值 nil 。
//   假如 key 储存的值不是字符串类型，返回一个错误，因为 GET 只能用于处理字符串值。
//
//   返回值：
//     当 key 不存在时，返回 nil ，否则，返回 key 的值。
//     如果 key 不是字符串类型，那么返回一个错误。
func (r *Redis) Get(key string) (NullString, error) {
	return r.run("GET", key).string()
}

// GetBit key offset
//
//   可用版本： >= 2.2.0
//   时间复杂度： O(1)
//
//   对 key 所储存的字符串值，获取指定偏移量上的位(bit)。
//   当 offset 比字符串值的长度大，或者 key 不存在时，返回 0 。
//
//   返回值：
//     字符串值指定偏移量上的位(bit)。
func (r *Redis) GetBit(key string, offset int) (int, error) {
	return r.run("GETBIT", key, strconv.Itoa(offset)).int()
}

// GetRange key start end
//
//   可用版本： >= 2.4.0
//   时间复杂度： O(N)， N 为要返回的字符串的长度。
//
//   返回 key 中字符串值的子字符串，字符串的截取范围由 start 和 end 两个偏移量决定(包括 start 和 end 在内)。
//   负数偏移量表示从字符串最后开始计数， -1 表示最后一个字符， -2 表示倒数第二个，以此类推。
//   GETRANGE 通过保证子字符串的值域(range)不超过实际字符串的值域来处理超出范围的值域请求。
//
//   复杂度最终由字符串的返回值长度决定，但因为从已有字符串中取出子字符串的操作非常廉价(cheap)，所以对于长度不大的字符串，该操作的复杂度也可看作O(1)。
//
//   返回值：
//     截取得出的子字符串。
func (r *Redis) GetRange(key string, start, end int) (string, error) {
	e, err := r.run("GETRANGE", key, strconv.Itoa(start), strconv.Itoa(end)).string()
	if err != nil {
		return "", err
	}
	return e.String, nil
}

// GetSet key value
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(1)
//
//   将给定 key 的值设为 value ，并返回 key 的旧值(old value)。
//   当 key 存在但不是字符串类型时，返回一个错误。
//
//   返回值：
//     返回给定 key 的旧值。
//     当 key 没有旧值时，也即是， key 不存在时，返回 nil 。
func (r *Redis) GetSet(key, value string) (NullString, error) {
	return r.run("GETSET", key, value).string()
}

// Incr key
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(1)
//
//   将 key 中储存的数字值增一。
//   如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCR 操作。
//   如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
//   本操作的值限制在 64 位(bit)有符号数字表示之内。
//   这是一个针对字符串的操作，因为 Redis 没有专用的整数类型，所以 key 内储存的字符串被解释为十进制 64 位有符号整数来执行 INCR 操作。
//
//   返回值：
//     执行 INCR 命令之后 key 的值。
func (r *Redis) Incr(key string) (int, error) {
	return r.run("INCR", key).int()
}

// IncrBy key
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(1)
//
//   将 key 所储存的值加上增量 increment 。
//   如果 key 不存在，那么 key 的值会先被初始化为 0 ，然后再执行 INCRBY 命令。
//   如果值包含错误的类型，或字符串类型的值不能表示为数字，那么返回一个错误。
//   本操作的值限制在 64 位(bit)有符号数字表示之内。
//   关于递增(increment) / 递减(decrement)操作的更多信息，参见 INCR 命令。
//
//   返回值：
//     加上 increment 之后， key 的值。
func (r *Redis) IncrBy(key string, increment int) (int, error) {
	return r.run("INCRBY", key, strconv.Itoa(increment)).int()
}

// IncrByFloat key increment
//
//   可用版本： >= 2.6.0
//   时间复杂度： O(1)
//
//   为 key 中所储存的值加上浮点数增量 increment 。
//   如果 key 不存在，那么 INCRBYFLOAT 会先将 key 的值设为 0 ，再执行加法操作。
//   如果命令执行成功，那么 key 的值会被更新为（执行加法之后的）新值，并且新值会以字符串的形式返回给调用者。
//   无论是 key 的值，还是增量 increment ，都可以使用像 2.0e7 、 3e5 、 90e-2 那样的指数符号(exponential notation)来表示，但是，执行 INCRBYFLOAT 命令之后的值总是以同样的形式储存，也即是，它们总是由一个数字，一个（可选的）小数点和一个任意位的小数部分组成（比如 3.14 、 69.768 ，诸如此类)，小数部分尾随的 0 会被移除，如果有需要的话，还会将浮点数改为整数（比如 3.0 会被保存成 3 ）。
//   除此之外，无论加法计算所得的浮点数的实际精度有多长， INCRBYFLOAT 的计算结果也最多只能表示小数点的后十七位。
//   当以下任意一个条件发生时，返回一个错误：
//
//   key 的值不是字符串类型(因为 Redis 中的数字和浮点数都以字符串的形式保存，所以它们都属于字符串类型）
//   key 当前的值或者给定的增量 increment 不能解释(parse)为双精度浮点数(double precision floating point number）
//
//   返回值：
//     执行命令之后 key 的值。
func (r *Redis) IncrByFloat(key string, increment float64) (float64, error) {
	p := r.run("INCRBYFLOAT", key, strconv.FormatFloat(increment, 'f', 10, 64))
	if p.err != nil {
		return 0, p.err
	}
	return strconv.ParseFloat(p.str, 64)
}

// MGet ...
//
//   可用版本： >= 1.0.0
//   时间复杂度: O(N) , N 为给定 key 的数量。
//
//   返回所有(一个或多个)给定 key 的值。
//   如果给定的 key 里面，有某个 key 不存在，那么这个 key 返回特殊值 nil 。因此，该命令永不失败。
//
//   返回值：
//     一个包含所有给定 key 的值的列表。
func (r *Redis) MGet(key string, keys ...string) ([]NullString, error) {
	p := r.run(append([]string{"MGET", key}, keys...)...)
	if p.err != nil {
		return nil, p.err
	}

	var ns []NullString
	for _, v := range p.replys {
		if v.err != nil {
			return nil, v.err // TODO 这里真的有error吗
		}
		n, _ := v.string()
		ns = append(ns, n)
	}
	return ns, nil
}

// MSet ...
//   可用版本： >= 1.0.1
//   时间复杂度： O(N)， N 为要设置的 key 数量。
//
//   同时设置一个或多个 key-value 对。
//   如果某个给定 key 已经存在，那么 MSET 会用新值覆盖原来的旧值，如果这不是你所希望的效果，请考虑使用 MSETNX 命令：它只会在所有给定 key 都不存在的情况下进行设置操作。
//   MSET 是一个原子性(atomic)操作，所有给定 key 都会在同一时间内被设置，某些给定 key 被更新而另一些给定 key 没有改变的情况，不可能发生。
//
//   返回值：
//     总是返回 OK (因为 MSET 不可能失败)
func (r *Redis) MSet(key, value string, kvs ...string) error {
	if len(kvs)%2 != 0 {
		return fmt.Errorf("key value pair, but got %d arguments", len(kvs)+2)
	}

	return r.run(append([]string{"MSET", key, value}, kvs...)...).err
}

// MSetNX key value [key value ...]
//
//   同时设置一个或多个 key-value 对，当且仅当所有给定 key 都不存在。
//   即使只有一个给定 key 已存在， MSETNX 也会拒绝执行所有给定 key 的设置操作。
//   MSETNX 是原子性的，因此它可以用作设置多个不同 key 表示不同字段(field)的唯一性逻辑对象(unique logic object)，所有字段要么全被设置，要么全不被设置。
//
//   可用版本：>= 1.0.1
//   时间复杂度：O(N)， N 为要设置的 key 的数量。
//
//   返回值：
//     当所有 key 都成功设置，返回 1 。
//     如果所有给定 key 都设置失败(至少有一个 key 已经存在)，那么返回 0 。
func (r *Redis) MSetNX(key, value string, kvs ...string) (bool, error) {
	if len(kvs)%2 != 0 {
		return false, fmt.Errorf("key value pair, but got %d arguments", len(kvs)+2)
	}

	return r.run(append([]string{"MSETNX", key, value}, kvs...)...).fixBool()
}

// Set key value [expiration EX seconds|PX milliseconds] [NX|XX]
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(1)
//
//   将字符串值 value 关联到 key 。
//   如果 key 已经持有其他值， SET 就覆写旧值，无视类型。
//
//   对于某个原本带有生存时间（TTL）的键来说， 当 SET 命令成功在这个键上执行时， 这个键原有的 TTL 将被清除。
//
//   可选参数
//
//     从 Redis 2.6.12 版本开始， SET 命令的行为可以通过一系列参数来修改：
//
//     EX second ：设置键的过期时间为 second 秒。 SET key value EX second 效果等同于 SETEX key second value 。
//     PX millisecond ：设置键的过期时间为 millisecond 毫秒。 SET key value PX millisecond 效果等同于 PSETEX key millisecond value 。
//     NX ：只在键不存在时，才对键进行设置操作。 SET key value NX 效果等同于 SETNX key value 。
//     XX ：只在键已经存在时，才对键进行设置操作。
//
//   返回值：
//
//     在 Redis 2.6.12 版本以前， SET 命令总是返回 OK 。
//     从 Redis 2.6.12 版本开始， SET 在设置操作成功完成时，才返回 OK 。
//     如果设置了 NX 或者 XX ，但因为条件没达到而造成设置操作未执行，那么命令返回空批量回复（NULL Bulk Reply）。
func (r *Redis) Set(key, value string, options ...SetOption) (bool, error) {
	if len(options) > 1 {
		return false, fmt.Errorf("must have 0 or 1 option")
	}

	args := []string{"SET", key, value}

	if len(options) > 0 {
		option := options[0]
		if option.Expire >= time.Millisecond {
			args = append(args, "PX", strconv.Itoa(int(option.Expire/time.Millisecond)))
		}
		if option.NX && option.XX {
			return false, fmt.Errorf("cannot set NX and XX option at the same time")
		} else if option.NX {
			args = append(args, "NX")
		} else if option.XX {
			args = append(args, "XX")
		}

	}

	p := r.run(args...)
	if p.err != nil {
		return false, p.err
	}

	return p.str == "OK" && !p.null, nil
}

// SetBit key offset value
//
//   可用版本： >= 2.2.0
//   时间复杂度: O(1)
//
//   对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)。
//   位的设置或清除取决于 value 参数，可以是 0 也可以是 1 。
//   当 key 不存在时，自动生成一个新的字符串值。
//   字符串会进行伸展(grown)以确保它可以将 value 保存在指定的偏移量上。当字符串值进行伸展时，空白位置以 0 填充。
//   offset 参数必须大于或等于 0 ，小于 2^32 (bit 映射被限制在 512 MB 之内)。
//
//   对使用大的 offset 的 SETBIT 操作来说，内存分配可能造成 Redis 服务器被阻塞。具体参考 SETRANGE 命令，warning(警告)部分。
//
//   返回值：
//     指定偏移量原来储存的位。
func (r *Redis) SetBit(key string, offset int, SetOrRemove bool) (int, error) {
	return r.run("SETBIT", key, strconv.Itoa(offset), boolToString(SetOrRemove)).int()
}

// SetRange key offset value
//
//   可用版本： >= 2.2.0
//   时间复杂度： 对小(small)的字符串，平摊复杂度O(1)。(关于什么字符串是”小”的，请参考 APPEND 命令)
//
//   用 value 参数覆写(overwrite)给定 key 所储存的字符串值，从偏移量 offset 开始。
//   不存在的 key 当作空白字符串处理。
//   SETRANGE 命令会确保字符串足够长以便将 value 设置在指定的偏移量上，如果给定 key 原来储存的字符串长度比偏移量小(比如字符串只有 5 个字符长，但你设置的 offset 是 10 )，那么原字符和偏移量之间的空白将用零字节(zerobytes, "\x00" )来填充。
//   注意你能使用的最大偏移量是 2^29-1(536870911) ，因为 Redis 字符串的大小被限制在 512 兆(megabytes)以内。如果你需要使用比这更大的空间，你可以使用多个 key 。
//
//   当生成一个很长的字符串时，Redis 需要分配内存空间，该操作有时候可能会造成服务器阻塞(block)。在2010年的Macbook Pro上，设置偏移量为 536870911(512MB 内存分配)，耗费约 300 毫秒， 设置偏移量为 134217728(128MB 内存分配)，耗费约 80 毫秒，设置偏移量 33554432(32MB 内存分配)，耗费约 30 毫秒，设置偏移量为 8388608(8MB 内存分配)，耗费约 8 毫秒。 注意若首次内存分配成功之后，再对同一个 key 调用 SETRANGE 操作，无须再重新内存。
//   否则为O(M)， M 为 value 参数的长度。
//
//   返回值：
//     被 SETRANGE 修改之后，字符串的长度。
func (r *Redis) SetRange(key string, offset int, value string) (int, error) {
	return r.run("SETRANGE", key, strconv.Itoa(offset), value).int()
}

// StrLen key
//
//   可用版本： >= 2.2.0
//   复杂度： O(1)
//
//   返回 key 所储存的字符串值的长度。
//   当 key 储存的不是字符串值时，返回一个错误。
//
//   返回值：
//     字符串值的长度。
//     当 key 不存在时，返回 0 。
func (r *Redis) StrLen(key string) (int, error) {
	return r.run("STRLEN", key).int()
}

// un-impl command: PSETEX / SETEX / SETNX
