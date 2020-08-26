package user

// UserSerializer : serializer interface which connects to the http transport
type Serializer interface {
	Decode(input []byte) (*User, error)
	Encode(input *User) ([]byte, error)
	EncodeAll(input interface{}) ([]byte, error)
}
