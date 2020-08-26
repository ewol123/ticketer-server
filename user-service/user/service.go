package user

// UserService : main service interface for user
type Service interface {
	Find(id string) (*User, error)
	FindAll(page int, rowsPerPage int, sortBy string, descending bool, filter string) (*[]User, int, error)
	Store(user *User) (*User, error)
	Update(user *User) error
	Delete(id string) error
}
