package redis

import (
	"bufio"
	"errors"
	"strconv"
)

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
		return bytesToReply(resp), nil
	case '-':
		message, err := readUntilCRCL(r.reader)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(message)) // TODO 错误类型
	case ':':
		length, err := readIntBeforeCRCL(r.reader)
		if err != nil {
			return nil, err
		}
		return intToReply(length), nil
	case '$':
		length, err := readIntBeforeCRCL(r.reader)
		if err != nil {
			return nil, err
		}

		if length == -1 {
			return nullReply(), nil
		}

		bs, err := readBytes(r.reader, int(length))
		if err != nil {
			return nil, err
		}

		readUntilCRCL(r.reader)

		return bytesToReply(bs), nil
	case '*':
		// multi bulk reply
		count, err := readIntBeforeCRCL(r.reader)
		if err != nil {
			return nil, err
		}

		var replys []*Reply
		for i := 0; i < int(count); i++ {
			reply, err := r.read()
			if err != nil {
				return nil, err
			}
			replys = append(replys, reply)
		}

		return &Reply{replys: replys}, nil
	}

	return nil, ErrUnSupportRespType
}

func readUntilCRCL(reader *bufio.Reader) ([]byte, error) {
	bs, err := reader.ReadBytes(LF)
	if err != nil {
		return bs, err
	}

	l := len(bs)
	if l >= 2 && bs[l-2] == CR {
		return bs[:l-2], nil
	}

	return bs, nil
}

func readIntBeforeCRCL(reader *bufio.Reader) (int64, error) {
	length, err := readUntilCRCL(reader)
	if err != nil {
		return 0, err
	}
	c, err := strconv.ParseInt(string(length), 10, 64)
	if err != nil {
		return 0, err
	}
	return c, nil
}

func readBytes(reader *bufio.Reader, length int) ([]byte, error) {
	bs := make([]byte, length)
	if _, err := reader.Read(bs); err != nil {
		return nil, err
	}
	return bs, nil
}
