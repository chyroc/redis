package redis

type Pubsub struct {
	redis *Redis
}

// Channels [pattern]
//
//   列出当前的活跃频道。
//
//   活跃频道指的是那些至少有一个订阅者的频道， 订阅模式的客户端不计算在内。
//
//   pattern 参数是可选的：
//
//     如果不给出 pattern 参数，那么列出订阅与发布系统中的所有活跃频道。
//     如果给出 pattern 参数，那么只列出和给定模式 pattern 相匹配的那些活跃频道。
//
//   复杂度： O(N) ， N 为活跃频道的数量（对于长度较短的频道和模式来说，将进行模式匹配的复杂度视为常数）。
//
//   返回值
//     一个由活跃频道组成的列表。
func (r *Pubsub) Channels(patterns ...string) ([]string, error) {
	return r.redis.run(buildSlice2("PUBSUB", "CHANNELS", patterns)...).fixStringSlice()
}

// NumSubscribe [channel-1 ... channel-N]
//
//   返回给定频道的订阅者数量， 订阅模式的客户端不计算在内。
//
//   复杂度： O(N) ， N 为给定频道的数量。
//
//   返回值
//     一个多条批量回复（Multi-bulk reply），回复中包含给定的频道，以及频道的订阅者数量。
//     格式为：频道 channel-1 ， channel-1 的订阅者数量，频道 channel-2 ， channel-2 的订阅者数量，诸如此类。
//     回复中频道的排列顺序和执行命令时给定频道的排列顺序一致。 不给定任何频道而直接调用这个命令也是可以的， 在这种情况下， 命令只返回一个空列表。
func (r *Pubsub) NumSubscribe(channels ...string) (map[string]int, error) {
	p := r.redis.run(buildSlice2("PUBSUB", "NUMSUB", channels)...)
	if p.errNotFromReply != nil {
		return nil, p.errNotFromReply
	}

	m := make(map[string]int)
	for k := 0; k+1 < len(p.replys); k += 2 {
		m[p.replys[k].str] = int(p.replys[k+1].integer)
	}
	return m, nil
}

// NumPattern ...
//
//   返回订阅模式的数量。
//
//   注意， 这个命令返回的不是订阅模式的客户端的数量， 而是客户端订阅的所有模式的数量总和。
//
//   复杂度： O(1) 。
//
//   返回值
//     一个整数回复（Integer reply）。
func (r *Pubsub) NumPattern() (int, error) {
	return r.redis.run("PUBSUB", "NUMPAT").int()
}
