package json

import (
	"github.com/ewol123/ticketer-server/user-service/serializer/mapdecoder"
	u "github.com/ewol123/ticketer-server/user-service/user"
	"github.com/fatih/structs"
	"testing"
	"time"
)

var encoded []byte

var usr = u.User{
	Id:        "8a5e9658-f954-45c0-a232-4dcbca0d4907",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	FullName:  "Test User",
	Email:     "test.user@test.com",
	Password:  "bcrypt",
	Roles:     roles,
}

var users = []u.User{
	usr,
}

var roles  = []u.Role{
	{Id: "7e3d3e49-b884-4803-852c-086f3a00b8ac", Name: "user" },
	{Id: "ef675295-68e2-4c8e-bf41-e05c99a46364", Name: "admin" },

}


func TestEncode(t *testing.T) {
	serializer := &User{}

	newMap := structs.Map(usr)

	enc, err := serializer.Encode(&newMap)

	if err != nil {
		t.Errorf("test encode failed, expected %v, got %v", nil, err)
	}

	t.Logf("encoded %v", enc)
	encoded = enc
}

func TestDecode(t *testing.T) {
	serializer := &User{}

	decoded, err := serializer.Decode(encoded)

	if err != nil {
		t.Errorf("test decode failed, expected %v, got %v", nil, err)
	}

	t.Logf("decoded %v", decoded)

	res := &u.User{}
	err = mapdecoder.Decode(*decoded, &res)
	if err != nil {
		t.Errorf("decoding to struct failed")
	}
	t.Logf("decoded to struct %v", res)

}
