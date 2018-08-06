package redis

import (
	"fmt"
	"strconv"
)

// HDel key field [field ...]
//
//   可用版本：>= 2.0.0
//   时间复杂度: O(N)， N 为要删除的域的数量。
//
//   删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。
//
//   在Redis2.4以下的版本里， HDEL 每次只能删除单个域，如果你需要在一个原子时间内删除多个域，请将命令包含在 MULTI / EXEC 块内。
//
//   返回值:
//     被成功移除的域的数量，不包括被忽略的域。
func (r *Redis) HDel(field string, fields ...string) (int, error) {
	return r.run(append([]string{"HDEL", field}, fields...)...).int()
}

// HExists key field
//
//   查看哈希表 key 中，给定域 field 是否存在。
//
//   可用版本： >= 2.0.0
//   时间复杂度： O(1)
//   返回值：
//     如果哈希表含有给定域，返回 1 。
//     如果哈希表不含有给定域，或 key 不存在，返回 0 。
func (r *Redis) HExists(key, field string) (bool, error) {
	return r.run("HEXISTS", key, field).fixBool()
}

// HGet key field
//
//   可用版本： >= 2.0.0
//   时间复杂度：O(1)
//
//   返回哈希表 key 中给定域 field 的值。
//
//   返回值：
//     给定域的值。
//     当给定域不存在或是给定 key 不存在时，返回 nil 。
func (r *Redis) HGet(key, field string) (NullString, error) {
	return r.run("HGET", key, field).string()
}

// HGetALL key
//
//   可用版本： >= 2.0.0
//   时间复杂度： O(N)， N 为哈希表的大小。
//
//   返回哈希表 key 中，所有的域和值。
//
//   在返回值里，紧跟每个域名(field name)之后是域的值(value)，所以返回值的长度是哈希表大小的两倍。
//
//   返回值：
//     以列表形式返回哈希表的域和域的值。
//     若 key 不存在，返回空列表。
func (r *Redis) HGetALL(key string) (map[string]string, error) {
	return r.run("HGETALL", key).fixMap()
}

// HIncrBy key field increment
//
//   可用版本： >= 2.0.0
//   时间复杂度： O(1)
//
//   为哈希表 key 中的域 field 的值加上增量 increment 。
//
//   增量也可以为负数，相当于对给定域进行减法操作。
//
//   如果 key 不存在，一个新的哈希表被创建并执行 HINCRBY 命令。
//
//   如果域 field 不存在，那么在执行命令前，域的值被初始化为 0 。
//
//   对一个储存字符串值的域 field 执行 HINCRBY 命令将造成一个错误。
//
//   本操作的值被限制在 64 位(bit)有符号数字表示之内。
//
//   返回值：
//     执行 HINCRBY 命令之后，哈希表 key 中域 field 的值。
func (r *Redis) HIncrBy(key, field string, increment int) (int, error) {
	return r.run("HINCRBY", key, field, strconv.Itoa(increment)).int()
}

// HIncrByFloat key field increment
//
//   可用版本： >= 2.6.0
//   时间复杂度： O(1)
//
//   为哈希表 key 中的域 field 加上浮点数增量 increment 。
//
//   如果哈希表中没有域 field ，那么 HINCRBYFLOAT 会先将域 field 的值设为 0 ，然后再执行加法操作。
//
//   如果键 key 不存在，那么 HINCRBYFLOAT 会先创建一个哈希表，再创建域 field ，最后再执行加法操作。
//
//   当以下任意一个条件发生时，返回一个错误：
//
//     域 field 的值不是字符串类型(因为 redis 中的数字和浮点数都以字符串的形式保存，所以它们都属于字符串类型）
//     域 field 当前的值或给定的增量 increment 不能解释(parse)为双精度浮点数(double precision floating point number)
//
//   HINCRBYFLOAT 命令的详细功能和 INCRBYFLOAT 命令类似，请查看 INCRBYFLOAT 命令获取更多相关信息。
//
//   返回值：
//     执行加法操作之后 field 域的值。
func (r *Redis) HIncrByFloat(key, field string, increment float64) (float64, error) {
	return r.run("HINCRBYFLOAT", key, field, float64ToString(increment)).fixFloat()
}

// HKeys key
//
//   可用版本： >= 2.0.0
//   时间复杂度： O(N)， N 为哈希表的大小。
//
//   返回哈希表 key 中的所有域。
//
//   返回值：
//     一个包含哈希表中所有域的表。
//     当 key 不存在时，返回一个空表。
func (r *Redis) HKeys(key string) ([]string, error) {
	return r.run("HKEYS", key).fixStringSlice()
}

// HLen key
//
//   时间复杂度： O(1)
//
//   返回哈希表 key 中域的数量。
//
//   返回值：
//     哈希表中域的数量。
//     当 key 不存在时，返回 0 。
func (r *Redis) HLen(key string) (int, error) {
	return r.run("HLEN", key).int()
}

// HMGet key field [field ...]
//
//   可用版本： >= 2.0.0
//   时间复杂度： O(N)， N 为给定域的数量。
//
//   返回哈希表 key 中，一个或多个给定域的值。
//
//   如果给定的域不存在于哈希表，那么返回一个 nil 值。
//
//   因为不存在的 key 被当作一个空哈希表来处理，所以对一个不存在的 key 进行 HMGET 操作将返回一个只带有 nil 值的表。
//
//   返回值：
//     一个包含多个给定域的关联值的表，表值的排列顺序和给定域参数的请求顺序一样。
func (r *Redis) HMGet(key, field string, fields ...string) ([]NullString, error) {
	return r.run(append([]string{"HMGET", key, field}, fields...)...).fixNullStringSlice()
}

// HMSet key field value [field value ...]
//
//   可用版本： >= 2.0.0
//   时间复杂度： O(N)， N 为 field-value 对的数量。
//
//   同时将多个 field-value (域-值)对设置到哈希表 key 中。
//
//   此命令会覆盖哈希表中已存在的域。
//
//   如果 key 不存在，一个空哈希表被创建并执行 HMSET 操作。
//
//   返回值：
//     如果命令执行成功，返回 OK 。
//     当 key 不是哈希表(hash)类型时，返回一个错误。
func (r *Redis) HMSet(key, field, value string, kvs ...string) error {
	if len(kvs)%2 != 0 {
		return fmt.Errorf("key value pair, but got %d arguments", len(kvs)+3)
	}

	return r.run(append([]string{"HMSET", key, field, value}, kvs...)...).errNotFromReply
}

// HSet key field value
//
//   可用版本： >= 2.0.0
//   时间复杂度： O(1)
//
//   将哈希表 key 中的域 field 的值设为 value 。
//
//   如果 key 不存在，一个新的哈希表被创建并进行 HSET 操作。
//
//   如果域 field 已经存在于哈希表中，旧值将被覆盖。
//
//   返回值：
//     如果 field 是哈希表中的一个新建域，并且值设置成功，返回 1 。
//     如果哈希表中域 field 已经存在且旧值已被新值覆盖，返回 0 。
func (r *Redis) HSet(key, field, value string) (bool, error) {
	return r.run("HSET", key, field, value).fixBool()
}

// HSetNX key field value
//
//   可用版本： >= 2.0.0
//   时间复杂度： O(1)
//
//   将哈希表 key 中的域 field 的值设置为 value ，当且仅当域 field 不存在。
//
//   若域 field 已经存在，该操作无效。
//
//   如果 key 不存在，一个新哈希表被创建并执行 HSETNX 命令。
//
//   返回值：
//     设置成功，返回 1 。
//     如果给定域已经存在且没有操作被执行，返回 0 。
func (r *Redis) HSetNX(key, field, value string) (bool, error) {
	return r.run("HSETNX", key, field, value).fixBool()
}

// HVals key
//
//   返回哈希表 key 中所有域的值。
//
//   可用版本： >= 2.0.0
//   时间复杂度： O(N)， N 为哈希表的大小。
//   返回值：
//     一个包含哈希表中所有值的表。
//     当 key 不存在时，返回一个空表。
func (r *Redis) HVals(key string) ([]string, error) {
	return r.run("HVALS", key).fixStringSlice()
}

// HScan ...
//
// 参见Scan的文档
func (r *Redis) HScan(key string, options ...ScanOption) *HScan {
	if len(options) > 1 {
		return &HScan{err: fmt.Errorf("must have 0 or 1 option")}
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
	return &HScan{redis: r, args: args, cursor: -1, result: make(map[string]string), key: key}
}

// HStrLen key field
//
//   可用版本： >= 3.2.0
//   时间复杂度： O(1)
//
//   返回哈希表 key 中， 与给定域 field 相关联的值的字符串长度（string length）。
//
//   如果给定的键或者域不存在， 那么命令返回 0 。
//
//   返回值：
//     一个整数。
func (r *Redis) HStrLen(key, field string) (int, error) {
	return r.run("HSTRLEN", key, field).int()
}
