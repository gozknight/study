package header

import "sync"

var (
	RequestPool  sync.Pool
	ResponsePool sync.Pool
)

func init() {
	RequestPool = sync.Pool{New: func() any {
		return &RequestHeader{}
	}}
	ResponsePool = sync.Pool{New: func() any {
		return &ResponseHeader{}
	}}
}

// ResetHeader reset request header
func (h *RequestHeader) ResetHeader() {
	h.Id = 0
	h.Checksum = 0
	h.Method = ""
	h.CompressType = 0
	h.RequestLen = 0
}

// ResetHeader reset response header
func (h *ResponseHeader) ResetHeader() {
	h.Error = ""
	h.Id = 0
	h.CompressType = 0
	h.Checksum = 0
	h.ResponseLen = 0
}
