package redis

import (
	"bufio"
	"fmt"
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
	fmt.Printf("try to conn %s\n", addr)
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
