package json

import (
	"encoding/json"
	"github.com/pkg/errors"
)

// User : define our user struct so we can implement user serializer interface on it
type User struct{}


// Decode : decode our input to a map
func (r *User) Decode(input []byte) (*map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal(input, &m); err != nil {
		return nil, errors.Wrap(err, "serializer.User.Decode")
	}
	return &m, nil
}


// Encode : encode a map to bytes
func (r *User) Encode(input *map[string]interface{}) ([]byte, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.User.EncodeAll")
	}

	return raw, nil
}
