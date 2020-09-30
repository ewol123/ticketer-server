package postgres

import (
	"github.com/ewol123/ticketer-server/ticketer-service/ticket"
	"reflect"
	"testing"
	"time"
)

var ticketModel = ticket.Ticket{
	Id:          "697b1e32-3f67-4539-b0cd-79470d38f449",
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
	UserId:      "5f058370-34b9-43e6-bcba-8686292f2b8d",
	WorkerId:    "d7b38eaf-a222-4ed3-a1f7-71a34050a7d8",
	FaultType:   "leak",
	Address:     "something",
	FullName:    "Test User",
	Phone:       "36300001111",
	GeoLocation: "1.1,-1.1",
	ImageUrl:    "http://somesite.com/images/1.jpg",
	Status:      ticket.DRAFT,
}

func TestNewPgRepository(t *testing.T) {

	repo, err := NewPgRepository("user=postgres password=test dbname=ticketer_test sslmode=disable")

	if err != nil {
		t.Errorf("test new pg repository failed, expected %v, got %v", nil, err)
	}

	i := reflect.TypeOf((*ticket.Repository)(nil)).Elem()
	isInterface := reflect.TypeOf(repo).Implements(i)

	if isInterface {
		t.Logf("test implements repository success, expected %v, got %v", true, isInterface)
	} else {
		t.Errorf("test implements repository  failed, expected %v, got %v", true, isInterface)
	}
}

func TestStore(t *testing.T) {

	repo, err := NewPgRepository("user=postgres password=test dbname=ticketer_test sslmode=disable")

	if err != nil {
		t.Errorf("test new pg repository failed, expected %v, got %v", nil, err)
	}

	res, err := repo.Store(&ticketModel)

	if err != nil {
		t.Errorf("test repository Store failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test repository Store success, expected %v, got %v", res, ticketModel)
	}
}

func TestFind(t *testing.T) {
	repo, err := NewPgRepository("user=postgres password=test dbname=ticketer_test sslmode=disable")

	if err != nil {
		t.Errorf("test new pg repository failed, expected %v, got %v", nil, err)
	}

	res, err := repo.Find("id", "697b1e32-3f67-4539-b0cd-79470d38f449")

	if err != nil {
		t.Errorf("test repository find failed, expected %v, got %v", ticketModel, err)
	}

	if res.FullName != ticketModel.FullName {
		t.Errorf("test repository find failed, expected %v, got %v", ticketModel, res)
	} else {
		t.Logf("test repository find success, expected %v, got %v", ticketModel, res)
	}
}

func TestFindAll(t *testing.T) {
	repo, err := NewPgRepository("user=postgres password=test dbname=ticketer_test sslmode=disable")

	expected := &[]ticket.Ticket{ticketModel}

	if err != nil {
		t.Errorf("test new pg repository failed, expected %v, got %v", nil, err)
	}

	rows, count, err := repo.FindAll(1, 10, "ticket_id", true, "", "","","","")

	if err != nil {
		t.Errorf("test repository find all failed, expected %v, got %v", nil, err)
	}

	if count == 0 {
		t.Errorf("test repository find all failed, expected %v, got %v", "not 0", count)
	}

	if len(*rows) == 0 {
		t.Errorf("test repository find all failed, expected %v, got %v", "slice with length", rows)
	} else {
		t.Logf("test repository find all success, expected %v, got %v", expected, rows)
	}

}

func TestUpdate(t *testing.T) {
	repo, err := NewPgRepository("user=postgres password=test dbname=ticketer_test sslmode=disable")

	if err != nil {
		t.Errorf("test new pg repository failed, expected %v, got %v", nil, err)
	}
	ticketModel.FullName = "new name"
	err = repo.Update(&ticketModel)

	if err != nil {
		t.Errorf("test repository update failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test repository update success, expected %v, got %v", nil, err)
	}
}

func TestDelete(t *testing.T) {
	repo, err := NewPgRepository("user=postgres password=test dbname=ticketer_test sslmode=disable")

	if err != nil {
		t.Errorf("test new pg repository failed, expected %v, got %v", nil, err)
	}

	err = repo.Delete("697b1e32-3f67-4539-b0cd-79470d38f449")

	if err != nil {
		t.Errorf("test repository delete failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test repository delete success, expected %v, got %v", nil, err)
	}
}
