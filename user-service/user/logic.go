package user

import (
	"errors"
	"github.com/google/uuid"
	validate "gopkg.in/dealancer/validate.v2"
	"time"

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
	return u.userRepo.Find(id)
}

func (u *userService) FindAll(page int, rowsPerPage int, sortBy string, descending bool) (*[]User, int, error) {
	return u.userRepo.FindAll(page,rowsPerPage,sortBy,descending)
}

func (u *userService) Store(user *User) error {
	if err := validate.Validate(user); err != nil {
		return errs.Wrap(ErrUserInvalid, "service.User.Store")
	}
	user.CreatedAt = time.Now().UTC().Unix()
	user.UpdatedAt = time.Now().UTC().Unix()
	id, err := uuid.NewRandom()
	if err != nil {
		return errs.Wrap(err, "service.User.Store")
	}
	user.Id = id.String()
	return u.userRepo.Store(user)

}

func (u *userService) Update(user *User) error {
	if err := validate.Validate(user); err != nil {
		return errs.Wrap(ErrUserInvalid, "service.User.Store")
	}

	user.UpdatedAt = time.Now().UTC().Unix()
	return u.userRepo.Update(user)
}

func (u *userService) Delete(id string) error {
	return u.userRepo.Delete(id)
}
