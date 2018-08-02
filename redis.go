package redis

import (
	"net"
	"bufio"
	"time"
)

type Redis struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func Dial(addr string) (*Redis, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return nil, err
	}
	r := new(Redis)

	r.conn = conn
	r.reader = bufio.NewReader(conn)
	r.writer = bufio.NewWriter(conn)

	return r, nil
}
