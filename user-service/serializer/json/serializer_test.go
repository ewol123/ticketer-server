package json

import (
	u "github.com/ewol123/ticketer-server/user-service/user"
	"testing"
)

var usr = u.User{
	Id:        "8a5e9658-f954-45c0-a232-4dcbca0d4907",
	CreatedAt: 0,
	UpdatedAt: 0,
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

	encoded, err := serializer.Encode(&usr)

	if err != nil {
		t.Errorf("test if encode failed, expected %v, got %v", nil, err)
	}

	t.Logf("encoded %v", encoded)
}

func TestEncodeAll(t *testing.T){
	serializer := &User{}

	encoded, err := serializer.EncodeAll(&roles)

	if err != nil {
		t.Errorf("test if encode all failed, expected %v, got %v", nil, err)
	}

	t.Logf("encoded all %v", encoded)
}

func TestDecode(t *testing.T) {
	serializer := &User{}

	bytes := []byte{123,34,105,100,34,58,34,56,97,53,101,57,54,53,56,45,102,57,53,52,45,52,53,99,48,45,97,50,51,50,45,52,100,99,98,99,97,48,100,52,57,48,55,34,44,34,99,114,101,97,116,101,100,95,97,116,34,58,48,44,34,117,112,100,97,116,101,100,95,97,116,34,58,48,44,34,102,117,108,108,95,110,97,109,101,34,58,34,84,101,115,116,32,85,115,101,114,34,44,34,101,109,97,105,108,34,58,34,116,101,115,116,46,117,115,101,114,64,116,101,115,116,46,99,111,109,34,44,34,112,97,115,115,119,111,114,100,34,58,34,98,99,114,121,112,116,34,44,34,114,111,108,101,115,34,58,91,123,34,105,100,34,58,34,55,101,51,100,51,101,52,57,45,98,56,56,52,45,52,56,48,51,45,56,53,50,99,45,48,56,54,102,51,97,48,48,98,56,97,99,34,44,34,110,97,109,101,34,58,34,117,115,101,114,34,125,44,123,34,105,100,34,58,34,101,102,54,55,53,50,57,53,45,54,56,101,50,45,52,99,56,101,45,98,102,52,49,45,101,48,53,99,57,57,97,52,54,51,54,52,34,44,34,110,97,109,101,34,58,34,97,100,109,105,110,34,125,93,125}
	decoded, err := serializer.Decode(bytes)

	if err != nil {
		t.Errorf("test if encode failed, expected %v, got %v", nil, err)
	}

	t.Logf("decoded %v", decoded)

}
