package routes

import (
	"context"
	"encoding/json"
	"github.com/ewol123/ticketer-server/ticketer-service/hack"
	"github.com/ewol123/ticketer-server/ticketer-service/ticket"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTicketWorkerBadRequest(t *testing.T) {
	hack.Init("../../hack/seed_test.sql")

	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetTicketWorker)
	req, err := http.NewRequest("GET", "/worker/v1/ticket/121214asd", nil)
	if err != nil {
		t.Errorf("GetTicketWorker test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetTicketWorkerNotFound(t *testing.T) {
	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetTicketWorker)
	req, err := http.NewRequest("GET", "/worker/v1/ticket/", nil)
	if err != nil {
		t.Errorf("GetTicketWorker test failed, error: %v", err)
	}

	req.Header.Add("Requester-Id", "60dd5185-6003-48da-9ff1-998a4477529c")

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "c655b6b9-4242-4ee9-910a-2560e8e49d6e")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestGetTicketWorkerFound(t *testing.T) {

	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.GetTicketWorker)
	req, err := http.NewRequest("GET", "/worker/v1/ticket/", nil)
	if err != nil {
		t.Errorf("GetTicketWorker test failed, error: %v", err)
	}

	req.Header.Add("Requester-Id", "60dd5185-6003-48da-9ff1-998a4477529c")

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "f735dd08-bbdd-4a65-9336-df21804eb47e")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
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
		t.Errorf("cannot decode response body of GetTicketWorker %v", err)
	}

	isEqual := m["Address"] == "test address" && m["FullName"] == "Peter"

	if isEqual != true {
		t.Errorf("handler returned wrong response body")
	} else {
		t.Logf("handler returned correct response body")
	}
	hack.TearDown()

}


