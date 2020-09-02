package user

import "time"

// UserService : main service interface for user
type Service interface {

	GetUser(model *GetUserRequestModel) (*GetUserResponseModel, error)
	GetAllUser(model *GetAllUserRequestModel) (*GetAllUserResponseModel, error)
	UpdateUser(model *UpdateUserRequestModel) error
	DeleteUser(model *DeleteUserRequestModel) error
	Register(model *RegisterRequestModel) (*RegisterReturnModel,error)
	ConfirmRegistration(model *ConfirmRegistrationRequestModel) error
	Login(model *LoginRequestModel) (*LoginResponseModel, error)
	ResetPassword(model *ResetPasswordRequestModel) error
	SendPasswdReset(model *SendPasswdResetRequestModel) (*SendPasswdResetResponseModel, error)

}


type GetUserRequestModel struct {
	Id string `validate:"empty=false & format=uuid4"`
}

type GetUserResponseModel struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time
	FullName  string
	Email     string
	Status	  StatusType
	RegistrationCode string
	ResetPasswordCode string
	Roles     []Role
}

type GetAllUserRequestModel struct{
	Page int `validate:"gte=0"`
	RowsPerPage int `validate:"gte=0"`
	SortBy string `validate:"empty=false"`
	Descending bool
	Filter string
}
type GetAllUserResponseModel struct {
	count int
	rows []GetUserResponseModel
}

type UpdateUserRequestModel struct {
	Id 		  string `validate:"empty=false"`
	FullName  string `validate:"empty=false"`
	Email     string `validate:"empty=false & format=email"`
	Status	  StatusType
}

type RegisterRequestModel struct {
	FullName string `validate:"empty=false"`
	Email    string `validate:"empty=false & format=email"`
	Password string `validate:"empty=false"`
}
type RegisterReturnModel struct {
	Id 				 string
	RegistrationCode string
}

type DeleteUserRequestModel struct {
	Id string `validate:"empty=false & format=uuid4"`
}
type ConfirmRegistrationRequestModel struct{
	Email string `validate:"empty=false & format=email"`
	RegistrationCode string  `validate:"empty=false"`
}
type LoginRequestModel struct{
	Email string `validate:"empty=false & format=email"`
	Password string `validate:"empty=false"`
}
type LoginResponseModel struct{
	Roles []Role
}

type ResetPasswordRequestModel struct{
	Email string `validate:"empty=false & format=email"`
	ResetPasswordCode string `validate:"empty=false"`
	Password string `validate:"empty=false"`
}
type SendPasswdResetRequestModel struct{
	Email string `validate:"empty=false & format=email"`
}
type SendPasswdResetResponseModel struct {
	ResetPasswordCode string
	Email string
	FullName string
}