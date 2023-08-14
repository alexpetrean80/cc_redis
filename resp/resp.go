package resp

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	SimpleString = byte('+')
	Error        = byte('-')
	Integer      = byte(':')
	BulkString   = byte('$')
	Array        = byte('*')
)

type RESPMessage struct {
	dataType byte
	data     []byte
}

func New(d []byte) *RESPMessage {
	return &RESPMessage{
		dataType: d[0],
		data:     d[1:],
	}
}

func (r *RESPMessage) Parse() (v interface{}, err error) {
	switch r.dataType {
	case SimpleString:
		return r.parseSimpleString(), nil
	case Error:
		return r.parseError(), nil
	case Integer:
		return r.parseInteger()
	case BulkString:
		return r.parseBulkString()
	case Array:
		return r.parseArray()
	default:
		return nil, errors.New("invalid data type")
	}
}

func (r RESPMessage) parseSimpleString() string {
	return trimStr(r.data)
}

func (r RESPMessage) parseError() error {
	return errors.New(r.parseSimpleString())
}

func (r RESPMessage) parseInteger() (interface{}, error) {
	v, err := strconv.ParseInt(trimStr(r.data), 10, 64)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (r RESPMessage) parseBulkString() (interface{}, error) {
	str := trimStr(r.data)
	split := strings.Split(str, "\r\n")
	byteCount, err := strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return nil, err
	}
	if byteCount == -1 {
		return nil, nil
	}
	if byteCount == 0 {
		return "", nil
	}
	if byteCount > 512 {
		return nil, errors.New("Bulk string too large")
	}
	return split[1], nil
}

func (r RESPMessage) parseArray() ([]interface{}, error) {
	str := trimStr(r.data)

	splits := strings.Split(str, "\r\n")
	byteCount, err := strconv.ParseInt(splits[0], 10, 64)
	if err != nil {
		return nil, err
	}

	str = ""
	for _, s := range splits[1:] {
		str += s
		str += "\r\n"
	}

	matches := regexp.MustCompile(`[\+\-:\$\*][[:ascii:]]*\r\n`).FindAllString(str, -1)

	if err != nil {
		return nil, err
	}
	if byteCount == -1 {
		return nil, nil
	}
	if byteCount == 0 {
		return []interface{}{}, nil
	}
	result := []interface{}{}
	for _, v := range matches {
		partial, err := New([]byte(v)).Parse()
		if err != nil {
			return nil, err
		}
		result = append(result, partial)
	}

	return result, nil
}

func trimStr(s []byte) string {
	return strings.Trim(string(s), "\r\n")
}
