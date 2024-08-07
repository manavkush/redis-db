package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// reads a full line from the reader. Will stop when we find a \r at the end of the line
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Resp) readArray() (Value, error) {
	val := Value{}
	val.typ = string("ARRAY")

	arrSize, _, err := r.readInteger()
	if err != nil {
		fmt.Printf("Error in reading array length. err: %v\n", err)
	}

	val.array = make([]Value, 0)
	for range arrSize {
		value, err := r.Read()
		if err != nil {
			fmt.Printf("Error in reading value. err: %v\n", err)
		}

		val.array = append(val.array, value)
	}

	return val, nil
}

func (r *Resp) readBulk() (Value, error) {
	val := Value{}
	val.typ = "BULK"

	bulkSize, _, err := r.readInteger()
	if err != nil {
		fmt.Printf("Error in reading bulk size. err: %v\n", err)
	}

	bulkBuf := make([]byte, bulkSize)

	r.reader.Read(bulkBuf)
	val.bulk = string(bulkBuf)

	// Read the ending \r\n
	r.readLine()

	return val, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v\n", string(_type))
		return Value{}, nil
	}
}
