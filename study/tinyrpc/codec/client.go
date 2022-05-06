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
	r io.Reader
	w io.Writer
	c io.Closer

	compressor compressor.CompressType // rpc compress type(raw,gzip,snappy,zlib)
	response   header.ResponseHeader   // rpc response header
	mutex      sync.Mutex              // protect pending map
	pending    map[uint64]string
}

func NewClientCodec(conn io.ReadWriteCloser,
	compressType compressor.CompressType) rpc.ClientCodec {
	return &clientCodec{
		r:          bufio.NewReader(conn),
		w:          bufio.NewWriter(conn),
		c:          conn,
		compressor: compressType,
		pending:    make(map[uint64]string),
	}
}

func (c *clientCodec) WriteRequest(r *rpc.Request, param any) error {
	c.mutex.Lock()
	c.pending[r.Seq] = r.ServiceMethod
	c.mutex.Unlock()
	if err := writeRequest(c.w, r, c.compressor, param); err != nil {
		return err
	}
	return nil
}

func writeRequest(w io.Writer, r *rpc.Request,
	compressType compressor.CompressType, param any) error {
	// 判断压缩器是否存在
	if _, ok := compressor.Compressors[compressType]; !ok {
		return errs.NotFoundCompressorError
	}
	// 用Protobuf序列化器进行编码
	reqBody, err := serializer.Serializers[serializer.Proto].Marshal(param)
	if err != nil {
		return err
	}
	// 压缩
	compressedReqBody, err := compressor.Compressors[compressType].Zip(reqBody)
	if err != nil {
		return err
	}
	// 从请求头部对象池取出请求头
	h := header.RequestPool.Get().(*header.RequestHeader)
	defer func() {
		h.ResetHeader()
		header.RequestPool.Put(h)
	}()
	h.Id = r.Seq
	h.Method = r.ServiceMethod
	h.RequestLen = uint32(len(compressedReqBody))
	h.CompressType = header.Compress(compressType)
	h.Checksum = crc32.ChecksumIEEE(compressedReqBody)
	// 请求头部使用protobuf 编码
	pbHeader, err := proto.Marshal(h)
	if err != err {
		return err
	}
	if err = sendFrame(w, pbHeader); err != nil { // 发送头部
		return err
	}
	if err = write(w, compressedReqBody); err != nil { //发送请求体
		return err
	}

	w.(*bufio.Writer).Flush()
	return nil
}

func (c *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	c.response.ResetHeader() // 重置clientCodec的响应头部
	err := readResponseHeader(c.r, &c.response)
	if err != nil {
		return err
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	r.Seq = c.response.Id              // 填充 r.Seq
	r.Error = c.response.Error         // 填充 r.Error
	r.ServiceMethod = c.pending[r.Seq] // 根据序号填充 r.ServiceMethod
	delete(c.pending, r.Seq)           // 删除pending里的序号
	return nil
}

func readResponseHeader(r io.Reader, h *header.ResponseHeader) error {
	pbHeader, err := recvFrame(r) //从IO中读取响应头部
	if err != nil {
		return err
	}
	//将字节串Unmarshal成ResponseHeader结构类型
	if err = proto.Unmarshal(pbHeader, h); err != nil {
		return err
	}
	return nil
}

func (c *clientCodec) ReadResponseBody(param any) error {
	if param == nil {
		if c.response.ResponseLen != 0 { // 废弃多余部分
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

func readResponseBody(r io.Reader, h *header.ResponseHeader, param any) error {
	respBody := make([]byte, h.ResponseLen) // 根据响应体长度，读取该长度的字节串
	err := read(r, respBody)
	if err != nil {
		return err
	}
	// 校验
	if h.Checksum != 0 {
		if crc32.ChecksumIEEE(respBody) != h.Checksum {
			return errs.UnexpectedChecksumError
		}
	}
	// 判断压缩器是否存在
	if _, ok := compressor.Compressors[compressor.CompressType(h.CompressType)]; !ok {
		return errs.NotFoundCompressorError
	}
	// 解压
	resp, err := compressor.Compressors[compressor.CompressType(h.CompressType)].Unzip(respBody)
	if err != nil {
		return err
	}
	// 把字节串反序列化成proto结构
	return serializer.Serializers[serializer.Proto].Unmarshal(resp, param)
}

func (c *clientCodec) Close() error {
	return c.c.Close()
}
