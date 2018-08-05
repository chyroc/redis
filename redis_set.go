package redis

// SAdd key member [member ...]
//
//   可用版本： >= 1.0.0
//   时间复杂度: O(N)， N 是被添加的元素的数量。
//
//   将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。
//
//   假如 key 不存在，则创建一个只包含 member 元素作成员的集合。
//
//   当 key 不是集合类型时，返回一个错误。
//
//   在Redis2.4版本以前， SADD 只接受单个 member 值。
//   返回值:
//     被添加到集合中的新元素的数量，不包括被忽略的元素。
func (r *Redis) SAdd(key, member string, members ...string) (int, error) {
	return r.run(append([]string{"SADD", key, member}, members...)...).int()
}

func (r *Redis) SCARD() {

}

func (r *Redis) SDIFF() {

}

func (r *Redis) SDIFFSTORE() {

}

func (r *Redis) SINTER() {

}

func (r *Redis) SINTERSTORE() {

}

func (r *Redis) SISMEMBER() {

}

func (r *Redis) SMEMBERS() {

}

func (r *Redis) SMOVE() {

}

func (r *Redis) SPOP() {

}

func (r *Redis) SRANDMEMBER() {

}

func (r *Redis) SREM() {

}

func (r *Redis) SUNION() {

}

func (r *Redis) SUNIONSTORE() {

}

func (r *Redis) SSCAN() {

}
