package redis

const (
	CL byte = 10 // \n
	CR byte = 13 // \r
)

func (r Redis) readUntilCRCL() ([]byte, error) {
	bs, err := r.reader.ReadBytes(CL)
	if err != nil {
		return bs, err
	}

	l := len(bs)
	if l >= 2 && bs[l-2] == CR {
		return bs[:l-2], nil
	}

	return bs, nil
}

func (r Redis) write(p []byte) error {
	_, err := r.writer.Write(p)
	return err
	//if _, err := conn.Write([]byte("*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n")); err != nil {
	//	return nil, err
	//}
	//
	//line, err := readUntilCRCL(reader)
	//if err != nil {
	//	return nil, err
	//}
	//
	//fmt.Printf("xxx [%s]\n", line)
}
