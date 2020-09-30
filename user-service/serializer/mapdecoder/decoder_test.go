package mapdecoder

import (
	u "github.com/ewol123/ticketer-server/user-service/user"
	"testing"
	"time"
)

func TestDecode(t *testing.T) {

	newMap := make(map[string]interface{})
	newMap["Id"] = "8a5e9658-f954-45c0-a232-4dcbca0d4907"
	newMap["CreatedAt"] = time.Now()
	newMap["UpdatedAt"] = time.Now()

	res := &u.User{}
	err := Decode(newMap, &res)
	if err != nil {
		t.Errorf("decoding to struct failed")
	}
	t.Logf("decoded to struct %v", res)
}