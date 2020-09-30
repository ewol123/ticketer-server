package postgres

import (
	"reflect"
	"testing"
	"time"

	u "github.com/ewol123/ticketer-server/user-service/user"
)

var usr = u.User{
	Id:               "8a5e9658-f954-45c0-a232-4dcbca0d4907",
	CreatedAt:        time.Now(),
	UpdatedAt:        time.Now(),
	FullName:         "Test User",
	Email:            "test.user@test.com",
	Password:         "bcrypt",
	Status:           u.PENDING,
	RegistrationCode: "123456",
}

func TestNewPgRepository(t *testing.T) {

	repo, err := NewPgRepository("user=postgres password=test dbname=user_test sslmode=disable")

	if err != nil {
		t.Errorf("test new pg repository failed, expected %v, got %v", nil, err)
	}

	i := reflect.TypeOf((*u.Repository)(nil)).Elem()
	isInterface := reflect.TypeOf(repo).Implements(i)

	if isInterface {
		t.Logf("test implements repository success, expected %v, got %v", true, isInterface)
	} else {
		t.Errorf("test implements repository  failed, expected %v, got %v", true, isInterface)
	}
}

func TestStore(t *testing.T) {

	repo, err := NewPgRepository("user=postgres password=test dbname=user_test sslmode=disable")

	if err != nil {
		t.Errorf("test new pg repository failed, expected %v, got %v", nil, err)
	}

	res, err := repo.Store(&usr)

	if err != nil {
		t.Errorf("test repository Store failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test repository Store success, expected %v, got %v", res, usr)
	}
}

func TestFind(t *testing.T) {
	repo, err := NewPgRepository("user=postgres password=test dbname=user_test sslmode=disable")

	if err != nil {
		t.Errorf("test new pg repository failed, expected %v, got %v", nil, err)
	}

	rows, err := repo.Find("id", "8a5e9658-f954-45c0-a232-4dcbca0d4907")

	if err != nil {
		t.Errorf("test repository find failed, expected %v, got %v", usr, err)
	}

	if rows.Email != usr.Email {
		t.Errorf("test repository find failed, expected %v, got %v", usr, rows)
	} else {
		t.Logf("test repository find success, expected %v, got %v", usr, rows)
	}
}

func TestFindAll(t *testing.T) {
	repo, err := NewPgRepository("user=postgres password=test dbname=user_test sslmode=disable")

	expected := &[]u.User{usr}

	if err != nil {
		t.Errorf("test new pg repository failed, expected %v, got %v", nil, err)
	}

	rows, count, err := repo.FindAll(1, 10, "user_id", true, "test")

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
	repo, err := NewPgRepository("user=postgres password=test dbname=user_test sslmode=disable")

	if err != nil {
		t.Errorf("test new pg repository failed, expected %v, got %v", nil, err)
	}
	usr.Email = "asd@gmail.com"
	err = repo.Update(&usr)

	if err != nil {
		t.Errorf("test repository update failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test repository update success, expected %v, got %v", nil, err)
	}
}

func TestDelete(t *testing.T) {
	repo, err := NewPgRepository("user=postgres password=test dbname=user_test sslmode=disable")

	if err != nil {
		t.Errorf("test new pg repository failed, expected %v, got %v", nil, err)
	}

	err = repo.Delete("8a5e9658-f954-45c0-a232-4dcbca0d4907")

	if err != nil {
		t.Errorf("test repository delete failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test repository delete success, expected %v, got %v", nil, err)
	}
}
