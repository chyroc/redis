package redis

// Discard ...
//
//   取消事务，放弃执行事务块内的所有命令。
//
//   如果正在使用 WATCH 命令监视某个(或某些) key，那么取消所有监视，等同于执行命令 UNWATCH 。
//
//   可用版本：>= 2.0.0
//   时间复杂度：O(1)。
//
//   返回值：
//     总是返回 OK 。
func (r *Redis) Discard() error {
	return r.run("DISCARD").errNotFromReply
}

// Exec ...
//
//   执行所有事务块内的命令。
//
//   假如某个(或某些) key 正处于 WATCH 命令的监视之下，且事务块中有和这个(或这些) key 相关的命令，那么 EXEC 命令只在这个(或这些) key 没有被其他命令所改动的情况下执行并生效，否则该事务被打断(abort)。
//
//   可用版本：>= 1.2.0
//   时间复杂度：事务块内所有命令的时间复杂度的总和。
//
//   返回值：
//     事务块内所有命令的返回值，按命令执行的先后顺序排列。
//     当操作被打断时，返回空值 nil 。
func (r *Redis) Exec() {
	r.run("EXEC")
	// TODO 返回值的格式？
	// list
	//   - String
	//   - Int
	//   - List
}

// Multi ...
//
//   标记一个事务块的开始。
//
//   事务块内的多条命令会按照先后顺序被放进一个队列当中，最后由 EXEC 命令原子性(atomic)地执行。
//
//   可用版本：>= 1.2.0
//   时间复杂度：O(1)。
//
//   返回值：
//     总是返回 OK 。
func (r *Redis) Multi() error {
	// TODO
	// 因为开启事务后，其他操作返回的都是QUEUED，所以这里应该返回一个struct，然后把其他的方法全部复制过来？结构体组合？
	return r.run("MULTI").errNotFromReply
}

// UnWatch ...
//
//   取消 WATCH 命令对所有 key 的监视。
//
//   如果在执行 WATCH 命令之后， EXEC 命令或 DISCARD 命令先被执行了的话，那么就不需要再执行 UNWATCH 了。
//
//   因为 EXEC 命令会执行事务，因此 WATCH 命令的效果已经产生了；而 DISCARD 命令在取消事务的同时也会取消所有对 key 的监视，因此这两个命令执行之后，就没有必要执行 UNWATCH 了。
//
//   可用版本：>= 2.2.0
//   时间复杂度：O(1)
//
//   返回值：
//     总是 OK 。
func (r *Redis) UnWatch() error {
	return r.run("UNWATCH").errNotFromReply
}

// Watch key [key ...]
//
//   监视一个(或多个) key ，如果在事务执行之前这个(或这些) key 被其他命令所改动，那么事务将被打断。
//
//   可用版本：>= 2.2.0
//   时间复杂度：O(1)。
//
//   返回值：
//     总是返回 OK 。
func (r *Redis) Watch(key string, keys ...string) error {
	return r.run(buildSlice2("WATCH", key, keys)...).errNotFromReply
}