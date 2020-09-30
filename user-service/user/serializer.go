package user

// UserSerializer : serializer interface which connects to the http transport
type Serializer interface {
	Decode(input []byte) (*map[string]interface{}, error)
	Encode(input *map[string]interface{}) ([]byte, error)
}
