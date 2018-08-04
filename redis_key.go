package redis

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Del key [key ...]
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(N)， N 为被删除的 key 的数量。
//
//   删除给定的一个或多个 key 。
//   不存在的 key 会被忽略。
//   删除单个字符串类型的 key ，时间复杂度为O(1)。
//   删除单个列表、集合、有序集合或哈希表类型的 key ，时间复杂度为O(M)， M 为以上数据结构内的元素数量。
//
//   返回值：
//     被删除 key 的数量
func (r *Redis) Del(key string, keys ...string) (int, error) {
	return r.run(append([]string{"DEL", key}, keys...)...).int()
}

// Dump key
//
//   可用版本： >= 2.6.0
//   时间复杂度：
//     查找给定键的复杂度为 O(1) ，对键进行序列化的复杂度为 O(N*M) ，其中 N 是构成 key 的 Redis 对象的数量，而 M 则是这些对象的平均大小。
//     如果序列化的对象是比较小的字符串，那么复杂度为 O(1) 。
//
//   序列化给定 key ，并返回被序列化的值，使用 RESTORE 命令可以将这个值反序列化为 Redis 键。
//   序列化生成的值有以下几个特点：
//
//   它带有 64 位的校验和，用于检测错误， RESTORE 在进行反序列化之前会先检查校验和。
//   值的编码格式和 RDB 文件保持一致。
//   RDB 版本会被编码在序列化值当中，如果因为 Redis 的版本不同造成 RDB 格式不兼容，那么 Redis 会拒绝对这个值进行反序列化操作。
//   序列化的值不包括任何生存时间信息。
//
//   返回值：
//     如果 key 不存在，那么返回 nil 。
//     否则，返回序列化之后的值。
func (r *Redis) Dump(key string) (NullString, error) {
	return r.run("DUMP", key).string()
}

// Exists key
//
//   可用版本：>= 1.0.0
//   时间复杂度： O(1)
//
//   检查给定 key 是否存在。
//
//   返回值：
//     若 key 存在，返回 1 ，否则返回 0
func (r *Redis) Exists(key string) (bool, error) {
	return r.run("EXISTS", key).fixBool()
}

// Expire key seconds
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(1)
//
//   为给定 key 设置生存时间，当 key 过期时(生存时间为 0 )，它会被自动删除。
//   在 Redis 中，带有生存时间的 key 被称为『易失的』(volatile)。
//   生存时间可以通过使用 DEL 命令来删除整个 key 来移除，或者被 SET 和 GETSET 命令覆写(overwrite)，这意味着，如果一个命令只是修改(alter)一个带生存时间的 key 的值而不是用一个新的 key 值来代替(replace)它的话，那么生存时间不会被改变。
//
//   比如说，对一个 key 执行 INCR 命令，对一个列表进行 LPUSH 命令，或者对一个哈希表执行 HSET 命令，这类操作都不会修改 key 本身的生存时间。
//
//   另一方面，如果使用 RENAME 对一个 key 进行改名，那么改名后的 key 的生存时间和改名前一样。
//
//   RENAME 命令的另一种可能是，尝试将一个带生存时间的 key 改名成另一个带生存时间的 another_key ，这时旧的 another_key (以及它的生存时间)会被删除，然后旧的 key 会改名为 another_key ，因此，新的 another_key 的生存时间也和原本的 key 一样。
//
//   使用 PERSIST 命令可以在不删除 key 的情况下，移除 key 的生存时间，让 key 重新成为一个『持久的』(persistent) key 。
//
//   更新生存时间
//
//     可以对一个已经带有生存时间的 key 执行 EXPIRE 命令，新指定的生存时间会取代旧的生存时间。
//
//   过期时间的精确度
//
//     在 Redis 2.4 版本中，过期时间的延迟在 1 秒钟之内 —— 也即是，就算 key 已经过期，但它还是可能在过期之后一秒钟之内被访问到，而在新的 Redis 2.6 版本中，延迟被降低到 1 毫秒之内。
//
//     Redis 2.1.3 之前的不同之处
//
//     在 Redis 2.1.3 之前的版本中，修改一个带有生存时间的 key 会导致整个 key 被删除，这一行为是受当时复制(replication)层的限制而作出的，现在这一限制已经被修复。
//
//   返回值：
//     设置成功返回 1 。
//     当 key 不存在或者不能为 key 设置生存时间时(比如在低于 2.1.3 版本的 Redis 中你尝试更新 key 的生存时间)，返回 0 。
//
//   注意：这里没有实现EXPIRE命令
func (r *Redis) Expire(key string, t time.Duration) (bool, error) {
	return r.run("PEXPIRE", key, durationToMillisecond(t)).fixBool()
}

// ExpireAt key timestamp
//
//   可用版本： >= 1.2.0
//   时间复杂度： O(1)
//
//   EXPIREAT 的作用和 EXPIRE 类似，都用于为 key 设置生存时间。
//
//   不同在于 EXPIREAT 命令接受的时间参数是 UNIX 时间戳(unix timestamp)。
//
//   返回值：
//     如果生存时间设置成功，返回 1 。
//     当 key 不存在或没办法设置生存时间，返回 0 。
//
//   注意：这里没有实现EXPIREAT方法
func (r *Redis) ExpireAt(key string, t time.Time) (bool, error) {
	return r.run("PEXPIREAT", key, strconv.Itoa(int(t.UnixNano()/int64(time.Millisecond)))).fixBool()
}

// Keys pattern
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(N)， N 为数据库中 key 的数量。
//
//   查找所有符合给定模式 pattern 的 key 。
//
//     KEYS * 匹配数据库中所有 key 。
//     KEYS h?llo 匹配 hello ， hallo 和 hxllo 等。
//     KEYS h*llo 匹配 hllo 和 heeeeello 等。
//     KEYS h[ae]llo 匹配 hello 和 hallo ，但不匹配 hillo 。
//     特殊符号用 \ 隔开
//
//   KEYS 的速度非常快，但在一个大的数据库中使用它仍然可能造成性能问题，如果你需要从一个数据集中查找特定的 key ，你最好还是用 Redis 的集合结构(set)来代替。
//
//   返回值：
//     符合给定模式的 key 列表。
func (r *Redis) Keys(pattern string) ([]string, error) {
	return r.run("KEYS", pattern).fixStringSlice()
}

// Migrate host port key destination-db timeout [COPY] [REPLACE]
//
//   可用版本： >= 2.6.0
//   时间复杂度：
//     这个命令在源实例上实际执行 DUMP 命令和 DEL 命令，在目标实例执行 RESTORE 命令，查看以上命令的文档可以看到详细的复杂度说明。
//     key 数据在两个实例之间传输的复杂度为 O(N) 。
//
//   将 key 原子性地从当前实例传送到目标实例的指定数据库上，一旦传送成功， key 保证会出现在目标实例上，而当前实例上的 key 会被删除。
//
//   这个命令是一个原子操作，它在执行的时候会阻塞进行迁移的两个实例，直到以下任意结果发生：迁移成功，迁移失败，等待超时。
//
//   命令的内部实现是这样的：它在当前实例对给定 key 执行 DUMP 命令 ，将它序列化，然后传送到目标实例，目标实例再使用 RESTORE 对数据进行反序列化，并将反序列化所得的数据添加到数据库中；当前实例就像目标实例的客户端那样，只要看到 RESTORE 命令返回 OK ，它就会调用 DEL 删除自己数据库上的 key 。
//
//   timeout 参数以毫秒为格式，指定当前实例和目标实例进行沟通的最大间隔时间。这说明操作并不一定要在 timeout 毫秒内完成，只是说数据传送的时间不能超过这个 timeout 数。
//
//   MIGRATE 命令需要在给定的时间规定内完成 IO 操作。如果在传送数据时发生 IO 错误，或者达到了超时时间，那么命令会停止执行，并返回一个特殊的错误： IOERR 。
//
//   当 IOERR 出现时，有以下两种可能：
//
//     key 可能存在于两个实例
//     key 可能只存在于当前实例
//     唯一不可能发生的情况就是丢失 key ，因此，如果一个客户端执行 MIGRATE 命令，并且不幸遇上 IOERR 错误，那么这个客户端唯一要做的就是检查自己数据库上的 key 是否已经被正确地删除。
//
//   如果有其他错误发生，那么 MIGRATE 保证 key 只会出现在当前实例中。（当然，目标实例的给定数据库上可能有和 key 同名的键，不过这和 MIGRATE 命令没有关系）。
//
//   可选项：
//
//     COPY ：不移除源实例上的 key 。
//     REPLACE ：替换目标实例上已存在的 key 。
//
//   返回值：
//     迁移成功时返回 OK ，否则返回相应的错误。
func (r *Redis) Migrate(host string, port int, key string, destinationDB int, timeout time.Duration, options ...MigrateOption) error {
	// MIGRATE 127.0.0.1 7777 greeting 0 1000
	// TODO test

	if len(options) > 1 {
		return fmt.Errorf("must have 0 or 1 option")
	}

	args := []string{"MIGRATE", host, strconv.Itoa(port), key, strconv.Itoa(destinationDB), strconv.Itoa(int(timeout / time.Millisecond))}
	if len(options) > 0 {
		if options[0].Copy {
			args = append(args, "COPY")
		}
		if options[0].Replace {
			args = append(args, "REPLACE")
		}
	}
	return r.run(args...).err
}

// Move key db
//
//   可用版本：>= 1.0.0
//   时间复杂度：O(1)
//
//   将当前数据库的 key 移动到给定的数据库 db 当中。
//
//   如果当前数据库(源数据库)和给定数据库(目标数据库)有相同名字的给定 key ，或者 key 不存在于当前数据库，那么 MOVE 没有任何效果。
//
//   因此，也可以利用这一特性，将 MOVE 当作锁(locking)原语(primitive)。
//
//   返回值：
//     移动成功返回 1 ，失败则返回 0 。
func (r *Redis) Move(key string, db int) (bool, error) {
	return r.run("MOVE", key, strconv.Itoa(db)).fixBool()
}

// Object subcommand [arguments [arguments]]
//
//   可用版本： >= 2.2.3
//   时间复杂度： O(1)
//
//   OBJECT 命令允许从内部察看给定 key 的 Redis 对象。
//
//   它通常用在除错(debugging)或者了解为了节省空间而对 key 使用特殊编码的情况。
//   当将Redis用作缓存程序时，你也可以通过 OBJECT 命令中的信息，决定 key 的驱逐策略(eviction policies)。
//   OBJECT 命令有多个子命令：
//
//     * OBJECT REFCOUNT <key> 返回给定 key 引用所储存的值的次数。此命令主要用于除错。
//     *  OBJECT ENCODING <key> 返回给定 key 锁储存的值所使用的内部表示(representation)。
//     * OBJECT IDLETIME <key> 返回给定 key 自储存以来的空闲时间(idle， 没有被读取也没有被写入)，以秒为单位。
//
//   对象可以以多种方式编码：
//
//     * 字符串可以被编码为 raw (一般字符串)或 int (为了节约内存，Redis 会将字符串表示的 64 位有符号整数编码为整数来进行储存）。
//    *   列表可以被编码为 ziplist 或 linkedlist 。 ziplist 是为节约大小较小的列表空间而作的特殊表示。
//    *   集合可以被编码为 intset 或者 hashtable 。 intset 是只储存数字的小集合的特殊表示。
//    *   哈希表可以编码为 zipmap 或者 hashtable 。 zipmap 是小哈希表的特殊表示。
//    *  有序集合可以被编码为 ziplist 或者 skiplist 格式。 ziplist 用于表示小的有序集合，而 skiplist 则用于表示任何大小的有序集合。
//
//   假如你做了什么让 Redis 没办法再使用节省空间的编码时(比如将一个只有 1 个元素的集合扩展为一个有 100 万个元素的集合)，特殊编码类型(specially encoded types)会自动转换成通用类型(general type)。
//
//   返回值：
//     REFCOUNT 和 IDLETIME 返回数字。
//     ENCODING 返回相应的编码类型。
func (r *Redis) Object(key string) *Object {
	return &Object{r, key}
}

// Persist key
//
//   可用版本： >= 2.2.0
//   时间复杂度： O(1)
//
//   移除给定 key 的生存时间，将这个 key 从『易失的』(带生存时间 key )转换成『持久的』(一个不带生存时间、永不过期的 key )。
//
//   返回值：
//     当生存时间移除成功时，返回 1 .
//     如果 key 不存在或 key 没有设置生存时间，返回 0 。
func (r *Redis) Persist(key string) (bool, error) {
	return r.run("PERSIST", key).fixBool()
}

// RandomKey ...
//
//   可用版本：>= 1.0.0
//   时间复杂度： O(1)
//
//   从当前数据库中随机返回(不删除)一个 key 。
//
//   返回值：
//     当数据库不为空时，返回一个 key 。
//     当数据库为空时，返回 nil 。
func (r *Redis) RandomKey() (NullString, error) {
	return r.run("RANDOMKEY").string()
}

// Rename key newkey
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(1)
//
//   将 key 改名为 newkey 。
//
//   当 key 和 newkey 相同，或者 key 不存在时，返回一个错误。
//
//   当 newkey 已经存在时， RENAME 命令将覆盖旧值。
//
//   返回值：
//     改名成功时提示 OK ，失败时候返回一个错误。
func (r *Redis) Rename(key, newkey string) error {
	return r.run("RENAME", key, newkey).err
}

// RenameNX key newkey
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(1)
//
//   当且仅当 newkey 不存在时，将 key 改名为 newkey 。
//
//   当 key 不存在时，返回一个错误。
//
//   返回值：
//     修改成功时，返回 1 。
//     如果 newkey 已经存在，返回 0 。
func (r *Redis) RenameNX(key, newkey string) (bool, error) {
	return r.run("RENAMENX", key, newkey).fixBool()
}

// Restore key ttl serialized-value [REPLACE]
//
//   可用版本： >= 2.6.0
//   时间复杂度：
//     查找给定键的复杂度为 O(1) ，对键进行反序列化的复杂度为 O(N*M) ，其中 N 是构成 key 的 Redis 对象的数量，而 M 则是这些对象的平均大小。
//     有序集合(sorted set)的反序列化复杂度为 O(N*M*log(N)) ，因为有序集合每次插入的复杂度为 O(log(N)) 。
//     如果反序列化的对象是比较小的字符串，那么复杂度为 O(1) 。
//
//   反序列化给定的序列化值，并将它和给定的 key 关联。
//
//   参数 ttl 以毫秒为单位为 key 设置生存时间；如果 ttl 为 0 ，那么不设置生存时间。
//
//   RESTORE 在执行反序列化之前会先对序列化值的 RDB 版本和数据校验和进行检查，如果 RDB 版本不相同或者数据不完整的话，那么 RESTORE 会拒绝进行反序列化，并返回一个错误。
//
//   如果键 key 已经存在， 并且给定了 REPLACE 选项， 那么使用反序列化得出的值来代替键 key 原有的值； 相反地， 如果键 key 已经存在， 但是没有给定 REPLACE 选项， 那么命令返回一个错误。
//
//   返回值：
//     如果反序列化成功那么返回 OK ，否则返回一个错误。
func (r *Redis) Restore(key string, ttl time.Duration, serializedValue string, Replace bool) error {
	args := []string{"RESTORE", key, durationToMillisecond(ttl), serializedValue}
	if Replace {
		args = append(args, "REPLACE")
	}
	return r.run(args...).err
}

// Sort key [BY pattern] [LIMIT offset count] [GET pattern [GET pattern ...]] [ASC | DESC] [ALPHA] [STORE destination]
//
//   返回或保存给定列表、集合、有序集合 key 中经过排序的元素。
//
//   排序默认以数字作为对象，值被解释为双精度浮点数，然后进行比较。
func (r *Redis) Sort() error {
	panic("not impl")
}

// TTL key
//
//   可用版本：>= 1.0.0
//   时间复杂度：O(1)
//
//   以秒为单位，返回给定 key 的剩余生存时间(TTL, time to live)。
//
//   返回值：
//     当 key 不存在时，返回 -2 。
//     当 key 存在但没有设置剩余生存时间时，返回 -1 。
//     否则，以秒为单位，返回 key 的剩余生存时间。
//
//   注意：这里没有实现TTL命令
func (r *Redis) TTL(key string) (*time.Duration, error) {
	c, err := r.run("PTTL", key).int()
	if err != nil {
		return nil, err
	}
	if c == -2 {
		return nil, ErrKeyNotExist
	} else if c == -1 {
		return nil, nil
	}

	t := time.Duration(c * int(time.Millisecond))
	return &t, nil
}

// Type key
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(1)
//
//   返回 key 所储存的值的类型。
//
//   返回值：
//     none (key不存在)
//     string (字符串)
//     list (列表)
//     set (集合)
//     zset (有序集)
//     hash (哈希表)
func (r *Redis) Type(key string) (KeyType, error) {
	p := r.run("TYPE", key)
	if p.err != nil {
		return "", p.err
	}
	switch strings.ToLower(p.str) {
	case "none":
		return KeyTypeNone, nil
	case "string":
		return KeyTypeString, nil
	case "list":
		return KeyTypeList, nil
	case "set":
		return KeyTypeSet, nil
	case "zset":
		return KeyTypeZSet, nil
	case "hash":
		return KeyTypeHash, nil
	default:
		return "", fmt.Errorf("invalid keytype %s", p.str)
	}
}

// Scan cursor [MATCH pattern] [COUNT count]
//
//   SCAN 命令及其相关的 SSCAN 命令、 HSCAN 命令和 ZSCAN 命令都用于增量地迭代（incrementally iterate）一集元素（a collection of elements）：
//
//     SCAN 命令用于迭代当前数据库中的数据库键。
//     SSCAN 命令用于迭代集合键中的元素。
//     HSCAN 命令用于迭代哈希键中的键值对。
//     ZSCAN 命令用于迭代有序集合中的元素（包括元素成员和元素分值）。
//
//   以上列出的四个命令都支持增量式迭代， 它们每次执行都只会返回少量元素， 所以这些命令可以用于生产环境， 而不会出现像 KEYS 命令、 SMEMBERS 命令带来的问题 —— 当 KEYS 命令被用于处理一个大的数据库时， 又或者 SMEMBERS 命令被用于处理一个大的集合键时， 它们可能会阻塞服务器达数秒之久。
//
//   不过， 增量式迭代命令也不是没有缺点的： 举个例子， 使用 SMEMBERS 命令可以返回集合键当前包含的所有元素， 但是对于 SCAN 这类增量式迭代命令来说， 因为在对键进行增量式迭代的过程中， 键可能会被修改， 所以增量式迭代命令只能对被返回的元素提供有限的保证 （offer limited guarantees about the returned elements）。
//
//   因为 SCAN 、 SSCAN 、 HSCAN 和 ZSCAN 四个命令的工作方式都非常相似， 所以这个文档会一并介绍这四个命令， 但是要记住：
//
//     SSCAN 命令、 HSCAN 命令和 ZSCAN 命令的第一个参数总是一个数据库键。
//     而 SCAN 命令则不需要在第一个参数提供任何数据库键 —— 因为它迭代的是当前数据库中的所有数据库键。
func (r *Redis) Scan(options ...ScanOption) *Scan {
	if len(options) > 1 {
		return &Scan{err: fmt.Errorf("must have 0 or 1 option")}
	}
	var args []string
	if len(options) > 0 {
		if options[0].Match != "" {
			args = append(args, "MATCH", options[0].Match)
		}
		if options[0].Count != 0 {
			args = append(args, "COUNT", strconv.Itoa(options[0].Count))
		}
	}
	return &Scan{redis: r, args: args, cursor: -1}
}
