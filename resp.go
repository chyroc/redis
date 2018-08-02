package redis

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
)

func (r *Redis) read() (interface{}, error) {
	respType, err := r.reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch respType {
	case '+':
		_, err := readUntilCRCL(r.reader)
		return nil, err
	case '-':
		message, err := readUntilCRCL(r.reader)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(message)) // TODO 错误类型
	case ':':

	case '$':
		length, err := readUntilCRCL(r.reader)
		if err != nil {
			return nil, err
		}
		c, err := strconv.Atoi(string(length))
		if err != nil {
			return nil, err
		}

		if c == -1 {
			return "", nil // TODO return Null Bulk String
		}

		bs := make([]byte, c)
		if _, err := r.reader.Read(bs); err != nil {
			return nil, err
		}

		readUntilCRCL(r.reader)

		return bs, nil
	case '*':
	}

	return nil, UnSupportRespType
}

func readUntilCRCL(reader *bufio.Reader) ([]byte, error) {
	bs, err := reader.ReadBytes(LF)
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
