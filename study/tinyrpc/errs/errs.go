package errs

import "errors"

var (
	NotFoundCompressorError       = errors.New("NotFoundCompressor")
	UnexpectedChecksumError       = errors.New("UnexpectedChecksumError")
	InvalidSequenceError          = errors.New("InvalidSequenceError")
	NotImplementProtoMessageError = errors.New("param does not implement proto.Message")
)
