package user

// UserRepository: interface to connect our business logic to our repository
type Repository interface {
	Find(id int) (*User, error)
	Store(user *User) error
}
