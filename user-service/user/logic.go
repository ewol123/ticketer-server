package user

import (
	"errors"
	"time"

	validate "gopkg.in/dealancer/validate.v2"

	errs "github.com/pkg/errors"
)

// errors
var (
	ErrUserNotFound = errors.New("user Not Found")
	ErrUserInvalid  = errors.New("user Invalid")
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

func (u *userService) Find(id string) (*User, error) {
	return u.Find(id)
}

func (u *userService) Store(user *User) error {
	if err := validate.Validate(user); err != nil {
		return errs.Wrap(ErrUserInvalid, "service.User.Store")
	}
	user.CreatedAt = time.Now().UTC().Unix()
	user.UpdatedAt = time.Now().UTC().Unix()

	return u.userRepo.Store(user)

}
