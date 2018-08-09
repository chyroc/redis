package redis

import (
	"bufio"
	"net"
	"sync"
	"time"
)

// Redis ...
type Redis struct {
	*sync.Mutex
	conn   net.Conn
	reader *bufio.Reader
}

// Dial conn redis
func Dial(addr string) (*Redis, error) {
	Log.Printf("try to conn %s\n", addr)
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return nil, err
	}

	r := new(Redis)
	r.conn = conn
	r.reader = bufio.NewReader(conn)
	r.Mutex = new(sync.Mutex)

	return r, nil
}
