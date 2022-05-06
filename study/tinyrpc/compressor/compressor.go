package compressor

type CompressType int32

const (
	Raw CompressType = iota
	Gzip
	Snappy
	Zlib
)

// Compressors 四种压缩器的实现
var Compressors = map[CompressType]Compressor{
	Raw:    RawCompressor{},
	Gzip:   GzipCompressor{},
	Snappy: SnappyCompressor{},
	Zlib:   ZlibCompressor{},
}

// Compressor 压缩器接口
type Compressor interface {
	Zip([]byte) ([]byte, error)
	Unzip([]byte) ([]byte, error)
}
