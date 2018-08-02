package redis

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
)

func (r *Redis) readToReply() *Reply {
	reply, err := r.read()
	if err != nil {
		return &Reply{err: err}
	}
	return reply
}

func (r *Redis) read() (*Reply, error) {
	respType, err := r.reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch respType {
	case '+':
		resp, err := readUntilCRCL(r.reader)
		if err != nil {
			return nil, err
		}
		return &Reply{str: nullString{String: string(resp)}}, nil
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
			return &Reply{str: nullString{Valid: false}}, nil
		}

		bs := make([]byte, c)
		if _, err := r.reader.Read(bs); err != nil {
			return nil, err
		}

		readUntilCRCL(r.reader)

		return &Reply{str: nullString{String: string(bs)}}, nil
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
