package user

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"testing"
)

var roles  = []Role{
	{ Id: "72daf87a-fda4-4c72-aff9-85edd68d155f", Name: "user" },
	{ Id: "336a3ff6-9fdb-496f-ac8c-e37759969cf2", Name: "admin" },
}

var user = User{
	Id:        "8a5e9658-f954-45c0-a232-4dcbca0d4907",
	FullName:  "Test User",
	Email:     "test.user@test.com",
	Password:  "$2y$10$2wcRIg0zbCk8b02HVgb3bui6Dkd8xeCZEmBhAbY8yfJ8NtzbzABk2",
	RegistrationCode: "212345",
	ResetPasswordCode: "123456",
	Status: ACTIVE,
	Roles: roles,
}

var user2 = User{
	Id:        "2a5e9658-f954-45c0-a232-4dcbca0d4904",
	FullName:  "Test User",
	Email:     "test.user2@test.com",
	Password:  "bcrypt",
	RegistrationCode: "111111",
	ResetPasswordCode: "333333",
	Status: ACTIVE,
}

var user3 = User{
	Id:        "2a5e9658-d954-45c0-a232-4dcbca0d4904",
	FullName:  "Test User",
	Email:     "test.user3@test.com",
	Password:  "something",
	RegistrationCode: "111111",
	ResetPasswordCode: "333333",
	Status: ACTIVE,
}

var user4 = User {
	Id:        "2a5e9658-d954-45c0-a232-2dcbca0d4904",
	FullName:  "Test User",
	Email:     "test.user@resetpw.com",
	Password:  "something",
	RegistrationCode: "111111",
	ResetPasswordCode: "333333",
	Status: ACTIVE,
}


var users = []User{
	user,
	user2,
	user3,
	user4,
}


type newRepo struct {}

func (r* newRepo) Find(column string, value string) (*User, error) {

	//hack for tests
	firstUpper := strings.Title(column)

	for i := range users {
		user := reflect.ValueOf(&users[i])

		colVal := reflect.Indirect(user).FieldByName(firstUpper)
		if colVal.String() == value {
			return &users[i], nil
		}
	}
	return nil, ErrUserNotFound
}

func (r* newRepo) FindAll(page int, rowsPerPage int, sortBy string, descending bool, filter string) (*[]User, int, error){
	offsetPage := page - 1

	if rowsPerPage * page > len(users) -1 {
		return &users, len(users), nil
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

func (r* newRepo) Store(user *User) (*User, error) {
	users = append(users, *user)
	return user, nil
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

func TestGetUser(t *testing.T){
	r := &newRepo{}
	service := NewUserService(r)

	getUserModel := GetUserRequestModel{Id: "8a5e9658-f954-45c0-a232-4dcbca0d4907" }

	shouldFind, err := service.GetUser(&getUserModel)

	if err != nil {
		t.Errorf("test found failed, expected %v, got %v", user, err)
	}

	if shouldFind.Id != user.Id {
		t.Errorf("test found failed, expected %v, got %v", user, shouldFind)
	} else {
		t.Logf("test found success, expected %v, got %v", user, shouldFind)
	}

	wrongModel := GetUserRequestModel{Id: "abc" }

	_, err = service.GetUser(&wrongModel)

	if err != nil {
		if errors.Cause(err) != ErrRequestInvalid {
			t.Errorf("test not found failed, expected %v, got %v", ErrRequestInvalid, err)
		} else {
			t.Logf("test not found success, expected %v, got %v", ErrRequestInvalid, err)
		}
	}
}

func TestRegister(t *testing.T) {
	r := &newRepo{}
	service := NewUserService(r)

	newUser := RegisterRequestModel{
		FullName:  "test 2",
		Email:     "test2@test.com",
		Password:  "bcrypt2",
	}

	invalidUser := RegisterRequestModel{
		FullName: "hallo",
	}
	
	user, err := service.Register(&newUser)

	if err != nil{
		t.Errorf("test user register failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test user register success, expected %v, got %v", user, user)
	}

	user, err = service.Register(&invalidUser)

	if err != nil{
		t.Logf("test user register invalid success, expected %v, got %v", ErrRequestInvalid, err)
	} else {
		t.Errorf("test user register invalid failed, expected %v, got %v", ErrRequestInvalid, err)
	}
}


func TestConfirmRegistration(t *testing.T) {
	r := &newRepo{}
	service := NewUserService(r)

	validUser := ConfirmRegistrationRequestModel{
		Email:     "test.user2@test.com",
		RegistrationCode: "111111",
	}

	invalidUser := ConfirmRegistrationRequestModel{
		Email: "hallo",
	}

	 err := service.ConfirmRegistration(&validUser)

	if err != nil{
		t.Errorf("test user confirm registration failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test user confirm registration success, expected %v, got %v", nil, nil)
	}

	 err = service.ConfirmRegistration(&invalidUser)

	if err != nil{
		t.Logf("test user confirm registration invalid success, expected %v, got %v", ErrRequestInvalid, err)
	} else {
		t.Errorf("test user confirm registration invalid failed, expected %v, got %v", ErrRequestInvalid, err)
	}
}


func TestGetAllUser(t *testing.T){
	r := &newRepo{}
	service := NewUserService(r)

	getAllUserModel := GetAllUserRequestModel{
		Page:        1,
		RowsPerPage: 10,
		SortBy:      "Id",
		Descending:  false,
		Filter:      "",
	}

	shouldFind, err := service.GetAllUser(&getAllUserModel)

	if err != nil {
		t.Errorf("test GetAllUser failed, expected %v, got %v", "GetAllUserResponseModel", err)
	}

	if shouldFind.rows[0].Id != "8a5e9658-f954-45c0-a232-4dcbca0d4907" {
		t.Errorf("test GetAllUser failed, expected %v, got %v", "8a5e9658-f954-45c0-a232-4dcbca0d4907", shouldFind.rows[0].Id)
	} else {
		t.Logf("test GetAllUser success, expected %v, got %v", "8a5e9658-f954-45c0-a232-4dcbca0d4907",  shouldFind.rows[0].Id)
	}

	if shouldFind.count > 0 {
		t.Logf("test GetALlUser success, expected %v, got %v", "greater than zero", shouldFind.count)
	} else {
		t.Errorf("test GetAllUser failed, expected %v, got %v", "greater than zero", shouldFind.count)
	}


	wrongModel := GetAllUserRequestModel{
		RowsPerPage: 0,
		SortBy:      "",
		Descending:  false,
		Filter:      "",
	}

	_, err = service.GetAllUser(&wrongModel)

	if err != nil {
		if errors.Cause(err) != ErrRequestInvalid {
			t.Errorf("test GetAllUser failed, expected %v, got %v", ErrRequestInvalid, err)
		} else {
			t.Logf("test GetAllUser success, expected %v, got %v", ErrRequestInvalid, err)
		}
	}
}

func TestLogin(t *testing.T){
	r := &newRepo{}
	service := NewUserService(r)

	loginModel := LoginRequestModel{
		Email:    "test.user@test.com",
		Password: "bcryptasdasd123",
	}

	loginResp, err := service.Login(&loginModel)

	if err != nil {
		t.Errorf("test Login failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test Login success, expected %v, got %v", nil, err)
	}

	if loginResp.Roles[0].Name != "user" {
		t.Errorf("test login failed, expected role %v, got %v","user", loginResp.Roles[0].Name)
	} else {
		t.Logf("test Login success, expected role %v, got %v", "user", loginResp.Roles[0].Name)
	}

	badLoginModel := LoginRequestModel{
		Email:    "asdsad",
	}

	loginResp, err = service.Login(&badLoginModel)

	if err != nil {
		if errors.Cause(err) == ErrRequestInvalid {
			t.Logf("test Login success, expected %v, got %v",ErrRequestInvalid, err)
		} else {
			t.Errorf("test Login failed, expected %v, got %v", ErrRequestInvalid, err)
		}
	}

	wrongPwModel := LoginRequestModel{
		Email:    "test.user@test.com",
		Password: "abaslé3",
	}

	loginResp, err = service.Login(&wrongPwModel)

	if err != nil {
		if errors.Cause(err) == ErrUserInvalid {
			t.Logf("test Login success, expected %v, got %v",ErrUserInvalid, err)
		} else {
			t.Errorf("test Login failed, expected %v, got %v", ErrUserInvalid, err)
		}
	}
}

func TestSendPasswdReset(t *testing.T){
	r := &newRepo{}
	service := NewUserService(r)

	newUser := SendPasswdResetRequestModel{
		Email:     "test.user3@test.com",
	}

	invalidUser := SendPasswdResetRequestModel{
		Email: "asd",
	}

	user, err := service.SendPasswdReset(&newUser)

	if err != nil{
		t.Errorf("test user SendPasswdReset failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test user SendPasswdReset success, expected %v, got %v", "SendPasswdResetResponseModel", user)
	}

	user, err = service.SendPasswdReset(&invalidUser)

	if err != nil{
		t.Logf("test user SendPasswdReset invalid success, expected %v, got %v", ErrRequestInvalid, err)
	} else {
		t.Errorf("test user SendPasswdReset invalid failed, expected %v, got %v", ErrRequestInvalid, err)
	}

}

func TestResetPassword(t *testing.T){
	r := &newRepo{}
	service := NewUserService(r)

	validUser := ResetPasswordRequestModel{
		Email:     "test.user@resetpw.com",
		ResetPasswordCode: "333333",
		Password: "asdasd1234",
	}

	invalidUser := ResetPasswordRequestModel{
		Email: "hallo",
	}

	err := service.ResetPassword(&validUser)

	if err != nil{
		t.Errorf("test user ResetPassword failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test user ResetPassword success, expected %v, got %v", nil, nil)
	}

	err = service.ResetPassword(&invalidUser)

	if err != nil{
		t.Logf("test user ResetPassword invalid success, expected %v, got %v", ErrRequestInvalid, err)
	} else {
		t.Errorf("test user ResetPassword invalid failed, expected %v, got %v", ErrRequestInvalid, err)
	}

}

func TestUpdateUser(t *testing.T) {
	r := &newRepo{}
	service := NewUserService(r)


	updateUser := UpdateUserRequestModel{
		Id: "8a5e9658-f954-45c0-a232-4dcbca0d4907",
		FullName:  "updated",
		Email:     "updated@test.com",
	}

	invalidUser := UpdateUserRequestModel{
		FullName: "hallo",
	}

	err := service.UpdateUser(&updateUser)

	if err != nil{
		t.Errorf("test user update failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test user update success, expected %v, got %v", nil, err)
	}

	gUserReqModel := GetUserRequestModel{
		Id: "8a5e9658-f954-45c0-a232-4dcbca0d4907",
	}

	shouldFind, err := service.GetUser(&gUserReqModel)

	if err != nil {
		t.Errorf("test user update failed, expected %v after update, got %v",updateUser, shouldFind)
	} else {
		if shouldFind.FullName != "updated" {
			t.Errorf("test user update failed, expected FullName = %v after update, got %v",updateUser.FullName, shouldFind.FullName)
		}
		if shouldFind.Email != "updated@test.com"{
			t.Errorf("test user update failed, expected Email = %v after update, got %v",updateUser.Email, shouldFind.Email)
		}
			t.Logf("test user update success, all values are as expected")
	}

	err = service.UpdateUser(&invalidUser)

	if err != nil{
		t.Logf("test user update invalid success, expected %v, got %v", ErrUserInvalid, err)
	} else {
		t.Errorf("test user update invalid failed, expected %v, got %v", ErrUserInvalid, err)
	}
}

func TestDeleteUser(t *testing.T){
	r := &newRepo{}
	service := NewUserService(r)

	deleteUserModel := DeleteUserRequestModel{Id:"8a5e9658-f954-45c0-a232-4dcbca0d4907" }

	err := service.DeleteUser(&deleteUserModel)

	if err != nil {
		t.Errorf("test delete failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test delete success, expected %v, got %v", nil, err)
	}


	invalidModel := DeleteUserRequestModel{Id:"abc" }

	err = service.DeleteUser(&invalidModel)

	if err != nil {
		if errors.Cause(err) != ErrRequestInvalid {
			t.Errorf("test delete failed, expected %v, got %v", ErrRequestInvalid, err)
		} else {
			t.Logf("test delete success, expected %v, got %v", ErrRequestInvalid, err)
		}
	}

}