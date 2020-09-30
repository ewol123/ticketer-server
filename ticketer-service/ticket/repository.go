package ticket

// TicketRepository: interface to connect our business logic to our repository
type Repository interface {
	Find(column string, value string) (*Ticket, error)
	FindAll(page int, rowsPerPage int, sortBy string, descending bool, filter string, workerId string, status string, lat string, long string) (*[]Ticket, int, error)
	Store(ticket *Ticket) (*Ticket, error)
	Update(ticket *Ticket) error
	Delete(id string) error
}