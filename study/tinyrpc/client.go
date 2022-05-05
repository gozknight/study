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

func WithCompress(c compressor.CompressType) Option {
	return compressOption(c)
}

func NewClient(conn io.ReadWriteCloser, opts ...Option) *Client {
	options := options{
		compressType: compressor.Raw,
	}
	for _, o := range opts {
		o.apply(&options)
	}
	return &Client{
		rpc.NewClientWithCodec(codec.NewClientCodec(conn, options.compressType)),
	}
}

func (c *Client) Call(serviceMethod string, args, reply any) error {
	return c.Client.Call(serviceMethod, args, reply)
}

func (c *Client) AsyncCall(serviceMethod string, arg, reply interface{}) chan *rpc.Call {
	return c.Go(serviceMethod, arg, reply, nil).Done
}
