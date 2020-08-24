package json

import (
	"encoding/json"
	"github.com/ewol123/ticketer-server/user-service/user"
	"github.com/pkg/errors"
)

// User : define our user struct so we can implement user serializer interface on it
type User struct{}

// Decode : decode our input to a User
func (r *User) Decode(input []byte) (*user.User, error) {
	user := &user.User{}
	if err := json.Unmarshal(input, user); err != nil {
		return nil, errors.Wrap(err, "serializer.User.Decode")
	}
	return user, nil
}

// Encode : encode our input to bytes
func (r *User) Encode(input *user.User) ([]byte, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.User.Encode")
	}
	return raw, nil
}

// EncodeAll : encode any interface to bytes
func (r *User) EncodeAll(input interface{}) ([]byte, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.User.EncodeAll")
	}

	return raw, nil
}
