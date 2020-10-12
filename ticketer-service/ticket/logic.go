package ticket

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	errs "github.com/pkg/errors"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrTicketNotFound = errors.New("ticket Not Found")
	ErrTicketInvalid  = errors.New("ticket Invalid")
	ErrRequestInvalid = errors.New("request payload is invalid")
)

type ticketService struct {
	ticketRepo Repository
}

func (t ticketService) SyncTicketWorker(model *SyncTicketRequestModelWorker) (*SyncTicketResponseModelWorker, error) {
	if err := validate.Validate(model); err != nil {
		return nil, errs.Wrap(ErrRequestInvalid, "service.Ticket.SyncTicketWorker")
	}

	workerId := model.RequesterId
	status := "draft"
	lat := model.Lat
	long := model.Long

	// TICKET VALIDATION AND UPDATE
	// ///////////////////////////

	for _, ticket := range model.Rows {
		checkTicket, err := t.ticketRepo.Find("id", ticket.Id)

		if err != nil {
			return nil, errs.Wrap(err, "service.Ticket.SyncTicketWorker")
		}

		//validation constraints which need to be true in order to update the ticket
		isWorkerValid := checkTicket.WorkerId == model.RequesterId
		isStatusValid := checkTicket.Status == DRAFT
		// TODO: [ticket sync] put more validation constraints after this line if needed

		if isWorkerValid && isStatusValid {
			// The client will send a base64 encoded string as image data, the ImageUrl here is the already uploaded
			// image based on that data. The business logic does not care how it gets the url, uploading the image to
			// a cloud provider / file server / image service is the responsibility of whatever implements the logic
			updateForm := Ticket{
				Id:       ticket.Id,
				ImageUrl: ticket.ImageUrl,
				Status:   ticket.Status,
			}

			err = t.ticketRepo.Update(&updateForm)

			if err != nil {
				return nil, errs.Wrap(err, "service.Ticket.SyncTicketWorker")
			}
		}
	}

	// GET THE FRESH STATE
	// ///////////////////////////

	// pagination with 10 results / page.
	// since a worker won't have more than 10 tickets at a time it's enough
	page := 1
	rowsPerPage := 10 // TODO: change this if more than 10 tickets need to be assigned at the same time
	sortBy := "ticket_created_at"
	descending := true
	filter := ""

	//we only want to return tickets
	// - assigned to the worker
	// - with draft status
	tickets, _, err := t.ticketRepo.FindAll(page, rowsPerPage, sortBy, descending, filter, workerId, status, "", "")

	if err != nil {
		return nil, errs.Wrap(err, "service.Ticket.SyncTicketWorker")
	}

	var syncTicketResponseModels SyncTicketResponseModelWorker

	if len(*tickets) > 0 {
		for _, ticket := range *tickets {
			gTicketModel := SyncTicketResp{
				Id:              ticket.Id,
				FaultType:       ticket.FaultType,
				Address:         ticket.Address,
				FullName:        ticket.FullName,
				Phone:           ticket.Phone,
				GeoLocation:     ticket.GeoLocation,
				Status:          ticket.Status,
				ServerUpdatedAt: ticket.UpdatedAt,
			}
			syncTicketResponseModels.Rows = append(syncTicketResponseModels.Rows, gTicketModel)
		}
	}

	// ASSIGN NEW TICKETS IF NEEDED
	// ///////////////////////////

	if len(syncTicketResponseModels.Rows) < 10 { // TODO: change this if more than 10 tickets need to be assigned at the same time

		// get tickets where
		// ticket is not assigned to anyone (workerid is null)
		// ticket is draft
		// worker is within X (10) km
		page = 1
		rowsPerPage = 10 - len(syncTicketResponseModels.Rows) // TODO: change this if more than 10 tickets need to be assigned at the same time
		sortBy = "ticket_created_at"
		descending = false // descending is false because older tickets will have higher priority to complete
		filter = ""

		newTickets, _, err := t.ticketRepo.FindAll(page, rowsPerPage, sortBy, descending, filter, "NULL", status, lat, long)

		if err != nil {
			return nil, errs.Wrap(err, "service.Ticket.SyncTicketWorker")
		}

		if len(*newTickets) > 0 {
			for _, newTicket := range *newTickets {
				gTicketModel := SyncTicketResp{
					Id:              newTicket.Id,
					FaultType:       newTicket.FaultType,
					Address:         newTicket.Address,
					FullName:        newTicket.FullName,
					Phone:           newTicket.Phone,
					GeoLocation:     newTicket.GeoLocation,
					Status:          newTicket.Status,
					ServerUpdatedAt: newTicket.UpdatedAt,
				}
				syncTicketResponseModels.Rows = append(syncTicketResponseModels.Rows, gTicketModel)

				// assign new ticket to the worker in the database too
				// usually we should only have one WRITE operation so everything is atomic
				// but because of the way the app is designed it won't cause any inconsistency
				updateForm := Ticket{
					Id:       newTicket.Id,
					WorkerId: workerId,
				}
				err = t.ticketRepo.Update(&updateForm)
				if err != nil {
					return nil, errs.Wrap(err, "service.Ticket.SyncTicketWorker")
				}

			}
		}
	}

	return &syncTicketResponseModels, nil

}

func (t ticketService) CreateTicketUser(model *CreateTicketRequestModelUser) error {
	if err := validate.Validate(model); err != nil {
		return errs.Wrap(ErrRequestInvalid, "service.Ticket.CreateTicketUser")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return errs.Wrap(err, "service.Ticket.CreateTicketUser")
	}

	geoLocation := fmt.Sprintf(`%v,%v`, model.Lat, model.Long)

	ticket := Ticket{
		Id:          id.String(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserId:      model.RequesterId,
		WorkerId:    "",
		FaultType:   model.FaultType,
		Address:     model.Address,
		FullName:    model.FullName,
		Phone:       model.Phone,
		GeoLocation: geoLocation,
		Status:      DRAFT,
	}

	_, err = t.ticketRepo.Store(&ticket)

	if err != nil {
		return errs.Wrap(err, "service.Ticket.CreateTicketUser")
	}

	return nil

}

func (t ticketService) GetTicketWorker(model *GetTicketRequestModelWorker) (*GetTicketResponseModel, error) {
	if err := validate.Validate(model); err != nil {
		return nil, errs.Wrap(ErrRequestInvalid, "service.Ticket.GetTicketWorker")
	}

	ticket, err := t.ticketRepo.Find("id", model.Id)
	if err != nil {
		return nil, errs.Wrap(err, "service.Ticket.GetTicketWorker")
	}

	if ticket.WorkerId != model.RequesterId {
		return nil, errs.Wrap(ErrRequestInvalid, "service.Ticket.GetTicketWorker")
	}

	getTicketResponseModel := GetTicketResponseModel{
		Id:          ticket.Id,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
		UserId:      ticket.UserId,
		WorkerId:    ticket.WorkerId,
		FaultType:   ticket.FaultType,
		Address:     ticket.Address,
		FullName:    ticket.FullName,
		Phone:       ticket.Phone,
		GeoLocation: ticket.GeoLocation,
		ImageUrl:    ticket.ImageUrl,
		Status:      ticket.Status,
	}

	return &getTicketResponseModel, nil
}

func (t ticketService) GetAllTicketWorker(model *GetAllTicketRequestModelWorker) (*GetAllTicketResponseModel, error) {
	if err := validate.Validate(model); err != nil {
		return nil, errs.Wrap(ErrRequestInvalid, "service.Ticket.GetAllTicketWorker")
	}

	page := model.Page
	rowsPerPage := model.RowsPerPage
	sortBy := model.SortBy
	descending := model.Descending
	filter := model.Filter
	workerId := model.RequesterId
	status := "draft"

	//we only want to return tickets assigned to the worker with draft status
	tickets, count, err := t.ticketRepo.FindAll(page, rowsPerPage, sortBy, descending, filter, workerId, status, "", "")

	if err != nil {
		return nil, errs.Wrap(err, "service.Ticket.GetAllTicketWorker")
	}

	var getTicketResponseModels []GetTicketResponseModel

	if len(*tickets) > 0 {
		for _, ticket := range *tickets {
			gTicketModel := GetTicketResponseModel{
				Id:          ticket.Id,
				CreatedAt:   ticket.CreatedAt,
				UpdatedAt:   ticket.UpdatedAt,
				UserId:      ticket.UserId,
				WorkerId:    ticket.WorkerId,
				FaultType:   ticket.FaultType,
				Address:     ticket.Address,
				FullName:    ticket.FullName,
				Phone:       ticket.Phone,
				GeoLocation: ticket.GeoLocation,
				ImageUrl:    ticket.ImageUrl,
				Status:      ticket.Status,
			}
			getTicketResponseModels = append(getTicketResponseModels, gTicketModel)
		}
	}

	getAllTicketResponseModel := GetAllTicketResponseModel{
		Count: count,
		Rows:  getTicketResponseModels,
	}

	return &getAllTicketResponseModel, nil
}

func (t ticketService) GetTicketAdmin(model *GetTicketRequestModelAdmin) (*GetTicketResponseModel, error) {
	if err := validate.Validate(model); err != nil {
		return nil, errs.Wrap(ErrRequestInvalid, "service.Ticket.GetTicketAdmin")
	}

	ticket, err := t.ticketRepo.Find("id", model.Id)
	if err != nil {
		return nil, errs.Wrap(err, "service.Ticket.GetTicketAdmin")
	}

	getTicketResponseModel := GetTicketResponseModel{
		Id:          ticket.Id,
		CreatedAt:   ticket.CreatedAt,
		UpdatedAt:   ticket.UpdatedAt,
		UserId:      ticket.UserId,
		WorkerId:    ticket.WorkerId,
		FaultType:   ticket.FaultType,
		Address:     ticket.Address,
		FullName:    ticket.FullName,
		Phone:       ticket.Phone,
		GeoLocation: ticket.GeoLocation,
		ImageUrl:    ticket.ImageUrl,
		Status:      ticket.Status,
	}

	return &getTicketResponseModel, nil
}

func (t ticketService) GetAllTicketAdmin(model *GetAllTicketRequestModelAdmin) (*GetAllTicketResponseModel, error) {
	if err := validate.Validate(model); err != nil {
		return nil, errs.Wrap(ErrRequestInvalid, "service.Ticket.GetAllTicketAdmin")
	}

	page := model.Page
	rowsPerPage := model.RowsPerPage
	sortBy := model.SortBy
	descending := model.Descending
	filter := model.Filter

	tickets, count, err := t.ticketRepo.FindAll(page, rowsPerPage, sortBy, descending, filter, "", "", "", "")

	if err != nil {
		return nil, errs.Wrap(err, "service.Ticket.GetAllTicketAdmin")
	}

	var getTicketResponseModels []GetTicketResponseModel

	if len(*tickets) > 0 {
		for _, ticket := range *tickets {
			gTicketModel := GetTicketResponseModel{
				Id:          ticket.Id,
				CreatedAt:   ticket.CreatedAt,
				UpdatedAt:   ticket.UpdatedAt,
				UserId:      ticket.UserId,
				WorkerId:    ticket.WorkerId,
				FaultType:   ticket.FaultType,
				Address:     ticket.Address,
				FullName:    ticket.FullName,
				Phone:       ticket.Phone,
				GeoLocation: ticket.GeoLocation,
				ImageUrl:    ticket.ImageUrl,
				Status:      ticket.Status,
			}
			getTicketResponseModels = append(getTicketResponseModels, gTicketModel)
		}
	}

	getAllTicketResponseModel := GetAllTicketResponseModel{
		Count: count,
		Rows:  getTicketResponseModels,
	}

	return &getAllTicketResponseModel, nil
}

func (t ticketService) UpdateTicketAdmin(model *UpdateTicketRequestModelAdmin) error {
	if err := validate.Validate(model); err != nil {
		return errs.Wrap(ErrRequestInvalid, "service.Ticket.UpdateTicketAdmin")
	}

	ticketUpdate := Ticket{
		Id:        model.Id,
		WorkerId:  model.WorkerId,
		FaultType: model.FaultType,
		Address:   model.Address,
		FullName:  model.FullName,
		Phone:     model.Phone,
		ImageUrl:  model.ImageUrl,
		Status:    model.Status,
	}

	if model.Lat != "" && model.Long != "" {
		ticketUpdate.GeoLocation = fmt.Sprintf(`%v,%v`, model.Lat, model.Long)
	}

	return t.ticketRepo.Update(&ticketUpdate)
}

func (t ticketService) DeleteTicketAdmin(model *DeleteTicketRequestModelAdmin) error {
	if err := validate.Validate(model); err != nil {
		return errs.Wrap(ErrRequestInvalid, "service.Ticket.DeleteTicketAdmin")
	}
	return t.ticketRepo.Delete(model.Id)
}

// NewTicketService : create a new ticket service
func NewTicketService(ticketRepo Repository) Service {
	return &ticketService{
		ticketRepo,
	}
}
