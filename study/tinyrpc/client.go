package tinyrpc

import (
	"gozknight.top/tinyrpc/codec"
	"gozknight.top/tinyrpc/compressor"
	"io"
	"net/rpc"
)

type Client struct {
	*rpc.Client
}

// Option 选项接口
type Option interface {
	apply(*options)
}

type options struct {
	compressType compressor.CompressType
}

type compressOption compressor.CompressType

func (c compressOption) apply(opt *options) {
	opt.compressType = compressor.CompressType(c)
}

// WithCompress 设置压缩类型
func WithCompress(c compressor.CompressType) Option {
	return compressOption(c)
}

func NewClient(conn io.ReadWriteCloser, opts ...Option) *Client {
	options := options{
		compressType: compressor.Raw,
	}
	for _, o := range opts { // 应用选项
		o.apply(&options)
	}
	return &Client{rpc.NewClientWithCodec(
		codec.NewClientCodec(conn, options.compressType))} // 使用TinyRPC的解码器
}

func (c *Client) Call(serviceMethod string, args any, reply any) error {
	return c.Client.Call(serviceMethod, args, reply)
}

func (c *Client) AsyncCall(serviceMethod string, args any, reply any) chan *rpc.Call {
	return c.Go(serviceMethod, args, reply, nil).Done
}
