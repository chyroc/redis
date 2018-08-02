package redis

import (
	"bytes"
	"strconv"
	"fmt"
)

const (
	LF byte = 10 // \n
	CR byte = 13 // \r
)

var CRLF = []byte{CR, LF}

func (r Redis) write(p []byte) error {
	_, err := r.conn.Write(p)
	return err
}

func (r Redis) readUntilCRCL() ([]byte, error) {
	bs, err := r.reader.ReadBytes(LF)
	if err != nil {
		fmt.Printf("[%s]\n", bs)
		return bs, err
	}

	l := len(bs)
	if l >= 2 && bs[l-2] == CR {
		fmt.Printf("[%s]\n", bs[:l-2])
		return bs[:l-2], nil
	}

	fmt.Printf("[%s]\n", bs)

	return bs, nil
}

/*

以下是这个协议的一般形式：
*<参数数量> CR LF
$<参数 1 的字节数量> CR LF
<参数 1 的数据> CR LF
...
$<参数 N 的字节数量> CR LF
<参数 N 的数据> CR LF
*/

func (r Redis) cmd(args ...string) error {
	buf := new(bytes.Buffer)

	buf.WriteString("*")
	buf.WriteString(strconv.Itoa(len(args)))
	buf.Write(CRLF)

	for _, arg := range args {
		p := []byte(arg)

		buf.WriteString("$")
		buf.WriteString(strconv.Itoa(len(p)))
		buf.Write(CRLF)

		buf.Write(p)
		buf.Write(CRLF)
	}
	fmt.Printf("send [%#v]\n", buf.String())

	return r.write(buf.Bytes())
}
