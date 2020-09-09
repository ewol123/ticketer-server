package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	seed "github.com/ewol123/ticketer-server/user-service/hack"
	"github.com/ewol123/ticketer-server/user-service/user"
)

func TestConfirmRegistrationBadRequest(t *testing.T) {
	hack.Init("../../hack/seed_test.sql")
	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ConfirmRegistration)

	req, err := http.NewRequest("POST", "/user/v1/confirm-registration",
		strings.NewReader(`{"RegistrationCode":"123457"}`))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Errorf("ConfirmRegistration test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestConfirmRegistrationNotFound(t *testing.T) {
	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ConfirmRegistration)
	req, err := http.NewRequest("POST", "/user/v1/confirm-registration",
		strings.NewReader(`{"Email":"test2@test.asd","RegistrationCode":"123457"}`))
	if err != nil {
		t.Errorf("ConfirmRegistration test failed, error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	} else {
		t.Logf("handler returned correct status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestConfirmRegistration(t *testing.T) {

	repo := ChooseRepo()
	service := user.NewUserService(repo)
	h := NewHandler(service)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ConfirmRegistration)
	req, err := http.NewRequest("GET", "/user/v1/confirm-registration",
		strings.NewReader(`{"Email":"test2@test.com","RegistrationCode":"123456"}`))
	if err != nil {
		t.Errorf("ConfirmRegistration test failed, error: %v", err)
	}

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
