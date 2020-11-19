package ticket

import "time"

// Service : main service interface for ticket. It has different requirements for ADMIN,WORKER,USER and UNSCOPED roles
type Service interface {
	// USER funcs
	CreateTicketUser(model *CreateTicketRequestModelUser) error

	// WORKER funcs
	GetTicketWorker(model *GetTicketRequestModelWorker) (*GetTicketResponseModel, error)
	GetAllTicketWorker(model *GetAllTicketRequestModelWorker) (*GetAllTicketResponseModel, error)
	SyncTicketWorker(model *SyncTicketRequestModelWorker) (*SyncTicketResponseModelWorker, error)

	// ADMIN funcs
	GetTicketAdmin(model *GetTicketRequestModelAdmin) (*GetTicketResponseModel, error)
	GetAllTicketAdmin(model *GetAllTicketRequestModelAdmin) (*GetAllTicketResponseModel, error)
	UpdateTicketAdmin(model *UpdateTicketRequestModelAdmin) error
	DeleteTicketAdmin(model *DeleteTicketRequestModelAdmin) error
}

// CREATE TICKET MODELS
// ///////////////////////////

type CreateTicketRequestModelUser struct {
	RequesterId string    `validate:"empty=false & format=uuid4"`
	FaultType   FaultType `validate:"empty=false"`
	Address     string    `validate:"empty=false"`
	FullName    string    `validate:"empty=false"`
	Phone       string    `validate:"empty=false"`
	Lat         string    `validate:"empty=false"`
	Long        string    `validate:"empty=false"`
}

// GET TICKET MODELS
// ///////////////////////////

type GetTicketRequestModelAdmin struct {
	Id string `validate:"empty=false & format=uuid4"`
}
type GetTicketRequestModelWorker struct {
	Id          string `validate:"empty=false & format=uuid4"`
	RequesterId string `validate:"empty=false & format=uuid4"` // who requests the resource? should be the worker's id
}
type GetTicketResponseModel struct {
	Id          string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserId      string
	WorkerId    string
	FaultType   FaultType
	Address     string
	FullName    string
	Phone       string
	GeoLocation string
	ImageUrl    string
	Status      StatusType
}

// GET ALL TICKET MODELS
// ///////////////////////////

type GetAllTicketRequestModelAdmin struct {
	Page        int    `validate:"gte=0"`
	RowsPerPage int    `validate:"gte=0"`
	SortBy      string `validate:"empty=false"`
	Descending  bool
	Filter      string
}

/* one worker probably will never have so many tickets at once that we need a pagination for them but
we can never know so just implement it */

type GetAllTicketRequestModelWorker struct {
	Page        int    `validate:"gte=0"`
	RowsPerPage int    `validate:"gte=0"`
	SortBy      string `validate:"empty=false"`
	Descending  bool
	Filter      string
	RequesterId string `validate:"empty=false & format=uuid4"` // who requests the resource? should be the worker's id
}

type GetAllTicketResponseModel struct {
	Count int
	Rows  []GetTicketResponseModel
}

// UPDATE TICKET MODELS
// ///////////////////////////

type UpdateTicketRequestModelAdmin struct {
	Id        string `validate:"empty=false & format=uuid4"`
	WorkerId  string
	FaultType FaultType
	Address   string
	FullName  string
	Phone     string
	Lat       string
	Long      string
	ImageUrl  string
	Status    StatusType
}

// DELETE TICKET MODELS
// ///////////////////////////

type DeleteTicketRequestModelAdmin struct {
	Id string `validate:"empty=false & format=uuid4"`
}

// SYNC TICKET MODELS

type SyncTicket struct {
	Id       string     `validate:"empty=false & format=uuid4"`
	ImageUrl string     `validate:"empty=false"`
	Status   StatusType `validate:"empty=false"`
}

type SyncTicketRequestModelWorker struct {
	RequesterId string `validate:"empty=false & format=uuid4"` // who requests the resource? should be the worker's id
	Lat         string `validate:"empty=false"`                // where is the worker when requesting the resource?
	Long        string `validate:"empty=false"`
	Rows        []SyncTicket
}

type SyncTicketResp struct {
	Id              string
	FaultType       FaultType
	Address         string
	FullName        string
	Phone           string
	GeoLocation     string
	Status          StatusType
	ServerUpdatedAt time.Time
}

type SyncTicketResponseModelWorker struct {
	Rows []SyncTicketResp
}
