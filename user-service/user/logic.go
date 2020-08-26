package user

import (
	"crypto/rand"
	"errors"
	"github.com/google/uuid"
	errs "github.com/pkg/errors"
	validate "gopkg.in/dealancer/validate.v2"
	"time"
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

func (u *userService) FindAll(page int, rowsPerPage int, sortBy string, descending bool, filter string) (*[]User, int, error) {
	return u.userRepo.FindAll(page,rowsPerPage,sortBy,descending, filter)
}

func (u *userService) Store(user *User) (*User, error) {
	if err := validate.Validate(user); err != nil {
		return nil, errs.Wrap(ErrUserInvalid, "service.User.Store")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Status = PENDING

	regCode, err := genCode(6)
	if err != nil {
		return nil, errs.Wrap(err, "service.User.Store")
	}
	user.RegistrationCode = regCode

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, errs.Wrap(err, "service.User.Store")
	}
	user.Id = id.String()
	return u.userRepo.Store(user)

}

func (u *userService) Update(user *User) error {
	if err := validate.Validate(user); err != nil {
		return errs.Wrap(ErrUserInvalid, "service.User.Store")
	}

	user.UpdatedAt = time.Now()
	return u.userRepo.Update(user)
}

func (u *userService) Delete(id string) error {
	return u.userRepo.Delete(id)
}


// UTILITY
const otpChars = "1234567890"
func genCode(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}


