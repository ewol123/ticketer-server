package json

import (
	"github.com/ewol123/ticketer-server/ticketer-service/serializer/mapdecoder"
	 "github.com/ewol123/ticketer-server/ticketer-service/ticket"
	"github.com/fatih/structs"
	"testing"
	"time"
)

var encoded []byte

var tick = ticket.Ticket{
	Id:        "8a5e9658-f954-45c0-a232-4dcbca0d4907",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	FullName:  "Test User",
	UserId:      "7e3d3e49-b884-4803-852c-086f3a00b8ac",
	WorkerId:    "ef675295-68e2-4c8e-bf41-e05c99a46364",
	FaultType:   "leak",
	Address:     "test",
	Phone:       "36300001111",
	GeoLocation: "1.1,-1.1",
	ImageUrl:    "http://image.com/1.jpg",
	Status:      "done",
}

func TestEncode(t *testing.T) {
	serializer := &Ticket{}

	newMap := structs.Map(tick)

	enc, err := serializer.Encode(&newMap)

	if err != nil {
		t.Errorf("test encode failed, expected %v, got %v", nil, err)
	}

	t.Logf("encoded %v", enc)
	encoded = enc
}

func TestDecode(t *testing.T) {
	serializer := &Ticket{}

	decoded, err := serializer.Decode(encoded)

	if err != nil {
		t.Errorf("test decode failed, expected %v, got %v", nil, err)
	}

	t.Logf("decoded %v", decoded)

	res := &ticket.Ticket{}
	err = mapdecoder.Decode(*decoded, &res)
	if err != nil {
		t.Errorf("decoding to struct failed")
	}
	t.Logf("decoded to struct %v", res)

}
