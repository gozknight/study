package serializer

type SerializeType int32

const (
	Proto SerializeType = iota
)

var Serializers = map[SerializeType]Serializer{
	Proto: ProtoSerializer{},
}

type Serializer interface {
	Marshal(message any) ([]byte, error)
	Unmarshal(data []byte, message any) error
}
