package routes

import (
	"github.com/ewol123/ticketer-server/ticketer-service/hack"
	"github.com/ewol123/ticketer-server/ticketer-service/ticket"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSyncTicketWorkerBadRequest(t *testing.T) {
	hack.Init("../../hack/seed_test.sql")
	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.SyncTicketWorker)

	req, err := http.NewRequest("POST", "/worker/v1/ticket/sync",
		strings.NewReader(`{"Something":"something"}`))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("SyncTicketWorker test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestSyncTicketWorkerUser(t *testing.T) {

	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.SyncTicketWorker)

	req, err := http.NewRequest("GET", "/worker/v1/ticket/sync",
		strings.NewReader(`
		{
		"Lat":"1.1", 
		"Long": "-1.1", 
		"Rows": [{"Id": "f735dd08-bbdd-4a65-9336-df21804eb47e", "ImageUrl": "http://image.com/1.jpg", "Status": "done" }]
		}`))

	if err != nil {
		t.Errorf("SyncTicketWorker test failed, error: %v", err)
	}

	req.Header.Add("Requester-Id", "60dd5185-6003-48da-9ff1-998a4477529c")

	handler.ServeHTTP(rr, req)

	// check status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusOK)
	}

	hack.TearDown()
}
