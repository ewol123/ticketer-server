package ticket

import (
	"time"
)

// FaultType : represents allowed fault types
type FaultType string
const (
	LEAK FaultType = "leak"
)

// StatusType : represents allowed ticket statuses
type StatusType string
const (
	DRAFT StatusType = "draft"
	DONE StatusType = "done"
	INACTIVE StatusType = "inactive"
)


// Ticket : ticket model
type Ticket struct {
	Id        string    `migrations:"id" json:"id" db:"ticket_id" `
	CreatedAt time.Time  `migrations:"created_at" json:"created_at" db:"ticket_created_at"`
	UpdatedAt time.Time `migrations:"updated_at" json:"updated_at" db:"ticket_updated_at"`
	UserId 	  string `migrations:"user_id" json:"user_id" db:"ticket_user_id"`
	WorkerId string `migrations:"worker_id" json:"worker_id" db:"ticket_worker_id"`
	FaultType FaultType `migrations:"fault_type" json:"fault_type" db:"ticket_fault_type"`
	Address string `migrations:"address" json:"address" db:"ticket_address"`
	FullName string `migrations:"full_name" json:"full_name" db:"ticket_full_name"`
	Phone string `migrations:"phone" json:"phone" db:"ticket_phone"`
	GeoLocation string `migrations:"geo_location" json:"geo_location" db:"ticket_geo_location"`
	ImageUrl string `migrations:"image_url" json:"image_url" db:"ticket_image_url"`
	Status StatusType `migrations:"status" json:"status" db:"ticket_status"`
}