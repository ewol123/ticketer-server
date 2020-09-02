package user

import (
	"time"
)


// StatusType : represents allowed user statuses
type StatusType string
const (
	PENDING StatusType = "pending"
	ACTIVE StatusType = "active"
	INACTIVE StatusType = "inactive"
)

var (
	USER = Role{Id: "72daf87a-fda4-4c72-aff9-85edd68d155f", Name: "user"}
	ADMIN  = Role{Id: "336a3ff6-9fdb-496f-ac8c-e37759969cf2", Name: "admin"}
	)



// User : user model
type User struct {
	Id        string    `migrations:"id" json:"id" db:"user_id" `
	CreatedAt time.Time  `migrations:"created_at" json:"created_at" db:"user_created_at"`
	UpdatedAt time.Time `migrations:"updated_at" json:"updated_at" db:"user_updated_at"`
	FullName  string `migrations:"full_name" json:"full_name" validate:"empty=false" db:"user_full_name"`
	Email     string `migrations:"email" json:"email" validate:"empty=false & format=email" db:"user_email"`
	Password  string `migrations:"password" json:"password" validate:"empty=false" db:"user_password"`
	Status	  StatusType `migrations:"status" json:"status" db:"user_status"`
	RegistrationCode string `migrations:"registration_code" json:"registration_code" db:"user_registration_code"`
	ResetPasswordCode string  `migrations:"reset_password_code" json:"reset_password_code" db:"user_reset_password_code"`
	Roles     []Role `json:"roles"`
}
// Role : role model
type Role struct {
	Id   string    `json:"id" migrations:"id" db:"roles_id"`
	Name string `json:"name" migrations:"name" db:"roles_name"`
}
