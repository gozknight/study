package bencode

import (
	"bufio"
	"io"
)

func Parse(r io.Reader) (*BObject, error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	b, err := br.Peek(1)
	if err != nil {
		return nil, err
	}
	var ans BObject
	switch {
	case b[0] >= '0' && b[0] <= '9':
		val, err := DecodeString(br)
		if err != nil {
			return nil, err
		}
		ans._type = BSTR
		ans._val = val
	case b[0] == 'i':
		val, err := DecodeInt(br)
		if err != nil {
			return nil, err
		}
		ans._type = BINT
		ans._val = val
	case b[0] == 'l':
		br.ReadByte()
		var list []*BObject
		for {
			if p, _ := br.Peek(1); p[0] == 'e' {
				br.ReadByte()
				break
			}
			elem, err := Parse(br)
			if err != nil {
				return nil, err
			}
			list = append(list, elem)
		}
		ans._type = BLIST
		ans._val = list
	case b[0] == 'd':
		br.ReadByte()
		dict := make(map[string]*BObject)
		for {
			if p, _ := br.Peek(1); p[0] == 'e' {
				br.ReadByte()
				break
			}
			key, err := DecodeString(br)
			if err != nil {
				return nil, err
			}
			val, err := Parse(br)
			if err != nil {
				return nil, err
			}
			dict[key] = val
		}
		ans._type = BDICT
		ans._val = dict
	default:
		return nil, ErrIvd
	}
	return &ans, nil
}
