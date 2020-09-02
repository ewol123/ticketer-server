package user

// UserRepository: interface to connect our business logic to our repository
type Repository interface {
	Find(column string, value string) (*User, error)
	FindAll(page int, rowsPerPage int, sortBy string, descending bool, filter string) (*[]User, int, error)
	Store(user *User) (*User, error)
	Update(user *User) error
	Delete(id string) error
}