package redis

import (
	"bytes"
	"fmt"
	"strconv"
)

func (r Redis) write(p []byte) error {
	_, err := r.conn.Write(p)
	return err
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
