package user

import (
	"reflect"
	"testing"
)

var user = User{
	Id:        "8a5e9658-f954-45c0-a232-4dcbca0d4907",
	CreatedAt: 0,
	UpdatedAt: 0,
	FullName:  "Test User",
	Email:     "test.user@test.com",
	Password:  "bcrypt",
	Roles:     roles,
}

var users = []User{
	user,
}

var roles  = []Role{
	{Id: "7e3d3e49-b884-4803-852c-086f3a00b8ac", Name: "user" },
	{Id: "ef675295-68e2-4c8e-bf41-e05c99a46364", Name: "admin" },

}

type newRepo struct {}

func (r* newRepo) Find(id string) (*User, error) {
	for i := range users {
		if users[i].Id == id {
			return &users[i], nil
		}
	}
	return nil, ErrUserNotFound
}

func (r* newRepo) FindAll(page int, rowsPerPage int, sortBy string, descending bool) (*[]User, int, error){
	offsetPage := page - 1

	if rowsPerPage * page > len(users) -1 {
		return nil, len(users), nil
	}

	pagination := users[offsetPage * rowsPerPage: page * rowsPerPage]
	return &pagination, len(users), nil

}

func (r* newRepo) Update(user *User) error {
	for i := range users {
		if users[i].Id == user.Id {
			users[i] = *user
			return nil
		}
	}
	return ErrUserNotFound
}

func (r* newRepo) Delete(id string) error{
	for i := range users {
		if users[i].Id == id {
			users[i] = users[len(users)-1]
			users[len(users)-1] = User{}
			users = users[:len(users)-1]
			return nil
		}
	}
	return ErrUserNotFound
}

func (r* newRepo) Store(user *User) error {
	users = append(users, *user)
	return nil
}


func TestRepository(t *testing.T) {

	r := &newRepo{}
	i := reflect.TypeOf((*Repository)(nil)).Elem()
	isInterface := reflect.TypeOf(r).Implements(i)

	if isInterface {
		t.Logf("test if implements repository success, expected %v, got %v", true, isInterface)
	} else {
		t.Errorf("test if implements repository failed, expected %v, got %v", true, isInterface)
	}
}

func TestNewUserService(t *testing.T) {

	r := &newRepo{}
	service := NewUserService(r)

	i := reflect.TypeOf((*Service)(nil)).Elem()
	isInterface := reflect.TypeOf(service).Implements(i)

	if isInterface {
		t.Logf("test if implements service success, expected %v, got %v", true, isInterface)
	} else {
		t.Errorf("test if implements service  failed, expected %v, got %v", true, isInterface)
	}
}

func TestFind(t *testing.T){
	r := &newRepo{}
	service := NewUserService(r)

	shouldFind, err := service.Find("8a5e9658-f954-45c0-a232-4dcbca0d4907")

	if err != nil {
		t.Errorf("test if found failed, expected %v, got %v", user, err)
	}

	if shouldFind.Id != user.Id {
		t.Errorf("test if found failed, expected %v, got %v", user, shouldFind)
	} else {
		t.Logf("test if found success, expected %v, got %v", user, shouldFind)

	}

	_, err = service.Find("abc")

	if err != nil {
		if err != ErrUserNotFound {
			t.Errorf("test if not found failed, expected %v, got %v", ErrUserNotFound, err)
		} else {
			t.Logf("test if not found success, expected %v, got %v", ErrUserNotFound, err)
		}
	}
}

func TestStore(t *testing.T) {
	r := &newRepo{}
	service := NewUserService(r)

	newUser := User{
		Id:        "94996a7a-312d-405b-9376-eb1850359632",
		CreatedAt: 0,
		UpdatedAt: 0,
		FullName:  "test 2",
		Email:     "test2@test.com",
		Password:  "bcrypt2",
		Roles:     roles,
	}

	invalidUser := User{
		FullName: "hallo",
	}
	
	err := service.Store(&newUser)

	if err != nil{
		t.Errorf("test user store failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test user store success, expected %v, got %v", nil, err)
	}

	err = service.Store(&invalidUser)

	if err != nil{
		t.Logf("test user store invalid success, expected %v, got %v", ErrUserInvalid, err)
	} else {
		t.Errorf("test user store invalid failed, expected %v, got %v", ErrUserInvalid, err)
	}
}

func TestUpdate(t *testing.T) {
	r := &newRepo{}
	service := NewUserService(r)

	updateUser := User{
		Id:        "8a5e9658-f954-45c0-a232-4dcbca0d4907",
		FullName:  "updated",
		Email:     "updated@test.com",
		Password:  "bcrypt2",
		Roles:     roles,
	}

	invalidUser := User{
		FullName: "hallo",
	}

	err := service.Update(&updateUser)

	if err != nil{
		t.Errorf("test user update failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test user update success, expected %v, got %v", nil, err)
	}


	shouldFind, err := service.Find("8a5e9658-f954-45c0-a232-4dcbca0d4907")

	if err != nil {
		t.Errorf("test user update failed, expected %v after update, got %v",updateUser, shouldFind)
	} else {
		if shouldFind.FullName != "updated" {
			t.Errorf("test user update failed, expected FullName = %v after update, got %v",updateUser.FullName, shouldFind.FullName)
		}
		if shouldFind.Email != "updated@test.com"{
			t.Errorf("test user update failed, expected Email = %v after update, got %v",updateUser.Email, shouldFind.Email)
		}
		if shouldFind.Password != "bcrypt2" {
			t.Errorf("test user update failed, expected Password = %v after update, got %v",updateUser.Password, shouldFind.Password)
		}
			t.Logf("test user update success, all values are as expected")
	}

	err = service.Update(&invalidUser)

	if err != nil{
		t.Logf("test user update invalid success, expected %v, got %v", ErrUserInvalid, err)
	} else {
		t.Errorf("test user update invalid failed, expected %v, got %v", ErrUserInvalid, err)
	}
}

func TestDelete(t *testing.T){
	r := &newRepo{}
	service := NewUserService(r)

	err := service.Delete("8a5e9658-f954-45c0-a232-4dcbca0d4907")

	if err != nil {
		t.Errorf("test delete failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test delete success, expected %v, got %v", nil, err)
	}


	err = service.Delete("abc")

	if err != nil {
		if err != ErrUserNotFound {
			t.Errorf("test delete failed, expected %v, got %v", ErrUserNotFound, err)
		} else {
			t.Logf("test delete success, expected %v, got %v", ErrUserNotFound, err)
		}
	}

}