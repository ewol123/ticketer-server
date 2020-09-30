package json

import (
	"encoding/json"
	"github.com/pkg/errors"
)

// Ticket : define our user struct so we can implement user serializer interface on it
type Ticket struct{}


// Decode : decode our input to a map
func (r *Ticket) Decode(input []byte) (*map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal(input, &m); err != nil {
		return nil, errors.Wrap(err, "serializer.Ticket.Decode")
	}
	return &m, nil
}


// Encode : encode a map to bytes
func (r *Ticket) Encode(input *map[string]interface{}) ([]byte, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Ticket.EncodeAll")
	}

	return raw, nil
}
