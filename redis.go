package redis

import (
	"bufio"
	"net"
	"time"
)

type Redis struct {
	conn   net.Conn
	reader *bufio.Reader
}

func Dial(addr string) (*Redis, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return nil, err
	}

	r := new(Redis)
	r.conn = conn
	r.reader = bufio.NewReader(conn)

	return r, nil
}
