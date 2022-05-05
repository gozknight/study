package compressor

type CompressType int32

const (
	Raw CompressType = iota
	Gzip
	Snappy
	Zlib
)

// Compressors TODO: implement compresssor
var Compressors = map[CompressType]Compressor{
	Raw:    RawCompressor{},
	Gzip:   GzipCompressor{},
	Snappy: SnappyCompressor{},
	Zlib:   ZlibCompressor{},
}

type Compressor interface {
	Zip([]byte) ([]byte, error)
	Unzip([]byte) ([]byte, error)
}
