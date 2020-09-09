package user

import (
	"crypto/rand"
	"errors"
	"github.com/google/uuid"
	errs "github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	validate "gopkg.in/dealancer/validate.v2"
	"time"
)


// errors
var (
	ErrUserNotFound = errors.New("user Not Found")
	ErrUserInvalid  = errors.New("user Invalid")
	ErrRequestInvalid = errors.New("request payload is invalid")
)


type userService struct {
	userRepo Repository
}



// NewUserService : create a new user service
func NewUserService(userRepo Repository) Service {
	return &userService{
		userRepo,
	}
}


func (u *userService) GetUser(model *GetUserRequestModel) (*GetUserResponseModel, error) {
	if err := validate.Validate(model); err != nil {
		return nil, errs.Wrap(ErrRequestInvalid, "service.User.GetUser")
	}

	user, err := u.userRepo.Find("id",model.Id)
	if err != nil {
		return nil, errs.Wrap(err, "service.User.GetUser")
	}

	getUserResponseModel := GetUserResponseModel{
		Id:        user.Id,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		FullName:  user.FullName,
		Email:     user.Email,
		Status:    user.Status,
		RegistrationCode: user.RegistrationCode,
		ResetPasswordCode: user.ResetPasswordCode,
		Roles:     user.Roles,
	}

	return &getUserResponseModel, nil

}


func (u *userService) GetAllUser(model *GetAllUserRequestModel) (*GetAllUserResponseModel, error) {
	if err := validate.Validate(model); err != nil {
		return nil, errs.Wrap(ErrRequestInvalid, "service.User.GetAllUser")
	}

	page := model.Page
	rowsPerPage := model.RowsPerPage
	sortBy := model.SortBy
	descending := model.Descending
	filter := model.Filter

	users,count,err := u.userRepo.FindAll(page,rowsPerPage,sortBy,descending,filter)

	if err != nil {
		return nil, errs.Wrap(err, "service.User.GetAllUser")
	}

	var getUserResponseModels []GetUserResponseModel

	if len(*users) > 0 {
		for _,user := range *users {
			gUserModel := GetUserResponseModel{
				Id:       user.Id,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				FullName:  user.FullName,
				Email:     user.Email,
				Status:    user.Status,
				Roles:     user.Roles,
			}
			getUserResponseModels = append(getUserResponseModels,gUserModel)
		}
	}

	getAllUserResponseModel := GetAllUserResponseModel{
		Count: count,
		Rows:  getUserResponseModels,
	}

	return &getAllUserResponseModel, nil

}


func (u *userService) Register(model *RegisterRequestModel) (*RegisterReturnModel,error) {
	if err := validate.Validate(model); err != nil {
		return nil, errs.Wrap(ErrRequestInvalid, "service.User.Register")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errs.Wrap(err, "service.User.Register")
	}

	regCode, err := genCode(6)
	if err != nil {
		return  nil, errs.Wrap(err, "service.User.Register")
	}

	password := []byte(model.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.Wrap(ErrUserInvalid, "service.User.Register")
	}

	user := User{
		Id:                id.String(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		FullName:          model.FullName,
		Email:             model.Email,
		Password:          string(hashedPassword),
		Status:            PENDING,
		RegistrationCode:  regCode,
		Roles:             []Role{USER},
	}

	_, err = u.userRepo.Store(&user)

	regReturnModel := RegisterReturnModel{
		Id: id.String(),
		RegistrationCode: regCode,
	}

	return &regReturnModel, nil
}

func (u *userService) ConfirmRegistration(model *ConfirmRegistrationRequestModel) error {
	if err := validate.Validate(model); err != nil {
		return errs.Wrap(ErrRequestInvalid, "service.User.ConfirmRegistration")
	}

	user, err := u.userRepo.Find("email",model.Email)
	if err != nil {
		return  errs.Wrap(err, "service.User.ConfirmRegistration")
	}

	if user.RegistrationCode != model.RegistrationCode {
		return errs.Wrap(ErrUserNotFound, "service.User.ConfirmRegistration")
	}

	updateUser := User{
		Id: user.Id,
		UpdatedAt: time.Now(),
		Status: ACTIVE,
		RegistrationCode: "",
	}

	return u.userRepo.Update(&updateUser)

}

func (u *userService) Login(model *LoginRequestModel) (*LoginResponseModel, error) {
	if err := validate.Validate(model); err != nil {
		return nil, errs.Wrap(ErrRequestInvalid, "service.User.Login")
	}

	user, err := u.userRepo.Find("email",model.Email)

	if err != nil {
		return  nil, errs.Wrap(ErrUserNotFound, "service.User.Login")
	}

	if user.Status != ACTIVE {
		return nil, errs.Wrap(ErrUserInvalid, "service.User.Login")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(model.Password))
	if err != nil {
		return nil, errs.Wrap(ErrUserInvalid, "service.User.Login")
	}

	respModel := LoginResponseModel{Roles: user.Roles}

	return &respModel, nil

}

func (u *userService) SendPasswdReset(model *SendPasswdResetRequestModel) (*SendPasswdResetResponseModel,error) {
	if err := validate.Validate(model); err != nil {
		return nil,errs.Wrap(ErrRequestInvalid, "service.User.SendPasswdReset")
	}

	user, err := u.userRepo.Find("email",model.Email)

	if err != nil {
		return nil, errs.Wrap(ErrUserNotFound, "service.User.SendPasswdReset")
	}

	if user.Status != ACTIVE {
		return nil, errs.Wrap(ErrUserInvalid, "service.User.SendPasswdReset")
	}

	regCode, err := genCode(6)
	if err != nil {
		return  nil, errs.Wrap(err, "service.User.SendPasswdReset")
	}

	updateUser := User{
		Id: user.Id,
		UpdatedAt:         time.Now(),
		ResetPasswordCode: regCode,
	}

	err = u.userRepo.Update(&updateUser)

	if err != nil {
		return nil, errs.Wrap(err, "service.User.SendPasswdReset")
	}

	pwdResetResponseModel := SendPasswdResetResponseModel{
		ResetPasswordCode: regCode,
		Email: user.Email,
		FullName: user.FullName}

	return &pwdResetResponseModel, nil
}


func (u *userService) ResetPassword(model *ResetPasswordRequestModel) error {
	if err := validate.Validate(model); err != nil {
		return errs.Wrap(ErrRequestInvalid, "service.User.ResetPassword")
	}

	user, err := u.userRepo.Find("email",model.Email)
	if err != nil {
		return  errs.Wrap(err, "service.User.ResetPassword")
	}

	if user.ResetPasswordCode != model.ResetPasswordCode {
		return errs.Wrap(ErrUserNotFound, "service.User.ResetPassword")
	}

	password := []byte(model.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return errs.Wrap(ErrUserInvalid, "service.User.ResetPassword")
	}

	updateUser := User{
		Id: user.Id,
		UpdatedAt: time.Now(),
		ResetPasswordCode: "",
		Password: string(hashedPassword),
	}

	return u.userRepo.Update(&updateUser)

}

func (u *userService) UpdateUser(model *UpdateUserRequestModel) error {
	if err := validate.Validate(model); err != nil {
		return errs.Wrap(ErrRequestInvalid, "service.User.UpdateUser")
	}

	user := User{
		Id: 			  model.Id,
		UpdatedAt:        time.Now(),
		FullName:         model.FullName,
		Email:            model.Email,
		Status:           model.Status,
	}

	return u.userRepo.Update(&user)
}


func (u *userService) DeleteUser(model *DeleteUserRequestModel) error {
	if err := validate.Validate(model); err != nil {
		return errs.Wrap(ErrRequestInvalid, "service.User.DeleteUser")
	}
	return u.userRepo.Delete(model.Id)
}


// UTILITY
const codeCars = "1234567890"
func genCode(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(codeCars)
	for i := 0; i < length; i++ {
		buffer[i] = codeCars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}


