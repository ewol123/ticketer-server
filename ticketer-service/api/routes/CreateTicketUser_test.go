package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ewol123/ticketer-server/ticketer-service/hack"
	"github.com/ewol123/ticketer-server/ticketer-service/ticket"
)

func TestCreateTicketUserBadRequest(t *testing.T) {
	hack.Init("../../hack/seed_test.sql")
	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.CreateTicketUser)

	req, err := http.NewRequest("POST", "/user/v1/ticket",
		strings.NewReader(`{"Something":"something"}`))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("CreateTicketUser test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestCreateTicketUser(t *testing.T) {

	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.CreateTicketUser)

	req, err := http.NewRequest("POST", "/user/v1/ticket",
		strings.NewReader(`
		{, 
		"FaultType": "leak", 
		"Address": "test", 
		"FullName": "Peter", 
		"Phone": "36300001111", 
		"Lat": "1.1", 
		"Long": "-1.1"}`))

	if err != nil {
		t.Errorf("CreateTicketUser test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	// check status
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusCreated)
	}

	hack.TearDown()
}
