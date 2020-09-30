package routes

import (
	"context"
	"github.com/ewol123/ticketer-server/ticketer-service/hack"
	"github.com/ewol123/ticketer-server/ticketer-service/ticket"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateTicketAdminBadRequest(t *testing.T) {
	hack.Init("../../hack/seed_test.sql")
	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UpdateTicketAdmin)

	req, err := http.NewRequest("PATCH", "/admin/v1/ticket",
		strings.NewReader(`{"x":"something.com"}`))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("UpdateTicketAdmin test failed, error: %v", err)
	}

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "-ec6c-45e1-869-27e14c3ed92d")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestUpdateTicketAdminNotFound(t *testing.T) {
	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UpdateTicketAdmin)
	req, err := http.NewRequest("PATCH", "/admin/v1/ticket",
		strings.NewReader(`{"Address":"test address 2"}`))
	if err != nil {
		t.Errorf("UpdateTicketAdmin test failed, error: %v", err)
	}

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "8a3f6121-054a-46a7-8a14-18c995a16204")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestUpdateTicketAdmin(t *testing.T) {

	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.UpdateTicketAdmin)
	req, err := http.NewRequest("PATCH", "/admin/v1/ticket",
		strings.NewReader(`{"Address":"test 2","FullName": "test new", "Status": "done"}`))
	if err != nil {
		t.Errorf("UpdateTicketAdmin test failed, error: %v", err)
	}

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "fc21847d-2cec-403d-9ba4-c64e8756a400")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	handler.ServeHTTP(rr, req)

	// check status
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusNoContent)
	}

	hack.TearDown()
}
