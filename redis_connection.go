package redis

import "strconv"

// Auth password
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(1)
//
//   通过设置配置文件中 requirepass 项的值(使用命令 CONFIG SET requirepass password )，可以使用密码来保护 Redis 服务器。
//
//   如果开启了密码保护的话，在每次连接 Redis 服务器之后，就要使用 AUTH 命令解锁，解锁之后才能使用其他 Redis 命令。
//
//   如果 AUTH 命令给定的密码 password 和配置文件中的密码相符的话，服务器会返回 OK 并开始接受命令输入。
//
//   另一方面，假如密码不匹配的话，服务器将返回一个错误，并要求客户端需重新输入密码。
//
//   因为 Redis 高性能的特点，在很短时间内尝试猜测非常多个密码是有可能的，因此请确保使用的密码足够复杂和足够长，以免遭受密码猜测攻击。
//
//   返回值：
//     密码匹配时返回 OK ，否则返回一个错误。
func (r *Redis) Auth(password string) error {
	return r.run("AUTH", password).errNotFromReply
}

// Echo message
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(1)
//
//   打印一个特定的信息 message ，测试时使用。
//
//   返回值：
//     message 自身。
func (r *Redis) Echo(message string) (NullString, error) {
	return r.run("ECHO", message).string()
}

// Ping ...
//
//   可用版本：>= 1.0.0
//   时间复杂度：O(1)
//
//   使用客户端向 Redis 服务器发送一个 PING ，如果服务器运作正常的话，会返回一个 PONG 。
//
//   通常用于测试与服务器的连接是否仍然生效，或者用于测量延迟值。
//
//   返回值：
//     如果连接正常就返回一个 PONG ，否则返回一个连接错误。
func (r *Redis) Ping() (NullString, error) {
	return r.run("PING").string()
}

// Quit ...
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(1)
//
//   请求服务器关闭与当前客户端的连接。
//
//   一旦所有等待中的回复(如果有的话)顺利写入到客户端，连接就会被关闭。
//
//   返回值：
//     总是返回 OK (但是不会被打印显示，因为当时 Redis-cli 已经退出)。
func (r *Redis) Quit() error {
	return r.run("QUIT").errNotFromReply
}

// Select index
//
//   可用版本： >= 1.0.0
//   时间复杂度： O(1)
//
//   切换到指定的数据库，数据库索引号 index 用数字值指定，以 0 作为起始索引值。
//
//   默认使用 0 号数据库。
//
//   返回值：
//     OK
func (r *Redis) Select(index int) error {
	return r.run("SELECT", strconv.Itoa(index)).errNotFromReply
}
