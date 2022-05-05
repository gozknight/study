package serializer

type SerializeType int32

const (
	Proto SerializeType = iota
)

var Serializers = map[SerializeType]Serializer{
	Proto: ProtoSerializer{},
}

type Serializer interface {
	Marshal(msg interface{}) ([]byte, error)
	Unmarshal(data []byte, msg interface{}) error
}
