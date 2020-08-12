package user

// UserRepository: interface to connect our business logic to our repository
type Repository interface {
	Find(id string) (*User, error)
	FindAll(page int, rowsPerPage int, sortBy string, descending bool) (*[]User, int, error)
	Store(user *User) error
	Update(user *User) (*User, error)
	Delete(id string) error
}
