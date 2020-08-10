package user

// UserService : main service interface for user
type Service interface {
	Find(id int) (*User, error)
	Store(user *User) error
}
