package codec

import (
	"bufio"
	"github.com/golang/protobuf/proto"
	"gozknight.top/tinyrpc/compressor"
	"gozknight.top/tinyrpc/errs"
	"gozknight.top/tinyrpc/header"
	"gozknight.top/tinyrpc/serializer"
	"hash/crc32"
	"io"
	"net/rpc"
	"sync"
)

type serverCodec struct {
	r       io.Reader
	w       io.Writer
	c       io.Closer
	request header.RequestHeader
	mutex   sync.Mutex
	seq     uint64
	pending map[uint64]uint64
}

func NewServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec {
	return &serverCodec{
		r:       bufio.NewReader(conn),
		w:       bufio.NewWriter(conn),
		c:       conn,
		pending: make(map[uint64]uint64),
	}
}

func (s *serverCodec) ReadRequestHeader(r *rpc.Request) error {
	s.request.ResetHeader()
	err := reeadRequestHeader(s.r, &s.request)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.seq++
	s.pending[s.seq] = s.request.Id
	r.ServiceMethod = s.request.Method
	r.Seq = s.seq
	return nil
}

func reeadRequestHeader(r io.Reader, h *header.RequestHeader) error {
	pbHeader, err := recvFrame(r)
	if err != nil {
		return err
	}
	if err = proto.Unmarshal(pbHeader, h); err != nil {
		return err
	}
	return nil
}

func (s *serverCodec) ReadRequestBody(param any) error {
	if param != nil {
		if s.request.RequestLen != 0 {
			if err := read(s.r, make([]byte, s.request.RequestLen)); err != nil {
				return err
			}
		}
		return nil
	}
	if err := readRequestBody(s.r, &s.request, param); err != nil {
		return err
	}
	return nil
}

func readRequestBody(r io.Reader, h *header.RequestHeader, param interface{}) error {
	reqBody := make([]byte, h.RequestLen)
	err := read(r, reqBody)
	if err != nil {
		return err
	}
	// 校验
	if h.Checksum != 0 {
		if crc32.ChecksumIEEE(reqBody) != h.Checksum {
			return errs.UnexpectedChecksumError
		}
	}
	// 判断压缩器是否存在
	if _, ok := compressor.Compressors[compressor.CompressType(h.CompressType)]; !ok {
		return errs.NotFoundCompressorError
	}
	// 解压
	req, err := compressor.Compressors[compressor.CompressType(h.CompressType)].Unzip(reqBody)
	if err != nil {
		return err
	}
	return serializer.Serializers[serializer.Proto].Unmarshal(req, param)
}

func (s *serverCodec) WriteResponse(r *rpc.Response, param any) error {
	s.mutex.Lock()
	id, ok := s.pending[r.Seq]
	if !ok {
		s.mutex.Unlock()
		return errs.InvalidSequenceError
	}
	delete(s.pending, r.Seq)
	s.mutex.Unlock()
	if err := writeResponse(s.w, id, r.Error, compressor.CompressType(s.request.CompressType), param); err != nil {
		return err
	}
	return nil
}

func writeResponse(w io.Writer, id uint64, serr string, compressType compressor.CompressType, param any) (err error) {
	// 如果RPC调用结果有误，把param置为nil
	if serr != "" {
		param = nil
	}
	// 判断压缩器是否存在
	if _, ok := compressor.Compressors[compressType]; !ok {
		return errs.NotFoundCompressorError
	}
	var respBody []byte
	if param != nil {
		respBody, err = serializer.Serializers[serializer.Proto].Marshal(param)
		if err != nil {
			return err
		}
	}
	// 压缩
	compressBody, err := compressor.Compressors[compressType].Zip(respBody)
	if err != nil {
		return err
	}
	h := header.ResponsePool.Get().(*header.ResponseHeader)
	defer func() {
		h.ResetHeader()
		header.ResponsePool.Put(h)
	}()
	h.Id = id
	h.Error = serr
	h.ResponseLen = uint32(len(compressBody))
	h.Checksum = crc32.ChecksumIEEE(compressBody)
	h.CompressType = header.Compress(compressType)
	pbHeader, err := proto.Marshal(h)
	if err != nil {
		return err
	}
	if err = sendFrame(w, pbHeader); err != nil {
		return err
	}
	if err = write(w, compressBody); err != nil {
		return err
	}
	if err = w.(*bufio.Writer).Flush(); err != nil {
		return err
	}
	w.(*bufio.Writer).Flush()
	return nil
}

func (s *serverCodec) Close() error {
	return s.c.Close()
}
