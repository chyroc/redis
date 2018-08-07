package redis

import (
	"fmt"
)

// PSubscribe ...
type PSubscribe struct {
	Err     error
	Pattern string
	Channel string
	Message string
}

// Subscribe ...
type Subscribe struct {
	Err     error
	Channel string
	Message string
}

// PSubscribe pattern [pattern ...]
//
//   订阅一个或多个符合给定模式的频道。
//
//   每个模式以 * 作为匹配符，比如 it* 匹配所有以 it 开头的频道( it.news 、 it.blog 、 it.tweets 等等)， news.* 匹配所有以 news. 开头的频道( news.it 、 news.global.today 等等)，诸如此类。
//
//   可用版本：>= 2.0.0
//   时间复杂度：O(N)， N 是订阅的模式的数量。
//
//   返回值：
//     接收到的信息(请参见下面的代码说明)。
func (r *Redis) PSubscribe(pattern string, patterns ...string) (chan PSubscribe, error) {
	r.Lock()
	defer r.Unlock()

	p := r.runWithLock(buildSlice2("PSUBSCRIBE", pattern, patterns)...)
	if p.errNotFromReply != nil {
		return nil, p.errNotFromReply
	}

	var s = make(chan PSubscribe, 1024)
	go func() {
		for {
			reply, err := r.read()
			if err != nil {
				s <- PSubscribe{Err: err}
			}
			fmt.Printf("psub: %s\n", reply)

			if len(reply.replys) < 1 {
				s <- PSubscribe{Err: fmt.Errorf("expect at least subscribe response")}
			}

			switch reply.replys[0].str {
			case "psubscribe":
			case "pmessage":
				if len(reply.replys) < 4 {
					s <- PSubscribe{Err: fmt.Errorf("expect 4 message for pmessage response")}
				}
				s <- PSubscribe{Pattern: reply.replys[1].str, Channel: reply.replys[2].str, Message: reply.replys[3].str}
			default:
				s <- PSubscribe{Err: fmt.Errorf("invalid reply: %s", reply)}
			}
		}
	}()

	return s, nil
}

// Publish channel message
//
//   将信息 message 发送到指定的频道 channel 。
//
//   可用版本：>= 2.0.0
//   时间复杂度：O(N+M)，其中 N 是频道 channel 的订阅者数量，而 M 则是使用模式订阅(subscribed patterns)的客户端的数量。
//
//   返回值：
//     接收到信息 message 的订阅者数量。
func (r *Redis) Publish(channel, message string) (int, error) {
	return r.run("PUBLISH", channel, message).int()
}

// Pubsub <subcommand> [argument [argument ...]]
func (r *Redis) Pubsub() *PubSub {
	return &PubSub{r}
}

// PUnSubscribe [pattern [pattern ...]]
//
//   指示客户端退订所有给定模式。
//
//   如果没有模式被指定，也即是，一个无参数的 PUNSUBSCRIBE 调用被执行，那么客户端使用 PSUBSCRIBE 命令订阅的所有模式都会被退订。在这种情况下，命令会返回一个信息，告知客户端所有被退订的模式。
//
//   可用版本：>= 2.0.0
//   时间复杂度：O(N+M) ，其中 N 是客户端已订阅的模式的数量， M 则是系统中所有客户端订阅的模式的数量。
//
//   返回值：
//     这个命令在不同的客户端中有不同的表现。
func (r *Redis) PUnSubscribe(patterns ...string) {
	// TODO: r.run(buildSlice1("PUNSUBSCRIBE",patterns)...)
}

// Subscribe channel [channel ...]
//
//   订阅给定的一个或多个频道的信息。
//
//   可用版本：>= 2.0.0
//   时间复杂度：O(N)，其中 N 是订阅的频道的数量。
//
//   返回值：
//     接收到的信息(请参见下面的代码说明)。
func (r *Redis) Subscribe(channel string, channels ...string) (chan Subscribe, error) {
	r.Lock()
	defer r.Unlock()

	p := r.runWithLock(buildSlice2("SUBSCRIBE", channel, channels)...)
	if p.errNotFromReply != nil {
		return nil, p.errNotFromReply
	}

	var s = make(chan Subscribe, 1024)
	go func() {
		for {
			reply, err := r.read()
			if err != nil {
				s <- Subscribe{Err: err}
			}
			fmt.Printf("sub: %s\n", reply)

			if len(reply.replys) < 1 {
				s <- Subscribe{Err: fmt.Errorf("expect at least subscribe response")}
			}

			switch reply.replys[0].str {
			case "subscribe":
			case "message":
				if len(reply.replys) < 3 {
					s <- Subscribe{Err: fmt.Errorf("expect 4 message for pmessage response")}
				}
				s <- Subscribe{Channel: reply.replys[1].str, Message: reply.replys[2].str}
			case "unsubscribe":
				close(s)
			default:
				s <- Subscribe{Err: fmt.Errorf("invalid reply: %s", reply)}
			}
		}
	}()

	return s, nil
}

// UnSubscribe [channel [channel ...]]
//
//   指示客户端退订给定的频道。
//
//   如果没有频道被指定，也即是，一个无参数的 UNSUBSCRIBE 调用被执行，那么客户端使用 SUBSCRIBE 命令订阅的所有频道都会被退订。在这种情况下，命令会返回一个信息，告知客户端所有被退订的频道。
//
//   可用版本：>= 2.0.0
//   时间复杂度：O(N) ， N 是客户端已订阅的频道的数量。
//
//   返回值：
//     这个命令在不同的客户端中有不同的表现。
func (r *Redis) UnSubscribe(channels ...string) (int, error) {
	// TODO close go channel
	return r.run(buildSlice1("UNSUBSCRIBE", channels)...).int()
}
