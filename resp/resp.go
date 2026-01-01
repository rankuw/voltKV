package resp

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
	Type  byte
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case STRING:
		return r.readString()
	case INTEGER:
		return r.readInteger()
	case BULK:
		return r.readBulk()
	case ARRAY:
		return r.readArray()
	case ERROR:
		return r.readError()
	default:
		return Value{}, nil
	}

}

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

func (r *Resp) readString() (Value, error) {
	line, _, err := r.readLine()
	if err != nil {
		return Value{}, err
	}

	return Value{Type: STRING, Str: string(line)}, nil
}

func (r *Resp) readError() (Value, error) {
	line, _, err := r.readLine()

	if err != nil {
		return Value{}, err
	}

	return Value{Type: ERROR, Str: string(line)}, nil
}

func (r *Resp) readInteger() (Value, error) {
	line, _, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	i64, err := strconv.Atoi(string(line))

	if err != nil {
		return Value{}, err
	}

	return Value{Type: INTEGER, Num: i64}, nil
}

func (r *Resp) readBulk() (Value, error) {
	line, _, err := r.readLine()

	if err != nil {
		return Value{}, err
	}

	size, err := strconv.Atoi(string(line))

	if err != nil {
		return Value{}, err
	}

	bytes := make([]byte, size)
	_, err = io.ReadFull(r.reader, bytes)

	if err != nil {
		return Value{}, err
	}

	return Value{Type: BULK, Str: string(bytes)}, nil

}

func (r *Resp) readArray() (Value, error) {
	line, _, err := r.readLine()

	if err != nil {
		return Value{}, err
	}

	ln, err := strconv.Atoi(string(line))

	if err != nil {
		return Value{}, err
	}

	value := Value{Type: ARRAY}
	for i := 0; i < ln; i++ {
		vl, err := r.Read()

		if err != nil {
			return Value{}, err
		}
		value.Array = append(value.Array, vl)
	}
	return value, nil
}

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) Write(v Value) error {
	fmt.Println("Value received to write", v)
	bytes := v.marshal()
	fmt.Println("marshaled bytes")
	fmt.Println(string(bytes))
	_, err := w.writer.Write(bytes)
	return err
}

func (v Value) marshal() []byte {
	switch v.Type {
	case STRING:
		return v.marshalString()
	case INTEGER:
		return v.marshalInteger()
	case ARRAY:
		return v.marshalArray()
	case BULK:
		return v.marshalBulk()
	case ERROR:
		return v.marshalError()
	default:
		return []byte{}
	}
}

func (v Value) marshalString() []byte {
	var bytes []byte

	bytes = append(bytes, STRING)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalInteger() []byte {
	var bytes []byte

	bytes = append(bytes, INTEGER)
	bytes = strconv.AppendInt(bytes, int64(v.Num), 10)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalBulk() []byte {
	var bytes []byte

	bytes = append(bytes, BULK)
	bytes = strconv.AppendInt(bytes, int64(len(v.Str)), 10)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalArray() []byte {
	var bytes []byte

	bytes = append(bytes, ARRAY)

	bytes = strconv.AppendInt(bytes, int64(len(v.Array)), 10)
	bytes = append(bytes, '\r', '\n')
	for _, val := range v.Array {
		bytes = append(bytes, val.marshal()...)
	}
	return bytes
}

func (v Value) marshalError() []byte {
	var bytes []byte

	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}
