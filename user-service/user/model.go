package user

// User : user model
type User struct {
	Id        string    `db:"id" json:"id"`
	CreatedAt int64  `db:"created_at" json:"created_at"`
	UpdatedAt int64  `db:"updated_at" json:"updated_at"`
	FullName  string `db:"full_name" json:"full_name" validate:"empty=false"`
	Email     string `db:"email" json:"email" validate:"empty=false & format=email"`
	Password  string `db:"password" json:"password" validate:"empty=false"`
	Roles     []Role `json:"roles" validate:"empty=false"`
}

type Role struct {
	Id   string    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
