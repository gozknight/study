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

type clientCodec struct {
	r          io.Reader
	w          io.Writer
	c          io.Closer
	compressor compressor.CompressType
	response   header.ResponseHeader
	mutex      sync.Mutex
	pending    map[uint64]string
}

func NewClientCodec(coon io.ReadWriteCloser, compressType compressor.CompressType) rpc.ClientCodec {
	return &clientCodec{
		r:          bufio.NewReader(coon),
		w:          bufio.NewWriter(coon),
		c:          coon,
		compressor: compressType,
		pending:    make(map[uint64]string),
	}
}

func (c *clientCodec) WriteRequest(r *rpc.Request, param interface{}) error {
	c.mutex.Lock()
	c.pending[r.Seq] = r.ServiceMethod
	c.mutex.Unlock()
	err := writeRequest(c.w, r, c.compressor, param)
	if err != nil {
		return err
	}
	return nil
}

func writeRequest(w io.Writer, r *rpc.Request, compressType compressor.CompressType, param interface{}) error {
	if _, ok := compressor.Compressors[compressType]; !ok {
		return errs.NotFoundCompressorError
	}
	reqBody, err := serializer.Serializers[serializer.Proto].Marshal(param)
	if err != nil {
		return err
	}
	compressBody, err := compressor.Compressors[compressType].Zip(reqBody)
	if err != nil {
		return err
	}
	h := header.RequestPool.Get().(*header.RequestHeader)
	defer func() {
		h.ResetHeader()
		header.RequestPool.Put(h)
	}()
	h.Id = r.Seq
	h.Method = r.ServiceMethod
	h.RequestLen = uint32(len(compressBody))
	h.CompressType = header.Compress(compressType)
	h.Checksum = crc32.ChecksumIEEE(compressBody)
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
	return nil
}

func (c *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	c.response.ResetHeader()
	err := readResponseHeader(c.r, &c.response)
	if err != nil {
		return err
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	r.Seq = c.response.Id
	r.Error = c.response.Error
	r.ServiceMethod = c.pending[r.Seq]
	delete(c.pending, r.Seq)
	return nil
}

func readResponseHeader(r io.Reader, h *header.ResponseHeader) error {
	pbHeader, err := recvFrame(r)
	if err != nil {
		return err
	}
	if err := proto.Unmarshal(pbHeader, h); err != nil {
		return err
	}
	return nil
}

func (c *clientCodec) ReadResponseBody(param interface{}) error {
	if param == nil {
		if c.response.ResponseLen != 0 {
			if err := read(c.r, make([]byte, c.response.ResponseLen)); err != nil {
				return err
			}
		}
		return nil
	}
	if err := readResponseBody(c.r, &c.response, param); err != nil {
		return err
	}
	return nil
}

func readResponseBody(r io.Reader, h *header.ResponseHeader, param interface{}) error {
	respBody := make([]byte, h.ResponseLen)
	err := read(r, respBody)
	if err != nil {
		return err
	}
	if h.Checksum != 0 {
		if crc32.ChecksumIEEE(respBody) != h.Checksum {
			return errs.UnexpectedChecksumError
		}
	}
	if _, ok := compressor.Compressors[compressor.CompressType(h.CompressType)]; !ok {
		return errs.NotFoundCompressorError
	}
	resp, err := compressor.Compressors[compressor.CompressType(h.CompressType)].Unzip(respBody)
	if err != nil {
		return err
	}
	return serializer.Serializers[serializer.Proto].Unmarshal(resp, param)
}

func (c *clientCodec) Close() error {
	return c.c.Close()
}
