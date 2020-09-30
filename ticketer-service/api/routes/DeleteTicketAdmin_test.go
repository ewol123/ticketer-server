package routes

import (
	"context"
	"github.com/ewol123/ticketer-server/ticketer-service/hack"
	"github.com/ewol123/ticketer-server/ticketer-service/ticket"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteTicketAdminBadRequest(t *testing.T) {
	hack.Init("../../hack/seed_test.sql")
	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.DeleteTicketAdmin)

	req, err := http.NewRequest("DELETE", "/admin/v1/ticket", nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("DeleteTicketAdmin test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestDeleteTicketAdminNotFound(t *testing.T){
	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.DeleteTicketAdmin)
	req, err := http.NewRequest("DELETE", "/admin/v1/ticket", nil)
	if err != nil {
		t.Errorf("DeleteTicketAdmin test failed, error: %v", err)
	}

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1dbbaa7b-0861-49c2-abde-c31722787166")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestDeleteTicketAdmin(t *testing.T){

	repo := ChooseRepo()
	service := ticket.NewTicketService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.DeleteTicketAdmin)
	req, err := http.NewRequest("DELETE", "/admin/v1/ticket", nil)
	if err != nil {
		t.Errorf("DeleteTicketAdmin test failed, error: %v", err)
	}

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "d076f530-2453-4af2-a9a2-52b54dc3d36f")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
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