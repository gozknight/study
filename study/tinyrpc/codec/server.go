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
	r io.Reader
	w io.Writer
	c io.Closer

	request header.RequestHeader

	mutex   sync.Mutex // 保护 seq, pending
	seq     uint64     // 一个自增的序号
	pending map[uint64]uint64
}

// NewServerCodec Create a new server codec
func NewServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec {
	return &serverCodec{
		r:       bufio.NewReader(conn),
		w:       bufio.NewWriter(conn),
		c:       conn,
		pending: make(map[uint64]uint64),
	}
}

func (s *serverCodec) ReadRequestHeader(r *rpc.Request) error {
	s.request.ResetHeader() // 重置serverCodec结构体的请求头部
	err := readRequestHeader(s.r, &s.request)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.seq++                            // 序号自增
	s.pending[s.seq] = s.request.Id    // 自增序号与请求头部的ID进行绑定
	r.ServiceMethod = s.request.Method // 填充 r.ServiceMethod
	r.Seq = s.seq                      // 填充 r.Seq
	return nil
}

func readRequestHeader(r io.Reader, h *header.RequestHeader) (err error) {
	pbHeader, err := recvFrame(r) // 读取请求头部文件
	if err != nil {
		return err
	}
	//将字节串Unmarshal成RequestHeader结构类型
	if err = proto.Unmarshal(pbHeader, h); err != nil {
		return err
	}
	return nil
}

func (s *serverCodec) ReadRequestBody(param any) error {
	if param == nil {
		if s.request.RequestLen != 0 { // 废弃多余部分
			if err := read(s.r, make([]byte, s.request.RequestLen)); err != nil {
				return err
			}
		}
		return nil
	}

	if err := readRequestBody(s.r, &s.request, param); err != nil {
		return nil
	}
	return nil
}

func readRequestBody(r io.Reader, h *header.RequestHeader, param any) error {
	reqBody := make([]byte, h.RequestLen)
	err := read(r, reqBody) // 根据请求体的大小，读取该大小的字节串
	if err != nil {
		return err
	}
	if h.Checksum != 0 { // 校验
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
	// 把字节串反序列化成proto结构
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

func writeResponse(w io.Writer, id uint64, serr string,
	compressType compressor.CompressType, param any) (err error) {
	if serr != "" {
		param = nil // 如果RPC调用结果有误，把param置为nil
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
	compressedRespBody, err := compressor.Compressors[compressType].Zip(respBody)
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
	h.ResponseLen = uint32(len(compressedRespBody))
	h.Checksum = crc32.ChecksumIEEE(compressedRespBody)
	h.CompressType = header.Compress(compressType)

	pbHeader, err := proto.Marshal(h)
	if err != err {
		return
	}
	// 发送响应头
	if err = sendFrame(w, pbHeader); err != nil {
		return
	}
	// 发送响应体
	if err = write(w, compressedRespBody); err != nil {
		return
	}
	w.(*bufio.Writer).Flush()
	return nil
}

func (s *serverCodec) Close() error {
	return s.c.Close()
}
