package serializer

import (
	"github.com/golang/protobuf/proto"
	"gozknight.top/tinyrpc/errs"
)

type ProtoSerializer struct{}

// Marshal .
func (_ ProtoSerializer) Marshal(message interface{}) ([]byte, error) {
	var body proto.Message
	if message == nil {
		return []byte{}, nil
	}
	var ok bool
	if body, ok = message.(proto.Message); !ok {
		return nil, errs.NotImplementProtoMessageError
	}
	return proto.Marshal(body)
}

// Unmarshal .
func (_ ProtoSerializer) Unmarshal(data []byte, message interface{}) error {
	var body proto.Message
	if message == nil {
		return nil
	}

	var ok bool
	body, ok = message.(proto.Message)
	if !ok {
		return errs.NotImplementProtoMessageError
	}

	return proto.Unmarshal(data, body)
}
