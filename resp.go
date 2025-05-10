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

func (r *Resp) ReadLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			fmt.Println("err:", err)
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		// carriage return
		if len(line) >= 2 && line[len(line)-2] == 13 {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) ReadInteger() (x int, n int, err error) {
	line, n, err := r.ReadLine()
	if err != nil {
		fmt.Println("Error reading the line :", err)
		return 0, 0, nil
	}
	i64, errInt := strconv.ParseInt(string(line), 10, 64)
	if errInt != nil {
		fmt.Println("Error converting to int :", err)
		return 0, n, nil
	}
	return int(i64), n, nil
}

func (r *Resp) ReadArray() (Value, error) {
	v := Value{}
	v.typ = "array"
	size, _, err := r.ReadInteger()
	if err != nil {
		return Value{}, err
	}
	v.array = make([]Value, size)
	for i := 0; i < size; i++ {
		val, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		v.array[i] = val
	}
	return v, nil
}

func (r *Resp) ReadBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"
	// get size of the letter
	len, _, err := r.ReadInteger()
	fmt.Println("size", len)
	if err != nil {
		return Value{}, err
	}
	bulk := make([]byte, len)
	r.reader.Read(bulk)

	v.bulk = string(bulk)
	r.ReadLine()
	return v, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		fmt.Println("Error reading byte", err)
		return Value{}, err
	}
	switch _type {
	case ARRAY:
		return r.ReadArray()
	case BULK:
		return r.ReadBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

// // test string
// // *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
// func test() {
// 	input := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
// 	r := NewResp(strings.NewReader(input))
// 	v, _ := r.Read()
// 	fmt.Println(v)
// 	// b, err := reader.ReadByte()
// 	// if err != nil {
// 	// 	fmt.Println("Error occured :", err)
// 	// 	os.Exit(1)
// 	// }
// 	// if b != '$' {
// 	// 	fmt.Println("Invalid first character")
// 	// 	os.Exit(1)
// 	// }
// 	// strSize, _ := reader.ReadByte()
// 	// size, _ := strconv.ParseInt(string(strSize), 10, 64)
// 	// // consume \r\n
// 	// reader.ReadByte()
// 	// reader.ReadByte()
// 	// name := make([]byte, size)
// 	// reader.Read(name)
// 	// fmt.Println(string(name))
// }
