package user

// UserService : main service interface for user
type Service interface {
	Find(id string) (*User, error)
	FindAll(page int, rowsPerPage int, sortBy string, descending bool) (*[]User, int, error)
	Store(user *User) error
	Update(user *User) (*User, error)
	Delete(id string) error
}
