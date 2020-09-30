package routes

import (
	"encoding/json"
	"github.com/ewol123/ticketer-server/ticketer-service/hack"
	"github.com/ewol123/ticketer-server/ticketer-service/ticket"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllTicketWorkerBadRequest(t *testing.T) {
	hack.Init("../../hack/seed_test.sql")

	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetAllTicketWorker)
	req, err := http.NewRequest("GET", "/worker/v1/ticket", nil)
	if err != nil {
		t.Errorf("GetAllTicketWorker test failed, error: %v", err)
	}

	q := req.URL.Query()
	q.Add("page", "1")
	q.Add("rowsPerPage", "10")
	q.Add("descending", "false")
	req.URL.RawQuery = q.Encode()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetAllTicketWorkerFound(t *testing.T) {

	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetAllTicketWorker)
	req, err := http.NewRequest("GET", "/worker/v1/ticket", nil)
	if err != nil {
		t.Errorf("GetAllTicketWorker test failed, error: %v", err)
	}

	req.Header.Add("Requester-Id", "60dd5185-6003-48da-9ff1-998a4477529c")

	q := req.URL.Query()
	q.Add("page", "1")
	q.Add("rowsPerPage", "10")
	q.Add("sortBy", "ticket_id")
	q.Add("descending", "false")
	q.Add("filter", "")
	req.URL.RawQuery = q.Encode()
	handler.ServeHTTP(rr, req)

	// check status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusOK)
	}

	// check response body
	m := make(map[string]interface{})
	if err := json.Unmarshal(rr.Body.Bytes(), &m); err != nil {
		t.Errorf("cannot decode response body of GetAllTicketWorker %v", err)
	}

	if m["Count"] != 0 {
		t.Logf("respone body looks OK %v", m)
	} else {
		t.Errorf("response body incorrect")
	}

	hack.TearDown()

}

