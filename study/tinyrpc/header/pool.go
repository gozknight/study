package header

import "sync"

var (
	RequestPool  sync.Pool
	ResponsePool sync.Pool
)

func init() {
	RequestPool = sync.Pool{New: func() interface{} {
		return &RequestHeader{}
	}}
	ResponsePool = sync.Pool{New: func() interface{} {
		return &ResponseHeader{}
	}}
}

func (req *RequestHeader) ResetHeader() {
	req.Id = 0
	req.Checksum = 0
	req.Method = ""
	req.RequestLen = 0
	req.CompressType = 0
}

func (resp *ResponseHeader) ResetHeader() {
	resp.ResponseLen = 0
	resp.CompressType = 0
	resp.Id = 0
	resp.Checksum = 0
	resp.Error = ""
}
