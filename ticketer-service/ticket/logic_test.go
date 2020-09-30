package ticket

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

var ticket = Ticket{
	Id:          "8a5e9658-f954-45c0-a232-4dcbca0d4907",
	CreatedAt:   time.Time{},
	UpdatedAt:   time.Time{},
	UserId:      "a3216484-4ddf-4966-93f5-85010fe88926",
	WorkerId:    "30eff849-07e2-48ce-ae37-fcb41f57f14a",
	FaultType:   "leak",
	Address:     "some address",
	FullName:    "Test User",
	Phone:       "36300000111",
	GeoLocation: "1.1,-1.1",
	ImageUrl:    "http://image.com/1.jpg",
	Status:      DRAFT,
}

var ticket2 = Ticket{
	Id:          "2a5e9658-f954-45c0-a232-4dcbca0d4904",
	CreatedAt:   time.Time{},
	UpdatedAt:   time.Time{},
	UserId:      "6651f933-7343-401b-8de2-2014ad58ab55",
	WorkerId:    "30eff849-07e2-48ce-ae37-fcb41f57f14a",
	FaultType:   "leak",
	Address:     "some address 2",
	FullName:    "Test User 2",
	Phone:       "36300000112",
	GeoLocation: "1.1,-1.1",
	ImageUrl:    "http://image.com/2.jpg",
	Status:      DRAFT,
}

var ticket3 = Ticket{
	Id:          "dc80936e-d015-4270-8df3-c9d24e651cf7",
	CreatedAt:   time.Time{},
	UpdatedAt:   time.Time{},
	UserId:      "93de681e-38c0-4b9b-bc27-8839e27aa19b",
	WorkerId:    "30eff849-07e2-48ce-ae37-fcb41f57f14a",
	FaultType:   "leak",
	Address:     "some address 3",
	FullName:    "Test User 3",
	Phone:       "36300000113",
	GeoLocation: "1.1,-1.1",
	ImageUrl:    "http://image.com/3.jpg",
	Status:      DRAFT,
}

var ticket4 = Ticket{
	Id:          "33c0df42-4372-4958-8cf0-da68b58eac46",
	CreatedAt:   time.Time{},
	UpdatedAt:   time.Time{},
	UserId:      "149e4ef6-f04b-4890-925b-9434f09c596e",
	WorkerId:    "NULL",
	FaultType:   "leak",
	Address:     "some address 4",
	FullName:    "Test User 4",
	Phone:       "36300000114",
	GeoLocation: "1.1,-1.1",
	ImageUrl:    "http://image.com/4.jpg",
	Status:      DRAFT,
}

var ticket5 = Ticket{
	Id:          "cc80936e-d015-4270-8df3-c9d24e651cf7",
	CreatedAt:   time.Time{},
	UpdatedAt:   time.Time{},
	UserId:      "93de681e-38c0-4b9b-bc27-8839e27aa19b",
	WorkerId:    "30eff849-07e2-48ce-ae37-fcb41f57f14a",
	FaultType:   "leak",
	Address:     "some address 5",
	FullName:    "Test User 5",
	Phone:       "36300000113",
	GeoLocation: "1.1,-1.1",
	ImageUrl:    "http://image.com/5.jpg",
	Status:      DRAFT,
}


var tickets = []Ticket{
	ticket,
	ticket2,
	ticket3,
	ticket4,
	ticket5,
}


type newRepo struct {}

func (r* newRepo) Find(column string, value string) (*Ticket, error) {

	//hack for tests
	firstUpper := strings.Title(column)

	for i := range tickets {
		user := reflect.ValueOf(&tickets[i])

		colVal := reflect.Indirect(user).FieldByName(firstUpper)
		if colVal.String() == value {
			return &tickets[i], nil
		}
	}
	return nil, ErrTicketNotFound
}

func (r* newRepo) FindAll(page int, rowsPerPage int, sortBy string, descending bool, filter string, workerId string, status string, lat string, long string) (*[]Ticket, int, error){
	offsetPage := page - 1


	var filteredTickets []Ticket

	for _, ticket := range tickets {
		workerRestriction := false
		statusRestriction := false
		latLongRestriction := false

		if workerId != "" && workerId != ticket.WorkerId {
			workerRestriction = true
		}


		if status != "" && StatusType(status) != ticket.Status {
			statusRestriction = true
		}

		if lat != "" && long != "" {
				// just return everything for tests
		}

		if workerRestriction == false && statusRestriction == false && latLongRestriction == false {
			filteredTickets = append(filteredTickets,ticket)
		}
 	}

	if rowsPerPage * page > len(filteredTickets) -1 {
		return &filteredTickets, len(filteredTickets), nil
	}

	pagination := filteredTickets[offsetPage * rowsPerPage: page * rowsPerPage]
	return &pagination, len(filteredTickets), nil

}

func (r* newRepo) Update(ticket *Ticket) error {
	for i := range tickets {
		if tickets[i].Id == ticket.Id {
			tickets[i] = *ticket
			return nil
		}
	}
	return ErrTicketNotFound
}

func (r* newRepo) Delete(id string) error{
	for i := range tickets {
		if tickets[i].Id == id {
			tickets[i] = tickets[len(tickets)-1]
			tickets[len(tickets)-1] = Ticket{}
			tickets = tickets[:len(tickets)-1]
			return nil
		}
	}
	return ErrTicketNotFound
}

func (r* newRepo) Store(ticket *Ticket) (*Ticket, error) {
	tickets = append(tickets, *ticket)
	return ticket, nil
}


func TestRepository(t *testing.T) {

	r := &newRepo{}
	i := reflect.TypeOf((*Repository)(nil)).Elem()
	isInterface := reflect.TypeOf(r).Implements(i)

	if isInterface {
		t.Logf("test if implements repository success, expected %v, got %v", true, isInterface)
	} else {
		t.Errorf("test if implements repository failed, expected %v, got %v", true, isInterface)
	}
}

func TestNewTicketService(t *testing.T) {

	r := &newRepo{}
	service := NewTicketService(r)

	i := reflect.TypeOf((*Service)(nil)).Elem()
	isInterface := reflect.TypeOf(service).Implements(i)

	if isInterface {
		t.Logf("test if implements service success, expected %v, got %v", true, isInterface)
	} else {
		t.Errorf("test if implements service  failed, expected %v, got %v", true, isInterface)
	}
}


func TestSyncTicketWorker(t *testing.T){
	r := &newRepo{}
	service := NewTicketService(r)


	model := SyncTicketRequestModelWorker{
		RequesterId: "6c596632-ed7a-47b3-8349-de42acb8425a",
		Lat:         "1.1",
		Long:        "-1.1",
		Rows: []SyncTicket{
			{
				Id:       "8a5e9658-f954-45c0-a232-4dcbca0d4907",
				ImageUrl: "asd",
				Status:   "done",
			},
		},
	}

	requestInvalid, err := service.SyncTicketWorker(&model)
	if err != nil {
		t.Errorf("test SyncTicketWorker failed, expected %v, got %v", "SyncTicketResponseModelWorker", err)
	}

	if requestInvalid.Rows[0].Id == "33c0df42-4372-4958-8cf0-da68b58eac46" {
		t.Logf("test SyncTicketWorker success, expected %v, got %v", "33c0df42-4372-4958-8cf0-da68b58eac46", requestInvalid.Rows[0].Id)
	} else {
		t.Errorf("test SyncTicketWorker failed, expected %v, got %v", "33c0df42-4372-4958-8cf0-da68b58eac46", requestInvalid.Rows[0].Id)
	}

	if tickets[0].Status == "draft" {
		t.Logf("test SyncTicketWorker success, expected %v to not update", "8a5e9658-f954-45c0-a232-4dcbca0d4907")
	} else {
		t.Errorf("test SyncTicketWorker failed, expected %v to not update", "8a5e9658-f954-45c0-a232-4dcbca0d4907")

	}


	model = SyncTicketRequestModelWorker{
		RequesterId: "30eff849-07e2-48ce-ae37-fcb41f57f14a",
		Lat:         "1.1",
		Long:        "-1.1",
		Rows:        []SyncTicket{
			{
				Id:       "8a5e9658-f954-45c0-a232-4dcbca0d4907",
				ImageUrl: "new_image",
				Status:   "done",
			},
			{
				Id:       "2a5e9658-f954-45c0-a232-4dcbca0d4904",
				ImageUrl: "new_image",
				Status:   "done",
			},
			{
				Id:       "dc80936e-d015-4270-8df3-c9d24e651cf7",
				ImageUrl: "new_image",
				Status:   "done",
			},
		},
	}

	requestValid, err := service.SyncTicketWorker(&model)
	if err != nil {
		t.Errorf("test SyncTicketWorker failed, expected %v, got %v", "SyncTicketResponseModelWorker", err)
	}

	if requestValid.Rows[0].Id == "cc80936e-d015-4270-8df3-c9d24e651cf7" {
		t.Logf("test SyncTicketWorker success, expected %v, got %v", "cc80936e-d015-4270-8df3-c9d24e651cf7", requestValid.Rows[0].Id)
	} else {
		t.Errorf("test SyncTicketWorker failed, expected %v, got %v", "cc80936e-d015-4270-8df3-c9d24e651cf7", requestValid.Rows[0].Id)
	}

	if tickets[0].Status == "done" && tickets[1].Status == "done" && tickets[2].Status == "done" {
		t.Logf("test SyncTicketWorker success, expected all tickets to update when request is correct")
	} else {
		t.Errorf("test SyncTicketWorker failed, expected all tickets to update when request is correct")
	}
}

func TestGetTicketAdmin(t *testing.T){
	r := &newRepo{}
	service := NewTicketService(r)

	model := GetTicketRequestModelAdmin{Id: "8a5e9658-f954-45c0-a232-4dcbca0d4907"}

	shouldFind, err := service.GetTicketAdmin(&model)

	if err != nil {
		t.Errorf("test GetTicketAdmin failed, expected %v, got %v", ticket, err)
	}

	if shouldFind.Id != ticket.Id {
		t.Errorf("test GetTicketAdmin failed, expected %v, got %v", ticket, shouldFind)
	} else {
		t.Logf("test GetTicketAdmin success, expected %v, got %v", ticket, shouldFind)
	}

	wrongModel := GetTicketRequestModelAdmin{Id: "abc" }

	_, err = service.GetTicketAdmin(&wrongModel)

	if err != nil {
		if errors.Cause(err) != ErrRequestInvalid {
			t.Errorf("test not found failed, expected %v, got %v", ErrRequestInvalid, err)
		} else {
			t.Logf("test not found success, expected %v, got %v", ErrRequestInvalid, err)
		}
	}
}

func TestGetTicketWorker(t *testing.T){
	r := &newRepo{}
	service := NewTicketService(r)

	model := GetTicketRequestModelWorker{Id: "cc80936e-d015-4270-8df3-c9d24e651cf7", RequesterId: "30eff849-07e2-48ce-ae37-fcb41f57f14a" }

	shouldFind, err := service.GetTicketWorker(&model)

	if err != nil {
		t.Errorf("test GetTicketWorker failed, expected %v, got %v", ticket5, err)
	}

	if shouldFind.Id != ticket5.Id {
		t.Errorf("test GetTicketWorker failed, expected %v, got %v", ticket5, shouldFind)
	} else {
		t.Logf("test GetTicketWorker success, expected %v, got %v", ticket5, shouldFind)
	}

	wrongModel := GetTicketRequestModelWorker{Id: "abc", RequesterId: "30eff849-07e2-48ce-ae37-fcb41f57f14a"  }

	_, err = service.GetTicketWorker(&wrongModel)

	if err != nil {
		if errors.Cause(err) != ErrRequestInvalid {
			t.Errorf("test not found failed, expected %v, got %v", ErrRequestInvalid, err)
		} else {
			t.Logf("test not found success, expected %v, got %v", ErrRequestInvalid, err)
		}
	}
}


func TestCreateTicketUser(t *testing.T) {
	r := &newRepo{}
	service := NewTicketService(r)

	model := CreateTicketRequestModelUser{
		UserId:    "77f4756f-7fe4-4c5f-a7b9-6c8e09627d0a",
		FaultType: "leak",
		Address:   "something",
		FullName:  "Peter",
		Phone:     "36308889999",
		Lat:       "1.1",
		Long:      "1.2",
	}

	invalidModel := CreateTicketRequestModelUser{
		UserId:    "77f4756f-7fe4-4c5f-a7b9-6c8e09627d0a",
		FaultType: "explosion",
		Address:   "wrong",
	}

	err := service.CreateTicketUser(&model)

	if err != nil{
		t.Errorf("test CreateTicketUser failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test CreateTicketUser success, expected %v, got %v", nil, err)
	}

	err = service.CreateTicketUser(&invalidModel)

	if err != nil{
		t.Logf("test CreateTicketUser invalid success, expected %v, got %v", ErrRequestInvalid, err)
	} else {
		t.Errorf("test CreateTicketUser invalid failed, expected %v, got %v", ErrRequestInvalid, err)
	}
}

func TestGetAllTicketAdmin(t *testing.T){
	r := &newRepo{}
	service := NewTicketService(r)

	model := GetAllTicketRequestModelAdmin{
		Page:        1,
		RowsPerPage: 10,
		SortBy:      "Id",
		Descending:  false,
		Filter:      "",
	}

	shouldFind, err := service.GetAllTicketAdmin(&model)

	if err != nil {
		t.Errorf("test GetAllTicketAdmin failed, expected %v, got %v", "GetAllTicketResponseModel", err)
	}

	if shouldFind.Rows[0].Id != "8a5e9658-f954-45c0-a232-4dcbca0d4907" {
		t.Errorf("test GetAllTicketAdmin failed, expected %v, got %v", "8a5e9658-f954-45c0-a232-4dcbca0d4907", shouldFind.Rows[0].Id)
	} else {
		t.Logf("test GetAllTicketAdmin success, expected %v, got %v", "8a5e9658-f954-45c0-a232-4dcbca0d4907",  shouldFind.Rows[0].Id)
	}

	if shouldFind.Count > 0 {
		t.Logf("test GetAllTicketAdmin success, expected %v, got %v", "greater than zero", shouldFind.Count)
	} else {
		t.Errorf("test GetAllTicketAdmin failed, expected %v, got %v", "greater than zero", shouldFind.Count)
	}


	wrongModel := GetAllTicketRequestModelAdmin{
		RowsPerPage: 0,
		SortBy:      "",
		Descending:  false,
		Filter:      "",
	}

	_, err = service.GetAllTicketAdmin(&wrongModel)

	if err != nil {
		if errors.Cause(err) != ErrRequestInvalid {
			t.Errorf("test GetAllTicketAdmin failed, expected %v, got %v", ErrRequestInvalid, err)
		} else {
			t.Logf("test GetAllTicketAdmin success, expected %v, got %v", ErrRequestInvalid, err)
		}
	}
}

func TestGetAllTicketWorker(t *testing.T){
	r := &newRepo{}
	service := NewTicketService(r)

	model := GetAllTicketRequestModelWorker{
		Page:        1,
		RowsPerPage: 10,
		SortBy:      "Id",
		Descending:  false,
		Filter:      "",
		RequesterId: "30eff849-07e2-48ce-ae37-fcb41f57f14a",
	}

	shouldFind, err := service.GetAllTicketWorker(&model)

	if err != nil {
		t.Errorf("test GetAllTicketWorker failed, expected %v, got %v", "GetAllTicketResponseModel", err)
	}

	if shouldFind.Rows[0].Id != "cc80936e-d015-4270-8df3-c9d24e651cf7" {
		t.Errorf("test GetAllTicketWorker failed, expected %v, got %v", "8a5e9658-f954-45c0-a232-4dcbca0d4907", shouldFind.Rows[0].Id)
	} else {
		t.Logf("test GetAllTicketWorker success, expected %v, got %v", "8a5e9658-f954-45c0-a232-4dcbca0d4907",  shouldFind.Rows[0].Id)
	}

	if shouldFind.Count > 0 {
		t.Logf("test GetAllTicketWorker success, expected %v, got %v", "greater than zero", shouldFind.Count)
	} else {
		t.Errorf("test GetAllTicketWorker failed, expected %v, got %v", "greater than zero", shouldFind.Count)
	}


	wrongModel := GetAllTicketRequestModelWorker{
		RowsPerPage: 0,
		SortBy:      "",
		Descending:  false,
		Filter:      "",
	}

	_, err = service.GetAllTicketWorker(&wrongModel)

	if err != nil {
		if errors.Cause(err) != ErrRequestInvalid {
			t.Errorf("test GetAllTicketWorker failed, expected %v, got %v", ErrRequestInvalid, err)
		} else {
			t.Logf("test GetAllTicketWorker success, expected %v, got %v", ErrRequestInvalid, err)
		}
	}
}


func TestUpdateTicketAdmin(t *testing.T) {
	r := &newRepo{}
	service := NewTicketService(r)


	model := UpdateTicketRequestModelAdmin{
		Id:        "8a5e9658-f954-45c0-a232-4dcbca0d4907",
		FullName:  "updated",
	}

	wrongModel := UpdateTicketRequestModelAdmin{
		FullName: "hallo",
	}

	err := service.UpdateTicketAdmin(&model)

	if err != nil{
		t.Errorf("test UpdateTicketAdmin failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test UpdateTicketAdmin success, expected %v, got %v", nil, err)
	}

	gTicketReqModel := GetTicketRequestModelAdmin{
		Id: "8a5e9658-f954-45c0-a232-4dcbca0d4907",
	}

	shouldFind, err := service.GetTicketAdmin(&gTicketReqModel)

	if err != nil {
		t.Errorf("test UpdateTicketAdmin failed, expected %v after update, got %v",shouldFind, err)
	} else {
		if shouldFind.FullName != "updated" {
			t.Errorf("test UpdateTicketAdmin failed, expected FullName = %v after update, got %v",model.FullName, shouldFind.FullName)
		}
		t.Logf("test UpdateTicketAdmin success, all values are as expected")
	}

	err = service.UpdateTicketAdmin(&wrongModel)

	if err != nil{
		t.Logf("test UpdateTicketAdmin invalid success, expected %v, got %v", ErrTicketInvalid, err)
	} else {
		t.Errorf("test UpdateTicketAdmin invalid failed, expected %v, got %v", ErrTicketInvalid, err)
	}
}

func TestDeleteTicketAdmin(t *testing.T) {
	r := &newRepo{}
	service := NewTicketService(r)

	model := DeleteTicketRequestModelAdmin{Id: "8a5e9658-f954-45c0-a232-4dcbca0d4907"}

	err := service.DeleteTicketAdmin(&model)

	if err != nil {
		t.Errorf("test DeleteTicketAdmin failed, expected %v, got %v", nil, err)
	} else {
		t.Logf("test DeleteTicketAdmin success, expected %v, got %v", nil, err)
	}

	wrongModel := DeleteTicketRequestModelAdmin{Id: "abc"}

	err = service.DeleteTicketAdmin(&wrongModel)

	if err != nil {
		if errors.Cause(err) != ErrRequestInvalid {
			t.Errorf("test DeleteTicketAdmin failed, expected %v, got %v", ErrRequestInvalid, err)
		} else {
			t.Logf("test DeleteTicketAdmin success, expected %v, got %v", ErrRequestInvalid, err)
		}
	}

}


