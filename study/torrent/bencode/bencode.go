package bencode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

var (
	ErrNum = errors.New("except num")
	ErrCol = errors.New("except colon")
	ErrEpi = errors.New("except i")
	ErrEpe = errors.New("except e")
	ErrTyp = errors.New("wrong type")
	ErrIvd = errors.New("invalid bencode")
)

type BType uint8

const (
	BSTR  BType = 0x01
	BINT  BType = 0x02
	BLIST BType = 0x03
	BDICT BType = 0x04
)

type BValue interface{}

type BObject struct {
	_type BType
	_val  BValue
}

func (o *BObject) Str() (string, error) {
	if o._type == BSTR {
		return o._val.(string), nil
	}
	return "", ErrTyp
}
func (o *BObject) Int() (int, error) {
	if o._type == BINT {
		return o._val.(int), nil
	}
	return 0, ErrTyp
}
func (o *BObject) List() ([]*BObject, error) {
	if o._type == BLIST {
		return o._val.([]*BObject), nil
	}
	return nil, ErrTyp
}
func (o *BObject) Dict() (map[string]*BObject, error) {
	if o._type == BDICT {
		return o._val.(map[string]*BObject), nil
	}
	return nil, ErrTyp
}
func (o *BObject) BEncode(w io.Writer) int {
	bw, ok := w.(*bufio.Writer)
	if !ok {
		bw = bufio.NewWriter(w)
	}
	wLen := 0
	switch o._type {
	case BSTR:
		str, _ := o.Str()
		wLen += EncodeString(bw, str)
	case BINT:
		val, _ := o.Int()
		wLen += EncodeInt(bw, val)
	case BLIST:
		bw.WriteByte('l')
		list, _ := o.List()
		for _, elem := range list {
			wLen += elem.BEncode(bw)
		}
		bw.WriteByte('e')
		wLen += 2
	case BDICT:
		bw.WriteByte('d')
		dict, _ := o.Dict()
		for k, v := range dict {
			wLen += EncodeString(bw, k)
			wLen += v.BEncode(bw)
		}
		bw.WriteByte('e')
		wLen += 2
	}
	bw.Flush()
	return wLen
}
func EncodeString(w io.Writer, val string) int {
	n := len(val)
	bw, ok := w.(*bufio.Writer)
	if !ok {
		bw = bufio.NewWriter(w)
	}
	str := fmt.Sprintf("%d:%s", n, val)
	_, err := bw.WriteString(str)
	if err != nil {
		return 0
	}
	err = bw.Flush()
	if err != nil {
		return 0
	}
	return len(str)
}
func DecodeString(r io.Reader) (val string, err error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	num, len := readDecimal(br)
	if len == 0 {
		return val, ErrNum
	}
	b, err := br.ReadByte()
	if b != ':' {
		return val, ErrCol
	}
	buf := make([]byte, num)
	_, err = io.ReadAtLeast(br, buf, num)
	val = string(buf)
	return
}
func EncodeInt(w io.Writer, val int) int {
	bw, ok := w.(*bufio.Writer)
	if !ok {
		bw = bufio.NewWriter(w)
	}
	str := fmt.Sprintf("i%de", val)
	_, err := bw.WriteString(str)
	if err != nil {
		return 0
	}
	n := len(str)
	err = bw.Flush()
	if err != nil {
		return 0
	}
	return n
}
func DecodeInt(r io.Reader) (val int, err error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	b, err := br.ReadByte()
	if b != 'i' {
		return val, ErrEpi
	}
	val, _ = readDecimal(br)
	b, err = br.ReadByte()
	if b != 'e' {
		return val, ErrEpe
	}
	return
}
func isNum(data byte) bool {
	return data >= '0' && data <= '9'
}

func readDecimal(r *bufio.Reader) (val int, len int) {
	sign := 1
	b, _ := r.ReadByte()
	len++
	if b == '-' {
		sign = -1
		b, _ = r.ReadByte()
		len++
	}
	for {
		if !isNum(b) {
			r.UnreadByte()
			len--
			return sign * val, len
		}
		val = val*10 + int(b-'0')
		b, _ = r.ReadByte()
		len++
	}
}
