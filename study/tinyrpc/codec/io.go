package codec

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
)

func sendFrame(w io.Writer, data []byte) (err error) {
	var size [binary.MaxVarintLen64]byte
	if data == nil || len(data) == 0 {
		n := binary.PutUvarint(size[:], uint64(0))
		if err = write(w, size[:n]); err != nil {
			return
		}
		return
	}
	n := binary.PutUvarint(size[:], uint64(len(data)))
	if err = write(w, size[:n]); err != nil {
		return
	}
	if err = write(w, data); err != nil {
		return
	}
	return
}

func write(w io.Writer, data []byte) error {
	for i := 0; i < len(data); {
		n, err := w.Write(data[i:])
		if _, ok := err.(net.Error); !ok {
			return err
		}
		i += n
	}
	return nil
}

func recvFrame(r io.Reader) (data []byte, err error) {
	size, err := binary.ReadUvarint(r.(*bufio.Reader))
	if err != nil {
		return nil, err
	}
	if size != 0 {
		data = make([]byte, size)
		if err = read(r, data); err != nil {
			return nil, err
		}
	}
	return data, nil
}

func read(r io.Reader, data []byte) error {
	for i := 0; i < len(data); {
		n, err := r.Read(data[i:])
		if err != nil {
			if _, ok := err.(net.Error); !ok {
				return err
			}
		}
		i += n
	}
	return nil
}
